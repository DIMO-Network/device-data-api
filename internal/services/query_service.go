package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awstypes "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
)

type QueryStorageService struct {
	es                      *elasticsearch.TypedClient
	storageSvcClient        *s3.Client
	log                     *zerolog.Logger
	AWSBucket               string
	ElasticIndex            string
	NATSDataDownloadSubject string
	MaxFileSize             int
}

type UserData struct {
	UserDeviceID     string                   `json:"userDeviceId"`
	RequestTimestamp string                   `json:"requestTimestamp"`
	Data             []map[string]interface{} `json:"data,omitempty"`
}

func NewQueryStorageService(es *elasticsearch.TypedClient, settings *config.Settings, log *zerolog.Logger) (*QueryStorageService, error) {

	ctx := log.WithContext(context.Background())

	resolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...any) (aws.Endpoint, error) {
			if settings.AWSEndpoint != "" {
				return aws.Endpoint{URL: settings.AWSEndpoint}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		},
	)

	awsconf, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithEndpointResolverWithOptions(resolver))
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsconf)

	return &QueryStorageService{
		es:                      es,
		storageSvcClient:        s3Client,
		AWSBucket:               settings.AWSBucketName,
		ElasticIndex:            settings.ElasticIndex,
		NATSDataDownloadSubject: settings.NATSDataDownloadSubject,
		MaxFileSize:             settings.MaxFileSize,
		log:                     log}, nil
}

func (uds *QueryStorageService) executeESQuery(q *search.Request) (string, error) {
	res, err := uds.es.Search().
		Index(uds.ElasticIndex).
		Request(q).
		Do(context.Background())
	if err != nil {
		uds.log.Err(err).Msg("Could not query Elasticsearch")
		return "", err
	}
	defer res.Body.Close()

	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		uds.log.Err(err).Msg("Could not parse Elasticsearch response body")
		return "", err
	}
	response := string(responseBytes)

	if res.StatusCode != 200 {
		uds.log.Info().RawJSON("elasticsearchResponseBody", responseBytes).Msg("Error from Elastic.")

		err := fmt.Errorf("invalid status code when querying elastic: %d", res.StatusCode)
		return response, err
	}

	return response, nil
}

func (uds *QueryStorageService) StreamDataToS3(ctx context.Context, userDeviceID string, startDate, endDate time.Time) ([]string, error) {
	respSize := pageSize
	var docCount int
	var fileSize int
	var newFile bool
	var keyName string
	downloadLinks := make([]string, 0)
	parts := make([]awstypes.CompletedPart, 0)

	query := uds.formatUserDataRequest(userDeviceID, startDate, endDate)

	expires := time.Now().Add(24 * time.Hour)
	keyName, docCount = generateKeyName(userDeviceID, docCount, startDate, endDate)
	upload, err := uds.storageSvcClient.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket:  aws.String(uds.AWSBucket),
		Key:     aws.String(keyName),
		Expires: &expires,
	})
	if err != nil {
		return downloadLinks, err
	}

	for respSize == pageSize {
		if newFile {
			final, err := uds.storageSvcClient.CompleteMultipartUpload(ctx,
				&s3.CompleteMultipartUploadInput{
					Bucket:   aws.String(uds.AWSBucket),
					Key:      aws.String(keyName),
					UploadId: upload.UploadId,
					MultipartUpload: &awstypes.CompletedMultipartUpload{
						Parts: parts,
					},
				},
			)
			if err != nil {
				return downloadLinks, err
			}

			downloadLinks = append(downloadLinks, *final.Location)
			keyName, docCount = generateKeyName(userDeviceID, docCount, startDate, endDate)
			upload, err = uds.storageSvcClient.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
				Bucket:  aws.String(uds.AWSBucket),
				Key:     aws.String(keyName),
				Expires: &expires,
			})
			if err != nil {
				return downloadLinks, err
			}

			parts = make([]awstypes.CompletedPart, 0)
			fmt.Println("new file: ", keyName)
			newFile = false
			fileSize = 0
		}

		response, err := uds.executeESQuery(query)
		if err != nil {
			uds.log.Err(err).Msg("user data download: unable to query elasticsearch")
			return downloadLinks, err
		}

		respSize = int(gjson.Get(response, "hits.hits.#").Int())
		if respSize == 0 {
			break
		}

		data := make([]map[string]interface{}, respSize)
		err = json.Unmarshal([]byte(gjson.Get(response, "hits.hits.#._source").Raw), &data)
		if err != nil {
			uds.log.Err(err).Msg("user data download: unable to unmarshal data")
			return downloadLinks, err
		}

		partNum := int32(len(parts) + 1)
		fmt.Println(partNum)
		dataString, err := uds.trimJSON(data)
		if err != nil {
			return downloadLinks, err
		}

		if partNum == 1 {
			opening := fmt.Sprintf(`{"userDeviceId": "%s","requestTimestamp": "%s", "data": [`, userDeviceID, time.Now().Format(time.RFC3339))
			dataString = opening + dataString
		}

		if respSize != pageSize {
			dataString = dataString + "]}"
		} else {
			dataString = dataString + ","
		}

		fileSize += len([]byte(dataString))
		if fileSize > uds.MaxFileSize {
			fmt.Println("should make another file")
			dataString = strings.Trim(dataString, ",")

			if !strings.HasSuffix(dataString, "]}") {
				dataString += "]}"
			}
			newFile = true
		}

		reader := bytes.NewReader([]byte(dataString))
		part, err := uds.storageSvcClient.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(uds.AWSBucket),
			Key:        aws.String(keyName),
			UploadId:   upload.UploadId,
			PartNumber: partNum,
			Body:       reader,
		})
		if err != nil {
			return downloadLinks, err
		}

		parts = append(parts, awstypes.CompletedPart{
			PartNumber: partNum,
			ETag:       part.ETag,
		})

		sA := gjson.Get(response, fmt.Sprintf("hits.hits.%d.sort.0", respSize-1))
		query.SearchAfter = []types.FieldValue{sA.String()}

	}

	final, err := uds.storageSvcClient.CompleteMultipartUpload(ctx,
		&s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(uds.AWSBucket),
			Key:      aws.String(keyName),
			UploadId: upload.UploadId,
			MultipartUpload: &awstypes.CompletedMultipartUpload{
				Parts: parts,
			},
		},
	)
	if err != nil {
		return downloadLinks, err
	}

	downloadLinks = append(downloadLinks, *final.Location)
	return downloadLinks, nil
}

// Elastic maximum.
var pageSize = 10000

func (uds *QueryStorageService) formatUserDataRequest(userDeviceID string, startDate, endDate time.Time) *search.Request {
	query := &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Filter: []types.Query{
					{Match: map[string]types.MatchQuery{"subject": {Query: userDeviceID}}},
					{Range: map[string]types.RangeQuery{"data.timestamp": types.DateRangeQuery{
						Gte: timeToEndpoint(startDate),
						Lte: timeToEndpoint(endDate),
					}}},
				},
			},
		},
		Sort: []types.SortCombinations{"data.timestamp"}, // Default is ascending.
		Size: &pageSize,
	}
	return query
}

func (uds *QueryStorageService) trimJSON(data []map[string]interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	s := strings.Trim(string(b), " ")
	s = strings.TrimLeft(s, "[")
	s = strings.TrimRight(s, "]")

	return s, nil
}

func timeToEndpoint(t time.Time) *string {
	if t.IsZero() {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

func generateKeyName(userDeviceID string, docCount int, startDate, endDate time.Time) (string, int) {
	var start, end string
	docCount = docCount + 1

	if !startDate.IsZero() {
		start = "_" + startDate.Format("2006-01-02")
	}

	if !endDate.IsZero() {
		end = "_" + startDate.Format("2006-01-02")
	}

	return fmt.Sprintf("userDownloads/%+v/%+v_DIMODeviceData_%+v%+v%+v.json", userDeviceID, docCount, userDeviceID, start, end), docCount
}

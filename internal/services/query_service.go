package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func (uds *QueryStorageService) StreamDataToS3(ctx context.Context, userDeviceID, startDate, endDate string) (string, error) {
	query := uds.formatUserDataRequest(userDeviceID, startDate, endDate)
	respSize := pageSize

	expires := time.Now().Add(24 * time.Hour)
	keyName := fmt.Sprintf("userDownloads/%+v/DIMODeviceData_%+v_%+v_%+v.json", userDeviceID, userDeviceID, startDate[:10], endDate[:10])
	upload, err := uds.storageSvcClient.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket:  aws.String(uds.AWSBucket),
		Key:     aws.String(keyName),
		Expires: &expires,
	})
	if err != nil {
		return "", err
	}

	parts := make([]awstypes.CompletedPart, 0)

	for respSize == pageSize {
		response, err := uds.executeESQuery(query)
		if err != nil {
			uds.log.Err(err).Msg("user data download: unable to query elasticsearch")
			return "", err
		}

		respSize = int(gjson.Get(response, "hits.hits.#").Int())
		if respSize == 0 {
			fmt.Println("this one!!!")
			break
		}

		data := make([]map[string]interface{}, respSize)
		err = json.Unmarshal([]byte(gjson.Get(response, "hits.hits.#._source").Raw), &data)
		if err != nil {
			uds.log.Err(err).Msg("user data download: unable to unmarshal data")
			return "", err
		}

		partNum := int32(len(parts) + 1)
		fmt.Println(partNum)
		dataString, err := uds.trimJSON(data)
		if err != nil {
			return "", err
		}

		if partNum == 1 {
			opening := fmt.Sprintf(`{"userDeviceId": "%s","requestTimestamp": "%s", "data":`, userDeviceID, time.Now().Format(time.RFC3339))
			dataString = opening + dataString
		}

		if respSize != pageSize {
			dataString = dataString + "]}"
		} else {
			dataString = dataString + ","
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
			return "", err
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
		return "", err
	}

	return *final.Location, nil
}

// Elastic maximum.
var pageSize = 10000

func (uds *QueryStorageService) formatUserDataRequest(userDeviceID, startDate, endDate string) *search.Request {
	query := &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Filter: []types.Query{
					{Match: map[string]types.MatchQuery{"subject": {Query: userDeviceID}}},
					{Range: map[string]types.RangeQuery{"data.timestamp": types.DateRangeQuery{
						Gte: &startDate,
						Lte: &endDate,
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
	s := string(b)
	return s[:len(s)-2], nil
}

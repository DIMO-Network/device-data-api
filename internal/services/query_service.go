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

func (ud *userData) executeESQuery(q *search.Request) (string, error) {
	res, err := ud.es.Search().
		Index(ud.ElasticIndex).
		Request(q).
		Do(context.Background())
	if err != nil {
		ud.log.Err(err).Msg("Could not query Elasticsearch")
		return "", err
	}
	defer res.Body.Close()

	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		ud.log.Err(err).Msg("Could not parse Elasticsearch response body")
		return "", err
	}
	response := string(responseBytes)

	if res.StatusCode != 200 {
		ud.log.Info().RawJSON("elasticsearchResponseBody", responseBytes).Msg("Error from Elastic.")

		err := fmt.Errorf("invalid status code when querying elastic: %d", res.StatusCode)
		return response, err
	}

	return response, nil
}

type userData struct {
	es               *elasticsearch.TypedClient
	storageSvcClient *s3.Client
	log              *zerolog.Logger
	AWSBucket        string
	ElasticIndex     string
	MaxFileSize      int
	keyName          string
	uploadObj        *s3.CreateMultipartUploadOutput
	fileSize         int
	userDeviceID     string
	uploadParts      []awstypes.CompletedPart
	downloadLinks    []string
	query            *search.Request
	docCount         int
}

func (uds *QueryStorageService) newS3Writer(ctx context.Context, query *search.Request, bucketName, userDeviceID string, startDate, endDate time.Time) (*userData, error) {

	exp := time.Now().Add(24 * time.Hour)

	keyName, docCount := generateKeyName(userDeviceID, 0, startDate, endDate)
	upload, err := uds.storageSvcClient.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket:  aws.String(uds.AWSBucket),
		Key:     aws.String(keyName),
		Expires: &exp,
	})
	if err != nil {
		return nil, err
	}

	return &userData{
		uploadParts:      make([]awstypes.CompletedPart, 0),
		AWSBucket:        bucketName,
		uploadObj:        upload,
		keyName:          keyName,
		docCount:         docCount,
		es:               uds.es,
		storageSvcClient: uds.storageSvcClient,
		log:              uds.log,
		ElasticIndex:     uds.ElasticIndex,
		MaxFileSize:      uds.MaxFileSize,
		userDeviceID:     userDeviceID,
		query:            query,
	}, nil

}

func (ud *userData) partNum() int32 {
	return int32(len(ud.uploadParts) + 1)
}

func (ud *userData) abortUploadHandleError(ctx context.Context, err error) {
	_, s3err := ud.storageSvcClient.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(ud.AWSBucket),
		Key:      aws.String(ud.keyName),
		UploadId: ud.uploadObj.UploadId,
	})
	if s3err != nil {
		ud.log.Err(s3err).Msgf("error aborting multipart upload: %+v", err)
	}
}

func (ud *userData) uploadPartToS3(ctx context.Context, reader *bytes.Reader, uploadParts []awstypes.CompletedPart) ([]awstypes.CompletedPart, error) {

	partNum := ud.partNum()

	part, err := ud.storageSvcClient.UploadPart(ctx, &s3.UploadPartInput{
		Bucket:     aws.String(ud.AWSBucket),
		Key:        aws.String(ud.keyName),
		UploadId:   ud.uploadObj.UploadId,
		PartNumber: partNum,
		Body:       reader,
	})
	if err != nil {
		ud.log.Err(err).Msg("error writing part to s3")
		ud.abortUploadHandleError(ctx, err)
		return uploadParts, err
	}

	uploadParts = append(uploadParts, awstypes.CompletedPart{
		PartNumber: partNum,
		ETag:       part.ETag,
	})

	return uploadParts, nil
}

func (ud *userData) finishWritingToS3(ctx context.Context, uploadParts []awstypes.CompletedPart) {
	final, err := ud.storageSvcClient.CompleteMultipartUpload(ctx,
		&s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(ud.AWSBucket),
			Key:      aws.String(ud.keyName),
			UploadId: ud.uploadObj.UploadId,
			MultipartUpload: &awstypes.CompletedMultipartUpload{
				Parts: uploadParts,
			},
		},
	)
	if err != nil {
		ud.log.Err(err).Msg("error finishing file write to s3")
		ud.abortUploadHandleError(ctx, err)
		return
	}

	ud.downloadLinks = append(ud.downloadLinks, *final.Location)
}

func (ud *userData) writeToS3(ctx context.Context, response string) error {

	respSize := int(gjson.Get(response, "hits.hits.#").Int())
	data := make([]map[string]interface{}, respSize)
	err := json.Unmarshal([]byte(gjson.Get(response, "hits.hits.#._source").Raw), &data)
	if err != nil {
		ud.log.Err(err).Msg("user data download: unable to unmarshal data")
		ud.abortUploadHandleError(ctx, err)
		return err
	}

	dataString, err := trimJSON(data)
	if err != nil {
		ud.log.Err(err).Msg("user data download: error trimming data")
		ud.abortUploadHandleError(ctx, err)
		return err
	}
	ud.fileSize += len([]byte(dataString))

	if ud.partNum() == 1 {
		opening := fmt.Sprintf(`{"userDeviceId": "%s","requestTimestamp": "%s", "data": [`, ud.userDeviceID, time.Now().Format(time.RFC3339))
		dataString = opening + dataString
	}

	// when we want to set a max file size, also check if ud.fileSize > ud.MaxFileSize here
	if respSize < pageSize {
		dataString = dataString + "]}"
		reader := bytes.NewReader([]byte(dataString))
		uploadParts, err := ud.uploadPartToS3(ctx, reader, ud.uploadParts)
		if err != nil {
			return err
		}
		ud.finishWritingToS3(ctx, uploadParts)
		return nil
	}

	dataString = dataString + ","
	reader := bytes.NewReader([]byte(dataString))
	ud.uploadParts, err = ud.uploadPartToS3(ctx, reader, ud.uploadParts)
	if err != nil {
		return err
	}

	sA := gjson.Get(response, fmt.Sprintf("hits.hits.%d.sort.0", respSize-1))
	ud.query.SearchAfter = []types.FieldValue{sA.String()}

	response, err = ud.executeESQuery(ud.query)
	if err != nil {
		ud.log.Err(err).Msg("user data download: unable to unmarshal data")
		ud.abortUploadHandleError(ctx, err)
		return err
	}

	return ud.writeToS3(ctx, response)

}

func (uds *QueryStorageService) StreamDataToS3(ctx context.Context, userDeviceID string, startDate, endDate time.Time) ([]string, error) {

	query := uds.formatUserDataRequest(userDeviceID, startDate, endDate)
	response, err := uds.executeESQuery(query)
	if err != nil {
		return []string{}, err
	}

	s3writer, err := uds.newS3Writer(ctx, query, uds.AWSBucket, userDeviceID, startDate, endDate)
	if err != nil {
		uds.log.Err(err).Msg("data streaming service: error creating s3 writer object")
		return []string{}, err
	}

	err = s3writer.writeToS3(ctx, response)
	if err != nil {
		uds.log.Err(err).Msg("data streaming service: unable to write data to s3")
		return []string{}, err
	}

	return s3writer.downloadLinks, nil

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

func trimJSON(data []map[string]interface{}) (string, error) {
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

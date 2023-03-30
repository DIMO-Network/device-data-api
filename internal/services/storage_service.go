package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog"
)

type StorageService struct {
	storageSvcClient *s3.Client
	log              *zerolog.Logger
	AWSBucket        string
	awsFileSize      int
}

func NewStorageService(settings *config.Settings, log *zerolog.Logger) (*StorageService, error) {

	ctx := log.WithContext(context.Background())
	awsconf, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...any) (aws.Endpoint, error) {
					return aws.Endpoint{URL: settings.AWSEndpoint}, nil
				},
			),
		),
	)

	if err != nil {
		return nil, err
	}
	s3Client := s3.NewFromConfig(awsconf)

	return &StorageService{
		storageSvcClient: s3Client,
		log:              log,
		AWSBucket:        settings.AWSBucketName,
		awsFileSize:      settings.AWSFileSize}, nil
}

func (ss *StorageService) generatePreSignedURL(ctx context.Context, keyName string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(ss.storageSvcClient)
	presignParams := &s3.GetObjectInput{
		Bucket: aws.String(ss.AWSBucket),
		Key:    aws.String(keyName),
	}
	presignDuration := func(po *s3.PresignOptions) {
		po.Expires = 5 * time.Minute
	}
	presignResult, err := presignClient.PresignGetObject(ctx, presignParams, presignDuration)
	return presignResult.URL, err
}

func (ss *StorageService) putObjectS3(ctx context.Context, keyName string, data []byte) error {

	_, err := ss.storageSvcClient.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(ss.AWSBucket),
		Key:    aws.String(keyName),
		Body:   bytes.NewReader(data),
	})
	return err

}

func (ss *StorageService) UploadUserData(ctx context.Context, params QueryValues, ud UserData) ([]string, error) {

	batchNum := math.Ceil(float64(len(ud.Data)) / float64(ss.awsFileSize))
	count := 1

	generatedURLs := make([]string, 0)

	for startIndex := 0; startIndex < len(ud.Data); startIndex += ss.awsFileSize {
		endIndex := startIndex + ss.awsFileSize
		if endIndex > len(ud.Data) {
			endIndex = len(ud.Data)
		}

		keyName := fmt.Sprintf("userDownloads/%+s/%+s_%+s_%vof%v.json", params.UserDeviceID, params.RangeStart, params.RangeEnd, count, batchNum)

		dataBytes, err := json.Marshal(UserData{
			User:             ud.User,
			RequestTimestamp: ud.RequestTimestamp,
			Data:             ud.Data[startIndex:endIndex],
		})
		if err != nil {
			return []string{}, err
		}

		err = ss.putObjectS3(ctx, keyName, dataBytes)
		if err != nil {
			return []string{}, err
		}

		url, err := ss.generatePreSignedURL(ctx, keyName, presignDuration*time.Hour)

		generatedURLs = append(generatedURLs, url)

	}
	return generatedURLs, nil
}

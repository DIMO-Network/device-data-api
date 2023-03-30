package services

import (
	"bytes"
	"context"
	"encoding/json"
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
		AWSBucket:        settings.AWSBucketName}, nil
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

func (ss *StorageService) UploadUserData(ctx context.Context, ud UserData, keyName string) (string, error) {
	dataBytes, err := json.Marshal(ud)
	if err != nil {
		return "", err
	}

	err = ss.putObjectS3(ctx, keyName, dataBytes)
	if err != nil {
		return "", err
	}
	return ss.generatePreSignedURL(ctx, keyName, presignDuration*time.Hour)
}

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
	settings         *config.Settings
	log              *zerolog.Logger
	context          context.Context
}

func NewStorageService(settings *config.Settings, log *zerolog.Logger) *StorageService {

	ctx := log.WithContext(context.Background())
	awsconf, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load AWS configuration.")
	}
	s3Client := s3.NewFromConfig(awsconf)

	return &StorageService{storageSvcClient: s3Client, log: log, settings: settings, context: ctx}
}

func (ss *StorageService) generatePreSignedURL(keyName string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(ss.storageSvcClient)
	presignParams := &s3.GetObjectInput{
		Bucket: aws.String(ss.settings.AWSBucketName),
		Key:    aws.String(keyName),
	}
	presignDuration := func(po *s3.PresignOptions) {
		po.Expires = 5 * time.Minute
	}
	presignResult, err := presignClient.PresignGetObject(ss.context, presignParams, presignDuration)
	return presignResult.URL, err
}

func (ss *StorageService) putObjectS3(keyName string, data []byte) error {

	_, err := ss.storageSvcClient.PutObject(ss.context, &s3.PutObjectInput{
		Bucket: aws.String(ss.settings.AWSBucketName),
		Key:    aws.String(keyName),
		Body:   bytes.NewReader(data),
	})
	return err

}

func (ss *StorageService) UploadUserData(ud UserData, keyName string) (string, error) {
	dataBytes, err := json.Marshal(ud)
	if err != nil {
		return "", err
	}

	err = ss.putObjectS3(keyName, dataBytes)
	if err != nil {
		return "", err
	}
	return ss.generatePreSignedURL(keyName, presignDuration*time.Hour)
}

package services

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog"
)

type StorageService struct {
	storageSvcSession *session.Session
	AWSBucketName     *string
	log               *zerolog.Logger
	AWSRegion         *string
}

func NewStorageService(settings *config.Settings, log *zerolog.Logger) *StorageService {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(settings.AWSRegion),
		Credentials: credentials.NewStaticCredentials(settings.AWSAccessKeyID, settings.AWSSecretAccessKey, ""),
	})
	if err != nil {
		return &StorageService{}
	}
	return &StorageService{storageSvcSession: sess, log: log, AWSRegion: aws.String(settings.AWSRegion), AWSBucketName: aws.String(settings.AWSBucketName)}
}

func (ss *StorageService) generatePreSignedURL(keyName string, session *s3.S3, expiration time.Duration) (string, error) {
	req, _ := session.GetObjectRequest(&s3.GetObjectInput{
		Bucket: ss.AWSBucketName,
		Key:    aws.String(keyName),
	})
	return req.Presign(expiration)
}

func (ss *StorageService) putObjectS3(keyname string, data []byte, svc *s3.S3) error {
	params := &s3.PutObjectInput{
		Bucket: ss.AWSBucketName,
		Key:    aws.String(keyname),
		Body:   bytes.NewReader(data),
	}
	_, err := svc.PutObject(params)
	return err

}

func (ss *StorageService) UploadUserData(ud UserData, keyName string) (string, error) {
	dataBytes, err := json.Marshal(ud)
	if err != nil {
		return "", err
	}
	svc := s3.New(ss.storageSvcSession)
	err = ss.putObjectS3(keyName, dataBytes, svc)
	if err != nil {
		return "", err
	}
	return ss.generatePreSignedURL(keyName, svc, presignDuration*time.Hour)
}

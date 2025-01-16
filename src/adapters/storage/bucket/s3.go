package bucket

import (
	"example/web-service-gin/src/infra/configuration"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang-migrate/migrate/v4/source/file"
	"log/slog"
	"mime/multipart"
)

// S3Storage implements port.StoragePort
type S3Storage struct {
	sess       *session.Session
	bucketName string
	s3Client   *s3.S3
}

func NewS3Bucket(configs *configuration.Aws) *S3Storage {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(configs.Region),
		Credentials: credentials.NewStaticCredentials(configs.ClientId, configs.Secret, configs.SessionToken),
	})

	if err != nil {
		slog.Info("Error authenticating with AWS")
	}

	return &S3Storage{sess, configs.BucketName, s3.New(sess)}
}

func (handler *S3Storage) UploadFile(file *multipart.FileHeader, fileKey string) (string, error) {

	fileData, _ := file.Open()

	// Upload input parameters
	upParams := &s3.PutObjectInput{
		Bucket:             aws.String(handler.bucketName),
		Key:                aws.String(fileKey),
		Body:               fileData,
		ContentDisposition: aws.String("attachment"),
	}

	// This will start an async running, and it will not wait to end
	go handler.s3Client.PutObject(upParams)

	return fileKey, nil
}

func (handler *S3Storage) DownloadFile(fileKey string) (*file.File, error) {
	return nil, nil
}

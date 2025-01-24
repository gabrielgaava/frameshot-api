package bucket

import (
	"context"
	"example/web-service-gin/src/infra/configuration"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang-migrate/migrate/v4/source/file"
	"mime/multipart"
	"strings"
)

// S3Storage implements port.StoragePort
type S3Storage struct {
	config     *configuration.Aws
	bucketName string
	s3Client   *s3.Client
	ctx        context.Context
}

func NewS3Bucket(configs *configuration.Aws, ctx context.Context) *S3Storage {
	s3Client := s3.NewFromConfig(configs.Config)
	return &S3Storage{
		configs,
		configs.BucketName,
		s3Client,
		ctx}
}

func (handler *S3Storage) UploadFile(file *multipart.FileHeader, fileKey string) (string, error) {

	fileData, _ := file.Open()

	// Upload input parameters
	upParams := &s3.PutObjectInput{
		Bucket: aws.String(handler.bucketName),
		Key:    aws.String(fileKey),
		Body:   fileData,
	}

	// This will start an async running, and it will not wait to end
	go handler.s3Client.PutObject(handler.ctx, upParams)

	return fileKey, nil
}

func (handler *S3Storage) DownloadFile(fileKey string) (*file.File, error) {
	return nil, nil
}

func (handler *S3Storage) GetFileUrl(fileKey string) string {

	template := "https://<bucket-name>.s3.<region>.amazonaws.com/<object-key>"
	template = strings.ReplaceAll(template, "<bucket-name>", handler.config.BucketName)
	template = strings.ReplaceAll(template, "<region>", handler.config.Config.Region)
	template = strings.ReplaceAll(template, "<object-key>", fileKey)

	return template
}

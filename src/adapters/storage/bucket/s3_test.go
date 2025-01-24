package bucket_test

import (
	"context"
	"example/web-service-gin/src/adapters/storage/bucket"
	"example/web-service-gin/src/infra/configuration"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setUp() *bucket.S3Storage {
	ctx := context.Background()
	config := configuration.Aws{
		Config:     aws.Config{Region: "us-east-1"},
		BucketName: "bucket-name",
	}

	return bucket.NewS3Bucket(&config, ctx)
}

func TestNewS3Bucket(t *testing.T) {
	storage := setUp()
	assert.NotNil(t, storage)
}

func TestDownloadFile(t *testing.T) {
	storage := setUp()
	file, err := storage.DownloadFile("/path/to/file")
	assert.Nil(t, file)
	assert.Nil(t, err)
}

func TestGetFileUrl(t *testing.T) {
	storage := setUp()
	expectedUrl := "https://bucket-name.s3.us-east-1.amazonaws.com/file.zip"
	url := storage.GetFileUrl("file.zip")
	assert.NotNil(t, url)
	assert.Equal(t, expectedUrl, url)
}

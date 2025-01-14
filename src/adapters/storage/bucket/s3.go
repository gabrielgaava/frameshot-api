package bucket

import (
	"example/web-service-gin/src/infra/configuration"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Bucket struct {
	sess      *session.Session
	s3Service *s3.S3
}

func NewS3Bucket(configs *configuration.Aws) (*S3Bucket, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(configs.Region),
		Credentials: credentials.NewStaticCredentials(configs.ClientId, configs.Secret, configs.SessionToken),
	})

	s3Service := s3.New(sess)
	return &S3Bucket{sess, s3Service}, nil
}

func UploadFile() (string, error) {
	return "", nil
}

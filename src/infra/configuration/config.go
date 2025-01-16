package configuration

import (
	"context"
	awslib "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joho/godotenv"
	"os"
)

// Container contains environment variables for the application, database, cache, token, and http server
type (
	Container struct {
		App  *App
		DB   *Database
		HTTP *HTTP
		AWS  *Aws
	}
	// App contains all the environment variables for the application
	App struct {
		Name string
		Env  string
	}

	// Database contains all the environment variables for the database
	Database struct {
		Connection string
		Host       string
		Port       string
		User       string
		Password   string
		Name       string
	}
	// HTTP contains all the environment variables for the http server
	HTTP struct {
		Env            string
		URL            string
		Port           string
		AllowedOrigins string
	}

	Aws struct {
		Config             awslib.Config
		BucketName         string
		CognitoJwksUrl     string
		S3QueueUrl         string
		VideoInputQueueUrl string
	}
)

// New creates a new container instance
func New() (*Container, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	app := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}

	db := &Database{
		Connection: os.Getenv("DB_CONNECTION"),
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		Name:       os.Getenv("DB_NAME"),
	}

	http := &HTTP{
		Env:            os.Getenv("APP_ENV"),
		URL:            os.Getenv("HTTP_URL"),
		Port:           os.Getenv("HTTP_PORT"),
		AllowedOrigins: os.Getenv("HTTP_ALLOWED_ORIGINS"),
	}

	awsConfiguration, _ := config.LoadDefaultConfig(context.Background())

	aws := &Aws{
		Config:             awsConfiguration,
		BucketName:         os.Getenv("AWS_BUCKET_NAME"),
		CognitoJwksUrl:     os.Getenv("AWS_COGNITO_JWKS_URL"),
		S3QueueUrl:         os.Getenv("AWS_S3_QUEUE_URL"),
		VideoInputQueueUrl: os.Getenv("AWS_VIDEO_INPUT_QUEUE_URL"),
	}

	return &Container{
		app,
		db,
		http,
		aws,
	}, nil
}

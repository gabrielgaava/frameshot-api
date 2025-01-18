package main

import (
	"context"
	"example/web-service-gin/src/adapters/handler/http"
	"example/web-service-gin/src/adapters/handler/queue"
	"example/web-service-gin/src/adapters/storage/bucket"
	"example/web-service-gin/src/adapters/storage/postgres"
	"example/web-service-gin/src/adapters/storage/postgres/repository"
	"example/web-service-gin/src/core/usecase"
	"example/web-service-gin/src/infra/configuration"
	"example/web-service-gin/src/infra/middleware"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	config := loadEnv()
	slog.Info("Starting the application", "app", config.App.Name, "env", config.App.Env)

	ctx := context.Background()
	db := loadDatabase(ctx, &config)
	defer db.Close()

	// Setting SQS
	queueHandler := queue.NewSQSHandler(config.AWS)
	sqsProducer := queue.NewSQSProducer(queueHandler, config.AWS.VideoInputQueueUrl)

	//Dependency Injection
	s3Storage := bucket.NewS3Bucket(config.AWS, ctx)
	requestRepository := repository.NewPGRequestRepository(db)
	requestUseCase := usecase.NewRequestUseCase(requestRepository, s3Storage, sqsProducer) // COLOCAR AQ O BAGUI
	requestHandler := http.NewRequestHandler(requestUseCase)

	// Starting Queue Consumers
	go queue.StartQueueConsumer(queueHandler, config.AWS.S3QueueUrl, requestUseCase.HandleUploadNotification, ctx)
	go queue.StartQueueConsumer(queueHandler, config.AWS.VideoOutputQueueUrl, requestUseCase.HandleVideoOutputNotification, ctx)

	// Routes and Middlewares Settings
	router := gin.Default()
	router.Use(middleware.JwtServiceMiddleware(config.AWS.CognitoJwksUrl))
	router.MaxMultipartMemory = 8 << 20
	router.POST("/requests", requestHandler.Register)
	router.GET("/requests", requestHandler.ListUsers)

	defer router.Run("localhost:8080")
}

// Load de .env file and create a pointer to all its configuration keys
func loadEnv() configuration.Container {
	config, err := configuration.New()

	if err != nil {
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	return *config
}

// Validate database connection and start the migrations
func loadDatabase(ctx context.Context, config *configuration.Container) *postgres.DB {
	db, err := postgres.New(ctx, config.DB)

	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}

	slog.Info("Successfully connected to the database", "db", config.DB.Connection)

	// Migrate database
	err = db.Migrate()
	if err != nil {
		slog.Error("Error migrating database", "error", err)
		os.Exit(1)
	}

	slog.Info("Successfully migrated the database")
	return db
}

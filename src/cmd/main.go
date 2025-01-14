package main

import (
	"context"
	"example/web-service-gin/src/adapters/handler/http"
	"example/web-service-gin/src/adapters/storage/postgres"
	"example/web-service-gin/src/adapters/storage/postgres/repository"
	"example/web-service-gin/src/core/usecase"
	"example/web-service-gin/src/infra/configuration"
	"github.com/gin-gonic/gin"
	"log/slog"
	HTTP "net/http"
	"os"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(HTTP.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(HTTP.StatusCreated, newAlbum)
}

func main() {

	config, err := configuration.New()

	if err != nil {
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	// Set logger
	// logger.Set(configuration.App)

	slog.Info("Starting the application", "app", config.App.Name, "env", config.App.Env)

	ctx := context.Background()
	db, err := postgres.New(ctx, config.DB)

	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}

	defer db.Close()
	slog.Info("Successfully connected to the database", "db", config.DB.Connection)

	// Migrate database
	err = db.Migrate()
	if err != nil {
		slog.Error("Error migrating database", "error", err)
		os.Exit(1)
	}

	slog.Info("Successfully migrated the database")

	//Dependency Injection
	requestRepository := repository.NewPGRequestRepository(db)
	requestUseCase := usecase.NewRequestUseCase(requestRepository)
	requestHandler := http.NewRequestHandler(requestUseCase)

	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.POST("/requests", requestHandler.Register)
	router.GET("/requests", requestHandler.ListUsers)

	router.Run("localhost:8080")
}

package main

import (
	"context"
	"killrvideo/go-backend-astra-dataapi/controllers"
	repo "killrvideo/go-backend-astra-dataapi/repositories"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	baseCtx := context.Background()

	// define DB connection
	cfg := repo.AstraConfig{
		Token:    os.Getenv("ASTRA_DB_APPLICATION_TOKEN"),
		Keyspace: os.Getenv("ASTRA_DB_KEYSPACE"),
		Endpoint: os.Getenv("ASTRA_DB_API_ENDPOINT"),
	}

	db, err := repo.NewAstraSession(cfg)

	if err != nil {
		log.Fatalf("Failed to connect to Astra: %v", err)
	}

	// controller definitions
	authController := controllers.NewAuthController(db, baseCtx)
	healthController := controllers.NewHealthController()
	//ratingsController := controllers.NewRatingsController(db, baseCtx)
	videoController := controllers.NewVideoController(db, baseCtx)

	// route definitions
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		health := api.Group("/health")
		{
			health.GET("", healthController.GetHealth)
		}
		auth := api.Group("/users")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/register", authController.Register)
			auth.GET("/me", authController.GetCurrentUser)
			auth.GET(":id", authController.GetUser)
		}
		videos := api.Group("/videos")
		{
			videos.GET("/id/:id", videoController.GetVideo)
			videos.POST("", videoController.SubmitVideo)
			videos.GET("/latest", videoController.GetLatestVideos)
			videos.GET("/id/:id/related", videoController.GetSimilarVideos)
			//			videos.GET("/:id/ratings", ratingsController.GetRatingsByVideoId)
			//			videos.POST("/id/:id/view", videoController.RecordVideoView)
			videos.GET("/:id/comments", videoController.GetComments)
			//			videos.POST("/:id/comments", videoController.SubmitComment)
		}
	}

	router.RunTLS("localhost:8443", "localhost.pem", "localhost-key.pem")
}

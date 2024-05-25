package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	cache = Cache{}
)

func main() {
	initialDatabaseConnection()

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8080"
	}

	ttl := os.Getenv("CACHE_TTL")
	if ttl == "" {
		ttl = "1"
	}
	ttlInt, err := strconv.Atoi(ttl)
	if err != nil {
		log.Fatalf("Failed to convert CACHE_TTL to int: %v", err)
	}

	go cache.Updater(time.Second * time.Duration(ttlInt))

	router := gin.New()
	router.Use(gin.Recovery())
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/ads", createAdHandler)
			v1.GET("/ads", getAdsHandler)
		}
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

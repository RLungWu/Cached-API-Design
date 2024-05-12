package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/RLungWu/Dcard-Backend-v2/db"
	"github.com/RLungWu/Dcard-Backend-v2/pkg/ad"
	"github.com/RLungWu/Dcard-Backend-v2/pkg/cache"
	"github.com/RLungWu/Dcard-Backend-v2/router"

	"github.com/gin-gonic/gin"
)

var (
	cacheAd = ad.Cachee{}
)

func main() {
	err := db.InitDatabaseConnection()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("AD_SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	ttl := os.Getenv("AD_CACHE_TTL")
	if ttl == "" {
		ttl = "1"
	}
	ttlInt, err := strconv.Atoi(ttl)
	if err != nil {
		panic(err)
	}

	go cacheAd.Updater(time.Second * time.Duration(ttlInt))

	r := gin.New()
	r.Use(gin.Recovery())
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/ad", router.ListAdsHandler)
			v1.POST("/ad", router.CreateAdHandler)
		}
	}

	log.Printf("Server is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

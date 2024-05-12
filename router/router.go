package router

import (
	"context"
	"net/http"

	"github.com/RLungWu/Dcard-Backend-v2/pkg/ad"
	"github.com/RLungWu/Dcard-Backend-v2/pkg/cache"

	"github.com/gin-gonic/gin"
)

func CreateAdHandler(c *gin.Context) {
	var ad ad.Ad
	if err := c.ShouldBindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := ads.InsertOne(context.TODO(), ad)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

func ListAdsHandler(c *gin.Context) {
	var query ad.AdQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results := cache.Filter(query)
	c.JSON(http.StatusOK, results)
}

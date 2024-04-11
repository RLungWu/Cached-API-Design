package main

import(
	"github.com/RLungWu/Dcard-Backend-HW/internal/api"
	
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	
	router.POST("/api/v1/ad", api.AdminCreateAD)
	router.GET("/api/v1/ad", api.GetAD)

	router.Run(":8080")
}
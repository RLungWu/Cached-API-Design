package main

import (
	"github.com/RLungWu/Dcard-Backend-HW/api/admin"
	"github.com/RLungWu/Dcard-Backend-HW/api/public"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/ad", admin.AdminCreateAD)
		v1.GET("/ad", public.GetAD)
	}

	router.Run(":8080")
}

package main

import(
	"github.com/RLungWu/Dcard-Backend-HW/internal/user"
	"github.com/RLungWu/Dcard-Backend-HW/internal/admin"
	
	"github.com/gin-gonic/gin"
)




func main() {
	router := gin.Default()
	

	router.POST("/api/v1/ad", AdminCreateAD)
	router.GET("/api/v1/ad", GetAD)


	router.Run(":8080")
}
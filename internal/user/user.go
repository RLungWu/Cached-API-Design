package user

import "github.com/gin-gonic/gin"

func GetAD(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Get AD",
	})
}
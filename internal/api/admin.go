package api

import "github.com/gin-gonic/gin"

func AdminCreateAD(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Admin Create AD",
	})
}
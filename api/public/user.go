package public

import (
	"github.com/RLungWu/Dcard-Backend-HW/api/ad"
	"github.com/gin-gonic/gin"
)

func GetAD(c *gin.Context) {
	//Get query from route
	var ad ad.AdRequest

	c.Bind(&ad)

	//Check request
	
}
package public

import (
	ad "github.com/RLungWu/Dcard-Backend-HW/api/Ad"
	"github.com/gin-gonic/gin"
)

func GetAD(c *gin.Context) {
	//Get query from route
	var ad ad.AdRequest

	c.Bind(&ad)

	//Check request
	
}
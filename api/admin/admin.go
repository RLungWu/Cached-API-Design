package admin

import (
	"net/http"

	"github.com/RLungWu/Dcard-Backend-HW/db/postgre"
	"github.com/RLungWu/Dcard-Backend-HW/db/redis"
	"github.com/gin-gonic/gin"
)

func AdminCreateAD(c *gin.Context) {
	var ad AdRequest
	if err := c.ShouldBindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := checkAdRequest(&ad); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	//Put in Redis
	db := redis.CreateRedisClient(0)
	defer db.Close()
	db.Set(redis.Context, ad.Title, ad, 0)

	//Put in Postgres
	pgdb := postgre.CreatePGClient()
	defer pgdb.Close()
	_, err := pgdb.Exec("INSERT INTO ad (title, start_at, end_at, conditions) VALUES ($1, $2, $3, $4)", ad.Title, ad.StartAt, ad.EndAt, ad.Conditions)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Insert into Postgres failed :" + err.Error(),
		})
	}



	c.JSON(200, gin.H{
		"message": "Admin Create AD",
		"ad":      ad,
	})
}

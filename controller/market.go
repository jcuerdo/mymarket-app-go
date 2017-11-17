package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jcuerdo/mymarket-app-go/database"
	"io/ioutil"
	"github.com/jcuerdo/mymarket-app-go/model"
	"encoding/json"
)

func GetMarkets() gin.HandlerFunc {
	return func(c *gin.Context) {
		marketRepository := database.GetMarketRepository()
		markets := marketRepository.GetMarkets()
		c.JSON(http.StatusOK, gin.H{
			"result": markets,
			"count":  len(markets),
		})
	}
}

func GetUserMarkets() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if  _ , ok := userId.(int) ; exists && ok{
			marketRepository := database.GetMarketRepository()
			markets := marketRepository.GetUserMarkets(userId.(int))
			c.JSON(http.StatusOK, gin.H{
				"result": markets,
				"count":  len(markets),
			})
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not set",
		})

	}
}

func AddMarket() gin.HandlerFunc {
	return func(c *gin.Context) {
		marketRepository := database.GetMarketRepository()

		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Writer.WriteHeader(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
		}

		market := model.Market{}

		if err := json.Unmarshal(data, &market); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
		}

		if market.Name == "" || market.Description == "" || market.Date == "" || market.Lat == 0 || market.Lon == 0{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "name,description,date,lat,lon are mandatory parameters",
			})
			c.Abort()
		}

		marketRepository.Create(market)
	}
}
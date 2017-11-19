package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jcuerdo/mymarket-app-go/database"
	"io/ioutil"
	"github.com/jcuerdo/mymarket-app-go/model"
	"encoding/json"
	"strconv"
)

func GetMarkets() gin.HandlerFunc {
	return func(c *gin.Context) {
		marketFilter := model.MarketFilter{}
		marketFilter.Lat, _ = strconv.ParseFloat(c.Query("lat"),64)
		marketFilter.Lon, _ = strconv.ParseFloat(c.Query("lon"),64)
		marketFilter.Radio, _ = strconv.ParseFloat(c.Query("radio"),64)
		marketFilter.Page, _ = strconv.ParseInt(c.Query("page"),10,64)

		marketRepository := database.GetMarketRepository()
		markets := marketRepository.GetMarkets(marketFilter)
		c.JSON(http.StatusOK, gin.H{
			"result": markets,
			"count":  len(markets),
			"page":  marketFilter.Page,
		})
	}
}

func GetUserMarkets() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if _, ok := userId.(int); exists && ok {
			marketRepository := database.GetMarketRepository()
			markets := marketRepository.GetUserMarkets(userId.(int))
			c.JSON(http.StatusOK, gin.H{
				"result": markets,
				"count":  len(markets),
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID not set",
			})
		}

	}
}

func AddMarket() gin.HandlerFunc {
	return func(c *gin.Context) {
		marketRepository := database.GetMarketRepository()
		userId, _ := c.Get("userId")
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

		if market.Name == "" || market.Description == "" || market.Date == "" || market.Lat == 0 || market.Lon == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "name,description,date,lat,lon are mandatory parameters",
			})
			c.Abort()
		}

		marketRepository.Create(market, userId.(int))
	}
}

func EditMarket() gin.HandlerFunc {
	return func(c *gin.Context) {

		//TODO: Function is owner

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
		if market.Id == 0 || market.Name == "" || market.Description == "" || market.Date == "" || market.Lat == 0 || market.Lon == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "name,description,date,lat,lon are mandatory parameters",
			})
			c.Abort()
		}
		marketRepository.Edit(market)
	}
}

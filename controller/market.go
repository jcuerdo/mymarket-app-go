package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jcuerdo/mymarket-app-go/database"
	"io/ioutil"
	"github.com/jcuerdo/mymarket-app-go/model"
	"encoding/json"
	"strconv"
	"time"
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

func GetMarket() gin.HandlerFunc {
	return func(c *gin.Context) {
		marketIdParameter := c.Param("marketId")
		marketId, err := strconv.ParseInt(marketIdParameter, 10, 64)
		if err == nil {
			marketRepository := database.GetMarketRepository()
			defer marketRepository.Db.Close()
			market := marketRepository.GetMarket(marketId)
			if market.Id > 0 {
				c.JSON(http.StatusOK, gin.H{
					"result": market,
				})
				c.Abort()
			}
		}
		c.AbortWithStatus(http.StatusNotFound)
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
				"error": "name,description,startdate,lat,lon are mandatory parameters",
			})
			c.Abort()
		}

		datetime, _ := time.Parse(time.RFC3339,market.Date)

		market.Date = datetime.Format("2006-01-02 15:04:05")

		market.UserId = userId.(int)
		lastInsertedId := marketRepository.Create(market)
		market.Id = lastInsertedId
		c.JSON(http.StatusCreated, gin.H{
			"result": market,
		})
		c.Abort()
	}
}

func EditMarket() gin.HandlerFunc {
	return func(c *gin.Context) {

		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Writer.WriteHeader(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
		}

		marketModifications := model.Market{}

		if err := json.Unmarshal(data, &marketModifications); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
		}
		if marketModifications.Id == 0 || marketModifications.Name == "" || marketModifications.Description == "" || marketModifications.Date == "" || marketModifications.Lat == 0 || marketModifications.Lon == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "name,description,date,lat,lon are mandatory parameters",
			})
			c.Abort()
		}

		datetime, _ := time.Parse(time.RFC3339,marketModifications.Date)

		marketModifications.Date = datetime.Format("2006-01-02 15:04:05")

		marketRepository := database.GetMarketRepository()
		marketDb := marketRepository.GetMarket(marketModifications.Id)

		marketDb.Name = marketModifications.Name
		marketDb.Description = marketModifications.Description
		marketDb.Date = marketModifications.Date
		marketDb.Lat = marketModifications.Lat
		marketDb.Lon = marketModifications.Lon

		userId, _ := c.Get("userId")
		if isOwner(userId.(int), marketDb) {
			marketRepository := database.GetMarketRepository()
			if marketRepository.Edit(marketDb){
				c.JSON(http.StatusCreated, gin.H{
					"result": marketDb,
				})
				c.Abort()
			} else{
				c.AbortWithStatus(http.StatusNotModified)
			}
		} else{
			c.AbortWithStatus(http.StatusUnauthorized)
		}


	}
}

func isOwner(userId int, market model.Market) (bool) {
	return userId == market.UserId
}

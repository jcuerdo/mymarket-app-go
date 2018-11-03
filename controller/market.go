package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jcuerdo/mymarket-app-go/database"
	"github.com/jcuerdo/mymarket-app-go/model"
	"io/ioutil"
	"net/http"
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
		marketId, _ := strconv.ParseInt(c.Param("marketId"),10,64)

		marketRepository := database.GetMarketRepository()
		defer marketRepository.Db.Close()
		market := marketRepository.GetMarket(marketId)
		if market.Id > 0 {
			c.JSON(http.StatusOK, gin.H{
				"result": market,
			})
			c.Abort()
			return
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
			return
		}

		market := model.MarketExportable{}

		if err := json.Unmarshal(data, &market); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
			return
		}

		if market.Name == "" || market.Description == "" || market.Date == "" || market.Lat == 0 || market.Lon == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "name,description,startdate,lat,lon are mandatory parameters",
			})
			c.Abort()
			return
		}

		datetime, _ := time.Parse(time.RFC3339,market.Date)

		market.Date = datetime.Format("2006-01-02 15:04:05")

		market.UserId = userId.(int)
		lastInsertedId := marketRepository.Create(market)

		if lastInsertedId < 0 {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": "Error creating market",
			})
			c.Abort()
			return
		}

		market.Id = lastInsertedId
		c.JSON(http.StatusCreated, gin.H{
			"result": market,
		})
		c.Abort()
		return
	}
}

func RepeatMarket() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := c.Get("userId")
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Writer.WriteHeader(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
			return
		}

		market := model.MarketExportable{}

		if err := json.Unmarshal(data, &market); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
			return
		}

		if market.Date == "" || market.Id == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "id and startdate  are mandatory parameters",
			})
			c.Abort()
			return
		}

		marketRepository := database.GetMarketRepository()

		marketDb := marketRepository.GetMarket(market.Id)

		if marketDb.UserId != userId {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "You have no privileges to clone this market",
			})
			c.Abort()
			return
		}

		marketRepository = database.GetMarketRepository()

		datetime, _ := time.Parse(time.RFC3339,market.Date)
		newDate := datetime.Format("2006-01-02 15:04:05")

		lastInsertedId := marketRepository.Repeat(marketDb, newDate)

		if lastInsertedId < 0 {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": "Error creating market repetition",
			})
			c.Abort()
			return
		}

		marketDb.Id = lastInsertedId
		c.JSON(http.StatusCreated, gin.H{
			"result": marketDb,
		})
		c.Abort()
		return
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
			return
		}

		marketModifications := model.MarketExportable{}

		if err := json.Unmarshal(data, &marketModifications); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
			return
		}
		if marketModifications.Id == 0 || marketModifications.Name == "" || marketModifications.Description == "" || marketModifications.Date == "" || marketModifications.Lat == 0 || marketModifications.Lon == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "id,name,description,date,lat,lon are mandatory parameters",
			})
			c.Abort()
			return
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

func DeleteMarket() gin.HandlerFunc {
	return func(c *gin.Context) {

		userId, existsUser := c.Get("userId")
		marketId, _ := strconv.Atoi(c.Param("marketId"))

		marketRepository := database.GetMarketRepository()
		marketDb := marketRepository.GetMarket(int64(marketId))

		if marketDb.Id == 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if _, ok := userId.(int); existsUser && ok && marketId > 0 {
			if isOwner(userId.(int), marketDb) {
				marketRepository := database.GetMarketRepository()
				if marketRepository.Delete(userId.(int), int64(marketId)) {
					c.AbortWithStatus(http.StatusCreated)
					return
				} else {
					c.AbortWithStatus(http.StatusBadGateway)
					return
				}
			} else {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

func isOwner(userId int, market model.MarketExportable) (bool) {
	return userId == market.UserId
}

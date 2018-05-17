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

func GetMarketPhotos() gin.HandlerFunc {
	return func(c *gin.Context) {
		marketIdParameter := c.Param("marketId")
		marketId, err := strconv.Atoi(marketIdParameter)
		if err == nil {
			photoRepository := database.GetPhotoRepository()
			markets := photoRepository.GetMarketPhotos(marketId)
			c.JSON(http.StatusOK, gin.H{
				"result": markets,
				"count":  len(markets),
			})
		}
	}
}

func GetMarketPhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		marketIdParameter := c.Param("marketId")
		marketId, err := strconv.Atoi(marketIdParameter)
		if err == nil {
			photoRepository := database.GetPhotoRepository()
			market := photoRepository.GetMarketPhoto(marketId)
			if market.Id > 0 {
				c.JSON(http.StatusOK, gin.H{
					"result": market,
				})
				c.Abort()
			} else {
				c.AbortWithStatus(http.StatusNotFound)
			}

		}
	}
}

func AddPhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		marketIdParameter := c.Param("marketId")
		marketId, err := strconv.Atoi(marketIdParameter)
		if err != nil{
			c.AbortWithStatus(http.StatusBadRequest)
		}
		photoRepository := database.GetPhotoRepository()
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Writer.WriteHeader(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
		}

		photo := model.Photo{}

		if err := json.Unmarshal(data, &photo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
		}

		if photo.Content == ""{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "content is a mandatory parameter",
			})
			c.Abort()
		}

		photoRepository.Create(photo,marketId)

		c.AbortWithStatus(http.StatusCreated)
	}
	}

func DeletePhotos() gin.HandlerFunc {
	return func(c *gin.Context) {
		marketIdParameter := c.Param("marketId")
		marketId, err := strconv.Atoi(marketIdParameter)
		if err != nil{
			c.AbortWithStatus(http.StatusBadRequest)
		}
		photoRepository := database.GetPhotoRepository()

		photoRepository.Delete(marketId)

		c.AbortWithStatus(http.StatusCreated)
	}
}
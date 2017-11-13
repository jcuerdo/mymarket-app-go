package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jcuerdo/mymarket-app-go/database"
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
		if exists{
			marketRepository := database.GetMarketRepository()
			markets := marketRepository.GetUserMarkets(userId.(int))
			c.JSON(http.StatusOK, gin.H{
				"result": markets,
				"count":  len(markets),
			})
		}
		//TODO: USER NOT FOUND
	}
}
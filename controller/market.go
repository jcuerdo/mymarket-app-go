package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jcuerdo/mymarket-app-go/repository"
	"net/http"
)

func GetMarkets() gin.HandlerFunc {
	return func(c *gin.Context) {
		//user, _ := c.Get("userId")
		markets := repository.GetMarkets()
		c.JSON(http.StatusOK, gin.H{
			"result": markets,
			"count":  len(markets),
		})
	}
}
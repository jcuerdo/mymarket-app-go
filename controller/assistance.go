package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jcuerdo/mymarket-app-go/service"
	"net/http"
	"github.com/jcuerdo/mymarket-app-go/database"
	"io/ioutil"
	"github.com/jcuerdo/mymarket-app-go/model"
	"encoding/json"
	"strconv"
	"log"
)

func GetMarketAssistances() gin.HandlerFunc {
	return func(c *gin.Context) {
		marketIdParameter := c.Param("marketId")
		marketId, err := strconv.Atoi(marketIdParameter)
		if err == nil {
			assistanceRepository := database.GetAssistanceRepository()
			assistances := assistanceRepository.GetMarketAssistance(marketId)
			c.JSON(http.StatusOK, gin.H{
				"result": assistances,
				"count":  len(assistances),
			})
		}
	}
}

func AddAssistance() gin.HandlerFunc {
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

		assistance := model.Assistance{}

		if err := json.Unmarshal(data, &assistance); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
			return
		}

		if assistance.MarketId == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "market_id is a mandatory parameter",
			})
			c.Abort()
			return
		}

		userId, exists := c.Get("userId")
		if _, ok := userId.(int); exists && ok {
			assistance.UserId = userId.(int)
			assistanceRepository := database.GetAssistanceRepository()
			if assistanceRepository.Create(assistance) {

				userRepository := database.GetUserRepository()
				ids := userRepository.GetUserTokensInvolvedInMarket(assistance.MarketId)
				owner := userRepository.GetUserTokenMarketOwner(assistance.MarketId)
				ids = append(ids, owner)

				notificationService := service.NewNotificatorService()
				notificationService.NotifyAssistanceToAll(ids, assistance)

				c.AbortWithStatus(http.StatusCreated)

			} else {
				c.AbortWithStatus(http.StatusNotModified)
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

	}
}

func DeleteAssistance() gin.HandlerFunc {
	return func(c *gin.Context) {

		userId, existsUser := c.Get("userId")
		assistanceId, _ := strconv.Atoi(c.Param("assistanceId"))

		if _, ok := userId.(int); existsUser && ok && assistanceId > 0 {
			assistanceRepository := database.GetAssistanceRepository()
			log.Println(assistanceId,userId)
			if assistanceRepository.Delete(userId.(int),assistanceId) {
				c.AbortWithStatus(http.StatusCreated)
			} else {
				c.AbortWithStatus(http.StatusNotFound)
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

	}
}
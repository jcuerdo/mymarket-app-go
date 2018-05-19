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

func GetMarketComments() gin.HandlerFunc {
	return func(c *gin.Context) {
		marketIdParameter := c.Param("marketId")
		marketId, err := strconv.Atoi(marketIdParameter)
		if err == nil {
			commentRepository := database.GetCommentRepository()
			comments := commentRepository.GetMarketComments(marketId)
			c.JSON(http.StatusOK, gin.H{
				"result": comments,
				"count":  len(comments),
			})
		}
	}
}

func AddComment() gin.HandlerFunc {
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

		comment := model.Comment{}

		if err := json.Unmarshal(data, &comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
			return
		}

		if comment.MarketId == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "market_id is a mandatory parameter",
			})
			c.Abort()
			return
		}

		if comment.Content == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "content is a mandatory parameter",
			})
			c.Abort()
			return
		}

		userId, exists := c.Get("userId")
		if _, ok := userId.(int); exists && ok {
			comment.UserId = userId.(int)
			commentRepository := database.GetCommentRepository()
			if commentRepository.Create(comment) {
				c.AbortWithStatus(http.StatusCreated)
			} else {
				c.AbortWithStatus(http.StatusServiceUnavailable)
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

	}
}

func DeleteComment() gin.HandlerFunc {
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

		comment := model.Comment{}

		if err := json.Unmarshal(data, &comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
			return
		}

		if comment.Id == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "id is a mandatory parameter",
			})
			c.Abort()
			return
		}

		userId, exists := c.Get("userId")
		if _, ok := userId.(int); exists && ok {
			comment.UserId = userId.(int)
			commentRepository := database.GetCommentRepository()
			if commentRepository.Delete(comment) {
				c.AbortWithStatus(http.StatusCreated)
			} else {
				c.AbortWithStatus(http.StatusNotFound)
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

	}
}
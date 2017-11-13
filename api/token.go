package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		_, exists := query["token"]
		if exists{
			c.Set("userId", 1)
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "'token' field for authorization is required",
		})
	}
}
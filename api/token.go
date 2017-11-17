package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		
		if token != ""{
			c.Set("userId", 1)
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "'token' field for authorization is required",
		})
		c.Abort()
	}
}
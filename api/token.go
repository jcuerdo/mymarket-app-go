package api

import (
	"github.com/gin-gonic/gin"
)

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		token, exists := query["token"]
		if exists{
			c.Set("userId", token[0])
			return
		}
        c.AbortWithStatus(403)
	}
}
package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jcuerdo/mymarket-app-go/database"
)

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")

		if token != ""{
			userRepository := database.GetUserRepository()
			userId := userRepository.GetUserIdByToken(token)
			if userId > 0 {
				c.Set("userId", userId)
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token not valid",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "token field for authorization is required",
		})
		c.Abort()
	}
}
package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jcuerdo/mymarket-app-go/database"
	"io/ioutil"
	"github.com/jcuerdo/mymarket-app-go/model"
	"encoding/json"
	"github.com/satori/go.uuid"
	"time"
	"strconv"
	"strings"
)

func AddUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRepository := database.GetUserRepository()
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Writer.WriteHeader(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
		}

		user := model.User{}

		if err := json.Unmarshal(data, &user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
		}

		if user.Email == "" || user.Password == ""{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "email and password is a mandatory parameter",
			})
			c.Abort()
		}

		userRepository.CreateUser(user)

		c.AbortWithStatus(http.StatusCreated)
	}
}

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Writer.WriteHeader(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
		}

		user := model.User{}

		if err := json.Unmarshal(data, &user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
		}

		if user.Email == "" || user.Password == ""{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "email and password is a mandatory parameter",
			})
			c.Abort()
		}

		userRepository := database.GetUserRepository()

		user = userRepository.GetUser(user.Email,user.Password)

		if user.Id != 0 {
			token := generateToken()
			if userRepository.CreateToken(user.Id, token) {
				c.JSON(http.StatusAccepted, gin.H{
					"result": token,
				})
				c.Abort()
				}
			}
			c.AbortWithStatus(http.StatusUnauthorized)

	}
}
func generateToken() string {
	uuid := uuid.NewV4().String()
	uuid = strconv.Itoa(int(time.Now().Second())) + strings.Replace(uuid, "-" , "" ,-1) + strconv.Itoa(int(time.Now().Nanosecond()))

	return uuid
}
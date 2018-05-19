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
	"log"
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
			return
		}

		user := model.User{}

		if err := json.Unmarshal(data, &user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
			return
		}

		if user.Email == "" || user.Password == ""{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "email and password is a mandatory parameter",
			})
			c.Abort()
			return
		}

		if userRepository.CreateUser(user) {
			c.AbortWithStatus(http.StatusCreated)
		}
		c.AbortWithStatus(http.StatusConflict)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRepository := database.GetUserRepository()
		userId, exists := c.Get("userId")
		if !exists{
			log.Println("User not exists")
			c.AbortWithStatus(http.StatusNotFound)
		}
		user := userRepository.GetUserById(userId.(int))
		if  user.Id > 0{
			log.Println("User found")
			userExportable := model.UserExportable{
				user.Id,
				user.Email,
				user.Password,
				user.FullName.String,
				user.Photo.String,
				user.Description.String,
				user.Role.String,
			}
			c.JSON(http.StatusOK, gin.H{
				"result": userExportable,
			})
			c.Abort()
			return
		}
		log.Println("User not found")
		c.AbortWithStatus(http.StatusNotFound)
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
			return
		}

		user := model.User{}

		if err := json.Unmarshal(data, &user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
			return
		}

		if user.Email == "" || user.Password == ""{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "email and password is a mandatory parameter",
			})
			c.Abort()
			return
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
				return
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
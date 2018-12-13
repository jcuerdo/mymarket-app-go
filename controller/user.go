package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
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

func UpdateUser() gin.HandlerFunc {
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

		user := model.UserUpdate{}

		if err := json.Unmarshal(data, &user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parameters " + err.Error(),
			})
			c.Abort()
			return
		}

		if user.Id <= 0 || user.Email == "" || user.Password == ""{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "id, email and password is a mandatory parameter",
			})
			c.Abort()
			return
		}

		if ok, error := userRepository.UpdateUser(user) ; ok == true {
			c.AbortWithStatus(http.StatusCreated)
			return
		} else {
			if error != nil{
				mySqlError , isSqlError := error.(*mysql.MySQLError)
				if isSqlError {
					if mySqlError.Number == 1062 {
						c.AbortWithStatus(http.StatusConflict)
						return
					}
				}

			}
		}
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

func obtainUser(userId int) model.User {
		userRepository := database.GetUserRepository()
		user := userRepository.GetUserById(userId)
		if  user.Id > 0{
			return user
		}
		return model.User{}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if !exists{
			log.Println("User not exists")
			c.AbortWithStatus(http.StatusNotFound)
		}
		user := obtainUser(userId.(int))
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

func GetUserPublic() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdParameter := c.Param("userId")
		userId, error := strconv.Atoi(userIdParameter)
		if error != nil{
			log.Println("User not exists")
			c.AbortWithStatus(http.StatusNotFound)
		}
		user := obtainUser(userId)
		if  user.Id > 0{
			log.Println("User found")
			userExportablePublic := model.UserExportablePublic{
				user.Id,
				user.Email,
				user.FullName.String,
				user.Photo.String,
				user.Description.String,
			}
			c.JSON(http.StatusOK, gin.H{
				"result": userExportablePublic,
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

func SendFirebaseToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		userId, exists := c.Get("userId")
		firebaseToken := c.Param("firebasetoken")

		if _, ok := userId.(int); exists && ok && firebaseToken != "" {
			userRepository := database.GetUserRepository()
			if userRepository.UpdateFirebaseToken(userId.(int), firebaseToken) {
				c.AbortWithStatus(http.StatusOK)
			} else{
				c.AbortWithStatus(http.StatusBadRequest)
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func generateToken() string {
	uuid := uuid.NewV4().String()
	uuid = strconv.Itoa(int(time.Now().Second())) + strings.Replace(uuid, "-" , "" ,-1) + strconv.Itoa(int(time.Now().Nanosecond()))

	return uuid
}
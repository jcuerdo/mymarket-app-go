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

func AddUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Todo create user
	}
}

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Todo login user
		//Todo If exist create token
	}
}
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jcuerdo/mymarket-app-go/api"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jcuerdo/mymarket-app-go/controller"
)

func main() {
	router := gin.Default()

	public := router.Group("/public/", api.Cors())
	private := router.Group("/private/", api.Cors(), api.ValidateToken())

	private.GET("/market", controller.GetMarkets())
	public.GET("/market", controller.GetMarkets())

	router.Run(":8080")
}

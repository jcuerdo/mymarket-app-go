package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jcuerdo/mymarket-app-go/api"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jcuerdo/mymarket-app-go/controller"
	"os"
	"io"
	"log"
)

func main() {

	router := gin.Default()

	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	gin.DefaultWriter = io.MultiWriter(f)
	gin.DefaultErrorWriter = io.MultiWriter(f)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	public, private := defineGroups(router)
	definePrivateRoutes(private)
	definePublicRoutes(public)
	router.Run(":8080")
}
func definePublicRoutes(public *gin.RouterGroup) {
	public.GET("/market/:marketId/photos", controller.GetMarketPhotos())
	public.GET("/market/:marketId/photo", controller.GetMarketPhoto())
	public.GET("/market", controller.GetMarkets())
	public.POST("/user/create", controller.AddUser())
	public.POST("/user/login", controller.LoginUser())
}
func definePrivateRoutes(private *gin.RouterGroup) {
	private.GET("/market", controller.GetUserMarkets())
	private.POST("/market", controller.AddMarket())
	private.OPTIONS("/market", api.Options())
	private.POST("/market/:marketId/edit", controller.EditMarket())
	private.POST("/market/:marketId/photo", controller.AddPhoto())
}
func defineGroups(router *gin.Engine) (*gin.RouterGroup, *gin.RouterGroup) {
	public := router.Group("/public/", api.Cors())
	private := router.Group("/private/", api.Cors(), api.ValidateToken())
	return public, private
}

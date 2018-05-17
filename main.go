package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jcuerdo/mymarket-app-go/api"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jcuerdo/mymarket-app-go/controller"
	"os"
	"io"
	"log"
	config2 "github.com/jcuerdo/mymarket-app-go/config"
)

const PARAMETERS_FILE = "parameters.yml"

func main() {
	loader := config2.Loader{PARAMETERS_FILE}
	config := loader.Load()

	router := gin.Default()

	f, err := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	multiWritter := io.MultiWriter(f, os.Stdout)
	gin.DefaultWriter = multiWritter
	gin.DefaultErrorWriter = multiWritter
	log.SetOutput(multiWritter)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	public, private := defineGroups(router)
	definePrivateRoutes(private)
	definePublicRoutes(public)
	router.Run(":8080")
}
func definePublicRoutes(public *gin.RouterGroup) {
	//Markets
	public.GET("/market/:marketId/photos", controller.GetMarketPhotos())
	public.GET("/market/:marketId/photo", controller.GetMarketPhoto())
	public.GET("/market", controller.GetMarkets())
	public.GET("/market/:marketId", controller.GetMarket())

	//User
	public.POST("/user/create", controller.AddUser())
	public.OPTIONS("/user/create", api.Options())
	public.POST("/user/login", controller.LoginUser())
	public.OPTIONS("/user/login", api.Options())
}
func definePrivateRoutes(private *gin.RouterGroup) {
	//Markets
	private.GET("/market", controller.GetUserMarkets())
	private.POST("/market", controller.AddMarket())
	private.OPTIONS("/market", api.Options())
	private.POST("/market/:marketId/edit", controller.EditMarket())
	private.OPTIONS("/market/:marketId/edit", api.Options())
	private.POST("/market/:marketId/photo", controller.AddPhoto())
	private.OPTIONS("/market/:marketId/photo", api.Options())
	private.DELETE("/market/:marketId/photo", controller.DeletePhotos())

	//Users
	private.GET("/user", controller.GetUser())

}
func defineGroups(router *gin.Engine) (*gin.RouterGroup, *gin.RouterGroup) {
	public := router.Group("/public/", api.Cors())
	private := router.Group("/private/", api.Cors(), api.ValidateToken())
	return public, private
}

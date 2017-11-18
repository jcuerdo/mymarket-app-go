package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jcuerdo/mymarket-app-go/api"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jcuerdo/mymarket-app-go/controller"
)

func main() {
	router := gin.Default()

	public, private := defineGroups(router)
	definePrivateRoutes(private)
	definePublicRoutes(public)
	router.Run(":8080")
}
func definePublicRoutes(public *gin.RouterGroup) {
	public.GET("/market/:marketId/photos", controller.GetMarketPhotos())
	public.GET("/market/:marketId/photo", controller.GetMarketPhoto())
	public.GET("/market", controller.GetMarkets())
}
func definePrivateRoutes(private *gin.RouterGroup) {
	private.GET("/market", controller.GetUserMarkets())
	private.POST("/market", controller.AddMarket())
	private.POST("/market/:marketId/edit", controller.EditMarket())
	private.POST("/market/:marketId/photo/add", controller.AddPhoto())
}
func defineGroups(router *gin.Engine) (*gin.RouterGroup, *gin.RouterGroup) {
	public := router.Group("/public/", api.Cors())
	private := router.Group("/private/", api.Cors(), api.ValidateToken())
	return public, private
}

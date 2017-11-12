package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jcuerdo/mymarket-app-go/api"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"net/http"
	"github.com/jcuerdo/mymarket-app-go/model"
)


func main() {
	router := gin.Default()

	public := router.Group("/public/", api.Cors())
	private := router.Group("/private/",  api.Cors(), api.ValidateToken())


	private.GET("/market", func(c *gin.Context) {
		user, _ := c.Get("userId")
		c.JSON(200, gin.H{
			"message": user,
		})
	})

	public.GET("/market", func(c *gin.Context) {
		db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/database")
		if err == nil{
			rows, queryError := db.Query("SELECT id,name,description from market")
			fmt.Println(queryError)
			if queryError == nil{
				var markets []model.Market
				for rows.Next() {
					var market model.Market
					err = rows.Scan(&market.Id, &market.Description, &market.Name)
					markets = append(markets, market)
					if err != nil {
						fmt.Print(err.Error())
					}
				}
				defer rows.Close()
				c.JSON(http.StatusOK, gin.H{
					"result": markets,
					"count":  len(markets),
				})
			}
		}
	})


	router.Run(":8080")
}


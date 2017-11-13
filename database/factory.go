package database

import (
	"database/sql"
	"log"
	"github.com/jcuerdo/mymarket-app-go/repository"
)

func getDatabase() (*sql.DB)  {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/database")
	if err == nil{
		return db
	} else {
		log.Fatal(err)
		panic(err)
		return nil
	}
}

func GetMarketRepository() (repository.MarketRepository){
	marketRepository := repository.MarketRepository{getDatabase()}
	return marketRepository
}

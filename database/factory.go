package database

import (
	"database/sql"
	"log"
	"github.com/jcuerdo/mymarket-app-go/repository"
	"os"
)

func getDatabase() (*sql.DB)  {
	dataSource := os.Getenv("DATASOURCE")
	db, err := sql.Open("mysql", dataSource)
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


func GetPhotoRepository() (repository.PhotoRepository){
	photoRepository := repository.PhotoRepository{getDatabase()}
	return photoRepository
}

func GetUserRepository() (repository.UserRepository){
	userRepository := repository.UserRepository{getDatabase()}
	return userRepository
}
package database

import (
	"database/sql"
	"log"
	"github.com/jcuerdo/mymarket-app-go/repository"
)

func getDatabase() (*sql.DB)  {
	db, err := sql.Open("mysql", "root:123456@tcp(ec2-34-215-191-148.us-west-2.compute.amazonaws.com:3306)/database")
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
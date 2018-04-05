package database

import (
	"database/sql"
	"github.com/jcuerdo/mymarket-app-go/repository"
	config2 "github.com/jcuerdo/mymarket-app-go/config"
	"log"
)

const PARAMETERS_FILE = "parameters.yml"

func getDatabase() (*sql.DB)  {
	loader := config2.Loader{PARAMETERS_FILE}
	config :=  loader.Load()
	db, err := sql.Open("mysql", config.Datasource)
	db.SetMaxOpenConns(10)

	if err == nil{
		return db
	} else {
		log.Fatal(err)
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
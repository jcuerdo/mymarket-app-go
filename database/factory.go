package database

import (
	"database/sql"
	"github.com/jcuerdo/mymarket-app-go/repository"
	config2 "github.com/jcuerdo/mymarket-app-go/config"
	"log"
)

const PARAMETERS_FILE = "parameters.yml"
const DRIVER = "mysql"
const MAX_CONNECTIONS = 10

func getDatabase() (*sql.DB)  {
	loader := config2.Loader{PARAMETERS_FILE}
	config :=  loader.Load()
	db, err := sql.Open(DRIVER, config.Datasource)
	db.SetMaxOpenConns(MAX_CONNECTIONS)

	if err == nil{
		pingErr := db.Ping()
		if pingErr != nil {
			log.Println(pingErr)
			panic(pingErr)
		}
		return db
	} else {
		log.Println(err)
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

func GetCommentRepository() (repository.CommentRepository){
	commentRepository := repository.CommentRepository{getDatabase()}
	return commentRepository
}

func GetAssistanceRepository() (repository.AssistanceRepository){
	assistanceRepository := repository.AssistanceRepository{getDatabase()}
	return assistanceRepository
}
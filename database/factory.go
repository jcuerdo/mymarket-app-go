package database

import (
	"database/sql"
	"log"
)

func GetDatabase() (*sql.DB)  {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/database")
	if err == nil{
		return db
	} else {
		log.Fatal(err)
		panic(err)
	}
}

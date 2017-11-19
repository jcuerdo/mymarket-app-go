package model

import "database/sql"

type User struct {
	Id           int 			`json:"id"`
	Email      	 string 		`json:"email"`
	Password     string 		`json:"password"`
	FullName     sql.NullString `json:"fullname"`
	Photo        sql.NullString `json:"photo"`
	Description  sql.NullString `json:"description"`
	Role         sql.NullString `json:"role"`
}
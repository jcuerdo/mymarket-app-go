package repository

import (
	"github.com/jcuerdo/mymarket-app-go/model"
	"database/sql"
	"fmt"
)

type UserRepository struct {
	Db *sql.DB
}

func (userRepository *UserRepository)GetUser(email string,password string) (model.User) {
	rows, error := 	userRepository.Db.Query("SELECT id, password, email, fullname, photo,description,role FROM user WHERE email = ? and password = ?", email, password)
	if error == nil{
		for rows.Next() {
			user, err := parseUserRow(rows)
			if err == nil {
				return user
			}
		}
	} else{
		fmt.Println(error)
	}
	return model.User{}
}

func (userRepository *UserRepository)CreateToken(userId int,token string) (bool) {
	rows, error := 	userRepository.Db.Query("INSERT INTO token (id,user_id,token) VALUES (null, ? , ?)", userId,token)
	defer rows.Close()
	if error != nil{
		fmt.Println(error)
		return true
	}
	return false
}

func (userRepository *UserRepository)GetUserIdByToken(token string) (int) {
	row := userRepository.Db.QueryRow(
		`
		SELECT user_id FROM token
		WHERE token = ?`,
		token,
	)

	var userId int
	error := row.Scan(&userId)

	if error != nil{
		fmt.Println(error)
	}
	return userId
}

func (userRepository *UserRepository)CreateUser(user model.User) (bool) {
	rows, error := 	userRepository.Db.Query(
		`
		INSERT INTO user
		(id,password,email,fullname,photo,description,role)
		VALUES
		(null,?,?,?,?,?,?)`,
		user.Password,
		user.Email,
		user.FullName,
		user.Photo,
		user.Description,
		"USER",
		)
	defer rows.Close()
	if error != nil{
		fmt.Println(error)
	}
	return error == nil
}

func parseUserRow(rows *sql.Rows) (model.User, error) {
	var user model.User
	err := rows.Scan(&user.Id, &user.Password, &user.Email,&user.FullName,&user.Photo,&user.Description,&user.Role)

	if err != nil{
		fmt.Println(err)
	}
	return user, err
}


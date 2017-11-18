package repository

import (
	"github.com/jcuerdo/mymarket-app-go/model"
	"database/sql"
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
	}
	return model.User{}
}

func (userRepository *UserRepository)CreateToken(userId int,token string) (bool) {
	rows, error := 	userRepository.Db.Query("INSERT INTO token (id,user_id,token) VALUES (null, ? , ?)", userId,token)
	defer rows.Close()
	if error == nil{
		return false
	}
	return true
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
	return error == nil
}

func parseUserRow(rows *sql.Rows) (model.User, error) {
	var user model.User
	err := rows.Scan(&user.Id, &user.Password, &user.Email,&user.FullName,&user.Photo,&user.Description,&user.Role)
	return user, err
}


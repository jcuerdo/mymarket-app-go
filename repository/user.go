package repository

import (
	"github.com/jcuerdo/mymarket-app-go/model"
	"database/sql"
	"log"
)

type UserRepository struct {
	Db *sql.DB
}

func (userRepository *UserRepository) GetUser(email string, password string) (model.User) {
	stmt, error := userRepository.Db.Prepare("SELECT id, password, email, fullname, photo,description,role FROM user WHERE email = ? and password = ?")
	defer stmt.Close()

	if error != nil {
		log.Println(error)
		return model.User{}
	}

	row := stmt.QueryRow(email, password)
	if error == nil {
		user, error := parseUserRow(row)
		if error == nil {
			return user
		}
	}

	log.Println(error)
	return model.User{}
}

func (userRepository *UserRepository) GetUserTokensInvolvedInMarket(marketId int) ([]string) {
	var result = make([]string, 0)
	stmt, error := userRepository.Db.Prepare(
		`
			SELECT DISTINCT(firebase_token)
			FROM user
			LEFT JOIN comment on user.id = comment.user_id
			LEFT JOIN assistance on user.id = assistance.user_id
			WHERE comment.market_id=? or assistance.market_id = ?;
	`)
	defer stmt.Close()

	if error != nil {
		log.Println(error)
		return result
	}

	rows, error := stmt.Query(marketId,marketId)
	if error == nil {
		for rows.Next() {
			user, error := parseUserTokenRow(rows)
			if error == nil {
				result = append(result, user.FirebaseToken.String)
			} else{
				log.Println(error)
			}
		}
	} else {
		log.Println(error)
	}

	return result
}

func (userRepository *UserRepository) GetUserTokenMarketOwner(marketId int) (string) {
	stmt, error := userRepository.Db.Prepare(
		`
			SELECT DISTINCT(firebase_token)
			FROM user
			LEFT JOIN market on user.id = market.user_id
			WHERE market.id=?;
	`)
	defer stmt.Close()

	if error != nil {
		log.Println(error)
		return ""
	}

	rows, error := stmt.Query(marketId)
	if error == nil {
		for rows.Next() {
			user, error := parseUserTokenRow(rows)
			if error == nil {
				return user.FirebaseToken.String
			} else{
				log.Println(error)
			}
		}
	} else {
		log.Println(error)
	}
	return ""
}

func (userRepository *UserRepository) GetUserById(id int) (model.User) {
	stmt, error := userRepository.Db.Prepare("SELECT id, password, email, fullname, photo,description,role FROM user WHERE id = ?")
	defer stmt.Close()
	defer userRepository.Db.Close()

	if error != nil {
		log.Println(error)
		return model.User{}
	}

	row := stmt.QueryRow(id)
	if error == nil {
		user, error := parseUserRow(row)
		if error == nil {
			return user
		}
	}

	log.Println(error)
	return model.User{}
}

func (userRepository *UserRepository) CreateToken(userId int, token string) (bool) {
	stmt, error := userRepository.Db.Prepare("INSERT INTO token (id,user_id,token) VALUES (null, ? , ?)")
	defer stmt.Close()
	defer userRepository.Db.Close()

	if error != nil {
		log.Println(error)
		return false
	}
	_, error = stmt.Exec(userId, token)

	if error != nil {
		log.Println(error)
		return false
	}
	return true
}

func (userRepository *UserRepository) UpdateFirebaseToken(userId int, firebaseToken string) (bool) {
	stmt, error := userRepository.Db.Prepare("UPDATE user SET firebase_token = ? where id = ?")
	defer stmt.Close()
	defer userRepository.Db.Close()

	if error != nil {
		log.Println(error)
		return false
	}
	_, error = stmt.Exec(firebaseToken, userId)

	if error != nil {
		log.Println(error)
		return false
	}
	return true
}

func (userRepository *UserRepository) GetUserIdByToken(token string) (int) {
	stmt, error := userRepository.Db.Prepare(
		`
		SELECT user_id FROM token
		WHERE token = ?`)

	row := stmt.QueryRow(token)

	defer stmt.Close()
	var userId int
	error = row.Scan(&userId)

	if error != nil {
		log.Println(error)
	}
	return userId
}

func (userRepository *UserRepository) CreateUser(user model.User) (bool) {
	stmt, error := userRepository.Db.Prepare(
		`
		INSERT INTO user
		(id,password,email,fullname,photo,description,role)
		VALUES
		(null,?,?,?,?,?,?)`)

	_, error = stmt.Exec(
		user.Password,
		user.Email,
		user.FullName,
		user.Photo,
		user.Description,
		"USER")

	defer userRepository.Db.Close()
	defer stmt.Close()
	if error != nil {
		log.Println(error)
	}
	return error == nil
}

func (userRepository *UserRepository) UpdateUser(user model.User) (bool) {
	stmt, error := userRepository.Db.Prepare(
		`
		UPDATE user
		set (id,password,email,fullname,photo,description)
		VALUES
		(null,?,?,?,?,?,?)
		WHERE id = ?`)

	_, error = stmt.Exec(
		user.Password,
		user.Email,
		user.FullName,
		user.Photo,
		user.Description,
		user.Id)

	defer userRepository.Db.Close()
	defer stmt.Close()
	if error != nil {
		log.Println(error)
	}
	return error == nil
}

func parseUserRow(row *sql.Row) (model.User, error) {
	var user model.User
	err := row.Scan(&user.Id, &user.Password, &user.Email, &user.FullName, &user.Photo, &user.Description, &user.Role)

	if err != nil {
		log.Println(err)
	}
	return user, err
}

func parseUserTokenRow(row *sql.Rows) (model.UserToken, error) {
	var usertoken model.UserToken
	err := row.Scan(&usertoken.FirebaseToken)

	if err != nil {
		log.Println(err)
	}
	return usertoken, err
}

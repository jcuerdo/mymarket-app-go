package repository

import (
	"github.com/jcuerdo/mymarket-app-go/model"
	"database/sql"
	"log"
)

type AssistanceRepository struct {
	Db *sql.DB
}

func (assistanceRepository *AssistanceRepository)GetMarketAssistance(market int) ([]model.AssistanceResult) {
	stmt, error := 	assistanceRepository.Db.Prepare(
		`
		SELECT a.id,a.market_id,u.id as user_id,u.email,u.fullname,u.photo,a.date FROM assistance a
		INNER JOIN user u on u.id = a.user_id
		WHERE market_id = ? ORDER BY a.date DESC `)
	defer stmt.Close()
	defer assistanceRepository.Db.Close()
	rows , error := stmt.Query(market)
	if error != nil{
		log.Println(error)
	}
	return parseAssistanceRows(rows, error)
}


func (assistanceRepository *AssistanceRepository)Create(assistance model.Assistance) (bool) {
	stmt, error := 	assistanceRepository.Db.Prepare(
		`
		INSERT INTO assistance
		(id,market_id,user_id)
		VALUES
		(null,?,?)`)

	if error != nil{
		log.Println(error)
	}

	_, error = stmt.Exec(
		assistance.MarketId,
		assistance.UserId,
		)

	defer stmt.Close()
	defer assistanceRepository.Db.Close()
	if error != nil{
		log.Println(error)
	}
	return error == nil
}

func (assistanceRepository *AssistanceRepository)Delete(userId int , assistanceId int) (bool) {
	stmt, error := 	assistanceRepository.Db.Prepare(`DELETE FROM assistance WHERE id = ? and user_id = ?`)

	result, _ := stmt.Exec(assistanceId, userId)

	defer stmt.Close()
	defer assistanceRepository.Db.Close()
	if error != nil{
		log.Println(error)
		return false
	}
	affectedRows, _ := result.RowsAffected()

	return affectedRows > 0
}

func parseAssistanceRows(rows *sql.Rows, error error) []model.AssistanceResult {
	var assistances []model.AssistanceResult
	if error == nil {
		for rows.Next() {
			assistance, err := parseAssistanceRow(rows)
			if err != nil {
				log.Println(err)
			} else {
				assistances = append(assistances, assistance)
			}
		}
	} else{
		log.Println(error)
	}
	return assistances

}
func parseAssistanceRow(rows *sql.Rows) (model.AssistanceResult, error) {
	var assistance model.AssistanceResult
	err := rows.Scan(&assistance.Id, &assistance.MarketId, &assistance.User.Id,&assistance.User.Email, &assistance.User.FullName,&assistance.User.Photo, &assistance.Date)

	if err != nil{
		log.Println(err)
	}
	return assistance, err
}

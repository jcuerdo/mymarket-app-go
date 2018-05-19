package repository

import (
	"github.com/jcuerdo/mymarket-app-go/model"
	"database/sql"
	"log"
)

type CommentRepository struct {
	Db *sql.DB
}

func (commentRepository *CommentRepository)GetMarketComments(market int) ([]model.Comment) {
	stmt, error := 	commentRepository.Db.Prepare("SELECT * FROM comment WHERE market_id = ?")
	defer stmt.Close()
	defer commentRepository.Db.Close()
	rows , error := stmt.Query(market)
	if error != nil{
		log.Println(error)
	}
	return parseCommentRows(rows, error)
}


func (commentRepository *CommentRepository)Create(comment model.Comment) (bool) {
	stmt, error := 	commentRepository.Db.Prepare(
		`
		INSERT INTO comment
		(id,market_id,user_id,content)
		VALUES
		(null,?,?,?)`)

	stmt.Exec(
		comment.MarketId,
		comment.UserId,
		comment.Content,
		)

	defer stmt.Close()
	defer commentRepository.Db.Close()
	if error != nil{
		log.Println(error)
	}
	return error == nil
}

func (commentRepository *CommentRepository)Delete(comment model.Comment) (bool) {
	stmt, error := 	commentRepository.Db.Prepare(`DELETE FROM comment WHERE id = ? and user_id = ?`)

	result, _ := stmt.Exec(comment.Id, comment.UserId)

	defer stmt.Close()
	defer commentRepository.Db.Close()
	if error != nil{
		log.Println(error)
		return false
	}
	affectedRows, _ := result.RowsAffected()

	return affectedRows > 0
}

func parseCommentRows(rows *sql.Rows, error error) []model.Comment {
	var comments []model.Comment
	if error == nil {
		for rows.Next() {
			comment, err := parseCommentRow(rows)
			if err != nil {
				log.Println(err)
			} else {
				comments = append(comments, comment)
			}
		}
	} else{
		log.Println(error)
	}
	return comments

}
func parseCommentRow(rows *sql.Rows) (model.Comment, error) {
	var comment model.Comment
	err := rows.Scan(&comment.Id, &comment.MarketId, &comment.UserId, &comment.Content)

	if err != nil{
		log.Println(err)
	}
	return comment, err
}

package repository

import (
	"github.com/jcuerdo/mymarket-app-go/model"
	"database/sql"
	"log"
)

type CommentRepository struct {
	Db *sql.DB
}

func (commentRepository *CommentRepository)GetMarketComments(market int) ([]model.CommentResult) {
	stmt, error := 	commentRepository.Db.Prepare(
		`
		SELECT c.id,c.market_id,u.id as user_id,u.email,u.fullname,u.photo,c.content,c.date FROM comment c
		INNER JOIN user u on u.id = c.user_id
		WHERE market_id = ? ORDER BY c.date DESC `)
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

func (commentRepository *CommentRepository)Delete(userId int , commentId int) (bool) {
	stmt, error := 	commentRepository.Db.Prepare(`DELETE FROM comment WHERE id = ? and user_id = ?`)

	result, _ := stmt.Exec(commentId, userId)

	defer stmt.Close()
	defer commentRepository.Db.Close()
	if error != nil{
		log.Println(error)
		return false
	}
	affectedRows, _ := result.RowsAffected()

	return affectedRows > 0
}

func parseCommentRows(rows *sql.Rows, error error) []model.CommentResult {
	var comments []model.CommentResult
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
func parseCommentRow(rows *sql.Rows) (model.CommentResult, error) {
	var comment model.CommentResult
	err := rows.Scan(&comment.Id, &comment.MarketId, &comment.User.Id,&comment.User.Email, &comment.User.FullName,&comment.User.Photo, &comment.Content, &comment.Date)

	if err != nil{
		log.Println(err)
	}
	return comment, err
}

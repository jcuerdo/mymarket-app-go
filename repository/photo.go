package repository

import (
	"github.com/jcuerdo/mymarket-app-go/model"
	"database/sql"
	"log"
)

type PhotoRepository struct {
	Db *sql.DB
}

func (photoRepository *PhotoRepository)GetMarketPhotos(market int) ([]model.Photo) {
	stmt, error := 	photoRepository.Db.Prepare("SELECT id,content FROM photo WHERE market_id = ?")
	defer stmt.Close()
	defer photoRepository.Db.Close()
	rows , error := stmt.Query(market)
	if error != nil{
		log.Println(error)
		return []model.Photo{}
	}
	return parsePhotoRows(rows, error)
}

func (photoRepository *PhotoRepository)GetMarketPhoto(market int) (model.Photo) {
	stmt,error := photoRepository.Db.Prepare("SELECT id,content FROM photo WHERE market_id = ? limit 1")
	defer stmt.Close()
	defer photoRepository.Db.Close()
	if error != nil{
		log.Println(error)
		return model.Photo{}
	} else{
		row := stmt.QueryRow(market)
		var photo model.Photo
		row.Scan(&photo.Id, &photo.Content)

		return photo
	}
}

func (photoRepository *PhotoRepository)Create(photo model.Photo,marketId int64) (bool) {
	stmt, error := 	photoRepository.Db.Prepare(
		`
		INSERT INTO photo
		(id,content,market_id)
		VALUES
		(null,?,?)`)

	stmt.Exec(
		photo.Content,
		marketId)

	defer stmt.Close()
	defer photoRepository.Db.Close()
	if error != nil{
		log.Println(error)
	}
	return error == nil
}

func (photoRepository *PhotoRepository)Delete(marketId int) (bool) {
	stmt, error := 	photoRepository.Db.Prepare(
		`DELETE FROM photo where market_id=?`)

	stmt.Exec(marketId)

	defer stmt.Close()
	defer photoRepository.Db.Close()
	if error != nil{
		log.Println(error)
	}
	return error == nil
}

func (photoRepository *PhotoRepository) Edit(market model.Market) (bool) {
	stmt, error := 	photoRepository.Db.Prepare(
		`
		UPDATE market SET
		name = ?, description = ? , startdate = ?,lat = ?,lon = ?
		WHERE id = ?`)
	stmt.Exec(
		market.Name,
		market.Description,
		market.Date,
		market.Lat,
		market.Lon,
		market.Id)
	defer stmt.Close()
	defer photoRepository.Db.Close()
	if error != nil{
		log.Println(error)
	}
	return error == nil
}

func parsePhotoRows(rows *sql.Rows, error error) []model.Photo {
	var photos []model.Photo
	if error == nil {
		for rows.Next() {
			photo, err := parsePhotoRow(rows)
			if err != nil {
				log.Println(err)
			} else {
				photos = append(photos, photo)
			}
		}
	} else{
		log.Println(error)
	}
	return photos

}
func parsePhotoRow(rows *sql.Rows) (model.Photo, error) {
	var photo model.Photo
	err := rows.Scan(&photo.Id, &photo.Content)

	if err != nil{
		log.Println(err)
	}
	return photo, err
}

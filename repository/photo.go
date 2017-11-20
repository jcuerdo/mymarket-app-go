package repository

import (
	"github.com/jcuerdo/mymarket-app-go/model"
	"database/sql"
	"fmt"
)

type PhotoRepository struct {
	Db *sql.DB
}

func (photoRepository *PhotoRepository)GetMarketPhotos(market int) ([]model.Photo) {
	rows, error := 	photoRepository.Db.Query("SELECT id,content FROM photo WHERE market_id = ?", market)
	defer rows.Close()
	if error != nil{
		fmt.Println(error)
	}
	return parsePhotoRows(rows, error)
}

func (photoRepository *PhotoRepository)GetMarketPhoto(market int) (model.Photo) {
	row := photoRepository.Db.QueryRow("SELECT id,content FROM photo WHERE market_id = ? limit 1", market)
	var photo model.Photo
	row.Scan(&photo.Id, &photo.Content)

	return photo
}

func (photoRepository *PhotoRepository)Create(photo model.Photo,marketId int) (bool) {
	rows, error := 	photoRepository.Db.Query(
		`
		INSERT INTO photo
		(id,content,market_id)
		VALUES
		(null,?,?)`,
		photo.Content,
		marketId)
	defer rows.Close()
	if error != nil{
		fmt.Println(error)
	}
	return error == nil
}

func (photoRepository *PhotoRepository) Edit(market model.Market) (bool) {
	rows, error := 	photoRepository.Db.Query(
		`
		UPDATE market SET
		name = ?, description = ? , startdate = ?,lat = ?,lon = ?
		WHERE id = ?`,
		market.Name,
		market.Description,
		market.Date,
		market.Lat,
		market.Lon,
		market.Id,)
	defer rows.Close()
	if error != nil{
		fmt.Println(error)
	}
	return error == nil
}

func parsePhotoRows(rows *sql.Rows, error error) []model.Photo {
	var photos []model.Photo
	if error == nil {
		for rows.Next() {
			photo, err := parsePhotoRow(rows)
			if err != nil {
				fmt.Println(err)
			} else {
				photos = append(photos, photo)
			}
		}
	} else{
		fmt.Println(error)
	}
	return photos

}
func parsePhotoRow(rows *sql.Rows) (model.Photo, error) {
	var photo model.Photo
	err := rows.Scan(&photo.Id, &photo.Content)

	if err != nil{
		fmt.Println(err)
	}
	return photo, err
}

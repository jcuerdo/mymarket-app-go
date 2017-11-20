package repository

import (
	"github.com/jcuerdo/mymarket-app-go/model"
	"log"
	"database/sql"
)

type PhotoRepository struct {
	Db *sql.DB
}

func (photoRepository *PhotoRepository)GetMarketPhotos(market int) ([]model.Photo) {
	rows, error := 	photoRepository.Db.Query("SELECT id,content FROM photo WHERE market_id = ?", market)
	defer rows.Close()
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
	return error == nil
}

func parsePhotoRows(rows *sql.Rows, error error) []model.Photo {
	var photos []model.Photo
	if error == nil {
		for rows.Next() {
			photo, err := parsePhotoRow(rows)
			if err != nil {
				log.Fatal(error)

			} else {
				photos = append(photos, photo)
			}
		}
	}
	return photos

}
func parsePhotoRow(rows *sql.Rows) (model.Photo, error) {
	var photo model.Photo
	err := rows.Scan(&photo.Id, &photo.Content)
	return photo, err
}

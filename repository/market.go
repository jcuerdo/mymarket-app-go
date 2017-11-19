package repository

import (
	"github.com/jcuerdo/mymarket-app-go/model"
	"log"
	"database/sql"
)

type MarketRepository struct {
	Db *sql.DB
}

func (marketRepository *MarketRepository)GetUserMarkets(user int) ([]model.Market) {
	rows, error := 	marketRepository.Db.Query("SELECT id,name,description,startdate,lat, lon FROM market WHERE active = 1 and user_id = ?", user)
	defer rows.Close()
	return parseRows(rows, error)
}

func (marketRepository *MarketRepository)GetMarkets() ([]model.Market) {
	rows, error := 	marketRepository.Db.Query("SELECT id,name,description,startdate,lat, lon FROM market WHERE active = 1")
	defer rows.Close()
	return parseRows(rows, error)
}

func (marketRepository *MarketRepository)Create(market model.Market, userId int) (bool) {
	rows, error := 	marketRepository.Db.Query(
		`
		INSERT INTO market
		(id,name,description,startdate,lat,lon,active,user_id)
		VALUES
		(null,?,?,?,?,?,?)`,
		market.Name,
		market.Description,
		market.Date,
 		market.Lat,
 		market.Lon,
 		userId,
 		true)
	defer rows.Close()
	return error == nil
}

func (marketRepository *MarketRepository) Edit(market model.Market) (bool) {
	rows, error := 	marketRepository.Db.Query(
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

func parseRows(rows *sql.Rows, error error) []model.Market {
	if error == nil {
		var markets []model.Market
		for rows.Next() {
			market, err := parseRow(rows)
			if err != nil {
				log.Fatal(error)

			} else {
				markets = append(markets, market)
			}
		}
		return markets
	}
	return nil

}
func parseRow(rows *sql.Rows) (model.Market, error) {
	var market model.Market
	err := rows.Scan(&market.Id, &market.Description, &market.Name, &market.Date, &market.Lat, &market.Lon)
	return market, err
}

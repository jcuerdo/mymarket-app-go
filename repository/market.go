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
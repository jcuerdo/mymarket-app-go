package repository

import (
	"github.com/jcuerdo/mymarket-app-go/model"
	"github.com/jcuerdo/mymarket-app-go/database"
	"log"
	"database/sql"
)

func GetMarkets() ([]model.Market) {
	db := database.GetDatabase()
	rows, error := db.Query("SELECT id,name,description,startdate,lat, lon FROM market WHERE active = 1")
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
				panic(error)
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

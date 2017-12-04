package repository

import (
	"github.com/jcuerdo/mymarket-app-go/model"
	"database/sql"
	"math"
	"fmt"
)

const EARTH_RATE  =  6371
const RADIO  =  180
const MAX_RESULTS  =  5

type MarketRepository struct {
	Db *sql.DB
}

func (marketRepository *MarketRepository)GetUserMarkets(user int) ([]model.Market) {
	stmt, error := 	marketRepository.Db.Prepare("SELECT id,user_id,name,description,startdate,lat, lon FROM market WHERE active = 1 and user_id = ?")
	rows, error := stmt.Query(user)
	defer rows.Close()
	defer marketRepository.Db.Close()
	return parseRows(rows, error)
}

func (marketRepository *MarketRepository)GetMarket(marketId int) (model.Market) {
	stmt, error := marketRepository.Db.Prepare("SELECT id,user_id,name,description,startdate,lat, lon FROM market WHERE active = 1 and id = ?")
	if error != nil{
		fmt.Println(error)
		return model.Market{}
	}
	row := stmt.QueryRow(marketId)
	defer marketRepository.Db.Close()
	var market model.Market
	row.Scan(&market.Id, &market.UserId, &market.Description, &market.Name, &market.Date, &market.Lat, &market.Lon)
	return market
}

func (marketRepository *MarketRepository)GetMarkets(marketFilter model.MarketFilter) ([]model.Market) {

	maxlat := marketFilter.Lat + rad2deg(marketFilter.Radio/EARTH_RATE)
	minlat := marketFilter.Lat - rad2deg(marketFilter.Radio/EARTH_RATE)
	maxlon := marketFilter.Lon + rad2deg(math.Asin(marketFilter.Radio/EARTH_RATE) / math.Cos(deg2rad(marketFilter.Lat)))
	minlon := marketFilter.Lon - rad2deg(math.Asin(marketFilter.Radio/EARTH_RATE) / math.Cos(deg2rad(marketFilter.Lat)))

	stmt, error := 	marketRepository.Db.Prepare(`
		SELECT
			id,user_id,name,description,startdate,lat, lon
		FROM
			market
		WHERE
			active = 1 AND
			lat <= ?   AND
			lat >= ?   AND
			lon <= ?   AND
			lon >= ?
		LIMIT ?,?
			`)
	defer stmt.Close()
	defer marketRepository.Db.Close()

	if error != nil {
		fmt.Println(error)
	}

	rows , error := stmt.Query(
		maxlat,
		minlat,
		maxlon,
		minlon,
		marketFilter.Page * MAX_RESULTS,
		MAX_RESULTS)

	if error != nil {
		fmt.Println(error)
	}

	if error != nil{
		fmt.Println(error)
		return nil
	}

	return parseRows(rows, error)
}

func deg2rad(deg float64) float64 {
	return deg * math.Pi / RADIO
}
func rad2deg(rad float64) float64 {
	return rad * RADIO / math.Pi
}

func (marketRepository *MarketRepository)Create(market model.Market) {
	stmt, error := 	marketRepository.Db.Prepare(
		`
		INSERT INTO market
		(id,name,description,startdate,lat,lon,active,user_id)
		VALUES
		(null,?,?,?,?,?,?,?)`)
	_ , error = stmt.Exec(
		market.Name,
		market.Description,
		market.Date,
		market.Lat,
		market.Lon,
		market.UserId,
		true)

	defer stmt.Close()
	defer marketRepository.Db.Close()

 	if error != nil{
 		fmt.Println(error)
	}
}

func (marketRepository *MarketRepository) Edit(market model.Market) (bool) {
	stmt, error := 	marketRepository.Db.Prepare(
		`
		UPDATE market SET
		name = ?, description = ? , startdate = ?,lat = ?,lon = ?
		WHERE id = ?`)

	_, error = stmt.Exec(
		market.Name,
		market.Description,
		market.Date,
		market.Lat,
		market.Lon,
		market.Id)

	defer stmt.Close()
	defer marketRepository.Db.Close()

	if error != nil{
		fmt.Println(error)
	}
	return error == nil
}

func parseRows(rows *sql.Rows, error error) []model.Market {
	if error == nil {
		var markets []model.Market
		for rows.Next() {
			market, err := parseRow(rows)
			if err != nil {
				fmt.Println(error)
			} else {
				markets = append(markets, market)
			}
		}
		return markets
	} else{
		fmt.Println(error)
	}

	return nil

}
func parseRow(rows *sql.Rows) (model.Market, error) {
	var market model.Market
	err := rows.Scan(&market.Id, &market.UserId, &market.Description, &market.Name, &market.Date, &market.Lat, &market.Lon)
	if err != nil{
		fmt.Println(err)
	}
	return market, err
}

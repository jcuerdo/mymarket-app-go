package repository

import (
	"database/sql"
	"github.com/jcuerdo/mymarket-app-go/model"
	"log"
	"math"
)

const EARTH_RATE  =  6371
const RADIO  =  180
const MAX_RESULTS  =  5

type MarketRepository struct {
	Db *sql.DB
}

func (marketRepository *MarketRepository)GetUserMarkets(user int) ([]model.MarketExportable) {
	stmt, error := 	marketRepository.Db.Prepare("SELECT id,user_id,name,description,startdate,lat, lon, market_type, flexible, place FROM market WHERE active = 1 and user_id = ?")
	if error != nil{
		log.Println(error)
		return []model.MarketExportable{}
	}
	rows, error := stmt.Query(user)
	if error != nil{
		log.Println(error)
		return []model.MarketExportable{}
	}
	defer rows.Close()
	defer marketRepository.Db.Close()
	return parseRows(rows, error)
}

func (marketRepository *MarketRepository)GetMarket(marketId int64) (model.MarketExportable) {
	defer marketRepository.Db.Close()
	stmt, error := marketRepository.Db.Prepare("SELECT id,user_id,name,description,startdate,lat, lon, market_type, flexible, place FROM market WHERE active = 1 and id = ?")
	if error != nil{
		log.Println(error)
		return model.MarketExportable{}
	}
	row := stmt.QueryRow(marketId)
	var market model.MarketExportable
	row.Scan(&market.Id, &market.UserId, &market.Name, &market.Description, &market.Date, &market.Lat, &market.Lon, &market.Type, &market.Flexible, &market.Place)
	return market
}

func (marketRepository *MarketRepository)GetMarkets(marketFilter model.MarketFilter) ([]model.MarketExportable) {
	if marketFilter.Radio > EARTH_RATE {
		marketFilter.Radio = EARTH_RATE
	}

	if marketFilter.Radio < 0 {
		marketFilter.Radio = 0
	}

	maxlat := marketFilter.Lat + rad2deg(marketFilter.Radio/EARTH_RATE)
	minlat := marketFilter.Lat - rad2deg(marketFilter.Radio/EARTH_RATE)
	maxlon := marketFilter.Lon + rad2deg(math.Asin(marketFilter.Radio/EARTH_RATE) / math.Cos(deg2rad(marketFilter.Lat)))
	minlon := marketFilter.Lon - rad2deg(math.Asin(marketFilter.Radio/EARTH_RATE) / math.Cos(deg2rad(marketFilter.Lat)))

	stmt, error := 	marketRepository.Db.Prepare(`
		SELECT
			id,user_id,name,description,startdate,lat, lon, market_type, flexible, place
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
	if error != nil{
		log.Println(error)
		return []model.MarketExportable{}
	}
	defer stmt.Close()
	defer marketRepository.Db.Close()

	if error != nil {
		log.Println(error)
	}

	rows , error := stmt.Query(
		maxlat,
		minlat,
		maxlon,
		minlon,
		marketFilter.Page * MAX_RESULTS,
		MAX_RESULTS)

	if error != nil {
		log.Println(error)
	}

	return parseRows(rows, error)
}

func deg2rad(deg float64) float64 {
	return deg * math.Pi / RADIO
}
func rad2deg(rad float64) float64 {
	return rad * RADIO / math.Pi
}

func (marketRepository *MarketRepository)Create(market model.MarketExportable) int64{
	stmt, error := 	marketRepository.Db.Prepare(
		`
		INSERT INTO market
		(id,name,description,startdate,lat,lon,active,user_id, market_type, flexible, place)
		VALUES
		(null,?,?,?,?,?,?,?,?,?,?)`)
	if error != nil{
		log.Println(error)
		return -1
	}
	 result , error := stmt.Exec(
		market.Name,
		market.Description,
		market.Date,
		market.Lat,
		market.Lon,
		true,
		market.UserId,
		market.Type,
		market.Flexible,
		market.Place,
		)
	if error != nil{
		log.Println(error)
		return -1
	}
	defer stmt.Close()
	defer marketRepository.Db.Close()

 	if error != nil{
 		log.Println(error)
	}

	lastInsertedId, _:= result.LastInsertId()

	return lastInsertedId
}

func (marketRepository *MarketRepository) Edit(market model.MarketExportable) (bool) {
	stmt, error := 	marketRepository.Db.Prepare(
		`
		UPDATE market SET
		name = ?, description = ? , startdate = ?,lat = ?,lon = ?, market_type = ?, flexible = ?, place = ?
		WHERE id = ?`)

	if error != nil{
		log.Println(error)
		return false
	}

	defer stmt.Close()
	defer marketRepository.Db.Close()
	if error != nil{
		log.Println(error)
		return false
	}
	_, error = stmt.Exec(
		market.Name,
		market.Description,
		market.Date,
		market.Lat,
		market.Lon,
		market.Id,
		market.Type,
		market.Flexible,
		market.Place,
		)

	if error != nil{
		log.Println(error)
	}
	return error == nil
}

func (marketRepository *MarketRepository) Repeat(market model.MarketExportable, newDate string) (int64) {
	market.Date = newDate

	return marketRepository.Create(market)
}

func (marketRepository *MarketRepository) Delete(userId int64, marketId int64) (bool) {
	stmt, error := 	marketRepository.Db.Prepare(
		`
		DELETE market
		WHERE id = ? and user_id = ?`)

	defer stmt.Close()
	defer marketRepository.Db.Close()

	_, error = stmt.Exec(marketId,userId)

	if error != nil{
		log.Println(error)
	}
	return error == nil
}

func parseRows(rows *sql.Rows, error error) []model.MarketExportable {
	if error == nil {
		var markets []model.MarketExportable
		for rows.Next() {
			market, err := parseRow(rows)
			if err != nil {
				log.Println(error)
			} else {
				markets = append(markets, market)
			}
		}
		return markets
	} else{
		log.Println(error)
		return nil
	}
}
func parseRow(rows *sql.Rows) (model.MarketExportable, error) {
	var market model.Market
	err := rows.Scan(&market.Id, &market.UserId, &market.Name, &market.Description, &market.Date, &market.Lat, &market.Lon, &market.Type, &market.Flexible, &market.Place)
	if err != nil{
		log.Println(err)
	}

	marketExportable := model.MarketExportable{
		Id :market.Id,
		UserId : market.UserId,
		Name: market.Name,
		Description : market.Description,
		Type: market.Type.String,
		Flexible : market.Flexible.Bool,
		Place : market.Place.String,
		Date : market.Date,
		Lat : market.Lat,
		Lon : market.Lon,
	}

	return marketExportable, err
}

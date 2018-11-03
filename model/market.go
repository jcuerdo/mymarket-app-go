package model

import "database/sql"

type Market struct {
	Id          int64          `json:"id"`
	UserId      int            `json:"user_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        sql.NullString `json:"type"`
	Flexible    sql.NullBool `json:"flexible"`
	Place       sql.NullString `json:"place"`
	Date        string         `json:"startdate"`
	Lat         float32        `json:"lat"`
	Lon         float32        `json:"lon"`
}

type MarketExportable struct {
	Id          int64   `json:"id"`
	UserId      int     `json:"user_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Flexible    bool    `json:"flexible"`
	Place       string  `json:"place"`
	Date        string  `json:"startdate"`
	Lat         float32 `json:"lat"`
	Lon         float32 `json:"lon"`
}

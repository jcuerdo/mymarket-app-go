package model

type Comment struct {
	Id       int    `json:"id"`
	MarketId int    `json:"market_id"`
	UserId   int    `json:"user_id"`
	Content  string `json:"content"`
}

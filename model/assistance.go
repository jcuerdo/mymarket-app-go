package model

type Assistance struct {
	Id       int    `json:"id"`
	MarketId int    `json:"market_id"`
	UserId   int    `json:"user_id"`
}

type AssistanceResult struct {
	Id       int    `json:"id"`
	MarketId int    `json:"market_id"`
	User     UserAssistance    `json:"user_id"`
	Date  string `json:"date"`
}
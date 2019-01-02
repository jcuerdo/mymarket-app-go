package model

type MarketFilter struct {
	Page	int64 	`json:"page"`
	Radio	float64 `json:"radio"`
	Lat		float64 `json:"lat"`
	Lon		float64 `json:"lon"`
	Privacy	string  `json:"privacy"`
}
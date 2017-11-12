package model

type Market struct {
	Id           int 	`json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Date		 string `json:"date"`
	Lat		 	 float32 `json:"lat"`
	Lon		 	 float32 `json:"lon"`
}
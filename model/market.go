package model

type Market struct {
	Id           int64   `json:"id"`
	UserId       int 	 `json:"user_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Type         string  `json:"type"`
	Flexible     string  `json:"flexible"`
	Place        string  `json:"place"`
	Date		 string  `json:"startdate"`
	Lat		 	 float32 `json:"lat"`
	Lon		 	 float32 `json:"lon"`
}
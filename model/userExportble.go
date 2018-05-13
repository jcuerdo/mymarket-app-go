package model

type UserExportable struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"fullname"`
	Photo       string `json:"photo"`
	Description string `json:"description"`
	Role        string `json:"role"`
}

package service

import (
	"github.com/jcuerdo/mymarket-app-go/model"
)



type Email struct {
}

func NewEmailNotificator() (Notificator){
	return Email{}
}

func (email Email) NotifyComment(ids []model.UserToken, comment model.Comment) (bool) {

	return true
}

func (email Email) NotifyAssistance(ids []model.UserToken, assistance model.Assistance) (bool) {

	return true
}

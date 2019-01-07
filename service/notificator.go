package service

import "github.com/jcuerdo/mymarket-app-go/model"

type Notificator interface {
	NotifyComment(users []model.UserToken, comment model.Comment) (bool)
}

type NotificatorService struct {
	Notificators []Notificator
}

func NewNotificatorService() (NotificatorService){
	return NotificatorService{
		Notificators : []Notificator{NewFireBaseNotificator(), NewEmailNotificator()},
	}
}

func (notificator *NotificatorService) NotifyCommentToAll(users []model.UserToken, comment model.Comment) (bool) {
	for _, notificator := range notificator.Notificators {
		notificator.NotifyComment(users, comment)
	}
	return true
}


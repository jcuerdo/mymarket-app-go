package service

import "github.com/jcuerdo/mymarket-app-go/model"

type FireBase struct {
}

func (firebase *FireBase) notifyComment(ids []string, comment model.Comment) (bool) {
	return true;
}


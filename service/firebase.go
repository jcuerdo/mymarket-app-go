package service

import (
	"github.com/jcuerdo/mymarket-app-go/model"
	"github.com/NaySoftware/go-fcm"
	"log"
	"strconv"
)

var serverKey = "AAAAykf79Vw:APA91bHxXCh8mee1m8ycjrVbF8PsnewZe0ZF5r3DLMpyyTstOrbYo5lIXxX4e9GX1VqBTM8vTMw6o_I3yZymXwnM_MINecGFSaBd65ar7qKDH5-4KDFseu6eHExVbl48WKfgQvGt0GIE"


type FireBase struct {
	serverKey string
}

func NewFireBaseNotificator() (Notificator){
	return FireBase{serverKey}
}

func (fireBase FireBase) NotifyComment(users []model.UserToken, comment model.Comment) (bool) {
	data := map[string]string{
		"msg": comment.Content,
		"marketID" : strconv.Itoa(comment.MarketId),
	}

	var ids = make([]string, 0)
	for _, user := range users {
		ids = append(ids, user.FirebaseToken.String)
	}

	firebaseClient := fcm.NewFcmClient(fireBase.serverKey)
	firebaseClient.NewFcmRegIdsMsg(ids, data)
	notificationPayload := &fcm.NotificationPayload{Title:comment.Content,Sound:"default",Color:"#0000FF"}
	firebaseClient.SetNotificationPayload(notificationPayload)
	resp, error := firebaseClient.Send()

	if error != nil {
		log.Println(error)
	}

	resp.PrintResults()

	return true
}


func (fireBase FireBase) NotifyAssistance(users []model.UserToken, assistance model.Assistance) (bool) {
	data := map[string]string{
		"marketID" : strconv.Itoa(assistance.MarketId),
	}

	var ids = make([]string, 0)
	for _, user := range users {
		ids = append(ids, user.FirebaseToken.String)
	}

	firebaseClient := fcm.NewFcmClient(fireBase.serverKey)
	firebaseClient.NewFcmRegIdsMsg(ids, data)
	notificationPayload := &fcm.NotificationPayload{Title:"Someone is going to the market", Sound:"default",Color:"#0000FF"}
	firebaseClient.SetNotificationPayload(notificationPayload)
	resp, error := firebaseClient.Send()

	if error != nil {
		log.Println(error)
	}

	resp.PrintResults()

	return true
}

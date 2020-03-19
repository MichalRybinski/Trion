package models

import (
	"time"
)

type AuthModel struct {
	UserID		string		`json:"uid" bson:"uid"`
	UUID			string		`json:"uuid" bson:"uuid"`
	Token			string		`json:"token" bson:"token"`
	CreatedAt time.Time	`json:"createdAt" bson:"createdAt"`
	UpdatedAt	time.Time	`json:"updatedAt" bson:"updatedAt"`
}

func NewAuth(uid string, uuid string, token string) AuthModel {
	now:=time.Now()
	return AuthModel{
		UserID: uid,
		UUID : uuid,
		Token : token,
		CreatedAt : now,
		UpdatedAt : now,
	}
}
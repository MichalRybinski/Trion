package models

import (
	//"encoding/json"
	"time"
)

type UserDBModel struct {
	Login			string		`json:"login" bson:"login"`
	Hash			string		`json:"hash" bson:"hash"`
	CreatedAt time.Time	`json:"createdAt" bson:"createdAt"`
	UpdatedAt	time.Time	`json:"updatedAt" bson:"updatedAt"`
}
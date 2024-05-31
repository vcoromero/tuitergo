package models

import "time"

type BodyRequestCreate struct {
	Message string `bson:"message" json:"message,omitempty"`
}

type Tuit struct {
	UserId     string    `bson:"user_id" json:"user_id,omitempty"`
	Message    string    `bson:"message" json:"message,omitempty"`
	Created_at time.Time `bson:"created_at" json:"created_at,omitempty"`
}

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BodyRequestCreate struct {
	Message string `bson:"message" json:"message,omitempty"`
}

type Tuit struct {
	UserId     string    `bson:"user_id" json:"user_id,omitempty"`
	Message    string    `bson:"message" json:"message,omitempty"`
	Created_at time.Time `bson:"created_at" json:"created_at,omitempty"`
}

type GetTuitsFromUser struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	UserId     string             `bson:"user_id" json:"user_id,omitempty"`
	Message    string             `bson:"message" json:"message,omitempty"`
	Created_at time.Time          `bson:"created_at" json:"created_at,omitempty"`
}

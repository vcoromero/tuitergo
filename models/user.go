package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name,omitempty"`
	LastNames string             `bson:"lastnames" json:"lastnames,omitempty"`
	Birthdate time.Time          `bson:"birthdate" json:"birthdate,omitempty"`
	Email     string             `bson:"email" json:"email"`
	Password  time.Time          `bson:"password" json:"password,omitempty"`
	Avatar    time.Time          `bson:"avatar" json:"avatar,omitempty"`
	Banner    time.Time          `bson:"banner" json:"banner,omitempty"`
	Biography time.Time          `bson:"biography" json:"biography,omitempty"`
	Location  time.Time          `bson:"location" json:"location,omitempty"`
	WebSite   time.Time          `bson:"website" json:"website,omitempty"`
}

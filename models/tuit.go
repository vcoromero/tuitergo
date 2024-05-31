package models

type Tuit struct {
	Message string `bson:"message" json:"message,omitempty"`
}

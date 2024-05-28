package db

import (
	"context"
	"time"

	"github.com/vcoromero/tuitergo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckedIfUserExist(email string) (models.User, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("users")
	condition := bson.M{"email": email}

	var result models.User

	err := col.FindOne(ctx, condition).Decode(&result)
	ID := ""
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, false, ID
		}
		return result, false, ID
	}
	ID = result.ID.Hex()
	return result, true, ID
}

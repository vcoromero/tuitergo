package db

import (
	"context"

	"github.com/vcoromero/tuitergo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProfile(id string) (models.User, error) {
	ctx := context.TODO()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("users")

	var profile models.User

	objId, _ := primitive.ObjectIDFromHex(id)

	condition := bson.M{
		"_id": objId,
	}

	err := col.FindOne(ctx, condition).Decode(&profile)
	profile.Password = ""
	if err != nil {
		return profile, err
	}
	return profile, err
}

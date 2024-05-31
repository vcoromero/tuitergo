package db

import (
	"context"

	"github.com/vcoromero/tuitergo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertTuit(t models.SaveTuit) (string, bool, error) {
	ctx := context.TODO()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("tuits")

	register := bson.M{
		"user_id":    t.UserId,
		"message":    t.Message,
		"created_at": t.Created_at,
	}

	result, err := col.InsertOne(ctx, register)
	if err != nil {
		return "", false, err
	}

	ObjtId, _ := result.InsertedID.(primitive.ObjectID)

	return ObjtId.String(), true, nil
}

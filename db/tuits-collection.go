package db

import (
	"context"

	"github.com/vcoromero/tuitergo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertTuit(t models.Tuit) (string, bool, error) {
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

func GetTuitsFromUser(id string, page int64) ([]*models.GetTuitsFromUser, bool) {
	ctx := context.TODO()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("tuits")

	var results []*models.GetTuitsFromUser

	condition := bson.M{
		"user_id": id,
	}

	options := options.Find()
	options.SetLimit(20)
	options.SetSort(bson.D{{Key: "created_at", Value: -1}})
	options.SetSkip((page - 1) * 20)

	cursor, err := col.Find(ctx, condition, options)

	if err != nil {
		return results, false
	}

	for cursor.Next(ctx) {

		var row models.GetTuitsFromUser
		err := cursor.Decode(&row)
		if err != nil {
			return results, false
		}

		results = append(results, &row)
	}

	return results, true
}

func DeleteTuit(id string, user_id string) error {
	ctx := context.TODO()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("tuits")

	objId, _ := primitive.ObjectIDFromHex(id)

	condition := bson.M{
		"_id":     objId,
		"user_id": user_id,
	}

	_, err := col.DeleteOne(ctx, condition)
	return err
}

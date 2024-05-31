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

func GetTuitsFromFollowers(id string, page int64) ([]models.GetTuitsFromFollowers, bool) {
	ctx := context.TODO()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("relationships")

	skip := (page - 1) * 20

	conditions := make([]bson.M, 0)
	conditions = append(conditions, bson.M{"$match": bson.M{"user_id": id}})
	conditions = append(conditions, bson.M{
		"$lookup": bson.M{
			"from":         "tuits",
			"localField":   "user_relationship_id",
			"foreignField": "user_id",
			"as":           "tuit",
		},
	})
	conditions = append(conditions, bson.M{"$unwind": "$tuit"})
	conditions = append(conditions, bson.M{"$sort": bson.M{"tuit.created_at": -1}})
	conditions = append(conditions, bson.M{"$skip": skip})
	conditions = append(conditions, bson.M{"$limit": 20})

	var result []models.GetTuitsFromFollowers
	cursor, err := col.Aggregate(ctx, conditions)
	if err != nil {
		return result, false
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return result, false
	}
	return result, true
}

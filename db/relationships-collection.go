package db

import (
	"context"

	"github.com/vcoromero/tuitergo/models"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateRelationship(t models.Relationship) (bool, error) {
	ctx := context.TODO()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("relationships")

	_, err := col.InsertOne(ctx, t)

	if err != nil {
		return false, err
	}

	return true, nil
}

func GetRelationship(t models.Relationship) bool {
	ctx := context.TODO()
	db := MongoCN.Database(DatabaseName)
	col := db.Collection("relationships")

	condition := bson.M{
		"user_id":              t.UserId,
		"user_relationship_id": t.UserRelationshipId,
	}

	var response models.Relationship

	err := col.FindOne(ctx, condition).Decode(&response)

	return err == nil
}

func DeleteRelationship(t models.Relationship) (bool, error) {
	ctx := context.TODO()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("relationships")

	_, err := col.DeleteOne(ctx, t)

	if err != nil {
		return false, err
	}

	return true, nil
}

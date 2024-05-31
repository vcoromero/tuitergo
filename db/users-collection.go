package db

import (
	"context"
	"fmt"
	"time"

	"github.com/vcoromero/tuitergo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func Login(email string, password string) (models.User, bool) {
	user, found, _ := CheckedIfUserExist(email)
	if !found {
		return user, false
	}

	passwordBytes := []byte(password)
	passwordDB := []byte(user.Password)

	err := bcrypt.CompareHashAndPassword(passwordDB, passwordBytes)
	if err != nil {
		return user, false
	}
	return user, true
}

func GetUser(id string) (models.User, error) {
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

func CreateUser(u models.User) (string, bool, error) {
	ctx := context.TODO()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("users")

	u.Password, _ = EncryptPassword(u.Password)

	result, err := col.InsertOne(ctx, u)
	if err != nil {
		return "", false, err
	}

	ObjtId, _ := result.InsertedID.(primitive.ObjectID)

	return ObjtId.String(), true, nil
}
func GetUsers(id string, page int64, search string, usertype string) ([]*models.User, bool) {
	ctx := context.TODO()
	db := MongoCN.Database(DatabaseName)
	col := db.Collection("users")

	var results []*models.User

	options := options.Find()
	options.SetLimit(20)
	options.SetSkip((page - 1) * 20)

	query := bson.M{
		"name": bson.M{"$regex": `(?i)` + search},
	}

	cursor, err := col.Find(ctx, query, options)
	if err != nil {
		return results, false
	}

	var join bool

	for cursor.Next(ctx) {
		var row models.User

		err := cursor.Decode(&row)
		if err != nil {
			fmt.Println("Decode = " + err.Error())
			return results, false
		}

		var r models.Relationship
		r.UserId = id
		r.UserRelationshipId = row.ID.Hex()

		join = false
		found := GetRelationship(r)

		if usertype == "new" && !found {
			join = true
		}
		if usertype == "follow" && found {
			join = true
		}

		if r.UserRelationshipId == id {
			join = false
		}

		if join {
			row.Password = ""
			results = append(results, &row)
		}
	}

	err = cursor.Err()

	if err != nil {
		fmt.Println("cur.Err() = " + err.Error())
		return results, false
	}

	cursor.Close(ctx)

	return results, true
}

func UpdateUser(u models.User, id string) (bool, error) {
	ctx := context.TODO()
	db := MongoCN.Database(DatabaseName)
	col := db.Collection("users")

	register := make(map[string]interface{})
	if len(u.Name) > 0 {
		register["name"] = u.Name
	}
	if len(u.LastNames) > 0 {
		register["lastnames"] = u.LastNames
	}
	register["birthdate"] = u.Birthdate
	if len(u.Avatar) > 0 {
		register["avatar"] = u.Avatar
	}
	if len(u.Banner) > 0 {
		register["banner"] = u.Banner
	}
	if len(u.Biography) > 0 {
		register["biography"] = u.Biography
	}
	if len(u.Location) > 0 {
		register["location"] = u.Location
	}
	if len(u.WebSite) > 0 {
		register["website"] = u.WebSite
	}

	updateString := bson.M{
		"$set": register,
	}

	objId, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{
		"_id": bson.M{"$eq": objId},
	}

	_, err := col.UpdateOne(ctx, filter, updateString)

	if err != nil {
		return false, err
	}

	return true, nil
}

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

func EncryptPassword(psswrd string) (string, error) {
	cost := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(psswrd), cost)
	if err != nil {
		return err.Error(), err
	}
	return string(bytes), nil
}

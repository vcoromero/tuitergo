package routers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vcoromero/tuitergo/db"
	"github.com/vcoromero/tuitergo/models"
)

func Register(ctx context.Context) models.ResponseAPI {
	var t models.User
	var r models.ResponseAPI
	r.Status = 400
	fmt.Println("Registering user")

	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		r.Message = "Error parsing request body: " + err.Error()
		fmt.Println(r.Message)
		return r
	}

	if len(t.Email) == 0 {
		r.Message = "Must have an Email"
		fmt.Println(r.Message)
		return r
	}

	if len(t.Password) < 6 {
		r.Message = "Password must be at least 6 characters long"
		fmt.Println(r.Message)
		return r
	}

	_, found, _ := db.CheckedIfUserExist(t.Email)
	if found {
		r.Message = "There is an user registered with this email"
		fmt.Println(r.Message)
		return r
	}

	_, status, err := db.InsertUser(t)
	if err != nil {
		r.Message = "An error occurred trying to insert the user: " + err.Error()
		fmt.Println(r.Message)
		return r
	}

	if !status {
		r.Message = "The user could not be inserted"
		fmt.Println(r.Message)
		return r
	}

	r.Status = 200
	r.Message = "User registered!"
	fmt.Println(r.Message)
	return r
}

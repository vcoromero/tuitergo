package routers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vcoromero/tuitergo/db"
	"github.com/vcoromero/tuitergo/models"
)

func UpdateProfile(ctx context.Context, claim models.Claim) models.ResponseAPI {
	var r models.ResponseAPI
	r.Status = 400
	fmt.Println("Entered to show profile")

	var t models.User

	body := ctx.Value(models.Key("body")).(string)

	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		r.Message = "Wrong data " + err.Error()
		return r
	}

	status, err := db.UpdateUser(t, claim.ID.Hex())
	if err != nil {
		r.Message = "Ocurred error tryng to update user" + err.Error()
		return r
	}
	if !status {
		r.Message = "Cannot update user"
		return r
	}

	r.Status = 200
	r.Message = "Update user succesfully"
	return r
}

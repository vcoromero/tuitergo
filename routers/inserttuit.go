package routers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/vcoromero/tuitergo/db"
	"github.com/vcoromero/tuitergo/models"
)

func InserTuit(ctx context.Context, claim models.Claim) models.ResponseAPI {
	var t models.Tuit
	var r models.ResponseAPI
	r.Status = 400
	fmt.Println("Registering tuit")

	user_id := claim.ID.Hex()
	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		r.Message = "Error parsing request body: " + err.Error()
		fmt.Println(r.Message)
		return r
	}

	register := models.SaveTuit{
		UserId:     user_id,
		Message:    t.Message,
		Created_at: time.Now(),
	}

	_, status, err := db.InsertTuit(register)
	if err != nil {
		r.Message = "An error occurred trying to insert tuit: " + err.Error()
		fmt.Println(r.Message)
		return r
	}

	if !status {
		r.Message = "Tuit could not be inserted"
		fmt.Println(r.Message)
		return r
	}

	r.Status = 200
	r.Message = "Tuit registered!"
	fmt.Println(r.Message)
	return r
}

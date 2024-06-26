package routers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/vcoromero/tuitergo/db"
	"github.com/vcoromero/tuitergo/models"
)

func CreateTuit(ctx context.Context, claim models.Claim) models.ResponseAPI {
	var t models.BodyRequestCreate
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

	register := models.Tuit{
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

func GetTuitsFromUser(request events.APIGatewayProxyRequest) models.ResponseAPI {
	var r models.ResponseAPI
	r.Status = 200
	fmt.Println("Entered to get tuits from user")

	ID := request.QueryStringParameters["id"]
	bodyPage := request.QueryStringParameters["page"]

	if len(ID) < 1 {
		r.Message = "id parameter is required"
		return r
	}
	if len(bodyPage) < 1 {
		bodyPage = "1"
	}

	page, err := strconv.Atoi(bodyPage)
	if err != nil {
		r.Message = "page parameter must be a number"
		return r
	}

	tuits, correct := db.GetTuitsFromUser(ID, int64(page))
	if !correct {
		r.Message = "Error occured trying to get tuits from user"
		return r
	}

	resJson, err := json.Marshal(tuits)
	if err != nil {
		r.Status = 500
		r.Message = "Error trying to parse the tuits data to json" + err.Error()
		return r
	}

	r.Status = 200
	r.Message = string(resJson)
	return r
}

func DeleteTuit(request events.APIGatewayProxyRequest, claim models.Claim) models.ResponseAPI {
	var r models.ResponseAPI
	r.Status = 200
	fmt.Println("Entered to get tuits from user")

	ID := request.QueryStringParameters["id"]

	if len(ID) < 1 {
		r.Message = "id parameter is required"
		return r
	}

	err := db.DeleteTuit(ID, claim.ID.Hex())
	if err != nil {
		r.Message = "Error occured trying to delete tuit " + err.Error()
		return r
	}

	r.Status = 200
	r.Message = "Tuit deleted!!"
	return r
}

func GetTuitsFromFollowers(request events.APIGatewayProxyRequest, claim models.Claim) models.ResponseAPI {
	var r models.ResponseAPI
	r.Status = 400

	user_id := claim.ID.Hex()

	bodyPage := request.QueryStringParameters["page"]
	if len(bodyPage) < 1 {
		bodyPage = "1"
	}

	page, err := strconv.Atoi(bodyPage)
	if err != nil {
		r.Message = "page parameter must be a number"
		return r
	}

	tuits, correct := db.GetTuitsFromFollowers(user_id, int64(page))
	if !correct {
		r.Message = "Error occured trying to get tuits from followers"
		return r
	}

	resJson, err := json.Marshal(tuits)
	if err != nil {
		r.Status = 500
		r.Message = "Error trying to parse the tuits data to json: " + err.Error()
		return r
	}

	r.Status = 200
	r.Message = string(resJson)
	fmt.Println("Response:", r) // Mensaje de depuración
	return r
}

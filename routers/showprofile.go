package routers

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/vcoromero/tuitergo/db"
	"github.com/vcoromero/tuitergo/models"
)

func ShowProfile(request events.APIGatewayProxyRequest) models.ResponseAPI {
	var r models.ResponseAPI
	r.Status = 200
	fmt.Println("Entered to show profile")
	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		r.Message = "id parameter is required"
		return r
	}

	profile, err := db.GetProfile(ID)
	if err != nil {
		r.Message = "Ocurred error to find profile" + err.Error()
		return r
	}

	resJson, err := json.Marshal(profile)
	if err != nil {
		r.Status = 500
		r.Message = "Error trying to parse the user data to json" + err.Error()
		return r
	}

	r.Status = 200
	r.Message = string(resJson)
	return r
}

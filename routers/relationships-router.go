package routers

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/vcoromero/tuitergo/db"
	"github.com/vcoromero/tuitergo/models"
)

func CreateRelationship(ctx context.Context, request events.APIGatewayProxyRequest, claim models.Claim) models.ResponseAPI {
	var r models.ResponseAPI
	r.Status = 400

	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		r.Message = "id parameter is required"
		return r
	}

	var t models.Relationship
	t.UserId = claim.ID.Hex()
	t.UserRelationshipId = ID

	status, err := db.CreateRelationship(t)

	if err != nil {
		r.Message = "Error occurred trying to create relationship " + err.Error()
		return r
	}

	if !status {
		r.Message = "Failed to create relationship "
		return r
	}

	r.Status = 200
	r.Message = "Relationship created!!"
	return r
}

func GetRelationship(request events.APIGatewayProxyRequest, claim models.Claim) models.ResponseAPI {
	var r models.ResponseAPI

	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		r.Message = "id parameter is required"
		return r
	}

	var t models.Relationship
	t.UserId = claim.ID.Hex()
	t.UserRelationshipId = ID

	var resp models.ReponseGetRelationship
	found := db.GetRelationship(t)

	if !found {
		resp.Status = !found
	} else {
		resp.Status = found
	}

	respJson, err := json.Marshal(found)
	if err != nil {
		r.Status = 500
		r.Message = "Error occurred trying to parse json " + err.Error()
		return r
	}

	r.Status = 200
	r.Message = string(respJson)

	return r
}

func DeleteRelationship(request events.APIGatewayProxyRequest, claim models.Claim) models.ResponseAPI {
	var r models.ResponseAPI
	r.Status = 400

	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		r.Message = "id parameter is required"
		return r
	}

	var t models.Relationship
	t.UserId = claim.ID.Hex()
	t.UserRelationshipId = ID

	status, err := db.DeleteRelationship(t)

	if err != nil {
		r.Message = "Error occurred trying to delete relationship " + err.Error()
		return r
	}

	if !status {
		r.Message = "Failed to delete relationship "
		return r
	}

	r.Status = 200
	r.Message = "Relationship deleted!!"
	return r
}

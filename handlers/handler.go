package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/vcoromero/tuitergo/jwt"
	"github.com/vcoromero/tuitergo/models"
	"github.com/vcoromero/tuitergo/routers"
)

func Handlers(ctx context.Context, request events.APIGatewayProxyRequest) models.ResponseAPI {
	path := ctx.Value(models.Key("path")).(string)
	method := ctx.Value(models.Key("method")).(string)
	fmt.Println("Voy a procesar", path, ">", method)

	var res models.ResponseAPI
	res.Status = 400

	isOk, statusCode, msg, claim := authorizationValidate(ctx, request)
	if !isOk {
		fmt.Println("Authorization failed")
		res.Status = statusCode
		res.Message = msg
		return res
	}

	switch method {
	case "POST":
		switch path {
		case "create-user":
			return routers.CreateUser(ctx)
		case "login":
			return routers.Login(ctx)
		case "create-tuit":
			return routers.CreateTuit(ctx, claim)
		case "upload-avatar":
			return routers.UploadImage(ctx, "A", request, claim)
		case "upload-banner":
			return routers.UploadImage(ctx, "B", request, claim)
		case "create-relationship":
			return routers.CreateRelationship(ctx, request, claim)
		default:
			fmt.Println("Unknown POST route")
		}
	case "GET":
		switch path {
		case "get-user":
			return routers.GetUser(request)
		case "get-tuits-from-user":
			return routers.GetTuitsFromUser(request)
		case "get-avatar-from-user":
			return routers.GetImage(ctx, "A", request, claim)
		case "get-banner-from-user":
			return routers.GetImage(ctx, "B", request, claim)
		case "get-relationship":
			return routers.GetRelationship(request, claim)
		case "get-users":
			return routers.GetUsers(request, claim)
		default:
			fmt.Println("Unknown GET route")
		}
	case "PUT":
		switch path {
		case "update-user":
			return routers.UpdateUser(ctx, claim)
		default:
			fmt.Println("Unknown PUT route")
		}
	case "DELETE":
		switch path {
		case "delete-tuit":
			return routers.DeleteTuit(request, claim)
		case "delete-relationship":
			return routers.DeleteRelationship(request, claim)
		default:
			fmt.Println("Unknown DELETE route")
		}
	default:
		fmt.Println("Invalid method")
	}
	res.Message = "Message invalid"
	return res
}

func authorizationValidate(ctx context.Context, request events.APIGatewayProxyRequest) (bool, int, string, models.Claim) {
	path := ctx.Value(models.Key("path")).(string)
	fmt.Println("Authorization validation for path:", path)
	if path == "create-user" || path == "login" || path == "getAvatar" || path == "getBanner" {
		fmt.Println("Public endpoint, no authorization required")
		return true, 200, "", models.Claim{}
	}

	token := request.Headers["Authorization"]
	fmt.Println("Authorization header:", token)
	if len(token) == 0 {
		fmt.Println("Authorization token required but not found")
		return false, 401, "required token", models.Claim{}
	}

	claim, isOK, msg, err := jwt.ProcessToken(token, ctx.Value(models.Key("jwtSign")).(string))
	if !isOK {
		if err != nil {
			fmt.Println("Token error: " + err.Error())
			return false, 401, err.Error(), models.Claim{}
		} else {
			fmt.Println("Token invalid: " + msg)
			return false, 401, msg, models.Claim{}
		}
	}
	fmt.Println("Token valid")
	return true, 200, msg, *claim
}

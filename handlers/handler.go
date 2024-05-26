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
	fmt.Println("Voy a procesar " + ctx.Value(models.Key("path")).(string) + " > " + ctx.Value(models.Key("method")).(string))

	var res models.ResponseAPI
	res.Status = 400

	isOk, statusCode, msg, _ := authorizationValidate(ctx, request)
	if !isOk {
		res.Status = statusCode
		res.Message = msg
		return res
	}

	switch ctx.Value(models.Key("method")).(string) {
	case "POST":
		switch ctx.Value(models.Key("path")).(string) {
		case "register":
			return routers.Register(ctx)
		}
	case "GET":
		switch ctx.Value(models.Key("path")).(string) {
		//
		}
	case "PUT":
		switch ctx.Value(models.Key("path")).(string) {
		//
		}
	case "DELETE":
		switch ctx.Value(models.Key("path")).(string) {
		//
		}
	}
	res.Message = "Message invalid"
	return res
}

func authorizationValidate(ctx context.Context, request events.APIGatewayProxyRequest) (bool, int, string, models.Claim) {
	path := ctx.Value(models.Key("path")).(string)
	if path == "register" || path == "login" || path == "getAvatar" || path == "getBanner" {
		return true, 200, "", models.Claim{}
	}

	token := request.Headers["Authorization"]
	if len(token) == 0 {
		return false, 401, "required token", models.Claim{}
	}

	claim, isOK, msg, err := jwt.ProcessToken(token, ctx.Value(models.Key("jwtSign")).(string))
	if !isOK {
		if err != nil {
			fmt.Println("token error: " + err.Error())
			return false, 401, err.Error(), models.Claim{}
		} else {
			fmt.Println("token error: " + err.Error())
			return false, 401, err.Error(), models.Claim{}
		}
	}
	fmt.Println("token ok!")
	return true, 200, msg, *claim
}

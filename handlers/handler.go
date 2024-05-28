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

	isOk, statusCode, msg, _ := authorizationValidate(ctx, request)
	if !isOk {
		fmt.Println("Authorization failed")
		res.Status = statusCode
		res.Message = msg
		return res
	}

	switch method {
	case "POST":
		fmt.Println("Processing POST request")
		switch path {
		case "register":
			fmt.Println("Handling register route")
			return routers.Register(ctx)
		case "login":
			return routers.Login(ctx)
		default:
			fmt.Println("Unknown POST route")
		}
	case "GET":
		fmt.Println("Processing GET request")
		switch path {
		// Añade tus rutas GET aquí
		default:
			fmt.Println("Unknown GET route")
		}
	case "PUT":
		fmt.Println("Processing PUT request")
		switch path {
		// Añade tus rutas PUT aquí
		default:
			fmt.Println("Unknown PUT route")
		}
	case "DELETE":
		fmt.Println("Processing DELETE request")
		switch path {
		// Añade tus rutas DELETE aquí
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
	if path == "register" || path == "login" || path == "getAvatar" || path == "getBanner" {
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

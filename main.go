package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vcoromero/tuitergo/awsgo"
	"github.com/vcoromero/tuitergo/db"
	"github.com/vcoromero/tuitergo/handlers"
	"github.com/vcoromero/tuitergo/models"
	"github.com/vcoromero/tuitergo/secretmanager"
)

func main() {
	lambda.Start(CallLambda)
}

func CallLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse

	awsgo.InitializeAWS()

	if !ValidatedParameters() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en las variables de entorno, deben incluir 'SecretName', 'BucketName' y 'UrlPrefix'",
			Headers: map[string]string{
				"Content-type": "application/json",
			},
		}
		return res, nil
	}

	SecretModel, err := secretmanager.GetSecret(os.Getenv("SecretName"))
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en la lectura de Secret " + err.Error(),
			Headers: map[string]string{
				"Content-type": "application/json",
			},
		}
		return res, nil
	}

	path := strings.Replace(request.PathParameters["tuitergo"], os.Getenv("UrlPrefix"), "", -1)
	fmt.Println("Resolved path:", path)
	ctx = context.WithValue(ctx, models.Key("path"), path)
	ctx = context.WithValue(ctx, models.Key("method"), request.HTTPMethod)
	ctx = context.WithValue(ctx, models.Key("user"), SecretModel.Username)
	ctx = context.WithValue(ctx, models.Key("password"), SecretModel.Password)
	ctx = context.WithValue(ctx, models.Key("host"), SecretModel.Host)
	ctx = context.WithValue(ctx, models.Key("database"), SecretModel.Database)
	ctx = context.WithValue(ctx, models.Key("jwtSign"), SecretModel.JWTSign)
	ctx = context.WithValue(ctx, models.Key("body"), request.Body)
	ctx = context.WithValue(ctx, models.Key("bucketName"), os.Getenv("BucketName"))

	// Check database connection
	err = db.ConnectDB(ctx)
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error conectando en la db " + err.Error(),
			Headers: map[string]string{
				"Content-type": "application/json",
			},
		}
		return res, nil
	}

	resAPI := handlers.Handlers(ctx, request)
	if resAPI.CustomResponse == nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: resAPI.Status,
			Body:       resAPI.Message,
			Headers: map[string]string{
				"Content-type": "application/json",
			},
		}
		return res, nil
	} else {
		return resAPI.CustomResponse, nil
	}
}

func ValidatedParameters() bool {
	_, gotParameter := os.LookupEnv("SecretName")
	if !gotParameter {
		return gotParameter
	}

	_, gotParameter = os.LookupEnv("BucketName")
	if !gotParameter {
		return gotParameter
	}

	_, gotParameter = os.LookupEnv("UrlPrefix")
	if !gotParameter {
		return gotParameter
	}
	return gotParameter
}

package main

import (
	"context"
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
	// Initialize AWS
	awsgo.InitializeAWS()

	// Create unique email index
	db.CreateUniqueEmailIndex()

	// Start Lambda
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

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("user"), SecretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("password"), SecretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("host"), SecretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("database"), SecretModel.Database)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("jwtSign"), SecretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("bucketName"), os.Getenv("BucketName"))

	// Check database connection
	err = db.ConnectDB(awsgo.Ctx)
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error  conectando en la db " + err.Error(),
			Headers: map[string]string{
				"Content-type": "application/json",
			},
		}
		return res, nil
	}

	resAPI := handlers.Handlers(awsgo.Ctx, request)
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

package models

import "github.com/aws/aws-lambda-go/events"

type ResponseAPI struct {
	Status         int
	Message        string
	CustomResponse *events.APIGatewayProxyResponse
}

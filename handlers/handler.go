package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/vcoromero/tuitergo/models"
)

func Handlers(ctx context.Context, request events.APIGatewayProxyRequest) models.ResponseAPI {
	fmt.Println("Voy a procesar " + ctx.Value(models.Key("path")).(string) + " > " + ctx.Value(models.Key("method")).(string))

	var res models.ResponseAPI
	res.Status = 400

	switch ctx.Value(models.Key("method")).(string) {
	case "POST":
		switch ctx.Value(models.Key("path")).(string) {
		//
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

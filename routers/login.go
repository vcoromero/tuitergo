package routers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/vcoromero/tuitergo/db"
	"github.com/vcoromero/tuitergo/jwt"
	"github.com/vcoromero/tuitergo/models"
)

func Login(ctx context.Context) models.ResponseAPI {
	var t models.User
	var r models.ResponseAPI
	r.Status = 400

	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		r.Message = "Invalid user or password " + err.Error()
		return r
	}

	if len(t.Email) == 0 {
		r.Message = "User email is required"
		return r
	}
	userData, exists := db.TryLogin(t.Email, t.Password)
	if !exists {
		r.Message = "User and password are invalid"
		return r
	}

	jwtKey, err := jwt.GenerateJWT(ctx, userData)

	if err != nil {
		r.Message = "Occured an error trying to generate the token " + err.Error()
		return r
	}

	response := models.ResponseLogin{
		Token: jwtKey,
	}

	token, err := json.Marshal(response)
	if err != nil {
		r.Message = "Occured an error trying to parse the token" + err.Error()
		return r
	}

	cookie := &http.Cookie{
		Name:    "token",
		Value:   jwtKey,
		Expires: time.Now().Add(24 * time.Hour),
	}

	cookieString := cookie.String()

	res := &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(token),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
			"SetCookie":                   cookieString,
		},
	}

	r.Status = 200
	r.Message = string(token)
	r.CustomResponse = res
	return r
}

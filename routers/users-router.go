package routers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
	userData, exists := db.Login(t.Email, t.Password)
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

func GetUser(request events.APIGatewayProxyRequest) models.ResponseAPI {
	var r models.ResponseAPI
	r.Status = 200
	fmt.Println("Entered to show profile")
	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		r.Message = "id parameter is required"
		return r
	}

	profile, err := db.GetUser(ID)
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

func CreateUser(ctx context.Context) models.ResponseAPI {
	var t models.User
	var r models.ResponseAPI
	r.Status = 400
	fmt.Println("Registering user")

	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		r.Message = "Error parsing request body: " + err.Error()
		fmt.Println(r.Message)
		return r
	}

	if len(t.Email) == 0 {
		r.Message = "Must have an Email"
		fmt.Println(r.Message)
		return r
	}

	if len(t.Password) < 6 {
		r.Message = "Password must be at least 6 characters long"
		fmt.Println(r.Message)
		return r
	}

	_, found, _ := db.CheckedIfUserExist(t.Email)
	if found {
		r.Message = "There is an user registered with this email"
		fmt.Println(r.Message)
		return r
	}

	_, status, err := db.CreateUser(t)
	if err != nil {
		r.Message = "An error occurred trying to insert the user: " + err.Error()
		fmt.Println(r.Message)
		return r
	}

	if !status {
		r.Message = "The user could not be inserted"
		fmt.Println(r.Message)
		return r
	}

	r.Status = 200
	r.Message = "User registered!"
	fmt.Println(r.Message)
	return r
}

func UpdateUser(ctx context.Context, claim models.Claim) models.ResponseAPI {
	var r models.ResponseAPI
	r.Status = 400
	fmt.Println("Entered to show profile")

	var t models.User

	body := ctx.Value(models.Key("body")).(string)

	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		r.Message = "Wrong data " + err.Error()
		return r
	}

	status, err := db.UpdateUser(t, claim.ID.Hex())
	if err != nil {
		r.Message = "Ocurred error tryng to update user" + err.Error()
		return r
	}
	if !status {
		r.Message = "Cannot update user"
		return r
	}

	r.Status = 200
	r.Message = "Update user succesfully"
	return r
}

type readSeeker struct {
	io.Reader
}

func (rs *readSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func UploadImage(ctx context.Context, uploadType string, request events.APIGatewayProxyRequest, claim models.Claim) models.ResponseAPI {
	var r models.ResponseAPI
	r.Status = 400
	UserId := claim.ID.Hex()

	var filename string
	var user models.User

	bucket := aws.String(ctx.Value(models.Key("bucketName")).(string))

	switch uploadType {
	case "A":
		filename = "avatars/" + UserId + ".jpg"
		user.Avatar = filename
	case "B":
		filename = "banners/" + UserId + ".jpg"
		user.Banner = filename
	}

	contentType := request.Headers["Content-Type"]
	fmt.Println("Content-Type: ", contentType)
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		r.Status = 500
		r.Message = "mediatype error: " + err.Error()
		return r
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		body, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			r.Status = 500
			r.Message = err.Error()
			return r
		}

		mr := multipart.NewReader(bytes.NewReader(body), params["boundary"])
		p, err := mr.NextPart()
		if err != nil && err != io.EOF {
			r.Status = 500
			r.Message = err.Error()
			return r
		}

		if err != io.EOF {
			if p.FileName() != "" {

				buff := bytes.NewBuffer(nil)
				if _, err := io.Copy(buff, p); err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				session, err := session.NewSession(&aws.Config{
					Region: aws.String("us-east-1"),
				})

				if err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				uploader := s3manager.NewUploader(session)
				_, err = uploader.Upload(&s3manager.UploadInput{
					Bucket: bucket,
					Key:    aws.String(filename),
					Body:   &readSeeker{buff},
				})

				if err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}

			}
		}

		status, err := db.UpdateUser(user, UserId)

		if err != nil || !status {
			r.Status = 400
			r.Message = "Error occured tryign to update user " + err.Error()
			return r
		}

	} else {
		r.Message = "Must be an image with 'Content-Type' type"
		r.Status = 400
		return r
	}

	r.Status = 200
	r.Message = "Image file uploaded!!"
	return r
}

package jwt

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vcoromero/tuitergo/models"
)

func GenerateJWT(ctx context.Context, t models.User) (string, error) {
	jwtSign := ctx.Value(models.Key("jwtSign")).(string)
	myKey := []byte(jwtSign)

	payload := jwt.MapClaims{
		"email":     t.Email,
		"name":      t.Name,
		"lastnames": t.LastNames,
		"birthdate": t.Birthdate,
		"biography": t.Biography,
		"location":  t.Location,
		"webSite":   t.WebSite,
		"_id":       t.ID.Hex(),
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := token.SignedString(myKey)
	if err != nil {
		return tokenStr, err
	}
	return tokenStr, nil
}

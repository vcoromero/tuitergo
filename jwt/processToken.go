package jwt

import (
	"errors"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/vcoromero/tuitergo/models"
)

var Email string
var UserID string

func ProcessToken(tk string, JWTSign string) (*models.Claim, bool, string, error) {
	myKey := []byte(JWTSign)
	var claims models.Claim

	splitToken := strings.Split(tk, "Bearer")
	if len(splitToken) != 2 {
		return &claims, false, string(""), errors.New("invalid token format")
	}

	tk = strings.TrimSpace(splitToken[1])

	tkn, err := jwt.ParseWithClaims(tk, &claims, func(t *jwt.Token) (interface{}, error) {
		return myKey, nil
	})

	if err == nil {
		//rutina que chequea contra la db
	}

	if !tkn.Valid {
		return &claims, false, string(""), errors.New("invalid token")
	}

	return &claims, false, string(""), nil

}

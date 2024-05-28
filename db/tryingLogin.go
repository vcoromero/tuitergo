package db

import (
	"github.com/vcoromero/tuitergo/models"
	"golang.org/x/crypto/bcrypt"
)

func TryLogin(email string, password string) (models.User, bool) {
	user, found, _ := CheckedIfUserExist(email)
	if !found {
		return user, false
	}

	passwordBytes := []byte(password)
	passwordDB := []byte(user.Password)

	err := bcrypt.CompareHashAndPassword(passwordDB, passwordBytes)
	if err != nil {
		return user, false
	}
	return user, true
}

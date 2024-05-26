package db

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(psswrd string) (string, error) {
	cost := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(psswrd), cost)
	if err != nil {
		return err.Error(), err
	}
	return string(bytes), nil
}

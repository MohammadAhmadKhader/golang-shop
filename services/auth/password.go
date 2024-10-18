package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	
	return string(hash), nil
}

func ComparePassword(hashedPassword string, plain []byte) (isEqual bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), plain)
	fmt.Println(err)
	return err == nil
}
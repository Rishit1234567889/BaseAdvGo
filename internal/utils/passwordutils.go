package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// 3.2
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password")
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(strorePassword, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(strorePassword), []byte(providedPassword))
	return err == nil
}

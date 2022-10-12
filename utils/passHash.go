package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashEncrypt(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), 14)

	if err != nil {
		log.Fatal("Error in encrypting - ", err)
		return string(bytes), err
	}

	return string(bytes), nil
}

func CheckPasswordMatch(formPass, dbHashedPass string) bool {
	invalid := bcrypt.CompareHashAndPassword([]byte(dbHashedPass), []byte(formPass))
	// if err == nil, then both passwords are same
	return invalid == nil
}

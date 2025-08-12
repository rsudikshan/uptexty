package service

import (
	"backend/internal/runtime_errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJwt(email string) (string,error) {

	var err error

	claims := jwt.MapClaims{
		"email":email,
		"exp":time.Now().Add(time.Hour*24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	keySecret,ok := os.LookupEnv("KEY_SECRET")

	if !ok {
		return "",&runtime_errors.InternalServerError{
			Message: "Key Secret Env not found",
		}
	}
	
	signedString,err := token.SignedString([]byte(keySecret))

	if err!=nil {
		return "",&runtime_errors.InternalServerError{
			Message: err.Error(),
		}
	}

	return signedString,nil
}
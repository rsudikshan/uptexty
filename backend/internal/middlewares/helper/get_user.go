package helper

import (
	"backend/internal/runtime_errors"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserIdFromClaims(claims jwt.MapClaims) (int,error) {
	id, ok := claims["id"].(float64)

	if !ok {
		return 0,&runtime_errors.InternalServerError{
			Message: "Couldnt extract from claims",
		}
	}

	return int(id),nil
}
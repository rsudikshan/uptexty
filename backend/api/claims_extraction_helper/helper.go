package claims_extraction_helper

import (
	"backend/internal/middlewares/helper"
	"backend/internal/runtime_errors"
	"github.com/golang-jwt/jwt/v5"

)

func ParseClaims(claimsValue any) (int,error) {

	if claimsValue == nil {

		return 0,&runtime_errors.InternalServerError{
			Message: "Couldnt parse claims",
		}
	}

	claims, ok := claimsValue.(jwt.MapClaims)
	if !ok {
		return 0,&runtime_errors.InternalServerError{
			Message: "Couldnt parse claims",
		}
	}

	return helper.GetUserIdFromClaims(claims)
}
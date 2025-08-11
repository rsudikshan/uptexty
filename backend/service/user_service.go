package service

import (
	"backend/internal/db"
	"backend/internal/runtime_errors"
	"backend/payloads/request"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(registerRequest request.RegisterRequest) error {
	
	var err error

	if registerRequest.Email == "" || registerRequest.Password == "" || registerRequest.Username == "" {
		return &runtime_errors.BadRequestError{
			Message: "Fields cannot be null",
		}
	}

	queryStr := "INSERT into user_table(email,username,password) VALUES ($1,$2,$3)"

	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password),12)

	if err!=nil {
		return &runtime_errors.InternalServerError{
			Message: err.Error(),
		}
	}

	_,err = db.DB.Exec(queryStr,registerRequest.Email,registerRequest.Username,hashedPassword)

	if err!=nil {
		return &runtime_errors.InternalServerError{
			Message: err.Error(),
		}
	}

	return nil
}
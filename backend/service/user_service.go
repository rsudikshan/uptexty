package service

import (
	"backend/internal/db"
	"backend/internal/runtime_errors"
	"backend/payloads/request"
	"backend/payloads/response"
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

func LoginUser(loginRequest request.LoginRequest) (response.LoginResponse, error) {

	var err error

	var response response.LoginResponse

	if loginRequest.Email == "" || loginRequest.Password == ""{
		return response,&runtime_errors.BadRequestError{
			Message: "Fields cannot be empty.",
		}
	}

	queryStr := "SELECT id,password FROM user_table WHERE email = $1"

	resultSet,err :=  db.DB.Query(queryStr,loginRequest.Email)

	if err!=nil{
		return response,&runtime_errors.BadRequestError{
			Message: err.Error(),
		}
	}

	if !resultSet.Next(){
		return response,&runtime_errors.UnauthorizedError{
			Message: "User not found.",
		}
	}

	var id int
	var password string

	err = resultSet.Scan(&id,&password)

	if err!=nil{
		return response,&runtime_errors.BadRequestError{
			Message: err.Error(),
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(password),[]byte(loginRequest.Password))

	if err !=nil{
		return response,&runtime_errors.UnauthorizedError{
			Message: "Invalid password.",
		}
	}

	response.Jwt,err = CreateJwt(id,loginRequest.Email)

	if err!= nil {
		return response,&runtime_errors.UnauthorizedError{
			Message: err.Error(),
		}
	}

	return response,nil
}
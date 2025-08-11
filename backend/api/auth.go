package api

import (
	"backend/global"
	"backend/payloads/request"
	"backend/service"
	"encoding/json"
	"net/http"
)

func Register(w http.ResponseWriter, req *http.Request) {
	var err error
	if req.Method != http.MethodPost{
		http.Error(w,"Invalid request method",http.StatusBadRequest)
		return
	}

	var registerRequest request.RegisterRequest

	err = json.NewDecoder(req.Body).Decode(&registerRequest)

	if err!=nil {
		http.Error(w,"Internal server error",http.StatusInternalServerError)
		return
	}

	err = service.RegisterUser(registerRequest)

	if err!=nil {
		global.HandleError(err,w)
		return
	}

	global.Success("User registered successfully",w)
}
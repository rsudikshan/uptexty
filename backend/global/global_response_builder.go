package global

import (
	"backend/models"
	"encoding/json"
	"net/http"
)

func Success(message string, w http.ResponseWriter) {
	sendResponse(&models.ResponseModel{
		Message: message,
		Success: true,
	},w)
}

func SuccessWithBody(message string,data any, w http.ResponseWriter) {
	sendResponse(&models.ResponseModel{
		Message: message,
		Success: true,
		Body: data,
	},w)
}

func sendResponse(data *models.ResponseModel, w http.ResponseWriter) {
	w.Header().Set("Content-type","application/json")
	json.NewEncoder(w).Encode(data)
}
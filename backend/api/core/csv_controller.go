package core

import (
	"backend/api/claims_extraction_helper"
	"backend/global"
	"backend/internal/middlewares"
	"backend/service/core_service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func UploadCsv(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost{
		http.Error(w,"Invalid request method",http.StatusBadRequest)
		return
	}

	req.ParseMultipartForm( 10 << 20)

	file,_,err := req.FormFile("file")
	filename := req.FormValue("filename")

	if err!=nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	claimsValue := req.Context().Value(middlewares.ClaimsKey)

	id,err := claims_extraction_helper.ParseClaims(claimsValue)

	if err!=nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}

	err = core_service.UploadCsvService(file,filename,id)

	if err!=nil {
		global.HandleError(err,w)
		return
	}

	global.Success("File uploaded successfully",w)
	defer file.Close()

}

func GetUploadedFiles( w http.ResponseWriter, req *http.Request){
	if req.Method != http.MethodGet{
		http.Error(w,"Invalid request method",http.StatusBadRequest)
		return
	}
	claimsValue := req.Context().Value(middlewares.ClaimsKey)

	id,err := claims_extraction_helper.ParseClaims(claimsValue)

	if err!=nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}

	response,err := core_service.GetUploadedFilesService(id)

	if err!=nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}

	global.SuccessWithBody("Success",response,w)
}


func GetRows( w http.ResponseWriter, req *http.Request){
	if req.Method != http.MethodGet{
		http.Error(w,"Invalid request method",http.StatusBadRequest)
		return
	}

	values := mux.Vars(req)
	claimsValue := req.Context().Value(middlewares.ClaimsKey)

	id,err := claims_extraction_helper.ParseClaims(claimsValue)

	if err!=nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}

	fileID, err := strconv.Atoi(values["id"])

	if err!=nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}

	response,err := core_service.GetRows(id,fileID)

	if err!=nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}

	global.SuccessWithBody("Success",response,w)
}

//new
func CreateRow(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	values := mux.Vars(req)
	claimsValue := req.Context().Value(middlewares.ClaimsKey)

	userID, err := claims_extraction_helper.ParseClaims(claimsValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileID, err := strconv.Atoi(values["id"])
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var reqBody struct {
		Position  float64 `json:"position"`
		InputText string  `json:"input_text"`
	}

	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Validate input
	if reqBody.InputText == "" {
		http.Error(w, "Input text cannot be empty", http.StatusBadRequest)
		return
	}

	response, err := core_service.CreateRowService(userID, fileID, reqBody.Position, reqBody.InputText)
	if err != nil {
		global.HandleError(err, w)
		return
	}

	global.SuccessWithBody("Row created successfully", response, w)
}

// UpdateRow handles PUT /files/{fileId}/rows/{rowId}
func UpdateRow(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	values := mux.Vars(req)
	claimsValue := req.Context().Value(middlewares.ClaimsKey)

	userID, err := claims_extraction_helper.ParseClaims(claimsValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileID, err := strconv.Atoi(values["fileId"])
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	rowID, err := strconv.Atoi(values["rowId"])
	if err != nil {
		http.Error(w, "Invalid row ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var reqBody struct {
		Position  float64 `json:"position"`
		InputText string  `json:"input_text"`
	}

	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Validate input
	if reqBody.InputText == "" {
		http.Error(w, "Input text cannot be empty", http.StatusBadRequest)
		return
	}

	response, err := core_service.UpdateRowService(userID, fileID, rowID, reqBody.Position, reqBody.InputText)
	if err != nil {
		global.HandleError(err, w)
		return
	}

	global.SuccessWithBody("Row updated successfully", response, w)
}

// DeleteRow handles DELETE /files/{fileId}/rows/{rowId}
func DeleteRow(w http.ResponseWriter, req *http.Request) {
	fmt.Print("hit")
	if req.Method != http.MethodDelete {
		fmt.Print("-1")
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	values := mux.Vars(req)
	claimsValue := req.Context().Value(middlewares.ClaimsKey)

	userID, err := claims_extraction_helper.ParseClaims(claimsValue)
	if err != nil {
		fmt.Print("0")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileID, err := strconv.Atoi(values["fileId"])
	if err != nil {
		fmt.Print("1")
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	rowID, err := strconv.Atoi(values["rowId"])
	if err != nil {
		fmt.Print("2")
		http.Error(w, "Invalid row ID", http.StatusBadRequest)
		return
	}

	err = core_service.DeleteRowService(userID, fileID, rowID)
	if err != nil {
		fmt.Print("3")
		global.HandleError(err, w)
		return
	}

	global.Success("Row deleted successfully", w)
}
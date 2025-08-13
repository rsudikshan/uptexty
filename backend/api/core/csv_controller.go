package core

import (
	"backend/api/claims_extraction_helper"
	"backend/global"
	"backend/internal/middlewares"
	"backend/service/core_service"
	"net/http"
)

func UploadCsv(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost{
		http.Error(w,"Invalid request method",http.StatusBadRequest)
		return
	}

	req.ParseMultipartForm( 10 << 20)

	file,_,err := req.FormFile("file")

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

	err = core_service.UploadCsvService(file,"test",id)

	if err!=nil {
		global.HandleError(err,w)
		return
	}

	global.Success("File uploaded successfully",w)
	defer file.Close()

}
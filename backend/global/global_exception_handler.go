package global

import (
	"backend/internal/runtime_errors"
	"backend/models"
	"encoding/json"
	"net/http"
)

func HandleError(e error, w http.ResponseWriter) {
	switch e.(type) {
	case *runtime_errors.InternalServerError:
		w.WriteHeader(http.StatusInternalServerError)

	case *runtime_errors.BadRequestError:
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-type","application/json")
	json.NewEncoder(w).Encode(models.ResponseModel{
		Success: false,
		Message: "Operation Failed: "+e.Error(),
	})
}

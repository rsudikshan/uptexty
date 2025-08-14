package global

import (
	"backend/internal/runtime_errors"
	"backend/models"
	"encoding/json"
	"net/http"
)

func HandleError(e error, w http.ResponseWriter) {

	w.Header().Set("Content-type","application/json")

	switch e.(type) {
	case *runtime_errors.InternalServerError:
		w.WriteHeader(http.StatusInternalServerError)

	case *runtime_errors.BadRequestError:
		w.WriteHeader(http.StatusBadRequest)

	
	case *runtime_errors.UnauthorizedError:
		w.WriteHeader(http.StatusUnauthorized)
	
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

	
	json.NewEncoder(w).Encode(models.ResponseModel{
		Success: false,
		Message: "Operation Failed: "+e.Error(),
	})
}

package helper

import (
	"backend/internal/db"
	"backend/internal/runtime_errors"
)

func EmailExists(email string) (bool, error) {
    queryStr := "SELECT EXISTS (SELECT 1 FROM user_table WHERE email = $1)"
    var exists bool
    err := db.DB.QueryRow(queryStr, email).Scan(&exists)
    if err != nil {
        return false, &runtime_errors.InternalServerError{
            Message: err.Error(),
        }
    }
    return exists, nil
}

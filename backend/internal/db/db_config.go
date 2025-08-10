package db

import (
	"backend/internal/runtime_errors"
	"database/sql"
	"os"
	_ "github.com/lib/pq"

)

var DB *sql.DB

func ConnectToDbServer() error {

	var err error

	dbUser,exists := os.LookupEnv("DB_USER")

	if !exists {
		return &runtime_errors.InternalServerError{
			Message: "Db username not found",
		}
	}

	dbName,exists := os.LookupEnv("DB_NAME")

	if !exists {
		return &runtime_errors.InternalServerError{
			Message: "Db username not found",
		}
	}

	dbPassword,exists := os.LookupEnv("DB_PASSWORD")

	if !exists {
		return &runtime_errors.InternalServerError{
			Message: "Db password not found",
		}
	}

	dbPort,exists := os.LookupEnv("DB_PORT")

	if !exists {
		return &runtime_errors.InternalServerError{
			Message: "Db port not found",
		}
	}

	connStr := 	"user="+dbUser+
				" password="+dbPassword+
				" port="+dbPort+
				" dbname="+dbName+
				" sslmode=disable"

	DB,err = sql.Open("postgres",connStr)

	if err!=nil{
		return &runtime_errors.InternalServerError{
			Message: err.Error(),
		}
	}

	err = DB.Ping()

	if err!=nil{
		return &runtime_errors.InternalServerError{
			Message: err.Error(),
		}
	}

	return nil

}
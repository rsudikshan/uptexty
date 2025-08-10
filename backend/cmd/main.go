package main

import (
	"backend/api"
	"backend/internal/db"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	var err error

	//env path
	err = godotenv.Load("C:/uptexty-assignment/backend/.env")

	if err!=nil{
		fmt.Println("Error starting server: "+err.Error())
		return
	}

	err = db.ConnectToDbServer()

	if err!=nil{
		fmt.Println("Error starting server: "+err.Error())
		return
	}

	port,exists := os.LookupEnv("PORT")

	if !exists {
		fmt.Println("Error starting server: Port address not specified.")
		return
	}

	LoadApis()

	//starting the server
	err = http.ListenAndServe(port,nil)

	if err!=nil{
		fmt.Println("Error starting server: "+err.Error())
		return
	}

}

func LoadApis(){
	http.HandleFunc("/test",api.Test)
}
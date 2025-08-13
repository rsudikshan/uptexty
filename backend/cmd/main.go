package main

import (
	"backend/api"
	"backend/api/core"
	"backend/internal/db"
	"backend/internal/middlewares"
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
	
	defer db.DB.Close()

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
	http.Handle("/test",middlewares.JwtFilter(http.HandlerFunc(api.Test)))
	http.HandleFunc("/register",api.Register)
	http.HandleFunc("/login",api.Login)
	http.Handle("/upload",middlewares.JwtFilter(http.HandlerFunc(core.UploadCsv)))
}
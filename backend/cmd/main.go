package main

import (
	"backend/api"
	"backend/api/core"
	"backend/internal/db"
	"backend/internal/middlewares"
	"fmt"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	var err error

	// Load env file
	err = godotenv.Load("C:/uptexty-assignment/backend/.env")
	if err != nil {
		fmt.Println("Error loading .env file: " + err.Error())
		return
	}

	// Connect to DB
	err = db.ConnectToDbServer()
	defer db.DB.Close()
	if err != nil {
		fmt.Println("Error connecting to DB: " + err.Error())
		return
	}

	port, exists := os.LookupEnv("PORT")
	if !exists {
		fmt.Println("PORT not specified in environment.")
		return
	}

	// Create router and load routes
	router := mux.NewRouter()
	LoadApis(router)

	// Start server
	fmt.Println("Server running on port " + port)
	err = http.ListenAndServe(port, router)
	if err != nil {
		fmt.Println("Error starting server: " + err.Error())
	}
}

func LoadApis(router *mux.Router) {
	// Apply CORS to all routes globally
	router.Use(middlewares.CorsMiddleware)

	// Public routes (no JWT)
	router.Handle("/register", http.HandlerFunc(api.Register))
	router.Handle("/login", http.HandlerFunc(api.Login))

	// Authenticated routes
	router.Handle("/test",
		middlewares.JwtFilter(http.HandlerFunc(api.Test)),
	)

	// File handling
	router.Handle("/upload",
		middlewares.JwtFilter(http.HandlerFunc(core.UploadCsv)),
	)
	router.Handle("/files",
		middlewares.JwtFilter(http.HandlerFunc(core.GetUploadedFiles)),
	)
	router.Handle("/files/{id}",
		middlewares.JwtFilter(http.HandlerFunc(core.GetRows)),
	)

	// Row operations
	router.Handle("/files/{id}/rows",
		middlewares.JwtFilter(http.HandlerFunc(core.CreateRow)),
	).Methods("POST")

	router.Handle("/files/{fileId}/rows/{rowId}",
		middlewares.JwtFilter(http.HandlerFunc(core.UpdateRow)),
	).Methods("PUT")

	router.Handle("/files/{fileId}/rows/{rowId}",
		middlewares.JwtFilter(http.HandlerFunc(core.DeleteRow)),
	).Methods("DELETE")

	// Global OPTIONS handler for CORS preflight
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

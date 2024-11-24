package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var database *mongo.Database

func init() {
	MONGODB_URL := GetDotEnvVariable("MONGODB_URL")
	client = ConnectToMongoDB(MONGODB_URL)
	database = client.Database("myNewDatabase")
}

func main() {
	router := mux.NewRouter()

	router.Use(LoggingMiddleware)

	router.HandleFunc("/books", GetBooks).Methods(http.MethodGet)
	router.HandleFunc("/books/{id}", GetBook).Methods(http.MethodGet)
	router.HandleFunc("/books", CreateBook).Methods(http.MethodPost)
	router.HandleFunc("/books/{id}", UpdateBook).Methods(http.MethodPut)
	router.HandleFunc("/books/{id}", DeleteBook).Methods(http.MethodDelete)

	PORT := GetDotEnvVariable("SERVER_PORT")
	if PORT == "" {
		PORT = ":8082"
	}
	log.Printf("Starting server on %v", PORT)
	if err := http.ListenAndServe(PORT, router); err != nil {
		log.Fatalf("Server failed to start on %v", PORT)
	}
}

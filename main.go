package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.Use(LoggingMiddleware)

	books = append(books, Book{
		ID: "1", Title: "Book One", Author: "John Doe",
		Publisher: &Company{Name: "Company One", Address: "Address One"}})
	books = append(books, Book{
		ID: "2", Title: "Book Two", Author: "Jane Doe",
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})

	router.HandleFunc("/books", GetBooks).Methods(http.MethodGet)
	router.HandleFunc("/books/{id}", GetBook).Methods(http.MethodGet)
	router.HandleFunc("/books", CreateBook).Methods(http.MethodPost)
	router.HandleFunc("/books/{id}", UpdateBook).Methods(http.MethodPut)
	router.HandleFunc("/books/{id}", DeleteBook).Methods(http.MethodDelete)

	PORT := os.Getenv("SERVER_PORT")
	if PORT == "" {
		PORT = ":8082"
	}
	log.Printf("Starting server on %v", PORT)
	if err := http.ListenAndServe(PORT, router); err != nil {
		log.Fatalf("Server failed to start on %v", PORT)
	}
}

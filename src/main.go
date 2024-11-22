package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Company struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

type Book struct {
	ID        string   `json:"id,omitempty"`
	Title     string   `json:"title,omitempty"`
	Author    string   `json:"author,omitempty"`
	Publisher *Company `json:"publisher"`
}

var books []Book

func GetBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Book not found"})
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

type logginResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (lwr *logginResponseWriter) WriteHeader(code int) {
	lwr.statusCode = code
	lwr.ResponseWriter.WriteHeader(code)
}
func (lwr *logginResponseWriter) Write(data []byte) (int, error) {
	size, err := lwr.ResponseWriter.Write(data)
	lwr.size += size
	return size, err
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &logginResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		lrw.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(lrw, r)

		// Log in access log format
		log.Printf("%s - %s [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"",
			r.RemoteAddr, // Client IP and port
			"-",          // User identifier (typically not available in web apps)
			start.Format("02/Jan/2006:15:04:05 -0700"), // Timestamp
			r.Method,       // HTTP method
			r.URL.Path,     // URL path
			r.Proto,        // HTTP version
			lrw.statusCode, // Response status code
			lrw.size,       // Response size in bytes
			r.Referer(),    // Referer header
			r.UserAgent(),  // User-Agent header
		)
	})
}

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

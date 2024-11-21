package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	Publisher *Company `json:"publisher,omitempty"`
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
	json.NewEncoder(w).Encode(&Book{})
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	books = append(books, book)
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

func main() {
	router := mux.NewRouter()

	books = append(books, Book{ID: "1", Title: "Book One", Author: "John Doe", Publisher: &Company{Name: "Company One", Address: "Address One"}})
	books = append(books, Book{ID: "2", Title: "Book Two", Author: "Jane Doe", Publisher: &Company{Name: "Company Two", Address: "Address Two"}})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})
	router.HandleFunc("/books", GetBooks).Methods(http.MethodGet)
	router.HandleFunc("/books{id}", GetBook).Methods(http.MethodGet)
	router.HandleFunc("/books", CreateBook).Methods(http.MethodPost)
	router.HandleFunc("/books{id}", UpdateBook).Methods(http.MethodPut)
	router.HandleFunc("/books{id}", DeleteBook).Methods(http.MethodDelete)

	http.ListenAndServe(":8082", router)
}

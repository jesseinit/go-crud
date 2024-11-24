package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	collection := database.Collection("books")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		http.Error(w, "Failed to fetch books: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var books []map[string]interface{}
	if err := cursor.All(ctx, &books); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to decode books: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	collection := database.Collection("books")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	params := mux.Vars(r)
	bookID := params["id"]

	var book map[string]interface{}
	err := collection.FindOne(ctx, bson.M{"_id": bookID}).Decode(&book)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Book not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch books: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	collection := database.Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	book.ID = primitive.NewObjectID().Hex()

	_, err = collection.InsertOne(ctx, book)
	if err != nil {
		http.Error(w, "Failed to create book: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookID := params["id"]

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	collection := database.Collection("books")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateByID(ctx, bookID, bson.M{"$set": updateData})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update Book"})
		return
	}

	var book map[string]interface{}
	err = collection.FindOne(ctx, bson.M{"_id": bookID}).Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update book: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookID := params["id"]

	collection := database.Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, delete_err := collection.DeleteOne(ctx, bson.M{"_id": bookID})
	if delete_err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Book Delete Failed"})
		return
	}

	if res.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Book has not been deleted"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

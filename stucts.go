package main

import "net/http"

type Company struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

type Book struct {
	ID        string   `json:"-" bson:"_id,omitempty"`
	Title     string   `json:"title,omitempty"`
	Author    string   `json:"author,omitempty"`
	Publisher *Company `json:"publisher"`
}

type logginResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

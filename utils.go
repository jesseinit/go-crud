package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
		log.Printf("%s - %s [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"",
			r.RemoteAddr,
			"-",
			start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.Path,
			r.Proto,
			lrw.statusCode,
			lrw.size,
			r.Referer(),
			r.UserAgent(),
		)
	})
}

func ConnectToMongoDB(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")

	return client
}

func GetDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

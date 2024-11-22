package main

import (
	"log"
	"net/http"
	"time"
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

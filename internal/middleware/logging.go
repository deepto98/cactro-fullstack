package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs the HTTP method, URI, and the time taken to process the request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s - %v", r.Method, r.RequestURI, time.Since(start))
	})
}

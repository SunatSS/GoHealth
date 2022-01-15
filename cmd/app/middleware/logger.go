package middleware

import (
	"log"
	"net/http"
)

// Logger is a middleware fuction that logs path and method of request
func Logger(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Start: %s %s", r.Method, r.URL.Path)
		handle.ServeHTTP(w, r)
		log.Printf("FINISH: %s %s", r.Method, r.URL.Path)
	})
}
package main

import (
	"fmt"
	"net/http"
)

func (app *application) enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func (app *application) authRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session ID (cookie value) from the request.
		_, err := app.GetSessionIDFromCookie(r)
		if err != nil {
			// If the user is not authenticated, redirect to a login page or return an error message.
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			fmt.Printf("User not authorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}

package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	handler := app.enableCORS(mux)

	mux.HandleFunc("/", app.HomeHandler)
	mux.HandleFunc("/register", app.RegisterHandler)
	mux.HandleFunc("/login", app.LoginHandler)

	return handler
}

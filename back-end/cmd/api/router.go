package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	handler := app.enableCORS(mux)

	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/register", app.Register)
	mux.HandleFunc("/login", app.Login)

	return handler
}

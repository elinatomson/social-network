package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	handler := app.enableCORS(mux)

	mux.HandleFunc("/", app.HomeHandler)
	mux.HandleFunc("/register", app.RegisterHandler)
	mux.HandleFunc("/login", app.LoginHandler)
	mux.HandleFunc("/logout", app.LogOutHandler)

	mux.Handle("/profile", app.authRequired(http.HandlerFunc(app.ProfileHandler)))
	mux.Handle("/social", app.authRequired(http.HandlerFunc(app.SocialHandler)))
	mux.Handle("/search", app.authRequired(http.HandlerFunc(app.SearchHandler)))

	return handler
}

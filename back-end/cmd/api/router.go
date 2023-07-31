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
	mux.Handle("/main", app.authRequired(http.HandlerFunc(app.MainPageHandler)))
	mux.Handle("/search", app.authRequired(http.HandlerFunc(app.SearchHandler)))
	mux.Handle("/user/", app.authRequired(http.HandlerFunc(app.UserHandler)))
	mux.Handle("/create-post", app.authRequired(http.HandlerFunc(app.CreatePostHandler)))
	mux.Handle("/all-posts", app.authRequired(http.HandlerFunc(app.AllPostsHandler)))
	mux.Handle("/create-comment", app.authRequired(http.HandlerFunc(app.CommentHandler)))

	return handler
}

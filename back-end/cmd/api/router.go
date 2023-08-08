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
	mux.Handle("/profile-type", app.authRequired(http.HandlerFunc(app.ProfileTypeHandler)))
	mux.Handle("/main", app.authRequired(http.HandlerFunc(app.MainPageHandler)))
	mux.Handle("/users", app.authRequired(http.HandlerFunc(app.GetUsersHandler)))
	mux.Handle("/search", app.authRequired(http.HandlerFunc(app.SearchHandler)))
	mux.Handle("/user/", app.authRequired(http.HandlerFunc(app.UserHandler)))
	mux.Handle("/create-post", app.authRequired(http.HandlerFunc(app.CreatePostHandler)))
	mux.Handle("/all-posts", app.authRequired(http.HandlerFunc(app.AllPostsHandler)))
	mux.Handle("/create-comment", app.authRequired(http.HandlerFunc(app.CommentHandler)))
	mux.Handle("/follow", app.authRequired(http.HandlerFunc(app.FollowHandler)))
	mux.Handle("/following", app.authRequired(http.HandlerFunc(app.FollowingHandler)))
	mux.Handle("/followers", app.authRequired(http.HandlerFunc(app.FollowersHandler)))
	mux.Handle("/follow-requests", app.authRequired(http.HandlerFunc(app.FollowRequestsHandler)))
	mux.Handle("/accept-follower", app.authRequired(http.HandlerFunc(app.AcceptFollowerHandler)))
	mux.Handle("/decline-follower", app.authRequired(http.HandlerFunc(app.DeclineFollowerHandler)))
	mux.Handle("/ws", app.authRequired(http.HandlerFunc(app.WebsocketHandler)))
	mux.Handle("/message", app.authRequired(http.HandlerFunc(app.AddMessageHandler)))
	mux.Handle("/messages", app.authRequired(http.HandlerFunc(app.GetMessagesHandler)))

	return handler
}

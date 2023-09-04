package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	handler := app.enableCORS(mux)

	mux.HandleFunc("/", app.HomeHandler)
	mux.HandleFunc("/register", app.RegisterHandler)
	mux.HandleFunc("/login", app.LoginHandler)
	mux.HandleFunc("/logout", app.LogOutHandler)

	fileServer := http.FileServer(http.Dir("./database/images"))
	mux.Handle("/images/", http.StripPrefix("/images/", fileServer))

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
	mux.Handle("/create-group", app.authRequired(http.HandlerFunc(app.CreateGroupHandler)))
	mux.Handle("/all-groups", app.authRequired(http.HandlerFunc(app.AllGroupsHandler)))
	mux.Handle("/group/", app.authRequired(http.HandlerFunc(app.GroupHandler)))
	mux.Handle("/group-posts", app.authRequired(http.HandlerFunc(app.GroupPostsHandler)))
	mux.Handle("/invite", app.authRequired(http.HandlerFunc(app.InviteNewMemberHandler)))
	mux.Handle("/group-invitations", app.authRequired(http.HandlerFunc(app.GroupInvitationHandler)))
	mux.Handle("/accept-group-invitation", app.authRequired(http.HandlerFunc(app.AcceptGroupInvitationHandler)))
	mux.Handle("/decline-group-invitation", app.authRequired(http.HandlerFunc(app.DeclineGroupInvitationHandler)))
	mux.Handle("/request-to-join-group", app.authRequired(http.HandlerFunc(app.RequestToJoinGroupHandler)))
	mux.Handle("/group-requests", app.authRequired(http.HandlerFunc(app.GroupRequestsHandler)))
	mux.Handle("/accept-group-request", app.authRequired(http.HandlerFunc(app.AcceptGroupRequestHandler)))
	mux.Handle("/decline-group-request", app.authRequired(http.HandlerFunc(app.DeclineGroupRequestHandler)))
	mux.Handle("/create-event", app.authRequired(http.HandlerFunc(app.CreateEventHandler)))
	mux.Handle("/group-event-notifications", app.authRequired(http.HandlerFunc(app.GroupEventNotificationsHandler)))
	mux.Handle("/group-event-seen", app.authRequired(http.HandlerFunc(app.EventSeenHandler)))
	mux.Handle("/group-events", app.authRequired(http.HandlerFunc(app.GroupEventsHandler)))
	mux.Handle("/group-event/", app.authRequired(http.HandlerFunc(app.GroupEventHandler)))
	mux.Handle("/going", app.authRequired(http.HandlerFunc(app.GoingHandler)))
	mux.Handle("/not-going", app.authRequired(http.HandlerFunc(app.NotGoingHandler)))
	mux.Handle("/ws", app.authRequired(http.HandlerFunc(app.WebsocketHandler)))
	mux.Handle("/message", app.authRequired(http.HandlerFunc(app.AddMessageHandler)))
	mux.Handle("/messages", app.authRequired(http.HandlerFunc(app.GetMessagesHandler)))

	return handler
}

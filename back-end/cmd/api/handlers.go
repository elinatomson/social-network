package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"back-end/models"
)

type SqliteDB struct {
	DB *sql.DB
}

func (app *application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Social Network up and running",
		Version: "1.0.0",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/register" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var userData models.UserData
	err := app.readJSON(w, r, &userData)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	stmt := `SELECT email FROM users WHERE email = ?`
	row := app.database.DB.QueryRow(stmt, userData.Email)
	var email string
	err = row.Scan(&email)
	if err != sql.ErrNoRows {
		app.errorJSON(w, fmt.Errorf("Email already taken"), http.StatusConflict)
		return
	}

	err = app.database.Register(&userData)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error adding data to the database"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, userData)
}

type appError struct {
	message string
}

func (e *appError) Error() string {
	return e.message
}

func (app *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/login" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var userData models.UserData
	err := app.readJSON(w, r, &userData)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	err = app.database.Login(&userData)
	if err != nil {
		app.errorJSON(w, &appError{message: "Email or password is not correct!"}, http.StatusUnauthorized)
		return
	}

	userId, email, firstName, lastName, err := app.database.DataFromUserData(&userData)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from user data"), http.StatusInternalServerError)
		return
	}

	cookieValue := app.addCookie(w, userId, email, firstName, lastName)

	app.writeJSON(w, http.StatusOK, map[string]string{"session": cookieValue})
}

func (app *application) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	app.deleteCookie(r)
	w.WriteHeader(http.StatusAccepted)
}

func (app *application) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/profile" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	_, email, firstName, lastName, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the email"), http.StatusInternalServerError)
		return
	}

	// Query the database to retrieve the user data based on the his email
	userData, err := app.database.GetUserDataByEmail(email)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to fetch user data"), http.StatusInternalServerError)
		return
	}

	allPosts, err := app.database.ProfilePosts(firstName, lastName)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from the database"), http.StatusInternalServerError)
		return
	}

	for i := range allPosts {
		postID := allPosts[i].PostID
		comments, err := app.database.GetCommentsByPostID(postID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Error getting comments from the database"), http.StatusInternalServerError)
			return
		}

		allPosts[i].Comments = comments
	}

	userDataWithPosts := struct {
		UserData *models.UserData `json:"user_data"`
		Posts    []models.Post    `json:"posts"`
	}{
		UserData: userData,
		Posts:    allPosts,
	}

	_ = app.writeJSON(w, http.StatusOK, userDataWithPosts)
}

func (app *application) MainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/main" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	_, email, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the email"), http.StatusInternalServerError)
		return
	}

	userData, err := app.database.GetUserDataByEmail(email)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to fetch user data"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, userData)
}

func (app *application) SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/search" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	users, err := app.database.SearchUsers(query)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error searching users: %s", err), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, users)
}

func (app *application) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	_, _, firstName, lastName, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from user sessions"), http.StatusInternalServerError)
		return
	}

	users, err := app.database.GetUsers()
	//setting the currentUser in the users db table as a true to add the current user's nickname to the response
	for i := range users {
		if users[i].FirstName == firstName && users[i].LastName == lastName {
			users[i].CurrentUser = true
			break
		}
	}

	_ = app.writeJSON(w, http.StatusOK, users)
}

func (app *application) UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/user/")
	id1, err := strconv.Atoi(id)

	user, err := app.database.GetUser(id1)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	allPosts, err := app.database.GetPostsByUserID(user.UserID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	userID, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from user sessions"), http.StatusInternalServerError)
		return
	}

	var filteredPosts []models.Post
	for i := range allPosts {
		post := allPosts[i]
		if post.Privacy == "public" {
			filteredPosts = append(filteredPosts, post)
		} else if post.Privacy == "for-selected-users" {
			if post.SelectedUserID == "" {
				continue
			}
			selectedUserIDs := strings.Split(post.SelectedUserID, ",")
			for _, user := range selectedUserIDs {
				selectedUserID, err := strconv.Atoi(user)
				if err != nil {
					app.errorJSON(w, fmt.Errorf("Error converting selected user ID to integer"), http.StatusInternalServerError)
					return
				}
				if selectedUserID == userID {
					filteredPosts = append(filteredPosts, post)
					break
				}
			}
		} else {
			isFollowing, err := app.database.IsFollowing(userID, post.UserID)
			if err != nil {
				app.errorJSON(w, fmt.Errorf("Error checking if the user is following the post author"), http.StatusInternalServerError)
				return
			}
			if isFollowing {
				filteredPosts = append(filteredPosts, post)
			}
		}
	}

	//comments for each post and add them to the filtered posts
	for i := range filteredPosts {
		postID := filteredPosts[i].PostID
		comments, err := app.database.GetCommentsByPostID(postID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Error getting comments from the database"), http.StatusInternalServerError)
			return
		}

		filteredPosts[i].Comments = comments
	}

	userDataWithPosts := struct {
		UserData *models.UserData `json:"user_data"`
		Posts    []models.Post    `json:"posts"`
	}{
		UserData: user,
		Posts:    filteredPosts,
	}

	_ = app.writeJSON(w, http.StatusOK, userDataWithPosts)
}

func (app *application) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/create-post" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var post models.Post
	err := app.readJSON(w, r, &post)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	userId, _, firstName, lastName, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from user sessions"), http.StatusInternalServerError)
		return
	}

	post.UserID = userId
	post.FirstName = firstName
	post.LastName = lastName

	err = app.database.CreatePost(&post)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error adding data to the database"), http.StatusInternalServerError)
		return
	}
	// Include an empty comments array for the newly created post.
	post.Comments = make([]models.Comment, 0)

	_ = app.writeJSON(w, http.StatusOK, post)
}

func (app *application) AllPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/all-posts" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	userID, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from user sessions"), http.StatusInternalServerError)
		return
	}

	var allPosts []models.Post

	allPosts, err = app.database.AllPosts()
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from the database"), http.StatusInternalServerError)
		return
	}

	// Filter posts based on privacy setting
	var filteredPosts []models.Post
	for _, post := range allPosts {
		if post.Privacy == "public" || post.UserID == userID {
			filteredPosts = append(filteredPosts, post)
		} else if post.Privacy == "for-selected-users" {
			// Splitting the comma-separated string of selected user IDs into a slice of strings
			selectedUserIDs := strings.Split(post.SelectedUserID, ",")
			// Checking if the logged-in user's ID is in the selectedUserIDs slice
			for _, user := range selectedUserIDs {
				selectedUserID, err := strconv.Atoi(user)
				if err != nil {
					app.errorJSON(w, fmt.Errorf("Error converting selected user ID to integer"), http.StatusInternalServerError)
					return
				}
				if selectedUserID == userID {
					filteredPosts = append(filteredPosts, post)
					break
				}
			}
		} else {
			// For private posts, include them only if the user is following the post's author
			isFollowing, err := app.database.IsFollowing(userID, post.UserID)
			if err != nil {
				app.errorJSON(w, fmt.Errorf("Error checking if the user is following the post author"), http.StatusInternalServerError)
				return
			}
			if isFollowing {
				filteredPosts = append(filteredPosts, post)
			}
		}
	}

	//comments for each post and add them to the filtered posts
	for i := range filteredPosts {
		postID := filteredPosts[i].PostID
		comments, err := app.database.GetCommentsByPostID(postID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Error getting comments from the database"), http.StatusInternalServerError)
			return
		}

		filteredPosts[i].Comments = comments
	}

	_ = app.writeJSON(w, http.StatusOK, filteredPosts)
}

func (app *application) CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/create-comment" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var comment models.Comment
	err := app.readJSON(w, r, &comment)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	userId, _, firstName, lastName, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from user sessions"), http.StatusInternalServerError)
		return
	}

	comment.UserID = userId
	comment.FirstName = firstName
	comment.LastName = lastName

	err = app.database.CreateComment(&comment)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error adding data to the database"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, comment)
}

func (app *application) ProfileTypeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/profile-type" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	userId, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		return
	}

	err = app.database.UpdateProfileType(userId)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to update the profile type"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, map[string]string{"message": "Profile type updated successfully"})
}

func (app *application) FollowHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/follow" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var request models.FollowRequest
	err := app.readJSON(w, r, &request)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	userId, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		return
	}

	isPublic, err := app.database.IsUserPublic(request.FollowingID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get user's public status"), http.StatusInternalServerError)
		return
	}

	isFollowing, err := app.database.IsFollowing(userId, request.FollowingID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to check if user is following"), http.StatusInternalServerError)
		return
	}

	if isFollowing {
		err = app.database.UnfollowUser(userId, request.FollowingID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to unfollow user: %w", err), http.StatusInternalServerError)
			return
		}
	} else if isPublic {
		err = app.database.FollowUser(userId, request.FollowingID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to follow user: %w", err), http.StatusInternalServerError)
			return
		}
	} else {
		err = app.database.FollowNotPublicUser(userId, request.FollowingID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to follow user: %w", err), http.StatusInternalServerError)
			return
		}
	}
	_ = app.writeJSON(w, http.StatusOK, request)
}

func (app *application) FollowingHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/following" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	userId, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		return
	}

	following, err := app.database.Following(userId)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the list of followed users: %w", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		Following []models.UserData `json:"following_users"`
	}{
		Following: following,
	}

	_ = app.writeJSON(w, http.StatusOK, response)
}

func (app *application) FollowersHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/followers" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	userId, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		return
	}

	followers, err := app.database.Followers(userId)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the list of followed users: %w", err), http.StatusInternalServerError)
		return
	}
	response := struct {
		Followers []models.UserData `json:"followers_users"`
	}{
		Followers: followers,
	}

	_ = app.writeJSON(w, http.StatusOK, response)
}

func (app *application) FollowRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/follow-requests" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	userID, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		userID = 0
		return
	}

	followRequests, err := app.database.FollowRequests(userID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get follow requests"), http.StatusInternalServerError)
		return
	}

	var usersData []*models.UserData
	for _, request := range followRequests {
		user, err := app.database.GetUserByID(request.FollowingID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to get user data for follower ID: %d", request.FollowerID), http.StatusInternalServerError)
			return
		}
		usersData = append(usersData, user)
	}

	_ = app.writeJSON(w, http.StatusOK, usersData)
}

func (app *application) AcceptFollowerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/accept-follower" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	userID, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		return
	}

	var request models.FollowRequest
	err = app.readJSON(w, r, &request)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	err = app.database.AcceptFollower(userID, request.FollowerID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to update follower status"), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Follower request accepted successfully"}
	_ = app.writeJSON(w, http.StatusOK, response)
}

func (app *application) DeclineFollowerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/decline-follower" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	userID, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		return
	}

	var request models.FollowRequest
	err = app.readJSON(w, r, &request)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	err = app.database.DeclineFollower(userID, request.FollowerID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to decline follower request"), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Follower request declined successfully"}
	_ = app.writeJSON(w, http.StatusOK, response)
}

func (app *application) CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/create-group" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var group models.Group
	err := app.readJSON(w, r, &group)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	userId, _, firstName, lastName, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from user sessions"), http.StatusInternalServerError)
		return
	}

	group.UserID = userId
	group.FirstName = firstName
	group.LastName = lastName

	groupID, err := app.database.CreateGroup(&group)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error adding data to the database"), http.StatusInternalServerError)
		return
	}

	var selectedUserIDs []string

	if len(group.SelectedUserID) > 0 {
		selectedUserIDs = strings.Split(group.SelectedUserID, ",")
	}

	// Add selected users to the grouprequests table
	for _, userIDStr := range selectedUserIDs {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Error converting string to int"), http.StatusInternalServerError)
			continue
		}
		groupMembers := models.GroupMembers{
			GroupID:        groupID,
			GroupTitle:     group.Title,
			GroupCreatorID: group.UserID,
			RequesterID:    userID,
		}

		err = app.database.AddGroupMembers(&groupMembers)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Error adding data to the database"), http.StatusInternalServerError)
			return
		}
	}

	_ = app.writeJSON(w, http.StatusOK, group)
}

func (app *application) AllGroupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/all-groups" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var allGroups []models.Group

	allGroups, err := app.database.AllGroups()
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from the database"), http.StatusInternalServerError)
		return
	}
	_ = app.writeJSON(w, http.StatusOK, allGroups)
}

func (app *application) GroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/group/")
	id1, err := strconv.Atoi(id)

	userID, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from user sessions"), http.StatusInternalServerError)
		return
	}

	group, err := app.database.GetGroup(id1)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	groupMembers, err := app.database.GetGroupMembers(id1)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var usersData []*models.UserData
	for _, memberID := range groupMembers {
		user, err := app.database.GetUserByID(memberID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to get user data"), http.StatusInternalServerError)
			return
		}
		usersData = append(usersData, user)
	}

	type GroupResponse struct {
		UserID       int                `json:"userID"`
		Group        *models.Group      `json:"group"`
		GroupMembers []int              `json:"group_members"`
		UserData     []*models.UserData `json:"userdata"`
	}

	groupResponse := GroupResponse{
		UserID:       userID,
		Group:        group,
		GroupMembers: groupMembers,
		UserData:     usersData,
	}

	app.writeJSON(w, http.StatusOK, groupResponse)
}

func (app *application) InviteNewMemberHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/invite" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var group models.Group
	err := app.readJSON(w, r, &group)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	groupID, err := app.database.CreateGroup(&group)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error adding data to the database"), http.StatusInternalServerError)
		return
	}

	var selectedUserID int

	groupMembers := models.GroupMembers{
		GroupID:        groupID,
		GroupTitle:     group.Title,
		GroupCreatorID: group.UserID,
		RequesterID:    selectedUserID,
	}

	err = app.database.InviteNewMember(&groupMembers)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error adding data to the database"), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Group invitation sent"}
	_ = app.writeJSON(w, http.StatusOK, response)
}

func (app *application) GroupInvitationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/group-invitations" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	userID, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		userID = 0
		return
	}

	groupInvitations, err := app.database.GroupInvitations(userID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get group requests"), http.StatusInternalServerError)
		return
	}

	type GroupInvitationWithUserData struct {
		GroupID        int              `json:"group_id"`
		GroupTitle     string           `json:"group_title"`
		GroupCreatorID int              `json:"group_creator_id"`
		InvitedUser    *models.UserData `json:"invited_user"`
	}

	var groupInvitationsWithUserData []GroupInvitationWithUserData

	for _, invitation := range groupInvitations {
		user, err := app.database.GetUserByID(invitation.RequesterID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to get user data"), http.StatusInternalServerError)
			return
		}

		invitationData := GroupInvitationWithUserData{
			GroupID:        invitation.GroupID,
			GroupTitle:     invitation.GroupTitle,
			GroupCreatorID: invitation.GroupCreatorID,
			InvitedUser:    user,
		}

		groupInvitationsWithUserData = append(groupInvitationsWithUserData, invitationData)
	}

	_ = app.writeJSON(w, http.StatusOK, groupInvitationsWithUserData)
}

func (app *application) AcceptGroupInvitationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/accept-group-invitation" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var invitation models.GroupMembers
	err := app.readJSON(w, r, &invitation)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	err = app.database.AcceptGroupInvitation(invitation.GroupID, invitation.RequesterID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to update invitation status"), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Group invitation accepted successfully"}
	_ = app.writeJSON(w, http.StatusOK, response)
}

func (app *application) DeclineGroupInvitationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/decline-group-invitation" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var invitation models.GroupMembers
	err := app.readJSON(w, r, &invitation)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	err = app.database.DeclineGroupInvitation(invitation.GroupID, invitation.RequesterID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to update invitation status"), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Group invitation declined successfully"}
	_ = app.writeJSON(w, http.StatusOK, response)
}

func (app *application) RequestToJoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/request-to-join-group" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var request models.GroupMembers
	err := app.readJSON(w, r, &request)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	group, err := app.database.GetGroup(request.GroupID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get group data"), http.StatusInternalServerError)
		return
	}

	groupTitle := group.Title
	groupCreatorID := group.UserID

	userId, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		return
	}

	isMember, err := app.database.IsMember(userId, request.GroupID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to check if user is following"), http.StatusInternalServerError)
		return
	}

	if isMember {
		err = app.database.LeaveGroup(userId, request.GroupID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to unfollow user: %w", err), http.StatusInternalServerError)
			return
		}
	} else {
		err = app.database.JoinGroup(userId, request.GroupID, groupTitle, groupCreatorID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to follow user: %w", err), http.StatusInternalServerError)
			return
		}
	}
	_ = app.writeJSON(w, http.StatusOK, request)
}

func (app *application) GroupRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/group-requests" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	userID, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		userID = 0
		return
	}

	groupRequests, err := app.database.GroupRequests(userID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get group requests"), http.StatusInternalServerError)
		return
	}

	type GroupRequestWithUserData struct {
		GroupID        int              `json:"group_id"`
		GroupTitle     string           `json:"group_title"`
		GroupCreatorID int              `json:"group_creator_id"`
		Requester      *models.UserData `json:"requester"`
	}

	var groupRequestsWithUserData []GroupRequestWithUserData

	for _, request := range groupRequests {
		user, err := app.database.GetUserByID(request.RequesterID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to get user data for requester ID: %d", request.RequesterID), http.StatusInternalServerError)
			return
		}

		requestData := GroupRequestWithUserData{
			GroupID:        request.GroupID,
			GroupTitle:     request.GroupTitle,
			GroupCreatorID: request.GroupCreatorID,
			Requester:      user,
		}

		groupRequestsWithUserData = append(groupRequestsWithUserData, requestData)
	}

	_ = app.writeJSON(w, http.StatusOK, groupRequestsWithUserData)
}

func (app *application) AcceptGroupRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/accept-group-request" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var request models.GroupMembers
	err := app.readJSON(w, r, &request)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	err = app.database.AcceptGroupRequest(request.GroupID, request.RequesterID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to update request status"), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Group request accepted successfully"}
	_ = app.writeJSON(w, http.StatusOK, response)
}

func (app *application) DeclineGroupRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/decline-group-request" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var request models.GroupMembers
	err := app.readJSON(w, r, &request)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	err = app.database.DeclineGroupRequest(request.GroupID, request.RequesterID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to update request status"), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Group request declined successfully"}
	_ = app.writeJSON(w, http.StatusOK, response)
}

func (app *application) AddMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/message" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var message models.Message
	err := app.readJSON(w, r, &message)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	message = models.Message{
		Message:       message.Message,
		FirstNameFrom: message.FirstNameFrom,
		FirstNameTo:   message.FirstNameTo,
		Date:          time.Now(),
	}
	_, _, message.FirstNameFrom, _, err = app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the name"), http.StatusInternalServerError)
		return
	}

	if message.Message != "" {
		err = app.database.AddMessage(message.Message, message.FirstNameFrom, message.FirstNameTo, message.Date)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to add message"), http.StatusInternalServerError)
			return
		}
	}
	_ = app.writeJSON(w, http.StatusCreated, message)
}

func (app *application) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/messages" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	firstNameTo := r.URL.Query().Get("firstNameTo")
	_, _, firstNameFrom, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the name"), http.StatusInternalServerError)
		return
	}

	messages, err := app.database.GetMessage(firstNameFrom, firstNameTo)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get messages"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, messages)
}

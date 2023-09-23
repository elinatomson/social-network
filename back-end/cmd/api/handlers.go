package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"social-network/back-end/models"

	"github.com/gorilla/websocket"
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

	// Parse the multipart form data to handle file uploads
	err := r.ParseMultipartForm(10) // 10 MB max file size
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error parsing form data"), http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	dateOfBirth := r.FormValue("date_of_birth")
	nickname := r.FormValue("nickname")
	aboutMe := r.FormValue("about_me")

	avatarFile, _, err := r.FormFile("avatar")
	var avatarFileName string
	if err != nil {
		avatarFileName = ""
	} else {
		defer avatarFile.Close()

		avatarFolderPath := "database/images/"
		avatarFileName = firstName + lastName + ".jpg"
		avatarFileData, err := ioutil.ReadAll(avatarFile)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Error reading avatar file"), http.StatusInternalServerError)
			return
		}

		err = ioutil.WriteFile(avatarFolderPath+avatarFileName, avatarFileData, 0644)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Error saving avatar file"), http.StatusInternalServerError)
			return
		}
	}

	userData := models.UserData{
		Email:       email,
		Password:    password,
		FirstName:   firstName,
		LastName:    lastName,
		DateOfBirth: dateOfBirth,
		Avatar:      avatarFileName,
		Nickname:    nickname,
		AboutMe:     aboutMe,
	}

	_, err = app.database.CheckEmail(userData.Email)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Email already taken"), http.StatusConflict)
		return
	}

	err = app.database.Register(&userData)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Email already taken"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, userData)

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
		app.errorJSON(w, fmt.Errorf("Email or password is not correct!"), http.StatusUnauthorized)
		return
	} else {
		userId, email, firstName, lastName, err := app.database.DataFromUserData(&userData)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Error getting data from user data"), http.StatusInternalServerError)
			return
		}
		cookieValue := app.addCookie(w, userId, email, firstName, lastName)

		app.writeJSON(w, http.StatusOK, map[string]string{"session": cookieValue})
	}
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

	userID, email, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the email"), http.StatusInternalServerError)
		return
	}

	userData, err := app.database.GetUserDataByEmail(email)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get user data"), http.StatusInternalServerError)
		return
	}

	allPosts, err := app.database.ProfilePosts(userID)
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
		app.errorJSON(w, fmt.Errorf("Failed to get user data"), http.StatusInternalServerError)
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
	userID, _, firstName, lastName, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from user sessions"), http.StatusInternalServerError)
		return
	}

	followers, err := app.database.Followers(userID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from the database"), http.StatusInternalServerError)
		return
	}

	followings, err := app.database.Following(userID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from the database"), http.StatusInternalServerError)
		return
	}

	//initalizing a map to store unique users
	uniqueUsers := make(map[int]models.UserData)

	for _, follower := range followers {
		//only those followers whose profile is public
		if follower.Public {
			uniqueUsers[follower.UserID] = follower
		}
	}

	for _, following := range followings {
		uniqueUsers[following.UserID] = following
	}

	var users []models.UserData
	for _, user := range uniqueUsers {
		users = append(users, user)
	}

	//setting the currentUser in the users db table as a true to add the current user's name to the response
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
		app.errorJSON(w, fmt.Errorf("Error getting user from the database"), http.StatusInternalServerError)
		return
	}

	allPosts, err := app.database.GetPostsByUserID(user.UserID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from the database"), http.StatusInternalServerError)
		return
	}

	userID, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from user sessions"), http.StatusInternalServerError)
		return
	}

	followers, err := app.database.Followers(user.UserID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from the database"), http.StatusInternalServerError)
		return
	}

	following, err := app.database.Following(user.UserID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from the database"), http.StatusInternalServerError)
		return
	}

	var filteredPosts []models.Post
	for i := range allPosts {
		post := allPosts[i]
		if post.GroupID == 0 && post.Privacy == "public" {
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
			if post.GroupID == 0 && isFollowing {
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
		CurrentUser int               `json:"current_user"`
		UserData    *models.UserData  `json:"user_data"`
		Followers   []models.UserData `json:"followers"`
		Following   []models.UserData `json:"following"`
		Posts       []models.Post     `json:"posts"`
	}{
		CurrentUser: userID,
		UserData:    user,
		Followers:   followers,
		Following:   following,
		Posts:       filteredPosts,
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

	err := r.ParseMultipartForm(10)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error parsing form data"), http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	privacy := r.FormValue("privacy")
	selectedUserID := r.FormValue("selected_user_id")
	groupID := r.FormValue("group_id")
	var groupIDInt int
	if groupID != "" && groupID != "undefined" {
		groupIDInt, err = strconv.Atoi(groupID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Error converting string to integer"), http.StatusInternalServerError)
			return
		}
	} else {
		groupIDInt = 0
	}

	imageFile, _, err := r.FormFile("image")
	var imageFileName string
	if err != nil {
		imageFileName = ""
	} else {
		defer imageFile.Close()

		imageFolderPath := "database/images/"
		//generating a random image name
		randomBytes := make([]byte, 16)
		_, err := rand.Read(randomBytes)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Error generating random image name"), http.StatusInternalServerError)
			return
		}

		//converting random bytes to a hexadecimal string
		imageName := hex.EncodeToString(randomBytes) + ".jpg"
		imageFileName = imageName

		imageFileData, err := ioutil.ReadAll(imageFile)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Error reading image file"), http.StatusInternalServerError)
			return
		}

		err = ioutil.WriteFile(imageFolderPath+imageFileName, imageFileData, 0644)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Error saving image file"), http.StatusInternalServerError)
			return
		}
	}

	post := models.Post{
		Content:        content,
		Privacy:        privacy,
		SelectedUserID: selectedUserID,
		GroupID:        groupIDInt,
		Image:          imageFileName,
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
	//including an empty comments array for the newly created post.
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

	var filteredPosts []models.Post

	for _, post := range allPosts {
		if post.GroupID == 0 && (post.Privacy == "public" || post.UserID == userID) {
			filteredPosts = append(filteredPosts, post)
		} else if post.Privacy == "for-selected-users" {
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
			if post.GroupID == 0 && isFollowing {
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

	err := r.ParseMultipartForm(10)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error parsing form data"), http.StatusBadRequest)
		return
	}

	commentContent := r.FormValue("comment")
	postID := r.FormValue("post_id")
	var postIDInt int
	postIDInt, err = strconv.Atoi(postID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error converting string to integer"), http.StatusInternalServerError)
		return
	}

	imageFile, _, err := r.FormFile("image")
	var imageFileName string
	if err != nil {
		imageFileName = ""
	} else {
		defer imageFile.Close()

		imageFolderPath := "database/images/"
		randomBytes := make([]byte, 16)
		_, err := rand.Read(randomBytes)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Error generating random image name"), http.StatusInternalServerError)
			return
		}

		imageName := hex.EncodeToString(randomBytes) + ".jpg"
		imageFileName = imageName

		imageFileData, err := ioutil.ReadAll(imageFile)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Error reading image file"), http.StatusInternalServerError)
			return
		}

		err = ioutil.WriteFile(imageFolderPath+imageFileName, imageFileData, 0644)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Error saving image file"), http.StatusInternalServerError)
			return
		}
	}

	comment := models.Comment{
		PostID:  postIDInt,
		Comment: commentContent,
		Image:   imageFileName,
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

	_ = app.writeJSON(w, http.StatusOK, err)
}

func (app *application) FollowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

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

	isPending, err := app.database.IsPending(userId, request.FollowingID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from the database"), http.StatusInternalServerError)
		return
	}

	if isFollowing || isPending {
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

func (app *application) FollowerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/follower-check" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	followingId := r.URL.Query().Get("userId")
	followingIdInt, err := strconv.Atoi(followingId)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Invalid groupId parameter"), http.StatusBadRequest)
		return
	}

	userId, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		return
	}

	isFollowing, err := app.database.IsFollowing(userId, followingIdInt)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to check if user is following"), http.StatusInternalServerError)
		return
	}

	isPending, err := app.database.IsPending(userId, followingIdInt)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from the database"), http.StatusInternalServerError)
		return
	}

	followData := struct {
		IsFollowing bool `json:"is_following"`
		IsPending   bool `json:"is_pending"`
	}{
		IsFollowing: isFollowing,
		IsPending:   isPending,
	}

	_ = app.writeJSON(w, http.StatusOK, followData)
}

func (app *application) FollowingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/following" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	userId, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		return
	}

	var following []models.UserData

	following, err = app.database.Following(userId)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the list of followed users: %w", err), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, following)
}

func (app *application) FollowersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/followers" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	userId, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		return
	}

	var followers []models.UserData

	followers, err = app.database.Followers(userId)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the list of followed users: %w", err), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, followers)
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

	_ = app.writeJSON(w, http.StatusOK, request)
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

	_ = app.writeJSON(w, http.StatusOK, request)
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

	//adding selected users to the groupmembers table
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
			MemberID:       userID,
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

	userID, _, firstName, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from user sessions"), http.StatusInternalServerError)
		return
	}

	group, err := app.database.GetGroup(id1)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting group data from database"), http.StatusInternalServerError)
		return
	}

	groupMembers, err := app.database.GetGroupMembers(id1)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting group members from database"), http.StatusInternalServerError)
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

	requestPending, err := app.database.CheckPending(userID, group.GroupID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error checking pending status"), http.StatusInternalServerError)
		return
	}

	type GroupResponse struct {
		UserID         int                `json:"userID"`
		CurrentUser    string             `json:"current_user"`
		Group          *models.Group      `json:"group"`
		GroupMembers   []int              `json:"group_members"`
		UserData       []*models.UserData `json:"userdata"`
		RequestPending bool               `json:"request_pending"`
	}

	groupResponse := GroupResponse{
		UserID:         userID,
		CurrentUser:    firstName,
		Group:          group,
		GroupMembers:   groupMembers,
		UserData:       usersData,
		RequestPending: requestPending,
	}

	app.writeJSON(w, http.StatusOK, groupResponse)
}

func (app *application) GroupPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/group-posts" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	groupId := r.URL.Query().Get("groupId")
	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Invalid groupId parameter"), http.StatusBadRequest)
		return
	}

	var allPosts []models.Post

	allPosts, err = app.database.AllPosts()
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from the database"), http.StatusInternalServerError)
		return
	}

	var filteredPosts []models.Post

	for _, post := range allPosts {
		if post.GroupID == groupIdInt {
			filteredPosts = append(filteredPosts, post)
		}
	}

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

func (app *application) InviteNewMemberHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/invite" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var groupMembers models.GroupMembers
	err := app.readJSON(w, r, &groupMembers)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	groupData, err := app.database.GetGroup(groupMembers.GroupID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get group data"), http.StatusInternalServerError)
		return
	}

	groupMembers.GroupTitle = groupData.Title
	groupMembers.GroupCreatorID = groupData.UserID

	err = app.database.InviteNewMember(&groupMembers)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error adding data to the database"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, err)
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
		user, err := app.database.GetUserByID(invitation.MemberID)
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

	err = app.database.AcceptGroupInvitation(invitation.GroupID, invitation.MemberID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to update invitation status"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, err)
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

	err = app.database.DeclineGroupInvitation(invitation.GroupID, invitation.MemberID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to update invitation status"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, err)
}

func (app *application) RequestToJoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

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
		Member         *models.UserData `json:"member"`
	}

	var groupRequestsWithUserData []GroupRequestWithUserData

	for _, request := range groupRequests {
		user, err := app.database.GetUserByID(request.MemberID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to get user data for member ID: %d", request.MemberID), http.StatusInternalServerError)
			return
		}

		requestData := GroupRequestWithUserData{
			GroupID:        request.GroupID,
			GroupTitle:     request.GroupTitle,
			GroupCreatorID: request.GroupCreatorID,
			Member:         user,
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

	err = app.database.AcceptGroupRequest(request.GroupID, request.MemberID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to update request status"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, err)
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

	err = app.database.DeclineGroupRequest(request.GroupID, request.MemberID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to update request status"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, err)
}

func (app *application) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/create-event" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var event models.Event
	err := app.readJSON(w, r, &event)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	userId, _, firstName, lastName, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from user sessions"), http.StatusInternalServerError)
		return
	}

	event.UserID = userId
	event.FirstName = firstName
	event.LastName = lastName

	eventID, err := app.database.CreateEvent(&event)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error adding data to the database"), http.StatusInternalServerError)
		return
	}

	groupMembers, err := app.database.GetGroupMembers(event.GroupID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting group members from database"), http.StatusInternalServerError)
		return
	}

	for _, memberID := range groupMembers {
		err = app.database.EventNotifications(eventID, memberID, event.GroupID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Error adding data to the database"), http.StatusInternalServerError)
			return
		}
	}

	groupCreatorID, err := app.database.GetGroupCreator(event.GroupID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting group members from database"), http.StatusInternalServerError)
		return
	}

	err = app.database.EventNotifications(eventID, groupCreatorID, event.GroupID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error adding data to the database"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, event)
}

func (app *application) GroupEventNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/group-event-notifications" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	userID, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		userID = 0
		return
	}

	eventNotifications, err := app.database.GetEventNotifications(userID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get event notifications"), http.StatusInternalServerError)
		return
	}

	type GroupRequestWithUserData struct {
		EventID    int    `json:"event_id"`
		EventTitle string `json:"event_title"`
		GroupID    int    `json:"group_id"`
		GroupTitle string `json:"group_title"`
	}

	var eventNotificationWithGroupData []GroupRequestWithUserData

	for _, notification := range eventNotifications {
		group, err := app.database.GetGroup(notification.GroupID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to get group data for group ID"), http.StatusInternalServerError)
			return
		}

		event, err := app.database.GetEvent(notification.EventID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to get group data for group ID"), http.StatusInternalServerError)
			return
		}

		requestData := GroupRequestWithUserData{
			EventID:    notification.EventID,
			EventTitle: event.Title,
			GroupID:    notification.GroupID,
			GroupTitle: group.Title,
		}

		eventNotificationWithGroupData = append(eventNotificationWithGroupData, requestData)
	}

	_ = app.writeJSON(w, http.StatusOK, eventNotificationWithGroupData)
}

func (app *application) EventSeenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/group-event-seen" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	userID, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		return
	}

	var eventNotification models.EventNotifications
	err = app.readJSON(w, r, &eventNotification)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	err = app.database.DeleteFromEventNotifications(eventNotification.EventID, userID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to delete from database"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, eventNotification)
}

func (app *application) GroupEventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/group-events" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	groupId := r.URL.Query().Get("groupId")
	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Invalid groupId parameter"), http.StatusBadRequest)
		return
	}

	var allEvents []models.Event

	allEvents, err = app.database.AllEvents()
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from the database"), http.StatusInternalServerError)
		return
	}

	var filteredEvents []models.Event

	for _, event := range allEvents {
		if event.GroupID == groupIdInt {
			filteredEvents = append(filteredEvents, event)
		}
	}

	_ = app.writeJSON(w, http.StatusOK, filteredEvents)
}

func (app *application) GroupEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/group-event/")
	id1, err := strconv.Atoi(id)

	userID, _, _, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from user sessions"), http.StatusInternalServerError)
		return
	}

	event, err := app.database.GetEvent(id1)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	isGroupMember, err := app.database.CheckMembership(userID, event.GroupID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	isGroupCreator, err := app.database.CheckCreator(userID, event.GroupID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	participants, err := app.database.GetParticipants(id1)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	isGoing, err := app.database.IsGoing(userID, id1)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to check if user is going"), http.StatusInternalServerError)
		return
	}

	isNotGoing, err := app.database.IsNotGoing(userID, id1)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to check if user is going"), http.StatusInternalServerError)
		return
	}

	response := struct {
		IsGroupMember  bool                       `json:"is_group_member"`
		IsGroupCreator bool                       `json:"is_group_creator"`
		Event          *models.Event              `json:"event"`
		Participants   []models.EventParticipants `json:"participants"`
		Going          bool                       `json:"going"`
		NotGoing       bool                       `json:"not_going"`
	}{
		IsGroupMember:  isGroupMember,
		IsGroupCreator: isGroupCreator,
		Event:          event,
		Participants:   participants,
		Going:          isGoing,
		NotGoing:       isNotGoing,
	}

	app.writeJSON(w, http.StatusOK, response)
}

func (app *application) GoingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/going" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var going models.EventParticipants
	err := app.readJSON(w, r, &going)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	userId, _, firstName, lastName, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		return
	}

	isNotGoing, err := app.database.IsNotGoing(userId, going.EventID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to check if user is going"), http.StatusInternalServerError)
		return
	}

	if isNotGoing {
		err = app.database.NotGoingToGoingEvent(userId, going.EventID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to mark as going: %w", err), http.StatusInternalServerError)
			return
		}
	} else {
		err = app.database.GoingToEvent(userId, going.EventID, firstName, lastName)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to mark as going: %w", err), http.StatusInternalServerError)
			return
		}
	}

	_ = app.writeJSON(w, http.StatusOK, going)
}

func (app *application) NotGoingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/not-going" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	var notGoing models.EventParticipants
	err := app.readJSON(w, r, &notGoing)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error decoding JSON data"), http.StatusBadRequest)
		return
	}

	userId, _, firstName, lastName, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the user ID from the session"), http.StatusInternalServerError)
		return
	}

	isGoing, err := app.database.IsGoing(userId, notGoing.EventID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to check if user is going"), http.StatusInternalServerError)
		return
	}

	if isGoing {
		err = app.database.GoingToNotGoingEvent(userId, notGoing.EventID)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to mark as not going: %w", err), http.StatusInternalServerError)
			return
		}
	} else {
		err = app.database.NotGoingToEvent(userId, notGoing.EventID, firstName, lastName)
		if err != nil {
			app.errorJSON(w, fmt.Errorf("Failed to mark as not going: %w", err), http.StatusInternalServerError)
			return
		}
	}

	_ = app.writeJSON(w, http.StatusOK, notGoing)
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

	if r.URL.Path != "/conversation-history/" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	firstNameTo := r.URL.Query().Get("firstNameTo")
	_, _, firstNameFrom, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the name"), http.StatusInternalServerError)
		return
	}

	messages, err := app.database.GetMessages(firstNameFrom, firstNameTo)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get messages"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, messages)
}

func (app *application) GetGroupMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/group-conversation-history/" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	groupName := r.URL.Query().Get("groupName")

	messages, err := app.database.GetGroupMessages(groupName)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get messages"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, messages)
}

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		allowedOrigin := "http://localhost:3000"
		return r.Header.Get("Origin") == allowedOrigin
	},
	}
	connections = make(map[string]*websocket.Conn)
	mutex       = sync.Mutex{}
)

func (app *application) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection in order to enable full-duplex communication and support WebSocket-specific features
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	_, _, firstName, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the name"), http.StatusInternalServerError)
		return
	}

	// Add the WebSocket connection to the connections map to maintain active WebSocket connections.
	mutex.Lock()
	connections[firstName] = conn
	mutex.Unlock()
	// The function enters a loop to continuously read messages from the client.
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// Convert the byte slice to a string
		messageStr := string(message)

		// Unmarshal the string into a Message struct
		var msg models.Message
		err = json.Unmarshal([]byte(messageStr), &msg)
		if err != nil {
			log.Println("Failed to unmarshal message:", err)
			break
		}

		// Calling the handleMessage function, passing the recipient user's name, writer user's name, and the message as parameters to handle the received message.
		app.handleMessage(r, w, msg.FirstNameFrom, msg.FirstNameTo, msg)
	}

	// Remove the WebSocket connection from the connections map when the connection is closed
	mutex.Lock()
	delete(connections, firstName)
	mutex.Unlock()
}

func (app *application) handleMessage(r *http.Request, w http.ResponseWriter, senderFirstName string, receiverFirstName string, message models.Message) {
	// Check if the recipient user has an active WebSocket connection
	mutex.Lock()
	recipientConn, recipientFound := connections[receiverFirstName]
	mutex.Unlock()
	// Check if the sender user has an active WebSocket connection
	_, _, senderFirstName, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the name"), http.StatusInternalServerError)
		return
	}
	mutex.Lock()
	senderConn, senderFound := connections[senderFirstName]
	mutex.Unlock()

	chatMessage := models.Message{
		Message:       message.Message,
		FirstNameFrom: senderFirstName,
		FirstNameTo:   receiverFirstName,
		Date:          message.Date,
	}

	if recipientFound {
		// Send the message to the recipient user's WebSocket connection
		data, err := json.Marshal(chatMessage)
		if err != nil {
			log.Println("Failed to marshal message:", err)
			return
		}
		err = recipientConn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Failed to write message to recipient:", err)
		}
	} else {
		log.Println("No active WebSocket connection found for recipient:", receiverFirstName)
	}
	if senderFound {
		// Send the message to the sender's WebSocket connection
		data, err := json.Marshal(chatMessage)
		if err != nil {
			log.Println("Failed to marshal message:", err)
			return
		}
		err = senderConn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Failed to write message to sender:", err)
		}
	} else {
		log.Println("No active WebSocket connection found for sender:", senderFirstName)
	}
}

var (
	upgraderGroup = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			allowedOrigin := "http://localhost:3000"
			return r.Header.Get("Origin") == allowedOrigin
		},
	}
	// Use a map to maintain active WebSocket connections for group chats.
	groupConnections = make(map[string]map[*websocket.Conn]bool)
	groupMutex       = sync.Mutex{}
)

func (app *application) GroupWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgraderGroup.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	groupName := r.URL.Query().Get("group")

	// Lock and add the WebSocket connection to the group's connection map
	groupMutex.Lock()
	if groupConnections[groupName] == nil {
		groupConnections[groupName] = make(map[*websocket.Conn]bool)
	}
	groupConnections[groupName][conn] = true
	groupMutex.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}

		messageStr := string(message)

		var groupMsg models.Message
		err = json.Unmarshal([]byte(messageStr), &groupMsg)
		if err != nil {
			log.Println("Failed to unmarshal message:", err)
			break
		}

		broadcastGroupMessage(groupName, groupMsg)
	}

	// Remove the WebSocket connection from the group's connection map when the connection is closed
	groupMutex.Lock()
	delete(groupConnections[groupName], conn)
	groupMutex.Unlock()
	conn.Close()
}

func broadcastGroupMessage(groupName string, message models.Message) {
	groupMutex.Lock()
	defer groupMutex.Unlock()

	// Iterate over all WebSocket connections in the group and send the message
	for conn := range groupConnections[groupName] {
		data, err := json.Marshal(message)
		if err != nil {
			log.Println("Failed to marshal message:", err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Failed to write message to connection:", err)
			continue
		}
	}
}

func (app *application) UnreadMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/unread-messages" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	_, _, firstName, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the name"), http.StatusInternalServerError)
		return
	}

	unreadMessages, err := app.database.GetUnreadMessages(firstName)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get messages"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, unreadMessages)
}

func (app *application) MarkMessagesAsReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorJSON(w, fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/mark-messages-as-read/" {
		app.errorJSON(w, fmt.Errorf("Error 404, page not found"), http.StatusNotFound)
		return
	}

	firstNameFrom := r.URL.Query().Get("firstNameFrom")
	_, _, firstNameto, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the name"), http.StatusInternalServerError)
		return
	}

	err = app.database.MarkMessagesAsRead(firstNameto, firstNameFrom)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to update messages"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, err)
}

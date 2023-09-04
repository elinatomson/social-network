package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io/ioutil"
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
	_, _, firstName, lastName, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting data from user sessions"), http.StatusInternalServerError)
		return
	}

	users, err := app.database.GetUsers()
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

	err := r.ParseMultipartForm(10) // 10 MB max file size
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

	err := r.ParseMultipartForm(10) // 10 MB max file size
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

	// Add selected users to the groupmembers table
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

	userID, _, _, _, err := app.database.DataFromSession(r)
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

	response := struct {
		IsGroupMember  bool                       `json:"is_group_member"`
		IsGroupCreator bool                       `json:"is_group_creator"`
		Event          *models.Event              `json:"event"`
		Participants   []models.EventParticipants `json:"participants"`
	}{
		IsGroupMember:  isGroupMember,
		IsGroupCreator: isGroupCreator,
		Event:          event,
		Participants:   participants,
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

package main

import (
	"database/sql"
	"fmt"
	"net/http"

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

	email, err := app.database.EmailFromUserData(&userData)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Error getting email from user data"), http.StatusInternalServerError)
		return
	}

	cookieValue := app.addCookie(w, email)

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

	app.writeJSON(w, http.StatusOK, map[string]string{})
}

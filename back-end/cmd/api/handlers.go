package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	//"time"
	"back-end/models"
)

type SqliteDB struct {
	DB *sql.DB
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Social Network up and running",
		Version: "1.0.0",
	}
	out, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *application) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse(w, http.StatusMethodNotAllowed, nil, "Invalid request method")
		return
	}

	if r.URL.Path != "/register" {
		jsonResponse(w, http.StatusNotFound, nil, "Error 404, page not found")
		return
	}

	var userData models.UserData
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, nil, "Error decoding JSON data")
		return
	}

	stmt := `SELECT email FROM users WHERE email = ?`
	row := app.database.DB.QueryRow(stmt, userData.Email)
	var email string
	err = row.Scan(&email)
	if err != sql.ErrNoRows {
		jsonResponse(w, http.StatusConflict, nil, "Email already taken")
		return
	}

	err = app.database.Register(&userData)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, nil, "Error adding data to the database")
		return
	}

	jsonResponse(w, http.StatusOK, userData, "")
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse(w, http.StatusMethodNotAllowed, nil, "Invalid request method")
		return
	}

	if r.URL.Path != "/login" {
		jsonResponse(w, http.StatusNotFound, nil, "Error 404, page not found")
		return
	}

	var userData models.UserData
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, nil, "Error decoding JSON data")
		return
	}

	err = app.database.Login(&userData)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, nil, "Email or password is not correct!")
		return
	}

	jsonResponse(w, http.StatusOK, userData, "")
}

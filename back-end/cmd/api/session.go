package main

import (
	"back-end/models"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

func (app *application) addCookie(w http.ResponseWriter, email string) {
	// Generate a new UUID for a session.
	uuid, _ := uuid.NewV4()
	value := uuid.String()
	expire := time.Now().Add(1 * time.Hour)
	cookie := http.Cookie{
		Name:    "sessionId",
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)

	session := &models.Session{
		Email:  email,
		Cookie: value,
	}

	err := app.database.Session(session)
	if err != nil {
		return
	}
}

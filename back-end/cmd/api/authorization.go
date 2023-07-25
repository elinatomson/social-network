package main

import (
	"back-end/models"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

func (app *application) addCookie(w http.ResponseWriter, email string) string {
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

	app.database.Session(session)

	return value
}

func (app *application) deleteCookie(r *http.Request) error {
	cookie, err := r.Cookie("sessionId")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil
		}
		return err
	}
	uuid, err := uuid.FromString(cookie.Value)
	if err != nil {
		return err
	}

	err = app.database.DeleteSession(uuid.String())
	if err != nil {
		return err
	}
	return nil
}

func (app *application) GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (int, error) {
	w.Header().Add("Vary", "Authorization")

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return 0, errors.New("no auth header")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return 0, errors.New("invalid auth header")
	}

	if headerParts[0] != "Bearer" {
		return 0, errors.New("invalid auth header")
	}

	cookie := headerParts[1]

	email, err := app.database.ValidateSession(cookie)
	if err != nil {
		return 0, err
	}

	return email, nil
}

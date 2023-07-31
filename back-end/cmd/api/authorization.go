package main

import (
	"back-end/models"
	"errors"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

func (app *application) addCookie(w http.ResponseWriter, userId int, email string, firstName string, lastName string) string {
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
		UserID:    userId,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Cookie:    value,
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

func (app *application) GetSessionIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("sessionId")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", errors.New("session cookie not found")
		}
		return "", err
	}

	return cookie.Value, nil
}

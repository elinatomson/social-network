package repository

import (
	"back-end/models"
	"context"
	"golang.org/x/crypto/bcrypt"
)

func (m *SqliteDB) Login(userData *models.UserData) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT password FROM users WHERE email = ?`
	row := m.DB.QueryRowContext(ctx, stmt, userData.Email)
	var hash string
	err := row.Scan(&hash)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(userData.Password))
	if err != nil {
		return err
	}

	return nil
}

func (m *SqliteDB) EmailFromUserData(userData *models.UserData) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT email FROM users WHERE email = ?`
	row := m.DB.QueryRowContext(ctx, stmt, userData.Email)

	var email string
	err := row.Scan(&email)
	if err != nil {
		return "", err
	}

	return email, nil
}




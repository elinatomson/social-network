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


	/*hash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = m.DB.ExecContext(ctx, stmt, userData.Email, hash, userData.FirstName, userData.LastName, userData.DateOfBirth, userData.Avatar, userData.Nickname, userData.AboutMe)
	if err != nil {
		return err
	}

	return nil*/


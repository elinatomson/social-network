package repository

import (
	"back-end/models"
	"context"
	"database/sql"
	"time"
	"golang.org/x/crypto/bcrypt"
)

type SqliteDB struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

func (m *SqliteDB) Connection() *sql.DB {
	return m.DB
}

func (m *SqliteDB) Register(userData *models.UserData) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = m.DB.ExecContext(ctx, stmt, userData.Email, hash, userData.FirstName, userData.LastName, userData.DateOfBirth, userData.Avatar, userData.Nickname, userData.AboutMe)
	if err != nil {
		return err
	}

	return nil
}

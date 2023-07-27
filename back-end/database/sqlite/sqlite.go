package sqlite

import (
	"back-end/models"
	"context"
	"database/sql"
	"time"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"github.com/gofrs/uuid"
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

func (m *SqliteDB) Session(session *models.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `INSERT INTO sessions (email, cookie) VALUES (?, ?)`

	_, err := m.DB.ExecContext(ctx, stmt, session.Email, session.Cookie)
	if err != nil {
		return err
	}
	return nil
}

func (m *SqliteDB) DeleteSession(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `DELETE FROM sessions WHERE cookie = ?`

	_, err := m.DB.ExecContext(ctx, stmt, uuid)
	if err != nil {
		return err
	}
	return nil
}

func (m *SqliteDB) EmailFromSession(r *http.Request) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	cookie, err := r.Cookie("sessionId")
	if err != nil {
		return "", err
	}

	uuid, err := uuid.FromString(cookie.Value)
	if err != nil {
		return "", err
	}

	stmt := `SELECT email FROM sessions WHERE cookie = ?`
	row := m.DB.QueryRowContext(ctx, stmt, uuid.String())
	var email string
	err = row.Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}

func (m *SqliteDB) GetUserDataByEmail(email string) (*models.UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT email, first_name, last_name, date_of_birth, avatar, nickname, about_me FROM users
		WHERE email = $1
		LIMIT 1
	`

	row := m.DB.QueryRowContext(ctx, stmt, email)
	userData := &models.UserData{}
	err := row.Scan(&userData.Email, &userData.FirstName, &userData.LastName, &userData.DateOfBirth, &userData.Avatar, &userData.Nickname, &userData.AboutMe)
	if err != nil {
		return nil, err
	}

	return userData, nil
}

func (m *SqliteDB) SearchUsers(query string) ([]models.UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT user_id, first_name, last_name, avatar FROM users WHERE first_name LIKE ? OR last_name LIKE ?`

	rows, err := m.DB.QueryContext(ctx, stmt, "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}

	var users []models.UserData
	for rows.Next() {
		var userData models.UserData
		err := rows.Scan(&userData.UserID, &userData.FirstName, &userData.LastName, &userData.Avatar)
		if err != nil {
			return nil, err
		}

		users = append(users, userData)
	}
	return users, nil
}
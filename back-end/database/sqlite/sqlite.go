package sqlite

import (
	"back-end/models"
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
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

func (m *SqliteDB) DataFromUserData(userData *models.UserData) (string, string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT email, first_name, last_name FROM users WHERE email = ?`
	row := m.DB.QueryRowContext(ctx, stmt, userData.Email)

	var email, firstName, lastName string
	err := row.Scan(&email, &firstName, &lastName)
	if err != nil {
		return "", "", "", err
	}

	return email, firstName, lastName, nil
}

func (m *SqliteDB) Session(session *models.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `INSERT INTO sessions (email, first_name, last_name, cookie) VALUES (?, ?, ?, ?)`

	_, err := m.DB.ExecContext(ctx, stmt, session.Email, session.FirstName, session.LastName, session.Cookie)
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

func (m *SqliteDB) DataFromSession(r *http.Request) (string, string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	cookie, err := r.Cookie("sessionId")
	if err != nil {
		return "", "", "", err
	}

	uuid, err := uuid.FromString(cookie.Value)
	if err != nil {
		return "", "", "", err
	}

	stmt := `SELECT email, first_name, last_name FROM sessions WHERE cookie = ?`
	row := m.DB.QueryRowContext(ctx, stmt, uuid.String())
	var email, firstName, lastName string
	err = row.Scan(&email, &firstName, &lastName)
	if err != nil {
		return "", "", "", err
	}
	return email, firstName, lastName, nil
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

	stmt := `SELECT user_id, first_name, last_name FROM users WHERE first_name LIKE ? OR last_name LIKE ?`

	rows, err := m.DB.QueryContext(ctx, stmt, "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}

	var users []models.UserData
	for rows.Next() {
		var userData models.UserData
		err := rows.Scan(&userData.UserID, &userData.FirstName, &userData.LastName)
		if err != nil {
			return nil, err
		}

		users = append(users, userData)
	}
	return users, nil
}

func (m *SqliteDB) GetUser(id int) (*models.UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT user_id, email, first_name, last_name, date_of_birth, avatar, nickname, about_me FROM users WHERE user_id = $1`

	row := m.DB.QueryRowContext(ctx, stmt, id)

	var user models.UserData

	err := row.Scan(
		&user.UserID, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.Nickname, &user.AboutMe,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *SqliteDB) CreatePost(post *models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	post.Date = time.Now()

	stmt := `INSERT INTO posts (content, first_name, last_name, privacy, image, date) VALUES (?, ?, ?, ?, ?, ?)`

	_, err := m.DB.ExecContext(ctx, stmt, post.Content, post.FirstName, post.LastName, post.Privacy, post.Image, post.Date)
	if err != nil {
		return err
	}

	return nil
}

func (m *SqliteDB) AllPosts() ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT post_id, content, first_name, last_name, image, date FROM posts`

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.PostID, &post.Content, &post.FirstName, &post.LastName, &post.Image, &post.Date)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (m *SqliteDB) UserPosts(firstName string, lastName string) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT post_id, content, first_name, last_name, image, date FROM posts WHERE first_name = ? AND last_name = ?`

	rows, err := m.DB.QueryContext(ctx, stmt, firstName, lastName)
	if err != nil {
		return nil, err
	}

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.PostID, &post.Content, &post.FirstName, &post.LastName, &post.Image, &post.Date)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

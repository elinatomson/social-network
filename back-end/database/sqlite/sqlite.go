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

func (m *SqliteDB) DataFromUserData(userData *models.UserData) (int, string, string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT user_id, email, first_name, last_name FROM users WHERE email = ?`
	row := m.DB.QueryRowContext(ctx, stmt, userData.Email)

	var userId int
	var email, firstName, lastName string
	err := row.Scan(&userId, &email, &firstName, &lastName)
	if err != nil {
		return 0, "", "", "", err
	}

	return userId, email, firstName, lastName, nil
}

func (m *SqliteDB) Session(session *models.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `INSERT INTO sessions (user_id, email, first_name, last_name, cookie) VALUES (?, ?, ?, ?, ?)`

	_, err := m.DB.ExecContext(ctx, stmt, session.UserID, session.Email, session.FirstName, session.LastName, session.Cookie)
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

func (m *SqliteDB) DataFromSession(r *http.Request) (int, string, string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	cookie, err := r.Cookie("sessionId")
	if err != nil {
		return 0, "", "", "", err
	}

	uuid, err := uuid.FromString(cookie.Value)
	if err != nil {
		return 0, "", "", "", err
	}

	stmt := `SELECT user_id, email, first_name, last_name FROM sessions WHERE cookie = ?`
	row := m.DB.QueryRowContext(ctx, stmt, uuid.String())
	var userId int
	var email, firstName, lastName string
	err = row.Scan(&userId, &email, &firstName, &lastName)
	if err != nil {
		return 0, "", "", "", err
	}
	return userId, email, firstName, lastName, nil
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

	stmt := `SELECT user_id, email, first_name, last_name, date_of_birth, avatar, nickname, about_me, public FROM users WHERE user_id = $1`

	row := m.DB.QueryRowContext(ctx, stmt, id)

	var user models.UserData

	err := row.Scan(
		&user.UserID, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.Nickname, &user.AboutMe, &user.Public,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *SqliteDB) GetUserByID(userID int) (*models.UserData, error) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    stmt := `SELECT user_id, first_name, last_name FROM users WHERE user_id = $1`

    row := m.DB.QueryRowContext(ctx, stmt, userID)

    var user models.UserData

    err := row.Scan(
        &user.UserID, &user.FirstName, &user.LastName,
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

	stmt := `INSERT INTO posts (user_id, content, first_name, last_name, privacy, image, date) VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := m.DB.ExecContext(ctx, stmt, post.UserID, post.Content, post.FirstName, post.LastName, post.Privacy, post.Image, post.Date)
	if err != nil {
		return err
	}

	return nil
}

func (m *SqliteDB) AllPosts() ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT post_id, user_id, content, first_name, last_name, image, date FROM posts`

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.PostID, &post.UserID, &post.Content, &post.FirstName, &post.LastName, &post.Image, &post.Date)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (m *SqliteDB) ProfilePosts(firstName string, lastName string) ([]models.Post, error) {
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

func (m *SqliteDB) GetPostsByUserID(userID int) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT post_id, user_id, content, first_name, last_name, privacy, image, date FROM posts WHERE user_id = ?`

	rows, err := m.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.PostID, &post.UserID, &post.Content, &post.FirstName, &post.LastName, &post.Privacy, &post.Image, &post.Date)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (m *SqliteDB) CreateComment(comment *models.Comment) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	comment.Date = time.Now()

	stmt := `INSERT INTO comments (post_id, user_id, comment, first_name, last_name, image, date) VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := m.DB.ExecContext(ctx, stmt, comment.PostID, comment.UserID, comment.Comment, comment.FirstName, comment.LastName, comment.Image, comment.Date)
	if err != nil {
		return err
	}

	return nil
}

func (m *SqliteDB) GetCommentsByPostID(postID int) ([]models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT comment_id, user_id, comment, first_name, last_name, image, date FROM comments WHERE post_id = ?`

	rows, err := m.DB.QueryContext(ctx, stmt, postID)
	if err != nil {
		return nil, err
	}

	var comments []models.Comment

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.CommentID, &comment.UserID, &comment.Comment, &comment.FirstName, &comment.LastName, &comment.Image, &comment.Date)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (m *SqliteDB) UpdateProfileType(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `UPDATE users SET public = NOT public WHERE user_id = ?`
	_, err := m.DB.ExecContext(ctx, stmt, userID)
	if err != nil {
		return err
	}

	return nil
}

func (m *SqliteDB) FollowUser(followerID, followingID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `INSERT INTO followers (follower_id, following_id, request_pending) VALUES (?, ?, 0)`

	_, err := m.DB.ExecContext(ctx, stmt, followerID, followingID)
	if err != nil {
		return err
	}

	return nil
}

func (m *SqliteDB) IsFollowing(userID, followingID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT EXISTS ( SELECT 1 FROM followers WHERE follower_id = $1 AND following_id = $2)`

	var isFollowing bool
	row := m.DB.QueryRowContext(ctx, stmt, userID, followingID)
	err := row.Scan(&isFollowing)
	if err != nil {
		return false, err
	}

	return isFollowing, nil
}

func (m *SqliteDB) FollowNotPublicUser(followerID, followingID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `INSERT INTO followers (follower_id, following_id, request_pending) VALUES (?, ?, 1)`

	_, err := m.DB.ExecContext(ctx, stmt, followerID, followingID)
	if err != nil {
		return err
	}

	return nil
}

func (m *SqliteDB) UnfollowUser(followerID, followingID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `DELETE FROM followers WHERE follower_id = ? AND following_id = ?`

	_, err := m.DB.ExecContext(ctx, stmt, followerID, followingID)
	if err != nil {
		return err
	}
	return nil
}

func (m *SqliteDB) IsUserPublic(userID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT public FROM users WHERE user_id = ?`

	var isPublic bool
	row := m.DB.QueryRowContext(ctx, stmt, userID)
	err := row.Scan(&isPublic)
	if err != nil {
		return false, err
	}

	return isPublic, nil
}

func (m *SqliteDB) Following(userID int) ([]models.UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		SELECT user_id, first_name, last_name FROM users
		JOIN followers ON user_id = following_id
		WHERE follower_id = $1 AND request_pending = false
	`

	rows, err := m.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}

	var following []models.UserData

	for rows.Next() {
		var user models.UserData
		err := rows.Scan(&user.UserID, &user.FirstName, &user.LastName)
		if err != nil {
			return nil, err
		}
		following = append(following, user)
	}

	return following, nil
}

func (m *SqliteDB) Followers(userID int) ([]models.UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		SELECT user_id, first_name, last_name FROM users
		JOIN followers ON user_id = follower_id
		WHERE following_id = $1 AND request_pending = false
	`

	rows, err := m.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}

	var followers []models.UserData

	for rows.Next() {
		var user models.UserData
		err := rows.Scan(&user.UserID, &user.FirstName, &user.LastName)
		if err != nil {
			return nil, err
		}
		followers = append(followers, user)
	}

	return followers, nil
}

func (m *SqliteDB) FollowRequests(userID int) ([]models.FollowRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT follower_id, request_pending FROM followers WHERE following_id = ? AND request_pending = true`

	rows, err := m.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}

	var followRequests []models.FollowRequest

	for rows.Next() {
		var followRequest models.FollowRequest
		err := rows.Scan(&followRequest.FollowingID, &followRequest.RequestPending)
		if err != nil {
			return nil, err
		}
		followRequests = append(followRequests, followRequest)
	}

	return followRequests, nil
}

func (m *SqliteDB) ApproveFollower(userID, followingID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `UPDATE followers SET request_pending = false WHERE (user_id = ? AND following_id = ?)`

	_, err := m.DB.ExecContext(ctx, stmt, userID, followingID)
	if err != nil {
		return err
	}

	return nil
}

package models

import "time"

type UserData struct {
	UserID      int    `json:"user_id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	Avatar      string `json:"avatar"`
	Nickname    string `json:"nickname"`
	AboutMe     string `json:"about_me"`
	Public      bool   `json:"public"`
}

type FollowRequest struct {
	FollowingID    int   `json:"following_id"`
	FollowerID     int   `json:"follower_id"`
	RequestPending *bool `json:"request_pending,omitempty"`
}

type Post struct {
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Privacy   string    `json:"privacy"`
	Names     []string    `json:"names"`
	Image     string    `json:"image"`
	Date      time.Time `json:"date"`
	Comments  []Comment `json:"comments"`
}

type Comment struct {
	CommentID int       `json:"comment_id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Comment   string    `json:"comment"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Image     string    `json:"image"`
	Date      time.Time `json:"date"`
}

type Session struct {
	UserID    int    `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Cookie    string `json:"cookie"`
}

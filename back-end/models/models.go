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
}

type Session struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Cookie    string `json:"cookie"`
}

type Post struct {
	PostID    int       `json:"post_id"`
	Content   string    `json:"content"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Privacy   string    `json:"privacy"`
	Image     string    `json:"image"`
	Date      time.Time `json:"date"`
}

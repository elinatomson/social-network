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
	CurrentUser bool   `json:"currentUser"`
	Online      bool   `json:"online"`
}

type FollowRequest struct {
	FollowingID    int   `json:"following_id"`
	FollowerID     int   `json:"follower_id"`
	RequestPending *bool `json:"request_pending,omitempty"`
}

type Post struct {
	PostID         int       `json:"post_id"`
	UserID         int       `json:"user_id"`
	Content        string    `json:"content"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Privacy        string    `json:"privacy"`
	SelectedUserID string    `json:"selected_user_id"`
	Image          string    `json:"image"`
	Date           time.Time `json:"date"`
	GroupID        int       `json:"group_id"`
	Comments       []Comment `json:"comments"`
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

type Message struct {
	MessageID     int
	Type          string    `json:"type"`
	Message       string    `json:"message"`
	FirstNameFrom string    `json:"first_name_from"`
	FirstNameTo   string    `json:"first_name_to"`
	Date          time.Time `json:"date"`
}

type Group struct {
	GroupID        int    `json:"group_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	UserID         int    `json:"user_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	SelectedUserID string `json:"selected_user_id"`
}

type GroupMembers struct {
	GroupID           int    `json:"group_id"`
	GroupTitle        string `json:"group_title"`
	GroupCreatorID    int    `json:"group_creator_id"`
	MemberID          int    `json:"member_id"`
	RequestPending    *bool  `json:"request_pending,omitempty"`
	InvitationPending *bool  `json:"invitation_pending,omitempty"`
}

type Event struct {
	EventID     int    `json:"event_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Time        string `json:"time"`
	GroupID     int    `json:"group_id"`
}

type EventParticipants struct {
	EventID       int    `json:"event_id"`
	ParticipantID int    `json:"participant_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Going         *bool  `json:"going,omitempty"`
}

type EventNotifications struct {
	EventID  int `json:"event_id"`
	MemberID int `json:"member_id"`
	GroupID  int `json:"group_id"`
}

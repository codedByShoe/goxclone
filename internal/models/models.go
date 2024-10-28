package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Username string `gorm:"unique"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"-"`
}

type Session struct {
	User      User   `gorm:"foreignKey:UserID" json:"user"`
	SessionID string `json:"session_id"`
	ID        uint   `gorm:"primaryKey" json:"id"`
	UserID    uint   `json:"user_id"`
}

type Post struct {
	User User
	gorm.Model
	Content  string
	Comments []Comment
	Likes    []Like
	Reposts  []Repost
	UserId   uint
}

type Like struct {
	gorm.Model
	PostId uint
}

type Repost struct {
	gorm.Model
	PostId uint
}

type Comment struct {
	gorm.Model
	Body   string
	User   User
	Post   Post
	UserId uint
	PostId uint
}

type Notification struct {
	gorm.Model
	Message string
	User    User
	UserId  uint
}

type UserRepo interface {
	CreateUser(*User) error
	GetUser(email string) (*User, error)
}

type SessionRepo interface {
	CreateSession(session *Session) (*Session, error)
	GetUserFromSession(sessionID string, userID string) (*User, error)
}

type PostRepo interface {
	CreatePost(*Post) error
	GetUsersPosts(userid uint) ([]Post, error)
	GetPost(postId uint) (Post, error)
	UpdatePost(postId uint, user_id uint, conent string) (*Post, error)
	DeletePost(postId uint, userId uint) error
}

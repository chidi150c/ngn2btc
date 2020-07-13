package data

import "time"

//DBType is the user database
type DBType map[string]*User

//type UserID string

//User is the person using the site
type User struct {
	Username      string
	Password      string
	Firstname     string
	Lastname      string
	Email         string
	ImageURL      string
	Token         string
	Url           string
	Authenticated bool
	CreatedDate   time.Time
	Expiry        int64
	Level         string
}

type UserService interface {
	AddUser(*User) error
	GetUser(username string) (*User, error)
	DeleteUser(username string) error
	ListUsers() ([]*User, error)
	UpdateUser(*User) error
}

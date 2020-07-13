package apicrypto

import (
	"context"
	"time"
)

type UserID int64

//DBType is the user modelbase
type UserDBType map[UserID]*User

//User is the person using the site
type User struct {
	ID            UserID
	Username      string
	Password      string
	Firstname     string
	Lastname      string
	Email         string
	ImageURL      string
	token         string
	userApiKey 	  string
	userApiSecret string
	Host          string
}

type UserServicer interface {
	AddUser(*User) error
	GetUser(username string) (*User, error)
	DeleteUser(username string) error
	ListUsers() ([]*User, error)
	UpdateUser(*User) error
}

package apiexch

import (
	"context"
	"time"
)

type UserID int64

//DBType is the user modelbase
type UDBType map[UserID]*User

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
	Url           string
	Authenticated bool
	CreatedTime   time.Time
	Expiry        int64
	Role          string
	Amount        string
	Bank          string
	BuyStatus     string
	SellStatus    string
	WalletAddress string
}

type UserServicer interface {
	GenerateWalletAddress(context.Context, User) WalletAddressAPI
	GetUserAddress(context.Context, User) string
	SetUserAddress(context.Context, User, string)
	AddUser(*User) error
	GetUser(username string) (*User, error)
	DeleteUser(username string) error
	ListUsers() ([]*User, error)
	UpdateUser(*User) error
}

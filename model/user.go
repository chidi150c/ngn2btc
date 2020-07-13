package model

import (
	"net/http"
	"time"
)

type UserID int64

//User is the person using the site
type User struct {
	ID            UserID
	Username      string
	Password      string
	Firstname     string
	Lastname      string
	Email         string
	ImageURL      string
	Token         string
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
	UserCryptoApiKey string
	UserCryptoApiSecret string
	CryptoHost string
}

type UserServicer interface {
	AddUser(*User, *User) (UserID, error)
	GetUser(UserID) (*User, error)
	DeleteUser(*User, UserID) error
	ListUsers() ([]*User, error)
	UpdateUser(*User, *User) error
	IsUserAlreadyInDB(string) (*User, bool, error)
}

type Sessioner interface {
	SetToken(http.ResponseWriter, *http.Request, string, UserID, string, string)
	Validate(http.HandlerFunc) http.HandlerFunc
	UserFromRequest(*http.Request) (string, *User, error)
	Userservice() UserServicer
}

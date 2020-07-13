package apiuser

type UserID int64

//UDBType is the user database
type UDBType map[UserID]*User

type User struct {
	ID            UserID
	Username      string
	Email         string
	ImageURL      string
	Url           string
	Authenticated bool
	Role          string
	WalletAddress string
	Token         string
}

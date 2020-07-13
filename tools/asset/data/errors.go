package data

// General errors.
const (
	ErrUnauthorized = Error("unauthorized")
	ErrInternal     = Error("internal error")
)

// User and Movie errors.
const (
	ErrUserNotFound     = Error("user not found")
	ErrMovieNotFound    = Error("movie not found")
	ErrExchNotFound     = Error("Exch not found")
	ErrUserExists       = Error("user already exists")
	ErrUserIDRequired   = Error("user id required")
	ErrUserNameRequired = Error("user's username required")
	ErrMovieIDRequired  = Error("movie id required")
	ErrExchIDRequired   = Error("Exch id required")
	ErrInvalidJSON      = Error("invalid json")
	ErrUserRequired     = Error("user required")
	ErrExchRequired     = Error("game required")
	ErrInvalidEntry     = Error("invalid Entry")
)

//login or Signup error
const (
	ErrUserNullPointer   = Error("User value is nill or User is Empty")
	ErrUserNotCached     = Error("Unable to save User in Cache or Session")
	ErrUserNameEmpty     = Error("Username is Empty please enter a Username")
	ErrUserPasswordEmpty = Error("Password is Empty please enter correct Password")
	ErrUsrDbUnreachable  = Error("Unable to get the UserDB into the Method")
	ErrMovDbUnreachable  = Error("Unable to get the MovieDB into the Method")
	ErrExcDbUnreachable  = Error("Unable to get the ExchDB into the Method")
)

//Session errors
const (
	ErrSessionCookieSaveError = Error("could not save cookie session please ensure cookie is enable on your browser")
	ErrIvalidRedirect         = Error("invalid redirect URL, Please try again")
	ErrSessionCookieError     = Error("could not create a cookie session please ensure cookie is enable on your browser")
)

// Error represents a User error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }

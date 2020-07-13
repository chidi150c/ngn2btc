package apiuser

import (
	"net/http"
	"time"
	"user-apiv2/model"

	"github.com/dgrijalva/jwt-go"
)

type key int

const (
	myKey key = 0
)

/*Session has the modelbase handle and current time so the services can reference them. By making the
UserService a non-pointer field we reduce the allocations required when creating
a new session.*/
type Session struct {
	db          UDBType
	userservice UserService
	session     model.Sessioner
}

func NewSession(uDB UDBType, msess model.Sessioner) *Session {
	s := &Session{
		db:         uDB,
		session:    msess,
		cachedUser: &User{},
	}
	s.userservice.session = s
	s.userservice.userservice = s.session.Userservice()
	return s
}

func (s *Session) UserFromRequest(r *http.Request) (string, *User, error) {
	red, modelusr, err := s.session.UserFromRequest(r)
	// modelusr, ok := usr.(*model.User)
	// if !ok {
	// 	return "", nil, errors.New("Wrong user type")
	// }
	otherusr := s.userservice.ModelToOtherUser(modelusr)
	s.cachedUser = otherusr
	//s.db[otherusr.ID] = s.cachedUser
	return red, otherusr, err
}

// deletes the cookie
func (s *Session) logout(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("Auth")
	if err != nil {
		return
	}
	deleteCookie := http.Cookie{Name: "Auth", Value: "none", Expires: time.Now()}
	http.SetCookie(w, &deleteCookie)
	return
}

// JWT schema of the model it will store
type claims struct {
	username    string `json:"username"`
	redirectURL string `json:"redirecturl"`
	level       string `json:"level"`
	jwt.StandardClaims
}

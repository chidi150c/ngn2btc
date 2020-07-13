package memory

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"user-apiv2/model"

	"github.com/dgrijalva/jwt-go"
)

//DBType is the user database
type MDBType map[model.UserID]*model.User

type key int

const myKey key = 0

/*Session has the database handle and current time so the services can reference them. By making the
UserService a non-pointer field we reduce the allocations required when creating
a new session.*/
type Session struct {
	db          MDBType
	UserService UserService
}

func NewSession(uDB MDBType) *Session {
	s := &Session{
		db: uDB,
	}
	s.UserService.session = s
	return s
}

// JWT schema of the data it will store
type claims struct {
	Username    string       `json:"username"`
	ID          model.UserID `json:"id"`
	RedirectURL string       `json:"redirecturl"`
	Level       string       `json:"level"`
	Token       string       `json:"-"`
	jwt.StandardClaims
}

// create a JWT and put in the clients cookie
func (s *Session) SetToken(w http.ResponseWriter, r *http.Request, username string, userid model.UserID, redirectURL, level string) {
	log.Println("*******************In SetToken1*******************")
	expireCookie := time.Now().Add(time.Minute * 20)
	expireToken := time.Now().Add(time.Minute * 20).Unix()
	userClaims := &claims{
		username,
		userid,
		redirectURL,
		level,
		"tokenkey",
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "localhost:8080",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	signedToken, _ := token.SignedString([]byte("sEcrEtPassWord!234"))
	cookie := http.Cookie{Name: "Auth", Value: signedToken, Expires: expireCookie, HttpOnly: true}
	http.SetCookie(w, &cookie)
	time.Sleep(300 * time.Millisecond)
	// redirect
	http.Redirect(w, r, userClaims.RedirectURL, http.StatusSeeOther)
	return
}

// middleware to protect private pages
func (s *Session) Validate(handler http.HandlerFunc) http.HandlerFunc {
	log.Println("*******************In Validate1*******************")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Auth")
		if err != nil {
			ctx := context.WithValue(r.Context(), myKey, &claims{"anonymous", 0, "", "", "", jwt.StandardClaims{}})
			handler(w, r.WithContext(ctx))
			// http.NotFound(w, r)
			return
		}
		token, err := jwt.ParseWithClaims(cookie.Value, &claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return []byte("sEcrEtPassWord!234"), nil
		})
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if userClaims, ok := token.Claims.(*claims); ok && token.Valid {
			userClaims.Token = token.Raw
			ctx := context.WithValue(r.Context(), myKey, *userClaims)
			handler(w, r.WithContext(ctx))
		} else {
			http.NotFound(w, r)
			return
		}
	})
}

func (s *Session) UserFromRequest(r *http.Request) (string, *model.User, error) {
	log.Println("*******************In UserFromRequest1*******************")
	userClaims, ok := r.Context().Value(myKey).(claims)
	if !ok || userClaims.Username == "anomnymous" {
		return "", nil, model.ErrUnauthorized
	}
	user, err := s.UserService.GetUser(userClaims.ID)
	if err != nil {
		return "", nil, model.ErrUnauthorized
	}

	user.Token = userClaims.Token
	return userClaims.RedirectURL, user, nil
}

// Authenticate returns the current authenticate user.
func (s *Session) Userservice() model.UserServicer {
	return &s.UserService
}


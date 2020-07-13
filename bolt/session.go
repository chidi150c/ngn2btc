package bolt

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"user-apiv2/model"

	bolt "github.com/coreos/bbolt"
	"github.com/dgrijalva/jwt-go"
)

type key int

const myKey key = 0

/*Session has the database handle and current time so the services can reference them. By making the
UserService a non-pointer field we reduce the allocations required when creating
a new session.*/
type Session struct {
	db *bolt.DB
	//db          UDBType
	UserService UserService
}

//uDB UDBType
func NewSession(udb *bolt.DB) *Session {
	s := &Session{
		db:         udb,
	}
	s.UserService.session = s

	// Start a writable transaction.
	tx, err := s.db.Begin(true)
	if err != nil {
		log.Fatalln("Error1 in bolt.Newsessio:38: ", err)
	}
	defer tx.Rollback()

	// Use the transaction...
	_, err = tx.CreateBucketIfNotExists([]byte("Users"))
	if err != nil {
		log.Fatalln("Error2 in bolt.Newsessio:45:", err)
	}

	// Commit the transaction and check for error.
	if err = tx.Commit(); err != nil {
		log.Fatalln("Error3 in bolt.Newsessio:50:", err)
	}
	return s
}

// Open opens and initializes the BoltDB database.
func (s *Session) Open() error {

	// Initialize top-level buckets.
	tx, err := s.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.CreateBucketIfNotExists([]byte("Users")); err != nil {
		return err
	}

	return tx.Commit()
}

// Close closes then underlying BoltDB database.
func (s *Session) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
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

var _ model.Sessioner = &Session{}

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
	log.Println("*********Started Session.UserFromRequestr********")
	defer log.Println("*********End Session.UserFromRequest********")

	userClaims, ok := r.Context().Value(myKey).(claims)
	if !ok || userClaims.Username == "anomnymous" {
		return "", nil, model.ErrUnauthorized
	}
	user, err := s.Userservice().GetUser(userClaims.ID)
	if err != nil {
		return "", nil, model.ErrUnauthorized
	}
	log.Println("user.Token = userClaims.Token", userClaims.Token)
	user.Token = userClaims.Token
	return userClaims.RedirectURL, user, nil
}

// Authenticate returns the current authenticate user.
func (s *Session) Userservice() model.UserServicer {
	return &s.UserService
}

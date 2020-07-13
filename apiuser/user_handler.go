package apiuser

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"user-apiv2/model"
	"user-apiv2/tools"

	"github.com/pkg/errors"
	"github.com/pressly/chi"
)

// UserHandler represents an HTTP API handler for users.
type UserHandler struct {
	mux *chi.Mux
	//redirectURL string
	session *Session
	//user        *User
	Logger *log.Logger
}

// NewUserHandler returns a new instance of UserHandler.
//mv UserService
func NewUserHandler(dbUser UDBType, msess model.Sessioner) *UserHandler {
	h := &UserHandler{
		mux:     chi.NewRouter(),
		Logger:  log.New(os.Stderr, "", log.LstdFlags),
		session: NewSession(dbUser, msess),
	}
	h.session.cachedUser = &User{}
	h.session.session.ClearCacheUser()
	// h.mux.Get("/", Validate(h.indexHandler))
	//h.mux.Get("/signup", h.signupHandler)
	h.mux.Post("/signup", h.signupPostHandler)
	h.mux.Get("/login", h.session.session.Validate(h.loginHandler))
	h.mux.Post("/login", h.loginPostHandler)
	h.mux.Post("/logout", h.logoutHandler)
	//h.mux.Get("/users/list", s.Validate(h.handleListUsers))
	//h.mux.Get("/user/delete/:username", Validate(h.handleDeleteUser))
	//h.mux.Post("/user", h.GetUsersHandler)
	//h.mux.Get("/users/get/:username", s.Validate(h.handleGetUser))
	//h.mux.Post("/users/update/:username", s.Validate(h.handleUpdateUser))
	h.mux.Get("/", h.session.session.Validate(h.indexHandler))
	return h
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("*********Started UserHandler.ServeHTTP********")
	defer log.Println("*********End UserHandler.ServeHTTP********")

	h.mux.ServeHTTP(w, r)
}

func (h *UserHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("*********Started UserHandler.indexHandler********")
	defer log.Println("*********End UserHandler.indexHandler********")

	var webuser *User
	_, user, err := h.session.UserFromRequest(r)
	if err != nil {
		webuser = nil
	} else if user.ImageURL == "" {
		user.ImageURL = "/tools/asset/images/user.png"
		webuser = user
	} else {
		webuser = user
	}
	if err := tools.IndexTmpl.Execute(w, r, nil, webuser, nil); err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	return
}

func (h *UserHandler) signupPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("*********Started UserHandler.signupPostHandler********")
	defer log.Println("*********End UserHandler.signupPostHandler********")

	redirectURL, err := validateRedirectURL(r.FormValue("redirect"))
	if err != nil {
		tools.Error(w, errors.Wrapf(err, "invalid redirect URL: %v", err), http.StatusInternalServerError, h.Logger)
		return
	}
	msg := &Message{
		Email:          r.FormValue("email"),
		Firstname:      r.FormValue("firstname"),
		Lastname:       r.FormValue("lastname"),
		Password:       r.FormValue("password"),
		PasswordVerify: r.FormValue("passwordVerify"),
		Username:       r.FormValue("username"),
		RedirectURL:    redirectURL,
	}
	msg.Errors = make(map[string]interface{})
	h.session.session.ClearCacheUser()

	// username taken or password or email incorrect ?
	if msg.Username == "" {
		msg.Errors["Username"] = model.ErrUserNameEmpty

		//get user with username and check if it exist
	} else if usr, oK, _ := h.session.session.Userservice().IsUserAlreadyInDB(msg.Username); oK == true && usr.Username != "" {
		msg.Errors["Username"] = fmt.Sprintf("Username \"%s\" is already taken, please try another", usr.Username)
	}
	if msg.Password == "" {
		msg.Errors["Password"] = model.ErrUserPasswordEmpty
	}
	_ = msg.Validate()
	if len(msg.Errors) != 0 {
		if err := tools.IndexTmpl.Execute(w, r, nil, nil, msg); err != nil {
			tools.Error(w, err, http.StatusInternalServerError, h.Logger)
			return
		}
		return
	}
	//*************************************************
	// get new valid values form request values
	usr := &model.User{
		Email:     r.FormValue("email"),
		Username:  r.FormValue("username"),
		Password:  r.FormValue("password"),
		Firstname: r.FormValue("firstname"),
		Lastname:  r.FormValue("lastname"),
	}
	userid, err := h.session.session.Userservice().AddUser(usr)
	if err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	//Set token for the user in cookie
	h.session.session.SetToken(w, r, usr.Username, userid, msg.RedirectURL, usr.Role)
	return
}

func (h *UserHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("*********Started UserHandler.loginHandler********")
	defer log.Println("*********End UserHandler.loginHandler********")

	redirectURL, err := validateRedirectURL(chi.URLParam(r, "redirect"))
	if err != nil {
		tools.Error(w, errors.Wrapf(err, "invalid redirect URL: %v", err), http.StatusInternalServerError, h.Logger)
		return
	}
	_, user, err := h.session.UserFromRequest(r)
	if err != nil {
		msg := &Message{
			RedirectURL: redirectURL,
			Errors:      map[string]interface{}{"Anonymous": "You need to login to continue access to this site"},
		}
		if err := tools.IndexTmpl.Execute(w, r, nil, nil, msg); err != nil {
			tools.Error(w, err, http.StatusInternalServerError, h.Logger)
			return
		}
		return
	}
	if err := tools.IndexTmpl.Execute(w, r, nil, user, nil); err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	return
}

func (h *UserHandler) loginPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("*********Started UserHandler.loginPostHandler********")
	defer log.Println("*********End UserHandler.loginPostHandler********")

	redirectURL, err := validateRedirectURL(r.FormValue("redirect"))
	if err != nil {
		tools.Error(w, errors.Wrapf(err, "invalid redirect URL: %v", err), http.StatusInternalServerError, h.Logger)
		return
	}
	msg := &Message{
		Password:    string(r.FormValue("password")),
		Username:    string(r.FormValue("username")),
		RedirectURL: redirectURL,
	}
	msg.Errors = make(map[string]interface{})
	// username taken or password or email incorrect ?
	var user = &model.User{}
	var oK bool
	if msg.Username == "" {
		msg.Errors["Username"] = model.ErrUserNameEmpty
		//get user with username and check if it exist
	} else if user, oK, _ = h.session.session.Userservice().IsUserAlreadyInDB(msg.Username); oK == false && user == nil {
		msg.Errors["Username"] = fmt.Sprintf("Unknown Username or Password, Try and signup")
	}
	if msg.Password == "" {
		msg.Errors["Password"] = model.ErrUserPasswordEmpty
	} else if user != nil {
		if user.Password != msg.Password {
			msg.Errors["Password"] = fmt.Sprintf("Unknown Username or Password, Try and signup")
		}
	}
	if len(msg.Errors) != 0 {
		if err := tools.IndexTmpl.Execute(w, r, nil, nil, msg); err != nil {
			tools.Error(w, err, http.StatusInternalServerError, h.Logger)
			return
		}
		return
	}
	//Set token for the user in cookie
	h.session.session.SetToken(w, r, user.Username, user.ID, redirectURL, user.Role)
	return
}

// deletes the cookie
func (h *UserHandler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("*********Started UserHandler.logoutHandler********")
	defer log.Println("*********End UserHandler.logoutHandler********")

	h.session.logout(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

// ValidateRedirectURL checks that the URL provided is valid.
// If the URL is missing, redirect the user to the application's root.
// The URL must not be absolute (i.e., the URL must refer to a path within this
// application).
func validateRedirectURL(path string) (string, error) {
	if path == "" {
		return "/", nil
	}

	// Ensure redirect URL is valid and not pointing to a different server.
	parsedURL, err := url.Parse(path)
	if err != nil {
		return "/", err
	}
	if parsedURL.IsAbs() {
		return "/", errors.New("URL invalid: URL must be absolute")
	}
	if strings.Contains(path, "/signup?redirect=") {
		return path[17:], nil
	}
	return path, nil
}

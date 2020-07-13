package memory

import (
	"log"
	"time"
	"user-apiv2/model"
)

var nextID int64

func init() {
	nextID = 1
}

type UserService struct {
	session *Session
}

var _ model.UserServicer = &UserService{}

// func (u *UserService) OtherUser(*model.User) interface{} {
// 	return nil
// }
// func (u *UserService) ModelToOtherUser(usr *model.User) interface{} {
// 	return nil
// }

func (u *UserService) AddUser(caller, usr *model.User) (model.UserID, error) {
	log.Println("*********Started Memory.UserService.AddUser********")
	defer log.Println("*********End Memory.UserService.AddUser********")

	if usr.Username == "" {
		return 0, model.ErrUserNameEmpty
	}
	// Retrieve current session user.
	if caller.Username != "Anonymous" && caller.Token != usr.Token{
		return 0, model.ErrUnauthorized
	}
	usr.CreatedTime = time.Now()
	id := model.UserID(nextID)
	usr.ID = id
	u.session.db[id] = usr
	nextID++
	return id, nil
}
func (u *UserService) IsUserAlreadyInDB(username string) (*model.User, bool, error) {
	if username == "" {
		return nil, false, model.ErrUserNameEmpty
	}
	for _, v := range u.session.db {
		if v.Username == username {
			return v, true, model.ErrUserExists
		}
	}
	return nil, false, model.ErrUserNotFound
}

func (u *UserService) GetUser(id model.UserID) (*model.User, error) {

	if id == 0 {
		return nil, model.ErrUserNameEmpty
	}
	usr := u.session.db[id]
	if usr == nil {
		return nil, model.ErrUserNotFound
	}
	return usr, nil
}

func (u *UserService) DeleteUser(caller *model.User, id model.UserID) error {
	// Retrieve current session user.
	if caller == nil {
		return model.ErrUnauthorized
	}
	userInDB, ok := u.session.db[id]
	if !ok {
		return model.ErrUserNotFound
	}
	if userInDB.ID != caller.ID {
		return model.ErrUnauthorized
	}
	delete(u.session.db, id)
	return nil
}

func (u *UserService) ListUsers() ([]*model.User, error) {

	var users []*model.User
	for _, b := range u.session.db {
		users = append(users, b)
	}
	return users, nil
}

func (u *UserService) UpdateUser(caller *model.User, usr *model.User) error {
	// Retrieve current session user.
	if caller == nil {
		return model.ErrUnauthorized
	}
	// Only allow owner to update user.
	userInDB, ok := u.session.db[usr.ID]
	if !ok {
		return model.ErrUserNotFound
	}
	if userInDB.ID != caller.ID && usr.Token != userInDB.Token {
		return model.ErrUnauthorized
	}

	// Update user.
	u.session.db[usr.ID] = usr
	return nil
}

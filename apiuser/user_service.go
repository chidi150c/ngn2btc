package apiuser

import (
	"user-apiv2/model"
)

type UserService struct {
	session     *Session
	userservice model.UserServicer
}

//var _ model.UserServicer = &UserService{}

func (u *UserService) ModelToOtherUser(user *model.User) *User {
	if user == nil {
		return nil
	}
	return &User{
		ID:            UserID(user.ID),
		Username:      user.Username,
		Email:         user.Email,
		ImageURL:      user.ImageURL,
		Url:           user.Url,
		Authenticated: user.Authenticated,
		Role:          user.Role,
		Token:         user.Token,
	}
}
func (u *UserService) OtherToModelUser(user *User, musr *model.User) *model.User {
	if user == nil {
		return nil
	}
	musr.ID = model.UserID(user.ID)
	musr.Username = user.Username
	musr.Email = user.Email
	musr.ImageURL = user.ImageURL
	musr.Url = user.Url
	musr.Authenticated = user.Authenticated
	musr.Role = user.Role
	return musr
}

func (u *UserService) AddUser(user *User) error {
	var muser = &model.User{}
	muser = u.OtherToModelUser(user, muser)
	err := u.userservice.UpdateUser(muser)
	if err != nil {
		return err
	}
	u.session.db[user.ID] = user
	u.session.caller = user
	return nil
}

func (u *UserService) UpdateUser(caller *User, user *User) error {
	if caller != nil {
		return model.ErrUnauthorized
	}
	// Only allow owner to update user.
	userInDB, ok := u.session.db[user.ID]
	if !ok {
		return model.ErrUserNotFound
	}
	if userInDB.ID != caller.ID && user.Token != userInDB.Token {
		return model.ErrUnauthorized
	}
	var musr = &model.User{}
	musr = u.OtherToModelUser(user, musr)
	err := u.userservice.UpdateUser(musr)
	if err != nil {
		return err
	}
	u.session.db[user.ID] = user
	return nil
}

func (u *UserService) DeleteUser(id UserID) error {
	caller, err := u.session.UserFromSession()
	if err != nil {
		return err
	}
	userInDB, ok := u.session.db[id]
	if !ok {
		return model.ErrUserNotFound
	}
	if userInDB.ID != caller.ID {
		return model.ErrUnauthorized
	}
	u.session.caller = &User{}
	delete(u.session.db, id)
	return nil
}

func (u *UserService) GetUser(id UserID) (*User, error) {
	if id == 0 {
		return nil, model.ErrUserNameEmpty
	}
	if usr := u.session.db[id]; usr != nil {
		return usr, nil
	}
	if usr := u.session.caller; usr != nil {
		return usr, nil
	}
	musr, err := u.userservice.GetUser(model.UserID(id))
	if err != nil {
		return nil, err
	}
	usr := u.ModelToOtherUser(musr)
	u.session.db[id] = usr
	u.session.caller = usr
	return usr, nil
}

func (u *UserService) ListUsers() ([]*User, error) {
	usr := []*User{}
	if len(u.session.db) > 0 {
		for _, v := range u.session.db {
			usr = append(usr, v)
		}
		return usr, nil
	}
	musers, err := u.userservice.ListUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range musers {
		usr = append(usr, u.ModelToOtherUser(user))
	}
	return usr, nil
}

package apicrypto

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"user-apiv2/model"
)

type UserService struct {
	session     *Session
	Userservice model.UserServicer
}

func (u *UserService) ModelToOtherUser(user *model.User) *User {
	if user == nil {
		return nil
	}
	transuser := &User{
		ID:            UserID(user.ID),
		Username:      user.Username,
		token:         user.Token,
		UserApiKey: user.UserCryptoApiKey,
		UserApiSecret: user.UserCryptoApiSecret,
		Host: user.CryptoHost,
	}
	return transuser
}


package apiexch

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
	exuser := &User{
		ID:            UserID(user.ID),
		Username:      user.Username,
		Password:      user.Password,
		Firstname:     user.Firstname,
		Lastname:      user.Lastname,
		Email:         user.Email,
		ImageURL:      user.ImageURL,
		token:         user.Token,
		Url:           user.Url,
		Authenticated: user.Authenticated,
		CreatedTime:   user.CreatedTime,
		Expiry:        user.Expiry,
		Role:          user.Role,
		Amount:        user.Amount,
		Bank:          user.Bank,
		BuyStatus:     user.BuyStatus,
		SellStatus:    user.SellStatus,
		WalletAddress: user.WalletAddress,
	}
	return exuser
}

func (u *UserService) GenerateWalletAddress(ctx context.Context, usr *User) WalletAddressAPI {
	fmt.Println()
	log.Println("**In GenerateWalletAddress Just entered************")
	log.Println("**In GenerateWalletAddress Just entered****usr.Password = ", usr.Password, " email = ", usr.Email)
	var walletAddress WalletAddressAPI
	url := "http://localhost:3000/api/v2/create/?password=" + usr.Password + "&api_code=fecf93e1-56a0-4eef-90de-a388f84c7e05&email=" + usr.Email + "&hd=true"
	//To get walletAddress
	walletAddressResp, err := http.Get(url)
	if err != nil {
		log.Printf("could not connecct to CurrencyLayer for NGN Rate to USD: Error is %v", err)
	}
	walletAddresser := json.NewDecoder(walletAddressResp.Body)
	err = walletAddresser.Decode(&walletAddress)
	if err != nil {
		log.Fatalf("could not json decode into walletAddress from walletAddress.Body Error is %v", err)
	}
	log.Println("**In GenerateWalletAddress Just entered********walletAddress = ", walletAddress)
	return walletAddress
}
func (u *UserService) GetUserAddress(ctx context.Context, usr *User) string {
	if usr.WalletAddress != "" {
		return usr.WalletAddress
	}
	a := u.GenerateWalletAddress(ctx, usr)
	usr.WalletAddress = a.Address
	return usr.WalletAddress
}

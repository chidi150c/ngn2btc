package internal

import (
	"time"
	"user-apiv2/model"

	proto "github.com/golang/protobuf/proto"
)

//go:generate protoc --gogo_out=. internal.proto

// MarshalUser encodes a User to binary format.
func MarshalUser(d *model.User) ([]byte, error) {
	return proto.Marshal(&User{
		ID:            int64(d.ID),
		Token:         d.Token,
		Username:      d.Username,
		Password:      d.Password,
		Firstname:     d.Firstname,
		Lastname:      d.Lastname,
		Email:         d.Email,
		ImageURL:      d.ImageURL,
		Url:           d.Url,
		Authenticated: d.Authenticated,
		CreatedTime:   d.CreatedTime.UnixNano(),
		Expiry:        d.Expiry,
		Amount:        d.Amount,
		Bank:          d.Bank,
		BuyStatus:     d.BuyStatus,
		SellStatus:    d.SellStatus,
		WalletAddress: d.WalletAddress,
		Role:          d.Role,
	})
}

// UnmarshalUser decodes a User from a binary data.
func UnmarshalUser(data []byte, d *model.User) error {
	var pb User
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	d.ID = model.UserID(pb.ID)
	d.Token = pb.Token
	d.Username = pb.Username
	d.Password = pb.Password
	d.Firstname = pb.Firstname
	d.Lastname = pb.Lastname
	d.Email = pb.Email
	d.ImageURL = pb.ImageURL
	d.Token = pb.Token
	d.Url = pb.Url
	d.Authenticated = pb.Authenticated
	d.CreatedTime = time.Unix(0, pb.CreatedTime).UTC()
	d.Expiry = pb.Expiry
	d.Role = pb.Role
	d.Amount = pb.Amount
	d.Bank = pb.Bank
	d.BuyStatus = pb.BuyStatus
	d.SellStatus = pb.SellStatus
	d.WalletAddress = pb.WalletAddress

	return nil
}

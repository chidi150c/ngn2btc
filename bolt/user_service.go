package bolt

import (
	"encoding/binary"
	"log"
	"user-apiv2/bolt/internal"
	"user-apiv2/model"

	bolt "github.com/coreos/bbolt"
)

var nextID int64

func init() {
	nextID = 1
}

type UserService struct {
	session *Session
}

var _ model.UserServicer = &UserService{}

// func (u *UserService) OtherUser(*User) interface{} {
// 	return nil
// }
// func (u *UserService) ModelToOtherUser(usr *User) interface{} {
// 	return nil
// }

func (u *UserService) AddUser(usr *model.User) (model.UserID, error) {
	log.Println("*********Started Bolt.UserService.AddUser********")
	defer log.Println("*********End Bolt.UserService.AddUser********")

	if usr.Username == "" {
		return 0, model.ErrUserNameEmpty
	}
	var musr = &model.User{}
	err := u.session.db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte("Users"))
		if err := b.ForEach(func(k, v []byte) error {
			if err := internal.UnmarshalUser(v, musr); err == nil {
				if musr.Username == usr.Username {
					return model.ErrUserExists
				}
			}
			return nil
		}); err != nil {
			return err
		}
		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := b.NextSequence()
		usr.ID = model.UserID(id)

		// Marshal user data into bytes.
		buf, err := internal.MarshalUser(usr)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put(itob(usr.ID), buf)
	})
	if err != nil {
		return 0, err
	}
	return usr.ID, nil
}

// itob returns an 8-byte big endian representation of v.
func itob(v model.UserID) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func (u *UserService) IsUserAlreadyInDB(username string) (*model.User, bool, error) {
	if username == "" {
		return nil, false, model.ErrUserNameEmpty
	}
	var musr = &model.User{}
	err := u.session.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("Users")).ForEach(func(k, v []byte) error {
			if err := internal.UnmarshalUser(v, musr); err == nil {
				if musr.Username == username {
					return model.ErrUserExists
				}
			}
			return nil
		})
	})
	if err != nil {
		return musr, true, err
	}
	return nil, false, model.ErrUserNotFound
}

func (u *UserService) GetUser(id model.UserID) (*model.User, error) {
	if id == 0 {
		return nil, model.ErrUserNameEmpty
	}
	var musr = &model.User{}
	err := u.session.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket([]byte("Users")).Get(itob(id))
		return internal.UnmarshalUser(v, musr)
	})
	if err != nil {
		return nil, err
	}
	if musr.Username == "" {
		return nil, model.ErrUserNotFound
	}
	return musr, nil
}

func (u *UserService) DeleteUser(id model.UserID) error {
	// Retrieve current session user.
	cachedUser, err := u.session.UserFromSession()
	if err != nil {
		return err
	}
	var musr = &model.User{}
	err = u.session.db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte("Users"))
		v := b.Get(itob(id))
		if err := internal.UnmarshalUser(v, musr); err != nil {
			return err
		}
		if musr.ID != cachedUser.ID && id != musr.ID {
			return model.ErrUnauthorized
		}
		return b.Delete(itob(id))
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) ListUsers() ([]*model.User, error) {

	var users []*model.User
	var usr = &model.User{}
	tx, err := u.session.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	b := tx.Bucket([]byte("Users"))
	if err := b.ForEach(func(_, v []byte) error {
		if err := internal.UnmarshalUser(v, usr); err != nil {
			return err
		}
		users = append(users, usr)
		return nil
	}); err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserService) UpdateUser(usr *model.User) error {
	// Retrieve current session user.
	cachedUser, err := u.session.UserFromSession()
	if err != nil {
		return err
	}
	// Only allow owner to update user.
	tx, err := u.session.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var usrInDB *model.User
	b := tx.Bucket([]byte("Users"))
	if v := b.Get(itob(usr.ID)); v == nil {
		return model.ErrUserNotFound
	} else if err := internal.UnmarshalUser(v, usrInDB); err != nil {
		return err
	}

	if usrInDB.ID != cachedUser.ID && usrInDB.Token != usr.Token {
		return model.ErrUnauthorized
	}
	v, err := internal.MarshalUser(usr)
	if err != nil {
		return err
	}
	if err := b.Put(itob(usr.ID), v); err != nil {
		return err
	}
	// Update user.
	u.session.CachedUser = usr
	return nil
}

package user

import (
	"fmt"
)

var ErrInvalidCredentials = fmt.Errorf("invalid credentials")

type User struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func Create(email string, password string) (*User, error) {
	u := new(User)
	u.Email = email
	pswHash, err := HashAndSaltPassword(password)
	if err != nil {
		return nil, err
	}
	u.PasswordHash = pswHash
	err = store().Create(u)
	return u, err
}

func VerifyPassword(email string, password string) (*User, error) {
	u, err := store().ByEmail(email)
	if err != nil {
		return nil, err
	}

	if !CheckPassword(u.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}
	return u, nil
}

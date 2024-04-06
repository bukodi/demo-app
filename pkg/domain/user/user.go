package user

import (
	"context"
	"fmt"
)

var ErrInvalidCredentials = fmt.Errorf("invalid credentials")

const (
	RoleEndUser string = "end_user"
	RoleAdmin   string = "admin"
)

type User struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}

func Create(ctx context.Context, email string, password string) (*User, error) {
	u := new(User)
	u.Email = email
	pswHash, err := HashAndSaltPassword(password)
	if err != nil {
		return nil, err
	}
	u.PasswordHash = pswHash
	err = store().Create(ctx, u)
	return u, err
}

func VerifyPassword(ctx context.Context, email string, password string) (*User, error) {
	u, err := store().ByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !CheckPassword(u.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}
	return u, nil
}

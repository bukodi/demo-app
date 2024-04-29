package user

import (
	"context"
	"fmt"
	"github.com/bukodi/demo-app/pkg/authn"
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

func List(ctx context.Context) ([]*User, error) {
	authnUser := authn.GetUser(ctx)
	if authnUser == nil || !authnUser.HasRole(RoleAdmin) {
		return nil, fmt.Errorf("unauthorized")
	}
	return store().List(ctx)
}

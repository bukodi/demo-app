package user

import (
	"context"
	"github.com/bukodi/demo-app/pkg/authn"
	"log/slog"
)

const builtinIDPName = "built-in"

func init() {
	authn.RegisterIdentityProvider("user", &authnIDP{})
}

type authnIDP struct {
}

func (a authnIDP) Name() string {
	return builtinIDPName
}

func (a authnIDP) FindById(ctx context.Context, userid string) authn.User {
	domainUser, err := ByEmail(ctx, userid)
	if err != nil {
		slog.Error("failed to find user by id", "userid", userid, "error", err)
		return nil
	}
	if domainUser == nil {
		return nil
	}

	return &authnUser{domainUser: domainUser}
}

type authnUser struct {
	domainUser *User
}

func (a authnUser) Id() string {
	return a.domainUser.Email
}

func (a authnUser) IDPName() string {
	return builtinIDPName
}

func (a authnUser) HasRole(role string) bool {
	return role == a.domainUser.Role
}

var _ authn.User = (*authnUser)(nil)

func (a authnIDP) FindAndAuthorize(ctx context.Context, userid string, password string) (authn.User, error) {
	u, err := VerifyPassword(ctx, userid, password)
	if err != nil {
		return nil, err
	}
	return &authnUser{domainUser: u}, nil
}

var _ authn.IdentityProvider = (*authnIDP)(nil)

package user

import (
	"context"
	"github.com/bukodi/demo-app/pkg/authn"
)

func init() {
	authn.RegisterIdentityProvider("user", &authnIDP{})
}

type authnIDP struct {
}

type authnUser struct {
	domainUser *User
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

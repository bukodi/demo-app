package user

import (
	"context"
	"github.com/bukodi/demo-app/pkg/authn"
)

type fakeUser struct {
	roles []string
}

func (u *fakeUser) HasRole(role string) bool {
	for _, r := range u.roles {
		if r == role {
			return true
		}
	}
	return false
}

func adminUserContext() context.Context {
	return authn.WithUser(context.Background(), &fakeUser{
		roles: []string{"admin"},
	})
}

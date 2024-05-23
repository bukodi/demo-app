package user

import (
	"context"
	"github.com/bukodi/demo-app/pkg/authn"
	"slices"
)

type fakeUser struct {
	roles []string
}

func (u *fakeUser) HasRole(role string) bool {
	return slices.Contains(u.roles, role)
}

func adminUserContext() context.Context {
	return authn.WithUser(context.Background(), &fakeUser{
		roles: []string{"admin"},
	})
}

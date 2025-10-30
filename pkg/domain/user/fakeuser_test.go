package user

import (
	"context"
	"slices"

	"github.com/bukodi/demo-app/pkg/authn"
)

type fakeUser struct {
	roles []string
}

func (u *fakeUser) Id() string {
	return "fakeUser"
}

func (u *fakeUser) IDPName() string {
	return "fakeIDP"
}

func (u *fakeUser) HasRole(role string) bool {
	return slices.Contains(u.roles, role)
}

func adminUserContext() context.Context {
	return authn.WithUser(context.Background(), &fakeUser{
		roles: []string{"admin"},
	})
}

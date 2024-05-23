package authn

import (
	"context"
	"fmt"
	"slices"
)

type testUser struct {
	userid   string
	password string
	roles    []string
}

func (tUsr *testUser) Id() string {
	return tUsr.userid
}

const testIDPName = "test_idp"

func (tUsr *testUser) IDPName() string {
	return testIDPName
}

func (tUsr *testUser) HasRole(role string) bool {
	return slices.Contains(tUsr.roles, role)
}

var _ User = &testUser{}

type testIDP struct {
	users []testUser
}

func (tIDP *testIDP) Name() string {
	return testIDPName
}

func (tIDP *testIDP) FindById(ctx context.Context, userid string) User {
	idx := slices.IndexFunc(tIDP.users, func(tu testUser) bool {
		return tu.userid == userid
	})
	if idx == -1 {
		return nil
	} else {
		return &tIDP.users[idx]
	}
}

func (tIDP *testIDP) FindAndAuthorize(ctx context.Context, userid string, password string) (User, error) {
	for _, user := range tIDP.users {
		if user.userid == userid {
			if user.password != password {
				return nil, fmt.Errorf("password not matches")
			}
			return &user, nil
		}
	}
	return nil, nil
}

var _ IdentityProvider = &testIDP{}

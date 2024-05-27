package user

import (
	"github.com/bukodi/demo-app/pkg/authz"
)

type authzType struct{} // dummy type to implement authz.Type
func (authzType) TypeName() string {
	return "user"
}

var AuthzType authz.Type = authzType{}

func (u *User) HasRole(role string) bool {
	return u.Role == role
}

func (u *User) String() string {
	return u.Email
}

func (u *User) Id() string {
	return u.Email
}

func (u *User) Type() authz.Type {
	return AuthzType
}

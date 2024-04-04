package user

import "github.com/bukodi/demo-app/pkg/authz"

type authzType struct{} // dummy type to implement authz.Type
func (authzType) TypeName() string {
	return "user"
}

var AuthzType authz.Type = authzType{}

var _ authz.User = (*User)(nil)

func (u *User) String() string {
	return u.Email
}

func (u *User) Id() string {
	return u.Email
}

func (u *User) Type() authz.Type {
	return AuthzType
}

package authz

import (
	"context"
	"github.com/bukodi/demo-app/pkg/authn"
)

type Relation interface {
	String() string
}
type Object interface {
	String() string
	Id() string
	Type() Type
}

type Type interface {
	TypeName() string
}

func Check(user authn.User, relation Relation, object Object) bool {
	return true
}

func CheckCtx(ctx context.Context, relation Relation, object Object) bool {
	user := authn.GetUser(ctx)
	return Check(user, relation, object)
}

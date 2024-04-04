package authz

type User interface {
	String() string
	Id() string
	Type() Type
}
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

func Check(user User, relation Relation, object Object) bool {
	return true
}

package authn

type User interface {
	HasRole(role string) bool
	Id() string
	IDPName() string
}

package user

import (
	"context"

	"github.com/bukodi/demo-app/pkg/util"
)

type Store interface {
	// Create inserts a new user into the store
	Create(ctx context.Context, user *User) error
	// List returns all users from the store
	List(ctx context.Context) ([]*User, error)
	// ByEmail returns a user by email
	ByEmail(ctx context.Context, email string) (*User, error)

	// Delete removes a user from the store. Returns true if the user was found and deleted
	Delete(ctx context.Context, email string) (bool, error)
}

type errWrapperStore struct {
	initErr error
}

var _ Store = (*errWrapperStore)(nil)

func (n errWrapperStore) Create(ctx context.Context, user *User) error { return n.initErr }
func (n errWrapperStore) List(ctx context.Context) ([]*User, error)    { return nil, n.initErr }
func (n errWrapperStore) ByEmail(ctx context.Context, email string) (*User, error) {
	return nil, n.initErr
}
func (n errWrapperStore) Delete(ctx context.Context, email string) (bool, error) {
	return false, n.initErr
}

var spi = util.NewGenericSpi[Store]()

func RegisterStoreProvider(providerName string, initFn func(context.Context) (Store, error)) {
	spi.RegisterProvider(providerName, initFn)
}

func ResetStore() {
	spi.Reset()
}

func store() Store {
	actual, initErr := spi.Actual()
	if initErr == nil {
		return actual
	}

	return &errWrapperStore{initErr}
}

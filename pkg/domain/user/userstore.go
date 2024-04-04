package user

import (
	"context"
	"fmt"
	"sync"
)

type UserStore interface {
	// Create inserts a new user into the store
	Create(ctx context.Context, user *User) error
	// List returns all users from the store
	List(ctx context.Context) ([]*User, error)
	// ByEmail returns a user by email
	ByEmail(ctx context.Context, email string) (*User, error)

	// Delete removes a user from the store. Returns true if the user was found and deleted
	Delete(ctx context.Context, email string) (bool, error)
}

var userStoreInstance UserStore
var userStoreLock sync.Mutex

func SetUserStore(us UserStore) {
	userStoreLock.Lock()
	defer userStoreLock.Unlock()
	if userStoreInstance != nil {
		panic(fmt.Sprintf("user store already set: (%T) %s", userStoreInstance, userStoreInstance))
	}
	userStoreInstance = us
}

func IsUserStoreSet() bool {
	userStoreLock.Lock()
	defer userStoreLock.Unlock()
	return userStoreInstance != nil
}

func store() UserStore {
	// Intentionally skip the sync here
	if userStoreInstance == nil {
		panic("user store not set")
	}
	return userStoreInstance
}

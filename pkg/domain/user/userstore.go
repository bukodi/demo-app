package user

import (
	"fmt"
	"sync"
)

type UserStore interface {
	// Create inserts a new user into the store
	Create(user *User) error
	// List returns all users from the store
	List() ([]*User, error)
	// ByEmail returns a user by email
	ByEmail(email string) (*User, error)
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

func store() UserStore {
	// Intentionally skip the sync here
	if userStoreInstance == nil {
		panic("user store not set")
	}
	return userStoreInstance
}

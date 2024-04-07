package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
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

type instanceWrapper struct {
	instance UserStore
}

type nilUserStore struct {
	initErr error
}

var _ UserStore = (*nilUserStore)(nil)

func (n nilUserStore) Create(ctx context.Context, user *User) error { return n.initErr }
func (n nilUserStore) List(ctx context.Context) ([]*User, error)    { return nil, n.initErr }
func (n nilUserStore) ByEmail(ctx context.Context, email string) (*User, error) {
	return nil, n.initErr
}
func (n nilUserStore) Delete(ctx context.Context, email string) (bool, error) {
	return false, n.initErr
}

var userStoreLock sync.Mutex
var userStoreProviders = make(map[string]func(context.Context) (UserStore, error))

var userStoreAtomic atomic.Value

func init() {
	userStoreAtomic.Store(&instanceWrapper{&nilUserStore{errStoreNotInitialized}})
}

var userStoreInitContext = context.Background()
var userStoreInitTimeout = time.Second * 5

var errBadProgramStructure = errors.New("bad program structure, initialization must be single-thread and must happen before any user store usage")
var errStoreNotInitialized = errors.New("user store not initialized")

func RegisterUserStoreProvider(providerName string, initFn func(context.Context) (UserStore, error)) {
	if ok := userStoreLock.TryLock(); !ok {
		panic(fmt.Errorf("can't register user store provider. %w", errBadProgramStructure))
	}
	defer userStoreLock.Unlock()
	if initFn == nil {
		return
	}

	if userStoreProviders[providerName] != nil {
		slog.Warn("user store provider already registered", "providerName", providerName)
		return
	} else {
		userStoreProviders[providerName] = initFn
	}
}

func ResetUserStore() {
	userStoreLock.Lock()
	defer userStoreLock.Unlock()
	prevStore := userStoreAtomic.Swap(&instanceWrapper{&nilUserStore{errStoreNotInitialized}}).(*instanceWrapper).instance
	if _, isNil := prevStore.(*nilUserStore); !isNil {
		if closer, ok := prevStore.(interface{ Close() error }); ok {
			closer.Close()
		}
	}
}

func IsUserStoreSet() (bool, error) {
	prevStore := userStoreAtomic.Load().(*instanceWrapper).instance
	if nus, isNil := prevStore.(*nilUserStore); isNil {
		return !errors.Is(nus.initErr, errStoreNotInitialized), nus.initErr
	} else {
		return true, nil
	}
}

func store() UserStore {
	actualStore := userStoreAtomic.Load().(*instanceWrapper).instance
	if nilStore, isNilStore := actualStore.(*nilUserStore); !isNilStore {
		// valid store set
		return actualStore
	} else if !errors.Is(nilStore.initErr, errStoreNotInitialized) {
		// Initialization already happened, but ended with error
		return nilStore
	}

	// Initialization part
	if ok := userStoreLock.TryLock(); !ok {
		panic(fmt.Errorf("can't acquire user store. %w", errBadProgramStructure))
	}
	defer userStoreLock.Unlock()

	var initResults []UserStore
	var initLogs []error
	var wasError bool
	for provName, initFn := range userStoreProviders {
		ctx, _ := context.WithTimeoutCause(userStoreInitContext, userStoreInitTimeout, fmt.Errorf("timeout for user store initialization (%s)", userStoreInitTimeout))
		us, err := initFn(ctx)
		var initLog error
		if err != nil {
			wasError = true
			initLog = fmt.Errorf("initializer %s failed with: %w", provName, err)
		} else if us == nil {
			initLog = fmt.Errorf("initializer %s returned with nil", provName)
		} else {
			initResults = append(initResults, us)
			initLog = fmt.Errorf("initializer %s succeeded: %v", provName, us)
		}
		initLogs = append(initLogs, initLog)
	}

	var initError error
	var usInstance UserStore
	if wasError {
		initError = fmt.Errorf("user store initialization failed: %w", errors.Join(initLogs...))
	} else if len(initResults) == 0 {
		if len(userStoreProviders) == 0 {
			initError = fmt.Errorf("user store not initialized, because no provider registered")
		} else {
			initError = fmt.Errorf("no user store initialized: %w", errors.Join(initLogs...))
		}
	} else if len(initResults) > 1 {
		for _, us := range initResults {
			if closer, ok := us.(interface{ Close() error }); ok {
				err := closer.Close()
				if err != nil {
					initLogs = append(initLogs, fmt.Errorf("error closing %v store: %w", us, err))
				}
			}
		}
		initError = fmt.Errorf("multiple user stores initialized: %w", errors.Join(initLogs...))
	} else if len(initResults) == 1 {
		slog.Debug("user store initialized", "initLog", errors.Join(initLogs...).Error())
		usInstance = initResults[0]
	} else {
		panic("invalid case")
	}

	if initError != nil {
		slog.Error("user store initialization failed", "err", initError)
		usInstance = &nilUserStore{initError}
	}

	userStoreAtomic.Store(&instanceWrapper{usInstance})
	return usInstance
}

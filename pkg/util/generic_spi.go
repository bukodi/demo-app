package util

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type GenericSpi[T any] struct {
	mutex        sync.Mutex
	providers    map[string]func(context.Context) (T, error)
	actualAtomic atomic.Value
	initContext  context.Context
	initTimeout  time.Duration
}

// instanceWrapper is a wrapper for the SPI instance, it is necessary for atomic.Value
type instanceWrapper[T any] struct {
	instance T
	initErr  error
}

func NewGenericSpi[T any]() *GenericSpi[T] {
	spi := &GenericSpi[T]{
		providers: make(map[string]func(context.Context) (T, error)),
	}
	spi.actualAtomic.Store(&instanceWrapper[T]{initErr: errStoreNotInitialized})
	spi.initContext = context.Background()
	spi.initTimeout = time.Second * 5
	return spi
}

var errBadProgramStructure = errors.New("bad program structure, provider registration must be single-thread and must happen before any other usage than Reset()")
var errStoreNotInitialized = errors.New("not initialized")

func (spi *GenericSpi[T]) isNil(t T) bool {
	ptr := (*uintptr)(unsafe.Pointer(&t))
	return *ptr == 0
}

func (spi *GenericSpi[T]) RegisterProvider(providerName string, initFn func(context.Context) (T, error)) {
	if ok := spi.mutex.TryLock(); !ok {
		panic(fmt.Errorf("can't register provider. %w", errBadProgramStructure))
	}
	defer spi.mutex.Unlock()
	if initFn == nil {
		return
	}

	if spi.providers[providerName] != nil {
		slog.Warn("provider already registered", "name", providerName)
		return
	} else {
		spi.providers[providerName] = initFn
	}
}

func (spi *GenericSpi[T]) Reset() {
	spi.mutex.Lock()
	defer spi.mutex.Unlock()
	prevStore := spi.actualAtomic.Swap(&instanceWrapper[T]{initErr: errStoreNotInitialized}).(*instanceWrapper[T]).instance
	if !spi.isNil(prevStore) {
		// TODO: implement closer
		// if closer, ok := prevStore.(interface{ Close() error }); ok {
		//	closer.Close()
		//}
	}
}

func (spi *GenericSpi[T]) IsSet() (bool, error) {
	w := spi.actualAtomic.Load().(*instanceWrapper[T])
	if spi.isNil(w.instance) {
		if errors.Is(w.initErr, errStoreNotInitialized) {
			return false, nil
		} else {
			return true, w.initErr
		}
	} else {
		return true, nil
	}
}

func (spi *GenericSpi[T]) Actual() (T, error) {
	prevWrapper := spi.actualAtomic.Load().(*instanceWrapper[T])
	if !spi.isNil(prevWrapper.instance) {
		// valid store set
		return prevWrapper.instance, nil
	} else if !errors.Is(prevWrapper.initErr, errStoreNotInitialized) {
		// Initialization already happened, but ended with error
		return prevWrapper.instance, prevWrapper.initErr
	}

	// Initialization part
	if ok := spi.mutex.TryLock(); !ok {
		panic(fmt.Errorf("can't acquire lock. %w", errBadProgramStructure))
	}
	defer spi.mutex.Unlock()

	var initResults []T
	var initLogs []error
	var wasError bool
	for provName, initFn := range spi.providers {
		ctx, _ := context.WithTimeoutCause(spi.initContext, spi.initTimeout, fmt.Errorf("timeout for initialization (%s)", spi.initTimeout))
		us, err := initFn(ctx)
		var initLog error
		if err != nil {
			wasError = true
			initLog = fmt.Errorf("initializer %s failed with: %w", provName, err)
		} else if spi.isNil(us) {
			initLog = fmt.Errorf("initializer %s returned with nil", provName)
		} else {
			initResults = append(initResults, us)
			initLog = fmt.Errorf("initializer %s succeeded: %v", provName, us)
		}
		initLogs = append(initLogs, initLog)
	}

	var initError error
	var usInstance T
	if wasError {
		initError = fmt.Errorf("initialization failed: %w", errors.Join(initLogs...))
	} else if len(initResults) == 0 {
		if len(spi.providers) == 0 {
			initError = fmt.Errorf("not initialized, because no provider registered")
		} else {
			initError = fmt.Errorf("no sucess initialization: %w", errors.Join(initLogs...))
		}
	} else if len(initResults) > 1 {
		for _, us := range initResults {
			// TODO: implement closer
			_ = us
			//if closer, ok := us.(interface{ Close() error }); ok {
			//	err := closer.Close()
			//	if err != nil {
			//		initLogs = append(initLogs, fmt.Errorf("error closing %v store: %w", us, err))
			//	}
		}
		initError = fmt.Errorf("multiple success initialization: %w", errors.Join(initLogs...))
	} else if len(initResults) == 1 {
		slog.Debug("initialized", "initLog", errors.Join(initLogs...).Error())
		usInstance = initResults[0]
	} else {
		panic("invalid case")
	}

	newWrapper := &instanceWrapper[T]{}
	if initError != nil {
		slog.Error("initialization failed", "err", initError)
		newWrapper.initErr = initError
	} else {
		newWrapper.instance = usInstance
	}

	spi.actualAtomic.Store(newWrapper)
	return newWrapper.instance, newWrapper.initErr
}

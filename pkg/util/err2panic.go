package util

import "testing"

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func MustT[T any](v T, err error, t *testing.T) T {
	if err != nil {
		t.Fatalf("%+v", err)
	}
	return v
}

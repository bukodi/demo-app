package util

import (
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func TestMustSimple(t *testing.T) {
	fn := func() (string, error) {
		return "success", nil
	}
	assert.Equal(t, Must(fn()), "success")

}

func TestMustMatrix(t *testing.T) {
	tests := []struct {
		name        string
		fn          func() (interface{}, error)
		want        interface{}
		shouldPanic bool
	}{
		{
			name: "success_return_value",
			fn: func() (interface{}, error) {
				return "success", nil
			},
			want:        "success",
			shouldPanic: false,
		},
		{
			name: "error_panics",
			fn: func() (interface{}, error) {
				return nil, errors.New("something went wrong")
			},
			want:        nil,
			shouldPanic: true,
		},
		{
			name: "success_int_value",
			fn: func() (interface{}, error) {
				return 42, nil
			},
			want:        42,
			shouldPanic: false,
		},
		{
			name: "nil_return_no_error",
			fn: func() (interface{}, error) {
				return nil, nil
			},
			want:        nil,
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				// Verify that Must panics when it is expected to
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Must() did not panic when it should have")
					}
				}()
				Must(tt.fn())
			} else {
				// Verify that Must returns the expected value on success
				got := Must(tt.fn())
				if got != tt.want {
					t.Errorf("Must() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

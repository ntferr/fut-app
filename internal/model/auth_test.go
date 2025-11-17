//go:build unit

package model

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	t.Parallel()
	is := require.New(t)
	tests := []struct {
		name        string
		model       AuthRequest
		hasErr      bool
		expectedErr error
	}{
		{
			name: "given validate correct AuthRequest, should be not error",
			model: AuthRequest{
				User:      "admin",
				Password:  "test123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			hasErr: false,
		},
		{
			name: "given an authRequest without user, should return an error",
			model: AuthRequest{
				Password:  "test123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			hasErr:      true,
			expectedErr: errors.New("user is required"),
		},
		{
			name: "given an authRequest without passowrd, should return an error",
			model: AuthRequest{
				User:      "admin",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.model.Validate()

			if test.hasErr {
				if test.model.User == "" {
					is.Equal(test.expectedErr, err)
					return
				}
				if test.model.Password == "" {
					is.Equal(test.expectedErr, err)
					return
				}
			}
			is.Equal(test.expectedErr, nil)
		})
	}
}

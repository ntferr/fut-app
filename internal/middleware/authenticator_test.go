//go:build unit

package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/fut-app/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestRemoveBearerPrefix(t *testing.T) {
	t.Parallel()
	is := require.New(t)
	tests := []struct {
		name          string
		authHeader    string
		hasError      bool
		expectedToken string
		expectedErr   error
	}{
		{
			name:          "when token is valid, should return token without bearer prefix",
			authHeader:    "Bearer xpto",
			hasError:      false,
			expectedToken: "xpto",
			expectedErr:   nil,
		},
		{
			name:        "when token is empty, should return error",
			authHeader:  "",
			hasError:    true,
			expectedErr: NewJWTErr(nil, "failed to split bearer from auth header"),
		},
		{
			name:        "when token has more than two parts, should return error",
			authHeader:  "Bearer translate programming",
			hasError:    true,
			expectedErr: NewJWTErr(nil, "failed to split bearer from auth header"),
		},
		{
			name:        "when token doesn't have bearer prefix, should return error",
			authHeader:  "xpto",
			hasError:    true,
			expectedErr: NewJWTErr(nil, "failed to split bearer from auth header"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			token, err := RemoveBearerPrefix(test.authHeader)
			if test.hasError {
				is.Equal(test.expectedErr, err)
				return
			}
			is.Nil(test.expectedErr)
			is.Equal(test.expectedToken, token)
		})
	}
}

func TestVerifyToken(t *testing.T) {
	t.Parallel()
	is := require.New(t)

	aReq := model.AuthRequest{
		ID:       2,
		User:     "test",
		Password: "abcd",
	}
	tokenString, err := aReq.GenerateToken("xpto")
	is.Nil(err)

	rsaToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"username": "testuser",
	})
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	is.Nil(err)
	strRSA, err := rsaToken.SignedString(privateKey)
	is.Nil(err)

	tests := []struct {
		name           string
		token          string
		secretKey      string
		hasErr         bool
		ExpectedErr    error
		ExpectedClaims jwt.MapClaims
	}{
		{
			name:        "when token and secret key is valid, should return jwt.MapClaims",
			token:       tokenString,
			secretKey:   "xpto",
			hasErr:      false,
			ExpectedErr: nil,
			ExpectedClaims: jwt.MapClaims{
				"user_id":  float64(2),
				"username": "test",
			},
		},
		{
			name:        "when token is empty, should return error",
			token:       "",
			secretKey:   "xpto",
			hasErr:      true,
			ExpectedErr: NewJWTErr(err, "this token isn't valid"),
		},
		{
			name:        "when secret key is empty, should return error",
			token:       tokenString,
			secretKey:   strRSA,
			hasErr:      true,
			ExpectedErr: NewJWTErr(err, "this token isn't valid"),
		},
		{
			name:        "when it's not a token, should return error",
			token:       strRSA,
			secretKey:   "xpto",
			hasErr:      true,
			ExpectedErr: NewJWTErr(nil, "this isn't a jwt token"),
		},
		{
			name:        "when secret is wrong, should return error",
			token:       tokenString,
			secretKey:   "test",
			hasErr:      true,
			ExpectedErr: NewJWTErr(err, "this token isn't valid"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			claims, err := VerifyToken(test.token, test.secretKey)

			if test.hasErr {
				is.Equal(test.ExpectedErr.Error(), err.Error())
				return
			}
			is.Nil(err)
			is.Equal(test.ExpectedClaims["username"], claims["username"])
			is.Equal(test.ExpectedClaims["user_id"], claims["user_id"])
			is.NotNil(claims["exp"])
		})
	}
}

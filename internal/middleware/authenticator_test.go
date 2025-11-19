//go:build unit

package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fut-app/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
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
			expectedErr: NewJWTErr(nil, "invalid authorization type: failed to split bearer from auth header"),
		},
		{
			name:        "when token has more than two parts, should return error",
			authHeader:  "Bearer translate programming",
			hasError:    true,
			expectedErr: NewJWTErr(nil, "invalid authorization type: failed to split bearer from auth header"),
		},
		{
			name:        "when token doesn't have bearer prefix, should return error",
			authHeader:  "xpto",
			hasError:    true,
			expectedErr: NewJWTErr(nil, "invalid authorization type: failed to split bearer from auth header"),
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
	fmt.Println(tokenString)
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
				"username": "test",
			},
		},
		{
			name:        "when token is empty, should return error",
			token:       "",
			secretKey:   "xpto",
			hasErr:      true,
			ExpectedErr: NewJWTErr(err, "invalid or expired token"),
		},
		{
			name:        "when secret key is empty, should return error",
			token:       tokenString,
			secretKey:   "",
			hasErr:      true,
			ExpectedErr: NewJWTErr(err, "invalid or expired token"),
		},
		{
			name:        "when it's not a token, should return error",
			token:       strRSA,
			secretKey:   "xpto",
			hasErr:      true,
			ExpectedErr: NewJWTErr(err, "invalid or expired token"),
		},
		{
			name:        "when secret is wrong, should return error",
			token:       tokenString,
			secretKey:   "test",
			hasErr:      true,
			ExpectedErr: NewJWTErr(err, "invalid or expired token"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			claims, err := VerifyToken(test.token, test.secretKey)
			if test.hasErr {
				is.Equal(test.ExpectedErr.Error(), err.Error())
				return
			}
			is.Nil(err)
			is.Equal(test.ExpectedClaims["username"], claims["username"])
			is.NotNil(claims["exp"])
		})
	}

	t.Run("when claims not map claims, should return error", func(t *testing.T) {
		t.Parallel()
		type CustomClaims struct {
			UserID int `json:"user_id"`
			jwt.RegisteredClaims
		}

		customClaims := CustomClaims{
			UserID: 789,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		}

		customToken := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
		customTokenString, _ := customToken.SignedString([]byte("xpto"))
		_, err := VerifyToken(customTokenString, "xpto")
		is.Equal(NewJWTErr(nil, "invalid token: claims are not MapClaims"), err)
	})
}

func TestJWTMiddleware(t *testing.T) {
	is := require.New(t)
	t.Run("when token and secret key are correct, should apply user into context", func(t *testing.T) {
		secretKey := "test_secret"
		var claims model.AuthClaims
		claims = model.AuthClaims{
			RegisteredClaims: jwt.RegisteredClaims{},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte(secretKey))
		is.Nil(err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		username := "adm"

		next := func(c echo.Context) error {
			return nil
		}

		req.Header.Set("Authorization", "Bearer "+tokenStr)

		middlewareFunc := JWTMiddleware(secretKey)
		err = middlewareFunc(next)(c)

		expectedUser := c.Get("user")

		is.Nil(err)
		is.Equal(http.StatusOK, rec.Code)
		is.Equal(username, expectedUser.(string))
	})

	t.Run("when token doesn't have username claim, should return error", func(t *testing.T) {
		type testModel struct {
			ID int
			jwt.RegisteredClaims
		}
		secretKey := "test_secret"
		claims := testModel{
			ID:               2,
			RegisteredClaims: jwt.RegisteredClaims{},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte(secretKey))
		is.Nil(err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		next := func(c echo.Context) error {
			return nil
		}

		req.Header.Set("Authorization", "Bearer "+tokenStr)

		middlewareFunc := JWTMiddleware(secretKey)
		err = middlewareFunc(next)(c)

		is.NotNil(err)
		is.Equal("code=401, message=final verify: invalid token", err.Error())
	})

	// TODO: break in prefix
	// TODO: break in verify
}

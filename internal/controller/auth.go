package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fut-app/internal/model"
	"github.com/fut-app/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	Service   service.AuthenticationRepository
	SecretKey string
}

func NewAuth(db service.Database, secretKey string) Auth {
	return Auth{
		Service:   &db,
		SecretKey: secretKey,
	}
}

func (a Auth) Authenticate(c echo.Context) error {
	var req model.AuthRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid JSON",
			},
		)
	}
	if err := req.Validate(); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			err,
		)
	}
	if err := req.HashPassword(); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "somenthing went wrong with given password",
			},
		)
	}
	storedAuth, err := a.Service.FindCredentials(c.Request().Context(), req)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			err,
		)
	}

	token, err := generateToken(*storedAuth, a.SecretKey)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": "something is wrong with provided credentials",
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		token,
	)
}

func generateToken(credentials model.AuthRequest, sk string) (string, error) {
	claims := &model.CustomClaims{
		AuthRequest: credentials,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 horas
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "fut-app",
			Subject:   "authentication",
			ID:        fmt.Sprintf("%d", credentials.ID),
		},
	}

	jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sk))
}

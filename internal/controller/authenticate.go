package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/fut-app/internal/model"
	"github.com/fut-app/internal/service"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	Service   service.CredentialsRepository
	SecretKey string
}

func NewAuth(db service.CredentialsDatabase, secretKey string) Auth {
	return Auth{
		Service:   &db,
		SecretKey: secretKey,
	}
}

func (a Auth) Authenticate(c echo.Context) error {
	var req model.AuthRequest
	if err := c.Bind(&req); err != nil {
		log.Println("failed to bind request\n", err)
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid JSON",
			},
		)
	}
	if err := req.Validate(); err != nil {
		log.Println("failed to validate auth request\n", err)
		return c.JSON(
			http.StatusBadRequest,
			err,
		)
	}

	credential, err := req.ParseAuthRequestToCredential()
	if err != nil {
		log.Println("failed to parse auth request to credential\n", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": "failed to create credential",
			},
		)
	}

	storedCredential, err := a.Service.FindCredentials(c.Request().Context(), credential)
	if err != nil {
		log.Println("failed to find credentials\n", err)
		return c.JSON(
			http.StatusInternalServerError,
			err,
		)
	}

	token, err := storedCredential.GenerateToken(a.SecretKey)
	if err != nil {
		log.Println("failed to generate token\n", err)
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

func (a Auth) CreateCredentials(c echo.Context) error {
	var req model.AuthRequest
	if err := c.Bind(&req); err != nil {
		log.Println("failed to bind request\n", err)
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid JSON",
			},
		)
	}
	if err := req.Validate(); err != nil {
		log.Println("failed to validate auth request\n", err)
		return c.JSON(
			http.StatusBadRequest,
			err,
		)
	}

	credential, err := req.ParseAuthRequestToCredential()
	if err != nil {
		log.Println("failed to parse auth request to credential\n", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": "failed to create credential",
			},
		)
	}

	credential.CreatedAt = time.Now()
	credential.UpdatedAt = time.Now()

	if err := a.Service.CreateCredentials(c.Request().Context(), credential); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": "failed to create credentials",
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		map[string]string{
			"message": "creation is succefull",
		},
	)
}

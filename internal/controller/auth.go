package controller

import (
	"net/http"

	"github.com/fut-app/internal/model"
	"github.com/fut-app/internal/service"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	Service service.AuthenticationRepository
}

func NewAuth(db service.Database) Auth {
	return Auth{
		Service: &db,
	}
}

func (a Auth) AuthHandler(c echo.Context) error {
	var req model.AuthRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid JSON",
			},
		)
	}
	if err := req.Validate(); err != nil {
		c.JSON(
			http.StatusBadRequest,
			err,
		)
	}
	if err := req.HashPassword(); err != nil {
		c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "somenthing went wrong with given password",
			},
		)
	}
	err := a.Service.FindCredentials(c.Request().Context(), req)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			err,
		)
	}

	return nil
}

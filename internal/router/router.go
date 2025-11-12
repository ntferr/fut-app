package router

import (
	"github.com/fut-app/internal/controller"
	"github.com/labstack/echo/v4"
)

func Setup(app *echo.Echo, control *controller.Controller) {
	app.POST("auth/login", control.Auth.AuthHandler)
}

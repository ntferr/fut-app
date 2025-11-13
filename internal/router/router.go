package router

import (
	"github.com/fut-app/internal/controller"
	"github.com/fut-app/internal/middleware"
	"github.com/labstack/echo/v4"
)

func Setup(app *echo.Echo, control *controller.Controller, secretKey string) {
	app.POST("auth/login", control.Auth.Authenticate)

	protected := app.Group("")
	protected.Use(middleware.JWTMiddleware(secretKey))
	protected.GET("/campeonatos", control.Champion.Championship)
}

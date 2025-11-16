package router

import (
	"github.com/fut-app/internal/controller"
	"github.com/fut-app/internal/middleware"
	"github.com/labstack/echo/v4"
)

func Setup(app *echo.Echo, control *controller.Controller, secretKey string) {
	app.POST("auth/login", control.Auth.Authenticate)
	app.POST("auth/create", control.Auth.CreateCredentials)

	protected := app.Group("/campeonatos", middleware.JWTMiddleware(secretKey))
	protected.GET("/", control.Champion.Championship)
}

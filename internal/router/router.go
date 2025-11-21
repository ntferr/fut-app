package router

import (
	"github.com/fut-app/internal/controller"
	"github.com/fut-app/internal/middleware"
	"github.com/labstack/echo/v4"
)

func Setup(app *echo.Echo, control *controller.Controller, secretKey string) {
	authEndpoints(
		app.Group("/auth"),
		&control.Auth,
		secretKey,
	)
	championshipEndpoints(
		app.Group("/campeonatos", middleware.JWTMiddleware(secretKey)),
		&control.Champion,
	)
}

func authEndpoints(auth *echo.Group, control *controller.Auth, secretKey string) {
	auth.POST("/login", control.Authenticate)
	auth.POST("/create", control.CreateCredentials, middleware.JWTMiddleware(secretKey))
}

func championshipEndpoints(champion *echo.Group, control *controller.Champion) {
	champion.GET("/", control.Championship)
	// TODO: endpoint for filters
}

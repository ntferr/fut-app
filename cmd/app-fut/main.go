package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fut-app/internal/config"
	"github.com/fut-app/internal/controller"
	"github.com/fut-app/internal/router"
	"github.com/fut-app/internal/service"
	"github.com/fut-app/pkg/gorm"
	"github.com/fut-app/pkg/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func shutdownAPP(app *echo.Echo) {
	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		<-quit
		app.Logger.Print("gracefully shutting down...")

		if err := app.Shutdown(ctx); err != nil {
			app.Logger.Errorf("error during shutdown: %v", err)
		}

		app.Logger.Info("server shutdown completed")
	}()
}

func addressFormater(host string, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}

func main() {
	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	cfg, err := config.NewConfig()
	if err != nil {
		app.Logger.Fatalf("failed to instatiate configs: %v", err)
	}

	dialector := postgres.NewPostgres(cfg.Postgres)
	db := gorm.NewGorm(dialector)
	authService := service.NewDatabase(db)
	authController := controller.NewAuth(authService)
	control := controller.NewController(authController)

	router.Setup(app, control)

	err = app.Start(
		addressFormater(
			cfg.App.Host,
			cfg.App.Port,
		),
	)

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}

	shutdownAPP(app)
}

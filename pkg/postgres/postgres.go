package postgres

import (
	"fmt"

	"github.com/fut-app/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(cfg config.PostgresConfig) gorm.Dialector {
	return postgres.New(postgres.Config{
		DSN:                  mountDNS(cfg),
		PreferSimpleProtocol: true,
	})
}

func mountDNS(cfg config.PostgresConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
	)
}

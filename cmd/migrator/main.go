package main

import (
	"log"

	"github.com/fut-app/internal/config"
	"github.com/fut-app/internal/model"
	"github.com/fut-app/pkg/gorm"
	"github.com/fut-app/pkg/postgres"
)

const filename = "payload.json"

func main() {
	log.Println("migrator: starting migration")
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("migrator: failed to initiate config")
	}
	dialector := postgres.NewPostgres(cfg.Postgres)
	db := gorm.NewGorm(dialector)

	log.Println("migrator: initiang automigrate")

	err = db.AutoMigrate(
		model.AuthRequest{},
	)
	if err != nil {
		log.Fatalf("migrator: failed to do automigrate: %s", err.Error())
	}

	log.Println("migrator: sucessfuly automigrate!")
}

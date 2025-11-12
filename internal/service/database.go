package service

import (
	"context"
	"errors"

	"github.com/fut-app/internal/model"
	"gorm.io/gorm"
)

type Database struct {
	Gorm *gorm.DB
}

func NewDatabase(gorm *gorm.DB) Database {
	return Database{
		Gorm: gorm,
	}
}

type AuthenticationRepository interface {
	FindCredentials(ctx context.Context, credentials model.AuthRequest) error
}

func (d *Database) FindCredentials(ctx context.Context, credentials model.AuthRequest) error {
	var dbCredentials model.AuthRequest
	result := d.Gorm.Where("user", dbCredentials).First(&dbCredentials)
	if result.Error != nil {
		return errors.New("user not found")
	}

	if !dbCredentials.CheckPassword(credentials.Password) {
		return errors.New("incorrect password")
	}

	return nil
}

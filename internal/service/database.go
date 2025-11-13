package service

import (
	"context"
	"errors"

	"github.com/fut-app/internal/model"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound    = errors.New("invalid credentials")
	ErrInvalidPassword = errors.New("invalid credentials")
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
	FindCredentials(ctx context.Context, credentials model.AuthRequest) (*model.AuthRequest, error)
}

func (d *Database) FindCredentials(ctx context.Context, credentials model.AuthRequest) (*model.AuthRequest, error) {
	var storedAuth model.AuthRequest
	result := d.Gorm.
		WithContext(ctx).
		Where("user", storedAuth).
		First(&storedAuth)
	if result.Error != nil {
		return nil, ErrUserNotFound
	}

	if !storedAuth.CheckPassword(credentials.Password) {
		return nil, ErrInvalidPassword
	}

	return &storedAuth, nil
}

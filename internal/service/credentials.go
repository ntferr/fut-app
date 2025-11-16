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
	ErrCreationFailed  = errors.New("failed to create new credentials")
)

type CredentialsDatabase struct {
	Gorm *gorm.DB
}

func NewDatabase(gorm *gorm.DB) CredentialsDatabase {
	return CredentialsDatabase{
		Gorm: gorm,
	}
}

type CredentialsRepository interface {
	FindCredentials(ctx context.Context, credentials model.AuthRequest) (*model.AuthRequest, error)
	CreateCredentials(ctx context.Context, credentials model.AuthRequest) error
}

func (d *CredentialsDatabase) FindCredentials(ctx context.Context, credentials model.AuthRequest) (*model.AuthRequest, error) {
	var storedCredentias model.AuthRequest
	result := d.Gorm.
		WithContext(ctx).
		Where("user", storedCredentias).
		First(&storedCredentias)
	if result.Error != nil {
		return nil, ErrUserNotFound
	}

	if !storedCredentias.CheckPassword(credentials.Password) {
		return nil, ErrInvalidPassword
	}

	return &storedCredentias, nil
}

func (d *CredentialsDatabase) CreateCredentials(ctx context.Context, credentials model.AuthRequest) error {
	err := d.Gorm.Create(&credentials).Error
	if err != nil {
		return ErrCreationFailed
	}
	return nil
}

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
	FindCredentials(ctx context.Context, credentials *model.Credential) (*model.Credential, error)
	CreateCredentials(ctx context.Context, credentials *model.Credential) error
}

func (d *CredentialsDatabase) FindCredentials(ctx context.Context, credential *model.Credential) (*model.Credential, error) {
	var storedCredential model.Credential
	result := d.Gorm.
		WithContext(ctx).
		Where("user", credential.User).
		First(&storedCredential)
	if result.Error != nil {
		return nil, ErrUserNotFound
	}

	return &storedCredential, nil
}

func (d *CredentialsDatabase) CreateCredentials(ctx context.Context, credential *model.Credential) error {
	err := d.Gorm.Create(&credential).Error
	if err != nil {
		return ErrCreationFailed
	}
	return nil
}

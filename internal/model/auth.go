package model

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthRequest struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	User      string    `json:"user" gorm:"user"`
	Password  string    `json:"password" gorm:"password"`
	CreatedAt time.Time `gorm:"createdAt"`
	UpdatedAt time.Time `gorm:"updatedAt"`
}

func (a *AuthRequest) Validate() error {
	if a.User == "" {
		return errors.New("user is required")
	}
	if a.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (a *AuthRequest) GenerateToken(secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"id":   a.ID,
		"user": a.User,
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSttring, err := token.SignedString(secretKey)
	if err != nil {
		return "", errors.New("failed to generate token")
	}
	return tokenSttring, nil
}

func (a *AuthRequest) CheckPassword(password string) bool {
	return a.Password == password
}

type CustomClaims struct {
	AuthRequest
	jwt.RegisteredClaims
}

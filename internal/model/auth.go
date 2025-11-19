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

type AuthClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
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
	now := time.Now()
	claims := AuthClaims{
		Username: a.User,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 1)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   a.User,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSttring, err := token.SignedString([]byte(secretKey))
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

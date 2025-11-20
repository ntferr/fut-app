package model

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
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

func (a *AuthRequest) ParseAuthRequestToCredential() (*Credential, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	credential := Credential{
		User:              a.User,
		EncryptedPassword: string(encryptedPassword),
	}

	return &credential, nil
}

type Credential struct {
	ID                int       `gorm:"primaryKey"`
	User              string    `gorm:"user"`
	EncryptedPassword string    `gorm:"encryptedPassword"`
	CreatedAt         time.Time `gorm:"createdAt"`
	UpdatedAt         time.Time `gorm:"updatedAt"`
}

func (a *Credential) GenerateToken(secretKey string) (string, error) {
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

func (c *Credential) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.EncryptedPassword), []byte(password))
	return err == nil
}

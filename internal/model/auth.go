package model

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
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

func (a *AuthRequest) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to hash password: ", err)
		return err
	}
	a.Password = string(hash)
	return nil
}

func (a *AuthRequest) CheckPassword(password string) bool {
	return a.Password == password
}

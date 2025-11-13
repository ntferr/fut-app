package gorm

import (
	"log"

	"gorm.io/gorm"
)

// NewGorm initialize session with database.
func NewGorm(dialector gorm.Dialector) *gorm.DB {
	db, err := gorm.Open(dialector)
	if err != nil {
		log.Fatalf("failed to initialize session with database: %v", err)
	}
	return db
}

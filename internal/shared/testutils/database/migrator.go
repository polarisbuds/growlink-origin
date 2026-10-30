package database

import (
	infrastructure "growlink/internal/infrastructure/postgresql"

	"gorm.io/gorm"
)

// Migrator applies database schema migrations.
type Migrator interface {
	Apply(db *gorm.DB) error
}

// AutoMigrator applies GORM AutoMigrate for all domain entities.
type AutoMigrator struct{}

func (AutoMigrator) Apply(db *gorm.DB) error {
	return infrastructure.AutoMigrate(db)
}

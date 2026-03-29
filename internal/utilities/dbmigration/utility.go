package dbmigration

import (
	"context"

	"gorm.io/gorm"
)

type MigrationRequest struct {
	ID        string
	Migration func(tx *gorm.DB) error
	Rollback  func(tx *gorm.DB) error
}

type Utility interface {
	RegisterMigration(request MigrationRequest)
	Migrate(ctx context.Context) error
}

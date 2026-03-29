package dbmigration

import (
	"context"
	"errors"
	"log/slog"

	"github.com/malcolmkzh/study-notifier/internal/utilities/db"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Implementation struct {
	db         *gorm.DB
	migrations []MigrationRequest
}

func NewMigrationUtility(dbUtility db.Utility) (*Implementation, error) {
	if dbUtility == nil {
		return nil, errors.New("db utility is required")
	}

	return &Implementation{
		db:         dbUtility.DB(),
		migrations: []MigrationRequest{},
	}, nil
}

func (m *Implementation) RegisterMigration(request MigrationRequest) {
	m.migrations = append(m.migrations, request)
}

func (m *Implementation) Migrate(ctx context.Context) error {
	migrations := make([]*gormigrate.Migration, len(m.migrations))
	for i, migration := range m.migrations {
		migrations[i] = &gormigrate.Migration{
			ID:       migration.ID,
			Migrate:  m.migrations[i].Migration,
			Rollback: m.migrations[i].Rollback,
		}
	}

	migrator := gormigrate.New(m.db.WithContext(ctx), gormigrate.DefaultOptions, migrations)
	if err := migrator.Migrate(); err != nil {
		return err
	}

	slog.Info("migration ran successfully")
	return nil
}

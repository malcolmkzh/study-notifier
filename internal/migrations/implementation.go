package migrations

import (
	"context"
	"fmt"
	"os"

	"github.com/malcolmkzh/study-notifier/internal/utilities/dbmigration"
	"gorm.io/gorm"
)

type Implementation struct{}

func NewMigration(ctx context.Context, dbMigrationUtility dbmigration.Utility) (*Implementation, error) {
	_ = ctx

	m := &Implementation{}

	dbMigrationUtility.RegisterMigration(dbmigration.MigrationRequest{
		ID:        "2026032800000000001",
		Migration: m.generateMigration("/internal/migrations/scripts/create_accounts_table.sql"),
		Rollback:  m.generateMigration("/internal/migrations/scripts/create_accounts_table_rollback.sql"),
	})

	dbMigrationUtility.RegisterMigration(dbmigration.MigrationRequest{
		ID:        "2026032800000000002",
		Migration: m.generateMigration("/internal/migrations/scripts/create_notes_table.sql"),
		Rollback:  m.generateMigration("/internal/migrations/scripts/create_notes_table_rollback.sql"),
	})

	dbMigrationUtility.RegisterMigration(dbmigration.MigrationRequest{
		ID:        "2026032900000000003",
		Migration: m.generateMigration("/internal/migrations/scripts/create_questions_table.sql"),
		Rollback:  m.generateMigration("/internal/migrations/scripts/create_questions_table_rollback.sql"),
	})

	dbMigrationUtility.RegisterMigration(dbmigration.MigrationRequest{
		ID:        "2026040500000000004",
		Migration: m.generateMigration("/internal/migrations/scripts/create_jobs_table.sql"),
		Rollback:  m.generateMigration("/internal/migrations/scripts/create_jobs_table_rollback.sql"),
	})

	dbMigrationUtility.RegisterMigration(dbmigration.MigrationRequest{
		ID:        "2026040500000000005",
		Migration: m.generateMigration("/internal/migrations/scripts/create_reminders_table.sql"),
		Rollback:  m.generateMigration("/internal/migrations/scripts/create_reminders_table_rollback.sql"),
	})

	dbMigrationUtility.RegisterMigration(dbmigration.MigrationRequest{
		ID:        "2026040600000000006",
		Migration: m.generateMigration("/internal/migrations/scripts/add_telegram_chat_id_to_accounts.sql"),
		Rollback:  m.generateMigration("/internal/migrations/scripts/add_telegram_chat_id_to_accounts_rollback.sql"),
	})

	dbMigrationUtility.RegisterMigration(dbmigration.MigrationRequest{
		ID:        "2026040600000000007",
		Migration: m.generateMigration("/internal/migrations/scripts/create_telegram_links_table.sql"),
		Rollback:  m.generateMigration("/internal/migrations/scripts/create_telegram_links_table_rollback.sql"),
	})

	return m, nil
}

func (m *Implementation) generateMigration(file string) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		path, err := os.Getwd()
		if err != nil {
			return err
		}

		query, err := os.ReadFile(path + file)
		if err != nil {
			return fmt.Errorf("read migration file %s: %w", file, err)
		}

		return tx.Exec(string(query)).Error
	}
}

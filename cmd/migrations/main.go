package main

import (
	"context"
	"log"

	"github.com/malcolmkzh/study-notifier/internal/migrations"
	"github.com/malcolmkzh/study-notifier/internal/utilities/config"
	"github.com/malcolmkzh/study-notifier/internal/utilities/db"
	"github.com/malcolmkzh/study-notifier/internal/utilities/dbmigration"
)

func main() {
	ctx := context.Background()

	configUtility, err := config.NewConfigUtility()
	if err != nil {
		log.Fatal("Failed to initialize configuration utility: ", err)
	}

	dbUtility, err := db.NewDbUtility(configUtility)
	if err != nil {
		log.Fatal("Failed to initialize database utility: ", err)
	}

	// Initialize migration utility
	migrationUtility, err := dbmigration.NewMigrationUtility(dbUtility)
	if err != nil {
		log.Fatal("Failed to initialize migration utility: ", err)
	}

	// Initialize migrations and
	_, err = migrations.NewMigration(ctx, migrationUtility)
	if err != nil {
		log.Fatal("Failed to initialize migrations: ", err)
	}

	// Run migrations
	err = migrationUtility.Migrate(ctx)
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
}

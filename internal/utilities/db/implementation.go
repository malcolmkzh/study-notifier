package db

import (
	"fmt"
	"github.com/malcolmkzh/study-notifier/internal/utilities/config"
	"log"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Implementation struct {
	db *gorm.DB
}

func NewDbUtility(configUtility config.Utility) (*Implementation, error) {
	cfg := configUtility.Config()

	databaseConnection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&multiStatements=true",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	if cfg.Verbose {
		slog.Info("Connecting to database", "host", cfg.DBHost, "port", cfg.DBPort, "name", cfg.DBName)
	}

	logLevel := gormlogger.Silent
	if cfg.EnableDBLog {
		logLevel = gormlogger.Info
	}

	db, err := gorm.Open(mysql.Open(databaseConnection), &gorm.Config{
		Logger: gormlogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gormlogger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logLevel,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	})
	if err != nil {
		return nil, err
	}

	if cfg.Verbose {
		slog.Info("Connected to database")
	}

	return &Implementation{db: db}, nil
}

func (m *Implementation) DB() *gorm.DB {
	return m.db
}

// Package database provides the PostgreSQL database connection via GORM.
// It configures connection pooling and registers auto-migrations for all domain models.
package database

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/yourusername/go-skeleton-code/internal/domain/contact"
	"github.com/yourusername/go-skeleton-code/internal/infrastructure/config"
)

// NewPostgresDB opens a GORM connection to PostgreSQL, configures the
// connection pool, and runs auto-migrations for all registered models.
func NewPostgresDB(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	// Use a silent GORM logger in production to avoid duplicate log output;
	// use info-level in development for SQL query visibility.
	gormLogLevel := gormlogger.Silent
	if !cfg.IsProd() {
		gormLogLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("database: failed to open connection: %w", err)
	}

	// Configure connection pool via the underlying *sql.DB.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("database: failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpen)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdle)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto-migrate all domain models — add new models here as the project grows.
	if err := autoMigrate(db); err != nil {
		return nil, err
	}

	log.Info("database: connected to PostgreSQL",
		zap.String("host", cfg.DBHost),
		zap.String("dbname", cfg.DBName),
	)

	return db, nil
}

// autoMigrate runs GORM AutoMigrate for all registered entity types.
// This creates or alters tables to match the current struct definitions.
func autoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&contact.Contact{},
		// register additional domain models here, e.g. &user.User{}
	); err != nil {
		return fmt.Errorf("database: auto-migrate failed: %w", err)
	}
	return nil
}

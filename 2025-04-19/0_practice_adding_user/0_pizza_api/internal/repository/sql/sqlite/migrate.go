package sqlite

import (
	"log/slog"

	"gorm.io/gorm"
)

func migrate(db *gorm.DB) error {
	slog.Debug("Migrating database")
	return db.AutoMigrate(Customer{}, Order{})
}

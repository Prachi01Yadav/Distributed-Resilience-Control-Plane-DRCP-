package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewSqliteDB initializes a new SQLite connection for local testing
func NewSqliteDB(filepath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

package config

import (
	"log"
	"work-distributor-system/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupDatabase initializes the SQLite database connection
// and performs automatic migration for the User and Task models.
// The `dsn` parameter represents the database file name (e.g., "tasks.db").
func SetupDatabase(dsn string) *gorm.DB {

	// Opening a new SQLite database connection using GORM
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Automatically creating tables for User and Task models
	err = db.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	return db // Returns the configured GORM database instance
}

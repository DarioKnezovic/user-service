package database

import (
	"fmt"
	"log"

	"github.com/DarioKnezovic/user-service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB establishes a connection to the database and returns a GORM DB instance.
func ConnectDB() (*gorm.DB, error) {
	cfg := config.LoadConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		cfg.DatabaseURL, cfg.DatabaseUsername, cfg.DatabasePassword, cfg.DatabaseName, cfg.DatabasePort)

	log.Printf("Attempting to connect on database...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}
	log.Printf("It's connected on %s", cfg.DatabaseURL)

	return db, nil
}

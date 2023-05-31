package database

import (
	"github.com/DarioKnezovic/user-service/internal/session"
	"github.com/DarioKnezovic/user-service/internal/user"
	"gorm.io/gorm"
	"log"
)

// PerformAutoMigrations runs auto migrations for all the models in the database.
func PerformAutoMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(&user.User{})
	err = db.AutoMigrate(&session.Session{})
	if err != nil {
		return err
	}
	log.Print("Auto migrations have been successfully finished")

	return nil
}

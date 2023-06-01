package repository

import (
	"github.com/DarioKnezovic/user-service/internal/session"
	"gorm.io/gorm"
)

type SessionRepository interface {
	StoreSession(sessionData session.Session) (*session.Session, error)
	DeleteSession(userId uint) error
	FindSessionById(token string) (bool, error)
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func (r *sessionRepository) StoreSession(sessionData session.Session) (*session.Session, error) {
	err := r.db.Create(&sessionData).Error
	if err != nil {
		return nil, err
	}

	return &sessionData, nil
}

func (r *sessionRepository) DeleteSession(userId uint) error {
	// Define the raw SQL query to delete sessions by user ID
	query := "DELETE FROM sessions WHERE user_id = ?"

	// Execute the raw SQL query
	result := r.db.Exec(query, userId)
	if result.Error != nil {
		// Handle the error
		return result.Error
	}

	// Check the number of rows affected by the delete operation
	rowsAffected := result.RowsAffected
	if rowsAffected == 0 {
		// No session found for the specified user ID
		return nil
	}

	// Deletion successful
	return nil
}

func (r *sessionRepository) FindSessionById(token string) (bool, error) {
	var foundedSession session.Session
	err := r.db.Where("session_id = ?", token).First(&foundedSession).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// User not found in the database
			return false, nil
		}
		// Other error occurred
		return false, err
	}

	return true, nil
}

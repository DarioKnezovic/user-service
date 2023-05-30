package repository

import (
	"github.com/DarioKnezovic/user-service/internal/user"
	"gorm.io/gorm"
)

// UserRepository represents the user repository interface.
type UserRepository interface {
	SaveUser(newUser user.User) (*user.User, error)
	FindUserByEmail(email string) (*user.User, error)
}

// userRepository represents an implementation of UserRepository.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	// Initialize and return the repository instance
	return &userRepository{
		db: db,
	}
}

// SaveUser saves a new user to the data storage.
func (r *userRepository) SaveUser(newUser user.User) (*user.User, error) {
	err := r.db.Create(&newUser).Error
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

// FindUserByEmail retrieves a user from the data storage based on the email.
func (r *userRepository) FindUserByEmail(email string) (*user.User, error) {
	var foundUser user.User
	err := r.db.Where("email = ?", email).First(&foundUser).Error
	if err != nil {
		return nil, err
	}

	return &foundUser, nil
}

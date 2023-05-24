package service

import (
	"github.com/DarioKnezovic/user-service/internal/user"
	"github.com/DarioKnezovic/user-service/internal/user/repository"
)

// UserService represents the user service implementation.
type UserService struct {
	userRepository repository.UserRepository
}

// NewUserService creates a new instance of UserService.
func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

// RegisterUser registers a new user.
func (s *UserService) RegisterUser(newUser user.User) (*user.User, error) {
	// Add any business logic or validation for user registration
	// Call the repository to save the new user
	savedUser, err := s.userRepository.SaveUser(newUser)
	if err != nil {
		return nil, err
	}
	return savedUser, nil
}

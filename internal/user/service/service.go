package service

import (
	"errors"
	"fmt"
	"github.com/DarioKnezovic/user-service/config"
	"github.com/DarioKnezovic/user-service/internal/user"
	"github.com/DarioKnezovic/user-service/internal/user/repository"
	"github.com/DarioKnezovic/user-service/pkg/util"
	"golang.org/x/crypto/bcrypt"
	"time"
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

// HashPassword hashes the provided password using bcrypt.
func (s *UserService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// RegisterUser registers a new user.
func (s *UserService) RegisterUser(newUser user.User) (*user.User, error) {
	// Hash the password before saving the user
	hashedPassword, err := s.HashPassword(newUser.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	newUser.Password = hashedPassword

	savedUser, err := s.userRepository.SaveUser(newUser)
	if err != nil {
		return nil, err
	}
	return savedUser, nil
}

// LoginUser performs user login and returns the authenticated user with the generated authentication token.
func (s *UserService) LoginUser(loginUser user.User) (string, error) {
	// Retrieve the user from the repository based on the email
	existingUser, err := s.userRepository.FindUserByEmail(loginUser.Email)
	if err != nil {
		return "", err
	}

	// Check if the user exists
	if existingUser == nil {
		return "", errors.New("user not found")
	}

	// Compare the provided password with the hashed password in the user object
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(loginUser.Password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	cfg := config.LoadConfig()
	// Password is correct, generate the authentication token
	token, err := util.GenerateJWT(existingUser.ID, existingUser.Email, []byte(cfg.JWTSecretKey), time.Hour*24) // Adjust the secret key and token expiration time as needed
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	// Return the authenticated user and the generated token
	return token, nil
}

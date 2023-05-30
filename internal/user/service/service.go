package service

import (
	"errors"
	"github.com/DarioKnezovic/user-service/config"
	"github.com/DarioKnezovic/user-service/internal/user"
	"github.com/DarioKnezovic/user-service/internal/user/repository"
	"github.com/DarioKnezovic/user-service/pkg/util"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var UserErrors = map[string]error{
	// ErrUserNotFound is returned when a user is not found in the repository.
	"ErrUserNotFound": errors.New("user not found"),

	// ErrInvalidPassword is returned when the provided password is invalid.
	"ErrInvalidPassword": errors.New("invalid password"),
}

const (
	ERR_USER_NOT_FOUND   = "ErrUserNotFound"
	ERR_INVALID_PASSWORD = "ErrInvalidPassword"
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

// GetError returns error by key stored in UserErrors
func (s *UserService) GetError(key string) error {
	return UserErrors[key]
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
		return nil, err
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
		return "", UserErrors[ERR_USER_NOT_FOUND]
	}

	// Compare the provided password with the hashed password in the user object
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(loginUser.Password))
	if err != nil {
		return "", UserErrors[ERR_INVALID_PASSWORD]
	}

	cfg := config.LoadConfig()
	// Password is correct, generate the authentication token
	token, err := util.GenerateJWT(existingUser.ID, existingUser.Email, []byte(cfg.JWTSecretKey), time.Hour*24) // Adjust the secret key and token expiration time as needed
	if err != nil {
		return "", err
	}

	// Return the authenticated user and the generated token
	return token, nil
}

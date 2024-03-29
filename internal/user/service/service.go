package service

import (
	"errors"
	"github.com/DarioKnezovic/user-service/internal/session/service"
	"github.com/DarioKnezovic/user-service/internal/user"
	"github.com/DarioKnezovic/user-service/internal/user/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var UserErrors = map[string]error{
	// ErrUserNotFound is returned when a user is not found in the repository.
	"ErrUserNotFound": errors.New("user not found"),

	// ErrInvalidPassword is returned when the provided password is invalid.
	"ErrInvalidPassword": errors.New("invalid password"),

	"ErrInternalServerError": errors.New("internal server error"),
}

const (
	ERR_USER_NOT_FOUND        = "ErrUserNotFound"
	ERR_INVALID_PASSWORD      = "ErrInvalidPassword"
	ERR_INTERNAL_SERVER_ERROR = "ErrInternalServerError"
)

// UserService represents the user service implementation.
type UserService struct {
	userRepository repository.UserRepository
	sessionService *service.SessionService
}

// NewUserService creates a new instance of UserService.
func NewUserService(userRepository repository.UserRepository, sessionService *service.SessionService) *UserService {
	return &UserService{
		userRepository: userRepository,
		sessionService: sessionService,
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
func (s *UserService) LoginUser(loginUser user.User) (user.LoginResponse, error) {
	// Retrieve the user from the repository based on the email
	existingUser, err := s.userRepository.FindUserByEmail(loginUser.Email)
	if err != nil {
		return user.LoginResponse{}, err
	}

	// Check if the user exists
	if existingUser == nil {
		return user.LoginResponse{}, UserErrors[ERR_USER_NOT_FOUND]
	}

	// Compare the provided password with the hashed password in the user object
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(loginUser.Password))
	if err != nil {
		return user.LoginResponse{}, UserErrors[ERR_INVALID_PASSWORD]
	}

	token, err := s.sessionService.CreateSession(existingUser)
	if err != nil {
		return user.LoginResponse{}, UserErrors[ERR_INTERNAL_SERVER_ERROR]
	}

	// Return the authenticated user and the generated token
	return user.LoginResponse{
		User:  existingUser,
		Token: token,
	}, nil
}

func (s *UserService) LogoutUser(userId uint) error {
	err := s.sessionService.EndSession(userId)
	if err != nil {
		return UserErrors[ERR_INTERNAL_SERVER_ERROR]
	}

	return nil
}

func (s *UserService) CheckIsUserExists(userId uint) (bool, error) {
	exists, err := s.userRepository.CheckUserExists(userId)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *UserService) GetUser(userId string) (user.User, error) {
	foundedUser, err := s.userRepository.GetUserById(userId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// User not found in the database
			return user.User{}, nil
		}
	}

	return foundedUser, nil
}

func (s *UserService) UpdateUser(userId string, payload user.User) error {
	foundedUser, err := s.userRepository.GetUserById(userId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("Record not found")
		}
		return err
	}

	hashedPassword, err := s.HashPassword(payload.Password)
	if err != nil {
		return err
	}

	foundedUser.FirstName = payload.FirstName
	foundedUser.LastName = payload.LastName
	foundedUser.Email = payload.Email
	foundedUser.Password = hashedPassword

	return s.userRepository.UpdateUserById(userId, foundedUser)
}

func (s *UserService) DeleteUser(userId string) error {
	foundedUser, err := s.userRepository.GetUserById(userId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("Record not found")
		}
		return err
	}

	return s.userRepository.DeleteUserById(foundedUser)
}

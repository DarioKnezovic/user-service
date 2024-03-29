package user

import (
	"gorm.io/gorm"
	"time"
)

// User represents a user entity.
type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	FirstName   string         `gorm:"not null" json:"first_name"`
	LastName    string         `gorm:"not null" json:"last_name"`
	Email       string         `gorm:"unique;not null" json:"email"`
	CustomerKey string         `json:"customer_key"`
	Password    string         `gorm:"not null" json:"password"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type LoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// UserService represents the user service interface.
type UserService interface {
	RegisterUser(newUser User) (*User, error)
	LoginUser(loginUser User) (LoginResponse, error)
	LogoutUser(usedId uint) error
	CheckIsUserExists(userId uint) (bool, error)
	GetError(key string) error
	GetUser(userId string) (User, error)
	UpdateUser(userId string, payload User) error
	DeleteUser(userId string) error
}

package session

import (
	"github.com/DarioKnezovic/user-service/internal/user"
	"time"
)

type Session struct {
	ID        uint `gorm:"primaryKey"`
	SessionID string
	UserID    uint
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SessionService interface {
	CreateSession(user *user.User) (string, error)
	EndSession(userId uint) error
	GetSessionByToken(token string) (bool, error)
}

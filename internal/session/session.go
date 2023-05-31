package session

import "time"

type Session struct {
	ID        uint `gorm:"primaryKey"`
	SessionID string
	UserID    uint
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

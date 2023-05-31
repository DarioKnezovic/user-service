package service

import (
	"github.com/DarioKnezovic/user-service/config"
	"github.com/DarioKnezovic/user-service/internal/session"
	"github.com/DarioKnezovic/user-service/internal/session/repository"
	"github.com/DarioKnezovic/user-service/internal/user"
	"github.com/DarioKnezovic/user-service/pkg/util"
	"time"
)

type SessionService struct {
	sessionRepository repository.SessionRepository
}

func NewSessionService(sessionRepository repository.SessionRepository) *SessionService {
	return &SessionService{
		sessionRepository: sessionRepository,
	}
}

func (s *SessionService) CreateSession(user *user.User) (string, error) {
	cfg := config.LoadConfig()
	token, err := util.GenerateJWT(user.ID, user.Email, []byte(cfg.JWTSecretKey), time.Minute*30)
	if err != nil {
		return "", err
	}

	sessionData := session.Session{
		SessionID: token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Minute * 30),
	}

	sessionResponse, err := s.sessionRepository.StoreSession(sessionData)
	if err != nil {
		return "", err
	}

	return sessionResponse.SessionID, nil
}

func (s *SessionService) EndSession(userId uint) error {
	return s.sessionRepository.DeleteSession(userId)
}

package grpc

import (
	"context"
	"errors"
	"github.com/DarioKnezovic/user-service/config"
	sessionService "github.com/DarioKnezovic/user-service/internal/session/service"
	"github.com/DarioKnezovic/user-service/internal/user/service"
	"github.com/DarioKnezovic/user-service/pkg/util"
	user "github.com/DarioKnezovic/user-service/proto"
	"log"
	"time"
)

type UserServiceServer struct {
	user.UnimplementedUserServiceServer
	UserService    *service.UserService
	SessionService *sessionService.SessionService
}

func IsTokenExpired(expiredAt int64) bool {
	expirationTime := time.Unix(expiredAt, 0)
	return time.Now().UTC().After(expirationTime)
}

func (s *UserServiceServer) CheckUserExists(ctx context.Context, req *user.UserExistsRequest) (*user.UserExistsResponse, error) {
	log.Println("Check user exists received for token: ", req.Token)
	token := req.Token
	cfg := config.LoadConfig()

	claims, err := util.VerifyJWT(token, []byte(cfg.JWTSecretKey))
	if err != nil {
		log.Println(err)
		return nil, errors.New("AUTH_TOKEN_NOT_VALID")
	}

	exists, err := s.UserService.CheckIsUserExists(claims.UserID)
	if err != nil {
		log.Printf("CheckUserExists function has thrown an error: %v", err)
		return nil, err
	}

	if !exists {
		log.Printf("User with ID %d does not exists", claims.UserID)
		return &user.UserExistsResponse{Exists: false}, nil
	}

	exists, err = s.SessionService.GetSessionByToken(token)
	if err != nil {
		log.Printf("GetSessionByToken function has thrown an error: %v", err)
		return nil, err
	}

	if !exists {
		log.Printf("Session does not exists for user %d", claims.UserID)
		return &user.UserExistsResponse{Exists: false}, err
	}

	return &user.UserExistsResponse{Exists: !IsTokenExpired(claims.ExpiresAt)}, nil
}

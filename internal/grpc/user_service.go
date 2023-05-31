package grpc

import (
	"context"
	"github.com/DarioKnezovic/user-service/internal/user/service"
	user "github.com/DarioKnezovic/user-service/proto"
	"log"
)

type UserServiceServer struct {
	user.UnimplementedUserServiceServer
	UserService *service.UserService
}

func (s *UserServiceServer) CheckUserExists(ctx context.Context, req *user.UserExistsRequest) (*user.UserExistsResponse, error) {
	log.Println("Check user exists received for ID: ", req.UserId)
	userID := req.UserId
	exists, err := s.UserService.CheckIsUserExists(userID)
	if err != nil {
		log.Printf("CheckUserExists function has thrown an error: %v", err)
		return nil, err
	}

	return &user.UserExistsResponse{Exists: exists}, nil
}

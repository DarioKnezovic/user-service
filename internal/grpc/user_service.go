package grpc

import (
	"context"
	"github.com/DarioKnezovic/user-service/internal/user/repository"
	user "github.com/DarioKnezovic/user-service/proto"
	"log"
)

type UserServiceServer struct {
	user.UnimplementedUserServiceServer
	UserRepository repository.UserRepository
}

func (s *UserServiceServer) CheckUserExists(ctx context.Context, req *user.UserExistsRequest) (*user.UserExistsResponse, error) {
	log.Println("Check user exists received for ID: ", req.UserId)
	userID := req.UserId
	exists := checkUserExistsInDatabase(userID, s.UserRepository)

	return &user.UserExistsResponse{Exists: exists}, nil
}

func checkUserExistsInDatabase(userID string, repo repository.UserRepository) bool {
	exist, err := repo.CheckUserExists(userID)
	if err != nil {
		log.Println(err)
		return false
	}

	return exist
}

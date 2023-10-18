package main

import (
	"fmt"
	"github.com/DarioKnezovic/user-service/api"
	"github.com/DarioKnezovic/user-service/config"
	userGrpc "github.com/DarioKnezovic/user-service/internal/grpc"
	sessionRepo "github.com/DarioKnezovic/user-service/internal/session/repository"
	session "github.com/DarioKnezovic/user-service/internal/session/service"
	"github.com/DarioKnezovic/user-service/internal/user/repository"
	"github.com/DarioKnezovic/user-service/internal/user/service"
	"github.com/DarioKnezovic/user-service/pkg/database"
	user "github.com/DarioKnezovic/user-service/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.LoadConfig()

	// Connect to the database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// Perform auto migrations
	err = database.PerformAutoMigrations(db)
	if err != nil {
		log.Fatalf("Failed to perform auto migrations: %v", err)
	}

	// Create a user repository
	userRepo := repository.NewUserRepository(db)
	// Create a session repository
	sessionRepository := sessionRepo.NewSessionRepository(db)
	// Create a session service
	sessionService := session.NewSessionService(sessionRepository)
	// Create a user service with the repository
	userService := service.NewUserService(userRepo, sessionService)

	// Create a new Gin Gonic router
	router := gin.Default()

	// Register the API routes
	api.RegisterRoutes(router, userService)

	// Start gRPC server
	server := grpc.NewServer()
	userServiceGrpc := &userGrpc.UserServiceServer{
		UserService:    userService,
		SessionService: sessionService,
	}
	user.RegisterUserServiceServer(server, userServiceGrpc)

	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		log.Printf("gRPC server is running on port %s", cfg.GRPCPort)
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Start the server
	log.Printf("Server listening on port %s", cfg.APIPort)
	log.Fatal(router.Run(fmt.Sprintf(":%s", cfg.APIPort)))
}

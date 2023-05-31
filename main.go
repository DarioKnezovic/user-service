package main

import (
	"fmt"
	"github.com/DarioKnezovic/user-service/config"
	userGrpc "github.com/DarioKnezovic/user-service/internal/grpc"
	"github.com/DarioKnezovic/user-service/internal/user/repository"
	"github.com/DarioKnezovic/user-service/internal/user/service"
	"github.com/DarioKnezovic/user-service/pkg/database"
	user "github.com/DarioKnezovic/user-service/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"

	"github.com/DarioKnezovic/user-service/api"
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
	// Create a user service with the repository
	userService := service.NewUserService(userRepo)

	// Register the API routes
	api.RegisterRoutes(userService)

	// Start gRPC server
	server := grpc.NewServer()
	userServiceGrpc := &userGrpc.UserServiceServer{
		UserService: userService,
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
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.APIPort), nil))
}

package main

import (
	"fmt"
	"github.com/DarioKnezovic/user-service/config"
	"log"
	"net/http"

	"github.com/DarioKnezovic/user-service/api"
	"github.com/DarioKnezovic/user-service/internal/user/repository"
	"github.com/DarioKnezovic/user-service/internal/user/service"
)

func main() {
	cfg := config.LoadConfig()
	// Create a user repository
	userRepo := repository.NewUserRepository()

	// Create a user service with the repository
	userService := service.NewUserService(userRepo)

	// Register the API routes
	api.RegisterRoutes(userService)

	// Start the server
	log.Printf("Server listening on port %s", cfg.APIPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.APIPort), nil))
}

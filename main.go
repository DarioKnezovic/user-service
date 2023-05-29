package main

import (
	"fmt"
	"github.com/DarioKnezovic/user-service/config"
	"github.com/DarioKnezovic/user-service/pkg/database"
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

	_, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// Create a user service with the repository
	userService := service.NewUserService(userRepo)

	// Register the API routes
	api.RegisterRoutes(userService)

	// Start the server
	log.Printf("Server listening on port %s", cfg.APIPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.APIPort), nil))
}

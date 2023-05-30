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

	// Start the server
	log.Printf("Server listening on port %s", cfg.APIPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.APIPort), nil))
}

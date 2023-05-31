package api

import (
	"net/http"

	"github.com/DarioKnezovic/user-service/api/handlers"
	"github.com/DarioKnezovic/user-service/internal/user"
)

// RegisterRoutes registers the API routes and their corresponding handlers.
func RegisterRoutes(userService user.UserService) {
	userHandler := &handlers.UserHandler{
		UserService: userService,
	}

	http.HandleFunc("/api/user/register", userHandler.RegisterUserHandler)
	http.HandleFunc("/api/user/login", userHandler.LoginUserHandler)
	http.HandleFunc("/api/user/logout", userHandler.LogoutUserHandler)
	// Add more routes and handlers as needed
}

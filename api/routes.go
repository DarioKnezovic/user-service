package api

import (
	"github.com/DarioKnezovic/user-service/api/handlers"
	"github.com/DarioKnezovic/user-service/internal/user"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the API routes and their corresponding handlers.
func RegisterRoutes(router *gin.Engine, userService user.UserService) {
	userHandler := &handlers.UserHandler{
		UserService: userService,
	}

	router.POST("/api/user/register", userHandler.RegisterUserHandler)
	router.POST("/api/user/login", userHandler.LoginUserHandler)
	router.POST("/api/user/logout", userHandler.LogoutUserHandler)

	router.GET("/api/users/:id", userHandler.GetUserDetailsHandler)
	router.PUT("/api/users/:id", userHandler.UpdateUserHandler)
}

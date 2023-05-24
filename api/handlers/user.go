package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/DarioKnezovic/user-service/internal/user"
	"net/http"
)

type UserHandler struct {
	UserService user.UserService
}

// RegisterUserHandler handles the registration of a new user.
func (h *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser user.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	fmt.Println("WOOOHOOOO | We are inside Register")

	registeredUser, err := h.UserService.RegisterUser(newUser)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registeredUser)
}

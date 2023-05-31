package handlers

import (
	"encoding/json"
	"errors"
	"github.com/DarioKnezovic/user-service/internal/user"
	"github.com/DarioKnezovic/user-service/pkg/util"
	"log"
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
		// TODO: add response body
		util.SendJSONResponse(w, http.StatusBadRequest, nil)
		return
	}

	registeredUser, err := h.UserService.RegisterUser(newUser)
	if err != nil {
		// TODO: add response body
		util.SendJSONResponse(w, http.StatusInternalServerError, nil)
		return
	}

	util.SendJSONResponse(w, http.StatusOK, registeredUser)
}

func (h *UserHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var loginUser user.User
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		util.SendJSONResponse(w, http.StatusBadRequest, nil)
		log.Println(err)
		return
	}

	token, err := h.UserService.LoginUser(loginUser)
	if err != nil {
		var statusCode int
		var responseBody interface{}

		switch {
		case errors.Is(err, h.UserService.GetError("ErrUserNotFound")):
			statusCode = http.StatusNotFound
			responseBody = map[string]string{
				"error": "User not found",
			}
		case errors.Is(err, h.UserService.GetError("ErrInvalidPassword")):
			statusCode = http.StatusUnauthorized
			responseBody = map[string]string{
				"error": "Invalid password",
			}
		default:
			statusCode = http.StatusInternalServerError
			responseBody = map[string]string{
				"error": "Internal server error",
			}
			log.Println(err)
		}

		util.SendJSONResponse(w, statusCode, responseBody)
		return
	}

	responseBody := map[string]string{
		"token": token,
	}

	util.SendJSONResponse(w, http.StatusOK, responseBody)
}

func (h *UserHandler) LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	var loggedUser user.User
	err := json.NewDecoder(r.Body).Decode(&loggedUser)
	if err != nil {
		log.Println(err)
		util.SendJSONResponse(w, http.StatusBadRequest, nil)
		return
	}

	err = h.UserService.LogoutUser(loggedUser.ID)
	if err != nil {
		log.Println(err)
		util.SendJSONResponse(w, http.StatusInternalServerError, nil)
		return
	}

	util.SendJSONResponse(w, http.StatusOK, nil)
}

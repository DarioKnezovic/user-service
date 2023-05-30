package handlers

import (
	"encoding/json"
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
		// TODO: add response body
		util.SendJSONResponse(w, http.StatusBadRequest, nil)
		log.Panic(err)
		return
	}

	token, err := h.UserService.LoginUser(loginUser)
	if err != nil {
		// TODO: add response body
		util.SendJSONResponse(w, http.StatusNotFound, nil)
		log.Panic(err)
		return
	}
	responseBody := map[string]string{
		"token": token,
	}

	util.SendJSONResponse(w, http.StatusOK, responseBody)
}

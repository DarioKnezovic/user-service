package handlers

import (
	"errors"
	"github.com/DarioKnezovic/user-service/internal/user"
	"github.com/DarioKnezovic/user-service/pkg/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type UserHandler struct {
	UserService user.UserService
}

// RegisterUserHandler handles the registration of a new user.
func (h *UserHandler) RegisterUserHandler(c *gin.Context) {
	var newUser user.User

	if err := c.BindJSON(&newUser); err != nil {
		// Handle JSON decoding error
		util.SendJSONResponse(c, http.StatusBadRequest, nil)
		return
	}

	registeredUser, err := h.UserService.RegisterUser(newUser)
	if err != nil {
		// TODO: add response body
		util.SendJSONResponse(c, http.StatusInternalServerError, nil)
		return
	}

	util.SendJSONResponse(c, http.StatusOK, registeredUser)
}

func (h *UserHandler) LoginUserHandler(c *gin.Context) {
	var loginUser user.User
	if err := c.BindJSON(&loginUser); err != nil {
		// Handle JSON decoding error
		util.SendJSONResponse(c, http.StatusBadRequest, nil)
		return
	}

	response, err := h.UserService.LoginUser(loginUser)
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

		util.SendJSONResponse(c, statusCode, responseBody)
		return
	}

	util.SendJSONResponse(c, http.StatusOK, response)
}

func (h *UserHandler) LogoutUserHandler(c *gin.Context) {
	var loggedUser user.User
	if err := c.BindJSON(&loggedUser); err != nil {
		// Handle JSON decoding error
		util.SendJSONResponse(c, http.StatusBadRequest, nil)
		return
	}

	err := h.UserService.LogoutUser(loggedUser.ID)
	if err != nil {
		log.Println(err)
		util.SendJSONResponse(c, http.StatusInternalServerError, nil)
		return
	}

	util.SendJSONResponse(c, http.StatusOK, nil)
}

func (h *UserHandler) GetUserDetailsHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.SendJSONResponse(c, http.StatusBadRequest, nil)
		return
	}

	fetchedUser, err := h.UserService.GetUser(id)
	if err != nil {
		log.Println(err)
		util.SendJSONResponse(c, http.StatusInternalServerError, nil)
		return
	}

	util.SendJSONResponse(c, http.StatusOK, fetchedUser)
}

func (h *UserHandler) UpdateUserHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.SendJSONResponse(c, http.StatusBadRequest, nil)
		return
	}

	var userPayload user.User
	if err := c.BindJSON(&userPayload); err != nil {
		// Handle JSON decoding error
		util.SendJSONResponse(c, http.StatusBadRequest, nil)
		return
	}

	err := h.UserService.UpdateUser(id, userPayload)
	if err != nil {
		if err.Error() == "Record not found" {
			util.SendJSONResponse(c, http.StatusNotFound, nil)
			return
		}

		util.SendJSONResponse(c, http.StatusInternalServerError, nil)
		return
	}

	util.SendJSONResponse(c, http.StatusOK, []interface{}{})
}

func (h *UserHandler) DeleteUserHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.SendJSONResponse(c, http.StatusBadRequest, nil)
		return
	}

	err := h.UserService.DeleteUser(id)
	if err != nil {
		if err.Error() == "Record not found" {
			util.SendJSONResponse(c, http.StatusNotFound, nil)
			return
		}

		util.SendJSONResponse(c, http.StatusInternalServerError, nil)
		return
	}
}

package util

import (
	"github.com/gin-gonic/gin"
)

// SendJSONResponse sends an HTTP response with the given status code and response body using Gin Gonic's context.
func SendJSONResponse(c *gin.Context, statusCode int, responseBody interface{}) {
	c.Header("Content-Type", "application/json")
	c.JSON(statusCode, responseBody)
}

package util

import (
	"encoding/json"
	"net/http"
)

// SendJSONResponse sends an HTTP response with the given status code and response body.
func SendJSONResponse(w http.ResponseWriter, statusCode int, responseBody interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if responseBody != nil {
		json.NewEncoder(w).Encode(responseBody)
	}
}

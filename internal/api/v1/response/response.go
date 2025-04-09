package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool           `json:"success"`
	Data    interface{}    `json:"data,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

// JSON sends a JSON response
func JSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		response := Response{
			Success: true,
			Data:    data,
		}

		json.NewEncoder(w).Encode(response)
	}
}

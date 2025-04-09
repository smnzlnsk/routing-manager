package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/smnzlnsk/routing-manager/internal/domain"
)

type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Error sends an error response
func Error(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errResp := &ErrorResponse{
		Code:    http.StatusText(status),
		Message: err.Error(),
	}

	// Map domain errors to appropriate status codes
	var domainErr *domain.Error
	if errors.As(err, &domainErr) {
		switch domainErr.Code {
		case domain.CodeNotFound:
			status = http.StatusNotFound
			errResp.Code = "not_found"
		case domain.CodeInterestAlreadyExists:
			status = http.StatusConflict
			errResp.Code = "interest_already_exists"
			// Add other domain error mappings
		}
	}

	response := Response{
		Success: false,
		Error:   errResp,
	}

	json.NewEncoder(w).Encode(response)
}

package httpErrors

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// purelyHttpError defines a custom error type
type HttpError struct {
	Code       string
	StatusCode int
	Message    string
}

// Error implements the error interface for purelyHttpError
func (e *HttpError) Error() string {
	return fmt.Sprintf("Code: %s, StatusCode: %d, Message: %s", e.Code, e.StatusCode, e.Message)
}

// NewPurelyHttpError creates a new purelyHttpError with default values
func HydrateHttpError(code string, statusCode int, message string) *HttpError {
	// Set default values
	if code == "" {
		code = "purely/requests/errors/internal"
	}
	if statusCode == 0 {
		statusCode = fiber.StatusInternalServerError // Default to 500
	}
	if message == "" {
		message = "Something went wrong"
	}

	return &HttpError{
		Code:       code,
		StatusCode: statusCode,
		Message:    message,
	}
}

package errors

import (
	"errors"
	"net/http"
)

// ErrorResponseDTO is the structure returned to clients on error
type ErrorResponseDTO struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// ToErrorResponse converts an error to a user-facing error response
// It hides technical details for internal/database errors (security principle)
func ToErrorResponse(err error) ErrorResponseDTO {
	kind := GetKind(err)

	switch kind {
	case ErrKindValidation:
		return ErrorResponseDTO{
			Error: "validation error",
			Code:  string(kind),
		}
	case ErrKindNotFound:
		return ErrorResponseDTO{
			Error: "not found",
			Code:  string(kind),
		}
	case ErrKindConflict:
		return ErrorResponseDTO{
			Error: "conflict",
			Code:  string(kind),
		}
	case ErrKindDatabase, ErrKindInternal, ErrKindConfig:
		// Don't expose technical details to client
		return ErrorResponseDTO{
			Error: "internal server error",
			Code:  string(kind),
		}
	default:
		return ErrorResponseDTO{
			Error: "unknown error",
			Code:  "unknown",
		}
	}
}

// StatusCodeForError returns the HTTP status code for an error
func StatusCodeForError(err error) int {
	var appErr *AppError
	if IsAppError(err, &appErr) {
		return appErr.HTTPStatusCode()
	}
	// Default to 500 for unknown errors
	return http.StatusInternalServerError
}

// IsAppError checks if err is an AppError and returns true if it is
// If it is, it also sets appErr to the unwrapped AppError
func IsAppError(err error, appErr **AppError) bool {
	var ae *AppError
	if err != nil {
		ok := As(err, &ae)
		if ok && appErr != nil {
			*appErr = ae
		}
		return ok
	}
	return false
}

// As unwraps an error using errors.As
// This is a convenience wrapper to check for *AppError in the error chain
func As(err error, target **AppError) bool {
	return errors.As(err, target)
}

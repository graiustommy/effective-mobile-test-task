package errors

import (
	"errors"
	"net/http"
)

type ErrorResponseDTO struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

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

func StatusCodeForError(err error) int {
	var appErr *AppError
	if IsAppError(err, &appErr) {
		return appErr.HTTPStatusCode()
	}
	return http.StatusInternalServerError
}

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

func As(err error, target **AppError) bool {
	return errors.As(err, target)
}

package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorKind string

const (
	ErrKindValidation ErrorKind = "validation"
	ErrKindNotFound   ErrorKind = "not_found"
	ErrKindConflict   ErrorKind = "conflict"
	ErrKindDatabase   ErrorKind = "database"
	ErrKindInternal   ErrorKind = "internal"
	ErrKindConfig     ErrorKind = "config"
)

type AppError struct {
	kind    ErrorKind
	message string
	err     error
}

func New(kind ErrorKind, message string) *AppError {
	return &AppError{
		kind:    kind,
		message: message,
		err:     errors.New(message),
	}
}

func Wrap(kind ErrorKind, message string, err error) *AppError {
	if err == nil {
		return &AppError{
			kind:    kind,
			message: message,
			err:     errors.New(message),
		}
	}
	return &AppError{
		kind:    kind,
		message: message,
		err:     fmt.Errorf("%s: %w", message, err),
	}
}

func (e *AppError) Kind() ErrorKind {
	return e.kind
}

func (e *AppError) Message() string {
	return e.message
}

func (e *AppError) Error() string {
	return e.err.Error()
}

func (e *AppError) Unwrap() error {
	return e.err
}

func (e *AppError) HTTPStatusCode() int {
	switch e.kind {
	case ErrKindValidation:
		return http.StatusBadRequest
	case ErrKindNotFound:
		return http.StatusNotFound
	case ErrKindConflict:
		return http.StatusConflict
	case ErrKindDatabase, ErrKindInternal, ErrKindConfig:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

var (
	ErrInvalidUserID         = New(ErrKindValidation, "invalid user id")
	ErrInvalidSubscriptionID = New(ErrKindValidation, "invalid subscription id")
	ErrServiceNameRequired   = New(ErrKindValidation, "service name is required")
	ErrInvalidDateFormat     = New(ErrKindValidation, "invalid date format, expected MM-YYYY")
	ErrSubscriptionNotFound  = New(ErrKindNotFound, "subscription not found")
	ErrSubscriptionExists    = New(ErrKindConflict, "subscription already exists")
)

func IsKind(err error, kind ErrorKind) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Kind() == kind
	}
	return false
}

func GetKind(err error) ErrorKind {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Kind()
	}
	return ErrKindInternal
}

func ValidateDateFormat(date string) error {
	if date == "" {
		return nil
	}
	if len(date) != 7 {
		return ErrInvalidDateFormat
	}
	if date[2] != '-' {
		return ErrInvalidDateFormat
	}
	month := date[0:2]
	year := date[3:7]

	if month < "01" || month > "12" {
		return ErrInvalidDateFormat
	}

	for _, c := range year {
		if c < '0' || c > '9' {
			return ErrInvalidDateFormat
		}
	}

	return nil
}

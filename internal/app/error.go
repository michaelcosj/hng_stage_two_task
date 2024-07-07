package app

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
)

var (
	ErrUserAlreadyExists    = errors.New("User already exists")
	ErrAuthenticationFailed = errors.New("Authentication Failed")
	ErrUserNotFound         = errors.New("User does not exist")
	ErrOrgNotFound          = errors.New("Organisation does not exist")
	ErrClientError          = errors.New("Client error")
)

type validationErrorItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ApiValidationError struct {
	Errors []validationErrorItem `json:"errors"`
}

func NewValidationError(problems map[string]string) ApiValidationError {
	err := ApiValidationError{}
	for field, message := range problems {
		err.Errors = append(err.Errors, validationErrorItem{field, message})
	}
	return err
}

func (e ApiValidationError) Error() string {
	buf := new(bytes.Buffer)
	for _, item := range e.Errors {
		fmt.Fprintf(buf, "{%s:%s}\t", item.Field, item.Message)
	}
	fmt.Fprintf(buf, "\n")
	return buf.String()
}

type ApiError struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	wrappedErr error
}

func (e ApiError) Error() string {
	return fmt.Sprintf("{status: %s, message: %s, code: %d, error: %s}\n", e.Status, e.Message, e.StatusCode, e.wrappedErr)
}

func (e ApiError) Is(target error) bool {
	_, ok := target.(ApiError)
	return ok
}

func NewApiError(status, message string, statusCode int) ApiError {
	return ApiError{status, message, statusCode, nil}
}

func InvalidJson() ApiError {
	return ApiError{
		Status:     "Bad Request",
		Message:    "Invalid json data",
		StatusCode: http.StatusBadRequest,
	}
}

func InvalidRequestData(err error) ApiError {
	return ApiError{
		Status:     "Bad Request",
		Message:    "Client Error",
		StatusCode: http.StatusBadRequest,
		wrappedErr: err,
	}
}

func ApiErrorFrom(err error) error {

	switch {
	case errors.Is(ApiError{}, err):
		return err
	case errors.Is(err, ErrUserAlreadyExists):
		return ApiError{
			Status:     "User already exists",
			Message:    "A user with this email already exists",
			StatusCode: http.StatusUnprocessableEntity,
			wrappedErr: err,
		}
	case errors.Is(err, ErrAuthenticationFailed):
		return ApiError{
			Status:     "Bad request",
			Message:    "Authentication failed",
			StatusCode: http.StatusUnauthorized,
			wrappedErr: err,
		}
	case errors.Is(err, ErrUserNotFound):
		return ApiError{
			Status:     "Not Found",
			Message:    "User not found",
			StatusCode: http.StatusNotFound,
			wrappedErr: err,
		}
	case errors.Is(err, ErrOrgNotFound):
		return ApiError{
			Status:     "Not found",
			Message:    "Organisation not found",
			StatusCode: http.StatusNotFound,
			wrappedErr: err,
		}
	case errors.Is(err, ErrClientError):
		return ApiError{
			Status:     "Bad request",
			Message:    "Client error",
			StatusCode: http.StatusBadRequest,
			wrappedErr: err,
		}
	default:
		log.Printf("unknown error encountered: %v", err)
		return err
	}
}

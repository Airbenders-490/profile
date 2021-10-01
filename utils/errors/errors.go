package errors

import "net/http"

type RestError struct {
	Code int
	Message string
}

func (e RestError) Error() string {
	return e.Message
}

func NewInternalServerError(message string) *RestError {
	return &RestError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}

func NewConflictError(message string) *RestError {
	return &RestError{
		Code:    http.StatusConflict,
		Message: message,
	}
}

func NewNotFoundError(message string) *RestError {
	return &RestError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

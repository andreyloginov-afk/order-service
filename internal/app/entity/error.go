package entity

import "net/http"

type AppError struct {
	code    int
	message string
}

func NewAppError(code int, message string) *AppError {
	return &AppError{code: code, message: message}
}

func (e *AppError) Error() string {
	return e.message
}

func (e *AppError) HTTPStatus() int {
	return e.code
}

var (
	ErrNotFound            = NewAppError(http.StatusNotFound, "not found")
	ErrAlreadyExists       = NewAppError(http.StatusConflict, "already exists")
	ErrIncorrectParameters = NewAppError(http.StatusBadRequest, "incorrect parameters")
)

package errors

import (
	"errors"
	"net/http"
)

type CustomError struct {
	Err           error // private message
	Status        int
	PublicMessage string
}

func (e CustomError) Error() string {
	return e.Err.Error()
}

var _ error = (*CustomError)(nil)

func New(status int, privateMessage string, publicMessage string) CustomError {
	return CustomError{
		Err:           errors.New(privateMessage),
		PublicMessage: publicMessage,
		Status:        status,
	}
}

func NewInternalServerError(err error) CustomError {
	return CustomError{
		Err:           err,
		PublicMessage: "Internal server error",
		Status:        http.StatusInternalServerError,
	}
}

func NewBadRequestError(err error, publicMessage string) CustomError {
	return CustomError{
		Err:           err,
		PublicMessage: publicMessage,
		Status:        http.StatusBadRequest,
	}
}

func NewForbiddenError(err error) CustomError {
	return CustomError{
		Err:           err,
		PublicMessage: "User is not authorized to access the requested resource",
		Status:        http.StatusForbidden,
	}
}

func NewUnauthorizedError(err error) CustomError {
	return CustomError{
		Err:           err,
		PublicMessage: "User is not authorized to access the requested resource",
		Status:        http.StatusUnauthorized,
	}
}

func NewNotFoundError(err error) CustomError {
	return CustomError{
		Err:           err,
		PublicMessage: "Not found",
		Status:        http.StatusNotFound,
	}
}

// predefined errors
var (
	ErrUnauthorized       CustomError = New(http.StatusUnauthorized, "unathurized", "Authentication Required - Please provide valid credentials to access this resource")
	ErrForbidden          CustomError = NewForbiddenError(errors.New("forbidden"))
	ErrMissingBearerToken CustomError = New(http.StatusUnauthorized, "missing bearer token", "Bearer token is required to access the requested resource")
)

package errors

import (
	"errors"
)

var (
	ErrWrongPWOrEmail        = errors.New("wrong email or password")
	ErrGenericMessage        = errors.New("an unexpected error has occurred, please try again later")
	ErrFailedToRetrieveToken = errors.New("something went wrong during retrieving user id from token")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden = errors.New("forbidden")

	ErrInvalidToken = errors.New("invalid token")
)
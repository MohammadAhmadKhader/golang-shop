package errors

import (
	"errors"
)

var (
	ErrWrongPWOrEmail        = errors.New("wrong email or password")
	ErrPasswordsNotMatching  = errors.New("passwords not matching")
	ErrGenericMessage        = errors.New("an unexpected error has occurred, please try again later")
	ErrFailedToRetrieveToken = errors.New("something went wrong during retrieving user id from token")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden = errors.New("forbidden")

	ErrInvalidToken = errors.New("invalid token")

	ErrNoFileFound = errors.New("no file was found")
	ErrUnexpectedDuringImageUpload = errors.New("an error has occurred during uploading image")
)
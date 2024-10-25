package errors

import "fmt"

type InvalidIdError struct {
	Resource string
	ID       uint
}

func (e *InvalidIdError) Error() string {
	return fmt.Sprintf("invalid %s id received '%d' ", e.Resource, e.ID)
}

func NewInvalidIDError(resource string, id uint) error {
	return &InvalidIdError{
		Resource: resource,
		ID:       id,
	}
}

// this error its message must contain a placeholder for its id
type ResourceWasNotFoundError struct {
	NotFoundErrMessage string
	ID              uint
}

func (e *ResourceWasNotFoundError) Error() string {
	return fmt.Sprintf(e.NotFoundErrMessage, e.ID)
}

func NewResourceWasNotFoundError(resource string, id uint) error {
	return &ResourceWasNotFoundError{
		NotFoundErrMessage: resource,
		ID:              id,
	}
}
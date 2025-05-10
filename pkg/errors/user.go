package errors

import "fmt"

type UserNotFoundError struct {
	Message string
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("UserNotFound error: %s", e.Message)
}

type UserAlreadyExistsError struct {
	Message string
}

func (e UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("UserAlreadyExists error: %s", e.Message)
}

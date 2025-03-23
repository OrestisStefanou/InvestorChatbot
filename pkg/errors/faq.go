package errors

import "fmt"

type TopicNotFoundError struct {
	Message string
}

func (e *TopicNotFoundError) Error() string {
	return fmt.Sprintf("TopicNotFoundError error: %s", e.Message)
}

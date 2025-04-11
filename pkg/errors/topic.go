package errors

import "fmt"

type InvalidTopicError struct {
	Message string
}

func (e InvalidTopicError) Error() string {
	return fmt.Sprintf("InvalidTopic error: %s", e.Message)
}

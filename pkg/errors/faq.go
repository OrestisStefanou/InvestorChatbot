package errors

import "fmt"

type TopicNotFoundError struct {
	Message string
}

func (e TopicNotFoundError) Error() string {
	return fmt.Sprintf("TopicNotFoundError error: %s", e.Message)
}

type FaqTopicNotFoundError struct {
	Message string
}

func (e FaqTopicNotFoundError) Error() string {
	return fmt.Sprintf("FaqTopicNotFoundError error: %s", e.Message)
}

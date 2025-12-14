package services

import "fmt"

type DataServiceError struct {
	Message string
}

func (e *DataServiceError) Error() string {
	return fmt.Sprintf("Data Service error: %s", e.Message)
}

type RagError struct {
	Message string
}

func (e *RagError) Error() string {
	return fmt.Sprintf("Rag error: %s", e.Message)
}

type SessionServiceError struct {
	Message string
}

func (e *SessionServiceError) Error() string {
	return fmt.Sprintf("Rag error: %s", e.Message)
}

package main

import "time"

// TaskStatus represents the possible states of a task
type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in-progress"
	StatusDone       TaskStatus = "done"
)

// Task represents a single task with all its properties
type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// Domain Errors
type TaskError struct {
	Code    string
	Message string
}

var (
	ErrTaskNotFound     = TaskError{Code: "NOT_FOUND", Message: "Task not found"}
	ErrInvalidStatus    = TaskError{Code: "INVALID_STATUS", Message: "Invalid task status"}
	ErrEmptyDescription = TaskError{
		Code:    "EMPTY_DESCRIPTION",
		Message: "Task description cannot be empty",
	}
	ErrInvalidID = TaskError{Code: "INVALID_ID", Message: "Invalid task ID"}
)

func (e TaskError) Error() string {
	return e.Message
}

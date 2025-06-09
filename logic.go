package main

import (
	"strings"
	"time"
)

// NewTask creates a new task with validation
func NewTask(id int, description string) (*Task, error) {
	if strings.TrimSpace(description) == "" {
		return nil, ErrEmptyDescription
	}

	now := time.Now()
	return &Task{
		ID:          id,
		Description: strings.TrimSpace(description),
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// UpdateDescription updates the task description
func (t *Task) UpdateDescription(description string) error {
	if strings.TrimSpace(description) == "" {
		return ErrEmptyDescription
	}

	t.Description = strings.TrimSpace(description)
	t.UpdatedAt = time.Now()
	return nil
}

// MarkInProgress changes task status to in-progress
func (t *Task) MarkInProgress() {
	t.Status = StatusInProgress
	t.UpdatedAt = time.Now()
}

// MarkDone changes task status to done
func (t *Task) MarkDone() {
	t.Status = StatusDone
	t.UpdatedAt = time.Now()
}

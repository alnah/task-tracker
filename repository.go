package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Repository Interface (Port)
type TaskRepository interface {
	Save(tasks []Task) error
	Load() ([]Task, error)
	GetNextID() (int, error)
}

// File Repository Implementation (Adapter)
type FileTaskRepository struct {
	filename string
}

func NewFileTaskRepository(filename string) *FileTaskRepository {
	return &FileTaskRepository{filename: filename}
}

func (r *FileTaskRepository) Save(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	err = os.WriteFile(r.filename, data, 0o600)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (r *FileTaskRepository) Load() ([]Task, error) {
	// Check if file exists
	if _, err := os.Stat(r.filename); os.IsNotExist(err) {
		// Return empty slice if file doesn't exist
		return []Task{}, nil
	}

	data, err := os.ReadFile(r.filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Handle empty file
	if len(data) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
	}

	return tasks, nil
}

func (r *FileTaskRepository) GetNextID() (int, error) {
	tasks, err := r.Load()
	if err != nil {
		return 0, err
	}

	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	return maxID + 1, nil
}

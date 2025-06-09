package main

import (
	"fmt"
	"slices"
)

// Application Service (Use Cases)
type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) AddTask(description string) (*Task, error) {
	nextID, err := s.repo.GetNextID()
	if err != nil {
		return nil, fmt.Errorf("failed to get next ID: %w", err)
	}

	task, err := NewTask(nextID, description)
	if err != nil {
		return nil, err
	}

	tasks, err := s.repo.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}

	tasks = append(tasks, *task)

	err = s.repo.Save(tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to save tasks: %w", err)
	}

	return task, nil
}

func (s *TaskService) UpdateTask(id int, description string) error {
	tasks, err := s.repo.Load()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	taskIndex := -1
	for i, task := range tasks {
		if task.ID == id {
			taskIndex = i
			break
		}
	}

	if taskIndex == -1 {
		return ErrTaskNotFound
	}

	err = tasks[taskIndex].UpdateDescription(description)
	if err != nil {
		return err
	}

	return s.repo.Save(tasks)
}

func (s *TaskService) DeleteTask(id int) error {
	tasks, err := s.repo.Load()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	taskIndex := -1
	for i, task := range tasks {
		if task.ID == id {
			taskIndex = i
			break
		}
	}

	if taskIndex == -1 {
		return ErrTaskNotFound
	}

	// Remove task from slice
	tasks = slices.Delete(tasks, taskIndex, taskIndex+1)

	return s.repo.Save(tasks)
}

func (s *TaskService) MarkTaskInProgress(id int) error {
	return s.updateTaskStatus(id, func(task *Task) {
		task.MarkInProgress()
	})
}

func (s *TaskService) MarkTaskDone(id int) error {
	return s.updateTaskStatus(id, func(task *Task) {
		task.MarkDone()
	})
}

func (s *TaskService) updateTaskStatus(id int, updateFn func(*Task)) error {
	tasks, err := s.repo.Load()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	taskIndex := -1
	for i, task := range tasks {
		if task.ID == id {
			taskIndex = i
			break
		}
	}

	if taskIndex == -1 {
		return ErrTaskNotFound
	}

	updateFn(&tasks[taskIndex])

	return s.repo.Save(tasks)
}

func (s *TaskService) ListTasks(status string) ([]Task, error) {
	tasks, err := s.repo.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}

	if status == "" {
		return tasks, nil
	}

	var filteredTasks []Task
	for _, task := range tasks {
		if string(task.Status) == status {
			filteredTasks = append(filteredTasks, task)
		}
	}

	return filteredTasks, nil
}

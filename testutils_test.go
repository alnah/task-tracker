package main

import (
	"errors"
	"testing"
	"time"
)

// Test Utilities for Task Tracker
// All utilities are in the same package - no import issues!

// TaskBuilder provides a fluent interface for creating test tasks
type TaskBuilder struct {
	id          int
	description string
	status      TaskStatus
	createdAt   time.Time
	updatedAt   time.Time
}

// NewTaskBuilder creates a new task builder with sensible defaults
func NewTaskBuilder() *TaskBuilder {
	now := time.Now()
	return &TaskBuilder{
		id:          1,
		description: "Default test task",
		status:      StatusTodo,
		createdAt:   now,
		updatedAt:   now,
	}
}

// WithID sets the task ID
func (b *TaskBuilder) WithID(id int) *TaskBuilder {
	b.id = id
	return b
}

// WithDescription sets the task description
func (b *TaskBuilder) WithDescription(desc string) *TaskBuilder {
	b.description = desc
	return b
}

// WithStatus sets the task status
func (b *TaskBuilder) WithStatus(status TaskStatus) *TaskBuilder {
	b.status = status
	return b
}

// WithTimestamps sets both created and updated timestamps
func (b *TaskBuilder) WithTimestamps(created, updated time.Time) *TaskBuilder {
	b.createdAt = created
	b.updatedAt = updated
	return b
}

// InProgress is a convenience method to set status to in-progress
func (b *TaskBuilder) InProgress() *TaskBuilder {
	return b.WithStatus(StatusInProgress)
}

// Done is a convenience method to set status to done
func (b *TaskBuilder) Done() *TaskBuilder {
	return b.WithStatus(StatusDone)
}

// BuildValid creates a valid Task using the domain constructor
func (b *TaskBuilder) BuildValid(t *testing.T) *Task {
	t.Helper()

	// Use the actual NewTask function from main package - no imports needed!
	task, err := NewTask(b.id, b.description)
	if err != nil {
		t.Fatalf("Failed to create valid test task: %v", err)
	}

	// Set custom timestamps if provided
	if !b.createdAt.IsZero() {
		task.CreatedAt = b.createdAt
	}
	if !b.updatedAt.IsZero() {
		task.UpdatedAt = b.updatedAt
	}

	// Apply status changes through domain methods
	switch b.status {
	case StatusInProgress:
		task.MarkInProgress()
	case StatusDone:
		task.MarkDone()
	}

	return task
}

// BuildInvalid creates a Task struct bypassing domain validation (for testing edge cases)
func (b *TaskBuilder) BuildInvalid() *Task {
	return &Task{
		ID:          b.id,
		Description: b.description,
		Status:      b.status,
		CreatedAt:   b.createdAt,
		UpdatedAt:   b.updatedAt,
	}
}

// Common test fixtures

// TodoTask creates a basic todo task
func TodoTask(t *testing.T) *Task {
	t.Helper()
	return NewTaskBuilder().WithDescription("Buy groceries").BuildValid(t)
}

// InProgressTask creates a task in progress
func InProgressTask(t *testing.T) *Task {
	t.Helper()
	return NewTaskBuilder().
		WithDescription("Write report").
		InProgress().
		BuildValid(t)
}

// DoneTask creates a completed task
func DoneTask(t *testing.T) *Task {
	t.Helper()
	return NewTaskBuilder().
		WithDescription("Call mom").
		Done().
		BuildValid(t)
}

// TaskWithID creates a task with specific ID
func TaskWithID(t *testing.T, id int) *Task {
	t.Helper()
	return NewTaskBuilder().
		WithID(id).
		WithDescription("Task " + string(rune(id+48))). // Convert to ASCII
		BuildValid(t)
}

// TaskSet creates multiple tasks for testing
func TaskSet(t *testing.T, count int) []Task {
	t.Helper()
	tasks := make([]Task, count)
	for i := range count {
		task := NewTaskBuilder().
			WithID(i + 1).
			WithDescription("Task " + string(rune(i+49))). // '1', '2', '3', etc.
			BuildValid(t)
		tasks[i] = *task
	}
	return tasks
}

// MixedStatusTasks creates a set of tasks with different statuses
func MixedStatusTasks(t *testing.T) []Task {
	t.Helper()
	return []Task{
		*NewTaskBuilder().WithID(1).WithDescription("Todo task").BuildValid(t),
		*NewTaskBuilder().WithID(2).WithDescription("In progress task").InProgress().BuildValid(t),
		*NewTaskBuilder().WithID(3).WithDescription("Done task").Done().BuildValid(t),
	}
}

// Time helpers for testing timestamps

// FixedTime returns a fixed time for consistent testing
func FixedTime() time.Time {
	return time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
}

// TimeAfter returns a time that's definitely after the given time
func TimeAfter(base time.Time) time.Time {
	return base.Add(1 * time.Minute)
}

// TimeBefore returns a time that's definitely before the given time
func TimeBefore(base time.Time) time.Time {
	return base.Add(-1 * time.Minute)
}

// Assertion helpers

// AssertTaskEquals compares two tasks for equality in tests
func AssertTaskEquals(t *testing.T, expected, actual *Task) {
	t.Helper()

	if actual.ID != expected.ID {
		t.Errorf("ID mismatch: got %d, want %d", actual.ID, expected.ID)
	}
	if actual.Description != expected.Description {
		t.Errorf("Description mismatch: got %q, want %q", actual.Description, expected.Description)
	}
	if actual.Status != expected.Status {
		t.Errorf("Status mismatch: got %q, want %q", actual.Status, expected.Status)
	}
	if !actual.CreatedAt.Equal(expected.CreatedAt) {
		t.Errorf("CreatedAt mismatch: got %v, want %v", actual.CreatedAt, expected.CreatedAt)
	}
	if !actual.UpdatedAt.Equal(expected.UpdatedAt) {
		t.Errorf("UpdatedAt mismatch: got %v, want %v", actual.UpdatedAt, expected.UpdatedAt)
	}
}

// AssertTasksEqual compares slices of tasks
func AssertTasksEqual(t *testing.T, expected, actual []Task) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Fatalf("Task count mismatch: got %d, want %d", len(actual), len(expected))
	}

	for i := range expected {
		AssertTaskEquals(t, &expected[i], &actual[i])
	}
}

// AssertTaskInSlice verifies a task exists in a slice
func AssertTaskInSlice(t *testing.T, task *Task, slice []Task) {
	t.Helper()

	for _, sliceTask := range slice {
		if sliceTask.ID == task.ID {
			AssertTaskEquals(t, task, &sliceTask)
			return
		}
	}

	t.Errorf("Task with ID %d not found in slice", task.ID)
}

// AssertTaskNotInSlice verifies a task does not exist in a slice
func AssertTaskNotInSlice(t *testing.T, taskID int, slice []Task) {
	t.Helper()

	for _, sliceTask := range slice {
		if sliceTask.ID == taskID {
			t.Errorf("Task with ID %d should not be in slice", taskID)
			return
		}
	}
}

// MockTaskRepository is an in-memory implementation for testing
type MockTaskRepository struct {
	tasks         []Task
	shouldError   bool
	errorToReturn error
	saveCallCount int
	loadCallCount int
}

// NewMockRepository creates a new mock repository
func NewMockRepository() *MockTaskRepository {
	return &MockTaskRepository{
		tasks: make([]Task, 0),
	}
}

// WithTasks preloads the repository with tasks
func (m *MockTaskRepository) WithTasks(tasks []Task) *MockTaskRepository {
	m.tasks = make([]Task, len(tasks))
	copy(m.tasks, tasks)
	return m
}

// WithError configures the repository to return errors
func (m *MockTaskRepository) WithError(err error) *MockTaskRepository {
	m.shouldError = true
	m.errorToReturn = err
	return m
}

// Save implements TaskRepository interface
func (m *MockTaskRepository) Save(tasks []Task) error {
	m.saveCallCount++

	if m.shouldError {
		return m.errorToReturn
	}

	m.tasks = make([]Task, len(tasks))
	copy(m.tasks, tasks)
	return nil
}

// Load implements TaskRepository interface
func (m *MockTaskRepository) Load() ([]Task, error) {
	m.loadCallCount++

	if m.shouldError {
		return nil, m.errorToReturn
	}

	result := make([]Task, len(m.tasks))
	copy(result, m.tasks)
	return result, nil
}

// GetNextID implements TaskRepository interface
func (m *MockTaskRepository) GetNextID() (int, error) {
	if m.shouldError {
		return 0, m.errorToReturn
	}

	maxID := 0
	for _, task := range m.tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	return maxID + 1, nil
}

// Test helpers for verification

// SaveCallCount returns how many times Save was called
func (m *MockTaskRepository) SaveCallCount() int {
	return m.saveCallCount
}

// LoadCallCount returns how many times Load was called
func (m *MockTaskRepository) LoadCallCount() int {
	return m.loadCallCount
}

// GetStoredTasks returns a copy of the currently stored tasks
func (m *MockTaskRepository) GetStoredTasks() []Task {
	result := make([]Task, len(m.tasks))
	copy(result, m.tasks)
	return result
}

// HasTask checks if a task with given ID exists
func (m *MockTaskRepository) HasTask(id int) bool {
	for _, task := range m.tasks {
		if task.ID == id {
			return true
		}
	}
	return false
}

// GetTask retrieves a task by ID (for test verification)
func (m *MockTaskRepository) GetTask(id int) (*Task, bool) {
	for _, task := range m.tasks {
		if task.ID == id {
			// Return a copy to avoid accidental modification
			taskCopy := task
			return &taskCopy, true
		}
	}
	return nil, false
}

// TaskCount returns the number of stored tasks
func (m *MockTaskRepository) TaskCount() int {
	return len(m.tasks)
}

// Clear removes all tasks (useful for test cleanup)
func (m *MockTaskRepository) Clear() {
	m.tasks = m.tasks[:0]
	m.saveCallCount = 0
	m.loadCallCount = 0
	m.shouldError = false
	m.errorToReturn = nil
}

// Predefined error scenarios for testing

// MockRepositoryWithLoadError returns a repository that fails on Load
func MockRepositoryWithLoadError() *MockTaskRepository {
	return NewMockRepository().WithError(errors.New("failed to load tasks"))
}

// MockRepositoryWithSaveError returns a repository that fails on Save
func MockRepositoryWithSaveError() *MockTaskRepository {
	return NewMockRepository().WithError(errors.New("failed to save tasks"))
}

// Repository test scenario builders

// RepositoryWithTasks creates a repository pre-loaded with tasks
func RepositoryWithTasks(tasks []Task) *MockTaskRepository {
	return NewMockRepository().WithTasks(tasks)
}

// EmptyRepository creates an empty repository
func EmptyRepository() *MockTaskRepository {
	return NewMockRepository()
}

// ErrorRepository creates a repository that always fails
func ErrorRepository(err error) *MockTaskRepository {
	return NewMockRepository().WithError(err)
}

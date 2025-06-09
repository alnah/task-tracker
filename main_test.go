package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

// Test Utilities and Fixtures

// TestTaskRepository is an in-memory implementation for testing
type TestTaskRepository struct {
	tasks []Task
}

func NewTestTaskRepository() *TestTaskRepository {
	return &TestTaskRepository{
		tasks: []Task{},
	}
}

func (r *TestTaskRepository) Save(tasks []Task) error {
	r.tasks = make([]Task, len(tasks))
	copy(r.tasks, tasks)
	return nil
}

func (r *TestTaskRepository) Load() ([]Task, error) {
	result := make([]Task, len(r.tasks))
	copy(result, r.tasks)
	return result, nil
}

func (r *TestTaskRepository) GetNextID() (int, error) {
	// Mimic the FileTaskRepository behavior
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

// Test fixtures
func createTestTask(t *testing.T, id int, description string) *Task {
	t.Helper()
	task, err := NewTask(id, description)
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}
	return task
}

func createTestTaskWithStatus(t *testing.T, id int, description string, status TaskStatus) *Task {
	t.Helper()
	task := createTestTask(t, id, description)
	task.Status = status
	return task
}

// Domain Layer Tests

func TestNewTask(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		description string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "valid task creation",
			id:          1,
			description: "Buy groceries",
			wantErr:     false,
		},
		{
			name:        "empty description",
			id:          1,
			description: "",
			wantErr:     true,
			expectedErr: ErrEmptyDescription,
		},
		{
			name:        "whitespace only description",
			id:          1,
			description: "   ",
			wantErr:     true,
			expectedErr: ErrEmptyDescription,
		},
		{
			name:        "valid description with whitespace",
			id:          1,
			description: "  Buy groceries  ",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewTask(tt.id, tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewTask() expected error but got none")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("NewTask() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewTask() unexpected error = %v", err)
				return
			}

			// Validate task properties
			if task.ID != tt.id {
				t.Errorf("NewTask() ID = %v, want %v", task.ID, tt.id)
			}

			expectedDesc := strings.TrimSpace(tt.description)
			if task.Description != expectedDesc {
				t.Errorf("NewTask() Description = %v, want %v", task.Description, expectedDesc)
			}

			if task.Status != StatusTodo {
				t.Errorf("NewTask() Status = %v, want %v", task.Status, StatusTodo)
			}

			if task.CreatedAt.IsZero() {
				t.Errorf("NewTask() CreatedAt should not be zero")
			}

			if task.UpdatedAt.IsZero() {
				t.Errorf("NewTask() UpdatedAt should not be zero")
			}

			if !task.CreatedAt.Equal(task.UpdatedAt) {
				t.Errorf("NewTask() CreatedAt and UpdatedAt should be equal for new task")
			}
		})
	}
}

func TestTask_UpdateDescription(t *testing.T) {
	task := createTestTask(t, 1, "Original description")
	originalCreatedAt := task.CreatedAt
	originalUpdatedAt := task.UpdatedAt

	// Small delay to ensure UpdatedAt changes
	time.Sleep(1 * time.Millisecond)

	tests := []struct {
		name        string
		description string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "valid update",
			description: "Updated description",
			wantErr:     false,
		},
		{
			name:        "empty description",
			description: "",
			wantErr:     true,
			expectedErr: ErrEmptyDescription,
		},
		{
			name:        "whitespace only description",
			description: "   ",
			wantErr:     true,
			expectedErr: ErrEmptyDescription,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh task for each test
			testTask := createTestTask(t, 1, "Original description")
			testTask.CreatedAt = originalCreatedAt
			testTask.UpdatedAt = originalUpdatedAt

			time.Sleep(1 * time.Millisecond) // Ensure time difference

			err := testTask.UpdateDescription(tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateDescription() expected error but got none")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("UpdateDescription() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateDescription() unexpected error = %v", err)
				return
			}

			expectedDesc := strings.TrimSpace(tt.description)
			if testTask.Description != expectedDesc {
				t.Errorf(
					"UpdateDescription() Description = %v, want %v",
					testTask.Description,
					expectedDesc,
				)
			}

			if testTask.CreatedAt != originalCreatedAt {
				t.Errorf("UpdateDescription() should not change CreatedAt")
			}

			if !testTask.UpdatedAt.After(originalUpdatedAt) {
				t.Errorf("UpdateDescription() should update UpdatedAt timestamp")
			}
		})
	}
}

func TestTask_MarkInProgress(t *testing.T) {
	task := createTestTask(t, 1, "Test task")
	originalCreatedAt := task.CreatedAt
	originalUpdatedAt := task.UpdatedAt

	time.Sleep(1 * time.Millisecond)
	task.MarkInProgress()

	if task.Status != StatusInProgress {
		t.Errorf("MarkInProgress() Status = %v, want %v", task.Status, StatusInProgress)
	}

	if task.CreatedAt != originalCreatedAt {
		t.Errorf("MarkInProgress() should not change CreatedAt")
	}

	if !task.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("MarkInProgress() should update UpdatedAt timestamp")
	}
}

func TestTask_MarkDone(t *testing.T) {
	task := createTestTask(t, 1, "Test task")
	originalCreatedAt := task.CreatedAt
	originalUpdatedAt := task.UpdatedAt

	time.Sleep(1 * time.Millisecond)
	task.MarkDone()

	if task.Status != StatusDone {
		t.Errorf("MarkDone() Status = %v, want %v", task.Status, StatusDone)
	}

	if task.CreatedAt != originalCreatedAt {
		t.Errorf("MarkDone() should not change CreatedAt")
	}

	if !task.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("MarkDone() should update UpdatedAt timestamp")
	}
}

// Repository Layer Tests

func TestFileTaskRepository_SaveAndLoad(t *testing.T) {
	// Create temporary file
	tmpFile := "test_tasks.json"
	defer os.Remove(tmpFile)

	repo := NewFileTaskRepository(tmpFile)

	// Test data
	tasks := []Task{
		*createTestTask(t, 1, "Task 1"),
		*createTestTaskWithStatus(t, 2, "Task 2", StatusInProgress),
		*createTestTaskWithStatus(t, 3, "Task 3", StatusDone),
	}

	// Test Save
	err := repo.Save(tasks)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Fatalf("Save() should create file")
	}

	// Test Load
	loadedTasks, err := repo.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify loaded data
	if len(loadedTasks) != len(tasks) {
		t.Errorf("Load() returned %d tasks, want %d", len(loadedTasks), len(tasks))
	}

	for i, task := range tasks {
		if i >= len(loadedTasks) {
			continue
		}
		loaded := loadedTasks[i]

		if loaded.ID != task.ID {
			t.Errorf("Load() task[%d].ID = %v, want %v", i, loaded.ID, task.ID)
		}
		if loaded.Description != task.Description {
			t.Errorf(
				"Load() task[%d].Description = %v, want %v",
				i,
				loaded.Description,
				task.Description,
			)
		}
		if loaded.Status != task.Status {
			t.Errorf("Load() task[%d].Status = %v, want %v", i, loaded.Status, task.Status)
		}
	}
}

func TestFileTaskRepository_LoadNonExistentFile(t *testing.T) {
	repo := NewFileTaskRepository("non_existent_file.json")

	tasks, err := repo.Load()
	if err != nil {
		t.Errorf("Load() on non-existent file should not error, got: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Load() on non-existent file should return empty slice, got %d tasks", len(tasks))
	}
}

func TestFileTaskRepository_LoadEmptyFile(t *testing.T) {
	tmpFile := "empty_test.json"
	defer os.Remove(tmpFile)

	// Create empty file
	err := os.WriteFile(tmpFile, []byte(""), 0o644)
	if err != nil {
		t.Fatalf("Failed to create empty test file: %v", err)
	}

	repo := NewFileTaskRepository(tmpFile)
	tasks, err := repo.Load()
	if err != nil {
		t.Errorf("Load() on empty file should not error, got: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Load() on empty file should return empty slice, got %d tasks", len(tasks))
	}
}

func TestFileTaskRepository_GetNextID(t *testing.T) {
	// Create temporary file for real file repository test
	tmpFile := "test_nextid.json"
	defer os.Remove(tmpFile)

	repo := NewFileTaskRepository(tmpFile)

	// Test with empty repository
	id, err := repo.GetNextID()
	if err != nil {
		t.Errorf("GetNextID() failed: %v", err)
	}
	if id != 1 {
		t.Errorf("GetNextID() on empty repo = %v, want 1", id)
	}

	// Add some tasks
	tasks := []Task{
		*createTestTask(t, 1, "Task 1"),
		*createTestTask(t, 3, "Task 3"),
		*createTestTask(t, 2, "Task 2"),
	}
	err = repo.Save(tasks)
	if err != nil {
		t.Fatalf("Failed to save test tasks: %v", err)
	}

	// Test with existing tasks
	id, err = repo.GetNextID()
	if err != nil {
		t.Errorf("GetNextID() failed: %v", err)
	}
	if id != 4 {
		t.Errorf("GetNextID() with max ID 3 = %v, want 4", id)
	}
}

func TestTestTaskRepository_GetNextID(t *testing.T) {
	repo := NewTestTaskRepository()

	// Test with empty repository
	id, err := repo.GetNextID()
	if err != nil {
		t.Errorf("GetNextID() failed: %v", err)
	}
	if id != 1 {
		t.Errorf("GetNextID() on empty repo = %v, want 1", id)
	}

	// Add some tasks and test max ID logic
	tasks := []Task{
		*createTestTask(t, 1, "Task 1"),
		*createTestTask(t, 3, "Task 3"),
		*createTestTask(t, 2, "Task 2"),
	}
	err = repo.Save(tasks)
	if err != nil {
		t.Fatalf("Failed to save test tasks: %v", err)
	}

	// Test with existing tasks - should return max ID + 1
	id, err = repo.GetNextID()
	if err != nil {
		t.Errorf("GetNextID() failed: %v", err)
	}
	if id != 4 {
		t.Errorf("GetNextID() with max ID 3 = %v, want 4", id)
	}
}

// Application Service Tests

func TestTaskService_AddTask(t *testing.T) {
	repo := NewTestTaskRepository()
	service := NewTaskService(repo)

	tests := []struct {
		name        string
		description string
		wantErr     bool
	}{
		{
			name:        "valid task",
			description: "Buy groceries",
			wantErr:     false,
		},
		{
			name:        "empty description",
			description: "",
			wantErr:     true,
		},
		{
			name:        "whitespace description",
			description: "   ",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := service.AddTask(tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("AddTask() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("AddTask() unexpected error = %v", err)
				return
			}

			if task == nil {
				t.Errorf("AddTask() returned nil task")
				return
			}

			// Verify task is saved
			tasks, err := repo.Load()
			if err != nil {
				t.Fatalf("Failed to load tasks for verification: %v", err)
			}
			found := false
			for _, savedTask := range tasks {
				if savedTask.ID == task.ID {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("AddTask() task not saved to repository")
			}
		})
	}
}

func TestTaskService_UpdateTask(t *testing.T) {
	repo := NewTestTaskRepository()
	service := NewTaskService(repo)

	// Add initial task
	task, err := service.AddTask("Original description")
	if err != nil {
		t.Fatalf("Failed to add initial task: %v", err)
	}

	tests := []struct {
		name        string
		taskID      int
		description string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "valid update",
			taskID:      task.ID,
			description: "Updated description",
			wantErr:     false,
		},
		{
			name:        "non-existent task",
			taskID:      999,
			description: "New description",
			wantErr:     true,
			expectedErr: ErrTaskNotFound,
		},
		{
			name:        "empty description",
			taskID:      task.ID,
			description: "",
			wantErr:     true,
			expectedErr: ErrEmptyDescription,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdateTask(tt.taskID, tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateTask() expected error but got none")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("UpdateTask() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateTask() unexpected error = %v", err)
				return
			}

			// Verify task is updated
			tasks, err := repo.Load()
			if err != nil {
				t.Fatalf("Failed to load tasks for verification: %v", err)
			}
			for _, savedTask := range tasks {
				if savedTask.ID == tt.taskID {
					expectedDesc := strings.TrimSpace(tt.description)
					if savedTask.Description != expectedDesc {
						t.Errorf(
							"UpdateTask() description = %v, want %v",
							savedTask.Description,
							expectedDesc,
						)
					}
					return
				}
			}

			t.Errorf("UpdateTask() task not found in repository")
		})
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	repo := NewTestTaskRepository()
	service := NewTaskService(repo)

	// Add initial tasks
	task1, err := service.AddTask("Task 1")
	if err != nil {
		t.Fatalf("Failed to add task1: %v", err)
	}

	task2, err := service.AddTask("Task 2")
	if err != nil {
		t.Fatalf("Failed to add task2: %v", err)
	}

	tests := []struct {
		name        string
		taskID      int
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "valid deletion",
			taskID:  task1.ID,
			wantErr: false,
		},
		{
			name:        "non-existent task",
			taskID:      999,
			wantErr:     true,
			expectedErr: ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteTask(tt.taskID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("DeleteTask() expected error but got none")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("DeleteTask() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("DeleteTask() unexpected error = %v", err)
				return
			}

			// Verify task is deleted
			tasks, err := repo.Load()
			if err != nil {
				t.Fatalf("Failed to load tasks for verification: %v", err)
			}
			for _, savedTask := range tasks {
				if savedTask.ID == tt.taskID {
					t.Errorf("DeleteTask() task still exists in repository")
					return
				}
			}

			// Verify other tasks still exist
			if tt.taskID == task1.ID {
				found := false
				for _, savedTask := range tasks {
					if savedTask.ID == task2.ID {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("DeleteTask() deleted wrong task")
				}
			}
		})
	}
}

func TestTaskService_MarkTaskInProgress(t *testing.T) {
	repo := NewTestTaskRepository()
	service := NewTaskService(repo)

	// Add initial task
	task, err := service.AddTask("Test task")
	if err != nil {
		t.Fatalf("Failed to add initial task: %v", err)
	}

	err = service.MarkTaskInProgress(task.ID)
	if err != nil {
		t.Errorf("MarkTaskInProgress() unexpected error = %v", err)
		return
	}

	// Verify status change
	tasks, err := repo.Load()
	if err != nil {
		t.Fatalf("Failed to load tasks for verification: %v", err)
	}
	for _, savedTask := range tasks {
		if savedTask.ID == task.ID {
			if savedTask.Status != StatusInProgress {
				t.Errorf(
					"MarkTaskInProgress() status = %v, want %v",
					savedTask.Status,
					StatusInProgress,
				)
			}
			return
		}
	}

	t.Errorf("MarkTaskInProgress() task not found in repository")
}

func TestTaskService_MarkTaskDone(t *testing.T) {
	repo := NewTestTaskRepository()
	service := NewTaskService(repo)

	// Add initial task
	task, err := service.AddTask("Test task")
	if err != nil {
		t.Fatalf("Failed to add initial task: %v", err)
	}

	err = service.MarkTaskDone(task.ID)
	if err != nil {
		t.Errorf("MarkTaskDone() unexpected error = %v", err)
		return
	}

	// Verify status change
	tasks, err := repo.Load()
	if err != nil {
		t.Fatalf("Failed to load tasks for verification: %v", err)
	}
	for _, savedTask := range tasks {
		if savedTask.ID == task.ID {
			if savedTask.Status != StatusDone {
				t.Errorf("MarkTaskDone() status = %v, want %v", savedTask.Status, StatusDone)
			}
			return
		}
	}

	t.Errorf("MarkTaskDone() task not found in repository")
}

func TestTaskService_ListTasks(t *testing.T) {
	repo := NewTestTaskRepository()
	service := NewTaskService(repo)

	// Add test data
	_, err := service.AddTask("Task 1") // todo
	if err != nil {
		t.Fatalf("Failed to add task1: %v", err)
	}

	task2, err := service.AddTask("Task 2")
	if err != nil {
		t.Fatalf("Failed to add task2: %v", err)
	}

	err = service.MarkTaskInProgress(task2.ID) // in-progress
	if err != nil {
		t.Fatalf("Failed to mark task2 in progress: %v", err)
	}

	task3, err := service.AddTask("Task 3")
	if err != nil {
		t.Fatalf("Failed to add task3: %v", err)
	}

	err = service.MarkTaskDone(task3.ID) // done
	if err != nil {
		t.Fatalf("Failed to mark task3 done: %v", err)
	}

	tests := []struct {
		name           string
		status         string
		expectedCount  int
		expectedStatus TaskStatus
	}{
		{
			name:          "list all tasks",
			status:        "",
			expectedCount: 3,
		},
		{
			name:           "list todo tasks",
			status:         "todo",
			expectedCount:  1,
			expectedStatus: StatusTodo,
		},
		{
			name:           "list in-progress tasks",
			status:         "in-progress",
			expectedCount:  1,
			expectedStatus: StatusInProgress,
		},
		{
			name:           "list done tasks",
			status:         "done",
			expectedCount:  1,
			expectedStatus: StatusDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := service.ListTasks(tt.status)
			if err != nil {
				t.Errorf("ListTasks() unexpected error = %v", err)
				return
			}

			if len(tasks) != tt.expectedCount {
				t.Errorf("ListTasks() returned %d tasks, want %d", len(tasks), tt.expectedCount)
				return
			}

			// If filtering by status, verify all tasks have that status
			if tt.status != "" {
				for _, task := range tasks {
					if task.Status != tt.expectedStatus {
						t.Errorf(
							"ListTasks() task status = %v, want %v",
							task.Status,
							tt.expectedStatus,
						)
					}
				}
			}
		})
	}
}

// Integration Tests

func TestTaskService_IntegrationWorkflow(t *testing.T) {
	repo := NewTestTaskRepository()
	service := NewTaskService(repo)

	// Test complete workflow

	// 1. Add tasks
	task1, err := service.AddTask("Buy groceries")
	if err != nil {
		t.Fatalf("Failed to add task1: %v", err)
	}

	task2, err := service.AddTask("Complete project")
	if err != nil {
		t.Fatalf("Failed to add task2: %v", err)
	}

	// 2. List all tasks
	tasks, err := service.ListTasks("")
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	// 3. Update task
	err = service.UpdateTask(task1.ID, "Buy groceries and cook dinner")
	if err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	// 4. Mark task in progress
	err = service.MarkTaskInProgress(task1.ID)
	if err != nil {
		t.Fatalf("Failed to mark task in progress: %v", err)
	}

	// 5. Mark task done
	err = service.MarkTaskDone(task2.ID)
	if err != nil {
		t.Fatalf("Failed to mark task done: %v", err)
	}

	// 6. List by status
	todoTasks, err := service.ListTasks("todo")
	if err != nil {
		t.Fatalf("Failed to list todo tasks: %v", err)
	}
	if len(todoTasks) != 0 {
		t.Errorf("Expected 0 todo tasks, got %d", len(todoTasks))
	}

	inProgressTasks, err := service.ListTasks("in-progress")
	if err != nil {
		t.Fatalf("Failed to list in-progress tasks: %v", err)
	}
	if len(inProgressTasks) != 1 {
		t.Errorf("Expected 1 in-progress task, got %d", len(inProgressTasks))
	}

	doneTasks, err := service.ListTasks("done")
	if err != nil {
		t.Fatalf("Failed to list done tasks: %v", err)
	}
	if len(doneTasks) != 1 {
		t.Errorf("Expected 1 done task, got %d", len(doneTasks))
	}

	// 7. Delete task
	err = service.DeleteTask(task1.ID)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	// 8. Verify final state
	finalTasks, err := service.ListTasks("")
	if err != nil {
		t.Fatalf("Failed to list final tasks: %v", err)
	}
	if len(finalTasks) != 1 {
		t.Errorf("Expected 1 final task, got %d", len(finalTasks))
	}
	if finalTasks[0].ID != task2.ID {
		t.Errorf("Wrong task remaining: got ID %d, want %d", finalTasks[0].ID, task2.ID)
	}
}

// Benchmark Tests

func BenchmarkTaskService_AddTask(b *testing.B) {
	repo := NewTestTaskRepository()
	service := NewTaskService(repo)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.AddTask("Benchmark task")
		if err != nil {
			b.Fatalf("AddTask failed: %v", err)
		}
	}
}

func BenchmarkTaskService_ListTasks(b *testing.B) {
	repo := NewTestTaskRepository()
	service := NewTaskService(repo)

	// Setup: Add 1000 tasks
	for i := 0; i < 1000; i++ {
		_, err := service.AddTask("Task " + string(rune(i)))
		if err != nil {
			b.Fatalf("Failed to add setup task %d: %v", i, err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.ListTasks("")
		if err != nil {
			b.Fatalf("ListTasks failed: %v", err)
		}
	}
}

// Property-based testing helper
func TestTaskJSONSerialization(t *testing.T) {
	original := createTestTask(t, 42, "Test task")
	original.MarkInProgress()

	// Serialize to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	// Deserialize from JSON
	var deserialized Task
	err = json.Unmarshal(data, &deserialized)
	if err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// Compare all fields
	if deserialized.ID != original.ID {
		t.Errorf("ID mismatch: got %d, want %d", deserialized.ID, original.ID)
	}
	if deserialized.Description != original.Description {
		t.Errorf(
			"Description mismatch: got %s, want %s",
			deserialized.Description,
			original.Description,
		)
	}
	if deserialized.Status != original.Status {
		t.Errorf("Status mismatch: got %s, want %s", deserialized.Status, original.Status)
	}
	if !deserialized.CreatedAt.Equal(original.CreatedAt) {
		t.Errorf("CreatedAt mismatch: got %v, want %v", deserialized.CreatedAt, original.CreatedAt)
	}
	if !deserialized.UpdatedAt.Equal(original.UpdatedAt) {
		t.Errorf("UpdatedAt mismatch: got %v, want %v", deserialized.UpdatedAt, original.UpdatedAt)
	}
}

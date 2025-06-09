package main

import (
	"os"
	"strings"
	"testing"
)

// TestFileTaskRepository_SaveAndLoad tests persistence functionality
func TestFileTaskRepository_SaveAndLoad(t *testing.T) {
	// Create temporary file
	tmpFile := "test_tasks.json"
	defer os.Remove(tmpFile)

	repo := NewFileTaskRepository(tmpFile)

	t.Run("save and load empty task list", func(t *testing.T) {
		emptyTasks := []Task{}

		err := repo.Save(emptyTasks)
		if err != nil {
			t.Fatalf("Save() failed: %v", err)
		}

		loadedTasks, err := repo.Load()
		if err != nil {
			t.Fatalf("Load() failed: %v", err)
		}

		if len(loadedTasks) != 0 {
			t.Errorf("Load() returned %d tasks, want 0", len(loadedTasks))
		}
	})

	t.Run("save and load single task", func(t *testing.T) {
		task := TodoTask(t)
		tasks := []Task{*task}

		err := repo.Save(tasks)
		if err != nil {
			t.Fatalf("Save() failed: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
			t.Fatalf("Save() should create file")
		}

		loadedTasks, err := repo.Load()
		if err != nil {
			t.Fatalf("Load() failed: %v", err)
		}

		if len(loadedTasks) != 1 {
			t.Errorf("Load() returned %d tasks, want 1", len(loadedTasks))
		}

		AssertTaskEquals(t, task, &loadedTasks[0])
	})

	t.Run("save and load multiple tasks", func(t *testing.T) {
		tasks := MixedStatusTasks(t)

		err := repo.Save(tasks)
		if err != nil {
			t.Fatalf("Save() failed: %v", err)
		}

		loadedTasks, err := repo.Load()
		if err != nil {
			t.Fatalf("Load() failed: %v", err)
		}

		AssertTasksEqual(t, tasks, loadedTasks)
	})

	t.Run("overwrite existing file", func(t *testing.T) {
		// Save initial tasks
		initialTasks := TaskSet(t, 2)
		err := repo.Save(initialTasks)
		if err != nil {
			t.Fatalf("Initial save failed: %v", err)
		}

		// Overwrite with new tasks
		newTasks := TaskSet(t, 3)
		err = repo.Save(newTasks)
		if err != nil {
			t.Fatalf("Overwrite save failed: %v", err)
		}

		// Verify only new tasks are present
		loadedTasks, err := repo.Load()
		if err != nil {
			t.Fatalf("Load after overwrite failed: %v", err)
		}

		if len(loadedTasks) != 3 {
			t.Errorf("Load() after overwrite returned %d tasks, want 3", len(loadedTasks))
		}
	})
}

// TestFileTaskRepository_LoadNonExistentFile tests graceful handling of missing files
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

// TestFileTaskRepository_LoadEmptyFile tests handling of empty files
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

// TestFileTaskRepository_LoadInvalidJSON tests error handling for corrupted files
func TestFileTaskRepository_LoadInvalidJSON(t *testing.T) {
	tmpFile := "invalid_test.json"
	defer os.Remove(tmpFile)

	// Create file with invalid JSON
	err := os.WriteFile(tmpFile, []byte("invalid json content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create invalid JSON test file: %v", err)
	}

	repo := NewFileTaskRepository(tmpFile)
	_, err = repo.Load()
	if err == nil {
		t.Errorf("Load() on invalid JSON should return error")
	}
}

// TestFileTaskRepository_GetNextID tests ID generation
func TestFileTaskRepository_GetNextID(t *testing.T) {
	tmpFile := "test_nextid.json"
	defer os.Remove(tmpFile)

	repo := NewFileTaskRepository(tmpFile)

	t.Run("empty repository returns ID 1", func(t *testing.T) {
		id, err := repo.GetNextID()
		if err != nil {
			t.Errorf("GetNextID() failed: %v", err)
		}
		if id != 1 {
			t.Errorf("GetNextID() on empty repo = %v, want 1", id)
		}
	})

	t.Run("returns max ID + 1", func(t *testing.T) {
		// Add tasks with non-sequential IDs
		tasks := []Task{
			*NewTaskBuilder().WithID(1).BuildValid(t),
			*NewTaskBuilder().WithID(5).BuildValid(t),
			*NewTaskBuilder().WithID(3).BuildValid(t),
		}
		err := repo.Save(tasks)
		if err != nil {
			t.Fatalf("Failed to save test tasks: %v", err)
		}

		id, err := repo.GetNextID()
		if err != nil {
			t.Errorf("GetNextID() failed: %v", err)
		}
		if id != 6 {
			t.Errorf("GetNextID() with max ID 5 = %v, want 6", id)
		}
	})

	t.Run("handles single task", func(t *testing.T) {
		tasks := []Task{*NewTaskBuilder().WithID(42).BuildValid(t)}
		err := repo.Save(tasks)
		if err != nil {
			t.Fatalf("Failed to save test task: %v", err)
		}

		id, err := repo.GetNextID()
		if err != nil {
			t.Errorf("GetNextID() failed: %v", err)
		}
		if id != 43 {
			t.Errorf("GetNextID() with single task ID 42 = %v, want 43", id)
		}
	})
}

// TestFileTaskRepository_FileFormat tests JSON format specifics
func TestFileTaskRepository_FileFormat(t *testing.T) {
	tmpFile := "format_test.json"
	defer os.Remove(tmpFile)

	repo := NewFileTaskRepository(tmpFile)
	task := TodoTask(t)
	tasks := []Task{*task}

	err := repo.Save(tasks)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Read raw file content
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	// Verify it's properly formatted JSON
	contentStr := string(content)
	if !strings.Contains(contentStr, "\"id\"") {
		t.Errorf("Saved JSON should contain id field")
	}
	if !strings.Contains(contentStr, "\"description\"") {
		t.Errorf("Saved JSON should contain description field")
	}
	if !strings.Contains(contentStr, "\"status\"") {
		t.Errorf("Saved JSON should contain status field")
	}

	// Verify it's indented (pretty-printed)
	if !strings.Contains(contentStr, "\n") {
		t.Errorf("Saved JSON should be pretty-printed with newlines")
	}
}

// TestMockRepository tests our mock implementation for consistency
func TestMockRepository(t *testing.T) {
	t.Run("mock behaves like file repository", func(t *testing.T) {
		mock := NewMockRepository()

		// Test empty repository
		tasks, err := mock.Load()
		if err != nil {
			t.Errorf("Mock Load() failed: %v", err)
		}
		if len(tasks) != 0 {
			t.Errorf("Mock Load() on empty repo should return empty slice")
		}

		// Test save and load
		testTasks := MixedStatusTasks(t)
		err = mock.Save(testTasks)
		if err != nil {
			t.Errorf("Mock Save() failed: %v", err)
		}

		loadedTasks, err := mock.Load()
		if err != nil {
			t.Errorf("Mock Load() after save failed: %v", err)
		}

		AssertTasksEqual(t, testTasks, loadedTasks)
	})

	t.Run("mock GetNextID works correctly", func(t *testing.T) {
		mock := NewMockRepository()

		// Empty repository
		id, err := mock.GetNextID()
		if err != nil {
			t.Errorf("Mock GetNextID() failed: %v", err)
		}
		if id != 1 {
			t.Errorf("Mock GetNextID() on empty repo = %v, want 1", id)
		}

		// With tasks
		tasks := []Task{
			*NewTaskBuilder().WithID(3).BuildValid(t),
			*NewTaskBuilder().WithID(1).BuildValid(t),
		}
		err = mock.Save(tasks)
		if err != nil {
			t.Errorf("Mock Save() failed: %v", err)
		}

		id, err = mock.GetNextID()
		if err != nil {
			t.Errorf("Mock GetNextID() failed: %v", err)
		}
		if id != 4 {
			t.Errorf("Mock GetNextID() with max ID 3 = %v, want 4", id)
		}
	})

	t.Run("mock error simulation", func(t *testing.T) {
		mock := NewMockRepository().WithError(ErrTaskNotFound)

		_, err := mock.Load()
		if err != ErrTaskNotFound {
			t.Errorf("Mock Load() should return configured error")
		}

		err = mock.Save([]Task{})
		if err != ErrTaskNotFound {
			t.Errorf("Mock Save() should return configured error")
		}

		_, err = mock.GetNextID()
		if err != ErrTaskNotFound {
			t.Errorf("Mock GetNextID() should return configured error")
		}
	})
}

// TestRepositoryInterface verifies both implementations satisfy the interface
func TestRepositoryInterface(t *testing.T) {
	implementations := []struct {
		name    string
		repo    TaskRepository
		cleanup func()
	}{
		{
			name:    "FileTaskRepository",
			repo:    NewFileTaskRepository("interface_test.json"),
			cleanup: func() { os.Remove("interface_test.json") },
		},
		{
			name:    "MockRepository",
			repo:    NewMockRepository(),
			cleanup: func() {},
		},
	}

	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			defer impl.cleanup()
			repo := impl.repo

			// Test interface compliance through usage
			tasks := TaskSet(t, 2)

			err := repo.Save(tasks)
			if err != nil {
				t.Errorf("%s Save() failed: %v", impl.name, err)
			}

			loadedTasks, err := repo.Load()
			if err != nil {
				t.Errorf("%s Load() failed: %v", impl.name, err)
			}

			if len(loadedTasks) != len(tasks) {
				t.Errorf("%s Load() returned %d tasks, want %d",
					impl.name, len(loadedTasks), len(tasks))
			}

			nextID, err := repo.GetNextID()
			if err != nil {
				t.Errorf("%s GetNextID() failed: %v", impl.name, err)
			}

			if nextID <= 0 {
				t.Errorf("%s GetNextID() returned non-positive ID: %d", impl.name, nextID)
			}
		})
	}
}

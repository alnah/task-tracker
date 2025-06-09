package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// TestTaskService_AddTask tests task creation orchestration
func TestTaskService_AddTask(t *testing.T) {
	t.Run("successful task addition", func(t *testing.T) {
		repo := NewMockRepository()
		service := NewTaskService(repo)

		task, err := service.AddTask("Buy groceries")
		if err != nil {
			t.Errorf("AddTask() unexpected error = %v", err)
			return
		}

		if task == nil {
			t.Errorf("AddTask() returned nil task")
			return
		}

		// Verify task properties
		if task.ID != 1 {
			t.Errorf("AddTask() ID = %v, want 1", task.ID)
		}
		if task.Description != "Buy groceries" {
			t.Errorf("AddTask() Description = %v, want 'Buy groceries'", task.Description)
		}
		if task.Status != StatusTodo {
			t.Errorf("AddTask() Status = %v, want %v", task.Status, StatusTodo)
		}

		// Verify task was saved
		if repo.SaveCallCount() != 1 {
			t.Errorf("AddTask() should call Save() once, called %d times", repo.SaveCallCount())
		}

		savedTasks := repo.GetStoredTasks()
		if len(savedTasks) != 1 {
			t.Errorf("AddTask() should save 1 task, saved %d", len(savedTasks))
		}

		AssertTaskEquals(t, task, &savedTasks[0])
	})

	t.Run("empty description validation", func(t *testing.T) {
		repo := NewMockRepository()
		service := NewTaskService(repo)

		task, err := service.AddTask("")
		if err == nil {
			t.Errorf("AddTask() with empty description should return error")
		}
		if task != nil {
			t.Errorf("AddTask() with error should return nil task")
		}
		if repo.SaveCallCount() != 0 {
			t.Errorf("AddTask() with validation error should not call Save()")
		}
	})

	t.Run("repository error handling", func(t *testing.T) {
		expectedErr := errors.New("repository failed")
		repo := NewMockRepository().WithError(expectedErr)
		service := NewTaskService(repo)

		task, err := service.AddTask("Valid description")
		if err == nil {
			t.Errorf("AddTask() should return error when repository fails")
		}
		if task != nil {
			t.Errorf("AddTask() should return nil task when repository fails")
		}
	})

	t.Run("sequential ID generation", func(t *testing.T) {
		repo := NewMockRepository()
		service := NewTaskService(repo)

		// Add multiple tasks
		task1, err := service.AddTask("Task 1")
		if err != nil {
			t.Fatalf("AddTask() failed for task 1: %v", err)
		}

		task2, err := service.AddTask("Task 2")
		if err != nil {
			t.Fatalf("AddTask() failed for task 2: %v", err)
		}

		task3, err := service.AddTask("Task 3")
		if err != nil {
			t.Fatalf("AddTask() failed for task 3: %v", err)
		}

		// Verify sequential IDs
		if task1.ID != 1 {
			t.Errorf("First task ID = %v, want 1", task1.ID)
		}
		if task2.ID != 2 {
			t.Errorf("Second task ID = %v, want 2", task2.ID)
		}
		if task3.ID != 3 {
			t.Errorf("Third task ID = %v, want 3", task3.ID)
		}

		// Verify all tasks are saved
		savedTasks := repo.GetStoredTasks()
		if len(savedTasks) != 3 {
			t.Errorf("Should have 3 saved tasks, got %d", len(savedTasks))
		}
	})
}

// TestTaskService_UpdateTask tests task modification orchestration
func TestTaskService_UpdateTask(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		existingTask := TodoTask(t)
		repo := NewMockRepository().WithTasks([]Task{*existingTask})
		service := NewTaskService(repo)

		err := service.UpdateTask(existingTask.ID, "Updated description")
		if err != nil {
			t.Errorf("UpdateTask() unexpected error = %v", err)
		}

		// Verify task was updated
		savedTasks := repo.GetStoredTasks()
		if len(savedTasks) != 1 {
			t.Fatalf("Should have 1 task, got %d", len(savedTasks))
		}

		updatedTask := savedTasks[0]
		if updatedTask.Description != "Updated description" {
			t.Errorf(
				"UpdateTask() description = %v, want 'Updated description'",
				updatedTask.Description,
			)
		}

		// Verify other properties unchanged
		if updatedTask.ID != existingTask.ID {
			t.Errorf("UpdateTask() should not change ID")
		}
		if updatedTask.Status != existingTask.Status {
			t.Errorf("UpdateTask() should not change status")
		}
		if updatedTask.CreatedAt != existingTask.CreatedAt {
			t.Errorf("UpdateTask() should not change CreatedAt")
		}
	})

	t.Run("task not found", func(t *testing.T) {
		repo := NewMockRepository()
		service := NewTaskService(repo)

		err := service.UpdateTask(999, "New description")
		if err != ErrTaskNotFound {
			t.Errorf("UpdateTask() error = %v, want %v", err, ErrTaskNotFound)
		}

		// Verify no save was attempted
		if repo.SaveCallCount() != 0 {
			t.Errorf("UpdateTask() with non-existent task should not call Save()")
		}
	})

	t.Run("empty description validation", func(t *testing.T) {
		existingTask := TodoTask(t)
		repo := NewMockRepository().WithTasks([]Task{*existingTask})
		service := NewTaskService(repo)

		err := service.UpdateTask(existingTask.ID, "")
		if err != ErrEmptyDescription {
			t.Errorf(
				"UpdateTask() with empty description error = %v, want %v",
				err,
				ErrEmptyDescription,
			)
		}

		// Verify task was not modified
		savedTasks := repo.GetStoredTasks()
		AssertTaskEquals(t, existingTask, &savedTasks[0])
	})
}

// TestTaskService_DeleteTask tests task removal orchestration
func TestTaskService_DeleteTask(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		tasks := MixedStatusTasks(t)
		repo := NewMockRepository().WithTasks(tasks)
		service := NewTaskService(repo)

		taskToDelete := tasks[1] // Middle task
		err := service.DeleteTask(taskToDelete.ID)
		if err != nil {
			t.Errorf("DeleteTask() unexpected error = %v", err)
		}

		// Verify task was removed
		savedTasks := repo.GetStoredTasks()
		if len(savedTasks) != 2 {
			t.Errorf("DeleteTask() should leave 2 tasks, got %d", len(savedTasks))
		}

		// Verify correct task was removed
		for _, task := range savedTasks {
			if task.ID == taskToDelete.ID {
				t.Errorf("DeleteTask() should remove task with ID %d", taskToDelete.ID)
			}
		}

		// Verify other tasks remain
		remainingIDs := make(map[int]bool)
		for _, task := range savedTasks {
			remainingIDs[task.ID] = true
		}

		for _, originalTask := range tasks {
			if originalTask.ID != taskToDelete.ID {
				if !remainingIDs[originalTask.ID] {
					t.Errorf("DeleteTask() should preserve task with ID %d", originalTask.ID)
				}
			}
		}
	})

	t.Run("task not found", func(t *testing.T) {
		repo := NewMockRepository()
		service := NewTaskService(repo)

		err := service.DeleteTask(999)
		if err != ErrTaskNotFound {
			t.Errorf("DeleteTask() error = %v, want %v", err, ErrTaskNotFound)
		}

		// Verify no save was attempted
		if repo.SaveCallCount() != 0 {
			t.Errorf("DeleteTask() with non-existent task should not call Save()")
		}
	})

	t.Run("delete from single task list", func(t *testing.T) {
		task := TodoTask(t)
		repo := NewMockRepository().WithTasks([]Task{*task})
		service := NewTaskService(repo)

		err := service.DeleteTask(task.ID)
		if err != nil {
			t.Errorf("DeleteTask() unexpected error = %v", err)
		}

		// Verify list is empty
		savedTasks := repo.GetStoredTasks()
		if len(savedTasks) != 0 {
			t.Errorf("DeleteTask() should result in empty list, got %d tasks", len(savedTasks))
		}
	})
}

// TestTaskService_MarkTaskInProgress tests status change orchestration
func TestTaskService_MarkTaskInProgress(t *testing.T) {
	t.Run("successful status change", func(t *testing.T) {
		task := TodoTask(t)
		repo := NewMockRepository().WithTasks([]Task{*task})
		service := NewTaskService(repo)

		err := service.MarkTaskInProgress(task.ID)
		if err != nil {
			t.Errorf("MarkTaskInProgress() unexpected error = %v", err)
		}

		// Verify status was changed
		savedTasks := repo.GetStoredTasks()
		if len(savedTasks) != 1 {
			t.Fatalf("Should have 1 task, got %d", len(savedTasks))
		}

		updatedTask := savedTasks[0]
		if updatedTask.Status != StatusInProgress {
			t.Errorf(
				"MarkTaskInProgress() status = %v, want %v",
				updatedTask.Status,
				StatusInProgress,
			)
		}

		// Verify other properties unchanged
		if updatedTask.ID != task.ID {
			t.Errorf("MarkTaskInProgress() should not change ID")
		}
		if updatedTask.Description != task.Description {
			t.Errorf("MarkTaskInProgress() should not change description")
		}
		if updatedTask.CreatedAt != task.CreatedAt {
			t.Errorf("MarkTaskInProgress() should not change CreatedAt")
		}
	})

	t.Run("task not found", func(t *testing.T) {
		repo := NewMockRepository()
		service := NewTaskService(repo)

		err := service.MarkTaskInProgress(999)
		if err != ErrTaskNotFound {
			t.Errorf("MarkTaskInProgress() error = %v, want %v", err, ErrTaskNotFound)
		}
	})

	t.Run("idempotent operation", func(t *testing.T) {
		task := InProgressTask(t)
		repo := NewMockRepository().WithTasks([]Task{*task})
		service := NewTaskService(repo)

		err := service.MarkTaskInProgress(task.ID)
		if err != nil {
			t.Errorf("MarkTaskInProgress() on already in-progress task should not error")
		}

		// Verify status remains in-progress
		savedTasks := repo.GetStoredTasks()
		if savedTasks[0].Status != StatusInProgress {
			t.Errorf("MarkTaskInProgress() should remain in-progress")
		}
	})
}

// TestTaskService_MarkTaskDone tests completion orchestration
func TestTaskService_MarkTaskDone(t *testing.T) {
	t.Run("successful completion", func(t *testing.T) {
		task := InProgressTask(t)
		repo := NewMockRepository().WithTasks([]Task{*task})
		service := NewTaskService(repo)

		err := service.MarkTaskDone(task.ID)
		if err != nil {
			t.Errorf("MarkTaskDone() unexpected error = %v", err)
		}

		// Verify status was changed
		savedTasks := repo.GetStoredTasks()
		updatedTask := savedTasks[0]
		if updatedTask.Status != StatusDone {
			t.Errorf("MarkTaskDone() status = %v, want %v", updatedTask.Status, StatusDone)
		}
	})

	t.Run("complete todo task directly", func(t *testing.T) {
		task := TodoTask(t)
		repo := NewMockRepository().WithTasks([]Task{*task})
		service := NewTaskService(repo)

		err := service.MarkTaskDone(task.ID)
		if err != nil {
			t.Errorf("MarkTaskDone() on todo task should work")
		}

		savedTasks := repo.GetStoredTasks()
		if savedTasks[0].Status != StatusDone {
			t.Errorf("MarkTaskDone() should change todo to done")
		}
	})

	t.Run("task not found", func(t *testing.T) {
		repo := NewMockRepository()
		service := NewTaskService(repo)

		err := service.MarkTaskDone(999)
		if err != ErrTaskNotFound {
			t.Errorf("MarkTaskDone() error = %v, want %v", err, ErrTaskNotFound)
		}
	})
}

// TestTaskService_ListTasks tests task retrieval and filtering
func TestTaskService_ListTasks(t *testing.T) {
	// Setup test data
	tasks := MixedStatusTasks(t)
	repo := NewMockRepository().WithTasks(tasks)
	service := NewTaskService(repo)

	t.Run("list all tasks", func(t *testing.T) {
		result, err := service.ListTasks("")
		if err != nil {
			t.Errorf("ListTasks() unexpected error = %v", err)
		}

		if len(result) != 3 {
			t.Errorf("ListTasks() returned %d tasks, want 3", len(result))
		}

		AssertTasksEqual(t, tasks, result)
	})

	t.Run("list todo tasks", func(t *testing.T) {
		result, err := service.ListTasks("todo")
		if err != nil {
			t.Errorf("ListTasks(todo) unexpected error = %v", err)
		}

		expectedCount := 1
		if len(result) != expectedCount {
			t.Errorf("ListTasks(todo) returned %d tasks, want %d", len(result), expectedCount)
		}

		for _, task := range result {
			if task.Status != StatusTodo {
				t.Errorf("ListTasks(todo) returned task with status %v", task.Status)
			}
		}
	})

	t.Run("list in-progress tasks", func(t *testing.T) {
		result, err := service.ListTasks("in-progress")
		if err != nil {
			t.Errorf("ListTasks(in-progress) unexpected error = %v", err)
		}

		expectedCount := 1
		if len(result) != expectedCount {
			t.Errorf(
				"ListTasks(in-progress) returned %d tasks, want %d",
				len(result),
				expectedCount,
			)
		}

		for _, task := range result {
			if task.Status != StatusInProgress {
				t.Errorf("ListTasks(in-progress) returned task with status %v", task.Status)
			}
		}
	})

	t.Run("list done tasks", func(t *testing.T) {
		result, err := service.ListTasks("done")
		if err != nil {
			t.Errorf("ListTasks(done) unexpected error = %v", err)
		}

		expectedCount := 1
		if len(result) != expectedCount {
			t.Errorf("ListTasks(done) returned %d tasks, want %d", len(result), expectedCount)
		}

		for _, task := range result {
			if task.Status != StatusDone {
				t.Errorf("ListTasks(done) returned task with status %v", task.Status)
			}
		}
	})

	t.Run("list from empty repository", func(t *testing.T) {
		emptyRepo := NewMockRepository()
		emptyService := NewTaskService(emptyRepo)

		result, err := emptyService.ListTasks("")
		if err != nil {
			t.Errorf("ListTasks() on empty repo unexpected error = %v", err)
		}

		if len(result) != 0 {
			t.Errorf("ListTasks() on empty repo returned %d tasks, want 0", len(result))
		}
	})

	t.Run("list with non-existent status", func(t *testing.T) {
		result, err := service.ListTasks("invalid-status")
		if err != nil {
			t.Errorf("ListTasks(invalid-status) unexpected error = %v", err)
		}

		if len(result) != 0 {
			t.Errorf("ListTasks(invalid-status) returned %d tasks, want 0", len(result))
		}
	})

	t.Run("repository error handling", func(t *testing.T) {
		expectedErr := errors.New("load failed")
		errorRepo := NewMockRepository().WithError(expectedErr)
		errorService := NewTaskService(errorRepo)

		_, err := errorService.ListTasks("")
		if err == nil {
			t.Errorf("ListTasks() should return error when repository fails")
		}
	})
}

// TestTaskService_EdgeCases tests unusual but valid scenarios
func TestTaskService_EdgeCases(t *testing.T) {
	t.Run("very long task description", func(t *testing.T) {
		repo := NewMockRepository()
		service := NewTaskService(repo)

		// Create a long description without trailing whitespace
		longDescription := strings.Repeat("Very long task description", 100)
		task, err := service.AddTask(longDescription)
		if err != nil {
			t.Errorf("AddTask() with long description should succeed: %v", err)
		}

		if task.Description != longDescription {
			t.Errorf("AddTask() should preserve long description")
		}
	})

	t.Run("task description with special characters", func(t *testing.T) {
		repo := NewMockRepository()
		service := NewTaskService(repo)

		specialDescription := "Task with émojis 🎯 and symbols @#$%^&*()"
		task, err := service.AddTask(specialDescription)
		if err != nil {
			t.Errorf("AddTask() with special characters should succeed: %v", err)
		}

		if task.Description != specialDescription {
			t.Errorf("AddTask() should preserve special characters")
		}
	})

	t.Run("operations on large task list", func(t *testing.T) {
		repo := NewMockRepository()
		service := NewTaskService(repo)

		// Add many tasks
		const taskCount = 100 // Reduced for faster tests
		for i := range taskCount {
			_, err := service.AddTask(fmt.Sprintf("Task %d", i+1))
			if err != nil {
				t.Fatalf("AddTask() %d failed: %v", i+1, err)
			}
		}

		// List all tasks
		tasks, err := service.ListTasks("")
		if err != nil {
			t.Fatalf("ListTasks() on large list failed: %v", err)
		}

		if len(tasks) != taskCount {
			t.Errorf("ListTasks() returned %d tasks, want %d", len(tasks), taskCount)
		}

		// Delete middle task
		middleID := taskCount / 2
		err = service.DeleteTask(middleID)
		if err != nil {
			t.Errorf("DeleteTask() from large list failed: %v", err)
		}

		// Verify deletion
		tasksAfterDelete, err := service.ListTasks("")
		if err != nil {
			t.Errorf("ListTasks() after delete failed: %v", err)
		}

		if len(tasksAfterDelete) != taskCount-1 {
			t.Errorf("After delete, expected %d tasks, got %d", taskCount-1, len(tasksAfterDelete))
		}
	})
}

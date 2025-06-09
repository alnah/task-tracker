package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestNewTask focuses on task creation business rules
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
			name:        "empty description rejected",
			id:          1,
			description: "",
			wantErr:     true,
			expectedErr: ErrEmptyDescription,
		},
		{
			name:        "whitespace only description rejected",
			id:          1,
			description: "   ",
			wantErr:     true,
			expectedErr: ErrEmptyDescription,
		},
		{
			name:        "whitespace trimmed from valid description",
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

			// Verify business invariants
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

			// Verify timestamps are set
			if task.CreatedAt.IsZero() {
				t.Errorf("NewTask() CreatedAt should not be zero")
			}

			if task.UpdatedAt.IsZero() {
				t.Errorf("NewTask() UpdatedAt should not be zero")
			}

			// For new tasks, timestamps should be equal
			if !task.CreatedAt.Equal(task.UpdatedAt) {
				t.Errorf("NewTask() CreatedAt and UpdatedAt should be equal for new task")
			}
		})
	}
}

// TestTask_UpdateDescription tests description update business rules
func TestTask_UpdateDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "valid description update",
			description: "Updated description",
			wantErr:     false,
		},
		{
			name:        "empty description rejected",
			description: "",
			wantErr:     true,
			expectedErr: ErrEmptyDescription,
		},
		{
			name:        "whitespace only description rejected",
			description: "   ",
			wantErr:     true,
			expectedErr: ErrEmptyDescription,
		},
		{
			name:        "whitespace trimmed from description",
			description: "  Trimmed description  ",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh task for each test
			task := NewTaskBuilder().
				WithDescription("Original description").
				WithTimestamps(FixedTime(), FixedTime()).
				BuildValid(t)

			originalCreatedAt := task.CreatedAt
			originalUpdatedAt := task.UpdatedAt

			// Small delay to ensure UpdatedAt changes
			time.Sleep(1 * time.Millisecond)

			err := task.UpdateDescription(tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateDescription() expected error but got none")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("UpdateDescription() error = %v, expectedErr %v", err, tt.expectedErr)
				}

				// Verify state unchanged on error
				if task.Description != "Original description" {
					t.Errorf("UpdateDescription() should not change description on error")
				}
				if task.UpdatedAt != originalUpdatedAt {
					t.Errorf("UpdateDescription() should not change UpdatedAt on error")
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateDescription() unexpected error = %v", err)
				return
			}

			// Verify business rules
			expectedDesc := strings.TrimSpace(tt.description)
			if task.Description != expectedDesc {
				t.Errorf(
					"UpdateDescription() Description = %v, want %v",
					task.Description,
					expectedDesc,
				)
			}

			// CreatedAt should remain unchanged
			if task.CreatedAt != originalCreatedAt {
				t.Errorf("UpdateDescription() should not change CreatedAt")
			}

			// UpdatedAt should be updated
			if !task.UpdatedAt.After(originalUpdatedAt) {
				t.Errorf("UpdateDescription() should update UpdatedAt timestamp")
			}
		})
	}
}

// TestTask_MarkInProgress tests state transition business rules
func TestTask_MarkInProgress(t *testing.T) {
	tests := []struct {
		name             string
		initialStatus    TaskStatus
		shouldTransition bool
	}{
		{
			name:             "todo to in-progress",
			initialStatus:    StatusTodo,
			shouldTransition: true,
		},
		{
			name:             "done to in-progress",
			initialStatus:    StatusDone,
			shouldTransition: true,
		},
		{
			name:             "in-progress to in-progress (idempotent)",
			initialStatus:    StatusInProgress,
			shouldTransition: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := NewTaskBuilder().
				WithStatus(tt.initialStatus).
				WithTimestamps(FixedTime(), FixedTime()).
				BuildValid(t)

			originalCreatedAt := task.CreatedAt
			originalUpdatedAt := task.UpdatedAt

			time.Sleep(1 * time.Millisecond)
			task.MarkInProgress()

			// Verify state transition
			if task.Status != StatusInProgress {
				t.Errorf("MarkInProgress() Status = %v, want %v", task.Status, StatusInProgress)
			}

			// Verify timestamps
			if task.CreatedAt != originalCreatedAt {
				t.Errorf("MarkInProgress() should not change CreatedAt")
			}

			if !task.UpdatedAt.After(originalUpdatedAt) {
				t.Errorf("MarkInProgress() should update UpdatedAt timestamp")
			}
		})
	}
}

// TestTask_MarkDone tests completion business rules
func TestTask_MarkDone(t *testing.T) {
	tests := []struct {
		name             string
		initialStatus    TaskStatus
		shouldTransition bool
	}{
		{
			name:             "todo to done",
			initialStatus:    StatusTodo,
			shouldTransition: true,
		},
		{
			name:             "in-progress to done",
			initialStatus:    StatusInProgress,
			shouldTransition: true,
		},
		{
			name:             "done to done (idempotent)",
			initialStatus:    StatusDone,
			shouldTransition: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := NewTaskBuilder().
				WithStatus(tt.initialStatus).
				WithTimestamps(FixedTime(), FixedTime()).
				BuildValid(t)

			originalCreatedAt := task.CreatedAt
			originalUpdatedAt := task.UpdatedAt

			time.Sleep(1 * time.Millisecond)
			task.MarkDone()

			// Verify state transition
			if task.Status != StatusDone {
				t.Errorf("MarkDone() Status = %v, want %v", task.Status, StatusDone)
			}

			// Verify timestamps
			if task.CreatedAt != originalCreatedAt {
				t.Errorf("MarkDone() should not change CreatedAt")
			}

			if !task.UpdatedAt.After(originalUpdatedAt) {
				t.Errorf("MarkDone() should update UpdatedAt timestamp")
			}
		})
	}
}

// TestTask_StateInvariants tests that task maintains valid state
func TestTask_StateInvariants(t *testing.T) {
	t.Run("new task has valid initial state", func(t *testing.T) {
		task := TodoTask(t)

		// Business invariants
		if task.ID <= 0 {
			t.Errorf("Task ID should be positive, got %d", task.ID)
		}

		if task.Description == "" {
			t.Errorf("Task description should not be empty")
		}

		if task.Status != StatusTodo {
			t.Errorf("New task should have todo status, got %v", task.Status)
		}

		if task.CreatedAt.IsZero() {
			t.Errorf("Task should have creation timestamp")
		}

		if task.UpdatedAt.IsZero() {
			t.Errorf("Task should have update timestamp")
		}

		if task.UpdatedAt.Before(task.CreatedAt) {
			t.Errorf("UpdatedAt should not be before CreatedAt")
		}
	})

	t.Run("multiple state changes maintain invariants", func(t *testing.T) {
		task := TodoTask(t)
		originalCreatedAt := task.CreatedAt

		// Chain multiple operations
		time.Sleep(1 * time.Millisecond)
		task.MarkInProgress()
		firstUpdate := task.UpdatedAt

		time.Sleep(1 * time.Millisecond)
		err := task.UpdateDescription("Updated description")
		if err != nil {
			t.Fatalf("UpdateDescription failed: %v", err)
		}
		secondUpdate := task.UpdatedAt

		time.Sleep(1 * time.Millisecond)
		task.MarkDone()
		thirdUpdate := task.UpdatedAt

		// Verify timestamp progression
		if task.CreatedAt != originalCreatedAt {
			t.Errorf("CreatedAt should never change")
		}

		if !firstUpdate.After(originalCreatedAt) {
			t.Errorf("First update should be after creation")
		}

		if !secondUpdate.After(firstUpdate) {
			t.Errorf("Second update should be after first")
		}

		if !thirdUpdate.After(secondUpdate) {
			t.Errorf("Third update should be after second")
		}

		// Verify final state
		if task.Status != StatusDone {
			t.Errorf("Final status should be done")
		}

		if task.Description != "Updated description" {
			t.Errorf("Description should be updated")
		}
	})
}

// TestTaskStatus_Values tests status value objects
func TestTaskStatus_Values(t *testing.T) {
	validStatuses := []TaskStatus{StatusTodo, StatusInProgress, StatusDone}

	for _, status := range validStatuses {
		t.Run(string(status), func(t *testing.T) {
			// Verify status string representation
			statusStr := string(status)
			if statusStr == "" {
				t.Errorf("Status should have string representation")
			}

			// Verify we can create tasks with this status
			task := NewTaskBuilder().
				WithStatus(status).
				BuildValid(t)

			if task.Status != status {
				t.Errorf("Task status = %v, want %v", task.Status, status)
			}
		})
	}
}

// TestTask_ImmutableCreationTime ensures creation time never changes
func TestTask_ImmutableCreationTime(t *testing.T) {
	task := TodoTask(t)
	originalCreatedAt := task.CreatedAt

	// Perform various operations
	operations := []func(){
		func() { task.MarkInProgress() },
		func() { task.MarkDone() },
		func() { _ = task.UpdateDescription("New description") },
	}

	for i, op := range operations {
		t.Run(fmt.Sprintf("operation_%d", i), func(t *testing.T) {
			time.Sleep(1 * time.Millisecond) // Ensure time difference
			op()

			if task.CreatedAt != originalCreatedAt {
				t.Errorf("CreatedAt should never change, was %v, now %v",
					originalCreatedAt, task.CreatedAt)
			}

			if !task.UpdatedAt.After(originalCreatedAt) {
				t.Errorf("UpdatedAt should be after CreatedAt")
			}
		})
	}
}

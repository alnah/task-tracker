package main

import (
	"fmt"
	"os"
	"slices"
	"testing"
	"time"
)

// TestFullTaskLifecycle tests complete task management workflow
func TestFullTaskLifecycle(t *testing.T) {
	// Use real file repository for true integration test
	tmpFile := "integration_test_tasks.json"
	defer os.Remove(tmpFile)

	repo := NewFileTaskRepository(tmpFile)
	service := NewTaskService(repo)

	t.Run("complete task lifecycle", func(t *testing.T) {
		// 1. Start with empty system
		tasks, err := service.ListTasks("")
		if err != nil {
			t.Fatalf("Initial ListTasks() failed: %v", err)
		}
		if len(tasks) != 0 {
			t.Errorf("Initial system should be empty, got %d tasks", len(tasks))
		}

		// 2. Add multiple tasks
		task1, err := service.AddTask("Buy groceries")
		if err != nil {
			t.Fatalf("AddTask() 1 failed: %v", err)
		}

		task2, err := service.AddTask("Complete project report")
		if err != nil {
			t.Fatalf("AddTask() 2 failed: %v", err)
		}

		task3, err := service.AddTask("Call mom")
		if err != nil {
			t.Fatalf("AddTask() 3 failed: %v", err)
		}

		// 3. Verify all tasks are present
		allTasks, err := service.ListTasks("")
		if err != nil {
			t.Fatalf("ListTasks() after adding failed: %v", err)
		}
		if len(allTasks) != 3 {
			t.Errorf("Should have 3 tasks after adding, got %d", len(allTasks))
		}

		// 4. Start working on first task
		err = service.MarkTaskInProgress(task1.ID)
		if err != nil {
			t.Fatalf("MarkTaskInProgress() failed: %v", err)
		}

		// 5. Complete second task directly
		err = service.MarkTaskDone(task2.ID)
		if err != nil {
			t.Fatalf("MarkTaskDone() failed: %v", err)
		}

		// 6. Update third task description
		err = service.UpdateTask(task3.ID, "Call mom and discuss weekend plans")
		if err != nil {
			t.Fatalf("UpdateTask() failed: %v", err)
		}

		// 7. Verify task statuses
		todoTasks, err := service.ListTasks("todo")
		if err != nil {
			t.Fatalf("ListTasks(todo) failed: %v", err)
		}
		if len(todoTasks) != 1 {
			t.Errorf("Should have 1 todo task, got %d", len(todoTasks))
		}
		if todoTasks[0].ID != task3.ID {
			t.Errorf("Wrong task in todo list")
		}
		if todoTasks[0].Description != "Call mom and discuss weekend plans" {
			t.Errorf("Task description not updated correctly")
		}

		inProgressTasks, err := service.ListTasks("in-progress")
		if err != nil {
			t.Fatalf("ListTasks(in-progress) failed: %v", err)
		}
		if len(inProgressTasks) != 1 {
			t.Errorf("Should have 1 in-progress task, got %d", len(inProgressTasks))
		}
		if inProgressTasks[0].ID != task1.ID {
			t.Errorf("Wrong task in in-progress list")
		}

		doneTasks, err := service.ListTasks("done")
		if err != nil {
			t.Fatalf("ListTasks(done) failed: %v", err)
		}
		if len(doneTasks) != 1 {
			t.Errorf("Should have 1 done task, got %d", len(doneTasks))
		}
		if doneTasks[0].ID != task2.ID {
			t.Errorf("Wrong task in done list")
		}

		// 8. Complete the in-progress task
		err = service.MarkTaskDone(task1.ID)
		if err != nil {
			t.Fatalf("MarkTaskDone() for task1 failed: %v", err)
		}

		// 9. Delete one completed task
		err = service.DeleteTask(task2.ID)
		if err != nil {
			t.Fatalf("DeleteTask() failed: %v", err)
		}

		// 10. Verify final state
		finalTasks, err := service.ListTasks("")
		if err != nil {
			t.Fatalf("Final ListTasks() failed: %v", err)
		}
		if len(finalTasks) != 2 {
			t.Errorf("Should have 2 tasks at end, got %d", len(finalTasks))
		}

		// Verify remaining tasks are correct
		remainingIDs := make(map[int]bool)
		for _, task := range finalTasks {
			remainingIDs[task.ID] = true
		}

		if !remainingIDs[task1.ID] {
			t.Errorf("Task 1 should remain")
		}
		if remainingIDs[task2.ID] {
			t.Errorf("Task 2 should be deleted")
		}
		if !remainingIDs[task3.ID] {
			t.Errorf("Task 3 should remain")
		}

		// 11. Verify persistence by creating new service instance
		newService := NewTaskService(NewFileTaskRepository(tmpFile))
		persistedTasks, err := newService.ListTasks("")
		if err != nil {
			t.Fatalf("ListTasks() with new service failed: %v", err)
		}

		if len(persistedTasks) != 2 {
			t.Errorf("Persisted tasks count = %d, want 2", len(persistedTasks))
		}

		AssertTasksEqual(t, finalTasks, persistedTasks)
	})
}

// TestConcurrentAccess tests behavior when multiple service instances access same file
func TestConcurrentAccess(t *testing.T) {
	tmpFile := "concurrent_test_tasks.json"
	defer os.Remove(tmpFile)

	t.Run("multiple service instances", func(t *testing.T) {
		// Create two service instances sharing the same file
		service1 := NewTaskService(NewFileTaskRepository(tmpFile))
		service2 := NewTaskService(NewFileTaskRepository(tmpFile))

		// Service 1 adds a task
		task1, err := service1.AddTask("Task from service 1")
		if err != nil {
			t.Fatalf("Service1 AddTask() failed: %v", err)
		}

		// Service 2 should see the task
		tasks2, err := service2.ListTasks("")
		if err != nil {
			t.Fatalf("Service2 ListTasks() failed: %v", err)
		}
		if len(tasks2) != 1 {
			t.Errorf("Service2 should see 1 task, got %d", len(tasks2))
		}

		// Service 2 adds another task
		task2, err := service2.AddTask("Task from service 2")
		if err != nil {
			t.Fatalf("Service2 AddTask() failed: %v", err)
		}

		// Service 1 should see both tasks
		tasks1, err := service1.ListTasks("")
		if err != nil {
			t.Fatalf("Service1 ListTasks() after service2 add failed: %v", err)
		}
		if len(tasks1) != 2 {
			t.Errorf("Service1 should see 2 tasks, got %d", len(tasks1))
		}

		// Verify task IDs are unique
		if task1.ID == task2.ID {
			t.Errorf("Tasks should have unique IDs, both have %d", task1.ID)
		}

		// Service 1 modifies task created by service 2
		err = service1.MarkTaskInProgress(task2.ID)
		if err != nil {
			t.Fatalf("Service1 MarkTaskInProgress() on service2's task failed: %v", err)
		}

		// Service 2 should see the modification
		updatedTasks, err := service2.ListTasks("")
		if err != nil {
			t.Fatalf("Service2 ListTasks() after modification failed: %v", err)
		}

		var modifiedTask *Task
		for _, task := range updatedTasks {
			if task.ID == task2.ID {
				modifiedTask = &task
				break
			}
		}

		if modifiedTask == nil {
			t.Fatalf("Modified task not found")
		}

		if modifiedTask.Status != StatusInProgress {
			t.Errorf("Task status should be in-progress, got %v", modifiedTask.Status)
		}
	})
}

// TestDataPersistence tests that data survives application restarts
func TestDataPersistence(t *testing.T) {
	tmpFile := "persistence_test_tasks.json"
	defer os.Remove(tmpFile)

	originalTasks := MixedStatusTasks(t)

	t.Run("data survives service recreation", func(t *testing.T) {
		// Create service and add tasks
		{
			service := NewTaskService(NewFileTaskRepository(tmpFile))
			for _, task := range originalTasks {
				_, err := service.AddTask(task.Description)
				if err != nil {
					t.Fatalf("AddTask() failed: %v", err)
				}
			}

			// Modify some tasks
			err := service.MarkTaskInProgress(2)
			if err != nil {
				t.Fatalf("MarkTaskInProgress() failed: %v", err)
			}

			err = service.MarkTaskDone(3)
			if err != nil {
				t.Fatalf("MarkTaskDone() failed: %v", err)
			}
		} // Service goes out of scope

		// Create new service instance
		{
			newService := NewTaskService(NewFileTaskRepository(tmpFile))

			// Verify data persisted
			tasks, err := newService.ListTasks("")
			if err != nil {
				t.Fatalf("ListTasks() with new service failed: %v", err)
			}

			if len(tasks) != len(originalTasks) {
				t.Errorf("Expected %d tasks, got %d", len(originalTasks), len(tasks))
			}

			// Verify statuses were persisted
			inProgressTasks, err := newService.ListTasks("in-progress")
			if err != nil {
				t.Fatalf("ListTasks(in-progress) failed: %v", err)
			}
			if len(inProgressTasks) != 1 {
				t.Errorf("Expected 1 in-progress task, got %d", len(inProgressTasks))
			}
			if inProgressTasks[0].ID != 2 {
				t.Errorf("Wrong task marked as in-progress")
			}

			doneTasks, err := newService.ListTasks("done")
			if err != nil {
				t.Fatalf("ListTasks(done) failed: %v", err)
			}
			if len(doneTasks) != 1 {
				t.Errorf("Expected 1 done task, got %d", len(doneTasks))
			}
			if doneTasks[0].ID != 3 {
				t.Errorf("Wrong task marked as done")
			}

			// Verify we can continue operations
			_, err = newService.AddTask("New task after restart")
			if err != nil {
				t.Errorf("AddTask() after restart failed: %v", err)
			}

			finalTasks, err := newService.ListTasks("")
			if err != nil {
				t.Fatalf("Final ListTasks() failed: %v", err)
			}
			if len(finalTasks) != len(originalTasks)+1 {
				t.Errorf(
					"Expected %d tasks after adding, got %d",
					len(originalTasks)+1,
					len(finalTasks),
				)
			}
		}
	})
}

// TestErrorRecovery tests system behavior after various error conditions
func TestErrorRecovery(t *testing.T) {
	tmpFile := "error_recovery_test_tasks.json"
	defer os.Remove(tmpFile)

	t.Run("recovery from file corruption", func(t *testing.T) {
		service := NewTaskService(NewFileTaskRepository(tmpFile))

		// Add some tasks normally
		_, err := service.AddTask("Task 1")
		if err != nil {
			t.Fatalf("Initial AddTask() failed: %v", err)
		}

		// Corrupt the file
		err = os.WriteFile(tmpFile, []byte("corrupted json"), 0o644)
		if err != nil {
			t.Fatalf("Failed to corrupt file: %v", err)
		}

		// Service should handle corruption gracefully
		_, err = service.ListTasks("")
		if err == nil {
			t.Errorf("ListTasks() should fail with corrupted file")
		}

		// Service should continue working if we remove the corrupted file
		err = os.Remove(tmpFile)
		if err != nil {
			t.Fatalf("Failed to remove corrupted file: %v", err)
		}

		// Now operations should work again (starting fresh)
		_, err = service.AddTask("Recovery task")
		if err != nil {
			t.Errorf("AddTask() should work after removing corrupted file: %v", err)
		}

		tasks, err := service.ListTasks("")
		if err != nil {
			t.Errorf("ListTasks() should work after recovery: %v", err)
		}
		if len(tasks) != 1 {
			t.Errorf("Should have 1 task after recovery, got %d", len(tasks))
		}
	})
}

// TestPerformanceCharacteristics tests system performance under load
func TestPerformanceCharacteristics(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	tmpFile := "performance_test_tasks.json"
	defer os.Remove(tmpFile)

	service := NewTaskService(NewFileTaskRepository(tmpFile))

	t.Run("large task list performance", func(t *testing.T) {
		const taskCount = 500 // Reasonable for integration test

		start := time.Now()

		// Add many tasks
		for i := range taskCount {
			_, err := service.AddTask(fmt.Sprintf("Performance test task %d", i+1))
			if err != nil {
				t.Fatalf("AddTask() %d failed: %v", i+1, err)
			}
		}

		addDuration := time.Since(start)
		t.Logf("Added %d tasks in %v (%.2f tasks/second)",
			taskCount, addDuration, float64(taskCount)/addDuration.Seconds())

		// Test list performance
		start = time.Now()
		tasks, err := service.ListTasks("")
		if err != nil {
			t.Fatalf("ListTasks() failed: %v", err)
		}
		listDuration := time.Since(start)

		if len(tasks) != taskCount {
			t.Errorf("Expected %d tasks, got %d", taskCount, len(tasks))
		}

		t.Logf("Listed %d tasks in %v", taskCount, listDuration)

		// Test filtered list performance
		start = time.Now()
		todoTasks, err := service.ListTasks("todo")
		if err != nil {
			t.Fatalf("ListTasks(todo) failed: %v", err)
		}
		filterDuration := time.Since(start)

		if len(todoTasks) != taskCount {
			t.Errorf("Expected %d todo tasks, got %d", taskCount, len(todoTasks))
		}

		t.Logf("Filtered %d tasks in %v", taskCount, filterDuration)

		// Performance assertions (generous limits for CI)
		if addDuration > 10*time.Second {
			t.Errorf("Adding %d tasks took too long: %v", taskCount, addDuration)
		}
		if listDuration > 500*time.Millisecond {
			t.Errorf("Listing %d tasks took too long: %v", taskCount, listDuration)
		}
		if filterDuration > 500*time.Millisecond {
			t.Errorf("Filtering %d tasks took too long: %v", taskCount, filterDuration)
		}
	})
}

// TestRealWorldScenarios tests common usage patterns
func TestRealWorldScenarios(t *testing.T) {
	tmpFile := "real_world_test_tasks.json"
	defer os.Remove(tmpFile)

	service := NewTaskService(NewFileTaskRepository(tmpFile))

	t.Run("daily task management workflow", func(t *testing.T) {
		// Morning: Add today's tasks
		morningTasks := []string{
			"Check emails",
			"Team standup meeting",
			"Review pull requests",
			"Work on feature X",
			"Update documentation",
		}

		var taskIDs []int
		for _, desc := range morningTasks {
			task, err := service.AddTask(desc)
			if err != nil {
				t.Fatalf("AddTask(%s) failed: %v", desc, err)
			}
			taskIDs = append(taskIDs, task.ID)
		}

		// Start with first task
		err := service.MarkTaskInProgress(taskIDs[0])
		if err != nil {
			t.Fatalf("MarkTaskInProgress() failed: %v", err)
		}

		// Complete first task, start second
		err = service.MarkTaskDone(taskIDs[0])
		if err != nil {
			t.Fatalf("MarkTaskDone() failed: %v", err)
		}

		err = service.MarkTaskInProgress(taskIDs[1])
		if err != nil {
			t.Fatalf("MarkTaskInProgress() failed: %v", err)
		}

		// Midday: Reprioritize - update a task description
		err = service.UpdateTask(taskIDs[3], "Work on feature X - implement user authentication")
		if err != nil {
			t.Fatalf("UpdateTask() failed: %v", err)
		}

		// Complete meeting
		err = service.MarkTaskDone(taskIDs[1])
		if err != nil {
			t.Fatalf("MarkTaskDone() failed: %v", err)
		}

		// Afternoon: Work on main task
		err = service.MarkTaskInProgress(taskIDs[3])
		if err != nil {
			t.Fatalf("MarkTaskInProgress() failed: %v", err)
		}

		// End of day: Check progress
		inProgressTasks, err := service.ListTasks("in-progress")
		if err != nil {
			t.Fatalf("ListTasks(in-progress) failed: %v", err)
		}
		if len(inProgressTasks) != 1 {
			t.Errorf("Should have 1 in-progress task, got %d", len(inProgressTasks))
		}

		doneTasks, err := service.ListTasks("done")
		if err != nil {
			t.Fatalf("ListTasks(done) failed: %v", err)
		}
		if len(doneTasks) != 2 {
			t.Errorf("Should have 2 done tasks, got %d", len(doneTasks))
		}

		todoTasks, err := service.ListTasks("todo")
		if err != nil {
			t.Fatalf("ListTasks(todo) failed: %v", err)
		}
		if len(todoTasks) != 2 {
			t.Errorf("Should have 2 todo tasks, got %d", len(todoTasks))
		}

		// Verify updated task description
		allTasks, err := service.ListTasks("")
		if err != nil {
			t.Fatalf("ListTasks() failed: %v", err)
		}

		var updatedTask *Task
		for _, task := range allTasks {
			if task.ID == taskIDs[3] {
				updatedTask = &task
				break
			}
		}

		if updatedTask == nil {
			t.Fatalf("Updated task not found")
		}

		expectedDesc := "Work on feature X - implement user authentication"
		if updatedTask.Description != expectedDesc {
			t.Errorf("Task description = %v, want %v", updatedTask.Description, expectedDesc)
		}
	})

	t.Run("project management scenario", func(t *testing.T) {
		// Setup project tasks
		projectTasks := []string{
			"Requirements gathering",
			"System design",
			"Database schema",
			"API implementation",
			"Frontend development",
			"Testing",
			"Deployment",
		}

		var projectTaskIDs []int
		for _, desc := range projectTasks {
			task, err := service.AddTask(desc)
			if err != nil {
				t.Fatalf("AddTask(%s) failed: %v", desc, err)
			}
			projectTaskIDs = append(projectTaskIDs, task.ID)
		}

		// Work through project phases
		phases := [][]int{
			{0, 1},    // Planning phase
			{2, 3},    // Development phase
			{4, 5, 6}, // Implementation phase
		}

		for phaseNum, taskIndices := range phases {
			// Start all tasks in phase
			for _, idx := range taskIndices {
				err := service.MarkTaskInProgress(projectTaskIDs[idx])
				if err != nil {
					t.Fatalf("Phase %d: MarkTaskInProgress() failed: %v", phaseNum, err)
				}
			}

			// Complete all tasks in phase
			for _, idx := range taskIndices {
				err := service.MarkTaskDone(projectTaskIDs[idx])
				if err != nil {
					t.Fatalf("Phase %d: MarkTaskDone() failed: %v", phaseNum, err)
				}
			}
		}

		// Verify project completion
		doneTasks, err := service.ListTasks("done")
		if err != nil {
			t.Fatalf("ListTasks(done) failed: %v", err)
		}

		// Count project tasks (excluding daily tasks from previous test)
		projectDoneCount := 0
		for _, task := range doneTasks {
			if slices.Contains(projectTasks, task.Description) {
				projectDoneCount++
			}
		}

		if projectDoneCount != len(projectTasks) {
			t.Errorf(
				"Should have %d project tasks done, got %d",
				len(projectTasks),
				projectDoneCount,
			)
		}
	})
}

package main

import (
	"fmt"
	"strconv"
	"strings"
)

// CLI Interface (Presentation Layer)
type CLI struct {
	service *TaskService
}

func NewCLI(service *TaskService) *CLI {
	return &CLI{service: service}
}

func (c *CLI) Run(args []string) {
	if len(args) < 2 {
		c.printUsage()
		return
	}

	command := args[1]

	switch command {
	case "add":
		c.handleAdd(args[2:])
	case "update":
		c.handleUpdate(args[2:])
	case "delete":
		c.handleDelete(args[2:])
	case "mark-in-progress":
		c.handleMarkInProgress(args[2:])
	case "mark-done":
		c.handleMarkDone(args[2:])
	case "list":
		c.handleList(args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		c.printUsage()
	}
}

func (c *CLI) handleAdd(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: Description is required")
		fmt.Println("Usage: task-cli add \"Task description\"")
		return
	}

	description := args[0]
	task, err := c.service.AddTask(description)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Printf("Task added successfully (ID: %d)\n", task.ID)
}

func (c *CLI) handleUpdate(args []string) {
	if len(args) < 2 {
		fmt.Println("Error: ID and description are required")
		fmt.Println("Usage: task-cli update <id> \"New description\"")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Error: Invalid task ID")
		return
	}

	description := args[1]
	err = c.service.UpdateTask(id, description)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Println("Task updated successfully")
}

func (c *CLI) handleDelete(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: ID is required")
		fmt.Println("Usage: task-cli delete <id>")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Error: Invalid task ID")
		return
	}

	err = c.service.DeleteTask(id)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Println("Task deleted successfully")
}

func (c *CLI) handleMarkInProgress(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: ID is required")
		fmt.Println("Usage: task-cli mark-in-progress <id>")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Error: Invalid task ID")
		return
	}

	err = c.service.MarkTaskInProgress(id)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Println("Task marked as in progress")
}

func (c *CLI) handleMarkDone(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: ID is required")
		fmt.Println("Usage: task-cli mark-done <id>")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Error: Invalid task ID")
		return
	}

	err = c.service.MarkTaskDone(id)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Println("Task marked as done")
}

func (c *CLI) handleList(args []string) {
	var status string
	if len(args) > 0 {
		status = args[0]
		// Validate status
		if status != "todo" && status != "in-progress" && status != "done" {
			fmt.Printf(
				"Error: Invalid status '%s'. Valid options: todo, in-progress, done\n",
				status,
			)
			return
		}
	}

	tasks, err := c.service.ListTasks(status)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	if len(tasks) == 0 {
		if status == "" {
			fmt.Println("No tasks found")
		} else {
			fmt.Printf("No tasks with status '%s' found\n", status)
		}
		return
	}

	c.printTasks(tasks)
}

func (c *CLI) printTasks(tasks []Task) {
	fmt.Println("Tasks:")
	fmt.Println("------")
	for _, task := range tasks {
		statusDisplay := strings.ToUpper(string(task.Status))
		fmt.Printf("ID: %d | Status: %s | Description: %s\n",
			task.ID, statusDisplay, task.Description)
		fmt.Printf("Created: %s | Updated: %s\n",
			task.CreatedAt.Format("2006-01-02 15:04:05"),
			task.UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("------")
	}
}

func (c *CLI) printUsage() {
	fmt.Println("Task Tracker CLI")
	fmt.Println("Usage:")
	fmt.Println("  task-cli add \"Task description\"")
	fmt.Println("  task-cli update <id> \"New description\"")
	fmt.Println("  task-cli delete <id>")
	fmt.Println("  task-cli mark-in-progress <id>")
	fmt.Println("  task-cli mark-done <id>")
	fmt.Println("  task-cli list [status]")
	fmt.Println("")
	fmt.Println("Status options for list command:")
	fmt.Println("  todo, in-progress, done")
}

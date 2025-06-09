package main

import "os"

// Main function - Application entry point
func main() {
	// Dependency injection
	repo := NewFileTaskRepository("tasks.json")
	service := NewTaskService(repo)
	cli := NewCLI(service)

	// Handle the case where no arguments are provided
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}

	// Run the CLI
	cli.Run(os.Args)
}

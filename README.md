# Task Tracker CLI

A simple and efficient command-line tool to manage your daily tasks. Built with Go for speed and reliability.

## Quick Start

```bash
# Get the code
git clone https://github.com/alnah/task-tracker
cd task-tracker

# Build and try it
make demo
```

## Installation

```bash
# Build locally
make build

# Or install system-wide
make install
```

## How to Use

### Basic Commands

```bash
# Add a new task
./task-cli add "Buy groceries"

# See all your tasks
./task-cli list

# Start working on a task
./task-cli mark-in-progress 1

# Complete a task
./task-cli mark-done 1

# Update a task description
./task-cli update 1 "Buy groceries and cook dinner"

# Remove a task
./task-cli delete 1
```

### Filter by Status

```bash
# See only pending tasks
./task-cli list todo

# See work in progress
./task-cli list in-progress

# See completed tasks
./task-cli list done
```

## Examples

### Daily Workflow

```bash
# Morning: Add today's tasks
./task-cli add "Review emails"
./task-cli add "Team meeting at 10am"
./task-cli add "Finish quarterly report"

# Start working
./task-cli mark-in-progress 1

# Complete tasks as you go
./task-cli mark-done 1
./task-cli mark-in-progress 2

# Check what's left
./task-cli list todo
```

### Project Management

```bash
# Add project tasks
./task-cli add "Design database schema"
./task-cli add "Implement user authentication"
./task-cli add "Write API documentation"
./task-cli add "Deploy to staging"

# Track progress
./task-cli list in-progress
./task-cli list done
```

## Project Structure

```
task-tracker/
├── main.go           # Application entry point
├── model.go          # Task data structure
├── logic.go          # Task operations
├── repository.go     # File storage
├── application.go    # Business logic
├── cli.go           # Command-line interface
├── main_test.go     # All tests
├── go.mod           # Go dependencies
├── Makefile         # Build commands
└── README.md        # This file
```

## Development

### Available Commands

```bash
make build      # Build the application
make test       # Run tests
make demo       # Try the application
make clean      # Clean up files
make install    # Install globally
make dev        # Development workflow
```

### Testing

The project includes comprehensive tests for all functionality:

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage
```

## Features

- ✅ **Simple**: Easy to learn and use
- ⚡ **Fast**: Built with Go for performance
- 💾 **Persistent**: Saves tasks to a JSON file
- 🔒 **Reliable**: Comprehensive error handling
- 🧪 **Tested**: Full test coverage
- 📱 **Portable**: Single binary, works anywhere

## How It Works

Your tasks are stored in a `tasks.json` file in the current directory. Each task has:

- **ID**: Unique number (auto-generated)
- **Description**: What you need to do
- **Status**: `todo`, `in-progress`, or `done`
- **Timestamps**: When created and last updated

## Requirements

- Go 1.19 or later
- No external dependencies

## Contributing

1. Fork the repository
2. Make your changes
3. Run `make test` to ensure tests pass
4. Submit a pull request

## License

MIT License - feel free to use this project for learning or in your own work.

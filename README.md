# Task Tracker CLI

A simple and efficient command-line tool to manage your daily tasks.
Built with Go using Clean Architecture principles for speed and reliability for my own educational purpose.

## Quick Start

```bash
# Get the code
git clone https://github.com/alnah/task-tracker
cd task-tracker

# Build and try it out
make demo
```

## Installation

```bash
# Build locally
make build

# Or install system-wide (requires sudo)
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
├── model.go          # Task data structure and domain errors
├── logic.go          # Domain business logic
├── repository.go     # Data persistence layer
├── application.go    # Application services (use cases)
├── cli.go           # Command-line interface
├── main_test.go     # Comprehensive test suite
├── go.mod           # Go module definition
├── Makefile         # Build automation
└── README.md        # This file
```

## Development

### Available Make Commands

```bash
make help           # Show all available commands
make build          # Build the application
make run            # Build and run with help command
make test           # Run all tests
make test-coverage  # Run tests with coverage report
make demo           # Run a complete interactive demo
make clean          # Clean up generated files
make install        # Install globally (requires sudo)
make dev            # Development workflow (clean, test, build)
make check          # Run tests and build verification
make quick-start    # Build and demo for new users
```

### Testing

The project includes comprehensive tests for all layers:

```bash
# Run all tests with verbose output
make test

# Generate coverage report (creates coverage.html)
make test-coverage

# Development workflow with all checks
make dev
```

### Development Workflow

```bash
# For active development
make dev        # Runs clean, test, and build

# Quick verification
make check      # Runs test and build

# New user experience
make quick-start # Runs build and demo
```

## Architecture

This project demonstrates **Clean Architecture** principles:

- **Domain Layer** (`model.go`, `logic.go`): Core business rules and entities
- **Application Layer** (`application.go`): Use cases and business workflows
- **Infrastructure Layer** (`repository.go`): Data persistence
- **Presentation Layer** (`cli.go`): User interface

## Features

- ✅ **Simple**: Easy to learn and use
- ⚡ **Fast**: Built with Go for performance
- 💾 **Persistent**: Saves tasks to a JSON file
- 🔒 **Reliable**: Comprehensive error handling and validation
- 🧪 **Tested**: Full test coverage with multiple test types
- 📱 **Portable**: Single binary, works anywhere
- 🏗️ **Clean Architecture**: Well-structured, maintainable code
- 🔧 **Developer Friendly**: Rich Makefile with helpful commands

## How It Works

Your tasks are stored in a `tasks.json` file in the current directory. Each task has:

- **ID**: Unique number (auto-generated)
- **Description**: What you need to do (validated, trimmed)
- **Status**: `todo`, `in-progress`, or `done`
- **Timestamps**: When created and last updated

The application follows domain-driven design with proper separation of concerns:

- Domain entities with business rules
- Repository pattern for data access
- Application services for use cases
- Clean CLI interface

## Requirements

- Go 1.24.2 or later
- No external dependencies (pure Go standard library)

### Testing Guidelines

- Write tests for all new features
- Maintain test coverage above 80%
- Use table-driven tests for multiple scenarios
- Include integration tests for workflows

## License

MIT License - feel free to use this project for learning or in your own work.

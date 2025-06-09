# Task Tracker CLI - Makefile

.PHONY: build test clean help run install demo

# Default target
help:
	@echo "Task Tracker CLI - Available commands:"
	@echo ""
	@echo "  build     Build the application"
	@echo "  run       Build and run with example command"
	@echo "  test      Run all tests"
	@echo "  demo      Run a complete demo"
	@echo "  clean     Clean up generated files"
	@echo "  install   Install globally (requires sudo)"
	@echo ""

# Build the application
build:
	@echo "Building task-cli..."
	@go build -o task-cli .
	@echo "✅ Build complete: ./task-cli"

# Quick run with help
run: build
	@echo "Running task-cli..."
	@./task-cli

# Run tests
test:
	@echo "Running tests..."
	@go test -v
	@echo "✅ All tests passed"

# Test with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

# Clean up
clean:
	@echo "Cleaning up..."
	@rm -f task-cli coverage.out coverage.html tasks.json test_*.json
	@go clean
	@echo "✅ Cleanup complete"

# Install globally
install: build
	@echo "Installing task-cli globally..."
	@sudo cp task-cli /usr/local/bin/
	@echo "✅ Installed! You can now use 'task-cli' from anywhere"

# Complete demo
demo: build
	@echo "=== Task Tracker Demo ==="
	@echo ""
	@echo "📝 Adding some tasks..."
	@./task-cli add "Buy groceries"
	@./task-cli add "Write project report"
	@./task-cli add "Call mom"
	@echo ""
	@echo "📋 Current tasks:"
	@./task-cli list
	@echo ""
	@echo "🚀 Starting work on first task..."
	@./task-cli mark-in-progress 1
	@echo ""
	@echo "✅ Completed second task..."
	@./task-cli mark-done 2
	@echo ""
	@echo "📝 Updating third task..."
	@./task-cli update 3 "Call mom and discuss weekend plans"
	@echo ""
	@echo "📊 Final status:"
	@./task-cli list
	@echo ""
	@echo "🎯 Tasks by status:"
	@echo "Todo tasks:"
	@./task-cli list todo
	@echo "In progress:"
	@./task-cli list in-progress
	@echo "Completed:"
	@./task-cli list done
	@echo ""
	@echo "=== Demo complete! ==="

# Development helpers
dev: clean test build
	@echo "✅ Development build ready"

check: test build
	@echo "✅ All checks passed"

# Quick start for new users
quick-start: build demo
	@echo ""
	@echo "🎉 Quick start complete!"
	@echo "Try these commands:"
	@echo "  ./task-cli add \"Your task\""
	@echo "  ./task-cli list"
	@echo "  ./task-cli help"

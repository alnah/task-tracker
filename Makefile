# Task Tracker CLI 

.PHONY: build test clean help run install demo
.PHONY: test-domain test-repository test-application test-integration test-fast test-slow
.PHONY: test-coverage

# Default target
help:
	@echo "Task Tracker CLI - Available commands:"
	@echo ""
	@echo "Build commands:"
	@echo "  build           Build the application"
	@echo "  run             Build and run with example command"
	@echo "  clean           Clean up generated files"
	@echo "  install         Install globally (requires sudo)"
	@echo ""
	@echo "Test commands:"
	@echo "  test            Run all tests"
	@echo "  test-fast       Run fast tests only (domain + application)"
	@echo "  test-slow       Run slow tests only (repository + integration)"
	@echo "  test-domain     Run domain logic tests"
	@echo "  test-repository Run repository tests"
	@echo "  test-application Run application service tests"
	@echo "  test-integration Run integration tests"
	@echo ""
	@echo "Coverage commands:"
	@echo "  test-coverage   Run tests with coverage report"
	@echo ""
	@echo "Demo commands:"
	@echo "  demo            Run a complete demo"
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


# Domain tests - fastest, pure business logic
test-domain:
	@echo "🧠 Testing domain logic..."
	@go test -v -run "TestNewTask|TestTask_" -count=1
	@echo "✅ Domain tests passed"

# Repository tests - database/file operations
test-repository:
	@echo "💾 Testing repository layer..."
	@go test -v -run "TestFileTaskRepository|TestMockRepository" -count=1
	@echo "✅ Repository tests passed"

# Application tests - service orchestration
test-application:
	@echo "⚙️  Testing application services..."
	@go test -v -run "TestTaskService" -count=1
	@echo "✅ Application tests passed"

# Integration tests - full stack
test-integration:
	@echo "🔗 Testing integration scenarios..."
	@go test -v -run "TestFull|TestConcurrent|TestData|TestPerformance|TestRealWorld" -count=1
	@echo "✅ Integration tests passed"

# Fast tests - no I/O operations
test-fast: test-domain test-application
	@echo "⚡ Fast test suite completed"

# Slow tests - involve I/O
test-slow: test-repository test-integration
	@echo "🐌 Slow test suite completed"

# All tests 
test:
	@echo "🧪 Running all tests..."
	@go test -v -count=1
	@echo "🎯 All tests completed successfully"

# Coverage analysis
test-coverage:
	@echo "📊 Running tests with coverage..."
	@go test -cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"
	@go tool cover -func=coverage.out | grep total:

# Test with specific patterns
test-pattern:
	@echo "🔍 Running tests matching pattern: $(PATTERN)"
	@go test -v -run "$(PATTERN)"

# Benchmark tests
benchmark:
	@echo "🏃 Running benchmarks..."
	@go test -bench=. -benchmem

# Race condition detection
test-race:
	@echo "🏁 Testing for race conditions..."
	@go test -race

# Clean up
clean:
	@echo "🧹 Cleaning up..."
	@rm -f task-cli
	@rm -f coverage.out coverage.html
	@rm -f tasks.json test_*.json
	@rm -f *_test_tasks.json
	@go clean
	@echo "✅ Cleanup complete"

# Install globally
install: build
	@echo "📦 Installing task-cli globally..."
	@sudo cp task-cli /usr/local/bin/
	@echo "✅ Installed! You can now use 'task-cli' from anywhere"

# Development workflows

# Quick development check
dev: clean test-fast build
	@echo "✅ Development check passed"

# Full verification before commit
verify: clean test test-coverage build
	@echo "✅ Full verification completed"

# Demo functionality
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

# Test-driven development helpers
tdd-domain:
	@echo "🔄 TDD Mode: Domain layer"
	@echo "Press Ctrl+C to stop..."
	@while true; do \
		make test-domain; \
		echo ""; \
		echo "Waiting for changes... (press Enter to run again)"; \
		read; \
	done


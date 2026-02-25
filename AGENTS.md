# AGENTS.md - Guide for AI Coding Agents

This document provides guidelines for AI coding agents (including Claude, GPT, Copilot, etc.) working in this repository.

## Repository Overview

This is the **awesome-skills** repository - a collection of skills, tools, or utilities (specific purpose TBD based on project evolution).

**Language**: Go (primary, based on .gitignore configuration)

## Build, Test, and Lint Commands

### Go Project Commands

```bash
# Build the project
go build ./...

# Build specific package
go build ./path/to/package

# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a single test file
go test ./path/to/package/file_test.go

# Run a specific test by name
go test -run TestFunctionName ./path/to/package

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Run tests with race detection
go test -race ./...

# Format code
go fmt ./...

# Run linter (if golangci-lint is installed)
golangci-lint run

# Vet code for suspicious constructs
go vet ./...

# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify
```

### Common Development Workflow

```bash
# Before committing
go fmt ./...
go vet ./...
go test ./...

# For comprehensive checks
golangci-lint run && go test -race -cover ./...
```

## Code Style Guidelines

### General Principles

1. **Follow Go conventions**: Use `gofmt` and standard Go idioms
2. **Simplicity over cleverness**: Write clear, maintainable code
3. **Error handling first**: Always handle errors explicitly
4. **Documentation**: Public APIs must have godoc comments

### Import Formatting

```go
import (
    // Standard library first
    "context"
    "fmt"
    "net/http"

    // External dependencies second
    "github.com/pkg/errors"
    "github.com/stretchr/testify/assert"

    // Internal packages last
    "github.com/pangu-studio/awesome-skills/internal/pkg"
)
```

- Use `goimports` for automatic import organization
- Group imports by category with blank lines between groups
- Avoid dot imports except in tests when appropriate

### Naming Conventions

#### Variables and Functions
- Use **camelCase** for private: `myVariable`, `helperFunction()`
- Use **PascalCase** for public: `PublicAPI`, `ExportedFunction()`
- Use **short, descriptive names** in small scopes: `i`, `err`, `ctx`
- Use **longer, descriptive names** in larger scopes: `requestContext`, `userRepository`

#### Constants
```go
const (
    MaxRetries = 3
    DefaultTimeout = 30 * time.Second
)
```

#### Interfaces
- Single-method interfaces: name with `-er` suffix: `Reader`, `Writer`, `Closer`
- Multi-method interfaces: descriptive name: `UserRepository`, `ConfigManager`

#### Files
- Use **snake_case** for filenames: `user_repository.go`, `http_server.go`
- Test files: `*_test.go`
- Use meaningful package names that match directory names

### Types and Interfaces

```go
// Good: Document exported types
// UserService handles user-related operations.
type UserService struct {
    repo Repository
    log  Logger
}

// Good: Small, focused interfaces
type Repository interface {
    Get(ctx context.Context, id string) (*User, error)
    Save(ctx context.Context, user *User) error
}

// Good: Pointer receivers for methods that modify state
func (s *UserService) UpdateUser(ctx context.Context, user *User) error {
    // Implementation
}

// Good: Value receivers for methods that don't modify state
func (u User) FullName() string {
    return u.FirstName + " " + u.LastName
}
```

### Error Handling

```go
// Good: Always check errors immediately
result, err := someOperation()
if err != nil {
    return fmt.Errorf("failed to perform operation: %w", err)
}

// Good: Wrap errors with context using %w
if err := saveData(data); err != nil {
    return fmt.Errorf("save data for user %s: %w", userID, err)
}

// Good: Use errors.Is and errors.As for error checking
if errors.Is(err, ErrNotFound) {
    return handleNotFound()
}

// Avoid: Ignoring errors
someOperation() // Bad

// Avoid: Generic error messages
return errors.New("error") // Bad
```

### Function Design

- **Keep functions small**: Ideally under 50 lines
- **Single responsibility**: Each function should do one thing well
- **Context first**: Context should be the first parameter
- **Return errors**: Return error as the last return value

```go
// Good
func ProcessUser(ctx context.Context, userID string) (*User, error) {
    // Implementation
}
```

### Comments and Documentation

```go
// Package-level documentation
// Package skills provides utilities for managing skills.
package skills

// Exported function documentation - describe what it does, parameters, and return values
// GetUserSkills retrieves all skills for a given user.
// Returns ErrNotFound if the user doesn't exist.
func GetUserSkills(ctx context.Context, userID string) ([]Skill, error) {
    // Implementation comments for complex logic
}
```

## Testing Guidelines

### Test File Structure

```go
package skills_test // Use package_test for black-box testing

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "github.com/pangu-studio/awesome-skills/skills"
)

func TestProcessUser(t *testing.T) {
    // Arrange
    user := &skills.User{ID: "123"}
    
    // Act
    result, err := skills.ProcessUser(context.Background(), user)
    
    // Assert
    require.NoError(t, err)
    assert.NotNil(t, result)
}
```

### Test Naming

- Use `Test` prefix: `TestFunctionName`
- Use table-driven tests for multiple cases
- Use subtests with `t.Run()` for clarity

## Git Commit Guidelines

- Use conventional commits: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`
- Keep commits atomic and focused
- Write clear, descriptive commit messages

## Dependencies

- Minimize external dependencies
- Use Go modules (`go.mod` and `go.sum`)
- Keep dependencies up to date
- Vendor dependencies if needed for reproducibility

## Performance Considerations

- Avoid premature optimization
- Use benchmarks to measure performance: `go test -bench=.`
- Use profiling tools when needed: `pprof`
- Consider memory allocations in hot paths

## Security

- Never commit secrets, API keys, or credentials
- Use environment variables or secret management for sensitive data
- Validate all inputs
- Use context for timeouts and cancellation

---

**Last Updated**: February 25, 2026
**Maintained By**: Pangu Studio Team

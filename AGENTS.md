# Agent Development Guide for aliyun-acme-hook

## Project Overview
This is a Go-based application that serves as an ACME (Automatic Certificate Management Environment) hook for Alibaba Cloud services, enabling automated SSL/TLS certificate management for domains using Alibaba Cloud's CDN, SLB (Server Load Balancer), and CAS (Certificate Authority Service).

## Build Commands
```bash
# Build the application
make app
# or directly with Go
go build -o dist/aliyun-acme-hook -ldflags="-s -w -extldflags \"-static\"" ./cmd/app

# Full build and install
make all

# Clean build artifacts
make clean
```

## Test Commands
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./internal/service/...
go test ./config/...

# Run a specific test (if tests existed - none currently in codebase)
# go test -run ^TestFunctionName$ ./path/to/package
```

## Lint Commands
```bash
# Standard Go tools
go fmt ./...
go vet ./...
gofmt -s -w .
golint ./...  # if golint is installed

# Using golangci-lint (install separately)
golangci-lint run
```

## Code Style Guidelines

### Imports
- Group imports with blank lines between groups
- Standard library imports first
- Third-party packages second
- Local packages third
- Use explicit imports, avoid dot imports
- Example:
```go
import (
    "context"
    "log"

    "github.com/urfave/cli/v2"

    "github.com/bububa/aliyun-acme-hook/config"
    "github.com/bububa/aliyun-acme-hook/internal/service"
)
```

### Formatting
- Use gofmt for consistent formatting
- 4-space indents (no tabs)
- Line length should be reasonable (typically < 120 chars)
- Struct tags should use camelCase for JSON/YAML
- Consistent naming conventions throughout

### Types
- Use descriptive type names
- Export types that are meant to be used outside the package
- Use pointer receivers for methods that modify the receiver
- Use value receivers for small structs and methods that don't modify the receiver

### Naming Conventions
- Use PascalCase for exported functions/types/variables
- Use camelCase for unexported functions/types/variables
- Use descriptive names (avoid single-letter variables except for loop counters)
- Use standard Go naming patterns (e.g., Marshal/Unmarshal, NewConstructor)

### Error Handling
- Always check and handle errors appropriately
- Use error wrapping with %w when returning errors
- Log errors when appropriate but don't over-log
- Return early on error conditions when possible

### Testing
- Test files should be named *_test.go
- Test functions should be named TestFunctionName for unit tests
- Use table-driven tests for multiple test cases
- Include benchmarks for performance-critical code when applicable
- Note: Current codebase has no test files yet

## Project Structure
- `cmd/app/main.go`: Application entry point
- `internal/app/`: Main application logic and CLI setup
- `internal/service/`: Core business logic services
- `internal/config/`: Configuration loading and parsing
- `internal/cas/`: Alibaba Cloud CAS integration
- `internal/cdn/`: Alibaba Cloud CDN integration
- `internal/slb/`: Alibaba Cloud SLB integration
- `internal/model/`: Data models and structures

## Dependencies
- github.com/urfave/cli/v2: CLI application framework
- github.com/alibabacloud-go/*: Alibaba Cloud SDK modules
- github.com/jinzhu/configor: Configuration loading
- Standard library for core functionality

## Environment Setup
- Go 1.25.0+ required (as specified in go.mod)
- Access to Alibaba Cloud account with appropriate permissions
- Configuration file at /etc/aliyun-acme-hook.toml or specified path

## Deployment
- Built binary is installed to /usr/local/bin/aliyun-acme-hook
- Requires TOML configuration file with Alibaba Cloud credentials
- Designed for use with acme.sh for automated certificate management
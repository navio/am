# Contributing to Alias Manager

Thank you for your interest in contributing to Alias Manager! This document provides guidelines and instructions for contributing.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/am.git
   cd am
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/navio/am.git
   ```

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git

### Installing Dependencies

```bash
go mod download
```

### Building

```bash
go build -o am
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run tests with verbose output
go test ./... -v

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Making Changes

1. **Create a branch** for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our coding standards (see below)

3. **Write tests** for your changes

4. **Run tests** to ensure everything passes:
   ```bash
   go test ./...
   ```

5. **Format your code**:
   ```bash
   go fmt ./...
   ```

6. **Commit your changes** with a clear message:
   ```bash
   git commit -m "Add feature: description of your changes"
   ```

## Coding Standards

- Follow standard Go conventions and idioms
- Use `gofmt` for formatting
- Write clear, self-documenting code
- Add comments for exported functions and types
- Keep functions focused and reasonably sized
- Handle errors explicitly

## Testing Guidelines

- Write tests for all new functionality
- Maintain or improve code coverage (aim for >80%)
- Test edge cases and error conditions
- Use table-driven tests where appropriate
- Name tests clearly: `TestFunctionName_Scenario`

## Pull Request Process

1. **Update documentation** if needed (README.md, code comments)

2. **Ensure all tests pass**:
   ```bash
   go test ./...
   ```

3. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

4. **Create a Pull Request** on GitHub:
   - Use a clear, descriptive title
   - Describe what changes you made and why
   - Reference any related issues
   - Include screenshots for UI changes (if applicable)

5. **Respond to feedback** from reviewers

6. **Wait for approval** and merge

## Code Review

All submissions require review. We use GitHub pull requests for this purpose. Your PR will be reviewed for:

- Code quality and style
- Test coverage
- Documentation
- Backward compatibility
- Performance implications

## Types of Contributions

### Bug Reports

- Use the GitHub issue tracker
- Describe the bug clearly
- Include steps to reproduce
- Specify your environment (OS, shell, Go version)

### Feature Requests

- Use the GitHub issue tracker
- Explain the feature and its benefits
- Discuss implementation approach if you have ideas

### Code Contributions

- Bug fixes
- New features
- Performance improvements
- Documentation improvements
- Test improvements

## Project Structure

```
am/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command and setup
│   ├── add.go             # Add command
│   ├── list.go            # List command
│   ├── delete.go          # Delete command
│   └── update.go          # Update command
├── internal/
│   ├── alias/             # Core alias management
│   │   ├── models.go      # Data structures
│   │   ├── parser.go      # Parse alias definitions
│   │   ├── manager.go     # CRUD operations
│   │   └── *_test.go      # Tests
│   ├── dotfile/           # File operations
│   │   ├── handler.go     # Read/write with backups
│   │   └── *_test.go      # Tests
│   └── shell/             # Shell detection
│       ├── detector.go    # Detect shell type
│       └── *_test.go      # Tests
└── main.go                # Entry point
```

## Adding a New Feature

1. **Discuss** the feature in an issue first
2. **Design** the implementation
3. **Implement** in a feature branch
4. **Test** thoroughly
5. **Document** in code and README
6. **Submit** a pull request

## Release Process

(For maintainers)

1. Update version in `cmd/root.go`
2. Update CHANGELOG.md
3. Create a git tag: `git tag v1.x.x`
4. Push tag: `git push origin v1.x.x`
5. Create GitHub release
6. Update Homebrew formula if needed

## Questions?

Feel free to ask questions by:
- Opening an issue
- Commenting on an existing issue or PR
- Reaching out to the maintainers

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

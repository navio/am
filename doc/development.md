---
title: Development
description: Build, test, and contribute to Alias Manager.
---

# Development

## Prerequisites

- Go 1.21 or later

## Build and test

```bash
go mod download
go build -o am
go test ./...
```

For coverage:

```bash
go test ./... -cover
```

## Project structure

```
am/
├── cmd/          # CLI commands
├── internal/     # alias logic, dotfile I/O, shell detection
└── main.go       # entry point
```

## Contributing

1. Fork and branch (`git checkout -b feature/your-topic`).
2. Make changes and add tests where useful.
3. Run `go test ./...` before opening a PR.
4. Submit a PR to `main` on GitHub.

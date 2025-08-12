# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MT To Hatena MD (mttohmd) is a Go utility that converts MovableType (MT) format backup files from Hatena Blog into individual Markdown files for each blog entry. The tool processes a single MT format text file and outputs separate Hatena Blog compatible Markdown files for each entry.

## Development Commands

### Build and Run
```bash
go build -o mttohmd
./mttohmd
```

### Run directly
```bash
go run main.go
```

### Testing
```bash
go test ./...
```

### Code formatting and linting
```bash
go fmt ./...
go vet ./...
```

## Architecture

The project is currently in early development with a minimal structure:

- `main.go` - Entry point containing the main application logic
- `go.mod` - Go module definition (Go 1.24.5)

The application will need to:
1. Parse MovableType format text files
2. Extract individual blog entries with metadata
3. Convert content to Hatena Blog Markdown format
4. Output separate files for each entry

## Language and Formatting

This is a Japanese language project. Comments and documentation should be in Japanese when appropriate for consistency with the existing codebase.
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
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -v -cover ./...

# Run specific package tests
go test ./entry
go test ./converter  
go test ./generator
```

### Code formatting and linting
```bash
go fmt ./...
go vet ./...
```

## Architecture

The project is organized into several packages:

- `main.go` - Entry point containing the main application logic
- `entry/` - MovableType format file parsing and entry structure definition
- `converter/` - HTML to Markdown conversion with Hatena Blog specific features
- `generator/` - File name generation and MT format output utilities
- `go.mod` - Go module definition (Go 1.24.5)

### Package Structure

#### `entry` package
- `Entry` struct - Defines blog entry structure with metadata fields
- `ParseEntries()` - Parses MT format files and extracts individual entries

#### `converter` package
- `ToMarkdown()` - Converts entries to Hatena Blog compatible Markdown
- `convertHTMLToMarkdown()` - Converts HTML tags to Markdown syntax
- Supports Hatena Blog ASIN tags: `<div class="hatena-asin-detail">` → `[asin:ID:detail]`
- Optimized with precompiled regex patterns for performance

#### `generator` package
- `GenerateFilename()` - Creates appropriate filenames from entry metadata
- `GenerateMTContent()` - Outputs entries in original MT format

The application:
1. Parses MovableType format text files into structured entries
2. Processes each entry through HTML-to-Markdown conversion
3. Handles Hatena Blog specific elements (ASIN product links, etc.)
4. Outputs both MT format (.txt) and Markdown format (.md) files
5. Supports comprehensive test coverage for all conversion functions

## Features

### HTML to Markdown Conversion
- Basic HTML tags: `<p>`, `<br>`, `<strong>`, `<em>`, `<a>`, `<img>`
- Header tags: `<h1>` through `<h6>`
- Lists: `<ul>`, `<ol>`, `<li>`
- Blockquotes: `<blockquote>`

### Hatena Blog Specific Features
- **ASIN Product Links**: Converts Amazon product detail boxes to Hatena notation
  - Input: `<div class="hatena-asin-detail">...Amazon link...</div>`
  - Output: `[asin:PRODUCT_ID:detail]`
  - Supports complex nested HTML structures and p-tag wrapped patterns

### Performance Optimizations
- Precompiled regex patterns for all HTML conversions
- Efficient pattern matching without runtime compilation overhead
- Optimized for processing large numbers of blog entries

## Language and Formatting

This is a Japanese language project. Comments and documentation should be in Japanese when appropriate for consistency with the existing codebase.
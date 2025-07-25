# Sigma CLI

A command-line tool for reading and analyzing SIQ (SIG Pack) files.

## Installation

```bash
go install github.com/minmaxmean/sigma@latest
```

Or run directly:

```bash
go run sigma.go [command]
```

## Usage

### Basic Commands

```bash
# Show help
sigma --help

# Read and display information from a SIQ file
sigma read game.siq

# Show help for the read command
sigma read --help
```

### Examples

```bash
# Process a SIQ file
sigma read data/example.siq

# Get help
sigma --help
```

## Features

The Sigma CLI provides comprehensive information about SIQ files including:

- **Package Metadata**: Name, ID, version, difficulty, language, publisher, date
- **Statistics**: Number of rounds, themes, and questions
- **File Listing**: All files contained in the SIQ archive
- **Question Details**: Question types, right/wrong answers, and content

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o sigma sigma.go
```

## Project Structure

- `sigma.go` - Main CLI application using Cobra
- `siq/` - SIQ file handling package
- `docs/` - Documentation for SIQ file formats
- `examples/` - Example applications
- `data/` - Sample data files 
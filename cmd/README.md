# Command Structure

This directory contains the individual command implementations for the `sigma` CLI tool.

## Files

- `read.go` - Implements the `read` command for displaying SIQ file information
- `markdown.go` - Implements the `markdown` command for converting SIQ files to markdown format

## Structure

Each command file follows this pattern:

1. **Command Definition**: A `cobra.Command` variable with usage, description, and arguments
2. **Run Function**: The main execution logic for the command
3. **Helper Functions**: Any additional functions needed for the command
4. **Export Function**: A function to return the command for registration with the root command

## Adding New Commands

To add a new command:

1. Create a new file in this directory (e.g., `newcommand.go`)
2. Follow the pattern established in existing command files
3. Export a function that returns the command (e.g., `GetNewCommandCmd()`)
4. Register the command in `sigma.go` by adding it to the `init()` function

## Example

```go
package cmd

import (
    "github.com/spf13/cobra"
)

var newCommand = &cobra.Command{
    Use:   "newcommand [args]",
    Short: "Description of the new command",
    Args:  cobra.ExactArgs(1),
    Run:   runNewCommand,
}

func runNewCommand(cmd *cobra.Command, args []string) {
    // Command implementation
}

func GetNewCommandCmd() *cobra.Command {
    return newCommand
}
``` 
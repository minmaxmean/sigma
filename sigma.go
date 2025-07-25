package main

import (
	"fmt"
	"os"

	"github.com/minmaxmean/sigma/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sigma",
	Short: "A CLI tool for processing SIQ files",
	Long: `Sigma is a command-line tool for reading and analyzing SIQ (SIG Pack) files.
It provides detailed information about SIG packages including questions, answers, and metadata.`,
}

func init() {
	rootCmd.AddCommand(cmd.GetReadCmd())
	rootCmd.AddCommand(cmd.GetMarkdownCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

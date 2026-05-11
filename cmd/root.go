package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "1.2.1"

var rootCmd = &cobra.Command{
	Use:   "am",
	Short: "Alias Manager - Manage your shell aliases with ease",
	Long: `Alias Manager (am) is a lightweight CLI tool for managing shell aliases.
It works directly with your .zshrc or .bashrc files, providing a simple
interface for adding, listing, updating, and deleting aliases.`,
	Version: version,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(`{{.Version}}
`)
}

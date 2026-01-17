package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/navio/am/internal/alias"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add <name> <command>",
	Aliases: []string{"a"},
	Short:   "Add a new alias",
	Long: `Add a new alias to your shell configuration file.
The alias will be added to your .zshrc or .bashrc depending on your current shell.

Example:
  am add ll 'ls -la'
  am add gs 'git status'
  am add server 'python -m http.server 8000'`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		command := strings.Join(args[1:], " ")

		manager := alias.NewManager()
		if err := manager.Add(name, command); err != nil {
			return err
		}

		// Success message
		green := color.New(color.FgGreen, color.Bold)
		cyan := color.New(color.FgCyan)

		green.Print("✓ ")
		fmt.Print("Alias added successfully!\n\n")
		fmt.Printf("  %s → %s\n\n", cyan.Sprint(name), command)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

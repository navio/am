package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/navio/am/internal/alias"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <name> <new-command>",
	Short: "Update an existing alias",
	Long: `Update the command for an existing alias.

Example:
  am update ll 'ls -lah'
  am update gs 'git status -sb'`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		newCommand := strings.Join(args[1:], " ")

		manager := alias.NewManager()

		// Get old alias details (for display)
		oldAlias, err := manager.Get(name)
		if err != nil {
			return err
		}

		// Update the alias
		if err := manager.Update(name, newCommand); err != nil {
			return err
		}

		// Success message
		green := color.New(color.FgGreen, color.Bold)
		cyan := color.New(color.FgCyan)
		gray := color.New(color.FgHiBlack)

		green.Print("✓ ")
		fmt.Print("Alias updated successfully!\n\n")
		fmt.Printf("  %s\n", cyan.Sprint(name))
		fmt.Printf("    %s %s\n", gray.Sprint("Before:"), oldAlias.Command)
		fmt.Printf("    %s %s\n\n", gray.Sprint("After: "), newCommand)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

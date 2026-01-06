package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/navio/am/internal/alias"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete an alias",
	Long: `Delete an alias from your shell configuration file.

Example:
  am delete ll
  am delete gs`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		manager := alias.NewManager()

		// Get the alias details before deleting (for display)
		aliasToDelete, err := manager.Get(name)
		if err != nil {
			return err
		}

		// Delete the alias
		if err := manager.Delete(name); err != nil {
			return err
		}

		// Success message
		green := color.New(color.FgGreen, color.Bold)
		cyan := color.New(color.FgCyan)

		green.Print("✓ ")
		fmt.Print("Alias deleted successfully!\n\n")
		fmt.Printf("  %s → %s\n\n", cyan.Sprint(aliasToDelete.Name), aliasToDelete.Command)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

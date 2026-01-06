package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/navio/am/internal/alias"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [search_term]",
	Short: "List all aliases",
	Long: `List all aliases defined in your shell configuration file.
Optionally, provide a search term to filter aliases by name.

Examples:
  am list          # List all aliases
  am list git      # List all aliases containing 'git' in the name`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		searchTerm := ""
		if len(args) > 0 {
			searchTerm = args[0]
		}

		manager := alias.NewManager()
		aliases, err := manager.List(searchTerm)
		if err != nil {
			return err
		}

		if len(aliases) == 0 {
			if searchTerm != "" {
				yellow := color.New(color.FgYellow)
				yellow.Printf("No aliases found matching '%s'\n", searchTerm)
			} else {
				yellow := color.New(color.FgYellow)
				yellow.Println("No aliases found. Use 'am add' to create one.")
			}
			return nil
		}

		// Display header
		cyan := color.New(color.FgCyan, color.Bold)
		if searchTerm != "" {
			cyan.Printf("Aliases matching '%s' (%d found):\n\n", searchTerm, len(aliases))
		} else {
			cyan.Printf("All aliases (%d total):\n\n", len(aliases))
		}

		// Calculate max name length for alignment
		maxNameLen := 0
		for _, a := range aliases {
			if len(a.Name) > maxNameLen {
				maxNameLen = len(a.Name)
			}
		}

		// Display aliases
		nameColor := color.New(color.FgGreen)
		for _, a := range aliases {
			padding := strings.Repeat(" ", maxNameLen-len(a.Name))
			fmt.Printf("  %s%s → %s\n", nameColor.Sprint(a.Name), padding, a.Command)
		}
		fmt.Println()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

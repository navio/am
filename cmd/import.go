package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/navio/am/internal/alias"
	"github.com/spf13/cobra"
)

var (
	importOverwrite bool
	importDryRun    bool
)

var importCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import aliases from a JSON file",
	Long: `Import aliases from a JSON file (as produced by 'am export').
The file must contain a JSON object mapping alias names to commands.

By default, existing aliases are kept (conflicts are reported).
Use --overwrite to replace existing aliases, or --dry-run to preview changes.

Use '-' as the filename to read from stdin.

Examples:
  am import aliases.json
  am import aliases.json --overwrite
  am import aliases.json --dry-run
  cat aliases.json | am import -`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		var data []byte
		var err error
		if path == "-" {
			data, err = io.ReadAll(os.Stdin)
		} else {
			data, err = os.ReadFile(path)
		}
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		entries := map[string]string{}
		if err := json.Unmarshal(data, &entries); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}

		if len(entries) == 0 {
			yellow := color.New(color.FgYellow)
			yellow.Println("No aliases found in input file.")
			return nil
		}

		manager := alias.NewManager()
		result, err := manager.ImportMany(entries, importOverwrite, importDryRun)
		if err != nil {
			return err
		}

		green := color.New(color.FgGreen, color.Bold)
		yellow := color.New(color.FgYellow)
		red := color.New(color.FgRed)
		cyan := color.New(color.FgCyan)

		if importDryRun {
			cyan.Println("Dry run - no changes were written.")
			fmt.Println()
		} else {
			green.Print("✓ ")
			fmt.Print("Import complete!\n\n")
		}

		fmt.Printf("  Added:     %d\n", len(result.Added))
		fmt.Printf("  Updated:   %d\n", len(result.Updated))
		fmt.Printf("  Skipped:   %d\n", len(result.Skipped))
		fmt.Printf("  Conflicts: %d\n", len(result.Conflicts))

		if len(result.Conflicts) > 0 {
			fmt.Println()
			yellow.Println("Conflicts (existing aliases differ; pass --overwrite to replace):")
			for _, name := range result.Conflicts {
				fmt.Printf("  - %s\n", name)
			}
		}

		if len(result.Skipped) > 0 && importDryRun {
			fmt.Println()
			red.Println("Skipped (invalid name or empty command, or identical to existing):")
			for _, name := range result.Skipped {
				fmt.Printf("  - %s\n", name)
			}
		}

		return nil
	},
}

func init() {
	importCmd.Flags().BoolVar(&importOverwrite, "overwrite", false, "Overwrite existing aliases on conflict")
	importCmd.Flags().BoolVar(&importDryRun, "dry-run", false, "Preview changes without writing to disk")
	rootCmd.AddCommand(importCmd)
}

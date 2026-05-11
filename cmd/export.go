package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/navio/am/internal/alias"
	"github.com/spf13/cobra"
)

var (
	exportFormat string
	exportOutput string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export aliases to a file or stdout",
	Long: `Export all aliases from your shell configuration to JSON (default) or shell format.

Examples:
  am export                          # JSON to stdout
  am export -o aliases.json          # JSON to file
  am export --format shell           # shell alias lines to stdout
  am export --format shell -o aliases.sh`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := alias.NewManager()
		aliases, err := manager.List("")
		if err != nil {
			return err
		}

		var output string
		switch strings.ToLower(exportFormat) {
		case "json":
			data := make(map[string]string, len(aliases))
			for _, a := range aliases {
				data[a.Name] = a.Command
			}
			b, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal aliases: %w", err)
			}
			output = string(b) + "\n"
		case "shell", "sh":
			sort.Slice(aliases, func(i, j int) bool { return aliases[i].Name < aliases[j].Name })
			var b strings.Builder
			for _, a := range aliases {
				fmt.Fprintf(&b, "alias %s='%s'\n", a.Name, strings.ReplaceAll(a.Command, "'", `'\''`))
			}
			output = b.String()
		default:
			return fmt.Errorf("unsupported format '%s' (use json or shell)", exportFormat)
		}

		if exportOutput == "" || exportOutput == "-" {
			fmt.Print(output)
			return nil
		}

		if err := os.WriteFile(exportOutput, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write to %s: %w", exportOutput, err)
		}

		green := color.New(color.FgGreen, color.Bold)
		green.Print("✓ ")
		fmt.Printf("Exported %d aliases to %s\n", len(aliases), exportOutput)
		return nil
	},
}

func init() {
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "json", "Output format: json or shell")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file (default: stdout)")
	rootCmd.AddCommand(exportCmd)
}

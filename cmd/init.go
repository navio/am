package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/navio/am/internal/dotfile"
	"github.com/navio/am/internal/shell"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize alias manager with auto-sourcing wrapper",
	Long: `Initialize alias manager by adding a shell function wrapper to your dotfile.
This wrapper automatically sources your shell configuration after add/update/delete
operations, so changes take effect immediately without manual sourcing.

This command:
- Detects your shell (bash, zsh, or fish)
- Adds the am() wrapper function to your dotfile
- Creates a backup before modifying (.bak)
- Only needs to be run once

After running this, restart your shell or run:
  source ~/.zshrc                      (for zsh)
  source ~/.bashrc                     (for bash)
  source ~/.config/fish/config.fish    (for fish)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		detector := shell.NewDetector()
		shellType, err := detector.GetShellType()
		if err != nil {
			return err
		}

		dotfilePath, err := detector.GetDotfilePath()
		if err != nil {
			return err
		}

		// Load current dotfile
		handler := dotfile.NewHandler()
		lines, err := handler.ReadOrCreate(dotfilePath, shellType)
		if err != nil {
			return fmt.Errorf("failed to read dotfile: %w", err)
		}

		// Check if wrapper already exists
		wrapperExists := false
		for _, line := range lines {
			if strings.Contains(line, "# Alias Manager auto-source wrapper") {
				wrapperExists = true
				break
			}
		}

		if wrapperExists {
			yellow := color.New(color.FgYellow)
			yellow.Println("✓ Alias Manager wrapper is already installed!")
			fmt.Printf("\nYour shell will automatically reload after add/update/delete operations.\n")
			return nil
		}

		// Generate wrapper function based on shell type
		wrapper := generateWrapper(shellType)

		// Add wrapper to dotfile
		lines = append(lines, "")
		lines = append(lines, "# Alias Manager auto-source wrapper")
		lines = append(lines, "# This function automatically sources your shell config after am modifications")
		lines = append(lines, wrapper...)

		// Write back to dotfile
		if err := handler.Write(dotfilePath, lines); err != nil {
			return fmt.Errorf("failed to write dotfile: %w", err)
		}

		// Success message
		green := color.New(color.FgGreen, color.Bold)
		cyan := color.New(color.FgCyan)

		green.Print("✓ ")
		fmt.Println("Alias Manager initialized successfully!")
		fmt.Println()
		fmt.Printf("Added auto-source wrapper to: %s\n", cyan.Sprint(dotfilePath))
		fmt.Printf("Backup created at: %s\n", cyan.Sprint(dotfilePath+".bak"))
		fmt.Println()

		yellow := color.New(color.FgYellow, color.Bold)
		yellow.Println("Next step:")
		fmt.Printf("  Restart your shell or run: %s\n\n", cyan.Sprint("source "+dotfilePath))

		fmt.Println("From now on, aliases will take effect immediately after add/update/delete!")

		return nil
	},
}

func generateWrapper(shellType shell.ShellType) []string {
	var wrapper []string

	switch shellType {
	case shell.Zsh:
		wrapper = []string{
			"am() {",
			"  command am \"$@\"",
			"  local exit_code=$?",
			"  # Auto-source after modifying commands",
			"  if [ $exit_code -eq 0 ] && [[ \"$1\" =~ ^(add|update|delete)$ ]]; then",
			"    source ~/.zshrc",
			"  fi",
			"  return $exit_code",
			"}",
		}
	case shell.Bash:
		wrapper = []string{
			"am() {",
			"  command am \"$@\"",
			"  local exit_code=$?",
			"  # Auto-source after modifying commands",
			"  if [ $exit_code -eq 0 ] && [[ \"$1\" =~ ^(add|update|delete)$ ]]; then",
			"    source ~/.bashrc",
			"  fi",
			"  return $exit_code",
			"}",
		}
	case shell.Fish:
		wrapper = []string{
			"function am",
			"  command am $argv",
			"  set exit_code $status",
			"  # Auto-source after modifying commands",
			"  if test $exit_code -eq 0; and string match -qr '^(add|update|delete)$' -- $argv[1]",
			"    source ~/.config/fish/config.fish",
			"  end",
			"  return $exit_code",
			"end",
		}
	default:
		wrapper = []string{
			"am() {",
			"  command am \"$@\"",
			"  local exit_code=$?",
			"  # Auto-source after modifying commands",
			"  if [ $exit_code -eq 0 ] && [[ \"$1\" =~ ^(add|update|delete)$ ]]; then",
			"    source ~/.zshrc",
			"  fi",
			"  return $exit_code",
			"}",
		}
	}

	return wrapper
}

func init() {
	rootCmd.AddCommand(initCmd)
}

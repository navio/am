package alias

import (
	"fmt"
	"strings"

	"github.com/navio/am/internal/dotfile"
	"github.com/navio/am/internal/shell"
)

// Manager handles CRUD operations for aliases
type Manager struct {
	handler  *dotfile.Handler
	parser   *Parser
	detector *shell.Detector
}

// NewManager creates a new Manager instance
func NewManager() *Manager {
	return &Manager{
		handler:  dotfile.NewHandler(),
		parser:   NewParser(),
		detector: shell.NewDetector(),
	}
}

// LoadAliases loads all aliases from the appropriate dotfile
func (m *Manager) LoadAliases() (*AliasFile, error) {
	dotfilePath, err := m.detector.GetDotfilePath()
	if err != nil {
		return nil, err
	}

	shellType, err := m.detector.GetShellType()
	if err != nil {
		return nil, err
	}

	lines, err := m.handler.ReadOrCreate(dotfilePath, shellType)
	if err != nil {
		return nil, fmt.Errorf("failed to read dotfile: %w", err)
	}

	aliasFile := NewAliasFile(dotfilePath, lines)
	if err := m.parser.Parse(aliasFile); err != nil {
		return nil, fmt.Errorf("failed to parse aliases: %w", err)
	}

	return aliasFile, nil
}

// Add adds a new alias to the dotfile
func (m *Manager) Add(name, command string) error {
	// Validate alias name
	if err := m.validateAliasName(name); err != nil {
		return err
	}

	if strings.TrimSpace(command) == "" {
		return fmt.Errorf("command cannot be empty")
	}

	aliasFile, err := m.LoadAliases()
	if err != nil {
		return err
	}

	// Check for duplicate
	if m.findAlias(aliasFile, name) != nil {
		return fmt.Errorf("alias '%s' already exists. Use 'update' to modify it", name)
	}

	// Format the new alias line
	aliasLine := m.parser.FormatAlias(name, command)

	// Append to the end of the file
	aliasFile.Lines = append(aliasFile.Lines, aliasLine)

	// Write back to file
	if err := m.handler.Write(aliasFile.Path, aliasFile.Lines); err != nil {
		return fmt.Errorf("failed to write dotfile: %w", err)
	}

	return nil
}

// List returns all aliases, optionally filtered by a search term
func (m *Manager) List(searchTerm string) ([]Alias, error) {
	aliasFile, err := m.LoadAliases()
	if err != nil {
		return nil, err
	}

	if searchTerm == "" {
		return aliasFile.Aliases, nil
	}

	// Filter aliases by search term (case-insensitive substring match)
	searchLower := strings.ToLower(searchTerm)
	var filtered []Alias
	for _, alias := range aliasFile.Aliases {
		if strings.Contains(strings.ToLower(alias.Name), searchLower) {
			filtered = append(filtered, alias)
		}
	}

	return filtered, nil
}

// Delete removes an alias from the dotfile
func (m *Manager) Delete(name string) error {
	aliasFile, err := m.LoadAliases()
	if err != nil {
		return err
	}

	alias := m.findAlias(aliasFile, name)
	if alias == nil {
		return fmt.Errorf("alias '%s' not found", name)
	}

	// Remove the line containing the alias
	aliasFile.Lines = append(
		aliasFile.Lines[:alias.LineNum],
		aliasFile.Lines[alias.LineNum+1:]...,
	)

	// Write back to file
	if err := m.handler.Write(aliasFile.Path, aliasFile.Lines); err != nil {
		return fmt.Errorf("failed to write dotfile: %w", err)
	}

	return nil
}

// Update modifies an existing alias
func (m *Manager) Update(name, newCommand string) error {
	if strings.TrimSpace(newCommand) == "" {
		return fmt.Errorf("command cannot be empty")
	}

	aliasFile, err := m.LoadAliases()
	if err != nil {
		return err
	}

	alias := m.findAlias(aliasFile, name)
	if alias == nil {
		return fmt.Errorf("alias '%s' not found. Use 'add' to create it", name)
	}

	// Replace the line with the updated alias
	aliasFile.Lines[alias.LineNum] = m.parser.FormatAlias(name, newCommand)

	// Write back to file
	if err := m.handler.Write(aliasFile.Path, aliasFile.Lines); err != nil {
		return fmt.Errorf("failed to write dotfile: %w", err)
	}

	return nil
}

// ImportResult summarizes the outcome of an import operation
type ImportResult struct {
	Added     []string
	Updated   []string
	Skipped   []string
	Conflicts []string
}

// ImportMany imports multiple aliases at once.
// If overwrite is true, existing aliases are updated; otherwise they are skipped.
// If dryRun is true, no changes are written to disk.
func (m *Manager) ImportMany(entries map[string]string, overwrite, dryRun bool) (*ImportResult, error) {
	aliasFile, err := m.LoadAliases()
	if err != nil {
		return nil, err
	}

	result := &ImportResult{}

	for name, command := range entries {
		if err := m.validateAliasName(name); err != nil {
			result.Skipped = append(result.Skipped, name)
			continue
		}
		if strings.TrimSpace(command) == "" {
			result.Skipped = append(result.Skipped, name)
			continue
		}

		existing := m.findAlias(aliasFile, name)
		if existing == nil {
			aliasFile.Lines = append(aliasFile.Lines, m.parser.FormatAlias(name, command))
			// Re-parse so subsequent lookups in this loop find the just-added alias
			if err := m.parser.Parse(aliasFile); err != nil {
				return nil, fmt.Errorf("failed to re-parse aliases: %w", err)
			}
			result.Added = append(result.Added, name)
			continue
		}

		if existing.Command == command {
			result.Skipped = append(result.Skipped, name)
			continue
		}

		if overwrite {
			aliasFile.Lines[existing.LineNum] = m.parser.FormatAlias(name, command)
			result.Updated = append(result.Updated, name)
		} else {
			result.Conflicts = append(result.Conflicts, name)
		}
	}

	if dryRun {
		return result, nil
	}

	if len(result.Added) == 0 && len(result.Updated) == 0 {
		return result, nil
	}

	if err := m.handler.Write(aliasFile.Path, aliasFile.Lines); err != nil {
		return nil, fmt.Errorf("failed to write dotfile: %w", err)
	}

	return result, nil
}

// Get retrieves a specific alias by name
func (m *Manager) Get(name string) (*Alias, error) {
	aliasFile, err := m.LoadAliases()
	if err != nil {
		return nil, err
	}

	alias := m.findAlias(aliasFile, name)
	if alias == nil {
		return nil, fmt.Errorf("alias '%s' not found", name)
	}

	return alias, nil
}

// findAlias searches for an alias by name in the alias file
func (m *Manager) findAlias(af *AliasFile, name string) *Alias {
	for i := range af.Aliases {
		if af.Aliases[i].Name == name {
			return &af.Aliases[i]
		}
	}
	return nil
}

// validateAliasName checks if the alias name is valid
func (m *Manager) validateAliasName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("alias name cannot be empty")
	}

	// Check for valid characters (alphanumeric, underscore, hyphen)
	for _, ch := range name {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '_' || ch == '-') {
			return fmt.Errorf("alias name '%s' contains invalid characters. Use only alphanumeric, underscore, or hyphen", name)
		}
	}

	return nil
}

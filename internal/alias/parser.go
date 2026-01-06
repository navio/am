package alias

import (
	"regexp"
	"strings"
)

var (
	// aliasRegex matches alias declarations in various formats:
	// alias name='command'
	// alias name="command"
	// alias name=command
	aliasRegex = regexp.MustCompile(`^\s*alias\s+([a-zA-Z0-9_-]+)=(.+)$`)
)

// Parser handles parsing of alias definitions from dotfiles
type Parser struct{}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{}
}

// Parse extracts all aliases from an AliasFile
func (p *Parser) Parse(af *AliasFile) error {
	for i, line := range af.Lines {
		alias, ok := p.ParseLine(line)
		if ok {
			alias.LineNum = i
			af.Aliases = append(af.Aliases, alias)
		}
	}
	return nil
}

// ParseLine attempts to parse a single line as an alias definition
// Returns the Alias and true if successful, empty Alias and false otherwise
func (p *Parser) ParseLine(line string) (Alias, bool) {
	matches := aliasRegex.FindStringSubmatch(line)
	if matches == nil {
		return Alias{}, false
	}

	name := matches[1]
	command := matches[2]

	// Remove surrounding quotes if present
	command = p.unquoteCommand(command)

	return Alias{
		Name:    name,
		Command: command,
	}, true
}

// unquoteCommand removes surrounding quotes from a command string
func (p *Parser) unquoteCommand(cmd string) string {
	cmd = strings.TrimSpace(cmd)

	// Remove surrounding single quotes
	if strings.HasPrefix(cmd, "'") && strings.HasSuffix(cmd, "'") && len(cmd) >= 2 {
		return cmd[1 : len(cmd)-1]
	}

	// Remove surrounding double quotes
	if strings.HasPrefix(cmd, "\"") && strings.HasSuffix(cmd, "\"") && len(cmd) >= 2 {
		return cmd[1 : len(cmd)-1]
	}

	return cmd
}

// QuoteCommand adds appropriate quotes to a command string for safe alias definition
func (p *Parser) QuoteCommand(cmd string) string {
	// Use single quotes by default as they preserve the command exactly
	// Escape any single quotes in the command
	cmd = strings.ReplaceAll(cmd, "'", "'\\''")
	return "'" + cmd + "'"
}

// FormatAlias formats an alias into the standard alias definition format
func (p *Parser) FormatAlias(name, command string) string {
	return "alias " + name + "=" + p.QuoteCommand(command)
}

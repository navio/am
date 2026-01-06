package alias

// Alias represents a shell alias with its name, command, and line number in the dotfile
type Alias struct {
	Name    string
	Command string
	LineNum int // Line number in the dotfile (0-indexed)
}

// AliasFile represents the content of a dotfile with aliases
type AliasFile struct {
	Path    string
	Lines   []string
	Aliases []Alias
}

// NewAliasFile creates a new AliasFile with the given path and lines
func NewAliasFile(path string, lines []string) *AliasFile {
	return &AliasFile{
		Path:    path,
		Lines:   lines,
		Aliases: []Alias{},
	}
}

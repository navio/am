package alias

import (
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		wantAlias   Alias
		wantSuccess bool
	}{
		{
			name:        "simple alias with single quotes",
			line:        "alias ll='ls -la'",
			wantAlias:   Alias{Name: "ll", Command: "ls -la"},
			wantSuccess: true,
		},
		{
			name:        "alias with double quotes",
			line:        `alias gs="git status"`,
			wantAlias:   Alias{Name: "gs", Command: "git status"},
			wantSuccess: true,
		},
		{
			name:        "alias without quotes",
			line:        "alias gs=git",
			wantAlias:   Alias{Name: "gs", Command: "git"},
			wantSuccess: true,
		},
		{
			name:        "alias with spaces before",
			line:        "  alias ll='ls -la'",
			wantAlias:   Alias{Name: "ll", Command: "ls -la"},
			wantSuccess: true,
		},
		{
			name:        "alias with complex command",
			line:        "alias serve='python -m http.server 8000'",
			wantAlias:   Alias{Name: "serve", Command: "python -m http.server 8000"},
			wantSuccess: true,
		},
		{
			name:        "alias with hyphen in name",
			line:        "alias git-log='git log --oneline'",
			wantAlias:   Alias{Name: "git-log", Command: "git log --oneline"},
			wantSuccess: true,
		},
		{
			name:        "alias with underscore in name",
			line:        "alias my_alias='echo hello'",
			wantAlias:   Alias{Name: "my_alias", Command: "echo hello"},
			wantSuccess: true,
		},
		{
			name:        "not an alias - comment",
			line:        "# alias ll='ls -la'",
			wantSuccess: false,
		},
		{
			name:        "not an alias - empty line",
			line:        "",
			wantSuccess: false,
		},
		{
			name:        "not an alias - export",
			line:        "export PATH=/usr/bin",
			wantSuccess: false,
		},
		{
			name:        "alias with quotes in command",
			line:        `alias echo-test='echo "hello world"'`,
			wantAlias:   Alias{Name: "echo-test", Command: `echo "hello world"`},
			wantSuccess: true,
		},
	}

	parser := NewParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAlias, gotSuccess := parser.ParseLine(tt.line)

			if gotSuccess != tt.wantSuccess {
				t.Errorf("ParseLine() success = %v, want %v", gotSuccess, tt.wantSuccess)
				return
			}

			if !tt.wantSuccess {
				return
			}

			if gotAlias.Name != tt.wantAlias.Name {
				t.Errorf("ParseLine() name = %v, want %v", gotAlias.Name, tt.wantAlias.Name)
			}

			if gotAlias.Command != tt.wantAlias.Command {
				t.Errorf("ParseLine() command = %v, want %v", gotAlias.Command, tt.wantAlias.Command)
			}
		})
	}
}

func TestUnquoteCommand(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single quotes",
			input: "'ls -la'",
			want:  "ls -la",
		},
		{
			name:  "double quotes",
			input: `"git status"`,
			want:  "git status",
		},
		{
			name:  "no quotes",
			input: "echo hello",
			want:  "echo hello",
		},
		{
			name:  "quotes with spaces",
			input: "  'ls -la'  ",
			want:  "ls -la",
		},
		{
			name:  "single quote only at start",
			input: "'ls -la",
			want:  "'ls -la",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
	}

	parser := NewParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.unquoteCommand(tt.input)
			if got != tt.want {
				t.Errorf("unquoteCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuoteCommand(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple command",
			input: "ls -la",
			want:  "'ls -la'",
		},
		{
			name:  "command with quotes",
			input: `echo "hello"`,
			want:  `'echo "hello"'`,
		},
		{
			name:  "command with single quote",
			input: "echo 'hello'",
			want:  `'echo '\''hello'\'''`,
		},
		{
			name:  "empty command",
			input: "",
			want:  "''",
		},
	}

	parser := NewParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.QuoteCommand(tt.input)
			if got != tt.want {
				t.Errorf("QuoteCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatAlias(t *testing.T) {
	tests := []struct {
		name    string
		aliasName string
		command string
		want    string
	}{
		{
			name:    "simple alias",
			aliasName: "ll",
			command: "ls -la",
			want:    "alias ll='ls -la'",
		},
		{
			name:    "complex command",
			aliasName: "serve",
			command: "python -m http.server 8000",
			want:    "alias serve='python -m http.server 8000'",
		},
	}

	parser := NewParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.FormatAlias(tt.aliasName, tt.command)
			if got != tt.want {
				t.Errorf("FormatAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	lines := []string{
		"# Shell configuration",
		"alias ll='ls -la'",
		"export PATH=/usr/bin",
		"alias gs='git status'",
		"",
		"alias serve='python -m http.server 8000'",
		"# Another comment",
	}

	af := NewAliasFile("/tmp/test", lines)
	parser := NewParser()

	err := parser.Parse(af)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(af.Aliases) != 3 {
		t.Errorf("Parse() found %d aliases, want 3", len(af.Aliases))
	}

	// Check first alias
	if af.Aliases[0].Name != "ll" {
		t.Errorf("First alias name = %v, want 'll'", af.Aliases[0].Name)
	}
	if af.Aliases[0].Command != "ls -la" {
		t.Errorf("First alias command = %v, want 'ls -la'", af.Aliases[0].Command)
	}
	if af.Aliases[0].LineNum != 1 {
		t.Errorf("First alias line number = %v, want 1", af.Aliases[0].LineNum)
	}

	// Check second alias
	if af.Aliases[1].Name != "gs" {
		t.Errorf("Second alias name = %v, want 'gs'", af.Aliases[1].Name)
	}

	// Check third alias
	if af.Aliases[2].Name != "serve" {
		t.Errorf("Third alias name = %v, want 'serve'", af.Aliases[2].Name)
	}
}

# Alias Manager (am)

A lightweight, fast CLI tool for managing shell aliases directly in your `.zshrc` or `.bashrc` files.

## Features

- **Simple Commands**: Easy-to-use interface for managing aliases
- **Direct File Manipulation**: Works directly with your dotfiles - no separate database
- **Automatic Backups**: Creates `.bak` files before any modification for safety
- **Shell Detection**: Automatically detects whether you're using bash or zsh
- **Search Functionality**: Filter aliases by name
- **Preserves Structure**: Maintains all comments, blank lines, and file organization
- **Zero Dependencies**: Single binary with no runtime dependencies
- **Cross-Platform**: Works on macOS, Linux, and other Unix-like systems

## Installation

### Via Homebrew (Recommended)

```bash
# One-step install
brew install navio/tap/am

# Or tap first, then install
brew tap navio/tap
brew install am
```

### From Source

```bash
# Clone the repository
git clone https://github.com/navio/am.git
cd am

# Build the binary
go build -o am

# Move to a directory in your PATH
sudo mv am /usr/local/bin/

# Or install to Go's bin directory
go install
```

### Download Binary (Coming Soon)

Download the latest release for your platform from the [releases page](https://github.com/navio/am/releases).

## Quick Start

### 1. Initialize (First Time Only)

Run this command once to enable auto-sourcing:

```bash
am init
```

This adds a shell wrapper function to your `.zshrc` or `.bashrc` that automatically reloads your shell configuration after adding, updating, or deleting aliases. After running `am init`, restart your shell or run:

```bash
source ~/.zshrc  # for zsh
source ~/.bashrc # for bash
```

**From this point forward, all alias changes take effect immediately!**

### 2. Start Managing Aliases

## Usage

### Add an Alias

```bash
am add ll 'ls -la'
am add gs 'git status'
am add serve 'python -m http.server 8000'
```

### List All Aliases

```bash
am list
```

Output:
```
All aliases (3 total):

  ll    → ls -la
  gs    → git status
  serve → python -m http.server 8000
```

### Search for Aliases

```bash
am list git
```

Output:
```
Aliases matching 'git' (1 found):

  gs → git status
```

### Update an Alias

```bash
am update ll 'ls -lah'
```

### Delete an Alias

```bash
am delete serve
```

### Get Help

```bash
am --help
am add --help
```

### Version

```bash
am --version
```

## How It Works

Alias Manager works directly with your shell configuration files:

1. **Shell Detection**: Automatically detects your shell (`$SHELL` environment variable)
2. **File Location**: Finds the appropriate dotfile:
   - Zsh: `~/.zshrc`
   - Bash: `~/.bashrc`
3. **Safe Modifications**: Creates a `.bak` backup before any write operation
4. **Structure Preservation**: Only modifies alias lines, preserving everything else

## Command Reference

| Command | Description | Example |
|---------|-------------|---------|
| `am init` | Initialize auto-sourcing wrapper (run once) | `am init` |
| `am add <name> <command>` | Add a new alias | `am add ll 'ls -la'` |
| `am list [search]` | List all aliases, optionally filtered | `am list git` |
| `am update <name> <command>` | Update an existing alias | `am update ll 'ls -lah'` |
| `am delete <name>` | Delete an alias | `am delete ll` |
| `am --version` | Show version | `am --version` |
| `am --help` | Show help | `am --help` |

## Examples

### Managing Development Aliases

```bash
# Add common git aliases
am add gs 'git status'
am add gc 'git commit -m'
am add gp 'git push'
am add gl 'git log --oneline --graph'

# Add navigation aliases
am add ..  'cd ..'
am add ... 'cd ../..'
am add home 'cd ~'

# Add utility aliases
am add serve 'python -m http.server 8000'
am add ports 'lsof -i -P | grep LISTEN'
```

### Searching and Managing

```bash
# Find all git-related aliases
am list git

# Update an alias
am update gs 'git status -sb'

# Remove unused aliases
am delete serve
```

## Safety Features

- **Automatic Backups**: Every modification creates a `.bak` file
- **Duplicate Detection**: Prevents adding aliases that already exist
- **Validation**: Checks for valid alias names (alphanumeric, underscore, hyphen only)
- **Atomic Writes**: Uses temporary files to ensure safe writes
- **Structure Preservation**: Never removes comments or formatting

## Backup Recovery

If something goes wrong, your backup file is always available:

```bash
# For zsh users
cp ~/.zshrc.bak ~/.zshrc
source ~/.zshrc

# For bash users
cp ~/.bashrc.bak ~/.bashrc
source ~/.bashrc
```

## Applying Changes

### With `am init` (Recommended)

If you've run `am init`, changes apply automatically! No manual sourcing needed.

### Without `am init`

If you haven't initialized the auto-source wrapper, manually reload your shell:

```bash
# For zsh
source ~/.zshrc

# For bash
source ~/.bashrc

# Or simply start a new shell session
```

## Development

### Prerequisites

- Go 1.21 or later

### Building from Source

```bash
# Clone the repository
git clone https://github.com/navio/am.git
cd am

# Install dependencies
go mod download

# Build
go build -o am

# Run tests
go test ./...

# Run tests with coverage
go test ./... -cover
```

### Project Structure

```
am/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command
│   ├── add.go             # Add command
│   ├── list.go            # List command
│   ├── delete.go          # Delete command
│   └── update.go          # Update command
├── internal/
│   ├── alias/             # Alias management logic
│   │   ├── models.go      # Data structures
│   │   ├── parser.go      # Alias parsing
│   │   └── manager.go     # CRUD operations
│   ├── dotfile/           # File I/O with backups
│   │   └── handler.go
│   └── shell/             # Shell detection
│       └── detector.go
├── main.go                # Entry point
└── *_test.go             # Tests
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run with coverage
go test ./... -cover

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Supported Shells

- **Zsh** (`.zshrc`)
- **Bash** (`.bashrc`)

## Limitations

- Only works with bash and zsh
- Aliases must be in the format: `alias name='command'`
- Does not support multi-line aliases (with backslash continuation)
- Does not support alias functions or complex shell scripting

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see LICENSE file for details

## Author

Created by [navio](https://github.com/navio)

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) - CLI framework
- Colored output with [color](https://github.com/fatih/color)

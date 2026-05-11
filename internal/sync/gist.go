package sync

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ErrGhMissing is returned when the gh CLI is not available on PATH.
var ErrGhMissing = errors.New("the GitHub CLI (gh) is required for 'am sync'. Install it from https://cli.github.com and run 'gh auth login'")

// GistClient talks to GitHub Gists via the gh CLI.
type GistClient struct {
	bin string
}

// NewGistClient verifies gh is installed and returns a client.
func NewGistClient() (*GistClient, error) {
	bin, err := exec.LookPath("gh")
	if err != nil {
		return nil, ErrGhMissing
	}
	return &GistClient{bin: bin}, nil
}

func (c *GistClient) run(stdin []byte, args ...string) ([]byte, error) {
	cmd := exec.Command(c.bin, args...)
	if stdin != nil {
		cmd.Stdin = bytes.NewReader(stdin)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("gh %s: %s", strings.Join(args, " "), msg)
	}
	return stdout.Bytes(), nil
}

// Create creates a new private gist with the given filename and content.
// Returns the gist ID.
func (c *GistClient) Create(filename string, content []byte, description string) (string, error) {
	tmp, err := os.CreateTemp("", "am-sync-*.json")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(content); err != nil {
		tmp.Close()
		return "", err
	}
	tmp.Close()

	args := []string{"gist", "create", tmp.Name(), "--filename", filename}
	if description != "" {
		args = append(args, "--desc", description)
	}
	out, err := c.run(nil, args...)
	if err != nil {
		return "", err
	}

	// gh prints the gist URL on the last non-empty line.
	url := ""
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "https://gist.github.com/") {
			url = line
		}
	}
	if url == "" {
		return "", fmt.Errorf("could not parse gist URL from gh output:\n%s", string(out))
	}
	parts := strings.Split(url, "/")
	return parts[len(parts)-1], nil
}

// View fetches the content of a single file from a gist.
func (c *GistClient) View(gistID, filename string) ([]byte, error) {
	out, err := c.run(nil, "gist", "view", gistID, "--filename", filename, "--raw")
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Edit replaces the content of a single file in a gist.
func (c *GistClient) Edit(gistID, filename string, content []byte) error {
	tmp, err := os.CreateTemp("", "am-sync-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(content); err != nil {
		tmp.Close()
		return err
	}
	tmp.Close()

	_, err = c.run(nil, "gist", "edit", gistID, "--filename", filename, tmp.Name())
	return err
}

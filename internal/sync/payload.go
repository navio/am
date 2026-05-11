package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/navio/am/internal/alias"
)

// DefaultGistFilename is the filename used inside a sync gist.
const DefaultGistFilename = "am-aliases.json"

// AliasesToMap converts a slice of aliases into a name->command map.
func AliasesToMap(aliases []alias.Alias) map[string]string {
	out := make(map[string]string, len(aliases))
	for _, a := range aliases {
		out[a.Name] = a.Command
	}
	return out
}

// Marshal serializes an alias map to canonical, sorted JSON bytes.
// The output is stable so hashes are reproducible across machines.
func Marshal(entries map[string]string) ([]byte, error) {
	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	ordered := make([]struct {
		K string
		V string
	}, 0, len(keys))
	for _, k := range keys {
		ordered = append(ordered, struct {
			K string
			V string
		}{k, entries[k]})
	}

	// json package does not preserve insertion order for maps; build manually.
	buf := []byte{'{', '\n'}
	for i, kv := range ordered {
		kb, err := json.Marshal(kv.K)
		if err != nil {
			return nil, err
		}
		vb, err := json.Marshal(kv.V)
		if err != nil {
			return nil, err
		}
		buf = append(buf, ' ', ' ')
		buf = append(buf, kb...)
		buf = append(buf, ':', ' ')
		buf = append(buf, vb...)
		if i < len(ordered)-1 {
			buf = append(buf, ',')
		}
		buf = append(buf, '\n')
	}
	buf = append(buf, '}', '\n')
	return buf, nil
}

// Unmarshal parses a sync payload back into an alias map.
func Unmarshal(data []byte) (map[string]string, error) {
	out := map[string]string{}
	if len(data) == 0 {
		return out, nil
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("failed to parse sync payload: %w", err)
	}
	return out, nil
}

// Hash returns a short SHA-256 fingerprint of the canonical payload bytes.
func Hash(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])[:12]
}

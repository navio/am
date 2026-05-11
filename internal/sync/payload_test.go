package sync

import (
	"strings"
	"testing"
)

func TestMarshalIsStable(t *testing.T) {
	a := map[string]string{"b": "2", "a": "1", "c": "3"}
	b := map[string]string{"c": "3", "a": "1", "b": "2"}

	ba, err := Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	bb, err := Marshal(b)
	if err != nil {
		t.Fatal(err)
	}
	if string(ba) != string(bb) {
		t.Fatalf("Marshal should be order-independent:\n%s\nvs\n%s", ba, bb)
	}

	idxA := strings.Index(string(ba), `"a"`)
	idxB := strings.Index(string(ba), `"b"`)
	idxC := strings.Index(string(ba), `"c"`)
	if !(idxA < idxB && idxB < idxC) {
		t.Fatalf("keys should be sorted alphabetically, got:\n%s", ba)
	}
}

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	in := map[string]string{"ll": "ls -la", "gs": "git status"}
	data, err := Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	out, err := Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != len(in) {
		t.Fatalf("length mismatch: %d vs %d", len(out), len(in))
	}
	for k, v := range in {
		if out[k] != v {
			t.Fatalf("entry %q: got %q want %q", k, out[k], v)
		}
	}
}

func TestHashIsDeterministicAndShort(t *testing.T) {
	data := []byte(`{"a":"1"}`)
	h1 := Hash(data)
	h2 := Hash(data)
	if h1 != h2 {
		t.Fatal("hash should be deterministic")
	}
	if len(h1) != 12 {
		t.Fatalf("hash length = %d, want 12", len(h1))
	}
	if Hash([]byte(`{"a":"2"}`)) == h1 {
		t.Fatal("different inputs should hash differently")
	}
}

func TestUnmarshalEmpty(t *testing.T) {
	out, err := Unmarshal(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}

package lockfile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSetInitializesEntries(t *testing.T) {
	l := &Lockfile{}
	l.Set("redis:7", "redis@sha256:abc")

	if len(l.entries) != 1 {
		t.Fatalf("expected entries to be initialized, got %d", len(l.entries))
	}
}

func TestSaveWritesRefWhenMissing(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "imgpin.lock")
	l := &Lockfile{
		path: tmp,
		entries: map[string]Entry{
			"redis:7": {Digest: "redis@sha256:abc"},
		},
	}

	if err := l.Save(); err != nil {
		t.Fatalf("save lockfile: %v", err)
	}

	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("read lockfile: %v", err)
	}

	var payload struct {
		Entries map[string]Entry `json:"entries"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal lockfile: %v", err)
	}
	if payload.Entries["redis:7"].Ref != "redis:7" {
		t.Fatalf("expected ref to be populated, got %q", payload.Entries["redis:7"].Ref)
	}
}

func TestSaveNilLockfile(t *testing.T) {
	var l *Lockfile
	if err := l.Save(); err == nil {
		t.Fatalf("expected error for nil lockfile")
	}
}

func TestEntriesReturnsCopy(t *testing.T) {
	l := &Lockfile{
		entries: map[string]Entry{
			"redis:7": {Ref: "redis:7", Digest: "sha256:abc"},
		},
	}

	cp := l.Entries()
	cp["redis:7"] = Entry{Ref: "redis:7", Digest: "mutated"}

	if l.entries["redis:7"].Digest == "mutated" {
		t.Fatalf("entries should return a copy")
	}
}

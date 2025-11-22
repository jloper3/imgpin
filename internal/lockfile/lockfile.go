package lockfile

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// Entry captures a resolved reference and its digest.
type Entry struct {
	Ref    string `json:"ref"`
	Digest string `json:"digest"`
}

// Lockfile stores image entries and the on-disk path backing them.
type Lockfile struct {
	path    string
	entries map[string]Entry
}

type filePayload struct {
	Entries map[string]Entry `json:"entries"`
}

// Load opens an existing lockfile or returns a new in-memory one if it does not exist yet.
func Load(path string) (*Lockfile, error) {
	lf := &Lockfile{
		path:    path,
		entries: map[string]Entry{},
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return lf, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return lf, nil
	}

	var payload filePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	for ref, entry := range payload.Entries {
		if entry.Ref == "" {
			entry.Ref = ref
		}
		lf.entries[ref] = entry
	}
	return lf, nil
}

// Set records or updates an image reference within the lockfile.
func (l *Lockfile) Set(ref, digest string) {
	if l.entries == nil {
		l.entries = map[string]Entry{}
	}
	l.entries[ref] = Entry{
		Ref:    ref,
		Digest: digest,
	}
}

// Get retrieves an entry if present.
func (l *Lockfile) Get(ref string) (Entry, bool) {
	entry, ok := l.entries[ref]
	return entry, ok
}

// Save writes the lockfile contents back to disk atomically.
func (l *Lockfile) Save() error {
	if l == nil {
		return errors.New("lockfile is nil")
	}
	if err := os.MkdirAll(filepath.Dir(l.path), 0o755); err != nil {
		return err
	}

	payload := filePayload{
		Entries: map[string]Entry{},
	}
	for ref, entry := range l.entries {
		if entry.Ref == "" {
			entry.Ref = ref
		}
		payload.Entries[ref] = entry
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}

	tempPath := filepath.Join(filepath.Dir(l.path), "."+filepath.Base(l.path)+".tmp")
	if err := os.WriteFile(tempPath, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tempPath, l.path)
}

// Entries returns a copy of the lockfile entries map.
func (l *Lockfile) Entries() map[string]Entry {
	if l == nil {
		return nil
	}
	copyMap := make(map[string]Entry, len(l.entries))
	for ref, entry := range l.entries {
		copyMap[ref] = entry
	}
	return copyMap
}

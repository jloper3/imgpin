package resolve

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"imgpin/internal/cache"

	"github.com/google/go-containerregistry/pkg/name"
)

func writeCacheFile(t *testing.T, path string, entries map[string]cache.Entry) {
	t.Helper()
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		t.Fatalf("marshal cache: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write cache: %v", err)
	}
}

func TestResolveUsesFreshCache(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "cache.json")
	writeCacheFile(t, tmp, map[string]cache.Entry{
		"repo/app:1.0": {
			Digest:    "repo/app@sha256:cached",
			FetchedAt: time.Now(),
		},
	})

	oldPath := cachePath
	cachePath = tmp
	t.Cleanup(func() { cachePath = oldPath })

	var remoteCalled bool
	prevFetcher := remoteFetchDigest
	remoteFetchDigest = func(ref name.Reference) (string, error) {
		remoteCalled = true
		return "sha256:new", nil
	}
	t.Cleanup(func() { remoteFetchDigest = prevFetcher })

	got, err := Resolve("repo/app:1.0")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if got != "repo/app@sha256:cached" {
		t.Fatalf("unexpected digest %q", got)
	}
	if remoteCalled {
		t.Fatalf("remote fetch should not have been called for fresh cache")
	}
}

func TestResolveFetchesAndPersists(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "cache.json")
	oldPath := cachePath
	cachePath = tmp
	t.Cleanup(func() { cachePath = oldPath })

	prevFetcher := remoteFetchDigest
	remoteFetchDigest = func(ref name.Reference) (string, error) {
		return "sha256:new", nil
	}
	t.Cleanup(func() { remoteFetchDigest = prevFetcher })

	got, err := Resolve("alpine:3.19")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}

	ref, _ := name.ParseReference("alpine:3.19")
	want := ref.Context().Name() + "@sha256:new"
	if got != want {
		t.Fatalf("unexpected digest, want %q got %q", want, got)
	}

	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("read cache: %v", err)
	}

	var cached map[string]cache.Entry
	if err := json.Unmarshal(data, &cached); err != nil {
		t.Fatalf("unmarshal cache: %v", err)
	}
	entry, ok := cached["alpine:3.19"]
	if !ok || entry.Digest != want {
		t.Fatalf("expected entry stored in cache, got %#v", cached["alpine:3.19"])
	}
	if entry.FetchedAt.IsZero() {
		t.Fatalf("expected fetched time recorded")
	}
}

func TestResolveRefreshesStaleEntry(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "cache.json")
	writeCacheFile(t, tmp, map[string]cache.Entry{
		"repo/app:1.0": {
			Digest:    "repo/app@sha256:old",
			FetchedAt: time.Now().Add(-2 * ttl),
		},
	})
	oldPath := cachePath
	cachePath = tmp
	t.Cleanup(func() { cachePath = oldPath })

	var calls int
	prevFetcher := remoteFetchDigest
	remoteFetchDigest = func(ref name.Reference) (string, error) {
		calls++
		return "sha256:new", nil
	}
	t.Cleanup(func() { remoteFetchDigest = prevFetcher })

	got, err := Resolve("repo/app:1.0")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected remote fetch for stale entry")
	}
	if !strings.Contains(got, "@sha256:new") {
		t.Fatalf("expected refreshed digest, got %q", got)
	}
}

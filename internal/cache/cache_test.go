package cache_test

import (
    "testing"
    "time"
    "path/filepath"
    cache "imgpin/internal/cache"
)

func TestCacheRoundtrip(t *testing.T) {
    path := filepath.Join(t.TempDir(), "cache.json")
    c, _ := cache.Load(path)
    c.Set("redis:7", "sha256:xyz")
    if err := c.Save(); err != nil { t.Fatal(err) }

    c2, err := cache.Load(path)
    if err != nil { t.Fatal(err) }

    e, ok := c2.Get("redis:7")
    if !ok { t.Fatalf("missing cache entry") }

    if e.Digest != "sha256:xyz" {
        t.Fatalf("wrong digest %s", e.Digest)
    }
    if time.Since(e.FetchedAt) > time.Minute {
        t.Fatalf("timestamp seems incorrect")
    }
}

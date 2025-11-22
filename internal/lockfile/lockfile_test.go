package lockfile_test

import (
    "testing"
    "path/filepath"
    lf "imgpin/internal/lockfile"
)

func TestLockfileRoundtrip(t *testing.T) {
    path := filepath.Join(t.TempDir(), "imgpin.lock")

    l, _ := lf.Load(path)
    l.Set("alpine:3.19", "sha256:zzz111")
    if err := l.Save(); err != nil { t.Fatal(err) }

    l2, err := lf.Load(path)
    if err != nil { t.Fatal(err) }

    e, ok := l2.Get("alpine:3.19")
    if !ok { t.Fatalf("missing entry") }

    if e.Digest != "sha256:zzz111" {
        t.Fatalf("unexpected digest: %s", e.Digest)
    }
}

package cli

import (
	"imgpin/internal/lockfile"
	"testing"
)

func TestDriftExists(t *testing.T) {
	expected := map[string]lockfile.Entry{
		"nginx:1.27": {Ref: "nginx:1.27", Digest: "sha256:aaa"},
	}

	actual := map[string]lockfile.Entry{
		"nginx:1.27": {Ref: "nginx:1.27", Digest: "sha256:bbb"},
	}

	if !driftExistsForTest(expected, actual) {
		t.Fatalf("expected drift but none detected")
	}
}

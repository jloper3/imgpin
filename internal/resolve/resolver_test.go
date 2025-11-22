package resolve_test

import (
	"imgpin/internal/resolve"
	"testing"
)

func TestResolveFormatOnly(t *testing.T) {
	// We cannot hit registries during tests, so just validate errors behave sanely.
	_, err := resolve.Resolve("not a ref")
	if err == nil {
		t.Fatalf("expected error for invalid ref")
	}

	// Ensure reference formatting rejection works
	_, err2 := resolve.Resolve("bad@@ref")
	if err2 == nil {
		t.Fatalf("expected parse error")
	}

	// NOTE: Full resolver tests require remote.Get mocking or wrapper injection.
}

package dockerfile_test

import (
	dkh "imgpin/internal/handlers/dockerfile"
	"testing"
)

func TestPinDockerfile(t *testing.T) {
	input := []byte("FROM python:3.11\nRUN echo hi")

	resolver := func(ref string) (string, error) {
		return "python@sha256:abc123", nil
	}

	out, changed, err := dkh.Pin(input, resolver)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatalf("expected change")
	}

	if string(out) == string(input) {
		t.Fatalf("expected modification, got unchanged output")
	}
}

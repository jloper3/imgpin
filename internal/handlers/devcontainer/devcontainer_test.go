package devcontainer_test

import (
    "testing"
    dcv "imgpin/internal/handlers/devcontainer"
)

func TestPinDevContainer(t *testing.T) {
    input := []byte(`{"image": "ubuntu:22.04"}`)

    resolver := func(ref string) (string, error) {
        return "ubuntu@sha256:def456", nil
    }

    out, changed, err := dcv.Pin(input, resolver)
    if err != nil { t.Fatal(err) }
    if !changed { t.Fatalf("expected change") }
    if string(out) == string(input) {
        t.Fatalf("expected updated devcontainer.json")
    }
}

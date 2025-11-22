package cli_test

import (
	"imgpin/internal/cli"
	"testing"
)

func TestExtractDockerfileRefs(t *testing.T) {
	refs := cli.ExtractDockerfileRefs([]byte("FROM node:18\nFROM alpine:3.19"))
	if len(refs) != 2 {
		t.Fatalf("expected 2 refs")
	}
}

func TestExtractDevcontainerRefs(t *testing.T) {
	refs := cli.ExtractDevcontainerRefs([]byte(`{"image": "ubuntu:22.04"}`))
	if len(refs) != 1 {
		t.Fatalf("expected 1 image")
	}
}

func TestExtractKubernetesRefs(t *testing.T) {
	sample := []byte(`apiVersion: v1
kind: Pod
spec:
  containers:
  - image: busybox:1.36
`)
	refs := cli.ExtractKubernetesRefs(sample)
	if len(refs) != 1 {
		t.Fatalf("expected 1 image")
	}
}

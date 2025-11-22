package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"imgpin/internal/cli"
)

func fakeResolver(mapping map[string]string) func(string) (string, error) {
	return func(ref string) (string, error) {
		if dgst, ok := mapping[ref]; ok {
			return dgst, nil
		}
		return "", fmt.Errorf("unknown ref %s", ref)
	}
}

func TestRunResolveCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	restore := cli.SetResolver(fakeResolver(map[string]string{
		"python:3.11": "python@sha256:aaa",
	}))
	defer restore()

	code := run([]string{"resolve", "python:3.11"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, stderr.String())
	}
	if strings.TrimSpace(stdout.String()) != "python@sha256:aaa" {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
}

func TestRunFileInPlace(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "Dockerfile")
	if err := os.WriteFile(path, []byte("FROM node:18\n"), 0o644); err != nil {
		t.Fatalf("write dockerfile: %v", err)
	}

	restore := cli.SetResolver(fakeResolver(map[string]string{
		"node:18": "node@sha256:zzz",
	}))
	defer restore()

	var stdout, stderr bytes.Buffer
	code := run([]string{"file", path, "--in-place"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("file command failed: %s", stderr.String())
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read dockerfile: %v", err)
	}
	if !strings.Contains(string(data), "node@sha256:zzz") {
		t.Fatalf("expected digest in dockerfile: %s", data)
	}
}

func TestRunLockAndCheck(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "Dockerfile")
	if err := os.WriteFile(path, []byte("FROM alpine:3.19\n"), 0o644); err != nil {
		t.Fatalf("write dockerfile: %v", err)
	}
	lockPath := filepath.Join(tmp, "imgpin.lock")

	restoreLock := cli.SetResolver(fakeResolver(map[string]string{
		"alpine:3.19": "alpine@sha256:old",
	}))
	var stdout, stderr bytes.Buffer
	if code := run([]string{"lock", "--lockfile", lockPath, tmp}, &stdout, &stderr); code != 0 {
		t.Fatalf("lock failed: %s", stderr.String())
	}
	restoreLock()

	restoreCheck := cli.SetResolver(fakeResolver(map[string]string{
		"alpine:3.19": "alpine@sha256:old",
	}))
	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"check", "--lockfile", lockPath}, &stdout, &stderr); code != 0 {
		t.Fatalf("check should pass: %s", stderr.String())
	}
	restoreCheck()

	restoreDrift := cli.SetResolver(fakeResolver(map[string]string{
		"alpine:3.19": "alpine@sha256:new",
	}))
	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"check", "--lockfile", lockPath}, &stdout, &stderr); code == 0 {
		t.Fatalf("expected drift failure")
	}
	if !strings.Contains(stderr.String(), "drifted") {
		t.Fatalf("expected drift message, got %q", stderr.String())
	}
	restoreDrift()
}

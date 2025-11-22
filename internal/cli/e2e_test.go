package cli_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"imgpin/internal/cli"
)

func runCLI(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	cmd := cli.NewRootCommand()
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return stdout.String(), stderr.String(), err
}

func fakeResolver(mapping map[string]string) func(string) (string, error) {
	return func(ref string) (string, error) {
		if dgst, ok := mapping[ref]; ok {
			return dgst, nil
		}
		return "", fmt.Errorf("unknown ref %s", ref)
	}
}

func TestResolveCommandPrintsDigest(t *testing.T) {
	restore := cli.SetResolver(fakeResolver(map[string]string{
		"python:3.11": "python@sha256:aaa111",
	}))
	defer restore()

	stdout, _, err := runCLI(t, "resolve", "python:3.11")
	if err != nil {
		t.Fatalf("resolve command failed: %v", err)
	}
	if strings.TrimSpace(stdout) != "python@sha256:aaa111" {
		t.Fatalf("unexpected stdout: %q", stdout)
	}
}

func TestFileCommandDryRunDockerfile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "Dockerfile")
	err := os.WriteFile(path, []byte("FROM node:18\nRUN echo hi\nFROM alpine:3.19\n"), 0o644)
	if err != nil {
		t.Fatalf("write dockerfile: %v", err)
	}

	restore := cli.SetResolver(fakeResolver(map[string]string{
		"node:18":     "node@sha256:node18",
		"alpine:3.19": "alpine@sha256:alpine319",
	}))
	defer restore()

	stdout, _, err := runCLI(t, "file", path)
	if err != nil {
		t.Fatalf("file command failed: %v", err)
	}
	want := "FROM node@sha256:node18\nRUN echo hi\nFROM alpine@sha256:alpine319\n"
	if stdout != want {
		t.Fatalf("unexpected stdout:\nwant: %q\ngot : %q", want, stdout)
	}

	got, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("read dockerfile: %v", readErr)
	}
	if string(got) != "FROM node:18\nRUN echo hi\nFROM alpine:3.19\n" {
		t.Fatalf("file should remain unchanged for dry-run")
	}
}

func TestFileCommandInPlaceKubernetes(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "deploy.yaml")
	src := "apiVersion: v1\nspec:\n  template:\n    spec:\n      containers:\n      - image: repo/app:1.0\n"
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	restore := cli.SetResolver(fakeResolver(map[string]string{
		"repo/app:1.0": "repo/app@sha256:aaa",
	}))
	defer restore()

	stdout, _, err := runCLI(t, "file", path, "--in-place")
	if err != nil {
		t.Fatalf("file command failed: %v", err)
	}
	if stdout != "" {
		t.Fatalf("expected no stdout for in-place run, got %q", stdout)
	}

	updated, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("read manifest: %v", readErr)
	}
	if !strings.Contains(string(updated), "repo/app@sha256:aaa") {
		t.Fatalf("expected manifest to be updated:\n%s", updated)
	}
}

func TestResolveCommandPropagatesErrors(t *testing.T) {
	restore := cli.SetResolver(fakeResolver(map[string]string{}))
	defer restore()

	_, _, err := runCLI(t, "resolve", "unknown:tag")
	if err == nil {
		t.Fatalf("expected error for unknown ref")
	}
}

func TestLockCommandCreatesLockfile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "Dockerfile")
	if err := os.WriteFile(path, []byte("FROM node:20\n"), 0o644); err != nil {
		t.Fatalf("write dockerfile: %v", err)
	}
	lockPath := filepath.Join(tmp, "imgpin.lock")

	restore := cli.SetResolver(fakeResolver(map[string]string{
		"node:20": "node@sha256:lock",
	}))
	defer restore()

	stdout, _, err := runCLI(t, "lock", "--lockfile", lockPath, tmp)
	if err != nil {
		t.Fatalf("lock command failed: %v", err)
	}
	if !strings.Contains(stdout, "locked 1 image") {
		t.Fatalf("unexpected stdout: %q", stdout)
	}

	data, readErr := os.ReadFile(lockPath)
	if readErr != nil {
		t.Fatalf("read lockfile: %v", readErr)
	}
	if !strings.Contains(string(data), "node@sha256:lock") {
		t.Fatalf("expected digest in lockfile: %s", data)
	}
}

func TestCheckCommandDetectsDrift(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "Dockerfile")
	if err := os.WriteFile(path, []byte("FROM alpine:3.19\n"), 0o644); err != nil {
		t.Fatalf("write dockerfile: %v", err)
	}
	lockPath := filepath.Join(tmp, "imgpin.lock")

	restore := cli.SetResolver(fakeResolver(map[string]string{
		"alpine:3.19": "alpine@sha256:old",
	}))
	stdout, _, err := runCLI(t, "lock", "--lockfile", lockPath, tmp)
	restore()
	if err != nil {
		t.Fatalf("lock command failed: %v", err)
	}
	if stdout == "" {
		t.Fatalf("expected lock command output")
	}

	restoreNew := cli.SetResolver(fakeResolver(map[string]string{
		"alpine:3.19": "alpine@sha256:new",
	}))
	defer restoreNew()

	_, stderr, checkErr := runCLI(t, "check", "--lockfile", lockPath)
	if checkErr == nil {
		t.Fatalf("expected drift error")
	}
	if !strings.Contains(stderr, "drifted") {
		t.Fatalf("expected drift message, got %q", stderr)
	}
}

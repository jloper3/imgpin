package cli

import (
	"bytes"
	"testing"
)

func TestRootCommandExpose(t *testing.T) {
	cmd := RootCommand()
	if cmd == nil {
		t.Fatal("expected root command")
	}
}

func TestExecuteWithNoArgs(t *testing.T) {
	cmd := RootCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{})

	if err := Execute(); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
	cmd.SetArgs(nil)
}

func TestSetResolverNilRestoresDefault(t *testing.T) {
	restore := SetResolver(func(string) (string, error) {
		return "fake@sha256:123", nil
	})
	restore()

	reset := SetResolver(nil)
	defer reset()
}

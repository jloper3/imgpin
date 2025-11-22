package main

import (
	"fmt"
	"io"
	"os"

	"imgpin/internal/cli"
)

func run(args []string, stdout, stderr io.Writer) int {
	cmd := cli.RootCommand()
	cmd.SetArgs(args)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	err := cli.Execute()
	cmd.SetArgs(nil)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

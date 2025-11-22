package cli

import "github.com/spf13/cobra"

var rootCmd = buildRootCommand()

func buildRootCommand() *cobra.Command {
	cmd := newRootCommand()
	cmd.AddCommand(newResolveCommand())
	cmd.AddCommand(newFileCommand())
	cmd.AddCommand(newLockCommand())
	cmd.AddCommand(newCheckCommand())
	return cmd
}

func newRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:           "imgpin",
		Short:         "Digest pinning CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
}

// Execute runs the root command for the compiled binary.
func Execute() error {
	return rootCmd.Execute()
}

// RootCommand returns the singleton root command used by the binary.
func RootCommand() *cobra.Command {
	return rootCmd
}

// NewRootCommand constructs an isolated command tree for tests.
func NewRootCommand() *cobra.Command {
	return buildRootCommand()
}

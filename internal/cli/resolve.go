package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newResolveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "resolve <image>",
		Short: "Resolve an image tag to its digest",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, a []string) error {
			dg, err := resolveImage(a[0])
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(c.OutOrStdout(), dg)
			return err
		},
	}
}

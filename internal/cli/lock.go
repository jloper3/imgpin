package cli

import (
	"fmt"

	"imgpin/internal/lockfile"

	"github.com/spf13/cobra"
)

func newLockCommand() *cobra.Command {
	var lockPath string

	cmd := &cobra.Command{
		Use:   "lock [paths...]",
		Short: "Resolve image refs and write an imgpin.lock file",
		Args:  cobra.ArbitraryArgs,
		RunE: func(c *cobra.Command, args []string) error {
			refs, err := collectImageRefs(args)
			if err != nil {
				return err
			}
			if len(refs) == 0 {
				fmt.Fprintln(c.OutOrStdout(), "no image references found")
				return nil
			}

			lf, err := lockfile.Load(lockPath)
			if err != nil {
				return err
			}

			for _, ref := range refs {
				dg, err := resolveImage(ref)
				if err != nil {
					return err
				}
				lf.Set(ref, dg)
			}

			if err := lf.Save(); err != nil {
				return err
			}

			fmt.Fprintf(c.OutOrStdout(), "locked %d image(s)\n", len(refs))
			return nil
		},
	}

	cmd.Flags().StringVar(&lockPath, "lockfile", "imgpin.lock", "path to lockfile")
	return cmd
}

package cli

import (
	"fmt"
	"sort"

	"imgpin/internal/lockfile"

	"github.com/spf13/cobra"
)

func newCheckCommand() *cobra.Command {
	var lockPath string

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Verify the imgpin.lock file against current registry digests",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			lf, err := lockfile.Load(lockPath)
			if err != nil {
				return err
			}

			entries := lf.Entries()
			if len(entries) == 0 {
				return fmt.Errorf("no lockfile entries found, run imgpin lock")
			}

			var drift []string
			refs := make([]string, 0, len(entries))
			for ref := range entries {
				refs = append(refs, ref)
			}
			sort.Strings(refs)

			for _, ref := range refs {
				entry := entries[ref]
				dg, err := resolveImage(ref)
				if err != nil {
					return err
				}
				if dg != entry.Digest {
					drift = append(drift, fmt.Sprintf("%s drifted: lock=%s actual=%s", ref, entry.Digest, dg))
				}
			}

			if len(drift) > 0 {
				for _, msg := range drift {
					fmt.Fprintln(c.ErrOrStderr(), msg)
				}
				return fmt.Errorf("drift detected")
			}

			fmt.Fprintln(c.OutOrStdout(), "lockfile is up to date")
			return nil
		},
	}

	cmd.Flags().StringVar(&lockPath, "lockfile", "imgpin.lock", "path to lockfile")
	return cmd
}

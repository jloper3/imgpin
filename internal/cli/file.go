package cli

import (
	"fmt"
	"os"
	"path/filepath"

	dcv "imgpin/internal/handlers/devcontainer"
	dkh "imgpin/internal/handlers/dockerfile"
	k8s "imgpin/internal/handlers/kubernetes"

	"github.com/spf13/cobra"
)

func newFileCommand() *cobra.Command {
	var inPlace bool

	cmd := &cobra.Command{
		Use:   "file <path>",
		Short: "Rewrite supported files with pinned digests",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, a []string) error {
			path := a[0]
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			base := filepath.Base(path)
			var (
				out     []byte
				changed bool
			)

			switch {
			case base == "Dockerfile":
				out, changed, err = dkh.Pin(b, resolveImage)
			case base == "devcontainer.json":
				out, changed, err = dcv.Pin(b, resolveImage)
			case filepath.Ext(base) == ".yaml" || filepath.Ext(base) == ".yml":
				out, changed, err = k8s.Pin(b, resolveImage)
			default:
				return fmt.Errorf("unsupported file %q", path)
			}

			if err != nil {
				return err
			}

			if inPlace {
				if changed {
					return os.WriteFile(path, out, 0o644)
				}
				return nil
			}

			_, err = c.OutOrStdout().Write(out)
			return err
		},
	}

	cmd.Flags().BoolVar(&inPlace, "in-place", false, "rewrite file in place")
	return cmd
}

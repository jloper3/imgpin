package cli

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

func collectImageRefs(paths []string) ([]string, error) {
	if len(paths) == 0 {
		paths = []string{"."}
	}

	seen := map[string]struct{}{}
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			if err := filepath.WalkDir(p, func(path string, d fs.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				if d.IsDir() {
					return nil
				}
				return collectFromFile(path, seen)
			}); err != nil {
				return nil, err
			}
		} else {
			if err := collectFromFile(p, seen); err != nil {
				return nil, err
			}
		}
	}

	out := make([]string, 0, len(seen))
	for ref := range seen {
		out = append(out, ref)
	}
	sort.Strings(out)
	return out, nil
}

func collectFromFile(path string, seen map[string]struct{}) error {
	refs, err := refsForFile(path)
	if err != nil {
		return err
	}
	for _, ref := range refs {
		if ref == "" {
			continue
		}
		seen[ref] = struct{}{}
	}
	return nil
}

func refsForFile(path string) ([]string, error) {
	base := filepath.Base(path)
	switch {
	case base == "Dockerfile":
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return ExtractDockerfileRefs(b), nil
	case base == "devcontainer.json":
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return ExtractDevcontainerRefs(b), nil
	}

	switch filepath.Ext(base) {
	case ".yaml", ".yml":
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return ExtractKubernetesRefs(b), nil
	default:
		return nil, nil
	}
}

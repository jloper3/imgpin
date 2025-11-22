package resolve

import (
	"os"
	"path/filepath"
	"time"

	"imgpin/internal/cache"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

var (
	cachePath         = filepath.Join(os.Getenv("HOME"), ".cache", "imgpin", "digests.json")
	ttl               = 24 * time.Hour
	remoteFetchDigest = func(ref name.Reference) (string, error) {
		desc, err := remote.Get(ref)
		if err != nil {
			return "", err
		}
		return desc.Digest.String(), nil
	}
)

// Resolve converts a human-readable image reference to a fully qualified digest reference.
func Resolve(r string) (string, error) {
	c, _ := cache.Load(cachePath)
	if e, ok := c.Get(r); ok && time.Since(e.FetchedAt) < ttl {
		return e.Digest, nil
	}

	ref, err := name.ParseReference(r)
	if err != nil {
		return "", err
	}

	digest, err := remoteFetchDigest(ref)
	if err != nil {
		return "", err
	}

	dg := ref.Context().Name() + "@" + digest
	c.Set(r, dg)
	_ = c.Save()
	return dg, nil
}

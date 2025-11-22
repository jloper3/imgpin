package cli

import "imgpin/internal/resolve"

var resolveImage = resolve.Resolve

// SetResolver overrides the resolver used by CLI commands. It returns a
// restore function to reinstate the previous resolver.
func SetResolver(fn func(string) (string, error)) func() {
	prev := resolveImage
	if fn == nil {
		resolveImage = resolve.Resolve
	} else {
		resolveImage = fn
	}
	return func() {
		resolveImage = prev
	}
}

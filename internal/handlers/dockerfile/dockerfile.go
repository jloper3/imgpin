package dockerfile

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`(?m)^\s*FROM\s+([^\s]+)`)

func Pin(b []byte, res func(string) (string, error)) ([]byte, bool, error) {
	t := string(b)
	changed := false
	out := re.ReplaceAllStringFunc(t, func(l string) string {
		m := re.FindStringSubmatch(l)
		if len(m) < 2 {
			return l
		}
		ref := m[1]
		if strings.Contains(ref, "@sha256:") {
			return l
		}
		dg, err := res(ref)
		if err != nil {
			return l
		}
		changed = true
		return strings.Replace(l, ref, dg, 1)
	})
	return []byte(out), changed, nil
}

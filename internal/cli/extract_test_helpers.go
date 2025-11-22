package cli

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v3"

	"imgpin/internal/lockfile"
)

func ExtractDockerfileRefs(b []byte) []string {
	lines := strings.Split(string(b), "\n")
	out := []string{}
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if strings.HasPrefix(ln, "FROM ") {
			fields := strings.Fields(ln)
			if len(fields) > 1 {
				out = append(out, fields[1])
			}
		}
	}
	return out
}

func ExtractDevcontainerRefs(b []byte) []string {
	var obj map[string]interface{}
	_ = json.Unmarshal(b, &obj)
	var out []string
	if img, ok := obj["image"].(string); ok {
		out = append(out, img)
	}
	return out
}

func ExtractKubernetesRefs(b []byte) []string {
	var obj map[string]interface{}
	_ = yaml.Unmarshal(b, &obj)
	rs := []string{}
	recurse(obj, &rs)
	return rs
}

func recurse(n interface{}, out *[]string) {
	switch v := n.(type) {
	case map[string]interface{}:
		for k, val := range v {
			if k == "image" {
				if img, ok := val.(string); ok {
					*out = append(*out, img)
				}
			}
			recurse(val, out)
		}
	case []interface{}:
		for _, it := range v {
			recurse(it, out)
		}
	}
}

func driftExistsForTest(expected, actual map[string]lockfile.Entry) bool {
	for ref, exp := range expected {
		if act, ok := actual[ref]; !ok || act.Digest != exp.Digest {
			return true
		}
	}
	return false
}

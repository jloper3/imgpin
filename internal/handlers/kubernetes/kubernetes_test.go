package kubernetes_test

import (
	k8s "imgpin/internal/handlers/kubernetes"
	"testing"
)

func TestPinK8sYAML(t *testing.T) {
	input := []byte(`
apiVersion: v1
kind: Pod
spec:
  containers:
    - name: test
      image: nginx:latest
`)

	resolver := func(ref string) (string, error) {
		return "nginx@sha256:999aaa", nil
	}

	out, changed, err := k8s.Pin(input, resolver)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatalf("expected change")
	}
	if string(out) == string(input) {
		t.Fatalf("expected rewrite of image field")
	}
}

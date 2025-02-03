package integrations

import "testing"

func TestIntegrationFileEnding(t *testing.T) {
	for f, b := range integrations {
		if b[len(b)-1] != '\n' {
			t.Errorf(`Integration '%s' wasn't terminated with line feed (\n)`, f)
		}
	}
}

package virtualterm

import (
	"strings"
	"testing"
)

type testCondition struct {
	Screen    string
	Operation func(t *testing.T, term *Term)
	Expected  string
}

type testTerm struct {
	Tests     []testCondition
	Operation func(t *testing.T, term *Term)
}

func (tt *testTerm) RunTests(t *testing.T) {
	t.Helper()

	for i, test := range tt.Tests {
		term := NewTerminal(nil)
		term._lfEnabled = false

		for _, r := range test.Screen {
			if r == '\n' {
				continue
			}
			term.writeCell(r)
		}

		begin := term.exportAsString()

		if tt.Operation != nil {
			tt.Operation(t, term)
		}

		if test.Operation != nil {
			test.Operation(t, term)
		}

		expected := _pad(test.Expected)
		actual := strings.ReplaceAll(term.exportAsString(), "·", ".")
		if actual != expected {
			t.Errorf("Expected doesn't match Actual in test %d:", i)
			t.Logf("  Raw Begin: '%s'", strings.ReplaceAll(begin, "\n", "↲"))
			t.Logf("  Raw End:   '%s'", strings.ReplaceAll(actual, "\n", "↲"))
			t.Logf("  Screen:%s", _indent(_pad(test.Screen)))
			t.Logf("  Expected:%s", _indent(expected))
			t.Logf("  Actual:%s", _indent(actual))
		}
	}
}

func _pad(s string) string {
	padded := []byte("..........\n..........\n..........\n..........\n")
	copy(padded, s)
	return string(padded)
}

func _indent(s string) string {
	s = strings.ReplaceAll("\n"+s[:len(s)-1], "\n", "|\n    |")
	return s[1:] + "|"
}

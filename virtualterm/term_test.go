package virtualterm

import (
	"bytes"
	"strings"
	"testing"

	"github.com/lmorg/mxtty/types"
)

const _testTermHeight = 5

// NewTestTerminal creates a new virtual term used for unit tests
func NewTestTerminal() *Term {
	size := &types.XY{X: 10, Y: _testTermHeight}

	term := &Term{}

	term.reset(size)

	return term
}

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
		term := NewTestTerminal()
		term.setJumpScroll()

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

var _padding = bytes.Repeat([]byte("..........\n"), _testTermHeight)

func _pad(s string) string {
	padded := bytes.Clone(_padding)
	copy(padded, s)
	return string(padded)
}

func _indent(s string) string {
	s = strings.ReplaceAll("\n"+s[:len(s)-1], "\n", "|\n    |")
	return s[1:] + "|"
}

package psuedotty

import "testing"

func TestRuneLength(t *testing.T) {
	tests := []struct {
		Rune   string
		Length int
	}{
		{
			Rune:   "0",
			Length: 1,
		},
		{
			Rune:   "世",
			Length: 3,
		},
		{
			Rune:   "界",
			Length: 3,
		},
		{
			Rune:   "🤗",
			Length: 4,
		},
	}

	for i, test := range tests {
		act := runeLength(test.Rune[0])
		if act != test.Length {
			t.Errorf("Unexpected length of rune i test %d", i)
			t.Logf("  Rune:     %s", test.Rune)
			t.Logf("  Expected: %d", test.Length)
			t.Logf("  Actual:   %d", act)
		}
	}
}

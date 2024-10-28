package virtualterm

import (
	"testing"

	"github.com/lmorg/mxtty/types"
)

func TestCellGetSetElementXY(t *testing.T) {
	tests := []*types.XY{
		{0, 0},
		{1, 1},
		{3, 7},
		{7, 3},
		{200, 0},
		{0, 200},
		{200, 200},
		{10000, 13},
		{13, 10000},
		{10000, 10000},
		{32767, 1},
		{1, 32767},
		{32767, 32767},
	}

	for i, expected := range tests {
		r := setElementXY(expected)
		actual := getElementXY(r)
		if expected.X != actual.X || expected.Y != actual.Y {
			t.Errorf("Mismatch in test %d", i)
			t.Logf("Expected: X: %d, Y%d", expected.X, expected.Y)
			t.Logf("Actual:   X: %d, Y%d", actual.X, actual.Y)
		}
	}
}
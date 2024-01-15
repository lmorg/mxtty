package types

import (
	"encoding/json"
	"testing"
)

func TestApcSlice(t *testing.T) {
	tests := []struct {
		ApcCode string
		Slice   []string
	}{
		{
			`BEGIN`,
			[]string{"BEGIN"},
		},
		{
			`BEGIN;TABLE`,
			[]string{"BEGIN", "TABLE"},
		},
		{
			`BEGIN;TABLE;{"foo":"bar"}`,
			[]string{"BEGIN", "TABLE", `{"foo":"bar"}`},
		},
		{
			`BEGIN;TABLE;{"foo":"bar","baz":";"}`,
			[]string{"BEGIN", "TABLE", `{"foo":"bar","baz":";"}`},
		},
	}

	for i, test := range tests {
		apc := NewApcSlice([]rune(test.ApcCode))
		if lazyJson(apc.slice) != lazyJson(test.Slice) {
			t.Errorf("slice mismatch in test %d:", i)
			t.Logf("  APC str: '%s'", test.ApcCode)
			t.Logf("  Expected: %s", lazyJson(test.Slice))
			t.Logf("  Actual:   %s", lazyJson(apc.slice))
		}
	}
}

func lazyJson(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func TestApcSliceIndexJson(t *testing.T) {
	ApcCode := `BEGIN;TABLE;{"foo":"bar","baz":";"}`
	IndexElements := []string{"BEGIN", "TABLE", `{"foo":"bar","baz":";"}`}
	apc := NewApcSlice([]rune(ApcCode))

	for i, test := range IndexElements {
		if apc.Index(i) != test {
			t.Errorf("slice mismatch in test %d:", i)
			t.Logf("  APC str: '%s'", ApcCode)
			t.Logf("  Expected: %s", test)
			t.Logf("  Actual:   %s", apc.Index(i))
		}
	}
}

func TestApcSliceIndexNoJson(t *testing.T) {
	ApcCode := `BEGIN;TABLE`
	IndexElements := []string{"BEGIN", "TABLE"}
	apc := NewApcSlice([]rune(ApcCode))

	for i, test := range IndexElements {
		if apc.Index(i) != test {
			t.Errorf("slice mismatch in test %d:", i)
			t.Logf("  APC str: '%s'", ApcCode)
			t.Logf("  Expected: %s", test)
			t.Logf("  Actual:   %s", apc.Index(i))
		}
	}
}

func TestApcSliceParameter(t *testing.T) {
	ApcCode := `BEGIN;TABLE;{"foo":"bar","baz":";"}`
	Parameters := map[string]string{"foo": "bar", "baz": ";"}
	apc := NewApcSlice([]rune(ApcCode))

	for key, value := range Parameters {
		if apc.Parameter(key) != value {
			t.Errorf("slice mismatch in test '%s':", key)
			t.Logf("  APC str: '%s'", ApcCode)
			t.Logf("  Expected: %s", value)
			t.Logf("  Actual:   %s", apc.Parameter(key))
		}
	}
}

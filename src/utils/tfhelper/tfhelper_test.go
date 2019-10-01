package tfhelper

import (
	"testing"
)

func TestDiffFuncSupressCaseSensitivity(t *testing.T) {

	type testParams struct {
		first     string
		second    string
		different bool
	}

	tests := []testParams{
		{"hello", "HELLO", true},  // logically the same
		{"hello", "hElLo", true},  // logically the same
		{"hello", "world", false}, // logically different
		{"hello", "WORLD", false}, // logically different
	}

	for _, test := range tests {
		if test.different != DiffFuncSupressCaseSensitivity("", test.first, test.second, nil) {
			t.Errorf("%s compared to %s got %v, but expected %v", test.first, test.second, !test.different, test.different)
		}
	}
}

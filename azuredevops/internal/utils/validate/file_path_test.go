// +build all utils path

package validate

import (
	"fmt"
	"testing"
)

func TestPathValidation(t *testing.T) {
	type TestCase struct {
		Value    string
		TestName string
		ErrCount int
	}
	cases := []TestCase{
		{
			Value:    `\`,
			TestName: "Default Path",
			ErrCount: 0,
		},
		{
			Value:    "",
			TestName: "Empty Path",
			ErrCount: 1,
		},
		{
			Value:    "A",
			TestName: "Wrong Starting Character",
			ErrCount: 1,
		},
	}

	illegalChars := []string{"<", ">", "|", ":", "$", "@", `"`, "/", "%", "+", "*", "?"}
	for _, c := range illegalChars {
		cases = append(cases, TestCase{
			Value:    fmt.Sprintf(`\%s`, c),
			TestName: fmt.Sprintf("Illegal Character - %s", c),
			ErrCount: 1,
		})
	}

	for _, tc := range cases {
		t.Run(tc.TestName, func(t *testing.T) {
			_, errors := Path(tc.Value, tc.TestName)
			if len(errors) != tc.ErrCount {
				t.Fatalf("Expected TestPathValidation to have %d not %d errors for %q", tc.ErrCount, len(errors), tc.TestName)
			}
		})
	}
}

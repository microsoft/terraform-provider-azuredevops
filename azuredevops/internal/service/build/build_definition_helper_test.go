package build

import (
	"reflect"
	"testing"
)

func TestDateToDays(t *testing.T) {
	cases := []struct {
		Input  []interface{}
		Expect int
		Valid  bool
	}{
		{
			Input:  []interface{}{"Mon"},
			Expect: 1,
			Valid:  true,
		},
		{
			Input:  []interface{}{"Mon", "Tue"},
			Expect: 1,
			Valid:  false,
		},
		{
			Input:  []interface{}{"Mon", "Wed"},
			Expect: 5,
			Valid:  true,
		},
		{
			Input:  []interface{}{"Mon", "Tue", "Wed"},
			Expect: 7,
			Valid:  true,
		},
		{
			Input:  []interface{}{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
			Expect: 127,
			Valid:  true,
		},
		{
			Input:  []interface{}{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
			Expect: 128,
			Valid:  false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		days := DateToDays(tc.Input)
		valid := days == tc.Expect

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestDaysToDate(t *testing.T) {
	cases := []struct {
		Input  int
		Expect []string
		Valid  bool
	}{
		{
			Input:  1,
			Expect: []string{"Mon"},
			Valid:  true,
		},
		{
			Input:  1,
			Expect: []string{"Mon", "Tue"},
			Valid:  false,
		},
		{
			Input:  5,
			Expect: []string{"Mon", "Wed"},
			Valid:  true,
		},
		{
			Input:  7,
			Expect: []string{"Mon", "Tue", "Wed"},
			Valid:  true,
		},
		{
			Input:  127,
			Expect: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
			Valid:  true,
		},
		{
			Input:  128,
			Expect: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
			Valid:  false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %d", tc.Input)
		days := DaysToDate(tc.Input)
		valid := reflect.DeepEqual(days, tc.Expect)

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

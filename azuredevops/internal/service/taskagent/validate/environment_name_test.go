package validate

import "testing"

func TestEnvironmentName(t *testing.T) {
	testData := []struct {
		Value string
		Error bool
	}{
		{
			Value: "a1",
			Error: false,
		},
		{
			Value: "11",
			Error: false,
		},
		{
			Value: "1a",
			Error: false,
		},
		{
			Value: "aa",
			Error: false,
		},
		{
			Value: "1-1",
			Error: false,
		},
		{
			Value: "a",
			Error: false,
		},
		{
			Value: "1",
			Error: false,
		},
		{
			Value: "1-",
			Error: false,
		},
		{
			Value: "a-",
			Error: false,
		},
		{
			Value: "a1-",
			Error: false,
		},
		{
			Value: "1a--1-1-a-",
			Error: false,
		},
		{
			Value: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde1234abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde123",
			Error: false,
		},
		{
			Value: "[]",
			Error: true,
		},
		{
			Value: "\"/\\[]:|<>+=;,?*.",
			Error: true,
		},
		{
			Value: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde1234abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde1234",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Value)

		_, err := EnvironmentName(v.Value, "unit test")
		if err != nil && !v.Error {
			t.Fatalf("Expected pass but got an error: %s", err)
		}
	}
}

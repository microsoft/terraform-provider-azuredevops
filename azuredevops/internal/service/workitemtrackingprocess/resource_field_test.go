//go:build (all || resource_workitemtrackingprocess || resource_workitemtrackingprocess_field) && !exclude_resource_workitemtrackingprocess
// +build all resource_workitemtrackingprocess resource_workitemtrackingprocess_field
// +build !exclude_resource_workitemtrackingprocess

package workitemtrackingprocess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandDefaultValue(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name:     "valid JSON string",
			input:    `"hello"`,
			expected: "hello",
		},
		{
			name:     "valid JSON number",
			input:    `42`,
			expected: float64(42),
		},
		{
			name:     "valid JSON object",
			input:    `{"key": "value"}`,
			expected: map[string]interface{}{"key": "value"},
		},
		{
			name:     "valid JSON array",
			input:    `["a", "b"]`,
			expected: []interface{}{"a", "b"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:        "invalid JSON",
			input:       "not valid json",
			expectError: true,
			errorMsg:    "invalid JSON for default_value_json",
		},
		{
			name:        "invalid JSON syntax",
			input:       `{"key": }`,
			expectError: true,
			errorMsg:    "invalid JSON for default_value_json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expandDefaultValue(tc.input)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestFlattenDefaultValue(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    string
		expectError bool
		errorMsg    string
	}{
		{
			name:     "nil value",
			input:    nil,
			expected: "",
		},
		{
			name:     "string value",
			input:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "number value",
			input:    42,
			expected: "42",
		},
		{
			name:     "map value",
			input:    map[string]interface{}{"key": "value"},
			expected: `{"key":"value"}`,
		},
		{
			name:     "array value",
			input:    []string{"a", "b"},
			expected: `["a","b"]`,
		},
		{
			name:     "bool value",
			input:    true,
			expected: "true",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := flattenDefaultValue(tc.input)
			if tc.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

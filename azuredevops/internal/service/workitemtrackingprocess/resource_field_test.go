//go:build (all || resource_workitemtrackingprocess || resource_workitemtrackingprocess_field) && !exclude_resource_workitemtrackingprocess
// +build all resource_workitemtrackingprocess resource_workitemtrackingprocess_field
// +build !exclude_resource_workitemtrackingprocess

package workitemtrackingprocess

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func getFieldResourceData(t *testing.T, input map[string]interface{}) *schema.ResourceData {
	r := ResourceField()
	return schema.TestResourceDataRaw(t, r.Schema, input)
}

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

func TestField_Create_InvalidDefaultValueJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	processId := uuid.New()
	d := getFieldResourceData(t, map[string]interface{}{
		"process_id":              processId.String(),
		"work_item_type_ref_name": "MyWorkItemType",
		"reference_name":          "Custom.MyField",
		"default_value_json":      "not valid json",
	})

	diags := resourceFieldCreate(context.Background(), d, clients)

	assert.NotEmpty(t, diags)
	assert.Contains(t, diags[0].Summary, "invalid JSON for default_value_json")
}

func TestField_Update_InvalidDefaultValueJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	processId := uuid.New()
	d := getFieldResourceData(t, map[string]interface{}{
		"process_id":              processId.String(),
		"work_item_type_ref_name": "MyWorkItemType",
		"reference_name":          "Custom.MyField",
		"default_value_json":      "not valid json",
	})
	d.SetId("Custom.MyField")

	diags := resourceFieldUpdate(context.Background(), d, clients)

	assert.NotEmpty(t, diags)
	assert.Contains(t, diags[0].Summary, "invalid JSON for default_value_json")
}

func TestField_Read_InvalidDefaultValueFromAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	processId := uuid.New()
	fieldRefName := "Custom.MyField"
	fieldName := "MyField"
	// Create a value that cannot be marshaled to JSON (channel)
	unmarshalableValue := make(chan int)

	mockClient.EXPECT().GetWorkItemTypeField(clients.Ctx, gomock.Any()).Return(
		&workitemtrackingprocess.ProcessWorkItemTypeField{
			ReferenceName: &fieldRefName,
			Name:          &fieldName,
			DefaultValue:  unmarshalableValue,
		}, nil,
	).Times(1)

	d := getFieldResourceData(t, map[string]interface{}{
		"process_id":              processId.String(),
		"work_item_type_ref_name": "MyWorkItemType",
		"reference_name":          fieldRefName,
	})
	d.SetId(fieldRefName)

	diags := resourceFieldRead(context.Background(), d, clients)

	assert.NotEmpty(t, diags)
	assert.Contains(t, diags[0].Summary, "failed to marshal default_value")
}

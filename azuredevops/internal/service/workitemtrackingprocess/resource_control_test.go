//go:build (all || resource_workitemtrackingprocess || resource_workitemtrackingprocess_control) && !exclude_resource_workitemtrackingprocess
// +build all resource_workitemtrackingprocess resource_workitemtrackingprocess_control
// +build !exclude_resource_workitemtrackingprocess

package workitemtrackingprocess

import (
	"testing"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
)

func TestExpandContribution(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected *workitemtrackingprocess.WitContribution
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty input",
			input:    []interface{}{},
			expected: nil,
		},
		{
			name:     "nil first element",
			input:    []interface{}{nil},
			expected: nil,
		},
		{
			name: "contribution_id only",
			input: []interface{}{
				map[string]interface{}{
					"contribution_id": "ms-devlabs.vsts-extensions-multivalue-control.multivalue-form-control",
				},
			},
			expected: &workitemtrackingprocess.WitContribution{
				ContributionId: converter.String("ms-devlabs.vsts-extensions-multivalue-control.multivalue-form-control"),
			},
		},
		{
			name: "all fields",
			input: []interface{}{
				map[string]interface{}{
					"contribution_id":           "ms-devlabs.vsts-extensions-multivalue-control.multivalue-form-control",
					"height":                    50,
					"show_on_deleted_work_item": true,
					"inputs": map[string]interface{}{
						"FieldName": "System.Tags",
						"Values":    "Option1;Option2;Option3",
					},
				},
			},
			expected: &workitemtrackingprocess.WitContribution{
				ContributionId:        converter.String("ms-devlabs.vsts-extensions-multivalue-control.multivalue-form-control"),
				Height:                converter.Int(50),
				ShowOnDeletedWorkItem: converter.Bool(true),
				Inputs: &map[string]interface{}{
					"FieldName": "System.Tags",
					"Values":    "Option1;Option2;Option3",
				},
			},
		},
		{
			name: "empty contribution_id is ignored",
			input: []interface{}{
				map[string]interface{}{
					"contribution_id": "",
					"height":          50,
				},
			},
			expected: &workitemtrackingprocess.WitContribution{
				Height: converter.Int(50),
			},
		},
		{
			name: "zero height is ignored",
			input: []interface{}{
				map[string]interface{}{
					"contribution_id": "my-contribution",
					"height":          0,
				},
			},
			expected: &workitemtrackingprocess.WitContribution{
				ContributionId: converter.String("my-contribution"),
			},
		},
		{
			name: "empty inputs is ignored",
			input: []interface{}{
				map[string]interface{}{
					"contribution_id": "my-contribution",
					"inputs":          map[string]interface{}{},
				},
			},
			expected: &workitemtrackingprocess.WitContribution{
				ContributionId: converter.String("my-contribution"),
			},
		},
		{
			name: "show_on_deleted_work_item false is set",
			input: []interface{}{
				map[string]interface{}{
					"contribution_id":           "my-contribution",
					"show_on_deleted_work_item": false,
				},
			},
			expected: &workitemtrackingprocess.WitContribution{
				ContributionId:        converter.String("my-contribution"),
				ShowOnDeletedWorkItem: converter.Bool(false),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandContribution(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFlattenContribution(t *testing.T) {
	tests := []struct {
		name     string
		input    *workitemtrackingprocess.WitContribution
		expected []interface{}
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name: "contribution_id only",
			input: &workitemtrackingprocess.WitContribution{
				ContributionId: converter.String("ms-devlabs.vsts-extensions-multivalue-control.multivalue-form-control"),
			},
			expected: []interface{}{
				map[string]interface{}{
					"contribution_id": "ms-devlabs.vsts-extensions-multivalue-control.multivalue-form-control",
				},
			},
		},
		{
			name: "all fields",
			input: &workitemtrackingprocess.WitContribution{
				ContributionId:        converter.String("ms-devlabs.vsts-extensions-multivalue-control.multivalue-form-control"),
				Height:                converter.Int(50),
				ShowOnDeletedWorkItem: converter.Bool(true),
				Inputs: &map[string]interface{}{
					"FieldName": "System.Tags",
					"Values":    "Option1;Option2;Option3",
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"contribution_id":           "ms-devlabs.vsts-extensions-multivalue-control.multivalue-form-control",
					"height":                    50,
					"show_on_deleted_work_item": true,
					"inputs": map[string]string{
						"FieldName": "System.Tags",
						"Values":    "Option1;Option2;Option3",
					},
				},
			},
		},
		{
			name:  "empty contribution",
			input: &workitemtrackingprocess.WitContribution{},
			expected: []interface{}{
				map[string]interface{}{},
			},
		},
		{
			name: "inputs with non-string values are filtered",
			input: &workitemtrackingprocess.WitContribution{
				ContributionId: converter.String("my-contribution"),
				Inputs: &map[string]interface{}{
					"StringVal": "hello",
					"IntVal":    123,
					"BoolVal":   true,
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"contribution_id": "my-contribution",
					"inputs": map[string]string{
						"StringVal": "hello",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenContribution(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFindControlInGroup(t *testing.T) {
	controlId := "System.Title"
	otherControlId := "System.Description"

	tests := []struct {
		name      string
		group     *workitemtrackingprocess.Group
		controlId string
		expected  *workitemtrackingprocess.Control
	}{
		{
			name: "found control",
			group: &workitemtrackingprocess.Group{
				Controls: &[]workitemtrackingprocess.Control{
					{Id: &controlId},
				},
			},
			controlId: controlId,
			expected:  &workitemtrackingprocess.Control{Id: &controlId},
		},
		{
			name: "found among multiple controls",
			group: &workitemtrackingprocess.Group{
				Controls: &[]workitemtrackingprocess.Control{
					{Id: &otherControlId},
					{Id: &controlId},
				},
			},
			controlId: controlId,
			expected:  &workitemtrackingprocess.Control{Id: &controlId},
		},
		{
			name: "not found",
			group: &workitemtrackingprocess.Group{
				Controls: &[]workitemtrackingprocess.Control{
					{Id: &otherControlId},
				},
			},
			controlId: controlId,
			expected:  nil,
		},
		{
			name: "nil controls",
			group: &workitemtrackingprocess.Group{
				Controls: nil,
			},
			controlId: controlId,
			expected:  nil,
		},
		{
			name: "empty controls",
			group: &workitemtrackingprocess.Group{
				Controls: &[]workitemtrackingprocess.Control{},
			},
			controlId: controlId,
			expected:  nil,
		},
		{
			name: "control with nil id",
			group: &workitemtrackingprocess.Group{
				Controls: &[]workitemtrackingprocess.Control{
					{Id: nil},
					{Id: &controlId},
				},
			},
			controlId: controlId,
			expected:  &workitemtrackingprocess.Control{Id: &controlId},
		},
		{
			name: "only control with nil id",
			group: &workitemtrackingprocess.Group{
				Controls: &[]workitemtrackingprocess.Control{
					{Id: nil},
				},
			},
			controlId: controlId,
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findControlInGroup(tt.group, tt.controlId)
			assert.Equal(t, tt.expected, result)
		})
	}
}

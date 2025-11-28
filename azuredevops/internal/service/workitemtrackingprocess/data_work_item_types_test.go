//go:build (all || data_sources || data_workitemtrackingprocess || data_workitemtrackingprocess_workitemtypes) && !exclude_data_sources
// +build all data_sources data_workitemtrackingprocess data_workitemtrackingprocess_workitemtypes
// +build !exclude_data_sources

package workitemtrackingprocess

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func toWorkItemTypeMap(wit *workitemtrackingprocess.ProcessWorkItemType) map[string]any {
	return map[string]any{
		"reference_name": *wit.ReferenceName,
		"name":           *wit.Name,
		"description":    converter.ToString(wit.Description, ""),
		"color":          "#" + *wit.Color,
		"icon":           *wit.Icon,
		"is_disabled":    *wit.IsDisabled,
		"inherits_from":  converter.ToString(wit.Inherits, ""),
		"customization":  string(*wit.Customization),
		"url":            *wit.Url,
	}
}

func TestDataWorkItemTypes_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	processId := "59788636-ed1e-4e20-a7d1-93ee382beba7"
	workItemType1 := createProcessWorkItemType("Custom.WorkItemType1")
	workItemType2 := createProcessWorkItemType("Custom.WorkItemType2")
	workItemType2.Description = nil

	testCases := []struct {
		name                  string
		input                 map[string]any
		workItemTypesToReturn []workitemtrackingprocess.ProcessWorkItemType
		returnError           error
		expectedWorkItemTypes []map[string]any
		expectedProcessId     string
	}{
		{
			name: "success",
			input: map[string]any{
				"process_id": processId,
			},
			workItemTypesToReturn: []workitemtrackingprocess.ProcessWorkItemType{*workItemType1, *workItemType2},
			expectedWorkItemTypes: []map[string]any{
				toWorkItemTypeMap(workItemType1),
				toWorkItemTypeMap(workItemType2),
			},
			expectedProcessId: processId,
		},
		{
			name: "error",
			input: map[string]any{
				"process_id": processId,
			},
			returnError: errors.New("GetProcessWorkItemTypes failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resourceData := schema.TestResourceDataRaw(t, DataWorkItemTypes().Schema, tc.input)

			ctx := context.Background()
			mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
			clients := &client.AggregatedClient{
				WorkItemTrackingProcessClient: mockClient,
				Ctx:                           ctx,
			}

			mockClient.
				EXPECT().
				GetProcessWorkItemTypes(ctx, gomock.Any()).
				DoAndReturn(
					func(ctx context.Context, args workitemtrackingprocess.GetProcessWorkItemTypesArgs) (*[]workitemtrackingprocess.ProcessWorkItemType, error) {
						assert.Equal(t, processId, args.ProcessId.String())
						assert.Equal(t, workitemtrackingprocess.GetWorkItemTypeExpandValues.None, *args.Expand)

						if tc.returnError != nil {
							return nil, tc.returnError
						}

						return &tc.workItemTypesToReturn, nil
					},
				).
				Times(1)

			err := readWorkItemTypes(ctx, resourceData, clients)

			if tc.returnError != nil {
				assert.True(t, err.HasError())
				return
			}
			assert.False(t, err.HasError())
			assert.Equal(t, tc.expectedProcessId, resourceData.Id())

			actualWorkItemTypes := resourceData.Get("work_item_types").(*schema.Set).List()

			// Convert to maps for comparison
			actualMaps := make([]map[string]any, len(actualWorkItemTypes))
			for i, wit := range actualWorkItemTypes {
				actualMaps[i] = wit.(map[string]any)
			}

			diffOptions := []cmp.Option{
				cmpopts.EquateEmpty(),
				cmpopts.SortSlices(func(a, b map[string]any) bool {
					return a["reference_name"].(string) < b["reference_name"].(string)
				}),
			}

			if diff := cmp.Diff(tc.expectedWorkItemTypes, actualMaps, diffOptions...); diff != "" {
				t.Errorf("Work item types mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

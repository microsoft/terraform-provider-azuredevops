//go:build (all || data_sources || data_workitemtrackingprocess || data_workitemtrackingprocess_workitemtypes) && !exclude_data_sources
// +build all data_sources data_workitemtrackingprocess data_workitemtrackingprocess_workitemtypes
// +build !exclude_data_sources

package workitemtrackingprocess

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDataWorkItemTypes_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	processId := "59788636-ed1e-4e20-a7d1-93ee382beba7"
	workItemType1 := createProcessWorkItemType("Custom.WorkItemType1")
	workItemType2 := createProcessWorkItemType("Custom.WorkItemType2")
	emptyWorkItemType := createEmptyProcessWorkItemType("Empty.WorkItemType")

	testCases := []struct {
		name                  string
		input                 map[string]any
		workItemTypesToReturn []workitemtrackingprocess.ProcessWorkItemType
		returnError           error
		expectedReturn        map[string]string
	}{
		{
			name: "success",
			input: map[string]any{
				"process_id": processId,
			},
			workItemTypesToReturn: []workitemtrackingprocess.ProcessWorkItemType{*workItemType1, *workItemType2, *emptyWorkItemType},
			expectedReturn: map[string]string{
				"id":                processId,
				"process_id":        processId,
				"work_item_types.#": "3",

				"work_item_types.3092183233.description":    *workItemType1.Description,
				"work_item_types.3092183233.inherits_from":  *workItemType1.Inherits,
				"work_item_types.3092183233.reference_name": *workItemType1.ReferenceName,
				"work_item_types.3092183233.name":           *workItemType1.Name,
				"work_item_types.3092183233.color":          "#" + *workItemType1.Color,
				"work_item_types.3092183233.icon":           *workItemType1.Icon,
				"work_item_types.3092183233.is_disabled":    strconv.FormatBool(*workItemType1.IsDisabled),
				"work_item_types.3092183233.customization":  string(*workItemType1.Customization),
				"work_item_types.3092183233.url":            *workItemType1.Url,

				"work_item_types.558344571.description":    *workItemType2.Description,
				"work_item_types.558344571.inherits_from":  *workItemType2.Inherits,
				"work_item_types.558344571.reference_name": *workItemType2.ReferenceName,
				"work_item_types.558344571.name":           *workItemType2.Name,
				"work_item_types.558344571.color":          "#" + *workItemType2.Color,
				"work_item_types.558344571.icon":           *workItemType2.Icon,
				"work_item_types.558344571.is_disabled":    strconv.FormatBool(*workItemType2.IsDisabled),
				"work_item_types.558344571.customization":  string(*workItemType2.Customization),
				"work_item_types.558344571.url":            *workItemType2.Url,

				"work_item_types.3142816001.color":          "",
				"work_item_types.3142816001.customization":  "",
				"work_item_types.3142816001.description":    "",
				"work_item_types.3142816001.icon":           "",
				"work_item_types.3142816001.inherits_from":  "",
				"work_item_types.3142816001.is_disabled":    "false",
				"work_item_types.3142816001.name":           "",
				"work_item_types.3142816001.reference_name": *emptyWorkItemType.ReferenceName,
				"work_item_types.3142816001.url":            "",
			},
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

			diffOptions := []cmp.Option{
				cmpopts.EquateEmpty(),
			}

			if diff := cmp.Diff(tc.expectedReturn, resourceData.State().Attributes, diffOptions...); diff != "" {
				t.Errorf("Work item types mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

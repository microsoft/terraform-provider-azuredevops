//go:build (all || data_sources || data_workitemtrackingprocess || data_workitemtrackingprocess_workitemtype) && !exclude_data_sources
// +build all data_sources data_workitemtrackingprocess data_workitemtrackingprocess_workitemtype
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
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func getDataWorkItemTypeResourceData(t *testing.T, input map[string]interface{}) *schema.ResourceData {
	r := DataWorkItemType()
	return schema.TestResourceDataRaw(t, r.Schema, input)
}

func TestDataWorkItemType_Get(t *testing.T) {
	processId := "59788636-ed1e-4e20-a7d1-93ee382beba7"
	referenceName := "Custom.WorkItemType"
	workItemType := createProcessWorkItemType(referenceName)

	testCases := []struct {
		name                 string
		input                map[string]any
		workItemTypeToReturn workitemtrackingprocess.ProcessWorkItemType
		returnError          error
		expectedReturn       map[string]string
	}{
		{
			name: "success",
			input: map[string]any{
				"process_id":     processId,
				"reference_name": referenceName,
			},
			workItemTypeToReturn: *workItemType,
			expectedReturn: map[string]string{
				"id":             referenceName,
				"process_id":     processId,
				"reference_name": referenceName,
				"name":           *workItemType.Name,
				"description":    *workItemType.Description,
				"color":          "#" + *workItemType.Color,
				"icon":           *workItemType.Icon,
				"is_disabled":    strconv.FormatBool(*workItemType.IsDisabled),
				"inherits_from":  *workItemType.Inherits,
				"customization":  string(*workItemType.Customization),
				"url":            *workItemType.Url,
			},
		},
		{
			name: "error from API",
			input: map[string]any{
				"process_id":     processId,
				"reference_name": referenceName,
			},
			returnError: errors.New("api failure"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
			clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: ctx}

			mockClient.EXPECT().
				GetProcessWorkItemType(ctx, gomock.Any()).
				DoAndReturn(
					func(ctx context.Context, args workitemtrackingprocess.GetProcessWorkItemTypeArgs) (*workitemtrackingprocess.ProcessWorkItemType, error) {
						assert.Equal(t, processId, args.ProcessId.String())
						assert.Equal(t, referenceName, *args.WitRefName)
						assert.Equal(t, workitemtrackingprocess.GetWorkItemTypeExpandValues.None, *args.Expand)

						if tc.returnError != nil {
							return nil, tc.returnError
						}

						return converter.ToPtr(tc.workItemTypeToReturn), nil
					},
				).
				Times(1)

			d := getDataWorkItemTypeResourceData(t, tc.input)

			err := readDataWorkItemType(ctx, d, clients)

			if tc.returnError != nil {
				assert.True(t, err.HasError())
				return
			}
			assert.False(t, err.HasError())

			diffOptions := []cmp.Option{
				cmpopts.EquateEmpty(),
			}
			if diff := cmp.Diff(tc.expectedReturn, d.State().Attributes, diffOptions...); diff != "" {
				t.Errorf("Resource data attributes mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func createProcessWorkItemType(referenceName string) *workitemtrackingprocess.ProcessWorkItemType {
	return &workitemtrackingprocess.ProcessWorkItemType{
		ReferenceName: converter.String(referenceName),
		Name:          converter.String("Work Item Type"),
		Description:   converter.String("A custom work item type"),
		Color:         converter.String("009CCC"),
		Icon:          converter.String("icon_clipboard"),
		IsDisabled:    converter.Bool(false),
		Inherits:      converter.String("System.WorkItemType"),
		Customization: &workitemtrackingprocess.CustomizationTypeValues.Custom,
		Url:           converter.String("https://dev.azure.com/org/_apis/work/processes/process-id/workitemtypes/Custom.WorkItemType"),
	}
}

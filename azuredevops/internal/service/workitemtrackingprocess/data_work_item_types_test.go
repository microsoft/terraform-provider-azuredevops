//go:build (all || data_sources || data_workitemtrackingprocess || data_workitemtrackingprocess_workitemtypes) && !exclude_data_sources
// +build all data_sources data_workitemtrackingprocess data_workitemtrackingprocess_workitemtypes
// +build !exclude_data_sources

package workitemtrackingprocess

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
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

func getValueOrDefault[T any](ptr *T, defaultValue T) T {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

func toExpectedWorkItemTypes(processId string, wits ...*workitemtrackingprocess.ProcessWorkItemType) map[string]string {
	m := map[string]string{
		"id":                processId,
		"process_id":        processId,
		"work_item_types.#": fmt.Sprintf("%d", len(wits)),
	}

	for _, wit := range wits {
		setWitAttribute := func(name string, value string) {
			m[fmt.Sprintf("work_item_types.%d.%s", crc32.ChecksumIEEE([]byte(*wit.ReferenceName)), name)] = value
		}

		setWitAttribute("description", getValueOrDefault(wit.Description, ""))
		setWitAttribute("inherits_from", getValueOrDefault(wit.Inherits, ""))
		setWitAttribute("reference_name", getValueOrDefault(wit.ReferenceName, ""))
		setWitAttribute("name", getValueOrDefault(wit.Name, ""))
		if wit.Color != nil {
			setWitAttribute("color", "#"+*wit.Color)
		} else {
			setWitAttribute("color", "")
		}
		setWitAttribute("icon", getValueOrDefault(wit.Icon, ""))
		setWitAttribute("is_disabled", strconv.FormatBool(getValueOrDefault(wit.IsDisabled, false)))
		if wit.Customization != nil {
			setWitAttribute("customization", string(*wit.Customization))
		} else {
			setWitAttribute("customization", "")
		}
		setWitAttribute("url", getValueOrDefault(wit.Url, ""))
	}

	return m
}

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
			expectedReturn:        toExpectedWorkItemTypes(processId, workItemType1, workItemType2, emptyWorkItemType),
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

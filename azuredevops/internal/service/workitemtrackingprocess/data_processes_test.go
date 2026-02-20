//go:build (all || data_sources || data_workitemtrackingprocess || data_workitemtrackingprocess_processes) && !exclude_data_sources
// +build all data_sources data_workitemtrackingprocess data_workitemtrackingprocess_processes
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

func getDataProcessesResourceData(t *testing.T, input map[string]interface{}) *schema.ResourceData {
	r := DataProcesses()
	return schema.TestResourceDataRaw(t, r.Schema, input)
}

func TestDataProcesses_ListProcesses(t *testing.T) {
	process1 := createProcessInfo("59788636-ed1e-4e20-a7d1-93ee382beba7", true)
	process2 := createProcessInfo("2166b5b0-8b17-4c9d-9360-d46526a021bf", false)
	minimalProcessWithExpand := createMinimalProcessInfo("57ad2818-2b64-4192-9f4d-2207dd8b0553", true)

	testCases := []struct {
		name              string
		input             map[string]any
		processesToReturn []workitemtrackingprocess.ProcessInfo
		returnError       error
		expectedReturn    map[string]string
	}{
		{
			name:              "success, no expand",
			input:             map[string]any{},
			processesToReturn: []workitemtrackingprocess.ProcessInfo{*process2},
			expectedReturn: map[string]string{
				"expand":                             "none",
				"id":                                 "none",
				"processes.#":                        "1",
				"processes.0.id":                     process2.TypeId.String(),
				"processes.0.name":                   *process2.Name,
				"processes.0.description":            *process2.Description,
				"processes.0.customization_type":     string(*process2.CustomizationType),
				"processes.0.parent_process_type_id": process2.ParentProcessTypeId.String(),
				"processes.0.is_default":             strconv.FormatBool(*process2.IsDefault),
				"processes.0.is_enabled":             strconv.FormatBool(*process2.IsEnabled),
				"processes.0.projects.#":             "0",
				"processes.0.reference_name":         "",
			},
		},
		{
			name: "success, with expand",
			input: map[string]any{
				"expand": "projects",
			},
			processesToReturn: []workitemtrackingprocess.ProcessInfo{*process1, *process2},
			expectedReturn: map[string]string{
				"expand":      "projects",
				"id":          "projects",
				"processes.#": "2",

				"processes.0.id":                     process1.TypeId.String(),
				"processes.0.name":                   *process1.Name,
				"processes.0.description":            *process1.Description,
				"processes.0.parent_process_type_id": process1.ParentProcessTypeId.String(),
				"processes.0.is_default":             strconv.FormatBool(*process1.IsDefault),
				"processes.0.is_enabled":             strconv.FormatBool(*process1.IsEnabled),
				"processes.0.customization_type":     string(*process1.CustomizationType),
				"processes.0.reference_name":         "",
				"processes.0.projects.#":             "1",
				"processes.0.projects.0.id":          (*process1.Projects)[0].Id.String(),
				"processes.0.projects.0.name":        *(*process1.Projects)[0].Name,
				"processes.0.projects.0.description": *(*process1.Projects)[0].Description,
				"processes.0.projects.0.url":         *(*process1.Projects)[0].Url,

				"processes.1.id":                     process2.TypeId.String(),
				"processes.1.name":                   *process2.Name,
				"processes.1.description":            *process2.Description,
				"processes.1.parent_process_type_id": process2.ParentProcessTypeId.String(),
				"processes.1.is_default":             strconv.FormatBool(*process2.IsDefault),
				"processes.1.is_enabled":             strconv.FormatBool(*process2.IsEnabled),
				"processes.1.customization_type":     string(*process2.CustomizationType),
				"processes.1.projects.#":             "0",
				"processes.1.reference_name":         "",
			},
		},
		{
			name: "success, with minimal attributes",
			input: map[string]any{
				"expand": "projects",
			},
			processesToReturn: []workitemtrackingprocess.ProcessInfo{*minimalProcessWithExpand},
			expectedReturn: map[string]string{
				"expand":      "projects",
				"id":          "projects",
				"processes.#": "1",

				"processes.0.id":                     minimalProcessWithExpand.TypeId.String(),
				"processes.0.name":                   "",
				"processes.0.description":            "",
				"processes.0.parent_process_type_id": "",
				"processes.0.is_default":             "false",
				"processes.0.is_enabled":             "false",
				"processes.0.customization_type":     "",
				"processes.0.reference_name":         "",
				"processes.0.projects.#":             "1",
				"processes.0.projects.0.id":          (*minimalProcessWithExpand.Projects)[0].Id.String(),
				"processes.0.projects.0.name":        "",
				"processes.0.projects.0.description": "",
				"processes.0.projects.0.url":         "",
			},
		},
		{
			name:        "error from API",
			input:       map[string]any{},
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
				GetListOfProcesses(ctx, gomock.Any()).
				DoAndReturn(
					func(ctx context.Context, args workitemtrackingprocess.GetListOfProcessesArgs) (*[]workitemtrackingprocess.ProcessInfo, error) {
						if expand, expandFound := tc.input["expand"]; expandFound {
							assert.Equal(t, expand.(string), string(*args.Expand))
						} else {
							assert.Equal(t, workitemtrackingprocess.GetProcessExpandLevelValues.None, *args.Expand)
						}

						if tc.returnError != nil {
							return nil, tc.returnError
						}

						return converter.ToPtr(tc.processesToReturn), nil
					},
				).
				Times(1)

			d := getDataProcessesResourceData(t, tc.input)

			err := readProcesses(ctx, d, clients)

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

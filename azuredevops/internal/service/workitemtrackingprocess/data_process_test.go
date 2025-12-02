//go:build (all || data_sources || data_workitemtrackingprocess || data_workitemtrackingprocess_process) && !exclude_data_sources
// +build all data_sources data_workitemtrackingprocess data_workitemtrackingprocess_process
// +build !exclude_data_sources

package workitemtrackingprocess

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func getDataProcessResourceData(t *testing.T, input map[string]interface{}) *schema.ResourceData {
	r := DataProcess()
	return schema.TestResourceDataRaw(t, r.Schema, input)
}

func TestDataProcess_Get(t *testing.T) {
	id := "59788636-ed1e-4e20-a7d1-93ee382beba7"
	processWithExpand := createProcessInfo(id, true)
	processWithoutExpand := createProcessInfo(id, false)
	minimalProcessWithExpand := createMinimalProcessInfo(id, true)

	testCases := []struct {
		name            string
		input           map[string]any
		processToReturn workitemtrackingprocess.ProcessInfo
		returnError     error
		expectedReturn  map[string]string
	}{
		{
			name: "success, no expand",
			input: map[string]any{
				"id": id,
			},
			processToReturn: *processWithoutExpand,
			expectedReturn: map[string]string{
				"expand":                 "none",
				"id":                     processWithoutExpand.TypeId.String(),
				"name":                   *processWithoutExpand.Name,
				"description":            *processWithoutExpand.Description,
				"customization_type":     string(*processWithoutExpand.CustomizationType),
				"parent_process_type_id": processWithoutExpand.ParentProcessTypeId.String(),
				"is_default":             strconv.FormatBool(*processWithoutExpand.IsDefault),
				"is_enabled":             strconv.FormatBool(*processWithoutExpand.IsEnabled),
				"projects.#":             "0",
			},
		},
		{
			name: "success, with expand",
			input: map[string]any{
				"id":     id,
				"expand": "projects",
			},
			processToReturn: *processWithExpand,
			expectedReturn: map[string]string{
				"expand": "projects",

				"id":                              processWithExpand.TypeId.String(),
				"name":                            *processWithExpand.Name,
				"description":                     *processWithExpand.Description,
				"parent_process_type_id":          processWithExpand.ParentProcessTypeId.String(),
				"is_default":                      strconv.FormatBool(*processWithExpand.IsDefault),
				"is_enabled":                      strconv.FormatBool(*processWithExpand.IsEnabled),
				"customization_type":              string(*processWithExpand.CustomizationType),
				"projects.#":                      "1",
				"projects.1967634730.id":          (*processWithExpand.Projects)[0].Id.String(),
				"projects.1967634730.name":        *(*processWithExpand.Projects)[0].Name,
				"projects.1967634730.description": *(*processWithExpand.Projects)[0].Description,
				"projects.1967634730.url":         *(*processWithExpand.Projects)[0].Url,
			},
		},
		{
			name: "success, minimal attributes returned",
			input: map[string]any{
				"id":     id,
				"expand": "projects",
			},
			processToReturn: *minimalProcessWithExpand,
			expectedReturn: map[string]string{
				"expand": "projects",

				"id":                              minimalProcessWithExpand.TypeId.String(),
				"projects.#":                      "1",
				"projects.1967634730.id":          (*minimalProcessWithExpand.Projects)[0].Id.String(),
				"projects.1967634730.name":        "",
				"projects.1967634730.description": "",
				"projects.1967634730.url":         "",
			},
		},
		{
			name: "error from API",
			input: map[string]any{
				"id": id,
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
				GetProcessByItsId(ctx, gomock.Any()).
				DoAndReturn(
					func(ctx context.Context, args workitemtrackingprocess.GetProcessByItsIdArgs) (*workitemtrackingprocess.ProcessInfo, error) {
						if expand, expandFound := tc.input["expand"]; expandFound {
							assert.Equal(t, expand.(string), string(*args.Expand))
						} else {
							assert.Equal(t, workitemtrackingprocess.GetProcessExpandLevelValues.None, *args.Expand)
						}

						if tc.returnError != nil {
							return nil, tc.returnError
						}

						return converter.ToPtr(tc.processToReturn), nil
					},
				).
				Times(1)

			d := getDataProcessResourceData(t, tc.input)

			err := readDataProcess(ctx, d, clients)

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

func createProcessInfo(id string, includeProject bool) *workitemtrackingprocess.ProcessInfo {
	process := workitemtrackingprocess.ProcessInfo{
		TypeId:              converter.UUID(id),
		Name:                converter.String("FirstProc"),
		Description:         converter.String("My first process"),
		CustomizationType:   &workitemtrackingprocess.CustomizationTypeValues.Inherited,
		IsDefault:           converter.Bool(false),
		IsEnabled:           converter.Bool(true),
		ParentProcessTypeId: converter.ToPtr(uuid.New()),
	}
	if includeProject {
		process.Projects = &[]workitemtrackingprocess.ProjectReference{
			{
				Id:          converter.UUID("382fe225-6483-4655-846f-4ac5f7654453"),
				Name:        converter.String("Project1"),
				Description: converter.String("My first project"),
				Url:         converter.String("vstfs:///Classification/TeamProject/6da06557-5456-48c8-b6dc-f111e39a023e"),
			},
		}
	}

	return &process
}

func createMinimalProcessInfo(id string, includeProject bool) *workitemtrackingprocess.ProcessInfo {
	process := workitemtrackingprocess.ProcessInfo{
		TypeId: converter.UUID(id),
	}
	if includeProject {
		process.Projects = &[]workitemtrackingprocess.ProjectReference{
			{
				Id: converter.UUID("382fe225-6483-4655-846f-4ac5f7654453"),
			},
		}
	}

	return &process
}

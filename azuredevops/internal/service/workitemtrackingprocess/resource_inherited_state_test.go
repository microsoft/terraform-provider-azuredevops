//go:build (all || resource_workitemtrackingprocess || resource_workitemtrackingprocess_inherited_state) && !exclude_resource_workitemtrackingprocess
// +build all resource_workitemtrackingprocess resource_workitemtrackingprocess_inherited_state
// +build !exclude_resource_workitemtrackingprocess

package workitemtrackingprocess

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func getInheritedStateResourceData(t *testing.T, input map[string]any) *schema.ResourceData {
	return schema.TestResourceDataRaw(t, ResourceInheritedState().Schema, input)
}

func TestInheritedState_Create(t *testing.T) {
	processId := uuid.New()
	witRefName := "MyProcess.MyWorkItemType"
	stateId := uuid.New()
	stateName := "New"

	tests := []struct {
		name          string
		returnStates  *[]workitemtrackingprocess.WorkItemStateResultModel
		returnError   error
		expectedError string
	}{
		{
			name:          "state not found",
			returnStates:  &[]workitemtrackingprocess.WorkItemStateResultModel{},
			expectedError: "not found",
		},
		{
			name: "custom state",
			returnStates: &[]workitemtrackingprocess.WorkItemStateResultModel{
				{Name: &stateName, Id: &stateId, CustomizationType: &workitemtrackingprocess.CustomizationTypeValues.Custom},
			},
			expectedError: "is a custom state",
		},
		{
			name: "customization type is nil",
			returnStates: &[]workitemtrackingprocess.WorkItemStateResultModel{
				{Name: &stateName, Id: &stateId, CustomizationType: nil},
			},
			expectedError: "has no customization type",
		},
		{
			name: "state ID is nil",
			returnStates: &[]workitemtrackingprocess.WorkItemStateResultModel{
				{Name: &stateName, Id: nil, CustomizationType: &workitemtrackingprocess.CustomizationTypeValues.System},
			},
			expectedError: "state ID is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
			clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

			mockClient.EXPECT().GetStateDefinitions(gomock.Any(), gomock.Any()).Return(tt.returnStates, tt.returnError)
			d := getInheritedStateResourceData(t, map[string]any{
				"process_id":                    processId.String(),
				"work_item_type_reference_name": witRefName,
				"name":                          stateName,
			})

			diags := createResourceInheritedState(context.Background(), d, clients)

			assert.NotEmpty(t, diags)
			assert.Contains(t, diags[0].Summary, tt.expectedError)
		})
	}

}

func TestInheritedState_Import(t *testing.T) {
	processId := uuid.New()
	witRefName := "MyProcess.MyWorkItemType"
	stateId := uuid.New()
	stateName := "New"

	tests := []struct {
		name          string
		returnStates  *[]workitemtrackingprocess.WorkItemStateResultModel
		returnError   error
		expectedError string
	}{
		{
			name:          "state not found",
			returnStates:  &[]workitemtrackingprocess.WorkItemStateResultModel{},
			expectedError: "not found",
		},
		{
			name: "custom state",
			returnStates: &[]workitemtrackingprocess.WorkItemStateResultModel{
				{Name: &stateName, Id: &stateId, CustomizationType: &workitemtrackingprocess.CustomizationTypeValues.Custom},
			},
			expectedError: "is a custom state",
		},
		{
			name: "customization type is nil",
			returnStates: &[]workitemtrackingprocess.WorkItemStateResultModel{
				{Name: &stateName, Id: &stateId, CustomizationType: nil},
			},
			expectedError: "has no customization type",
		},
		{
			name: "state ID is nil",
			returnStates: &[]workitemtrackingprocess.WorkItemStateResultModel{
				{Name: &stateName, Id: nil, CustomizationType: &workitemtrackingprocess.CustomizationTypeValues.System},
			},
			expectedError: "state ID is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
			clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

			mockClient.EXPECT().GetStateDefinitions(gomock.Any(), gomock.Any()).Return(tt.returnStates, tt.returnError)
			d := getInheritedStateResourceData(t, map[string]any{})
			d.SetId(fmt.Sprintf("%s/%s/%s", processId.String(), witRefName, stateName))

			_, err := importResourceInheritedState(context.Background(), d, clients)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

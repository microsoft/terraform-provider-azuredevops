//go:build (all || resource_workitemtrackingprocess || resource_workitemtrackingprocess_inherited_control) && !exclude_resource_workitemtrackingprocess
// +build all resource_workitemtrackingprocess resource_workitemtrackingprocess_inherited_control
// +build !exclude_resource_workitemtrackingprocess

package workitemtrackingprocess

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func getInheritedControlResourceData(t *testing.T, input map[string]any) *schema.ResourceData {
	r := ResourceInheritedControl()

	attributes := make(map[string]string)
	rawConfigValues := make(map[string]cty.Value)
	for k, v := range input {
		if s, ok := v.(string); ok {
			attributes[k] = s
			rawConfigValues[k] = cty.StringVal(s)
		}
	}

	state := &terraform.InstanceState{
		Attributes: attributes,
		RawConfig:  cty.ObjectVal(rawConfigValues),
	}

	return r.Data(state)
}

func createProcessWorkItemTypeWithControl(witRefName, groupId string, control workitemtrackingprocess.Control) *workitemtrackingprocess.ProcessWorkItemType {
	return &workitemtrackingprocess.ProcessWorkItemType{
		ReferenceName: &witRefName,
		Layout: &workitemtrackingprocess.FormLayout{
			Pages: &[]workitemtrackingprocess.Page{
				{
					Id: converter.String("page-1"),
					Sections: &[]workitemtrackingprocess.Section{
						{
							Id: converter.String("section-1"),
							Groups: &[]workitemtrackingprocess.Group{
								{
									Id: &groupId,
									Controls: &[]workitemtrackingprocess.Control{
										control,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestInheritedControl_Update_NilAttributesSentWhenNotConfigured(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	processId := uuid.New()
	witRefName := "MyProcess.MyWorkItemType"
	groupId := "group-1"
	controlId := "System.Description"
	label := "Description"
	visible := true
	inherited := true

	returnControl := workitemtrackingprocess.Control{
		Id:        &controlId,
		Label:     &label,
		Visible:   &visible,
		Inherited: &inherited,
	}

	// The key assertion: when visible and label are not configured, they should be nil
	mockClient.EXPECT().UpdateControl(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.UpdateControlArgs) (*workitemtrackingprocess.Control, error) {
			assert.Equal(t, processId, *args.ProcessId)
			assert.Equal(t, witRefName, *args.WitRefName)
			assert.Equal(t, groupId, *args.GroupId)
			assert.Equal(t, controlId, *args.ControlId)
			// These are the critical assertions: Visible and Label should be nil when not configured
			assert.Nil(t, args.Control.Visible, "Visible should be nil when not configured")
			assert.Nil(t, args.Control.Label, "Label should be nil when not configured")
			return &returnControl, nil
		},
	).Times(1)

	returnWorkItemType := createProcessWorkItemTypeWithControl(witRefName, groupId, returnControl)

	mockClient.EXPECT().GetProcessWorkItemType(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.GetProcessWorkItemTypeArgs) (*workitemtrackingprocess.ProcessWorkItemType, error) {
			return returnWorkItemType, nil
		},
	).Times(1)

	d := getInheritedControlResourceData(t, map[string]any{
		"process_id":                    processId.String(),
		"work_item_type_reference_name": witRefName,
		"group_id":                      groupId,
		"control_id":                    controlId,
	})
	d.SetId(controlId)

	diags := updateResourceInheritedControl(context.Background(), d, clients)
	assert.Empty(t, diags)
}

func TestInheritedControl_Create_Validation(t *testing.T) {
	processId := uuid.New()
	witRefName := "MyProcess.MyWorkItemType"
	existingGroupId := "group-1"
	existingControlId := "System.Title"
	inherited := true
	notInherited := false

	tests := []struct {
		name               string
		groupId            string
		controlId          string
		returnWorkItemType *workitemtrackingprocess.ProcessWorkItemType
		returnError        error
		expectedError      string
	}{
		{
			name:          "API error",
			groupId:       existingGroupId,
			controlId:     existingControlId,
			returnError:   fmt.Errorf("API error"),
			expectedError: "getting work item type",
		},
		{
			name:          "nil work item type",
			groupId:       existingGroupId,
			controlId:     existingControlId,
			expectedError: "work item type or layout is nil",
		},
		{
			name:      "group not found",
			groupId:   "non-existent-group",
			controlId: existingControlId,
			returnWorkItemType: createProcessWorkItemTypeWithControl(witRefName, existingGroupId, workitemtrackingprocess.Control{
				Id:        &existingControlId,
				Inherited: &inherited,
			}),
			expectedError: "group non-existent-group not found in layout",
		},
		{
			name:      "control not found",
			groupId:   existingGroupId,
			controlId: "System.NonExistent",
			returnWorkItemType: createProcessWorkItemTypeWithControl(witRefName, existingGroupId, workitemtrackingprocess.Control{
				Id:        &existingControlId,
				Inherited: &inherited,
			}),
			expectedError: "control System.NonExistent not found in group group-1",
		},
		{
			name:      "control not inherited",
			groupId:   existingGroupId,
			controlId: existingControlId,
			returnWorkItemType: createProcessWorkItemTypeWithControl(witRefName, existingGroupId, workitemtrackingprocess.Control{
				Id:        &existingControlId,
				Inherited: &notInherited,
			}),
			expectedError: "control System.Title is not inherited",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
			clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

			mockClient.EXPECT().GetProcessWorkItemType(clients.Ctx, gomock.Any()).Return(tt.returnWorkItemType, tt.returnError).Times(1)

			d := getInheritedControlResourceData(t, map[string]any{
				"process_id":                    processId.String(),
				"work_item_type_reference_name": witRefName,
				"group_id":                      tt.groupId,
				"control_id":                    tt.controlId,
			})

			diags := createResourceInheritedControl(context.Background(), d, clients)
			assert.NotEmpty(t, diags)
			assert.Contains(t, diags[0].Summary, tt.expectedError)
		})
	}
}

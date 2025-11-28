//go:build (all || resource_workitemtrackingprocess || resource_workitemtrackingprocess_process) && !exclude_resource_workitemtrackingprocess
// +build all resource_workitemtrackingprocess resource_workitemtrackingprocess_process
// +build !exclude_resource_workitemtrackingprocess

package workitemtrackingprocess

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func getWorkItemTypeResourceData(t *testing.T, input map[string]interface{}) *schema.ResourceData {
	r := ResourceWorkItemType()
	return schema.TestResourceDataRaw(t, r.Schema, input)
}

func TestWorkItemType_Create_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	processId := uuid.New()
	name := "MyWorkItemType"
	description := "My Process Description"
	icon := "icon_crown"
	color := "#009ccc"
	isDisabled := false
	inheritsFrom := "MyParent"
	referenceName := "MyNewAgileProcess.MyWorkItemType"
	url := "https://dev.azure.com/foo/_apis/work/processes/4bab314e-358e-4bf3-9508-806ba6ac0c30/workItemTypes/MyNewAgileProcess.MyWorkItemType"

	colorWithoutHash := strings.ReplaceAll(color, "#", "")
	returnWorkItemType := &workitemtrackingprocess.ProcessWorkItemType{
		Icon:          &icon,
		Color:         &colorWithoutHash,
		Inherits:      &inheritsFrom,
		IsDisabled:    &isDisabled,
		Customization: &workitemtrackingprocess.CustomizationTypeValues.Custom,
		Description:   &description,
		Name:          &name,
		ReferenceName: &referenceName,
		Url:           &url,
	}

	mockClient.EXPECT().CreateProcessWorkItemType(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.CreateProcessWorkItemTypeArgs) (*workitemtrackingprocess.ProcessWorkItemType, error) {
			assert.Equal(t, name, *args.WorkItemType.Name)
			assert.Equal(t, description, *args.WorkItemType.Description)
			assert.Equal(t, colorWithoutHash, *args.WorkItemType.Color)
			assert.Equal(t, icon, *args.WorkItemType.Icon)
			assert.Equal(t, inheritsFrom, *args.WorkItemType.InheritsFrom)
			assert.Equal(t, isDisabled, *args.WorkItemType.IsDisabled)

			return returnWorkItemType, nil
		},
	).Times(1)

	d := getWorkItemTypeResourceData(t, map[string]any{
		"process_id":    processId.String(),
		"name":          name,
		"color":         color,
		"icon":          icon,
		"inherits_from": inheritsFrom,
		"is_disabled":   isDisabled,
		"description":   description,
	})

	diags := createResourceWorkItemType(context.Background(), d, clients)
	assert.Empty(t, diags)

	expectedWorkItem := map[string]string{
		"process_id":     processId.String(),
		"name":           name,
		"description":    description,
		"icon":           icon,
		"color":          color,
		"inherits_from":  inheritsFrom,
		"is_disabled":    strconv.FormatBool(isDisabled),
		"id":             referenceName,
		"reference_name": referenceName,
		"url":            url,
	}
	diffOptions := []cmp.Option{
		cmpopts.EquateEmpty(),
	}
	if diff := cmp.Diff(expectedWorkItem, d.State().Attributes, diffOptions...); diff != "" {
		t.Errorf("Resource data attributes mismatch (-want +got):\n%s", diff)
	}
}

func TestWorkItemType_Delete_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	processId := uuid.New()
	referenceName := "MyNewAgileProcess.MyWorkItemType"

	mockClient.EXPECT().DeleteProcessWorkItemType(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.DeleteProcessWorkItemTypeArgs) error {
			assert.Equal(t, processId, *args.ProcessId)
			assert.Equal(t, referenceName, *args.WitRefName)
			return nil
		},
	).Times(1)

	d := getWorkItemTypeResourceData(t, map[string]any{
		"name":       "MyWorkItemType",
		"process_id": processId.String(),
	})
	d.SetId(referenceName)

	diags := deleteResourceWorkItemType(context.Background(), d, clients)

	assert.Empty(t, diags)
}

func TestWorkItemType_Read_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	processId := uuid.New()
	name := "MyWorkItemType"
	description := "My Process Description"
	icon := "icon_crown"
	color := "#009ccc"
	isDisabled := false
	inheritsFrom := "MyParent"
	referenceName := "MyNewAgileProcess.MyWorkItemType"
	url := "https://dev.azure.com/foo/_apis/work/processes/4bab314e-358e-4bf3-9508-806ba6ac0c30/workItemTypes/MyNewAgileProcess.MyWorkItemType"

	colorWithoutHash := strings.ReplaceAll(color, "#", "")
	returnWorkItemType := &workitemtrackingprocess.ProcessWorkItemType{
		Icon:          &icon,
		Color:         &colorWithoutHash,
		Inherits:      &inheritsFrom,
		IsDisabled:    &isDisabled,
		Customization: &workitemtrackingprocess.CustomizationTypeValues.Custom,
		Description:   &description,
		Name:          &name,
		ReferenceName: &referenceName,
		Url:           &url,
	}

	mockClient.EXPECT().GetProcessWorkItemType(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.GetProcessWorkItemTypeArgs) (*workitemtrackingprocess.ProcessWorkItemType, error) {
			assert.Equal(t, workitemtrackingprocess.GetWorkItemTypeExpandValues.None, *args.Expand)
			assert.Equal(t, processId, *args.ProcessId)
			assert.Equal(t, referenceName, *args.WitRefName)

			return returnWorkItemType, nil
		},
	).Times(1)

	d := getWorkItemTypeResourceData(t, map[string]any{
		"process_id": processId.String(),
		"name":       name,
	})
	d.SetId(referenceName)

	diags := readResourceWorkItemType(context.Background(), d, clients)
	assert.Empty(t, diags)

	expectedWorkItem := map[string]string{
		"process_id":     processId.String(),
		"name":           name,
		"description":    description,
		"icon":           icon,
		"color":          color,
		"inherits_from":  inheritsFrom,
		"is_disabled":    strconv.FormatBool(isDisabled),
		"id":             referenceName,
		"reference_name": referenceName,
		"url":            url,
	}
	diffOptions := []cmp.Option{
		cmpopts.EquateEmpty(),
	}
	if diff := cmp.Diff(expectedWorkItem, d.State().Attributes, diffOptions...); diff != "" {
		t.Errorf("Resource data attributes mismatch (-want +got):\n%s", diff)
	}
}

func TestWorkItemType_Update_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	processId := uuid.New()
	name := "MyWorkItemType"
	description := "My Process Description"
	icon := "icon_crown"
	color := "#009ccc"
	isDisabled := false
	inheritsFrom := "MyParent"
	referenceName := "MyNewAgileProcess.MyWorkItemType"
	url := "https://dev.azure.com/foo/_apis/work/processes/4bab314e-358e-4bf3-9508-806ba6ac0c30/workItemTypes/MyNewAgileProcess.MyWorkItemType"

	colorWithoutHash := strings.ReplaceAll(color, "#", "")
	returnWorkItemType := &workitemtrackingprocess.ProcessWorkItemType{
		Icon:          &icon,
		Color:         &colorWithoutHash,
		Inherits:      &inheritsFrom,
		IsDisabled:    &isDisabled,
		Customization: &workitemtrackingprocess.CustomizationTypeValues.Custom,
		Description:   &description,
		Name:          &name,
		ReferenceName: &referenceName,
		Url:           &url,
	}

	mockClient.EXPECT().UpdateProcessWorkItemType(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.UpdateProcessWorkItemTypeArgs) (*workitemtrackingprocess.ProcessWorkItemType, error) {
			assert.Equal(t, processId, *args.ProcessId)
			assert.Equal(t, referenceName, *args.WitRefName)
			assert.Equal(t, description, *args.WorkItemTypeUpdate.Description)
			assert.Equal(t, colorWithoutHash, *args.WorkItemTypeUpdate.Color)
			assert.Equal(t, icon, *args.WorkItemTypeUpdate.Icon)
			assert.Equal(t, isDisabled, *args.WorkItemTypeUpdate.IsDisabled)

			return returnWorkItemType, nil
		},
	).Times(1)

	d := getWorkItemTypeResourceData(t, map[string]any{
		"process_id":    processId.String(),
		"name":          name,
		"color":         color,
		"icon":          icon,
		"inherits_from": inheritsFrom,
		"is_disabled":   isDisabled,
		"description":   description,
	})
	d.SetId(referenceName)

	diags := updateResourceWorkItemType(context.Background(), d, clients)
	assert.Empty(t, diags)

	expectedWorkItem := map[string]string{
		"process_id":     processId.String(),
		"name":           name,
		"description":    description,
		"icon":           icon,
		"color":          color,
		"inherits_from":  inheritsFrom,
		"is_disabled":    strconv.FormatBool(isDisabled),
		"id":             referenceName,
		"reference_name": referenceName,
		"url":            url,
	}
	diffOptions := []cmp.Option{
		cmpopts.EquateEmpty(),
	}
	if diff := cmp.Diff(expectedWorkItem, d.State().Attributes, diffOptions...); diff != "" {
		t.Errorf("Resource data attributes mismatch (-want +got):\n%s", diff)
	}
}

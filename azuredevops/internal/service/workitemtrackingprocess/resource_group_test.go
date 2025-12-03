//go:build (all || resource_workitemtrackingprocess || resource_workitemtrackingprocess_group) && !exclude_resource_workitemtrackingprocess
// +build all resource_workitemtrackingprocess resource_workitemtrackingprocess_group
// +build !exclude_resource_workitemtrackingprocess

package workitemtrackingprocess

import (
	"context"
	"strconv"
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

func getGroupResourceData(t *testing.T, input map[string]any) *schema.ResourceData {
	r := ResourceGroup()
	return schema.TestResourceDataRaw(t, r.Schema, input)
}

// Simulates an already created group (persisted in state)
func getPersistedGroupResourceData(t *testing.T, id string, input map[string]any) *schema.ResourceData {
	r := ResourceGroup()
	d := schema.TestResourceDataRaw(t, r.Schema, input)
	d.SetId(id)
	return r.Data(d.State())
}

func TestGroup_Create_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	processId := uuid.New()
	witRefName := "MyNewAgileProcess.MyWorkItemType"
	pageId := "page-1"
	sectionId := "section-1"
	groupId := "group-1"
	label := "My Group"
	order := 1
	visible := true

	returnGroup := &workitemtrackingprocess.Group{
		Id:      &groupId,
		Label:   &label,
		Order:   &order,
		Visible: &visible,
	}

	mockClient.EXPECT().AddGroup(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.AddGroupArgs) (*workitemtrackingprocess.Group, error) {
			assert.Equal(t, processId, *args.ProcessId)
			assert.Equal(t, witRefName, *args.WitRefName)
			assert.Equal(t, pageId, *args.PageId)
			assert.Equal(t, sectionId, *args.SectionId)
			assert.Equal(t, label, *args.Group.Label)
			assert.Equal(t, order, *args.Group.Order)
			assert.Equal(t, visible, *args.Group.Visible)

			return returnGroup, nil
		},
	).Times(1)

	d := getGroupResourceData(t, map[string]any{
		"process_id":                    processId.String(),
		"work_item_type_reference_name": witRefName,
		"page_id":                       pageId,
		"section_id":                    sectionId,
		"label":                         label,
		"order":                         order,
		"visible":                       visible,
	})

	diags := createResourceGroup(context.Background(), d, clients)
	assert.Empty(t, diags)

	expectedGroup := map[string]string{
		"process_id":                    processId.String(),
		"work_item_type_reference_name": witRefName,
		"page_id":                       pageId,
		"section_id":                    sectionId,
		"label":                         label,
		"order":                         strconv.Itoa(order),
		"visible":                       strconv.FormatBool(visible),
		"id":                            groupId,
	}
	diffOptions := []cmp.Option{
		cmpopts.EquateEmpty(),
	}
	if diff := cmp.Diff(expectedGroup, d.State().Attributes, diffOptions...); diff != "" {
		t.Errorf("Resource data attributes mismatch (-want +got):\n%s", diff)
	}
}

func TestGroup_Delete_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	processId := uuid.New()
	witRefName := "MyNewAgileProcess.MyWorkItemType"
	pageId := "page-1"
	sectionId := "section-1"
	groupId := "group-1"

	mockClient.EXPECT().RemoveGroup(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.RemoveGroupArgs) error {
			assert.Equal(t, processId, *args.ProcessId)
			assert.Equal(t, witRefName, *args.WitRefName)
			assert.Equal(t, pageId, *args.PageId)
			assert.Equal(t, sectionId, *args.SectionId)
			assert.Equal(t, groupId, *args.GroupId)
			return nil
		},
	).Times(1)

	d := getGroupResourceData(t, map[string]any{
		"process_id":                    processId.String(),
		"work_item_type_reference_name": witRefName,
		"page_id":                       pageId,
		"section_id":                    sectionId,
		"label":                         "My Group",
	})
	d.SetId(groupId)

	diags := deleteResourceGroup(context.Background(), d, clients)

	assert.Empty(t, diags)
}

func TestGroup_Read_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	processId := uuid.New()
	witRefName := "MyNewAgileProcess.MyWorkItemType"
	pageId := "page-1"
	sectionId := "section-1"
	groupId := "group-1"
	label := "My Group"
	order := 1
	visible := true

	returnWorkItemType := &workitemtrackingprocess.ProcessWorkItemType{
		ReferenceName: &witRefName,
		Layout: &workitemtrackingprocess.FormLayout{
			Pages: &[]workitemtrackingprocess.Page{
				{
					Id: &pageId,
					Sections: &[]workitemtrackingprocess.Section{
						{
							Id: &sectionId,
							Groups: &[]workitemtrackingprocess.Group{
								{
									Id:      &groupId,
									Label:   &label,
									Order:   &order,
									Visible: &visible,
								},
							},
						},
					},
				},
			},
		},
	}

	mockClient.EXPECT().GetProcessWorkItemType(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.GetProcessWorkItemTypeArgs) (*workitemtrackingprocess.ProcessWorkItemType, error) {
			assert.Equal(t, workitemtrackingprocess.GetWorkItemTypeExpandValues.Layout, *args.Expand)
			assert.Equal(t, processId, *args.ProcessId)
			assert.Equal(t, witRefName, *args.WitRefName)

			return returnWorkItemType, nil
		},
	).Times(1)

	d := getGroupResourceData(t, map[string]any{
		"process_id":                    processId.String(),
		"work_item_type_reference_name": witRefName,
		"page_id":                       pageId,
		"section_id":                    sectionId,
		"label":                         label,
	})
	d.SetId(groupId)

	diags := readResourceGroup(context.Background(), d, clients)
	assert.Empty(t, diags)

	expectedGroup := map[string]string{
		"process_id":                    processId.String(),
		"work_item_type_reference_name": witRefName,
		"page_id":                       pageId,
		"section_id":                    sectionId,
		"label":                         label,
		"order":                         strconv.Itoa(order),
		"visible":                       strconv.FormatBool(visible),
		"id":                            groupId,
	}
	diffOptions := []cmp.Option{
		cmpopts.EquateEmpty(),
	}
	if diff := cmp.Diff(expectedGroup, d.State().Attributes, diffOptions...); diff != "" {
		t.Errorf("Resource data attributes mismatch (-want +got):\n%s", diff)
	}
}

func TestGroup_Update_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	processId := uuid.New()
	witRefName := "MyNewAgileProcess.MyWorkItemType"
	pageId := "page-1"
	sectionId := "section-1"
	groupId := "group-1"
	initialLabel := "My Group"
	initialOrder := 1
	initialVisible := true
	updatedLabel := "Updated Group Label"
	updatedOrder := 2
	updatedVisible := false

	returnGroup := &workitemtrackingprocess.Group{
		Id:      &groupId,
		Label:   &updatedLabel,
		Order:   &updatedOrder,
		Visible: &updatedVisible,
	}

	mockClient.EXPECT().UpdateGroup(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.UpdateGroupArgs) (*workitemtrackingprocess.Group, error) {
			assert.Equal(t, processId, *args.ProcessId)
			assert.Equal(t, witRefName, *args.WitRefName)
			assert.Equal(t, pageId, *args.PageId)
			assert.Equal(t, sectionId, *args.SectionId)
			assert.Equal(t, groupId, *args.GroupId)
			assert.Equal(t, updatedLabel, *args.Group.Label)
			assert.Equal(t, updatedOrder, *args.Group.Order)
			assert.Equal(t, updatedVisible, *args.Group.Visible)

			return returnGroup, nil
		},
	).Times(1)

	d := getPersistedGroupResourceData(t, groupId, map[string]any{
		"process_id":                    processId.String(),
		"work_item_type_reference_name": witRefName,
		"page_id":                       pageId,
		"section_id":                    sectionId,
		"label":                         initialLabel,
		"order":                         initialOrder,
		"visible":                       initialVisible,
	})

	// Set the new values
	d.Set("label", updatedLabel)
	d.Set("order", updatedOrder)
	d.Set("visible", updatedVisible)

	diags := updateResourceGroup(context.Background(), d, clients)
	assert.Empty(t, diags)

	expectedGroup := map[string]string{
		"process_id":                    processId.String(),
		"work_item_type_reference_name": witRefName,
		"page_id":                       pageId,
		"section_id":                    sectionId,
		"label":                         updatedLabel,
		"order":                         strconv.Itoa(updatedOrder),
		"visible":                       strconv.FormatBool(updatedVisible),
		"id":                            groupId,
	}
	diffOptions := []cmp.Option{
		cmpopts.EquateEmpty(),
	}
	if diff := cmp.Diff(expectedGroup, d.State().Attributes, diffOptions...); diff != "" {
		t.Errorf("Resource data attributes mismatch (-want +got):\n%s", diff)
	}
}

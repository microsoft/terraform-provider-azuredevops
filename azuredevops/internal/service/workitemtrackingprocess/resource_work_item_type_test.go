//go:build (all || resource_workitemtrackingprocess || resource_workitemtrackingprocess_workitemtype) && !exclude_resource_workitemtrackingprocess
// +build all resource_workitemtrackingprocess resource_workitemtrackingprocess_workitemtype
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

	pageId := "page-1"
	sectionId := "section-1"
	groupId := "group-1"

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

	withLayout := func(wit *workitemtrackingprocess.ProcessWorkItemType) *workitemtrackingprocess.ProcessWorkItemType {
		wit.Layout = &workitemtrackingprocess.FormLayout{
			Pages: &[]workitemtrackingprocess.Page{
				{
					Id:       &pageId,
					PageType: &workitemtrackingprocess.PageTypeValues.Custom,
					Sections: &[]workitemtrackingprocess.Section{
						{
							Id: &sectionId,
							Groups: &[]workitemtrackingprocess.Group{
								{
									Id: &groupId,
								},
							},
						},
					},
				},
			},
		}
		return wit
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

	// Create calls read to get the full state including layout
	mockClient.EXPECT().GetProcessWorkItemType(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.GetProcessWorkItemTypeArgs) (*workitemtrackingprocess.ProcessWorkItemType, error) {
			assert.Equal(t, workitemtrackingprocess.GetWorkItemTypeExpandValues.Layout, *args.Expand)
			assert.Equal(t, processId, *args.ProcessId)
			assert.Equal(t, referenceName, *args.WitRefName)

			return withLayout(returnWorkItemType), nil
		},
	).Times(1)

	d := getWorkItemTypeResourceData(t, map[string]any{
		"process_id":                      processId.String(),
		"name":                            name,
		"color":                           color,
		"icon":                            icon,
		"parent_work_item_reference_name": inheritsFrom,
		"is_enabled":                      !isDisabled,
		"description":                     description,
	})

	diags := createResourceWorkItemType(context.Background(), d, clients)
	assert.Empty(t, diags)

	expectedWorkItem := map[string]string{
		"process_id":                      processId.String(),
		"name":                            name,
		"description":                     description,
		"icon":                            icon,
		"color":                           color,
		"parent_work_item_reference_name": inheritsFrom,
		"is_enabled":                      strconv.FormatBool(!isDisabled),
		"id":                              referenceName,
		"reference_name":                  referenceName,
		"url":                             url,

		"pages.#":                                "1",
		"pages.0.id":                             pageId,
		"pages.0.page_type":                      "custom",
		"pages.0.sections.#":                     "1",
		"pages.0.sections.0.id":                  sectionId,
		"pages.0.sections.0.groups.#":            "1",
		"pages.0.sections.0.groups.0.id":         groupId,
		"pages.0.sections.0.groups.0.controls.#": "0",
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

	// Create full pages structure
	controlId1 := "control-1"
	controlId2 := "control-2"
	groupId1 := "group-1"
	groupId2 := "group-2"
	sectionId1 := "section-1"
	sectionId2 := "section-2"
	pageId1 := "page-1"
	pageId2 := "page-2"

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
		Layout: &workitemtrackingprocess.FormLayout{
			Pages: &[]workitemtrackingprocess.Page{
				{
					Id:       &pageId1,
					PageType: &workitemtrackingprocess.PageTypeValues.Custom,
					Sections: &[]workitemtrackingprocess.Section{
						{
							Id: &sectionId1,
							Groups: &[]workitemtrackingprocess.Group{
								{
									Id: &groupId1,
									Controls: &[]workitemtrackingprocess.Control{
										{
											Id: &controlId1,
										},
										{
											Id: &controlId2,
										},
										{},
									},
								},
								{
									Id:       &groupId2,
									Controls: &[]workitemtrackingprocess.Control{},
								},
								{},
							},
						},
						{
							Id:     &sectionId2,
							Groups: &[]workitemtrackingprocess.Group{},
						},
						{},
					},
				},
				{
					Id:       &pageId2,
					PageType: &workitemtrackingprocess.PageTypeValues.History,
					Sections: &[]workitemtrackingprocess.Section{},
				},
				{},
			},
		},
	}

	mockClient.EXPECT().GetProcessWorkItemType(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.GetProcessWorkItemTypeArgs) (*workitemtrackingprocess.ProcessWorkItemType, error) {
			assert.Equal(t, workitemtrackingprocess.GetWorkItemTypeExpandValues.Layout, *args.Expand)
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
		"process_id":                      processId.String(),
		"name":                            name,
		"description":                     description,
		"icon":                            icon,
		"color":                           color,
		"parent_work_item_reference_name": inheritsFrom,
		"is_enabled":                      strconv.FormatBool(!isDisabled),
		"id":                              referenceName,
		"reference_name":                  referenceName,
		"url":                             url,

		"pages.#":                                   "3",
		"pages.0.id":                                pageId1,
		"pages.0.page_type":                         "custom",
		"pages.0.sections.#":                        "3",
		"pages.0.sections.0.id":                     sectionId1,
		"pages.0.sections.0.groups.#":               "3",
		"pages.0.sections.0.groups.0.id":            groupId1,
		"pages.0.sections.0.groups.0.controls.#":    "3",
		"pages.0.sections.0.groups.0.controls.0.id": controlId1,
		"pages.0.sections.0.groups.0.controls.1.id": controlId2,
		"pages.0.sections.0.groups.1.id":            groupId2,
		"pages.0.sections.0.groups.1.controls.#":    "0",
		"pages.0.sections.1.id":                     sectionId2,
		"pages.0.sections.1.groups.#":               "0",
		"pages.1.id":                                pageId2,
		"pages.1.page_type":                         "history",
		"pages.1.sections.#":                        "0",
	}
	diffOptions := []cmp.Option{
		cmpopts.EquateEmpty(),
	}
	if diff := cmp.Diff(expectedWorkItem, d.State().Attributes, diffOptions...); diff != "" {
		t.Errorf("Resource data attributes mismatch (-want +got):\n%s", diff)
	}
}

func TestWorkItemType_Read_APIReturnsNoProperties(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	processId := uuid.New()
	referenceName := "MyNewAgileProcess.MyWorkItemType"

	returnWorkItemType := &workitemtrackingprocess.ProcessWorkItemType{
		ReferenceName: &referenceName,
	}

	mockClient.EXPECT().GetProcessWorkItemType(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.GetProcessWorkItemTypeArgs) (*workitemtrackingprocess.ProcessWorkItemType, error) {
			assert.Equal(t, workitemtrackingprocess.GetWorkItemTypeExpandValues.Layout, *args.Expand)
			assert.Equal(t, processId, *args.ProcessId)
			assert.Equal(t, referenceName, *args.WitRefName)

			return returnWorkItemType, nil
		},
	).Times(1)

	d := getWorkItemTypeResourceData(t, map[string]any{
		"process_id":     processId.String(),
		"reference_name": referenceName,
	})
	d.SetId(referenceName)

	diags := readResourceWorkItemType(context.Background(), d, clients)
	assert.Empty(t, diags)

	// When API returns nil for all properties except reference_name, state should reflect that
	expectedState := map[string]string{
		"id":                              referenceName,
		"process_id":                      processId.String(),
		"name":                            "",
		"description":                     "",
		"is_enabled":                      "true",
		"color":                           "",
		"icon":                            "",
		"parent_work_item_reference_name": "",
		"reference_name":                  referenceName,
		"url":                             "",
		"pages.#":                         "0",
	}
	if diff := cmp.Diff(expectedState, d.State().Attributes); diff != "" {
		t.Errorf("expected resource attributes to correspond to the API response: (-expected +got):\n%s", diff)
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

	pageId := "page-1"
	sectionId := "section-1"
	groupId := "group-1"

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

	withLayout := func(wit *workitemtrackingprocess.ProcessWorkItemType) *workitemtrackingprocess.ProcessWorkItemType {
		wit.Layout = &workitemtrackingprocess.FormLayout{
			Pages: &[]workitemtrackingprocess.Page{
				{
					Id:       &pageId,
					PageType: &workitemtrackingprocess.PageTypeValues.Custom,
					Sections: &[]workitemtrackingprocess.Section{
						{
							Id: &sectionId,
							Groups: &[]workitemtrackingprocess.Group{
								{
									Id: &groupId,
								},
							},
						},
					},
				},
			},
		}
		return wit
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

	// Update calls read to get the full state including layout
	mockClient.EXPECT().GetProcessWorkItemType(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.GetProcessWorkItemTypeArgs) (*workitemtrackingprocess.ProcessWorkItemType, error) {
			assert.Equal(t, workitemtrackingprocess.GetWorkItemTypeExpandValues.Layout, *args.Expand)
			assert.Equal(t, processId, *args.ProcessId)
			assert.Equal(t, referenceName, *args.WitRefName)

			return withLayout(returnWorkItemType), nil
		},
	).Times(1)

	d := getWorkItemTypeResourceData(t, map[string]any{
		"process_id":                      processId.String(),
		"name":                            name,
		"color":                           color,
		"icon":                            icon,
		"parent_work_item_reference_name": inheritsFrom,
		"is_enabled":                      !isDisabled,
		"description":                     description,
	})
	d.SetId(referenceName)

	diags := updateResourceWorkItemType(context.Background(), d, clients)
	assert.Empty(t, diags)

	expectedWorkItem := map[string]string{
		"process_id":                      processId.String(),
		"name":                            name,
		"description":                     description,
		"icon":                            icon,
		"color":                           color,
		"parent_work_item_reference_name": inheritsFrom,
		"is_enabled":                      strconv.FormatBool(!isDisabled),
		"id":                              referenceName,
		"reference_name":                  referenceName,
		"url":                             url,

		"pages.#":                                "1",
		"pages.0.id":                             pageId,
		"pages.0.page_type":                      "custom",
		"pages.0.sections.#":                     "1",
		"pages.0.sections.0.id":                  sectionId,
		"pages.0.sections.0.groups.#":            "1",
		"pages.0.sections.0.groups.0.id":         groupId,
		"pages.0.sections.0.groups.0.controls.#": "0",
	}
	diffOptions := []cmp.Option{
		cmpopts.EquateEmpty(),
	}
	if diff := cmp.Diff(expectedWorkItem, d.State().Attributes, diffOptions...); diff != "" {
		t.Errorf("Resource data attributes mismatch (-want +got):\n%s", diff)
	}
}

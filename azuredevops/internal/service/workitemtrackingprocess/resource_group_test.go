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
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func getGroupResourceData(t *testing.T, input map[string]any) *schema.ResourceData {
	r := ResourceGroup()
	return schema.TestResourceDataRaw(t, r.Schema, input)
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
	order := 0
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

func TestGroup_FindGroupById(t *testing.T) {
	groupId := "target-group"
	pageId1 := "page-1"
	pageId2 := "page-2"
	sectionId1 := "section-1"
	sectionId2 := "section-2"

	tests := []struct {
		name     string
		layout   *workitemtrackingprocess.FormLayout
		groupId  string
		expected bool
	}{
		{
			name: "found in first page first section",
			layout: &workitemtrackingprocess.FormLayout{
				Pages: &[]workitemtrackingprocess.Page{
					{
						Id: &pageId1,
						Sections: &[]workitemtrackingprocess.Section{
							{
								Id: &sectionId1,
								Groups: &[]workitemtrackingprocess.Group{
									{Id: &groupId},
								},
							},
						},
					},
				},
			},
			groupId:  groupId,
			expected: true,
		},
		{
			name: "found in second page",
			layout: &workitemtrackingprocess.FormLayout{
				Pages: &[]workitemtrackingprocess.Page{
					{
						Id: &pageId1,
						Sections: &[]workitemtrackingprocess.Section{
							{
								Id:     &sectionId1,
								Groups: &[]workitemtrackingprocess.Group{},
							},
						},
					},
					{
						Id: &pageId2,
						Sections: &[]workitemtrackingprocess.Section{
							{
								Id: &sectionId2,
								Groups: &[]workitemtrackingprocess.Group{
									{Id: &groupId},
								},
							},
						},
					},
				},
			},
			groupId:  groupId,
			expected: true,
		},
		{
			name: "found among multiple groups",
			layout: &workitemtrackingprocess.FormLayout{
				Pages: &[]workitemtrackingprocess.Page{
					{
						Id: &pageId1,
						Sections: &[]workitemtrackingprocess.Section{
							{
								Id: &sectionId1,
								Groups: &[]workitemtrackingprocess.Group{
									{Id: converter.String("other-1")},
									{Id: converter.String("other-2")},
									{Id: &groupId},
									{Id: converter.String("other-3")},
								},
							},
						},
					},
				},
			},
			groupId:  groupId,
			expected: true,
		},
		{
			name: "not found",
			layout: &workitemtrackingprocess.FormLayout{
				Pages: &[]workitemtrackingprocess.Page{
					{
						Id: &pageId1,
						Sections: &[]workitemtrackingprocess.Section{
							{
								Id: &sectionId1,
								Groups: &[]workitemtrackingprocess.Group{
									{Id: converter.String("other-group")},
								},
							},
						},
					},
				},
			},
			groupId:  "nonexistent",
			expected: false,
		},
		{
			name:     "nil layout",
			layout:   nil,
			groupId:  groupId,
			expected: false,
		},
		{
			name: "nil pages",
			layout: &workitemtrackingprocess.FormLayout{
				Pages: nil,
			},
			groupId:  groupId,
			expected: false,
		},
		{
			name: "empty pages",
			layout: &workitemtrackingprocess.FormLayout{
				Pages: &[]workitemtrackingprocess.Page{},
			},
			groupId:  groupId,
			expected: false,
		},
		{
			name: "nil sections",
			layout: &workitemtrackingprocess.FormLayout{
				Pages: &[]workitemtrackingprocess.Page{
					{
						Id:       &pageId1,
						Sections: nil,
					},
				},
			},
			groupId:  groupId,
			expected: false,
		},
		{
			name: "nil groups",
			layout: &workitemtrackingprocess.FormLayout{
				Pages: &[]workitemtrackingprocess.Page{
					{
						Id: &pageId1,
						Sections: &[]workitemtrackingprocess.Section{
							{
								Id:     &sectionId1,
								Groups: nil,
							},
						},
					},
				},
			},
			groupId:  groupId,
			expected: false,
		},
		{
			name: "group with nil id",
			layout: &workitemtrackingprocess.FormLayout{
				Pages: &[]workitemtrackingprocess.Page{
					{
						Id: &pageId1,
						Sections: &[]workitemtrackingprocess.Section{
							{
								Id: &sectionId1,
								Groups: &[]workitemtrackingprocess.Group{
									{Id: nil},
								},
							},
						},
					},
				},
			},
			groupId:  groupId,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findGroupById(tt.layout, tt.groupId)
			if tt.expected {
				assert.NotNil(t, result, "expected to find group")
				assert.Equal(t, tt.groupId, *result.Id)
			} else {
				assert.Nil(t, result, "expected not to find group")
			}
		})
	}
}

func TestGroup_Import(t *testing.T) {
	tests := []struct {
		name                          string
		importId                      string
		expectError                   bool
		errorContains                 string
		expectedProcessId             string
		expectedWorkItemTypeReference string
		expectedPageId                string
		expectedSectionId             string
		expectedGroupId               string
	}{
		{
			name:                          "valid import id",
			importId:                      "00000000-0000-0000-0000-000000000001/MyProcess.MyWorkItemType/page-1/section-1/group-1",
			expectError:                   false,
			expectedProcessId:             "00000000-0000-0000-0000-000000000001",
			expectedWorkItemTypeReference: "MyProcess.MyWorkItemType",
			expectedPageId:                "page-1",
			expectedSectionId:             "section-1",
			expectedGroupId:               "group-1",
		},
		{
			name:          "missing parts",
			importId:      "process-id/wit-ref-name/page-id",
			expectError:   true,
			errorContains: "invalid import ID format",
		},
		{
			name:          "too many parts",
			importId:      "process-id/wit-ref-name/page-id/section-id/group-id/extra",
			expectError:   true,
			errorContains: "invalid import ID format",
		},
		{
			name:          "empty string",
			importId:      "",
			expectError:   true,
			errorContains: "invalid import ID format",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			d := getGroupResourceData(t, map[string]any{})
			d.SetId(testCase.importId)

			result, err := importResourceGroup(context.Background(), d, nil)

			if testCase.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.errorContains)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, 1)
				assert.Equal(t, testCase.expectedProcessId, d.Get("process_id"))
				assert.Equal(t, testCase.expectedWorkItemTypeReference, d.Get("work_item_type_reference_name"))
				assert.Equal(t, testCase.expectedPageId, d.Get("page_id"))
				assert.Equal(t, testCase.expectedSectionId, d.Get("section_id"))
				assert.Equal(t, testCase.expectedGroupId, d.Id())
			}
		})
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
	label := "My Group"
	order := 1
	visible := true

	returnGroup := &workitemtrackingprocess.Group{
		Id:      &groupId,
		Label:   &label,
		Order:   &order,
		Visible: &visible,
	}

	// TestResourceDataRaw treats all fields as changed, so HasChange will return true
	// This means MoveGroupToPage will be called even though we're not actually moving
	// MoveGroupToPage also updates the group, so UpdateGroup is not called
	mockClient.EXPECT().MoveGroupToPage(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.MoveGroupToPageArgs) (*workitemtrackingprocess.Group, error) {
			assert.Equal(t, processId, *args.ProcessId)
			assert.Equal(t, witRefName, *args.WitRefName)
			assert.Equal(t, pageId, *args.PageId)
			assert.Equal(t, sectionId, *args.SectionId)
			assert.Equal(t, groupId, *args.GroupId)
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
	d.SetId(groupId)

	diags := updateResourceGroup(context.Background(), d, clients)
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

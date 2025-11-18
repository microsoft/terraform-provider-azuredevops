//go:build (all || resource_workitemtrackingprocess_process) && !resource_workitemtrackingprocess
// +build all resource_workitemtrackingprocess_process
// +build !resource_workitemtrackingprocess

package workitemtrackingprocess

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func getProcessResourceData(t *testing.T, input map[string]interface{}) *schema.ResourceData {
	r := ResourceProcess()
	return schema.TestResourceDataRaw(t, r.Schema, input)
}

func TestProcesses_Create_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	parentID := uuid.New()
	typeID := uuid.New()
	name := "MyProcess"
	description := "My Process Description"
	referenceName := "My.Process.ReferenceName"
	expectedProcessInfo := &workitemtrackingprocess.ProcessInfo{
		TypeId:              &typeID,
		Name:                &name,
		Description:         &description,
		IsDefault:           converter.Bool(false),
		IsEnabled:           converter.Bool(true),
		CustomizationType:   &workitemtrackingprocess.CustomizationTypeValues.Inherited,
		ParentProcessTypeId: &parentID,
		ReferenceName:       &referenceName,
	}

	mockClient.EXPECT().CreateNewProcess(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.CreateNewProcessArgs) (*workitemtrackingprocess.ProcessInfo, error) {
			assert.Equal(t, name, *args.CreateRequest.Name)
			assert.Equal(t, description, *args.CreateRequest.Description)
			assert.Equal(t, parentID, *args.CreateRequest.ParentProcessTypeId)
			assert.Equal(t, referenceName, *args.CreateRequest.ReferenceName)
			return expectedProcessInfo, nil
		},
	).Times(1)

	mockClient.EXPECT().GetProcessByItsId(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.GetProcessByItsIdArgs) (*workitemtrackingprocess.ProcessInfo, error) {
			assert.Equal(t, typeID, *args.ProcessTypeId)
			return expectedProcessInfo, nil
		},
	).Times(1)

	d := getProcessResourceData(t, map[string]any{
		"name":                   name,
		"parent_process_type_id": parentID.String(),
		"description":            description,
		"reference_name":         referenceName,
	})

	diags := createResourceProcess(context.Background(), d, clients)

	assert.Empty(t, diags)
	assert.Equal(t, typeID.String(), d.Id())
	assert.Equal(t, name, d.Get("name"))
	assert.Equal(t, description, d.Get("description"))
	assert.Equal(t, referenceName, d.Get("reference_name"))
	assert.Equal(t, parentID.String(), d.Get("parent_process_type_id"))
	assert.Equal(t, false, d.Get("is_default"))
	assert.Equal(t, true, d.Get("is_enabled"))
	assert.Equal(t, "inherited", d.Get("customization_type"))
	projects := d.Get("projects").(*schema.Set)
	assert.NotNil(t, projects)
	assert.Equal(t, 0, projects.Len())
}

func TestProcesses_CreateWithDefault_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	parentID := uuid.New()
	typeID := uuid.New()
	name := "MyProcess"
	referenceName := "My.Process.ReferenceName"
	expectedProcessInfo := &workitemtrackingprocess.ProcessInfo{
		TypeId:              &typeID,
		Name:                &name,
		Description:         converter.String(""),
		IsDefault:           converter.Bool(false),
		IsEnabled:           converter.Bool(true),
		CustomizationType:   &workitemtrackingprocess.CustomizationTypeValues.Inherited,
		ParentProcessTypeId: &parentID,
		ReferenceName:       &referenceName,
	}

	mockClient.EXPECT().CreateNewProcess(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.CreateNewProcessArgs) (*workitemtrackingprocess.ProcessInfo, error) {
			assert.Equal(t, name, *args.CreateRequest.Name)
			assert.Equal(t, parentID, *args.CreateRequest.ParentProcessTypeId)
			assert.Nil(t, args.CreateRequest.ReferenceName)
			return expectedProcessInfo, nil
		},
	).Times(1)
	mockClient.EXPECT().EditProcess(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.EditProcessArgs) (*workitemtrackingprocess.ProcessInfo, error) {
			assert.Equal(t, typeID, *args.ProcessTypeId)
			assert.Equal(t, name, *args.UpdateRequest.Name)
			assert.True(t, *args.UpdateRequest.IsDefault)
			assert.True(t, *args.UpdateRequest.IsEnabled)
			expectedProcessInfo.IsDefault = args.UpdateRequest.IsDefault
			return expectedProcessInfo, nil
		},
	).Times(1)

	mockClient.EXPECT().GetProcessByItsId(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.GetProcessByItsIdArgs) (*workitemtrackingprocess.ProcessInfo, error) {
			assert.Equal(t, typeID, *args.ProcessTypeId)
			return expectedProcessInfo, nil
		},
	).Times(1)

	d := getProcessResourceData(t, map[string]any{
		"name":                   name,
		"parent_process_type_id": parentID.String(),
		"is_default":             true,
	})

	diags := createResourceProcess(context.Background(), d, clients)

	assert.Empty(t, diags)
	assert.Equal(t, typeID.String(), d.Id())
	assert.Equal(t, name, d.Get("name"))
	assert.Equal(t, referenceName, d.Get("reference_name"))
	assert.Equal(t, parentID.String(), d.Get("parent_process_type_id"))
	assert.Equal(t, true, d.Get("is_default"))
	assert.Equal(t, true, d.Get("is_enabled"))
	assert.Equal(t, "inherited", d.Get("customization_type"))
	projects := d.Get("projects").(*schema.Set)
	assert.NotNil(t, projects)
	assert.Equal(t, 0, projects.Len())
}

func TestProcesses_Update_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	parentID := uuid.New()
	typeID := uuid.New()
	name := "MyProcess"
	description := "My Process Description"
	referenceName := "My.Process.ReferenceName"
	expectedProcessInfo := &workitemtrackingprocess.ProcessInfo{
		TypeId:              &typeID,
		Name:                &name,
		Description:         &description,
		IsDefault:           converter.Bool(true),
		IsEnabled:           converter.Bool(false),
		CustomizationType:   &workitemtrackingprocess.CustomizationTypeValues.Inherited,
		ParentProcessTypeId: &parentID,
		ReferenceName:       &referenceName,
	}

	mockClient.EXPECT().EditProcess(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.EditProcessArgs) (*workitemtrackingprocess.ProcessInfo, error) {
			assert.Equal(t, typeID, *args.ProcessTypeId)
			assert.Equal(t, name, *args.UpdateRequest.Name)
			assert.Equal(t, description, *args.UpdateRequest.Description)
			assert.True(t, *args.UpdateRequest.IsDefault)
			assert.False(t, *args.UpdateRequest.IsEnabled)
			return expectedProcessInfo, nil
		},
	).Times(1)

	mockClient.EXPECT().GetProcessByItsId(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtrackingprocess.GetProcessByItsIdArgs) (*workitemtrackingprocess.ProcessInfo, error) {
			assert.Equal(t, typeID, *args.ProcessTypeId)
			return expectedProcessInfo, nil
		},
	).Times(1)

	d := getProcessResourceData(t, map[string]any{
		"name":                   name,
		"parent_process_type_id": parentID.String(),
		"description":            description,
		"reference_name":         referenceName,
		"is_default":             true,
		"is_enabled":             false,
	})
	d.SetId(typeID.String())

	diags := updateResourceProcess(context.Background(), d, clients)

	assert.Empty(t, diags)
	assert.Equal(t, typeID.String(), d.Id())
	assert.Equal(t, name, d.Get("name"))
	assert.Equal(t, description, d.Get("description"))
	assert.Equal(t, referenceName, d.Get("reference_name"))
	assert.Equal(t, parentID.String(), d.Get("parent_process_type_id"))
	assert.Equal(t, true, d.Get("is_default"))
	assert.Equal(t, false, d.Get("is_enabled"))
	assert.Equal(t, "inherited", d.Get("customization_type"))
	projects := d.Get("projects").(*schema.Set)
	assert.NotNil(t, projects)
	assert.Equal(t, 0, projects.Len())
}

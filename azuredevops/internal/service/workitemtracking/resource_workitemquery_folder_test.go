package workitemtracking

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// Helper to build *schema.ResourceData for the query folder resource
func getQueryFolderResourceData(t *testing.T, input map[string]interface{}) *schema.ResourceData {
	r := ResourceQueryFolder()
	return schema.TestResourceDataRaw(t, r.Schema, input)
}

// Test successful create + read flow for a query folder using parent_id.
func TestResourceQueryFolder_CreateRead_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	projectID := uuid.NewString()
	parentID := uuid.NewString()
	folderID := uuid.NewString()
	name := "MyFolder"

	// Expect CreateQuery. We only assert that we receive the correct parent Query reference.
	mockClient.EXPECT().CreateQuery(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtracking.CreateQueryArgs) (*workitemtracking.QueryHierarchyItem, error) {
			assert.Equal(t, parentID, *args.Query)
			// Return a response with Id so create succeeds
			return &workitemtracking.QueryHierarchyItem{Id: converter.UUID(folderID)}, nil
		},
	).Times(1)

	// Expect GetQuery called during the read after create. Return IsFolder = true so read passes.
	mockClient.EXPECT().GetQuery(clients.Ctx, gomock.Any()).Return(&workitemtracking.QueryHierarchyItem{
		Id:       converter.UUID(folderID),
		Name:     converter.String(name),
		IsFolder: converter.Bool(true),
	}, nil).Times(1)

	d := getQueryFolderResourceData(t, map[string]interface{}{
		"name":       name,
		"project_id": projectID,
		"parent_id":  parentID,
	})

	diags := resourceQueryFolderCreate(context.Background(), d, clients)
	assert.Empty(t, diags, "expected no diagnostics on successful create")
	assert.Equal(t, folderID, d.Id())
	assert.Equal(t, name, d.Get("name"))
}

// Test that specifying both area and parent_id returns an error diagnostic.
func TestResourceQueryFolder_Create_Error_AreaAndParentBothSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	d := getQueryFolderResourceData(t, map[string]interface{}{
		"name":       "Invalid",
		"project_id": uuid.NewString(),
		"parent_id":  uuid.NewString(),
		"area":       "Shared Queries",
	})

	diags := resourceQueryFolderCreate(context.Background(), d, clients)

	assert.Len(t, diags, 1)
	assert.Contains(t, diags[0].Summary, "Only one of 'area' or 'parent_id'")
}

// Test read sets ID to empty when underlying API indicates the folder isn't found (404).
func TestResourceQueryFolder_Read_NotFound_ClearsID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	projectID := uuid.NewString()
	folderID := uuid.NewString()

	d := getQueryFolderResourceData(t, map[string]interface{}{
		"name":       "MyFolder",
		"project_id": projectID,
		"parent_id":  uuid.NewString(),
	})
	d.SetId(folderID)

	code := 404
	mockClient.EXPECT().GetQuery(clients.Ctx, gomock.Any()).Return(nil, azuredevops.WrappedError{StatusCode: &code}).Times(1)

	diags := resourceQueryFolderRead(context.Background(), d, clients)

	assert.Empty(t, diags, "not found should not return diags")
	assert.Empty(t, d.Id(), "ID should be cleared when not found")
}

// Test read produces an error diagnostic when the returned item is not a folder.
func TestResourceQueryFolder_Read_Error_WhenNotFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	projectID := uuid.NewString()
	folderID := uuid.NewString()

	d := getQueryFolderResourceData(t, map[string]interface{}{
		"name":       "MyFolder",
		"project_id": projectID,
		"parent_id":  uuid.NewString(),
	})
	d.SetId(folderID)

	mockClient.EXPECT().GetQuery(clients.Ctx, gomock.Any()).Return(&workitemtracking.QueryHierarchyItem{
		Id:       converter.UUID(folderID),
		Name:     converter.String("my folder"),
		IsFolder: converter.Bool(false), // Not a folder -> should trigger error
	}, nil).Times(1)

	diags := resourceQueryFolderRead(context.Background(), d, clients)

	assert.Len(t, diags, 1)
	assert.Contains(t, diags[0].Summary, "not a folder")
}

// Test delete invokes DeleteQuery without diagnostics.
func TestResourceQueryFolder_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	projectID := uuid.NewString()
	folderID := uuid.NewString()

	d := getQueryFolderResourceData(t, map[string]interface{}{
		"name":       "DeleteMe",
		"project_id": projectID,
		"parent_id":  uuid.NewString(),
	})
	d.SetId(folderID)

	mockClient.EXPECT().DeleteQuery(clients.Ctx, gomock.Any()).Return(nil).Times(1)

	diags := resourceQueryFolderDelete(context.Background(), d, clients)

	assert.Empty(t, diags, "expected delete to succeed without diagnostics")
}

// Defensive test: calling read with no ID should produce a diagnostic error.
func TestResourceQueryFolder_Read_Error_NoID(t *testing.T) {
	d := getQueryFolderResourceData(t, map[string]interface{}{
		"name":       "NoId",
		"project_id": uuid.NewString(),
		"parent_id":  uuid.NewString(),
	})
	diags := resourceQueryFolderRead(context.Background(), d, &client.AggregatedClient{})

	assert.Len(t, diags, 1)
	assert.IsType(t, diag.Diagnostics{}, diags)
}

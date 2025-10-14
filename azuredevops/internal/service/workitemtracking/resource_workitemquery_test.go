package workitemtracking

import (
	"context"
	"fmt"
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

// Helper to build *schema.ResourceData for the query resource
func getQueryResourceData(t *testing.T, input map[string]interface{}) *schema.ResourceData {
	r := ResourceQuery()
	return schema.TestResourceDataRaw(t, r.Schema, input)
}

// Test create and read succeeds.
func TestResourceQuery_CreateRead_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	projectID := uuid.NewString()
	parentID := uuid.NewString()
	queryID := uuid.NewString()
	name := "MyQuery"
	wiql := "SELECT [System.Id] FROM WorkItems"

	mockClient.EXPECT().CreateQuery(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtracking.CreateQueryArgs) (*workitemtracking.QueryHierarchyItem, error) {
			assert.Equal(t, parentID, *args.Query)
			return &workitemtracking.QueryHierarchyItem{Id: converter.UUID(queryID)}, nil
		},
	).Times(1)

	mockClient.EXPECT().GetQuery(clients.Ctx, gomock.Any()).Return(&workitemtracking.QueryHierarchyItem{
		Id:       converter.UUID(queryID),
		Name:     converter.String(name),
		Wiql:     converter.String(wiql),
		IsFolder: converter.Bool(false),
	}, nil).Times(1)

	d := getQueryResourceData(t, map[string]interface{}{
		"name":      name,
		"parent_id": parentID,
		"wiql":      wiql,
		// project_id is accessed dynamically; some resources rely on provider level. Simulate presence.
		"project_id": projectID,
	})

	diags := resourceQueryCreate(context.Background(), d, clients)

	assert.Empty(t, diags)
	assert.Equal(t, queryID, d.Id())
	assert.Equal(t, wiql, d.Get("wiql"))
}

// Test read marks resource removed on 404.
func TestResourceQuery_Read_NotFound_ClearsID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	projectID := uuid.NewString()
	queryID := uuid.NewString()

	d := getQueryResourceData(t, map[string]interface{}{
		"name":       "MyQuery",
		"wiql":       "SELECT 1",
		"parent_id":  uuid.NewString(),
		"project_id": projectID,
	})

	d.SetId(queryID)
	code := 404

	mockClient.EXPECT().GetQuery(clients.Ctx, gomock.Any()).Return(nil, azuredevops.WrappedError{StatusCode: &code}).Times(1)

	diags := resourceQueryRead(context.Background(), d, clients)

	assert.Empty(t, diags)
	assert.Empty(t, d.Id())
}

// Test read errors when a folder is returned instead of a query.
func TestResourceQuery_Read_Error_WhenFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	projectID := uuid.NewString()
	queryID := uuid.NewString()

	d := getQueryResourceData(t, map[string]interface{}{
		"name":       "MyQuery",
		"wiql":       "SELECT 1",
		"parent_id":  uuid.NewString(),
		"project_id": projectID,
	})
	d.SetId(queryID)

	mockClient.EXPECT().GetQuery(clients.Ctx, gomock.Any()).Return(&workitemtracking.QueryHierarchyItem{
		Id:       converter.UUID(queryID),
		Name:     converter.String("MyQuery"),
		Wiql:     converter.String("SELECT 1"),
		IsFolder: converter.Bool(true),
	}, nil).Times(1)

	diags := resourceQueryRead(context.Background(), d, clients)

	assert.Len(t, diags, 1)
	assert.Contains(t, diags[0].Summary, "is a folder")
}

// Test update changes wiql and name when both changed.
func TestResourceQuery_Update_ModifiesFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	projectID := uuid.NewString()
	queryID := uuid.NewString()

	oldWiql := "SELECT 1"
	newWiql := "SELECT 2"
	oldName := "OldName"
	newName := "NewName"

	d := getQueryResourceData(t, map[string]interface{}{
		"name":       oldName,
		"wiql":       oldWiql,
		"parent_id":  uuid.NewString(),
		"project_id": projectID,
	})
	d.SetId(queryID)

	// Simulate change
	d.Set("wiql", newWiql)
	d.Set("name", newName)

	// First GetQuery to fetch existing
	mockClient.EXPECT().GetQuery(clients.Ctx, gomock.Any()).Return(&workitemtracking.QueryHierarchyItem{
		Id:       converter.UUID(queryID),
		Name:     converter.String(oldName),
		Wiql:     converter.String(oldWiql),
		IsFolder: converter.Bool(false),
	}, nil).Times(1)

	// UpdateQuery should be called; validate mutated content
	mockClient.EXPECT().UpdateQuery(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtracking.UpdateQueryArgs) (*workitemtracking.QueryHierarchyItem, error) {
			assert.Equal(t, newName, *args.QueryUpdate.Name)
			assert.Equal(t, newWiql, *args.QueryUpdate.Wiql)
			return args.QueryUpdate, nil
		},
	).Times(1)

	// Final read after update
	mockClient.EXPECT().GetQuery(clients.Ctx, gomock.Any()).Return(&workitemtracking.QueryHierarchyItem{
		Id:       converter.UUID(queryID),
		Name:     converter.String(newName),
		Wiql:     converter.String(newWiql),
		IsFolder: converter.Bool(false),
	}, nil).Times(1)

	diags := resourceQueryUpdate(context.Background(), d, clients)

	assert.Empty(t, diags)
	assert.Equal(t, newWiql, d.Get("wiql"))
	assert.Equal(t, newName, d.Get("name"))
}

// Test delete calls DeleteQuery.
func TestResourceQuery_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	projectID := uuid.NewString()
	queryID := uuid.NewString()

	d := getQueryResourceData(t, map[string]interface{}{
		"name":       "del",
		"wiql":       "SELECT 1",
		"parent_id":  uuid.NewString(),
		"project_id": projectID,
	})
	d.SetId(queryID)

	mockClient.EXPECT().DeleteQuery(clients.Ctx, gomock.Any()).Return(nil).Times(1)

	diags := resourceQueryDelete(context.Background(), d, clients)

	assert.Empty(t, diags)
}

// Defensive test: read with no ID returns diagnostic.
func TestResourceQuery_Read_Error_NoID(t *testing.T) {
	d := getQueryResourceData(t, map[string]interface{}{
		"name":       "test",
		"wiql":       "SELECT 1",
		"parent_id":  uuid.NewString(),
		"project_id": uuid.NewString(),
	})
	diags := resourceQueryRead(context.Background(), d, &client.AggregatedClient{})
	assert.Len(t, diags, 1)
	assert.IsType(t, diag.Diagnostics{}, diags)
}

// Error path for update when UpdateQuery returns failure.
func TestResourceQuery_Update_Error_OnUpdateFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	projectID := uuid.NewString()
	queryID := uuid.NewString()
	oldWiql := "SELECT 1"

	d := getQueryResourceData(t, map[string]interface{}{
		"name":       "Name",
		"wiql":       oldWiql,
		"parent_id":  uuid.NewString(),
		"project_id": projectID,
	})
	d.SetId(queryID)
	d.Set("wiql", "SELECT 2")

	mockClient.EXPECT().GetQuery(clients.Ctx, gomock.Any()).Return(&workitemtracking.QueryHierarchyItem{
		Id:       converter.UUID(queryID),
		Name:     converter.String("Name"),
		Wiql:     converter.String(oldWiql),
		IsFolder: converter.Bool(false),
	}, nil).Times(1)

	mockClient.EXPECT().UpdateQuery(clients.Ctx, gomock.Any()).Return(nil, fmt.Errorf("update failed")).Times(1)

	diags := resourceQueryUpdate(context.Background(), d, clients)
	assert.Len(t, diags, 1)
	assert.Contains(t, diags[0].Summary, "Updating query")
}

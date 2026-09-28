package workitemtracking

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func getIterationResourceData(t *testing.T, input map[string]interface{}) *schema.ResourceData {
	return schema.TestResourceDataRaw(t, ResourceIteration().Schema, input)
}

func TestResourceIteration_Create_SendsDates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}
	projectID := uuid.NewString()

	mockClient.EXPECT().CreateOrUpdateClassificationNode(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtracking.CreateOrUpdateClassificationNodeArgs) (*workitemtracking.WorkItemClassificationNode, error) {
			assert.Equal(t, workitemtracking.TreeStructureGroupValues.Iterations, *args.StructureGroup)
			assert.Equal(t, "", *args.Path)
			assert.Equal(t, "Sprint 1", *args.PostedNode.Name)
			assert.Equal(t, map[string]interface{}{
				"startDate":  "2026-01-05T00:00:00Z",
				"finishDate": "2026-01-16T00:00:00Z",
			}, *args.PostedNode.Attributes)
			return &workitemtracking.WorkItemClassificationNode{Id: converter.Int(42)}, nil
		}).Times(1)

	mockClient.EXPECT().GetClassificationNodes(clients.Ctx, gomock.Any()).Return(&[]workitemtracking.WorkItemClassificationNode{{
		Id:   converter.Int(42),
		Name: converter.String("Sprint 1"),
		Path: converter.String(`\example-project\Iteration\Sprint 1`),
		Attributes: &map[string]interface{}{
			"startDate":  "2026-01-05T00:00:00Z",
			"finishDate": "2026-01-16T00:00:00Z",
		},
	}}, nil).Times(1)

	d := getIterationResourceData(t, map[string]interface{}{
		"project_id":  projectID,
		"name":        "Sprint 1",
		"start_date":  "2026-01-05",
		"finish_date": "2026-01-16",
	})

	diags := resourceIterationCreate(context.Background(), d, clients)
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, "42", d.Id())
	assert.Equal(t, `example-project\Sprint 1`, d.Get("path"))
	assert.Equal(t, "2026-01-05", d.Get("start_date"))
	assert.Equal(t, "2026-01-16", d.Get("finish_date"))
}

func TestResourceIteration_Create_WithoutDates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	mockClient.EXPECT().CreateOrUpdateClassificationNode(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtracking.CreateOrUpdateClassificationNodeArgs) (*workitemtracking.WorkItemClassificationNode, error) {
			assert.Nil(t, args.PostedNode.Attributes)
			return &workitemtracking.WorkItemClassificationNode{Id: converter.Int(42)}, nil
		}).Times(1)
	mockClient.EXPECT().GetClassificationNodes(clients.Ctx, gomock.Any()).Return(&[]workitemtracking.WorkItemClassificationNode{{
		Id:   converter.Int(42),
		Name: converter.String("Sprint 1"),
		Path: converter.String(`\example-project\Iteration\Sprint 1`),
	}}, nil).Times(1)

	d := getIterationResourceData(t, map[string]interface{}{
		"project_id": uuid.NewString(),
		"name":       "Sprint 1",
	})

	diags := resourceIterationCreate(context.Background(), d, clients)
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, "", d.Get("start_date"))
	assert.Equal(t, "", d.Get("finish_date"))
}

func TestResourceIteration_Read_NotFound_ClearsID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	mockClient.EXPECT().GetClassificationNodes(clients.Ctx, gomock.Any()).Return(nil, azuredevops.WrappedError{
		StatusCode: converter.Int(http.StatusNotFound),
	}).Times(1)

	d := getIterationResourceData(t, map[string]interface{}{
		"project_id": uuid.NewString(),
		"name":       "Sprint 1",
	})
	d.SetId("42")

	diags := resourceIterationRead(context.Background(), d, clients)
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, "", d.Id())
}

func TestResourceIteration_Delete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	mockClient.EXPECT().DeleteClassificationNode(clients.Ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, args workitemtracking.DeleteClassificationNodeArgs) error {
			assert.Equal(t, workitemtracking.TreeStructureGroupValues.Iterations, *args.StructureGroup)
			assert.Equal(t, "Sprint 1", *args.Path)
			return errors.New("boom")
		}).Times(1)

	d := getIterationResourceData(t, map[string]interface{}{
		"project_id": uuid.NewString(),
		"name":       "Sprint 1",
	})
	d.SetId("42")

	diags := resourceIterationDelete(context.Background(), d, clients)
	require.True(t, diags.HasError())
	assert.Contains(t, diags[0].Summary, "boom")
}

func TestResourceIteration_ExpandIterationDates_ClearsUnset(t *testing.T) {
	d := getIterationResourceData(t, map[string]interface{}{
		"project_id": uuid.NewString(),
		"name":       "Sprint 1",
	})
	assert.Equal(t, map[string]interface{}{"startDate": nil, "finishDate": nil}, *expandIterationDates(d))
}

func TestResourceIteration_ValidateIterationDate(t *testing.T) {
	for _, v := range []string{"2026-01-05", "2024-02-29"} {
		_, errs := validateIterationDate(v, "start_date")
		assert.Empty(t, errs, v)
	}
	for _, v := range []string{"", "2026-1-5", "05/01/2026", "2026-01-05T00:00:00Z", "2026-02-30"} {
		_, errs := validateIterationDate(v, "start_date")
		assert.NotEmpty(t, errs, v)
	}
}

func TestResourceIteration_FlattenIterationDate(t *testing.T) {
	for in, want := range map[interface{}]string{
		nil:                    "",
		"":                     "",
		"2026-01-05T00:00:00Z": "2026-01-05",
	} {
		got, err := flattenIterationDate(in)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	}
	_, err := flattenIterationDate("not-a-date")
	assert.Error(t, err)
}

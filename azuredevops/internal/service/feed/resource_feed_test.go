//go:build (all || resource_feed) && !exclude_feed
// +build all resource_feed
// +build !exclude_feed

package feed

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/feed"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/testhelper"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	FeedName      = "some-feed-name"
	FeedProjectId = uuid.New().String()
)

// verifies that if an error is produced on create, the error is not swallowed

func TestFeed_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceFeed()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":       FeedName,
		"project_id": FeedProjectId,
	})

	feedClient := azdosdkmocks.NewMockFeedClient(ctrl)
	clients := &client.AggregatedClient{FeedClient: feedClient, Ctx: context.Background()}

	expectedArgs := feed.CreateFeedArgs{
		Feed:    &feed.Feed{Name: &FeedName},
		Project: &FeedProjectId,
	}

	feedClient.
		EXPECT().
		CreateFeed(clients.Ctx, expectedArgs).
		Return(nil, fmt.Errorf("Name already exists")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Name already exists")
}

// verifies that if an error is produced on update, the error is not swallowed

func TestFeed_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceFeed()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":       FeedName,
		"project_id": FeedProjectId,
	})

	feedClient := azdosdkmocks.NewMockFeedClient(ctrl)
	clients := &client.AggregatedClient{FeedClient: feedClient, Ctx: context.Background()}

	expectedArgs := feed.UpdateFeedArgs{
		Feed:    &feed.FeedUpdate{},
		FeedId:  &FeedName,
		Project: &FeedProjectId,
	}

	feedClient.
		EXPECT().
		UpdateFeed(clients.Ctx, expectedArgs).
		Return(nil, fmt.Errorf("Feed with given name not found")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Feed with given name not found")
}

// verifies that if an error is produced on delete, the error is not swallowed

func TestFeed_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceFeed()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":       FeedName,
		"project_id": FeedProjectId,
	})

	feedClient := azdosdkmocks.NewMockFeedClient(ctrl)
	clients := &client.AggregatedClient{FeedClient: feedClient, Ctx: context.Background()}

	expectedDeleteArgs := feed.DeleteFeedArgs{
		FeedId:  &FeedName,
		Project: &FeedProjectId,
	}

	feedClient.
		EXPECT().
		DeleteFeed(clients.Ctx, expectedDeleteArgs).
		Return(fmt.Errorf("Feed with given name not found")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Feed with given name not found")
}

func TestFeed_Create_WithUpstreamSources(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceFeed()

	serviceEndpointID := uuid.New()
	serviceEndpointProjectID := uuid.New()
	internalViewID := uuid.New()
	internalFeedID := uuid.New()

	resourceData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":       FeedName,
		"project_id": FeedProjectId,
		"upstream_sources": []interface{}{
			map[string]interface{}{
				"name":                        "upstream-test",
				"protocol":                    "nuget",
				"location":                    "https://upstream.test",
				"upstream_source_type":        "internal",
				"service_endpoint_id":         serviceEndpointID.String(),
				"service_endpoint_project_id": serviceEndpointProjectID.String(),
				"internal_upstream_view_id":   internalViewID.String(),
				"internal_upstream_feed_id":   internalFeedID.String(),
			},
		},
	})

	feedClient := azdosdkmocks.NewMockFeedClient(ctrl)
	clients := &client.AggregatedClient{FeedClient: feedClient, Ctx: context.Background()}

	expectedArgs := feed.CreateFeedArgs{
		Feed: &feed.Feed{
			Name:            &FeedName,
			UpstreamEnabled: converter.Bool(true),
			UpstreamSources: &[]feed.UpstreamSource{
				{
					Name:                     converter.String("upstream-test"),
					Protocol:                 converter.String("nuget"),
					Location:                 converter.String("https://upstream.test"),
					UpstreamSourceType:       &feed.UpstreamSourceTypeValues.Internal,
					ServiceEndpointId:        &serviceEndpointID,
					ServiceEndpointProjectId: &serviceEndpointProjectID,
					InternalUpstreamViewId:   &internalViewID,
					InternalUpstreamFeedId:   &internalFeedID,
				},
			},
		},
		Project: &FeedProjectId,
	}

	feedID := uuid.New()
	createdFeed := feed.Feed{
		Id:              &feedID,
		Name:            &FeedName,
		Project:         &feed.ProjectReference{Id: converter.UUID(FeedProjectId)},
		UpstreamSources: expectedArgs.Feed.UpstreamSources,
	}

	feedClient.
		EXPECT().
		CreateFeed(clients.Ctx, expectedArgs).
		Return(&createdFeed, nil).
		Times(1)

	feedClient.
		EXPECT().
		GetFeed(clients.Ctx, feed.GetFeedArgs{FeedId: converter.String(feedID.String()), Project: &FeedProjectId}).
		Return(&createdFeed, nil).
		Times(1)

	err := r.Create(resourceData, clients)
	require.NoError(t, err)
	require.Equal(t, feedID.String(), resourceData.Id())
}

func TestFeed_Update_WithUpstreamSources(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceFeed()

	serviceEndpointID := uuid.New()

	resourceData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":       FeedName,
		"project_id": FeedProjectId,
		"upstream_sources": []interface{}{
			map[string]interface{}{
				"name":                 "upstream-update",
				"protocol":             "npm",
				"location":             "https://registry.npmjs.org",
				"upstream_source_type": "public",
				"service_endpoint_id":  serviceEndpointID.String(),
			},
		},
	})
	resourceData.SetId(FeedName)

	feedClient := azdosdkmocks.NewMockFeedClient(ctrl)
	clients := &client.AggregatedClient{FeedClient: feedClient, Ctx: context.Background()}

	expectedArgs := feed.UpdateFeedArgs{
		Feed: &feed.FeedUpdate{
			UpstreamEnabled: converter.Bool(true),
			UpstreamSources: &[]feed.UpstreamSource{
				{
					Name:               converter.String("upstream-update"),
					Protocol:           converter.String("npm"),
					Location:           converter.String("https://registry.npmjs.org"),
					UpstreamSourceType: &feed.UpstreamSourceTypeValues.Public,
					ServiceEndpointId:  &serviceEndpointID,
				},
			},
		},
		FeedId:  &FeedName,
		Project: &FeedProjectId,
	}

	updatedFeed := feed.Feed{
		Id:              testhelper.CreateUUID(),
		Name:            &FeedName,
		Project:         &feed.ProjectReference{Id: converter.UUID(FeedProjectId)},
		UpstreamSources: expectedArgs.Feed.UpstreamSources,
	}

	feedClient.
		EXPECT().
		UpdateFeed(clients.Ctx, expectedArgs).
		Return(&updatedFeed, nil).
		Times(1)

	feedClient.
		EXPECT().
		GetFeed(clients.Ctx, feed.GetFeedArgs{FeedId: &FeedName, Project: &FeedProjectId}).
		Return(&updatedFeed, nil).
		Times(1)

	err := r.Update(resourceData, clients)
	require.NoError(t, err)
}

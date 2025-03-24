//go:build (all || resource_feed) && !exclude_feed
// +build all resource_feed
// +build !exclude_feed

package feed

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/feed"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var FeedName = "some-feed-name"
var FeedProjectId = uuid.New().String()

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

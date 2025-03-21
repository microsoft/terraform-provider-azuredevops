//go:build (all || resource_feed) && !exclude_feed
// +build all resource_feed
// +build !exclude_feed

package feed

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/feed"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var FeedId = uuid.New().String()
var ProjectId = uuid.New().String()
var IdentityDescriptor = "some-identity-descriptor"
var IdentityLegacyDescriptor = "some-legacy-identity-descriptor"
var IdentityId = uuid.New()
var Role = "reader"

// verifies that if an error is produced on create, the error is not swallowed

func TestFeedPermission_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceFeedPermission()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"feed_id":             FeedId,
		"project_id":          ProjectId,
		"identity_descriptor": IdentityDescriptor,
		"role":                Role,
	})

	feedClient := azdosdkmocks.NewMockFeedClient(ctrl)
	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{
		FeedClient:     feedClient,
		IdentityClient: identityClient,
		GraphClient:    graphClient,
		Ctx:            context.Background(),
	}

	feed_id := resourceData.Get("feed_id").(string)
	identity_descriptor := resourceData.Get("identity_descriptor").(string)
	role := feed.FeedRole(resourceData.Get("role").(string))
	project_id := resourceData.Get("project_id").(string)
	display_name := resourceData.Get("display_name").(string)

	graphClient.
		EXPECT().
		GetStorageKey(clients.Ctx, graph.GetStorageKeyArgs{
			SubjectDescriptor: &identity_descriptor,
		}).
		Return(&graph.GraphStorageKeyResult{
			Value: &IdentityId,
		}, nil).
		Times(1)

	identityClient.
		EXPECT().
		ReadIdentity(clients.Ctx, identity.ReadIdentityArgs{
			IdentityId: converter.String(IdentityId.String()),
		}).
		Return(&identity.Identity{
			Id:         &IdentityId,
			Descriptor: &IdentityLegacyDescriptor,
		}, nil).
		Times(1)

	feedClient.
		EXPECT().
		GetFeedPermissions(clients.Ctx, feed.GetFeedPermissionsArgs{
			FeedId:  &feed_id,
			Project: &project_id,
		}).
		Return(&[]feed.FeedPermission{}, nil).
		Times(1)

	expectedArgs := feed.SetFeedPermissionsArgs{
		FeedId:  &feed_id,
		Project: &project_id,
		FeedPermission: &[]feed.FeedPermission{
			{
				DisplayName:        &display_name,
				IdentityDescriptor: &IdentityLegacyDescriptor,
				IdentityId:         &IdentityId,
				Role:               &role,
			},
		},
	}

	feedClient.
		EXPECT().
		SetFeedPermissions(clients.Ctx, expectedArgs).
		Return(nil, fmt.Errorf("Something unexpected happened")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Something unexpected happened")
}

// verifies that if an error is produced on update, the error is not swallowed

func TestFeedPermission_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceFeedPermission()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"feed_id":             FeedId,
		"project_id":          ProjectId,
		"identity_descriptor": IdentityDescriptor,
		"role":                Role,
	})

	feedClient := azdosdkmocks.NewMockFeedClient(ctrl)
	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{
		FeedClient:     feedClient,
		IdentityClient: identityClient,
		GraphClient:    graphClient,
		Ctx:            context.Background(),
	}

	feed_id := resourceData.Get("feed_id").(string)
	identity_descriptor := resourceData.Get("identity_descriptor").(string)
	role := feed.FeedRole(resourceData.Get("role").(string))
	project_id := resourceData.Get("project_id").(string)
	display_name := resourceData.Get("display_name").(string)

	graphClient.
		EXPECT().
		GetStorageKey(clients.Ctx, graph.GetStorageKeyArgs{
			SubjectDescriptor: &identity_descriptor,
		}).
		Return(&graph.GraphStorageKeyResult{
			Value: &IdentityId,
		}, nil).
		Times(1)

	identityClient.
		EXPECT().
		ReadIdentity(clients.Ctx, identity.ReadIdentityArgs{
			IdentityId: converter.String(IdentityId.String()),
		}).
		Return(&identity.Identity{
			Id:         &IdentityId,
			Descriptor: &IdentityLegacyDescriptor,
		}, nil).
		Times(1)

	feedClient.
		EXPECT().
		GetFeedPermissions(clients.Ctx, feed.GetFeedPermissionsArgs{
			FeedId:  &feed_id,
			Project: &project_id,
		}).
		Return(&[]feed.FeedPermission{{
			DisplayName:        &display_name,
			IdentityDescriptor: &IdentityLegacyDescriptor,
			IdentityId:         &IdentityId,
			Role:               &role,
		}}, nil).
		Times(1)

	expectedArgs := feed.SetFeedPermissionsArgs{
		FeedId:  &feed_id,
		Project: &project_id,
		FeedPermission: &[]feed.FeedPermission{
			{
				DisplayName:        &display_name,
				IdentityDescriptor: &IdentityLegacyDescriptor,
				IdentityId:         &IdentityId,
				Role:               &role,
			},
		},
	}

	feedClient.
		EXPECT().
		SetFeedPermissions(clients.Ctx, expectedArgs).
		Return(nil, fmt.Errorf("Something unexpected happened")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Something unexpected happened")
}

// verifies that if an error is produced on read, the error is not swallowed

func TestFeedPermission_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceFeedPermission()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"feed_id":             FeedId,
		"project_id":          ProjectId,
		"identity_descriptor": IdentityDescriptor,
		"role":                Role,
	})

	feedClient := azdosdkmocks.NewMockFeedClient(ctrl)
	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{
		FeedClient:     feedClient,
		IdentityClient: identityClient,
		GraphClient:    graphClient,
		Ctx:            context.Background(),
	}

	feed_id := resourceData.Get("feed_id").(string)
	identity_descriptor := resourceData.Get("identity_descriptor").(string)
	project_id := resourceData.Get("project_id").(string)

	graphClient.
		EXPECT().
		GetStorageKey(clients.Ctx, graph.GetStorageKeyArgs{
			SubjectDescriptor: &identity_descriptor,
		}).
		Return(&graph.GraphStorageKeyResult{
			Value: &IdentityId,
		}, nil).
		Times(1)

	identityClient.
		EXPECT().
		ReadIdentity(clients.Ctx, identity.ReadIdentityArgs{
			IdentityId: converter.String(IdentityId.String()),
		}).
		Return(&identity.Identity{
			Id:         &IdentityId,
			Descriptor: &IdentityLegacyDescriptor,
		}, nil).
		Times(1)

	feedClient.
		EXPECT().
		GetFeedPermissions(clients.Ctx, feed.GetFeedPermissionsArgs{
			FeedId:  &feed_id,
			Project: &project_id,
		}).
		Return(nil, fmt.Errorf("Something unexpected happened")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Something unexpected happened")
}

// verifies that if an error is produced on delete, the error is not swallowed

func TestFeedPermission_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceFeedPermission()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"feed_id":             FeedId,
		"project_id":          ProjectId,
		"identity_descriptor": IdentityDescriptor,
		"role":                Role,
	})

	feedClient := azdosdkmocks.NewMockFeedClient(ctrl)
	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{
		FeedClient:     feedClient,
		IdentityClient: identityClient,
		GraphClient:    graphClient,
		Ctx:            context.Background(),
	}

	feed_id := resourceData.Get("feed_id").(string)
	identity_descriptor := resourceData.Get("identity_descriptor").(string)
	role := feed.FeedRoleValues.None
	project_id := resourceData.Get("project_id").(string)

	graphClient.
		EXPECT().
		GetStorageKey(clients.Ctx, graph.GetStorageKeyArgs{
			SubjectDescriptor: &identity_descriptor,
		}).
		Return(&graph.GraphStorageKeyResult{
			Value: &IdentityId,
		}, nil).
		Times(1)

	identityClient.
		EXPECT().
		ReadIdentity(clients.Ctx, identity.ReadIdentityArgs{
			IdentityId: converter.String(IdentityId.String()),
		}).
		Return(&identity.Identity{
			Id:         &IdentityId,
			Descriptor: &IdentityLegacyDescriptor,
		}, nil).
		Times(1)

	expectedArgs := feed.SetFeedPermissionsArgs{
		FeedId:  &feed_id,
		Project: &project_id,
		FeedPermission: &[]feed.FeedPermission{
			{
				IdentityDescriptor: &IdentityLegacyDescriptor,
				Role:               &role,
			},
		},
	}

	feedClient.
		EXPECT().
		SetFeedPermissions(clients.Ctx, expectedArgs).
		Return(nil, fmt.Errorf("Something unexpected happened")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Something unexpected happened")
}

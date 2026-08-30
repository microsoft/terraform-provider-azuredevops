package git

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGitRefResource_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &client.AggregatedClient{
		GitReposClient: repoClient,
		Ctx:            context.Background(),
	}

	repoId := "00000000-0000-0000-0000-000000000000"
	refName := "refs/heads/new-branch"
	sourceObjectId := "abcd1234"

	// Mock GetRefs to find the source branch
	repoClient.
		EXPECT().
		GetRefs(clients.Ctx, gomock.Any()).
		Return(&git.GetRefsResponseValue{
			Value: []git.GitRef{
				{
					Name:     converter.String("refs/heads/main"),
					ObjectId: &sourceObjectId,
				},
			},
		}, nil).
		Times(1)

	// Mock UpdateRefs to create the new branch
	repoClient.
		EXPECT().
		UpdateRefs(clients.Ctx, gomock.Any()).
		Return(&[]git.GitRefUpdateResult{
			{
				Success:      converter.Bool(true),
				UpdateStatus: &git.GitRefUpdateStatusValues.Succeeded,
			},
		}, nil).
		Times(1)

	// Mock GetRefs for the Read call after Create
	repoClient.
		EXPECT().
		GetRefs(clients.Ctx, gomock.Any()).
		Return(&git.GetRefsResponseValue{
			Value: []git.GitRef{
				{
					Name:     &refName,
					ObjectId: &sourceObjectId,
				},
			},
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceGitRef().Schema, nil)
	resourceData.Set("repository_id", repoId)
	resourceData.Set("name", refName)
	resourceData.Set("ref_branch", "main")

	err := resourceGitRefCreate(clients.Ctx, resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("%s:%s", repoId, refName), resourceData.Id())
}

func TestGitRefResource_Create_WithProjectID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &client.AggregatedClient{
		GitReposClient: repoClient,
		Ctx:            context.Background(),
	}

	repoId := "00000000-0000-0000-0000-000000000000"
	projectId := "11111111-1111-1111-1111-111111111111"
	refName := "refs/heads/new-branch"
	sourceObjectId := "abcd1234"

	// Mock GetRefs to find the source branch
	repoClient.
		EXPECT().
		GetRefs(clients.Ctx, gomock.Any()).
		Do(func(ctx context.Context, args git.GetRefsArgs) {
			require.Equal(t, projectId, *args.Project)
		}).
		Return(&git.GetRefsResponseValue{
			Value: []git.GitRef{
				{
					Name:     converter.String("refs/heads/main"),
					ObjectId: &sourceObjectId,
				},
			},
		}, nil).
		Times(1)

	// Mock UpdateRefs to create the new branch
	repoClient.
		EXPECT().
		UpdateRefs(clients.Ctx, gomock.Any()).
		Do(func(ctx context.Context, args git.UpdateRefsArgs) {
			require.Equal(t, projectId, *args.Project)
		}).
		Return(&[]git.GitRefUpdateResult{
			{
				Success:      converter.Bool(true),
				UpdateStatus: &git.GitRefUpdateStatusValues.Succeeded,
			},
		}, nil).
		Times(1)

	// Mock GetRefs for the Read call after Create
	repoClient.
		EXPECT().
		GetRefs(clients.Ctx, gomock.Any()).
		Do(func(ctx context.Context, args git.GetRefsArgs) {
			require.Equal(t, projectId, *args.Project)
		}).
		Return(&git.GetRefsResponseValue{
			Value: []git.GitRef{
				{
					Name:     &refName,
					ObjectId: &sourceObjectId,
				},
			},
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceGitRef().Schema, nil)
	resourceData.Set("repository_id", repoId)
	resourceData.Set("project_id", projectId)
	resourceData.Set("name", refName)
	resourceData.Set("ref_branch", "main")

	err := resourceGitRefCreate(clients.Ctx, resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("%s:%s", repoId, refName), resourceData.Id())
}

func TestGitRefResource_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &client.AggregatedClient{
		GitReposClient: repoClient,
		Ctx:            context.Background(),
	}

	repoId := "00000000-0000-0000-0000-000000000000"
	refName := "refs/heads/branch-to-delete"
	objectId := "1234abcd"

	resourceData := schema.TestResourceDataRaw(t, ResourceGitRef().Schema, nil)
	resourceData.SetId(fmt.Sprintf("%s:%s", repoId, refName))
	resourceData.Set("object_id", objectId)

	repoClient.
		EXPECT().
		UpdateRefs(clients.Ctx, gomock.Any()).
		Return(&[]git.GitRefUpdateResult{
			{
				Success:      converter.Bool(true),
				UpdateStatus: &git.GitRefUpdateStatusValues.Succeeded,
			},
		}, nil).
		Times(1)

	err := resourceGitRefDelete(clients.Ctx, resourceData, clients)
	require.Nil(t, err)
}

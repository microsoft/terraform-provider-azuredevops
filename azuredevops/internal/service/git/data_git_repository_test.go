//go:build (all || git || data_sources || data_git_repository) && (!exclude_data_sources || !exclude_git || !exclude_data_git_repository)
// +build all git data_sources data_git_repository
// +build !exclude_data_sources !exclude_git !exclude_data_git_repository

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
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/testhelper"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var gitRepo = git.GitRepository{
	Links:         nil,
	DefaultBranch: converter.String("master"),
	Id:            testhelper.CreateUUID(),
	IsFork:        converter.Bool(true),
	Name:          converter.String("repo-02"),
	ParentRepository: &git.GitRepositoryRef{
		Id:   testhelper.CreateUUID(),
		Name: converter.String("repo-parent-02"),
	},
	Project:         azProjectRef,
	RemoteUrl:       nil,
	Size:            converter.UInt64(0),
	SshUrl:          nil,
	Url:             nil,
	ValidRemoteUrls: nil,
	WebUrl:          nil,
	IsDisabled:      nil,
}

func TestGitRepositoryDataSource_Read_DontSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	repoClient := azdosdkmocks.NewMockGitClient(ctrl)

	clients := &client.AggregatedClient{
		GitReposClient: repoClient,
		GraphClient:    graphClient,
		Ctx:            context.Background(),
	}

	repoClient.
		EXPECT().
		GetRepository(clients.Ctx, git.GetRepositoryArgs{
			RepositoryId: gitRepo.Name,
			Project:      converter.String(gitRepo.Project.Id.String()),
		}).
		Return(nil, fmt.Errorf("@@GetRepository@@failed")).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataGitRepository().Schema, nil)
	resourceData.Set("name", gitRepo.Name)
	resourceData.Set("project_id", gitRepo.Project.Id.String())

	err := dataSourceGitRepositoryRead(resourceData, clients)
	require.NotNil(t, err)
}

func TestGitRepositoryDataSource_Read_Repository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	repoClient := azdosdkmocks.NewMockGitClient(ctrl)

	clients := &client.AggregatedClient{
		GitReposClient: repoClient,
		GraphClient:    graphClient,
		Ctx:            context.Background(),
	}

	expectedGetRepositoryArgs := git.GetRepositoryArgs{
		RepositoryId: gitRepo.Name,
		Project:      converter.String(gitRepo.Project.Id.String()),
	}
	repoClient.
		EXPECT().
		GetRepository(clients.Ctx, expectedGetRepositoryArgs).
		Return(&gitRepo, nil)

	resourceData := schema.TestResourceDataRaw(t, DataGitRepository().Schema, nil)
	resourceData.Set("name", gitRepo.Name)
	resourceData.Set("project_id", gitRepo.Project.Id.String())

	err := dataSourceGitRepositoryRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, resourceData.Id(), gitRepo.Id.String())
	require.Equal(t, resourceData.Get("name"), *gitRepo.Name)
	require.Equal(t, resourceData.Get("project_id"), gitRepo.Project.Id.String())
}

func TestGitRepositoryDataSource_Read_RepositoryNotFoundError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	repoClient := azdosdkmocks.NewMockGitClient(ctrl)

	clients := &client.AggregatedClient{
		GitReposClient: repoClient,
		GraphClient:    graphClient,
		Ctx:            context.Background(),
	}

	repoClient.
		EXPECT().
		GetRepository(clients.Ctx, gomock.Any()).
		Return(nil, fmt.Errorf("TF200016: The project does not exist"))

	resourceData := schema.TestResourceDataRaw(t, DataGitRepository().Schema, nil)
	resourceData.Set("name", "@@invalid@@")
	resourceData.Set("project_id", gitRepo.Project.Id.String())

	err := dataSourceGitRepositoryRead(resourceData, clients)
	require.NotNil(t, err)
}

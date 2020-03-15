// +build all core data_git_repositories

package azuredevops

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/testhelper"
	"github.com/stretchr/testify/require"
)

func init() {
	/* add code for test setup here */
}

var gitRepoListEmpty = []git.GitRepository{}

var azProjectRef = &core.TeamProjectReference{
	Id:   testhelper.CreateUUID(),
	Name: converter.String("project-01"),
}

var gitRepoList = []git.GitRepository{
	{
		Links:            nil,
		DefaultBranch:    nil,
		Id:               testhelper.CreateUUID(),
		IsFork:           converter.Bool(false),
		Name:             converter.String("repo-01"),
		ParentRepository: nil,
		Project:          azProjectRef,
		RemoteUrl:        nil,
		Size:             nil,
		SshUrl:           nil,
		Url:              nil,
		ValidRemoteUrls:  nil,
		WebUrl:           nil,
	},
	{
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
	},
	{
		Links:            nil,
		DefaultBranch:    converter.String("dev"),
		Id:               testhelper.CreateUUID(),
		IsFork:           nil,
		Name:             converter.String("repo-03"),
		ParentRepository: nil,
		Project: &core.TeamProjectReference{
			Id:   testhelper.CreateUUID(),
			Name: converter.String("project-02"),
		},
		RemoteUrl:       nil,
		Size:            converter.UInt64(1234),
		SshUrl:          nil,
		Url:             nil,
		ValidRemoteUrls: nil,
		WebUrl:          nil,
	},
}

func TestGitRepositoriesDataSource_Read_TestHandleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	repoClient := azdosdkmocks.NewMockGitClient(ctrl)

	clients := &config.AggregatedClient{
		GitReposClient: repoClient,
		GraphClient:    graphClient,
		Ctx:            context.Background(),
	}

	expectedGetRepositoriesArgs := git.GetRepositoriesArgs{
		IncludeHidden: converter.Bool(false),
	}

	repoClient.
		EXPECT().
		GetRepositories(clients.Ctx, expectedGetRepositoriesArgs).
		Return(nil, errors.New("GetRepositories() Failed")).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, dataGitRepositories().Schema, nil)

	err := dataSourceGitRepositoriesRead(resourceData, clients)
	require.NotNil(t, err)
	require.Zero(t, resourceData.Id())
	repos := resourceData.Get("repositories").(*schema.Set)
	require.NotNil(t, repos)
	require.Zero(t, repos.Len())
}

func TestGitRepositoriesDataSource_Read_TestHandleErrorWithSpecificRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	repoClient := azdosdkmocks.NewMockGitClient(ctrl)

	clients := &config.AggregatedClient{
		GitReposClient: repoClient,
		GraphClient:    graphClient,
		Ctx:            context.Background(),
	}

	repo := gitRepoList[2]
	expectedGetRepositoryArgs := git.GetRepositoryArgs{
		RepositoryId: repo.Name,
		Project:      converter.String(repo.Project.Id.String()),
	}
	repoClient.
		EXPECT().
		GetRepository(clients.Ctx, expectedGetRepositoryArgs).
		Return(nil, errors.New("GetRepository() Failed"))

	resourceData := schema.TestResourceDataRaw(t, dataGitRepositories().Schema, nil)
	resourceData.Set("name", *repo.Name)
	resourceData.Set("project_id", repo.Project.Id.String())

	err := dataSourceGitRepositoriesRead(resourceData, clients)
	require.NotNil(t, err)
	require.Zero(t, resourceData.Id())
	repos := resourceData.Get("repositories").(*schema.Set)
	require.NotNil(t, repos)
	require.Zero(t, repos.Len())
}

func TestGitRepositoriesDataSource_Read_NoRepositories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	repoClient := azdosdkmocks.NewMockGitClient(ctrl)

	clients := &config.AggregatedClient{
		GitReposClient: repoClient,
		GraphClient:    graphClient,
		Ctx:            context.Background(),
	}

	expectedGetRepositoriesArgs := git.GetRepositoriesArgs{
		IncludeHidden: converter.Bool(false),
	}

	repoClient.
		EXPECT().
		GetRepositories(clients.Ctx, expectedGetRepositoriesArgs).
		Return(&[]git.GitRepository{}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, dataGitRepositories().Schema, nil)

	err := dataSourceGitRepositoriesRead(resourceData, clients)
	require.Nil(t, err)
	repos := resourceData.Get("repositories").(*schema.Set)
	require.NotNil(t, repos)
	require.Zero(t, repos.Len())
}

func TestGitRepositoriesDataSource_Read_AllRepositories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	repoClient := azdosdkmocks.NewMockGitClient(ctrl)

	clients := &config.AggregatedClient{
		GitReposClient: repoClient,
		GraphClient:    graphClient,
		Ctx:            context.Background(),
	}

	expectedGetRepositoriesArgs := git.GetRepositoriesArgs{
		IncludeHidden: converter.Bool(false),
	}

	repoClient.
		EXPECT().
		GetRepositories(clients.Ctx, expectedGetRepositoriesArgs).
		Return(&gitRepoList, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, dataGitRepositories().Schema, nil)

	err := dataSourceGitRepositoriesRead(resourceData, clients)
	require.Nil(t, err)
	repos := resourceData.Get("repositories").(*schema.Set)
	require.NotNil(t, repos)
	require.Equal(t, repos.Len(), 3)
}

func TestGitRepositoriesDataSource_Read_AllRepositoriesByProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	repoClient := azdosdkmocks.NewMockGitClient(ctrl)

	clients := &config.AggregatedClient{
		GitReposClient: repoClient,
		GraphClient:    graphClient,
		Ctx:            context.Background(),
	}

	expectedGetRepositoriesArgs := git.GetRepositoriesArgs{
		Project:       converter.String(azProjectRef.Id.String()),
		IncludeHidden: converter.Bool(false),
	}

	repoClient.
		EXPECT().
		GetRepositories(clients.Ctx, expectedGetRepositoriesArgs).
		Return(&[]git.GitRepository{
			gitRepoList[0],
			gitRepoList[1],
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, dataGitRepositories().Schema, nil)
	resourceData.Set("project_id", azProjectRef.Id.String())

	err := dataSourceGitRepositoriesRead(resourceData, clients)
	require.Nil(t, err)
	repos := resourceData.Get("repositories").(*schema.Set)
	require.NotNil(t, repos)
	require.Equal(t, repos.Len(), 2)
	repoMap := make(map[string]interface{})
	for _, item := range repos.List() {
		repoData := item.(map[string]interface{})
		repoMap[repoData["name"].(string)] = repoData
	}

	for i := 0; i < 2; i++ {
		require.Contains(t, repoMap, *gitRepoList[i].Name)
		repo := repoMap[*gitRepoList[i].Name].(map[string]interface{})
		require.Equal(t, gitRepoList[i].Project.Id.String(), repo["project_id"])
		require.Equal(t, gitRepoList[i].Id.String(), repo["id"])
	}
}

func TestGitRepositoriesDataSource_Read_SingleRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	repoClient := azdosdkmocks.NewMockGitClient(ctrl)

	clients := &config.AggregatedClient{
		GitReposClient: repoClient,
		GraphClient:    graphClient,
		Ctx:            context.Background(),
	}

	repo := gitRepoList[1]
	expectedGetRepositoryArgs := git.GetRepositoryArgs{
		RepositoryId: repo.Name,
		Project:      converter.String(repo.Project.Id.String()),
	}
	repoClient.
		EXPECT().
		GetRepository(clients.Ctx, expectedGetRepositoryArgs).
		Return(&repo, nil)

	resourceData := schema.TestResourceDataRaw(t, dataGitRepositories().Schema, nil)
	resourceData.Set("name", *repo.Name)
	resourceData.Set("project_id", repo.Project.Id.String())

	err := dataSourceGitRepositoriesRead(resourceData, clients)
	require.Nil(t, err)
	repos := resourceData.Get("repositories").(*schema.Set)
	require.NotNil(t, repos)
	require.Equal(t, repos.Len(), 1)
}

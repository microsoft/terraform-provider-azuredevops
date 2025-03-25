//go:build (all || git || data_sources || data_git_repository_file) && (!exclude_data_sources || !exclude_git || !exclude_data_git_repository_file)
// +build all git data_sources data_git_repository_file
// +build !exclude_data_sources !exclude_git !exclude_data_git_repository_file

package git

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/testhelper"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var gitFileRepo = git.GitRepository{
	Links:         nil,
	DefaultBranch: converter.String("master"),
	Id:            testhelper.CreateUUID(),
	IsFork:        converter.Bool(true),
	Name:          converter.String("repo-02"),
	ParentRepository: &git.GitRepositoryRef{
		Id:   testhelper.CreateUUID(),
		Name: converter.String("repo-parent-02"),
	},
}

var gitItem = git.GitItem{
	Path:     converter.String("MyFile.txt"),
	CommitId: converter.String("ca82a6dff817ec66f44342007202690a93763949"),
	Content:  converter.String("hello-world"),
}

var gitCommit = git.GitCommit{
	CommitId:         converter.String("ca82a6dff817ec66f44342007202690a93763949"),
	Comment:          converter.String("Commit message"),
	CommentTruncated: converter.Bool(false),
}

func TestGitRepositoryFileDataSource_ReadBranch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoClient := azdosdkmocks.NewMockGitClient(ctrl)

	clients := &client.AggregatedClient{
		GitReposClient: repoClient,
		Ctx:            context.Background(),
	}

	repoClient.
		EXPECT().
		GetCommit(clients.Ctx, git.GetCommitArgs{
			RepositoryId: converter.String(gitFileRepo.Id.String()),
			CommitId:     gitItem.CommitId,
		}).
		Return(&gitCommit, nil).
		After(
			repoClient.
				EXPECT().
				GetItem(clients.Ctx, git.GetItemArgs{
					RepositoryId:   converter.String(gitFileRepo.Id.String()),
					Path:           gitItem.Path,
					IncludeContent: converter.Bool(true),
					VersionDescriptor: &git.GitVersionDescriptor{
						Version:     converter.String("master"),
						VersionType: &git.GitVersionTypeValues.Branch,
					},
				}).
				Return(&gitItem, nil))

	resourceData := schema.TestResourceDataRaw(t, DataGitRepositoryFile().Schema, nil)
	resourceData.Set("repository_id", gitFileRepo.Id.String())
	resourceData.Set("file", gitItem.Path)
	resourceData.Set("branch", converter.String("master")) // Uses short branch ref

	err := dataSourceGitRepositoryFileRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, gitFileRepo.Id.String()+"/"+*gitItem.Path+":branch:master", resourceData.Id())
	require.Equal(t, *gitItem.Content, resourceData.Get("content"))
	require.Equal(t, *gitCommit.Comment, resourceData.Get("commit_message"))
}

func TestGitRepositoryFileDataSource_ReadTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoClient := azdosdkmocks.NewMockGitClient(ctrl)

	clients := &client.AggregatedClient{
		GitReposClient: repoClient,
		Ctx:            context.Background(),
	}

	repoClient.
		EXPECT().
		GetCommit(clients.Ctx, git.GetCommitArgs{
			RepositoryId: converter.String(gitFileRepo.Id.String()),
			CommitId:     gitItem.CommitId,
		}).
		Return(&gitCommit, nil).
		After(
			repoClient.
				EXPECT().
				GetItem(clients.Ctx, git.GetItemArgs{
					RepositoryId:   converter.String(gitFileRepo.Id.String()),
					Path:           gitItem.Path,
					IncludeContent: converter.Bool(true),
					VersionDescriptor: &git.GitVersionDescriptor{
						Version:     converter.String("v1.2.3"),
						VersionType: &git.GitVersionTypeValues.Tag,
					},
				}).
				Return(&gitItem, nil))

	resourceData := schema.TestResourceDataRaw(t, DataGitRepositoryFile().Schema, nil)
	resourceData.Set("repository_id", gitFileRepo.Id.String())
	resourceData.Set("file", gitItem.Path)
	resourceData.Set("tag", converter.String("refs/tags/v1.2.3")) // Uses full tag ref to validate tag splitting

	err := dataSourceGitRepositoryFileRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, gitFileRepo.Id.String()+"/"+*gitItem.Path+":tag:v1.2.3", resourceData.Id())
	require.Equal(t, *gitItem.Content, resourceData.Get("content"))
	require.Equal(t, *gitCommit.Comment, resourceData.Get("commit_message"))
}

func TestGitRepositoryFileDataSource_ReadCommitFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoClient := azdosdkmocks.NewMockGitClient(ctrl)

	clients := &client.AggregatedClient{
		GitReposClient: repoClient,
		Ctx:            context.Background(),
	}

	repoClient.
		EXPECT().
		GetCommit(clients.Ctx, git.GetCommitArgs{
			RepositoryId: converter.String(gitFileRepo.Id.String()),
			CommitId:     gitItem.CommitId,
		}).
		Return(nil, fmt.Errorf("Failed to read commit")).
		After(
			repoClient.
				EXPECT().
				GetItem(clients.Ctx, git.GetItemArgs{
					RepositoryId:   converter.String(gitFileRepo.Id.String()),
					Path:           gitItem.Path,
					IncludeContent: converter.Bool(true),
					VersionDescriptor: &git.GitVersionDescriptor{
						Version:     converter.String("master"),
						VersionType: &git.GitVersionTypeValues.Branch,
					},
				}).
				Return(&gitItem, nil))

	resourceData := schema.TestResourceDataRaw(t, DataGitRepositoryFile().Schema, nil)
	resourceData.Set("repository_id", gitFileRepo.Id.String())
	resourceData.Set("file", gitItem.Path)
	resourceData.Set("branch", converter.String("master")) // Uses short branch ref

	err := dataSourceGitRepositoryFileRead(resourceData, clients)
	require.NotNil(t, err)
	require.Equal(t, "Get commit failed, repositoryID: "+gitFileRepo.Id.String()+", branch: master, file: "+*gitItem.Path+". Error:  Failed to read commit", err.Error())
}

func TestGitRepositoryFileDataSource_ReadItemFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoClient := azdosdkmocks.NewMockGitClient(ctrl)

	clients := &client.AggregatedClient{
		GitReposClient: repoClient,
		Ctx:            context.Background(),
	}

	repoClient.
		EXPECT().
		GetItem(clients.Ctx, git.GetItemArgs{
			RepositoryId:   converter.String(gitFileRepo.Id.String()),
			Path:           gitItem.Path,
			IncludeContent: converter.Bool(true),
			VersionDescriptor: &git.GitVersionDescriptor{
				Version:     converter.String("master"),
				VersionType: &git.GitVersionTypeValues.Branch,
			},
		}).
		Return(nil, fmt.Errorf("Failed to get item"))

	resourceData := schema.TestResourceDataRaw(t, DataGitRepositoryFile().Schema, nil)
	resourceData.Set("repository_id", gitFileRepo.Id.String())
	resourceData.Set("file", gitItem.Path)
	resourceData.Set("branch", converter.String("master")) // Uses short branch ref

	err := dataSourceGitRepositoryFileRead(resourceData, clients)
	require.NotNil(t, err)
	require.Equal(t, "Get item failed, repositoryID: "+gitFileRepo.Id.String()+", branch: master, file: "+*gitItem.Path+". Error: Failed to get item", err.Error())
}

func TestGitRepositoryFileDataSource_ReadItemNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoClient := azdosdkmocks.NewMockGitClient(ctrl)

	clients := &client.AggregatedClient{
		GitReposClient: repoClient,
		Ctx:            context.Background(),
	}

	repoClient.
		EXPECT().
		GetItem(clients.Ctx, git.GetItemArgs{
			RepositoryId:   converter.String(gitFileRepo.Id.String()),
			Path:           gitItem.Path,
			IncludeContent: converter.Bool(true),
			VersionDescriptor: &git.GitVersionDescriptor{
				Version:     converter.String("master"),
				VersionType: &git.GitVersionTypeValues.Branch,
			},
		}).
		Return(nil, azuredevops.WrappedError{
			StatusCode: converter.Int(http.StatusNotFound),
		})

	resourceData := schema.TestResourceDataRaw(t, DataGitRepositoryFile().Schema, nil)
	resourceData.Set("repository_id", gitFileRepo.Id.String())
	resourceData.Set("file", gitItem.Path)
	resourceData.Set("branch", converter.String("master")) // Uses short branch ref

	err := dataSourceGitRepositoryFileRead(resourceData, clients)
	require.NotNil(t, err)
	require.Equal(t, "Item not found, repositoryID: "+gitFileRepo.Id.String()+", branch: master, file: "+*gitItem.Path+". Error: REST call returned status code 404", err.Error())
}

func TestGitRepositoryFileDataSource_NoVersionType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		Ctx: context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, DataGitRepositoryFile().Schema, nil)
	resourceData.Set("repository_id", gitFileRepo.Id.String())
	resourceData.Set("file", gitItem.Path)

	err := dataSourceGitRepositoryFileRead(resourceData, clients)
	require.NotNil(t, err)
	require.Equal(t, "One of 'branch' or 'tag' must be specified", err.Error())
}

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
	require.Equal(t, fmt.Sprintf("%s/%s:branch:master", gitFileRepo.Id.String(), *gitItem.Path), resourceData.Id())
	require.Equal(t, *gitItem.Content, resourceData.Get("content"))
	require.Equal(t, *gitCommit.Comment, resourceData.Get("last_commit_message"))
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
	require.Equal(t, fmt.Sprintf("%s/%s:tag:v1.2.3", gitFileRepo.Id.String(), *gitItem.Path), resourceData.Id())
	require.Equal(t, *gitItem.Content, resourceData.Get("content"))
	require.Equal(t, *gitCommit.Comment, resourceData.Get("last_commit_message"))
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
	require.Equal(t, fmt.Sprintf("Get commit failed, repositoryID: %s, commitID: %s. Error:  Failed to read commit", gitFileRepo.Id.String(), *gitItem.CommitId), err.Error())
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
	require.Equal(t, fmt.Sprintf("Get item failed, repositoryID: %s, branch: master, file: %s. Error: Failed to get item", gitFileRepo.Id.String(), *gitItem.Path), err.Error())
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
	require.Equal(t, fmt.Sprintf("Item not found, repositoryID: %s, branch: master, file: %s. Error: REST call returned status code 404", gitFileRepo.Id.String(), *gitItem.Path), err.Error())
}

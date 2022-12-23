//go:build (all || git || resource_git_repository) && (!exclude_git || !exclude_resource_git_repository)
// +build all git resource_git_repository
// +build !exclude_git !exclude_resource_git_repository

package git

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/git"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

var testRepoProjectID = uuid.New()
var testRepoID = uuid.New()

// This definition matches the overall structure of what a configured git repository would
// look like. Note that the ID and Name attributes match -- this is the service-side behavior
// when configuring a GitHub repo.
var testGitRepository = git.GitRepository{
	Id:   &testRepoID,
	Name: converter.String("RepoName"),
	Project: &core.TeamProjectReference{
		Id:   &testRepoProjectID,
		Name: converter.String("ProjectName"),
	},
}

// verifies that the create operation is considered failed if the initial API
// call fails.
func TestGitRepo_Create_DoesNotSwallowErrorFromFailedCreateCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := schema.TestResourceDataRaw(t, ResourceGitRepository().Schema, nil)
	resourceData.SetId(testGitRepository.Id.String())
	flattenGitRepository(resourceData, &testGitRepository)
	configureCleanInitialization(resourceData)

	reposClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &client.AggregatedClient{GitReposClient: reposClient, Ctx: context.Background()}

	expectedArgs := git.CreateRepositoryArgs{
		GitRepositoryToCreate: &git.GitRepositoryCreateOptions{
			Name: testGitRepository.Name,
			Project: &core.TeamProjectReference{
				Id: &testRepoProjectID,
			},
		},
	}
	reposClient.
		EXPECT().
		CreateRepository(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateGitRepository() Failed")).
		Times(1)

	err := resourceGitRepositoryCreate(resourceData, clients)
	require.Regexp(t, ".*CreateGitRepository\\(\\) Failed$", err.Error())
}

// verifies that the update operation is considered failed if the initial API
// call fails.
func TestGitRepo_Update_DoesNotSwallowErrorFromFailedCreateCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := schema.TestResourceDataRaw(t, ResourceGitRepository().Schema, nil)
	resourceData.SetId(testGitRepository.Id.String())
	flattenGitRepository(resourceData, &testGitRepository)
	configureCleanInitialization(resourceData)

	reposClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &client.AggregatedClient{GitReposClient: reposClient, Ctx: context.Background()}

	reposClient.
		EXPECT().
		UpdateRepository(clients.Ctx, gomock.Any()).
		Return(nil, errors.New("UpdateGitRepository() Failed")).
		Times(1)

	err := resourceGitRepositoryUpdate(resourceData, clients)
	require.Regexp(t, ".*UpdateGitRepository\\(\\) Failed$", err.Error())
}

func configureCleanInitialization(d *schema.ResourceData) {
	d.Set("initialization", &[]map[string]interface{}{
		{
			"init_type": "Clean",
		},
	})
}

// verifies that a round-trip flatten/expand sequence will not result in data loss of non-computed properties.
//
//	Note: there is no need to expand computed properties, so they won't be tested here.
func TestGitRepo_FlattenExpand_RoundTrip(t *testing.T) {
	projectID := uuid.New()
	project := core.TeamProjectReference{Id: &projectID}

	repoID := uuid.New()
	repoName := "name"
	gitRepo := git.GitRepository{Id: &repoID, Name: &repoName, Project: &project}

	resourceData := schema.TestResourceDataRaw(t, ResourceGitRepository().Schema, nil)
	resourceData.SetId(gitRepo.Id.String())
	configureCleanInitialization(resourceData)
	flattenGitRepository(resourceData, &gitRepo)

	expandedGitRepo, repoInitialization, expandedProjectID, err := expandGitRepository(resourceData)

	require.Nil(t, err)
	require.NotNil(t, expandedGitRepo)
	require.NotNil(t, expandedGitRepo.Id)
	require.Equal(t, *expandedGitRepo.Id, repoID)
	require.NotNil(t, expandedProjectID)
	require.Equal(t, *expandedProjectID, projectID)
	require.NotNil(t, repoInitialization)
	require.Equal(t, repoInitialization.initType, "Clean")
	require.Equal(t, repoInitialization.sourceType, "")
	require.Equal(t, repoInitialization.sourceURL, "")
}

// verifies that the read operation is considered failed if the initial API
// call fails.
func TestGitRepo_Read_DoesNotSwallowErrorFromFailedReadCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reposClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &client.AggregatedClient{
		GitReposClient: reposClient,
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, ResourceGitRepository().Schema, nil)
	resourceData.SetId("an-id")
	resourceData.Set("project_id", "a-project")

	expectedArgs := git.GetRepositoryArgs{RepositoryId: converter.String("an-id"), Project: converter.String("a-project")}
	reposClient.
		EXPECT().
		GetRepository(clients.Ctx, expectedArgs).
		Return(nil, fmt.Errorf("GetRepository() Failed")).
		Times(1)

	err := resourceGitRepositoryRead(resourceData, clients)
	require.Contains(t, err.Error(), "GetRepository() Failed")
}

// verifies that 'Clean' repo initalization uses default branch name
func TestGitRepo_Initialize_UsesTheDefaultBranch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reposClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &client.AggregatedClient{
		GitReposClient: reposClient,
		Ctx:            context.Background(),
	}

	repoName := "test-repository"
	projectName := "test-project"

	repo := git.GitRepository{
		Name: &repoName,
		Project: &core.TeamProjectReference{
			Name: &projectName,
		},
	}

	defaultBranch := "refs/heads/main"

	expectedArgs := git.CreatePushArgs{
		RepositoryId: repo.Name,
		Project:      repo.Project.Name,
		Push: &git.GitPush{
			RefUpdates: &[]git.GitRefUpdate{
				{
					Name:        converter.String(defaultBranch),
					OldObjectId: converter.String("0000000000000000000000000000000000000000"),
				},
			},
			Commits: &[]git.GitCommitRef{
				{
					Comment: converter.String("Initial commit."),
					Changes: &[]interface{}{
						git.Change{
							ChangeType: &git.VersionControlChangeTypeValues.Add,
							Item: git.GitItem{
								Path: converter.String("/readme.md"),
							},
							NewContent: &git.ItemContent{
								ContentType: &git.ItemContentTypeValues.RawText,
								Content:     repo.Project.Name,
							},
						},
					},
				},
			},
		},
	}
	reposClient.
		EXPECT().
		CreatePush(clients.Ctx, expectedArgs).
		Return(nil, nil).
		Times(1)

	err := initializeGitRepository(clients, &repo, &defaultBranch)
	require.Nil(t, err, fmt.Sprintf("Error was %v", err))
}

// verifies that the resource ID is used for reads if the ID is set
func TestGitRepo_Read_UsesIdIfSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reposClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &client.AggregatedClient{
		GitReposClient: reposClient,
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, ResourceGitRepository().Schema, nil)
	resourceData.SetId("an-id")
	resourceData.Set("project_id", "a-project")

	expectedArgs := git.GetRepositoryArgs{RepositoryId: converter.String("an-id"), Project: converter.String("a-project")}
	reposClient.
		EXPECT().
		GetRepository(clients.Ctx, expectedArgs).
		Return(nil, fmt.Errorf("error")).
		Times(1)

	resourceGitRepositoryRead(resourceData, clients)
}

func TestGitRepo_Delete_ChecksForValidUUID(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceGitRepository().Schema, nil)
	resourceData.SetId("not-a-uuid-id")

	err := resourceGitRepositoryDelete(resourceData, &client.AggregatedClient{})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "Invalid repositoryId UUID")
}

func TestGitRepo_Delete_DoesNotSwallowErrorFromFailedDeleteCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reposClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &client.AggregatedClient{
		GitReposClient: reposClient,
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, ResourceGitRepository().Schema, nil)
	id := uuid.New()
	resourceData.SetId(id.String())

	expectedArgs := git.DeleteRepositoryArgs{RepositoryId: &id}
	reposClient.
		EXPECT().
		DeleteRepository(clients.Ctx, expectedArgs).
		Return(fmt.Errorf("DeleteRepository() Failed")).
		Times(1)

	err := resourceGitRepositoryDelete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteRepository() Failed")
}

// verifies that the name is used for reads if the ID is not set
func TestGitRepo_Read_UsesNameIfIdNotSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reposClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &client.AggregatedClient{
		GitReposClient: reposClient,
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, ResourceGitRepository().Schema, nil)
	resourceData.Set("name", "a-name")
	resourceData.Set("project_id", "a-project")

	expectedArgs := git.GetRepositoryArgs{RepositoryId: converter.String("a-name"), Project: converter.String("a-project")}
	reposClient.
		EXPECT().
		GetRepository(clients.Ctx, expectedArgs).
		Return(nil, fmt.Errorf("error")).
		Times(1)

	resourceGitRepositoryRead(resourceData, clients)
}

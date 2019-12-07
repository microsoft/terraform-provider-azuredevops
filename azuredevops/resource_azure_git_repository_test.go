// +build all core resource_git_repository

package azuredevops

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/testhelper"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/stretchr/testify/require"
)

var testRepoProjectID = uuid.New()
var testRepoID = uuid.New()

// This definition matches the overall structure of what a configured git repository would
// look like. Note that the ID and Name attributes match -- this is the service-side behavior
// when configuring a GitHub repo.
var testAzureGitRepository = git.GitRepository{
	Id:   &testRepoID,
	Name: converter.String("RepoName"),
	Project: &core.TeamProjectReference{
		Id:   &testRepoProjectID,
		Name: converter.String("ProjectName"),
	},
}

/**
 * Begin unit tests
 */

// verifies that the create operation is considered failed if the initial API
// call fails.
func TestAzureGitRepo_Create_DoesNotSwallowErrorFromFailedCreateCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := schema.TestResourceDataRaw(t, resourceAzureGitRepository().Schema, nil)
	flattenAzureGitRepository(resourceData, &testAzureGitRepository)
	configureCleanInitialization(resourceData)

	reposClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &config.AggregatedClient{GitReposClient: reposClient, Ctx: context.Background()}

	expectedArgs := git.CreateRepositoryArgs{
		GitRepositoryToCreate: &git.GitRepositoryCreateOptions{
			Name: testAzureGitRepository.Name,
			Project: &core.TeamProjectReference{
				Id: &testRepoProjectID,
			},
		},
	}
	reposClient.
		EXPECT().
		CreateRepository(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateAzureGitRepository() Failed")).
		Times(1)

	err := resourceAzureGitRepositoryCreate(resourceData, clients)
	require.Regexp(t, ".*CreateAzureGitRepository\\(\\) Failed$", err.Error())
}

// verifies that the update operation is considered failed if the initial API
// call fails.
func TestAzureGitRepo_Update_DoesNotSwallowErrorFromFailedCreateCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := schema.TestResourceDataRaw(t, resourceAzureGitRepository().Schema, nil)
	flattenAzureGitRepository(resourceData, &testAzureGitRepository)
	configureCleanInitialization(resourceData)

	reposClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &config.AggregatedClient{GitReposClient: reposClient, Ctx: context.Background()}

	reposClient.
		EXPECT().
		UpdateRepository(clients.Ctx, gomock.Any()).
		Return(nil, errors.New("UpdateAzureGitRepository() Failed")).
		Times(1)

	err := resourceAzureGitRepositoryUpdate(resourceData, clients)
	require.Regexp(t, ".*UpdateAzureGitRepository\\(\\) Failed$", err.Error())
}

func configureCleanInitialization(d *schema.ResourceData) {
	d.Set("initialization", &[]map[string]interface{}{
		{
			"init_type": "Clean",
		},
	})
}

// verifies that a round-trip flatten/expand sequence will not result in data loss of non-computed properties.
//	Note: there is no need to expand computed properties, so they won't be tested here.
func TestAzureGitRepo_FlattenExpand_RoundTrip(t *testing.T) {
	projectID := uuid.New()
	project := core.TeamProjectReference{Id: &projectID}

	repoID := uuid.New()
	repoName := "name"
	gitRepo := git.GitRepository{Id: &repoID, Name: &repoName, Project: &project}

	resourceData := schema.TestResourceDataRaw(t, resourceAzureGitRepository().Schema, nil)
	flattenAzureGitRepository(resourceData, &gitRepo)
	configureCleanInitialization(resourceData)

	expandedGitRepo, repoInitialization, expandedProjectID, err := expandAzureGitRepository(resourceData)

	require.Nil(t, err)
	require.Equal(t, *expandedGitRepo.Id, repoID)
	require.Equal(t, *expandedProjectID, projectID)
	require.Equal(t, repoInitialization.initType, "Clean")
	require.Equal(t, repoInitialization.sourceType, "")
	require.Equal(t, repoInitialization.sourceURL, "")
}

// verifies that the read operation is considered failed if the initial API
// call fails.
func TestAzureGitRepo_Read_DoesNotSwallowErrorFromFailedReadCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reposClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &config.AggregatedClient{
		GitReposClient: reposClient,
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, resourceAzureGitRepository().Schema, nil)
	resourceData.SetId("an-id")
	resourceData.Set("project_id", "a-project")

	expectedArgs := git.GetRepositoryArgs{RepositoryId: converter.String("an-id"), Project: converter.String("a-project")}
	reposClient.
		EXPECT().
		GetRepository(clients.Ctx, expectedArgs).
		Return(nil, fmt.Errorf("GetRepository() Failed")).
		Times(1)

	err := resourceAzureGitRepositoryRead(resourceData, clients)
	require.Contains(t, err.Error(), "GetRepository() Failed")
}

// verifies that the resource ID is used for reads if the ID is set
func TestAzureGitRepo_Read_UsesIdIfSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reposClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &config.AggregatedClient{
		GitReposClient: reposClient,
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, resourceAzureGitRepository().Schema, nil)
	resourceData.SetId("an-id")
	resourceData.Set("project_id", "a-project")

	expectedArgs := git.GetRepositoryArgs{RepositoryId: converter.String("an-id"), Project: converter.String("a-project")}
	reposClient.
		EXPECT().
		GetRepository(clients.Ctx, expectedArgs).
		Return(nil, fmt.Errorf("error")).
		Times(1)

	resourceAzureGitRepositoryRead(resourceData, clients)
}

func TestAzureGitRepo_Delete_ChecksForValidUUID(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceAzureGitRepository().Schema, nil)
	resourceData.SetId("not-a-uuid-id")

	err := resourceAzureGitRepositoryDelete(resourceData, &config.AggregatedClient{})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "Invalid repositoryId UUID")
}

func TestAzureGitRepo_Delete_DoesNotSwallowErrorFromFailedDeleteCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reposClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &config.AggregatedClient{
		GitReposClient: reposClient,
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, resourceAzureGitRepository().Schema, nil)
	id := uuid.New()
	resourceData.SetId(id.String())

	expectedArgs := git.DeleteRepositoryArgs{RepositoryId: &id}
	reposClient.
		EXPECT().
		DeleteRepository(clients.Ctx, expectedArgs).
		Return(fmt.Errorf("DeleteRepository() Failed")).
		Times(1)

	err := resourceAzureGitRepositoryDelete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteRepository() Failed")
}

// verifies that the name is used for reads if the ID is not set
func TestAzureGitRepo_Read_UsesNameIfIdNotSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reposClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &config.AggregatedClient{
		GitReposClient: reposClient,
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, resourceAzureGitRepository().Schema, nil)
	resourceData.Set("name", "a-name")
	resourceData.Set("project_id", "a-project")

	expectedArgs := git.GetRepositoryArgs{RepositoryId: converter.String("a-name"), Project: converter.String("a-project")}
	reposClient.
		EXPECT().
		GetRepository(clients.Ctx, expectedArgs).
		Return(nil, fmt.Errorf("error")).
		Times(1)

	resourceAzureGitRepositoryRead(resourceData, clients)
}

/**
 * Begin acceptance tests
 */

// Verifies that the following sequence of events occurrs without error:
//	(1) TF apply creates resource
//	(2) TF state values are set
//	(3) resource can be queried by ID and has expected name
// 	(4) TF destroy deletes resource
//	(5) resource can no longer be queried by ID
func TestAccAzureGitRepo_CreateAndUpdate(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	gitRepoNameFirst := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	gitRepoNameSecond := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	tfRepoNode := "azuredevops_azure_git_repository.gitrepo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testhelper.TestAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccAzureGitRepoCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testhelper.TestAccAzureGitRepoResource(projectName, gitRepoNameFirst, "Uninitialized"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfRepoNode, "name", gitRepoNameFirst),
					testAccCheckAzureGitRepoResourceExists(gitRepoNameFirst),
					resource.TestCheckResourceAttrSet(tfRepoNode, "is_fork"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "remote_url"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "size"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "ssh_url"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "url"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "web_url"),
				),
			},
			{
				Config: testhelper.TestAccAzureGitRepoResource(projectName, gitRepoNameSecond, "Uninitialized"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfRepoNode, "name", gitRepoNameSecond),
					testAccCheckAzureGitRepoResourceExists(gitRepoNameSecond),
					resource.TestCheckResourceAttrSet(tfRepoNode, "is_fork"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "remote_url"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "size"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "ssh_url"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "url"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "web_url"),
				),
			},
		},
	})
}

// Given the name of an AzDO git repository, this will return a function that will check whether
// or not the definition (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func testAccCheckAzureGitRepoResourceExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clients := testAccProvider.Meta().(*config.AggregatedClient)

		gitRepo, ok := s.RootModule().Resources["azuredevops_azure_git_repository.gitrepo"]
		if !ok {
			return fmt.Errorf("Did not find a repo definition in the TF state")
		}

		repoID := gitRepo.Primary.ID
		projectID := gitRepo.Primary.Attributes["project_id"]

		repo, err := azureGitRepositoryRead(clients, repoID, "", projectID)
		if err != nil {
			return err
		}

		if *repo.Name != expectedName {
			return fmt.Errorf("AzDO Git Repository has Name=%s, but expected Name=%s", *repo.Name, expectedName)
		}

		return nil
	}
}

func testAccAzureGitRepoCheckDestroy(s *terraform.State) error {
	clients := testAccProvider.Meta().(*config.AggregatedClient)

	// verify that every repository referenced in the state does not exist in AzDO
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_azure_git_repository" {
			continue
		}

		repoID := resource.Primary.ID
		projectID := resource.Primary.Attributes["project_id"]

		// indicates the git repository still exists - this should fail the test
		if _, err := azureGitRepositoryRead(clients, repoID, "", projectID); err == nil {
			return fmt.Errorf("repository with ID %s should not exist", repoID)
		}
	}

	return nil
}

// Verifies that a newly created repo with init_type of "Clean" has the expected
// master branch available
func TestAccAzureGitRepo_RepoInitialization_Clean(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	gitRepoName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	tfRepoNode := "azuredevops_azure_git_repository.gitrepo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testhelper.TestAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccAzureGitRepoCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testhelper.TestAccAzureGitRepoResource(projectName, gitRepoName, "Clean"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfRepoNode, "name", gitRepoName),
					testAccCheckAzureGitRepoResourceExists(gitRepoName),
					resource.TestCheckResourceAttr(tfRepoNode, "default_branch", "refs/heads/master"),
				),
			},
		},
	})
}

// Verifies that a newly created repo with init_type of "Uninitialized" does NOT
// have a master branch established
func TestAccAzureGitRepo_RepoInitialization_Uninitialized(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	gitRepoName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	tfRepoNode := "azuredevops_azure_git_repository.gitrepo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testhelper.TestAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccAzureGitRepoCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testhelper.TestAccAzureGitRepoResource(projectName, gitRepoName, "Uninitialized"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfRepoNode, "name", gitRepoName),
					testAccCheckAzureGitRepoResourceExists(gitRepoName),
					resource.TestCheckResourceAttr(tfRepoNode, "default_branch", ""),
				),
			},
		},
	})
}

func init() {
	InitProvider()
}

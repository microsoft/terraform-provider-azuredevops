// +build all permissions resource_git_permissions
// +build !exclude_permissions !exclude_resource_project_permissions

package permissions

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
)

/**
 * Begin unit tests
 */

var gitProjectID = "9083e944-8e9e-405e-960a-c80180aa71e6"
var gitTokenProject = fmt.Sprintf("repoV2/%s", gitProjectID)
var gitRepositoryID = "c629a0a4-926d-45d1-8095-6e2499cf3938"
var gitTokenRepository = fmt.Sprintf("%s/%s", gitTokenProject, gitRepositoryID)
var gitTokenBranchAll = fmt.Sprintf("%s/refs/heads", gitTokenRepository)
var gitBranchNameValid = "master"
var gitTokenBranch = fmt.Sprintf("%s/refs/heads/%s", gitTokenRepository, encodeBranchName(gitBranchNameValid))
var gitBranchNameInValid = "@@invalid@@"

func TestGitPermissions_CreateGitToken(t *testing.T) {
	var d *schema.ResourceData
	var token string
	var err error

	d = getGitPermissionsResource(t, gitProjectID, "", "")
	token, err = createGitToken(d, nil)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, gitTokenProject, token)

	d = getGitPermissionsResource(t, gitProjectID, gitRepositoryID, "")
	token, err = createGitToken(d, nil)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, gitTokenRepository, token)

	d = getGitPermissionsResource(t, "", gitRepositoryID, "")
	token, err = createGitToken(d, nil)
	assert.Empty(t, token)
	assert.NotNil(t, err)
}

func TestGitPermissions_CreateGitTokenWithBranch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &client.AggregatedClient{
		GitReposClient: gitClient,
		Ctx:            context.Background(),
	}

	gitClient.EXPECT().
		GetRefs(clients.Ctx, gomock.Any()).
		Return(&git.GetRefsResponseValue{
			Value: []git.GitRef{
				{
					Name: &gitBranchNameValid,
				},
			},
			ContinuationToken: "",
		}, nil).
		Times(1)

	var d *schema.ResourceData
	var token string
	var err error

	d = getGitPermissionsResource(t, gitProjectID, "", gitBranchNameValid)
	token, err = createGitToken(d, clients)
	assert.Empty(t, token)
	assert.NotNil(t, err)

	d = getGitPermissionsResource(t, gitProjectID, gitRepositoryID, gitBranchNameValid)
	token, err = createGitToken(d, clients)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, gitTokenBranch, token)
}

func TestGitPermissions_CreateGitTokenWithBranch_HandleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &client.AggregatedClient{
		GitReposClient: gitClient,
		Ctx:            context.Background(),
	}

	errMsg := "@@GetRefs@@failed"
	gitClient.EXPECT().
		GetRefs(clients.Ctx, gomock.Any()).
		Return(nil, fmt.Errorf(errMsg)).
		Times(1)

	d := getGitPermissionsResource(t, gitProjectID, gitRepositoryID, gitBranchNameValid)
	token, err := createGitToken(d, clients)
	assert.Empty(t, token)
	assert.NotNil(t, err)
	assert.EqualError(t, err, errMsg)
}

func encodeBranchName(branchName string) string {
	ret, _ := converter.EncodeUtf16HexString(branchName)
	return ret
}

func getGitPermissionsResource(t *testing.T, gitProjectID string, repoID string, branchName string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, ResourceGitPermissions().Schema, nil)
	if gitProjectID != "" {
		d.Set("project_id", gitProjectID)
	}
	if repoID != "" {
		d.Set("repository_id", repoID)
	}
	if branchName != "" {
		d.Set("branch_name", branchName)
	}
	return d
}

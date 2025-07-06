//go:build (all || permissions || resource_git_permissions) && (!exclude_permissions || !exclude_resource_project_permissions)
// +build all permissions resource_git_permissions
// +build !exclude_permissions !exclude_resource_project_permissions

package permissions

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

/**
 * Begin unit tests
 */

var (
	gitProjectID          = "9083e944-8e9e-405e-960a-c80180aa71e6"
	gitTokenProject       = fmt.Sprintf("repoV2/%s", gitProjectID)
	gitRepositoryID       = "c629a0a4-926d-45d1-8095-6e2499cf3938"
	gitTokenRepository    = fmt.Sprintf("%s/%s", gitTokenProject, gitRepositoryID)
	gitTokenBranchAll     = fmt.Sprintf("%s/refs/heads", gitTokenRepository)
	gitBranchNameValid    = "master"
	gitTokenBranch        = fmt.Sprintf("%s/refs/heads/%s", gitTokenRepository, encodeBranchName(gitBranchNameValid))
	gitSubBranchNameValid = "1.0.0"
	gitTokenSubBranch     = fmt.Sprintf("%s/refs/heads/%s", gitTokenRepository, encodeBranchName(gitBranchNameValid)+"/"+encodeBranchName(gitSubBranchNameValid))
	gitBranchNameInValid  = "@@invalid@@"
)

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

	clients := &client.AggregatedClient{
		Ctx: context.Background(),
	}

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

	d = getGitPermissionsResource(t, gitProjectID, gitRepositoryID, "/refs/heads/"+gitBranchNameValid)
	token, err = createGitToken(d, clients)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, gitTokenBranch, token)

	d = getGitPermissionsResource(t, gitProjectID, gitRepositoryID, gitBranchNameValid+"/"+gitSubBranchNameValid)
	token, err = createGitToken(d, clients)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, gitTokenSubBranch, token)

	d = getGitPermissionsResource(t, gitProjectID, gitRepositoryID, "refs/heads/"+gitBranchNameValid+"/"+gitSubBranchNameValid)
	token, err = createGitToken(d, clients)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, gitTokenSubBranch, token)
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

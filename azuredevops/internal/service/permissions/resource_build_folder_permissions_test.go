//go:build (all || permissions || resource_build_folder_permissions) && (!exclude_permissions || !resource_build_folder_permissions)
// +build all permissions resource_build_folder_permissions
// +build !exclude_permissions !resource_build_folder_permissions

package permissions

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

/**
 * Begin unit tests
 */

var buildFolderProjectID = "9083e944-8e9e-405e-960a-c80180aa71e6"

var buildFolderToken = fmt.Sprintf("%s", buildFolderProjectID)

var buildFolderPath = "a/b/c"
var buildFolderTokenPath = fmt.Sprintf("%s/%s", buildFolderProjectID, buildFolderPath)

func TestBuildFolderPermissions_CreateBuildFolderToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	buildClient := azdosdkmocks.NewMockBuildClient(ctrl)
	clients := &client.AggregatedClient{
		BuildClient: buildClient,
		Ctx:         context.Background(),
	}

	folder := build.Folder{
		Description: converter.String("Test Folder"),
		Path:        converter.String("\\"),
	}

	mockFolders := []build.Folder{folder}

	buildClient.EXPECT().
		GetFolders(clients.Ctx, gomock.Any()).
		Return(&mockFolders, nil).
		Times(1)

	var d *schema.ResourceData
	var token string
	var err error

	d = getBuildFolderPermissionsResource(t, buildFolderProjectID, "\\")
	token, err = createBuildFolderToken(d, clients)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, "9083e944-8e9e-405e-960a-c80180aa71e6", token)

	d = getBuildFolderPermissionsResource(t, "", "")
	token, err = createBuildFolderToken(d, clients)
	assert.Empty(t, token)
	assert.NotNil(t, err)
}

func TestBuildFolderPermissions_CreateBuildTokenWithPaths(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	buildClient := azdosdkmocks.NewMockBuildClient(ctrl)
	clients := &client.AggregatedClient{
		BuildClient: buildClient,
		Ctx:         context.Background(),
	}

	path := "\\a\\b\\c"

	folder := build.Folder{
		Description: converter.String("Test Folder"),
		Path:        converter.String(path),
	}

	mockFolders := []build.Folder{folder}

	buildClient.EXPECT().
		GetFolders(clients.Ctx, gomock.Any()).
		Return(&mockFolders, nil).
		Times(1)

	var d *schema.ResourceData
	var token string
	var err error

	d = getBuildFolderPermissionsResource(t, buildFolderProjectID, path)
	token, err = createBuildFolderToken(d, clients)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, buildFolderTokenPath, token)
}

func getBuildFolderPermissionsResource(t *testing.T, projectID string, buildFolderPath string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, ResourceBuildFolderPermissions().Schema, nil)
	if projectID != "" {
		d.Set("project_id", projectID)
	}
	if buildFolderPath != "" {
		d.Set("path", buildFolderPath)
	}
	return d
}

// +build all permissions resource_build_definition_permissions
// +build !exclude_permissions !resource_build_definition_permissions

package permissions

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
)

/**
 * Begin unit tests
 */

var buildPermissionsID = "9083e944-8e9e-405e-960a-c80180aa71e6"
var buildDefinitionID = "5"

var buildToken = fmt.Sprintf("%s/%s", buildPermissionsID, buildDefinitionID)

var buildDefinitionPath = "a/b/c"
var buildTokenPath = fmt.Sprintf("%s/%s/%s", buildPermissionsID, buildDefinitionPath, buildDefinitionID)

func TestBuildDefinitionPermissions_CreateBuildToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	buildClient := azdosdkmocks.NewMockBuildClient(ctrl)
	clients := &client.AggregatedClient{
		BuildClient: buildClient,
		Ctx:         context.Background(),
	}

	buildClient.EXPECT().
		GetDefinition(clients.Ctx, gomock.Any()).
		Return(&build.BuildDefinition{
			Id:   converter.Int(5),
			Path: converter.String("\\"),
		}, nil).
		Times(1)

	var d *schema.ResourceData
	var token string
	var err error

	d = getBuildDefinitionPermissionsResource(t, buildPermissionsID, buildDefinitionID, "")
	token, err = createBuildToken(d, clients)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, buildToken, token)

	d = getBuildDefinitionPermissionsResource(t, "", "", "")
	token, err = createBuildToken(d, clients)
	assert.Empty(t, token)
	assert.NotNil(t, err)
}

func TestBuildDefinitionPermissions_CreateBuildTokenWithPaths(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	buildClient := azdosdkmocks.NewMockBuildClient(ctrl)
	clients := &client.AggregatedClient{
		BuildClient: buildClient,
		Ctx:         context.Background(),
	}

	path := "\\a\\b\\c"

	buildClient.EXPECT().
		GetDefinition(clients.Ctx, gomock.Any()).
		Return(&build.BuildDefinition{
			Id:   converter.Int(5),
			Path: converter.String(path),
		}, nil).
		Times(1)

	var d *schema.ResourceData
	var token string
	var err error

	d = getBuildDefinitionPermissionsResource(t, buildPermissionsID, buildDefinitionID, path)
	token, err = createBuildToken(d, clients)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, buildTokenPath, token)

	d = getBuildDefinitionPermissionsResource(t, "", "", "")
	token, err = createBuildToken(d, clients)
	assert.Empty(t, token)
	assert.NotNil(t, err)
}

func getBuildDefinitionPermissionsResource(t *testing.T, projectID string, buildDefinitionID string, buildDefinitionPath string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, ResourceBuildDefinitionPermissions().Schema, nil)
	if projectID != "" {
		d.Set("project_id", projectID)
	}
	if buildDefinitionID != "" {
		d.Set("build_definition_id", buildDefinitionID)
	}
	if buildDefinitionPath != "" {
		d.Set("build_definition_path", buildDefinitionPath)
	}
	return d
}

func TestTransformPath(t *testing.T) {
	assert.Equal(t, transformPath("\\a"), "a")
	assert.Equal(t, transformPath("\\a\\"), "a")
	assert.Equal(t, transformPath("\\a\\cc"), "a/cc")
}

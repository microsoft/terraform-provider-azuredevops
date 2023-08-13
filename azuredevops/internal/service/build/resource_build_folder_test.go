//go:build (all || resource_build_folder) && !exclude_resource_build_folder
// +build all resource_build_folder
// +build !exclude_resource_build_folder

package build

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

var testProjectUUID = uuid.New()
var testProjectReference = core.TeamProjectReference{
	Id: &testProjectUUID,
}

var testPath = "\\"

// This definition matches the overall structure of what a configured folder would look like
var testBuildFolder = build.Folder{
	Description: converter.String("My Folder Description"),
	Path:        converter.String(testPath),
	Project:     &testProjectReference,
}

// validates that an error is thrown if any of the un-supported path characters are used
func TestBuildFolder_PathInvalidCharacterListIsError(t *testing.T) {
	expectedInvalidPathCharacters := []string{"<", ">", "|", ":", "$", "@", "\"", "/", "%", "+", "*", "?"}
	pathSchema := ResourceBuildFolder().Schema["path"]

	for _, invalidCharacter := range expectedInvalidPathCharacters {
		_, errors := pathSchema.ValidateFunc(`\`+invalidCharacter, "")
		require.Equal(t, "<>|:$@\"/%+*? are not allowed in path", errors[0].Error())
	}
}

// validates that an error is thrown if path does not start with slash
func TestBuildFolder_PathInvalidStartingSlashIsError(t *testing.T) {
	pathSchema := ResourceBuildFolder().Schema["path"]
	_, errors := pathSchema.ValidateFunc("dir\\dir", "")
	require.Equal(t, "path must start with backslash", errors[0].Error())
}

// verifies that an expand will fail if there is insufficient configuration data found in the resource
func TestBuildFolder_Expand_FailsIfNotEnoughData(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceBuildFolder().Schema, nil)
	_, _, err := expandBuildFolder(resourceData)
	require.NotNil(t, err)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestBuildFolder_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := schema.TestResourceDataRaw(t, ResourceBuildFolder().Schema, nil)
	flattenBuildFolder(resourceData, &testBuildFolder, testProjectID)

	buildClient := azdosdkmocks.NewMockBuildClient(ctrl)
	clients := &client.AggregatedClient{BuildClient: buildClient, Ctx: context.Background()}

	buildClient.
		EXPECT().
		CreateFolder(clients.Ctx, gomock.Any()).
		Return(nil, errors.New("CreateFolder() Failed")).
		Times(1)

	err := resourceBuildFolderCreate(resourceData, clients)
	require.Contains(t, err.Error(), "CreateFolder() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestBuildFolder_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := schema.TestResourceDataRaw(t, ResourceBuildFolder().Schema, nil)
	flattenBuildFolder(resourceData, &testBuildFolder, testProjectID)

	buildClient := azdosdkmocks.NewMockBuildClient(ctrl)
	clients := &client.AggregatedClient{BuildClient: buildClient, Ctx: context.Background()}

	buildClient.
		EXPECT().
		GetFolders(clients.Ctx, gomock.Any()).
		Return(nil, errors.New("GetFolder() Failed")).
		Times(1)

	err := resourceBuildFolderRead(resourceData, clients)
	require.Equal(t, "GetFolder() Failed", err.Error())
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestBuildFolder_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := schema.TestResourceDataRaw(t, ResourceBuildFolder().Schema, nil)
	flattenBuildFolder(resourceData, &testBuildFolder, testProjectID)

	buildClient := azdosdkmocks.NewMockBuildClient(ctrl)
	clients := &client.AggregatedClient{BuildClient: buildClient, Ctx: context.Background()}

	buildClient.
		EXPECT().
		DeleteFolder(clients.Ctx, gomock.Any()).
		Return(errors.New("DeleteFolder() Failed")).
		Times(1)

	err := resourceBuildFolderDelete(resourceData, clients)
	require.Equal(t, "DeleteFolder() Failed", err.Error())
}

// verifies that if an error is produced on an update, it is not swallowed
func TestBuildFolder_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := schema.TestResourceDataRaw(t, ResourceBuildFolder().Schema, nil)
	flattenBuildFolder(resourceData, &testBuildFolder, testProjectID)

	buildClient := azdosdkmocks.NewMockBuildClient(ctrl)
	clients := &client.AggregatedClient{BuildClient: buildClient, Ctx: context.Background()}

	buildClient.
		EXPECT().
		UpdateFolder(clients.Ctx, gomock.Any()).
		Return(nil, errors.New("UpdateFolder() Failed")).
		Times(1)

	err := resourceBuildFolderUpdate(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateFolder() Failed")
}

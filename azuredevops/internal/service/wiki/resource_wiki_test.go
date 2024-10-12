//go:build (all || wiki || resource_wiki) && !exclude_resource_wiki
// +build all wiki resource_wiki
// +build !exclude_resource_wiki

package wiki

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/wiki"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

var testWikiProjectID = uuid.New()
var testWikiID = uuid.New()
var testWikiType = wiki.WikiType("codeWiki")

func TestWiki_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wikiClient := azdosdkmocks.NewMockWikiClient(ctrl)
	clients := &client.AggregatedClient{WikiClient: wikiClient, Ctx: context.Background()}

	resourceData := schema.TestResourceDataRaw(t, ResourceWiki().Schema, nil)
	resourceData.SetId(testWikiID.String())
	resourceData.Set("name", "testwiki")
	resourceData.Set("project_id", testWikiProjectID.String())
	resourceData.Set("type", testWikiType)

	expectedArgs := wiki.CreateWikiArgs{
		WikiCreateParams: &wiki.WikiCreateParametersV2{
			Name:      converter.String("testwiki"),
			ProjectId: &testWikiProjectID,
			Type:      &testWikiType,
		},
	}

	wikiClient.
		EXPECT().
		CreateWiki(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateWiki() Failed")).
		Times(1)

	err := resourceWikiCreate(resourceData, clients)
	require.Regexp(t, ".*CreateWiki\\(\\) Failed$", err.Error())
}

func TestWiki_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wikiClient := azdosdkmocks.NewMockWikiClient(ctrl)
	clients := &client.AggregatedClient{WikiClient: wikiClient, Ctx: context.Background()}

	resourceData := schema.TestResourceDataRaw(t, ResourceWiki().Schema, nil)
	resourceData.SetId(testWikiID.String())

	expectedArgs := wiki.GetWikiArgs{WikiIdentifier: converter.String(testWikiID.String())}

	wikiClient.
		EXPECT().
		GetWiki(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetWiki() Failed")).
		Times(1)

	err := resourceWikiRead(resourceData, clients)
	require.Regexp(t, ".*GetWiki\\(\\) Failed$", err.Error())
}

func TestWiki_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wikiClient := azdosdkmocks.NewMockWikiClient(ctrl)
	clients := &client.AggregatedClient{WikiClient: wikiClient, Ctx: context.Background()}

	resourceData := schema.TestResourceDataRaw(t, ResourceWiki().Schema, nil)
	resourceData.SetId(testWikiID.String())
	resourceData.Set("name", "Something")
	resourceData.Set("project_id", testWikiProjectID.String())
	resourceData.Set("type", testWikiType)

	expectedArgs := wiki.DeleteWikiArgs{
		WikiIdentifier: converter.String(testWikiID.String()),
		Project:        converter.String(testWikiProjectID.String())}

	wikiClient.
		EXPECT().
		DeleteWiki(clients.Ctx, expectedArgs).
		Return(nil, errors.New("DeleteWiki() Failed")).
		Times(1)

	err := resourceWikiDelete(resourceData, clients)
	require.Regexp(t, ".*DeleteWiki\\(\\) Failed$", err.Error())
}

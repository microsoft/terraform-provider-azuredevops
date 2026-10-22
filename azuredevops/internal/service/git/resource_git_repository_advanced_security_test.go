package git

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/sdk/advancedsecurity"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	advSecProjectID    = uuid.New().String()
	advSecRepositoryID = uuid.New().String()
)

func newAdvancedSecurityResourceData(t *testing.T, config map[string]interface{}) *schema.ResourceData {
	raw := map[string]interface{}{
		"project_id":    advSecProjectID,
		"repository_id": advSecRepositoryID,
	}
	for key, value := range config {
		raw[key] = value
	}
	d := schema.TestResourceDataRaw(t, ResourceGitRepositoryAdvancedSecurity().Schema, raw)
	d.SetId(advSecProjectID + "/" + advSecRepositoryID)
	return d
}

func advancedSecurityClients(ctrl *gomock.Controller) (*client.AggregatedClient, *azdosdkmocks.MockAdvancedsecurityClient) {
	advSecClient := azdosdkmocks.NewMockAdvancedsecurityClient(ctrl)
	return &client.AggregatedClient{AdvancedSecurityClient: advSecClient, Ctx: context.Background()}, advSecClient
}

func enabledRepoEnablement() *advancedsecurity.RepoEnablementSettings {
	return &advancedsecurity.RepoEnablementSettings{
		SecretProtectionFeatures: &advancedsecurity.SecretProtectionFeatures{
			SecretProtectionEnabled: converter.Bool(true),
			BlockPushes:             converter.Bool(true),
		},
		CodeSecurityFeatures: &advancedsecurity.CodeSecurityFeatures{
			CodeSecurityEnabled:                converter.Bool(true),
			DependencyScanningInjectionEnabled: converter.Bool(true),
			CodeQLEnabled:                      converter.Bool(false),
			AutofixEnabled:                     nil,
		},
	}
}

func TestGitRepositoryAdvancedSecurity_Expand_SendsOnlyConfiguredSettings(t *testing.T) {
	settings := expandRepoEnablement(map[string]cty.Value{
		"project_id":                cty.StringVal(advSecProjectID),
		"secret_protection_enabled": cty.True,
		"block_pushes":              cty.NullVal(cty.Bool),
		"code_security_enabled":     cty.True,
		"codeql_enabled":            cty.False,
	})

	require.Equal(t, &advancedsecurity.RepoEnablementSettings{
		SecretProtectionFeatures: &advancedsecurity.SecretProtectionFeatures{SecretProtectionEnabled: converter.Bool(true)},
		CodeSecurityFeatures: &advancedsecurity.CodeSecurityFeatures{
			CodeSecurityEnabled: converter.Bool(true),
			CodeQLEnabled:       converter.Bool(false),
		},
	}, settings)
}

func TestGitRepositoryAdvancedSecurity_Expand_OmitsUnconfiguredFeatures(t *testing.T) {
	settings := expandRepoEnablement(map[string]cty.Value{"block_pushes": cty.False})

	require.Nil(t, settings.CodeSecurityFeatures)
	require.Equal(t, &advancedsecurity.SecretProtectionFeatures{BlockPushes: converter.Bool(false)}, settings.SecretProtectionFeatures)
}

func TestGitRepositoryAdvancedSecurity_Create_UpdatesAndReads(t *testing.T) {
	ctrl := gomock.NewController(t)
	clients, advSecClient := advancedSecurityClients(ctrl)
	d := newAdvancedSecurityResourceData(t, nil)
	d.SetId("")

	advSecClient.EXPECT().UpdateRepoEnablement(clients.Ctx, advancedsecurity.UpdateRepoEnablementArgs{
		Project:    converter.String(advSecProjectID),
		Repository: converter.String(advSecRepositoryID),
		Settings:   &advancedsecurity.RepoEnablementSettings{},
	}).Return(nil).Times(1)
	advSecClient.EXPECT().GetRepoEnablement(clients.Ctx, gomock.Any()).Return(enabledRepoEnablement(), nil).Times(1)

	diags := resourceGitRepositoryAdvancedSecurityCreateUpdate(clients.Ctx, d, clients)

	require.False(t, diags.HasError(), "%v", diags)
	require.Equal(t, advSecProjectID+"/"+advSecRepositoryID, d.Id())
	require.True(t, d.Get("secret_protection_enabled").(bool))
}

func TestGitRepositoryAdvancedSecurity_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	clients, advSecClient := advancedSecurityClients(ctrl)
	d := newAdvancedSecurityResourceData(t, nil)

	advSecClient.EXPECT().UpdateRepoEnablement(clients.Ctx, gomock.Any()).Return(errors.New("@@UpdateRepoEnablement@@failed@@")).Times(1)

	diags := resourceGitRepositoryAdvancedSecurityCreateUpdate(clients.Ctx, d, clients)

	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "@@UpdateRepoEnablement@@failed@@")
}

func TestGitRepositoryAdvancedSecurity_Read_SetsAllSettings(t *testing.T) {
	ctrl := gomock.NewController(t)
	clients, advSecClient := advancedSecurityClients(ctrl)
	d := newAdvancedSecurityResourceData(t, nil)

	advSecClient.EXPECT().GetRepoEnablement(clients.Ctx, advancedsecurity.GetRepoEnablementArgs{
		Project:              converter.String(advSecProjectID),
		Repository:           converter.String(advSecRepositoryID),
		IncludeAllProperties: converter.Bool(true),
	}).Return(enabledRepoEnablement(), nil).Times(1)

	diags := resourceGitRepositoryAdvancedSecurityRead(clients.Ctx, d, clients)

	require.False(t, diags.HasError(), "%v", diags)
	require.True(t, d.Get("secret_protection_enabled").(bool))
	require.True(t, d.Get("block_pushes").(bool))
	require.True(t, d.Get("code_security_enabled").(bool))
	require.True(t, d.Get("dependency_scanning_injection_enabled").(bool))
	require.False(t, d.Get("codeql_enabled").(bool))
	require.False(t, d.Get("autofix_enabled").(bool))
}

func TestGitRepositoryAdvancedSecurity_Read_RemovesResourceWhenNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	clients, advSecClient := advancedSecurityClients(ctrl)
	d := newAdvancedSecurityResourceData(t, nil)

	advSecClient.EXPECT().GetRepoEnablement(clients.Ctx, gomock.Any()).Return(nil, azuredevops.WrappedError{StatusCode: converter.Int(http.StatusNotFound)}).Times(1)

	diags := resourceGitRepositoryAdvancedSecurityRead(clients.Ctx, d, clients)

	require.False(t, diags.HasError(), "%v", diags)
	require.Equal(t, "", d.Id())
}

func TestGitRepositoryAdvancedSecurity_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	clients, advSecClient := advancedSecurityClients(ctrl)
	d := newAdvancedSecurityResourceData(t, nil)

	advSecClient.EXPECT().GetRepoEnablement(clients.Ctx, gomock.Any()).Return(nil, errors.New("@@GetRepoEnablement@@failed@@")).Times(1)

	diags := resourceGitRepositoryAdvancedSecurityRead(clients.Ctx, d, clients)

	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "@@GetRepoEnablement@@failed@@")
}

func TestGitRepositoryAdvancedSecurity_Delete_DisablesAdvancedSecurity(t *testing.T) {
	ctrl := gomock.NewController(t)
	clients, advSecClient := advancedSecurityClients(ctrl)
	d := newAdvancedSecurityResourceData(t, nil)

	advSecClient.EXPECT().UpdateRepoEnablement(clients.Ctx, advancedsecurity.UpdateRepoEnablementArgs{
		Project:    converter.String(advSecProjectID),
		Repository: converter.String(advSecRepositoryID),
		Settings: &advancedsecurity.RepoEnablementSettings{
			SecretProtectionFeatures: &advancedsecurity.SecretProtectionFeatures{SecretProtectionEnabled: converter.Bool(false)},
			CodeSecurityFeatures:     &advancedsecurity.CodeSecurityFeatures{CodeSecurityEnabled: converter.Bool(false)},
		},
	}).Return(nil).Times(1)

	diags := resourceGitRepositoryAdvancedSecurityDelete(clients.Ctx, d, clients)

	require.False(t, diags.HasError(), "%v", diags)
	require.Equal(t, "", d.Id())
}

func TestGitRepositoryAdvancedSecurity_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	clients, advSecClient := advancedSecurityClients(ctrl)
	d := newAdvancedSecurityResourceData(t, nil)

	advSecClient.EXPECT().UpdateRepoEnablement(clients.Ctx, gomock.Any()).Return(errors.New("@@UpdateRepoEnablement@@failed@@")).Times(1)

	diags := resourceGitRepositoryAdvancedSecurityDelete(clients.Ctx, d, clients)

	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "@@UpdateRepoEnablement@@failed@@")
}

func TestGitRepositoryAdvancedSecurity_FailsWithoutClient(t *testing.T) {
	clients := &client.AggregatedClient{OrganizationURL: "https://dev.azure.com/example", Ctx: context.Background()}
	d := newAdvancedSecurityResourceData(t, nil)

	diags := resourceGitRepositoryAdvancedSecurityRead(clients.Ctx, d, clients)

	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "Advanced Security is not available")
}

func TestGitRepositoryAdvancedSecurity_Import(t *testing.T) {
	for _, tc := range []struct {
		id      string
		wantErr bool
	}{
		{id: advSecProjectID + "/" + advSecRepositoryID},
		{id: advSecRepositoryID, wantErr: true},
		{id: advSecProjectID + "/", wantErr: true},
		{id: "a/b/c", wantErr: true},
	} {
		d := schema.TestResourceDataRaw(t, ResourceGitRepositoryAdvancedSecurity().Schema, nil)
		d.SetId(tc.id)

		result, err := resourceGitRepositoryAdvancedSecurityImport(context.Background(), d, nil)

		if tc.wantErr {
			require.Error(t, err, tc.id)
			continue
		}
		require.NoError(t, err)
		require.Equal(t, advSecProjectID, result[0].Get("project_id"))
		require.Equal(t, advSecRepositoryID, result[0].Get("repository_id"))
	}
}

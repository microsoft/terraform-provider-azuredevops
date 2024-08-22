//go:build all || resource_personal_access_token
// +build all resource_personal_access_token

package tokens

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/tokens"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

func TestResourcePersonalAccessToken_CreateWithInvalidValidToDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	valid_to := "invalid-date"
	resourceData := generateResourceData(t, nil, nil, nil, &valid_to, false)
	_, clients := generateMocks(ctrl)

	err := resourceAzurePersonalAccessTokenCreate(resourceData, clients)
	require.Contains(t, err.Error(), "parsing valid to date:")
}

func TestResourcePersonalAccessToken_CreateFailsWithExpectedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tokenName := "test-token"
	scopes := []string{"vso.tokens", "vso.packages"}
	valid_to := "1969-12-31T00:00:00Z"
	resourceData := generateResourceData(t, nil, &tokenName, &scopes, &valid_to, true)
	tokenClient, clients := generateMocks(ctrl)

	create_scopes := "vso.tokens vso.packages"
	valid_to_time, _ := time.Parse(time.RFC3339, valid_to)
	tokenClient.
		EXPECT().
		CreatePat(clients.Ctx, tokens.CreatePatArgs{
			Token: &tokens.PatTokenCreateRequest{
				AllOrgs:     converter.Bool(true),
				DisplayName: &tokenName,
				Scope:       &create_scopes,
				ValidTo:     &azuredevops.Time{Time: valid_to_time},
			},
		}).
		Return(nil, errors.New("create-error")).
		Times(1)

	err := resourceAzurePersonalAccessTokenCreate(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "creating pat token in Azure DevOps:")
}

func TestResourcePersonalAccessToken_CreateFailsWithExpectedReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tokenName := "test-token"
	scopes := []string{"vso.tokens", "vso.packages"}
	valid_to := "1969-12-31T00:00:00Z"
	resourceData := generateResourceData(t, nil, &tokenName, &scopes, &valid_to, true)
	tokenClient, clients := generateMocks(ctrl)

	authorization_id := uuid.New()
	create_scopes := "vso.tokens vso.packages"
	valid_to_time, _ := time.Parse(time.RFC3339, valid_to)
	tokenClient.
		EXPECT().
		CreatePat(clients.Ctx, tokens.CreatePatArgs{
			Token: &tokens.PatTokenCreateRequest{
				AllOrgs:     converter.Bool(true),
				DisplayName: &tokenName,
				Scope:       &create_scopes,
				ValidTo:     &azuredevops.Time{Time: valid_to_time},
			},
		}).
		Return(&tokens.PatTokenResult{
			PatToken: tokens.PatToken{
				AuthorizationId: &authorization_id,
			},
			PatTokenError: nil,
		}, nil).
		Times(1)

	tokenClient.
		EXPECT().
		GetPat(clients.Ctx, tokens.GetPatArgs{AuthorizationId: &authorization_id}).
		Return(nil, errors.New("GetPat() Failed")).
		Times(1)

	err := resourceAzurePersonalAccessTokenCreate(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "GetPat() Failed")
}

func TestResourcePersonalAccessToken_ReadWithInvalidAuthorizationId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authorization_id := &uuid.Nil
	resourceData := generateResourceData(t, authorization_id, nil, nil, nil, false)
	_, clients := generateMocks(ctrl)

	err := resourceAzurePersonalAccessTokenRead(resourceData, clients)
	require.Contains(t, err.Error(), "parse token authorization ID")
}

func TestResourcePersonalAccessToken_ReadReturnsWithExpectedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authorization_id := uuid.New()
	resourceData := generateResourceData(t, &authorization_id, nil, nil, nil, false)
	tokenClient, clients := generateMocks(ctrl)

	tokenClient.
		EXPECT().
		GetPat(clients.Ctx, tokens.GetPatArgs{AuthorizationId: &authorization_id}).
		Return(nil, errors.New("GetPat() Failed")).
		Times(1)

	err := resourceAzurePersonalAccessTokenRead(resourceData, clients)
	require.Contains(t, err.Error(), "GetPat() Failed")
}

func TestResourcePersonalAccessToken_UpdateWithInvalidAuthorizationId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authorization_id := &uuid.Nil
	resourceData := generateResourceData(t, authorization_id, nil, nil, nil, false)
	_, clients := generateMocks(ctrl)

	err := resourceAzurePersonalAccessTokenUpdate(resourceData, clients)
	require.Contains(t, err.Error(), "parse token authorization ID")
}

func TestResourcePersonalAccessToken_UpdateWithInvalidValidToDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	autorization_id := uuid.New()
	valid_to := "invalid-date"
	resourceData := generateResourceData(t, &autorization_id, nil, nil, &valid_to, false)
	_, clients := generateMocks(ctrl)

	err := resourceAzurePersonalAccessTokenUpdate(resourceData, clients)
	require.Contains(t, err.Error(), "parsing valid to date:")
}

func TestResourcePersonalAccessToken_UpdateFailsWithExpectedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authorization_id := uuid.New()
	tokenName := "test-token"
	scopes := []string{"vso.tokens", "vso.packages"}
	valid_to := "1969-12-31T00:00:00Z"
	resourceData := generateResourceData(t, &authorization_id, &tokenName, &scopes, &valid_to, true)
	tokenClient, clients := generateMocks(ctrl)

	create_scopes := "vso.tokens vso.packages"
	valid_to_time, _ := time.Parse(time.RFC3339, valid_to)
	tokenClient.
		EXPECT().
		UpdatePat(clients.Ctx, tokens.UpdatePatArgs{
			Token: &tokens.PatTokenUpdateRequest{
				AllOrgs:         converter.Bool(true),
				AuthorizationId: &authorization_id,
				DisplayName:     &tokenName,
				Scope:           &create_scopes,
				ValidTo:         &azuredevops.Time{Time: valid_to_time},
			},
		}).
		Return(nil, errors.New("update-error")).
		Times(1)

	err := resourceAzurePersonalAccessTokenUpdate(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "updating Personal Access Token in Azure DevOps:")
}

func TestResourcePersonalAccessToken_UpdateFailsWithExpectedReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authorization_id := uuid.New()
	tokenName := "test-token"
	scopes := []string{"vso.tokens", "vso.packages"}
	valid_to := "1969-12-31T00:00:00Z"
	resourceData := generateResourceData(t, &authorization_id, &tokenName, &scopes, &valid_to, true)
	tokenClient, clients := generateMocks(ctrl)

	create_scopes := "vso.tokens vso.packages"
	valid_to_time, _ := time.Parse(time.RFC3339, valid_to)
	tokenClient.
		EXPECT().
		UpdatePat(clients.Ctx, tokens.UpdatePatArgs{
			Token: &tokens.PatTokenUpdateRequest{
				AllOrgs:         converter.Bool(true),
				AuthorizationId: &authorization_id,
				DisplayName:     &tokenName,
				Scope:           &create_scopes,
				ValidTo:         &azuredevops.Time{Time: valid_to_time},
			},
		}).
		Return(&tokens.PatTokenResult{
			PatToken: tokens.PatToken{
				AuthorizationId: &uuid.Nil,
			},
			PatTokenError: nil,
		}, nil).
		Times(1)

	tokenClient.
		EXPECT().
		GetPat(clients.Ctx, tokens.GetPatArgs{AuthorizationId: &authorization_id}).
		Return(nil, errors.New("GetPat() Failed")).
		Times(1)

	err := resourceAzurePersonalAccessTokenUpdate(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "GetPat() Failed")
}

func TestResourcePersonalAccessToken_RevokeWithInvalidAuthorizationId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authorization_id := &uuid.Nil
	resourceData := generateResourceData(t, authorization_id, nil, nil, nil, false)
	_, clients := generateMocks(ctrl)

	err := resourceAzurePersonalAccessTokenRevoke(resourceData, clients)
	require.Contains(t, err.Error(), "parse token authorization ID")
}

func TestResourcePersonalAccessToken_RevokeFailsWithExpectedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authorization_id := uuid.New()
	resourceData := generateResourceData(t, &authorization_id, nil, nil, nil, true)
	tokenClient, clients := generateMocks(ctrl)

	tokenClient.
		EXPECT().
		Revoke(clients.Ctx, tokens.RevokeArgs{AuthorizationId: &authorization_id}).
		Return(errors.New("revoke-error")).
		Times(1)

	err := resourceAzurePersonalAccessTokenRevoke(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "revoke Personal Access Token in Azure DevOps:")
}

func generateResourceData(t *testing.T, authorizationId *uuid.UUID, tokenName *string, scopes *[]string, validTo *string, allOrgs bool) *schema.ResourceData {
	resourceData := schema.TestResourceDataRaw(t, ResourcePersonalAccessToken().Schema, nil)
	resourceData.Set("all_orgs", allOrgs)

	if authorizationId != nil {
		resourceData.SetId(authorizationId.String())
		resourceData.Set("authorization_id", authorizationId.String())
	}

	if tokenName != nil {
		resourceData.Set("name", *tokenName)
	}

	if scopes != nil {
		resourceData.Set("scope", scopes)
	}

	if validTo != nil {
		resourceData.Set("valid_to", *validTo)
	}

	return resourceData
}

func generateMocks(ctrl *gomock.Controller) (*azdosdkmocks.MockTokenClient, *client.AggregatedClient) {
	tokenClient := azdosdkmocks.NewMockTokenClient(ctrl)
	return tokenClient, &client.AggregatedClient{
		TokensClient: tokenClient,
		Ctx:          context.Background(),
	}
}

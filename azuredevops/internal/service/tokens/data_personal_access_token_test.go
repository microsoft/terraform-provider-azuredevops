//go:build all || core || data_projects
// +build all core data_projects

package tokens

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/tokens"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/stretchr/testify/require"
)

func TestDataSourcePersonalAccessToken_Read_TestPersonalAccessTokeInvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tokensClient := azdosdkmocks.NewMockTokenClient(ctrl)
	clients := &client.AggregatedClient{
		TokensClient: tokensClient,
		Ctx:          context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, DataPersonalAccessToken().Schema, nil)
	resourceData.Set("authorization_id", "invalid-uuid")
	err := dataPersonalAccessTokenRead(resourceData, clients)
	require.Contains(t, err.Error(), "parse token authorization ID:")
}

func TestDataSourcePersonalAccessToken_Read_TestPersonalAccessTokenNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tokensClient := azdosdkmocks.NewMockTokenClient(ctrl)
	clients := &client.AggregatedClient{
		TokensClient: tokensClient,
		Ctx:          context.Background(),
	}

	authorization_id := uuid.New()
	tokensClient.EXPECT().
		GetPat(clients.Ctx, tokens.GetPatArgs{AuthorizationId: &authorization_id}).
		Return(nil, errors.New("unauthorized")).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataPersonalAccessToken().Schema, nil)
	resourceData.Set("authorization_id", authorization_id.String())
	err := dataPersonalAccessTokenRead(resourceData, clients)
	require.Contains(t, err.Error(), "Error getting personal access token by authorization ID: unauthorized")
}

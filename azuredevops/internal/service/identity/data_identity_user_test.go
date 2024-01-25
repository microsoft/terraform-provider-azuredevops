//go:build (all || core || data_sources || data_users) && (!exclude_data_sources || !exclude_data_users)
// +build all core data_sources data_users
// +build !exclude_data_sources !exclude_data_users

package identity

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

var id, _ = uuid.Parse("00000000-0000-0000-0000-000000000000")

var usrList1 = []identity.Identity{
	{
		Descriptor:          converter.String("aad.YWU4YzJkMTktOThmYS03ZDhmLWJhNTAtOWI4MWQzYTUxZjcy"),
		ProviderDisplayName: converter.String("DesireeMCollins@jourrapide.com"),
	},
	{
		Descriptor:          converter.String("svc.Nzc3NGFjMDMtOGEyOS00NGFjLTg2ZjEtZmE0YmRlZDc4ZGUyOkdpdEh1YiBBcHA6MjE3OGVhMGItMmNlMi00ZDA4LTg4YTMtODdiYWRjYjkzNjE1"),
		ProviderDisplayName: converter.String("GitHub"),
	},
	{
		Descriptor:          converter.String("aad.NTJmM2YzZmMtZmE4NS00ZGNhLTkwYjUtMDNiY2U0NjBlMmIy"),
		ProviderDisplayName: converter.String("WalterMBrooks@rhyta.com"),
	},
	{
		Descriptor:          converter.String("msa.MmRkMWZkMWMtNWE0Mi00NDM1LThhM2ItMWQ1NWUwOTA2NzUx"),
		ProviderDisplayName: converter.String("KerryMRaymond@teleworm.us"),
	},
}

var usrList2 = []identity.Identity{
	{
		Descriptor:          converter.String("svc.Nzc3NGFjMDMtOGEyOS00NGFjLTg2ZjEtZmE0YmRlZDc4ZGUyOkJ1aWxkOmYwNzg2ZGZkLTA4OGYtNDYxOS1hY2NjLTJlYzc1ZTEyNmRjZQ"),
		ProviderDisplayName: converter.String("Project Collection Build Service (unittest)"),
	},
	{
		Descriptor:          converter.String("bnd.dXBuOmIwYjg0Yzc1LTI4NjctNGQ0Mi1iNWFhLTg1YjE5Nzc1ZDdlN1xBbGxlbkJNY0tpbm5vbkBkYXlyZXAuY29t"),
		ProviderDisplayName: converter.String("AllenBMcKinnon@dayrep.com"),
	},
	{
		Descriptor:          converter.String("win.Uy0xLTUtMjEtMjA4NTA0NzkxOS0yODAyODgxNjk5LTE5MDg1MTA1MTctNTAw"),
		ProviderDisplayName: converter.String("Local Admin"),
	},
}

// verfies that the data source propagates an error from the API correctly
func TestDataSourceIdentityUser_Read_TestDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{
		IdentityClient: identityClient,
		Ctx:            context.Background(),
	}
	projectID := uuid.New()
	projectIDstring := projectID.String()
	expectedArgs := identity.ReadIdentitiesArgs{
		IdentityIds: &projectIDstring,
	}
	identityClient.
		EXPECT().
		ReadIdentities(clients.Ctx, expectedArgs).
		Return(nil, errors.New("ReadIdentities() Failed"))

	resourceData := schema.TestResourceDataRaw(t, DataIdentityUser().Schema, nil)
	err := dataIdentitySourceUserRead(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "ReadIdentities() Failed")
}

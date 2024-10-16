//go:build (all || resource_serviceendpoint_github_enterprise) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_github_enterprise
// +build !exclude_serviceendpoints

package serviceendpoint

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

var ghesTestServiceEndpointID = uuid.New()
var ghesRandomServiceEndpointProjectID = uuid.New()
var ghesTestServiceEndpointProjectID = &ghesRandomServiceEndpointProjectID

var ghesTestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": "UNIT_TEST_ACCESS_TOKEN",
		},
		Scheme: converter.String("Token"),
	},
	Id:          &ghesTestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_NAME"),
	Owner:       converter.String("library"),
	Type:        converter.String("githubenterprise"),
	Url:         converter.String("https://github.contoso.com"),
	Description: converter.String("UNIT_TEST_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: ghesTestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_NAME"),
			Description: converter.String("UNIT_TEST_DESCRIPTION"),
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func TestServiceEndpointGitHubEnterprise_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointGitHubEnterprise().Schema, nil)
	configureGhesAuthPersonal(resourceData)
	flattenServiceEndpointGitHubEnterprise(resourceData, &ghesTestServiceEndpoint)

	serviceEndpointAfterRoundTrip, err := expandServiceEndpointGitHubEnterprise(resourceData)

	require.Nil(t, err)
	require.Equal(t, ghesTestServiceEndpoint, *serviceEndpointAfterRoundTrip)
	require.Equal(t, ghesTestServiceEndpointProjectID, (*serviceEndpointAfterRoundTrip.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestServiceEndpointGitHubEnterprise_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointGitHubEnterprise()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureGhesAuthPersonal(resourceData)
	flattenServiceEndpointGitHubEnterprise(resourceData, &ghesTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &ghesTestServiceEndpoint}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointGitHubEnterprise_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointGitHubEnterprise()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointGitHubEnterprise(resourceData, &ghesTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: ghesTestServiceEndpoint.Id,
		Project:    converter.String(ghesTestServiceEndpointProjectID.String()),
	}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestServiceEndpointGitHubEnterprise_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointGitHubEnterprise()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointGitHubEnterprise(resourceData, &ghesTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: ghesTestServiceEndpoint.Id,
		ProjectIds: &[]string{
			ghesTestServiceEndpointProjectID.String(),
		},
	}

	buildClient.
		EXPECT().
		DeleteServiceEndpoint(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteServiceEndpoint() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteServiceEndpoint() Failed")
}

// verifies that if an error is produced on an update, it is not swallowed
func TestServiceEndpointGitHubEnterprise_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointGitHubEnterprise()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureGhesAuthPersonal(resourceData)
	flattenServiceEndpointGitHubEnterprise(resourceData, &ghesTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &ghesTestServiceEndpoint,
		EndpointId: ghesTestServiceEndpoint.Id,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

func configureGhesAuthPersonal(d *schema.ResourceData) {
	d.Set("auth_personal", &[]map[string]interface{}{
		{
			"personal_access_token": "UNIT_TEST_ACCESS_TOKEN",
		},
	})
}

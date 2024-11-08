//go:build (all || resource_serviceendpoint_azurecr) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_azurecr
// +build !exclude_serviceendpoints

package serviceendpoint

import (
	"context"
	"errors"
	"fmt"
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

var azureCRTestServiceEndpointID = uuid.New()
var azureCRRandomServiceEndpointProjectID = uuid.New()
var azureCRTestServiceEndpointProjectID = &azureCRRandomServiceEndpointProjectID
var subscription_id = "42125daf-72fd-417c-9ea7-080690625ad3"
var scope = fmt.Sprintf(
	"/subscriptions/%s/resourceGroups/testrg/providers/Microsoft.ContainerRegistry/registries/testacr",
	subscription_id,
)

var azureCRTestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"authenticationType": "spnKey",
			"tenantId":           "aba07645-051c-44b4-b806-c34d33f3dcd1", //fake value
			"loginServer":        "testacr.azurecr.io",
			"scope":              scope,
			"serviceprincipalid": "00000000-0000-0000-0000-000000000000",
		},
		Scheme: converter.String("ServicePrincipal"),
	},
	Data: &map[string]string{
		"registryId":               scope,
		"subscriptionId":           subscription_id,
		"subscriptionName":         "testS",
		"registrytype":             "ACR",
		"appObjectId":              "00000000-0000-0000-0000-000000000000",
		"spnObjectId":              "00000000-0000-0000-0000-000000000000",
		"azureSpnPermissions":      "[{\\\"roleAssignmentId\\\":\\\"00000000 - 0000-0000-0000-000000000000\\\",\\\"resourceProvider\\\":\\\"Microsoft.RoleAssignment\\\",\\\"provisioned\\\":true}]",
		"azureSpnRoleAssignmentId": "00000000-0000-0000-0000-000000000000",
	},
	Id:          &azureCRTestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("dockerregistry"),
	Url:         converter.String("https://testacr.azurecr.io"),
	Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: azureCRTestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func TestServiceEndpointAzureCR_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointAzureCR().Schema, nil)
	resourceData.Set("project_id", (*azureCRTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointAzureCR(resourceData, &azureCRTestServiceEndpoint)

	serviceEndpointAfterRoundTrip, err := expandServiceEndpointAzureCR(resourceData)

	require.Equal(t, azureCRTestServiceEndpoint, *serviceEndpointAfterRoundTrip)
	require.Equal(t, azureCRTestServiceEndpointProjectID, (*serviceEndpointAfterRoundTrip.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
	require.Nil(t, err)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestServiceEndpointAzureCR_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAzureCR()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*azureCRTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointAzureCR(resourceData, &azureCRTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &azureCRTestServiceEndpoint}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointAzureCR_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAzureCR()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*azureCRTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointAzureCR(resourceData, &azureCRTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: azureCRTestServiceEndpoint.Id,
		Project:    converter.String(azureCRTestServiceEndpointProjectID.String()),
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
func TestServiceEndpointAzureCR_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAzureCR()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*azureCRTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointAzureCR(resourceData, &azureCRTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: azureCRTestServiceEndpoint.Id,
		ProjectIds: &[]string{
			azureCRTestServiceEndpointProjectID.String(),
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
func TestServiceEndpointAzureCR_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAzureCR()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*azureCRTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointAzureCR(resourceData, &azureCRTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &azureCRTestServiceEndpoint,
		EndpointId: azureCRTestServiceEndpoint.Id,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

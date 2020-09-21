// +build all resource_serviceendpoint_devops
// +build !exclude_serviceendpoints

package serviceendpoint

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

var azdoTestServiceEndpointID = uuid.New()
var azdoRandomServiceEndpointProjectID = uuid.New().String()
var azdoTestServiceEndpointProjectID = &azdoRandomServiceEndpointProjectID

var azdoTestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": "UNIT_TEST_ACCESS_TOKEN",
		},
		Scheme: converter.String("Token"),
	},
	Data: &map[string]string{
		"releaseUrl": "https://vsrm.dev.azure.com/example",
	},
	Id:          &azdoTestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_NAME"),
	Description: converter.String("UNIT_TEST_DESCRIPTION"),
	Owner:       converter.String("library"),
	Type:        converter.String("azdoapi"),
	Url:         converter.String("https://dev.azure.com/example"),
}

// verifies that the flatten/expand round trip yields the same service endpoint
func TestServiceEndpointAzureDevOps_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointAzureDevOps().Schema, nil)
	azdoConfigureExtraFields(resourceData)
	flattenServiceEndpointAzureDevOps(resourceData, &azdoTestServiceEndpoint, azdoTestServiceEndpointProjectID)

	serviceEndpointAfterRoundTrip, projectID, err := expandServiceEndpointAzureDevOps(resourceData)

	require.Nil(t, err)
	require.Equal(t, azdoTestServiceEndpoint, *serviceEndpointAfterRoundTrip)
	require.Equal(t, azdoTestServiceEndpointProjectID, projectID)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestServiceEndpointAzureDevOps_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAzureDevOps()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	azdoConfigureExtraFields(resourceData)
	flattenServiceEndpointAzureDevOps(resourceData, &azdoTestServiceEndpoint, azdoTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &azdoTestServiceEndpoint, Project: azdoTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointAzureDevOps_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAzureDevOps()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointAzureDevOps(resourceData, &azdoTestServiceEndpoint, azdoTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{EndpointId: azdoTestServiceEndpoint.Id, Project: azdoTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestServiceEndpointAzureDevOps_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAzureDevOps()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointAzureDevOps(resourceData, &azdoTestServiceEndpoint, azdoTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{EndpointId: azdoTestServiceEndpoint.Id, Project: azdoTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		DeleteServiceEndpoint(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteServiceEndpoint() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteServiceEndpoint() Failed")
}

// verifies that if an error is produced on an update, it is not swallowed
func TestServiceEndpointAzureDevOps_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAzureDevOps()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	azdoConfigureExtraFields(resourceData)
	flattenServiceEndpointAzureDevOps(resourceData, &azdoTestServiceEndpoint, azdoTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &azdoTestServiceEndpoint,
		EndpointId: azdoTestServiceEndpoint.Id,
		Project:    azdoTestServiceEndpointProjectID,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

func azdoConfigureExtraFields(d *schema.ResourceData) {
	d.Set("organization_name", "example")
	d.Set("auth_personal", &[]map[string]interface{}{
		{
			"personal_access_token": "UNIT_TEST_ACCESS_TOKEN",
		},
	})
}

// +build all resource_serviceendpoint_rapid7
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

var rapid7TestServiceEndpointID = uuid.New()
var rapid7RandomServiceEndpointProjectID = uuid.New().String()
var rapid7TestServiceEndpointProjectID = &rapid7RandomServiceEndpointProjectID

var rapid7TestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": "UNIT_TEST_auth_token",
		},
		Scheme: converter.String("Token"),
	},
	Data: &map[string]string{
		"region": "UNIT_TEST_region",
	},
	Id:          &rapid7TestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("ias"),
	Url:         converter.String("https://rapid7.com"),
}

// verifies that the flatten/expand round trip yields the same service endpoint
func TestServiceEndpointRapid7_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointRapid7().Schema, nil)
	flattenServiceEndpoinRapid7(resourceData, &rapid7TestServiceEndpoint, rapid7TestServiceEndpointProjectID)

	serviceEndpointAfterRoundTrip, projectID, err := expandServiceEndpointRapid7qube(resourceData)

	require.Equal(t, rapid7TestServiceEndpoint, *serviceEndpointAfterRoundTrip)
	require.Equal(t, rapid7TestServiceEndpointProjectID, projectID)
	require.Nil(t, err)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestServiceEndpointRapid7_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointRapid7()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointRapid7(resourceData, &rapid7TestServiceEndpoint, rapid7TestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &rapid7TestServiceEndpoint, Project: rapid7TestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointRapid7_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointRapid7()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointRapid7(resourceData, &rapid7TestServiceEndpoint, rapid7TestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{EndpointId: rapid7TestServiceEndpoint.Id, Project: rapid7TestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestServiceEndpointRapid7_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointRapid7()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointRapid7(resourceData, &rapid7TestServiceEndpoint, rapid7TestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{EndpointId: rapid7TestServiceEndpoint.Id, Project: rapid7TestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		DeleteServiceEndpoint(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteServiceEndpoint() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteServiceEndpoint() Failed")
}

func TestServiceEndpointRapid7_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointRapid7()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointRapid7(resourceData, &rapid7TestServiceEndpoint, rapid7TestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &rapid7TestServiceEndpoint,
		EndpointId: rapid7TestServiceEndpoint.Id,
		Project:    rapid7TestServiceEndpointProjectID,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

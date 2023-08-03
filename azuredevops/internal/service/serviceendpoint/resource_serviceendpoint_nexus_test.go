//go:build (all || resource_serviceendpoint_nexus) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_nexus
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

var nexusTestServiceEndpointIDpassword = uuid.New()
var nexusRandomServiceEndpointProjectIDpassword = uuid.New()
var nexusTestServiceEndpointProjectIDpassword = &nexusRandomServiceEndpointProjectIDpassword

var nexusTestServiceEndpointPassword = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": "",
			"password": "",
		},
		Scheme: converter.String("UsernamePassword"),
	},
	Id:          &nexusTestServiceEndpointIDpassword,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("NexusIqServiceConnection"),
	Url:         converter.String("https://www.nexus.com"),
	Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: nexusTestServiceEndpointProjectIDpassword,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func testServiceEndpointNexus_ExpandFlatten_Roundtrip(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	for _, ep := range []*serviceendpoint.ServiceEndpoint{ep, ep} {
		resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointNexus().Schema, nil)
		flattenServiceEndpointNexus(resourceData, ep, id)

		serviceEndpointAfterRoundTrip, projectID, err := expandServiceEndpointNexus(resourceData)

		require.Nil(t, err)
		require.Equal(t, *ep, *serviceEndpointAfterRoundTrip)
		require.Equal(t, id, projectID)
	}
}
func TestServiceEndpointNexus_ExpandFlatten_RoundtripPassword(t *testing.T) {
	testServiceEndpointNexus_ExpandFlatten_Roundtrip(t, &nexusTestServiceEndpointPassword, nexusTestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on create, the error is not swallowed
func testServiceEndpointNexus_Create_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointNexus()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointNexus(resourceData, ep, id)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: ep}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}
func TestServiceEndpointNexus_Create_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointNexus_Create_DoesNotSwallowError(t, &nexusTestServiceEndpointPassword, nexusTestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on read, the error is not swallowed
func testServiceEndpointNexus_Read_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointNexus()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointNexus(resourceData, ep, id)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: ep.Id,
		Project:    converter.String(id.String()),
	}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
}
func TestServiceEndpointNexus_Read_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointNexus_Read_DoesNotSwallowError(t, &nexusTestServiceEndpointPassword, nexusTestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on a delete, it is not swallowed
func testServiceEndpointNexus_Delete_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointNexus()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointNexus(resourceData, ep, id)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: ep.Id,
		ProjectIds: &[]string{
			id.String(),
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
func TestServiceEndpointNexus_Delete_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointNexus_Delete_DoesNotSwallowError(t, &nexusTestServiceEndpointPassword, nexusTestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on a update, it is not swallowed
func testServiceEndpointNexus_Update_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointNexus()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointNexus(resourceData, ep, id)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   ep,
		EndpointId: ep.Id,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}
func TestServiceEndpointNexus_Update_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointNexus_Delete_DoesNotSwallowError(t, &nexusTestServiceEndpointPassword, nexusTestServiceEndpointProjectIDpassword)
}

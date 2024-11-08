//go:build (all || resource_serviceendpoint_incomingwebhook) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_incomingwebhook
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

var incomingWebhookTestServiceEndpointID = uuid.New()
var incomingWebhookRandomServiceEndpointProjectID = uuid.New()
var incomingWebhookTestServiceEndpointProjectID = &incomingWebhookRandomServiceEndpointProjectID

var incomingWebhookTestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"webhookname": "myTestWebhook",
			"secret":      "",
			"header":      "X-Test-Header",
		},
		Scheme: converter.String("None"),
	},
	Id:          &incomingWebhookTestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("incomingwebhook"),
	Url:         converter.String("https://dev.azure.com"),
	Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: incomingWebhookTestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func TestServiceEndpointIncomingWebhook_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointIncomingWebhook().Schema, nil)
	resourceData.Set("project_id", (*incomingWebhookTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointIncomingWebhook(resourceData, &incomingWebhookTestServiceEndpoint)

	serviceEndpointAfterRoundTrip, err := expandServiceEndpointIncomingWebhook(resourceData)

	require.Equal(t, incomingWebhookTestServiceEndpoint, *serviceEndpointAfterRoundTrip)
	require.Equal(t, incomingWebhookTestServiceEndpointProjectID, (*serviceEndpointAfterRoundTrip.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
	require.Nil(t, err)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestServiceEndpointIncomingWebhook_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointIncomingWebhook()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*incomingWebhookTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointIncomingWebhook(resourceData, &incomingWebhookTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &incomingWebhookTestServiceEndpoint}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointIncomingWebhook_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointIncomingWebhook()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*incomingWebhookTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointIncomingWebhook(resourceData, &incomingWebhookTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: incomingWebhookTestServiceEndpoint.Id,
		Project:    converter.String(incomingWebhookTestServiceEndpointProjectID.String()),
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
func TestServiceEndpointIncomingWebhook_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointIncomingWebhook()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*incomingWebhookTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointIncomingWebhook(resourceData, &incomingWebhookTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: incomingWebhookTestServiceEndpoint.Id,
		ProjectIds: &[]string{
			incomingWebhookTestServiceEndpointProjectID.String(),
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
func TestServiceEndpointIncomingWebhook_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointIncomingWebhook()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*incomingWebhookTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointIncomingWebhook(resourceData, &incomingWebhookTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &incomingWebhookTestServiceEndpoint,
		EndpointId: incomingWebhookTestServiceEndpoint.Id,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

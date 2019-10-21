package azuredevops

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"

	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
)

var testServiceEndpointID = uuid.New()
var randomServiceEndpointProjectID = uuid.New().String()
var testServiceEndpointProjectID = &randomServiceEndpointProjectID

var testServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"accessToken": "UNIT_TEST_ACCESS_TOKEN",
		},
		Scheme: converter.String("PersonalAccessToken"),
	},
	Id:    &testServiceEndpointID,
	Name:  converter.String("UNIT_TEST_NAME"),
	Owner: converter.String("library"), // Supported values are "library", "agentcloud"
	Type:  converter.String("UNIT_TEST_TYPE"),
	Url:   converter.String("UNIT_TEST_URL"),
}

/**
 * Begin unit tests
 */

// verifies that the flatten/expand round trip yields the same build definition
func TestAzureDevOpsServiceEndpoint_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceServiceEndpoint().Schema, nil)
	flattenServiceEndpoint(resourceData, &testServiceEndpoint, testServiceEndpointProjectID)

	serviceEndpointAfterRoundTrip, projectID := expandServiceEndpoint(resourceData)

	require.Equal(t, testServiceEndpoint, *serviceEndpointAfterRoundTrip)
	require.Equal(t, testServiceEndpointProjectID, projectID)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestAzureDevOpsServiceEndpoint_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := schema.TestResourceDataRaw(t, resourceServiceEndpoint().Schema, nil)
	flattenServiceEndpoint(resourceData, &testServiceEndpoint, testServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &aggregatedClient{ServiceEndpointClient: buildClient, ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &testServiceEndpoint, Project: testServiceEndpointProjectID}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := resourceServiceEndpointCreate(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestAzureDevOpsServiceEndpoint_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := schema.TestResourceDataRaw(t, resourceServiceEndpoint().Schema, nil)
	flattenServiceEndpoint(resourceData, &testServiceEndpoint, testServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &aggregatedClient{ServiceEndpointClient: buildClient, ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{EndpointId: testServiceEndpoint.Id, Project: testServiceEndpointProjectID}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := resourceServiceEndpointRead(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestAzureDevOpsServiceEndpoint_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := schema.TestResourceDataRaw(t, resourceServiceEndpoint().Schema, nil)
	flattenServiceEndpoint(resourceData, &testServiceEndpoint, testServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &aggregatedClient{ServiceEndpointClient: buildClient, ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{EndpointId: testServiceEndpoint.Id, Project: testServiceEndpointProjectID}
	buildClient.
		EXPECT().
		DeleteServiceEndpoint(clients.ctx, expectedArgs).
		Return(errors.New("DeleteServiceEndpoint() Failed")).
		Times(1)

	err := resourceServiceEndpointDelete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteServiceEndpoint() Failed")
}

// verifies that if an error is produced on an update, it is not swallowed
func TestAzureDevOpsServiceEndpoint_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := schema.TestResourceDataRaw(t, resourceServiceEndpoint().Schema, nil)
	flattenServiceEndpoint(resourceData, &testServiceEndpoint, testServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &aggregatedClient{ServiceEndpointClient: buildClient, ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &testServiceEndpoint,
		EndpointId: testServiceEndpoint.Id,
		Project:    testServiceEndpointProjectID,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := resourceServiceEndpointUpdate(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

// +build all resource_serviceendpoint_artifactory
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

var artifactoryTestServiceEndpointID = uuid.New()
var artifactoryRandomServiceEndpointProjectID = uuid.New().String()
var artifactoryTestServiceEndpointProjectID = &artifactoryRandomServiceEndpointProjectID

var artifactoryTestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": "AR_TEST_username",
			"password": "AR_TEST_password",
		},
		Scheme: converter.String("UsernamePassword"),
	},
	Id:          &artifactoryTestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("artifactoryService"),
	Url:         converter.String("https://www.artifactory.com"),
}
var artifactoryTestServiceEndpointToken = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": "AR_TEST_token",
		},
		Scheme: converter.String("Token"),
	},
	Id:          &artifactoryTestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("artifactoryService"),
	Url:         converter.String("https://www.artifactory.com"),
}

// verifies that the flatten/expand round trip yields the same service endpoint
func TestServiceEndpointArtifactory_ExpandFlatten_Roundtrip(t *testing.T) {
	for _, ep := range []*serviceendpoint.ServiceEndpoint{&artifactoryTestServiceEndpointToken, &artifactoryTestServiceEndpoint} {

		resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointArtifactory().Schema, nil)
		flattenServiceEndpointArtifactory(resourceData, ep, artifactoryTestServiceEndpointProjectID)

		serviceEndpointAfterRoundTrip, projectID, err := expandServiceEndpointArtifactory(resourceData)

		require.Equal(t, *ep, *serviceEndpointAfterRoundTrip)
		require.Equal(t, artifactoryTestServiceEndpointProjectID, projectID)
		require.Nil(t, err)
	}
}

// verifies that if an error is produced on create, the error is not swallowed
func TestServiceEndpointArtifactory_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointArtifactory()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointArtifactory(resourceData, &artifactoryTestServiceEndpoint, artifactoryTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &artifactoryTestServiceEndpoint, Project: artifactoryTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointArtifactory_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointArtifactory()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointArtifactory(resourceData, &artifactoryTestServiceEndpoint, artifactoryTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{EndpointId: artifactoryTestServiceEndpoint.Id, Project: artifactoryTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestServiceEndpointArtifactory_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointArtifactory()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointArtifactory(resourceData, &artifactoryTestServiceEndpoint, artifactoryTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{EndpointId: artifactoryTestServiceEndpoint.Id, Project: artifactoryTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		DeleteServiceEndpoint(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteServiceEndpoint() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteServiceEndpoint() Failed")
}

// verifies that if an error is produced on an update, it is not swallowed
func TestServiceEndpointArtifactory_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointArtifactory()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointArtifactory(resourceData, &artifactoryTestServiceEndpoint, artifactoryTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &artifactoryTestServiceEndpoint,
		EndpointId: artifactoryTestServiceEndpoint.Id,
		Project:    artifactoryTestServiceEndpointProjectID,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

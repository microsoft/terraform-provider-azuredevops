//go:build (all || resource_serviceendpoint_artifactory) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_artifactory
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

var artifactoryTestServiceEndpointIDpassword = uuid.New()
var artifactoryRandomServiceEndpointProjectIDpassword = uuid.New()
var artifactoryTestServiceEndpointProjectIDpassword = &artifactoryRandomServiceEndpointProjectIDpassword

var artifactoryTestServiceEndpointPassword = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"password": "",
			"username": "",
		},
		Scheme: converter.String("UsernamePassword"),
	},
	Id:    &artifactoryTestServiceEndpointIDpassword,
	Name:  converter.String("UNIT_TEST_CONN_NAME"),
	Owner: converter.String("library"), // Supported values are "library", "agentcloud"
	Type:  converter.String("artifactoryService"),
	Url:   converter.String("https://www.artifactory.com"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: artifactoryTestServiceEndpointProjectIDpassword,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

var artifactoryTestServiceEndpointID = uuid.New()
var artifactoryRandomServiceEndpointProjectID = uuid.New()
var artifactoryTestServiceEndpointProjectID = &artifactoryRandomServiceEndpointProjectID

var artifactoryTestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": "",
		},
		Scheme: converter.String("Token"),
	},
	Id:    &artifactoryTestServiceEndpointID,
	Name:  converter.String("UNIT_TEST_CONN_NAME"),
	Owner: converter.String("library"), // Supported values are "library", "agentcloud"
	Type:  converter.String("artifactoryService"),
	Url:   converter.String("https://www.artifactory.com"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: artifactoryTestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func testServiceEndpointArtifactory_ExpandFlatten_Roundtrip(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	for _, ep := range []*serviceendpoint.ServiceEndpoint{ep, ep} {

		resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointArtifactory().Schema, nil)
		flattenServiceEndpointArtifactory(resourceData, ep, id)

		serviceEndpointAfterRoundTrip, projectID, err := expandServiceEndpointArtifactory(resourceData)
		require.Nil(t, err)
		require.Equal(t, *ep, *serviceEndpointAfterRoundTrip)
		require.Equal(t, id, projectID)
	}
}

func TestServiceEndpointArtifactory_ExpandFlatten_RoundtripPassword(t *testing.T) {
	testServiceEndpointArtifactory_ExpandFlatten_Roundtrip(t, &artifactoryTestServiceEndpointPassword, artifactoryTestServiceEndpointProjectIDpassword)
}

func TestServiceEndpointArtifactory_ExpandFlatten_RoundtripToken(t *testing.T) {
	testServiceEndpointArtifactory_ExpandFlatten_Roundtrip(t, &artifactoryTestServiceEndpoint, artifactoryTestServiceEndpointProjectID)
}

// verifies that if an error is produced on create, the error is not swallowed
func testServiceEndpointArtifactory_Create_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointArtifactory()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointArtifactory(resourceData, ep, id)

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
func TestServiceEndpointArtifactory_Create_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointArtifactory_Create_DoesNotSwallowError(t, &artifactoryTestServiceEndpoint, artifactoryTestServiceEndpointProjectID)
}
func TestServiceEndpointArtifactory_Create_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointArtifactory_Create_DoesNotSwallowError(t, &artifactoryTestServiceEndpointPassword, artifactoryTestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on a read, it is not swallowed
func testServiceEndpointArtifactory_Read_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointArtifactory()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointArtifactory(resourceData, ep, id)

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
func TestServiceEndpointArtifactory_Read_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointArtifactory_Read_DoesNotSwallowError(t, &artifactoryTestServiceEndpoint, artifactoryTestServiceEndpointProjectID)
}
func TestServiceEndpointArtifactory_Read_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointArtifactory_Read_DoesNotSwallowError(t, &artifactoryTestServiceEndpointPassword, artifactoryTestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on a delete, it is not swallowed
func testServiceEndpointArtifactory_Delete_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointArtifactory()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointArtifactory(resourceData, ep, id)

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
func TestServiceEndpointArtifactory_Delete_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointArtifactory_Delete_DoesNotSwallowError(t, &artifactoryTestServiceEndpoint, artifactoryTestServiceEndpointProjectID)
}
func TestServiceEndpointArtifactory_Delete_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointArtifactory_Delete_DoesNotSwallowError(t, &artifactoryTestServiceEndpointPassword, artifactoryTestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on an update, it is not swallowed
func testServiceEndpointArtifactory_Update_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointArtifactory()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointArtifactory(resourceData, ep, id)

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
func TestServiceEndpointArtifactory_Update_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointArtifactory_Delete_DoesNotSwallowError(t, &artifactoryTestServiceEndpoint, artifactoryTestServiceEndpointProjectID)
}
func TestServiceEndpointArtifactory_Update_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointArtifactory_Delete_DoesNotSwallowError(t, &artifactoryTestServiceEndpointPassword, artifactoryTestServiceEndpointProjectIDpassword)
}

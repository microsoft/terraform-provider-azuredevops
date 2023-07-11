//go:build (all || resource_serviceendpoint_jfrog_artifactory_v2) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_jfrog_artifactory_v2
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

var artifactoryV2TestServiceEndpointIDpassword = uuid.New()
var artifactoryV2RandomServiceEndpointProjectIDpassword = uuid.New()
var artifactoryV2TestServiceEndpointProjectIDpassword = &artifactoryRandomServiceEndpointProjectIDpassword

var artifactoryV2TestServiceEndpointPassword = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": "",
			"password": "",
		},
		Scheme: converter.String("UsernamePassword"),
	},
	Id:          &artifactoryV2TestServiceEndpointIDpassword,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("jfrogArtifactoryService"),
	Url:         converter.String("https://www.artifactory.com"),
	Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: artifactoryV2TestServiceEndpointProjectIDpassword,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

var artifactoryV2TestServiceEndpointID = uuid.New()
var artifactoryV2RandomServiceEndpointProjectID = uuid.New()
var artifactoryV2TestServiceEndpointProjectID = &artifactoryRandomServiceEndpointProjectID

var artifactoryV2TestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": "",
		},
		Scheme: converter.String("Token"),
	},
	Id:          &artifactoryV2TestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("jfrogArtifactoryService"),
	Url:         converter.String("https://www.artifactory.com"),
	Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: artifactoryV2TestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func testServiceEndpointArtifactoryV2_ExpandFlatten_Roundtrip(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	for _, ep := range []*serviceendpoint.ServiceEndpoint{ep, ep} {

		resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointJFrogArtifactoryV2().Schema, nil)
		flattenServiceEndpointArtifactory(resourceData, ep, id)

		serviceEndpointAfterRoundTrip, projectID, err := expandServiceEndpointJFrogArtifactoryV2(resourceData)
		require.Nil(t, err)
		require.Equal(t, *ep, *serviceEndpointAfterRoundTrip)
		require.Equal(t, id, projectID)

	}
}
func TestServiceEndpointArtifactoryV2_ExpandFlatten_RoundtripPassword(t *testing.T) {
	testServiceEndpointArtifactoryV2_ExpandFlatten_Roundtrip(t, &artifactoryV2TestServiceEndpointPassword, artifactoryV2TestServiceEndpointProjectIDpassword)
}

func TestServiceEndpointArtifactoryV2_ExpandFlatten_RoundtripToken(t *testing.T) {
	testServiceEndpointArtifactoryV2_ExpandFlatten_Roundtrip(t, &artifactoryV2TestServiceEndpoint, artifactoryV2TestServiceEndpointProjectID)
}

// verifies that if an error is produced on create, the error is not swallowed
func testServiceEndpointArtifactoryV2_Create_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointJFrogArtifactoryV2()
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
func TestServiceEndpointArtifactoryV2_Create_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointArtifactoryV2_Create_DoesNotSwallowError(t, &artifactoryV2TestServiceEndpoint, artifactoryV2TestServiceEndpointProjectID)
}
func TestServiceEndpointArtifactoryV2_Create_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointArtifactoryV2_Create_DoesNotSwallowError(t, &artifactoryV2TestServiceEndpointPassword, artifactoryV2TestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on a read, it is not swallowed
func testServiceEndpointArtifactoryV2_Read_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointJFrogArtifactoryV2()
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
func TestServiceEndpointArtifactoryV2_Read_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointArtifactoryV2_Read_DoesNotSwallowError(t, &artifactoryV2TestServiceEndpoint, artifactoryV2TestServiceEndpointProjectID)
}
func TestServiceEndpointArtifactoryV2_Read_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointArtifactoryV2_Read_DoesNotSwallowError(t, &artifactoryV2TestServiceEndpointPassword, artifactoryV2TestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on a delete, it is not swallowed
func testServiceEndpointArtifactoryV2_Delete_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointJFrogArtifactoryV2()
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
func TestServiceEndpointArtifactoryV2_Delete_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointArtifactoryV2_Delete_DoesNotSwallowError(t, &artifactoryV2TestServiceEndpoint, artifactoryV2TestServiceEndpointProjectID)
}
func TestServiceEndpointArtifactoryV2_Delete_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointArtifactoryV2_Delete_DoesNotSwallowError(t, &artifactoryV2TestServiceEndpointPassword, artifactoryV2TestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on an update, it is not swallowed
func testServiceEndpointArtifactoryV2_Update_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointJFrogArtifactoryV2()
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
func TestServiceEndpointArtifactoryV2_Update_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointArtifactoryV2_Delete_DoesNotSwallowError(t, &artifactoryV2TestServiceEndpoint, artifactoryV2TestServiceEndpointProjectID)
}
func TestServiceEndpointArtifactoryV2_Update_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointArtifactoryV2_Delete_DoesNotSwallowError(t, &artifactoryV2TestServiceEndpointPassword, artifactoryV2TestServiceEndpointProjectIDpassword)
}

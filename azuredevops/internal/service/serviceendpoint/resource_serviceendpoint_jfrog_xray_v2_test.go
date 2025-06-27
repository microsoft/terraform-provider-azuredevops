//go:build (all || resource_serviceendpoint_jfrog_xray_v2) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_jfrog_xray_v2
// +build !exclude_serviceendpoints

package serviceendpoint

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	xrayV2TestServiceEndpointIDpassword          = uuid.New()
	xrayV2RandomServiceEndpointProjectIDpassword = uuid.New()
	xrayV2TestServiceEndpointProjectIDpassword   = &artifactoryRandomServiceEndpointProjectIDpassword
)

var xrayV2TestServiceEndpointPassword = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": "",
			"password": "",
		},
		Scheme: converter.String("UsernamePassword"),
	},
	Id:          &xrayV2TestServiceEndpointIDpassword,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("jfrogXrayService"),
	Url:         converter.String("https://www.artifactory.com"),
	Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: xrayV2TestServiceEndpointProjectIDpassword,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

var (
	xrayV2TestServiceEndpointID          = uuid.New()
	xrayV2RandomServiceEndpointProjectID = uuid.New()
	xrayV2TestServiceEndpointProjectID   = &artifactoryRandomServiceEndpointProjectID
)

var xrayV2TestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": "",
		},
		Scheme: converter.String("Token"),
	},
	Id:          &xrayV2TestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("jfrogXrayService"),
	Url:         converter.String("https://www.artifactory.com"),
	Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: xrayV2TestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func testServiceEndpointxrayV2_ExpandFlatten_Roundtrip(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	for _, ep := range []*serviceendpoint.ServiceEndpoint{ep, ep} {

		resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointJFrogXRayV2().Schema, nil)
		resourceData.Set("project_id", (*ep.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
		flattenServiceEndpointArtifactory(resourceData, ep)

		serviceEndpointAfterRoundTrip, err := expandServiceEndpointJFrogXRayV2(resourceData)
		require.Nil(t, err)
		require.Equal(t, *ep, *serviceEndpointAfterRoundTrip)
		require.Equal(t, id, (*serviceEndpointAfterRoundTrip.ServiceEndpointProjectReferences)[0].ProjectReference.Id)

	}
}

func TestServiceEndpointxrayV2_ExpandFlatten_RoundtripPassword(t *testing.T) {
	testServiceEndpointxrayV2_ExpandFlatten_Roundtrip(t, &xrayV2TestServiceEndpointPassword, xrayV2TestServiceEndpointProjectIDpassword)
}

func TestServiceEndpointxrayV2_ExpandFlatten_RoundtripToken(t *testing.T) {
	testServiceEndpointxrayV2_ExpandFlatten_Roundtrip(t, &xrayV2TestServiceEndpoint, xrayV2TestServiceEndpointProjectID)
}

// verifies that if an error is produced on create, the error is not swallowed
func testServiceEndpointxrayV2_Create_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointJFrogXRayV2()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*ep.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointArtifactory(resourceData, ep)

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

func TestServiceEndpointxrayV2_Create_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointxrayV2_Create_DoesNotSwallowError(t, &xrayV2TestServiceEndpoint, xrayV2TestServiceEndpointProjectID)
}

func TestServiceEndpointxrayV2_Create_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointxrayV2_Create_DoesNotSwallowError(t, &xrayV2TestServiceEndpointPassword, xrayV2TestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on a read, it is not swallowed
func testServiceEndpointxrayV2_Read_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointJFrogXRayV2()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*ep.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointArtifactory(resourceData, ep)

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

func TestServiceEndpointxrayV2_Read_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointxrayV2_Read_DoesNotSwallowError(t, &xrayV2TestServiceEndpoint, xrayV2TestServiceEndpointProjectID)
}

func TestServiceEndpointxrayV2_Read_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointxrayV2_Read_DoesNotSwallowError(t, &xrayV2TestServiceEndpointPassword, xrayV2TestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on a delete, it is not swallowed
func testServiceEndpointxrayV2_Delete_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointJFrogXRayV2()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*ep.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointArtifactory(resourceData, ep)

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

func TestServiceEndpointxrayV2_Delete_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointxrayV2_Delete_DoesNotSwallowError(t, &xrayV2TestServiceEndpoint, xrayV2TestServiceEndpointProjectID)
}

func TestServiceEndpointxrayV2_Delete_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointxrayV2_Delete_DoesNotSwallowError(t, &xrayV2TestServiceEndpointPassword, xrayV2TestServiceEndpointProjectIDpassword)
}

func TestServiceEndpointxrayV2_Update_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointxrayV2_Delete_DoesNotSwallowError(t, &xrayV2TestServiceEndpoint, xrayV2TestServiceEndpointProjectID)
}

func TestServiceEndpointxrayV2_Update_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointxrayV2_Delete_DoesNotSwallowError(t, &xrayV2TestServiceEndpointPassword, xrayV2TestServiceEndpointProjectIDpassword)
}

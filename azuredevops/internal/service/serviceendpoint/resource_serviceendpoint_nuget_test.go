//go:build (all || resource_serviceendpoint_nuget) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_nuget
// +build !exclude_serviceendpoints

package serviceendpoint

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

var nugetTestServiceEndpointIDpassword = uuid.New()
var nugetRandomServiceEndpointProjectIDpassword = uuid.New()
var nugetTestServiceEndpointProjectIDpassword = &nugetRandomServiceEndpointProjectID

var nugetTestServiceEndpointPassword = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": "",
			"password": "",
		},
		Scheme: converter.String("UsernamePassword"),
	},
	Id:    &nugetTestServiceEndpointIDpassword,
	Name:  converter.String("UNIT_TEST_CONN_NAME"),
	Owner: converter.String("library"), // Supported values are "library", "agentcloud"
	Type:  converter.String("externalnugetfeed"),
	Url:   converter.String("https://api.nuget.org/v3/index.json"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: nugetTestServiceEndpointProjectIDpassword,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

var nugetTestServiceEndpointID = uuid.New()
var nugetRandomServiceEndpointProjectID = uuid.New()
var nugetTestServiceEndpointProjectID = &nugetRandomServiceEndpointProjectID

var nugetTestServiceEndpointToken = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": "",
		},
		Scheme: converter.String("Token"),
	},
	Id:    &nugetTestServiceEndpointID,
	Name:  converter.String("UNIT_TEST_CONN_NAME"),
	Owner: converter.String("library"), // Supported values are "library", "agentcloud"
	Type:  converter.String("externalnugetfeed"),
	Url:   converter.String("https://api.nuget.org/v3/index.json"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: nugetTestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

var nugetTestServiceEndpointKey = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"nugetkey": "",
		},
		Scheme: converter.String("None"),
	},
	Id:    &nugetTestServiceEndpointID,
	Name:  converter.String("UNIT_TEST_CONN_NAME"),
	Owner: converter.String("library"), // Supported values are "library", "agentcloud"
	Type:  converter.String("externalnugetfeed"),
	Url:   converter.String("https://api.nuget.org/v3/index.json"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: nugetTestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func testServiceEndpointNuget_ExpandFlatten_Roundtrip(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	for _, ep := range []*serviceendpoint.ServiceEndpoint{ep, ep} {

		resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointNuget().Schema, nil)
		flattenServiceEndpointNuget(resourceData, ep, id)

		serviceEndpointAfterRoundTrip, projectID, err := expandServiceEndpointNuget(resourceData)
		require.Nil(t, err)
		require.Equal(t, *ep, *serviceEndpointAfterRoundTrip)
		require.Equal(t, id, projectID)

	}
}
func TestServiceEndpointNuget_ExpandFlatten_RoundtripPassword(t *testing.T) {
	testServiceEndpointNuget_ExpandFlatten_Roundtrip(t, &nugetTestServiceEndpointPassword, nugetTestServiceEndpointProjectIDpassword)
}
func TestServiceEndpointNuget_ExpandFlatten_RoundtripToken(t *testing.T) {
	testServiceEndpointNuget_ExpandFlatten_Roundtrip(t, &nugetTestServiceEndpointToken, nugetTestServiceEndpointProjectID)
}
func TestServiceEndpointNuget_ExpandFlatten_RoundtripKey(t *testing.T) {
	testServiceEndpointNuget_ExpandFlatten_Roundtrip(t, &nugetTestServiceEndpointKey, nugetTestServiceEndpointProjectID)
}

// verifies that if an error is produced on create, the error is not swallowed
func testServiceEndpointNuget_Create_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointNuget()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointNuget(resourceData, ep, id)

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
func TestServiceEndpointNuget_Create_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointNuget_Create_DoesNotSwallowError(t, &nugetTestServiceEndpointToken, nugetTestServiceEndpointProjectID)
}
func TestServiceEndpointNuget_Create_DoesNotSwallowErrorKey(t *testing.T) {
	testServiceEndpointNuget_Create_DoesNotSwallowError(t, &nugetTestServiceEndpointKey, nugetTestServiceEndpointProjectID)
}
func TestServiceEndpointNuget_Create_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointNuget_Create_DoesNotSwallowError(t, &nugetTestServiceEndpointPassword, nugetTestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on a read, it is not swallowed
func testServiceEndpointNuget_Read_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointNuget()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointNuget(resourceData, ep, id)

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
func TestServiceEndpointNuget_Read_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointNuget_Read_DoesNotSwallowError(t, &nugetTestServiceEndpointToken, nugetTestServiceEndpointProjectID)
}
func TestServiceEndpointNuget_Read_DoesNotSwallowErrorKey(t *testing.T) {
	testServiceEndpointNuget_Read_DoesNotSwallowError(t, &nugetTestServiceEndpointKey, nugetTestServiceEndpointProjectID)
}
func TestServiceEndpointNuget_Read_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointNuget_Read_DoesNotSwallowError(t, &nugetTestServiceEndpointPassword, nugetTestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on a delete, it is not swallowed
func testServiceEndpointNuget_Delete_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointNuget()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointNuget(resourceData, ep, id)

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
func TestServiceEndpointNuget_Delete_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointNuget_Delete_DoesNotSwallowError(t, &nugetTestServiceEndpointToken, nugetTestServiceEndpointProjectID)
}
func TestServiceEndpointNuget_Delete_DoesNotSwallowErrorKey(t *testing.T) {
	testServiceEndpointNuget_Delete_DoesNotSwallowError(t, &nugetTestServiceEndpointKey, nugetTestServiceEndpointProjectID)
}
func TestServiceEndpointNuget_Delete_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointNuget_Delete_DoesNotSwallowError(t, &nugetTestServiceEndpointPassword, nugetTestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on an update, it is not swallowed
func testServiceEndpointNuget_Update_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointNuget()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointNuget(resourceData, ep, id)

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
func TestServiceEndpointNuget_Update_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointNuget_Delete_DoesNotSwallowError(t, &nugetTestServiceEndpointToken, nugetTestServiceEndpointProjectID)
}
func TestServiceEndpointNuget_Update_DoesNotSwallowErrorKey(t *testing.T) {
	testServiceEndpointNuget_Delete_DoesNotSwallowError(t, &nugetTestServiceEndpointKey, nugetTestServiceEndpointProjectID)
}
func TestServiceEndpointNuget_Update_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointNuget_Delete_DoesNotSwallowError(t, &nugetTestServiceEndpointPassword, nugetTestServiceEndpointProjectIDpassword)
}

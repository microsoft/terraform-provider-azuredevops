//go:build (all || resource_serviceendpoint_dockerregistry) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_dockerregistry
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

var dockerRegistryTestServiceEndpointID = uuid.New()
var dockerRegistryRandomServiceEndpointProjectID = uuid.New()
var dockerRegistryTestServiceEndpointProjectID = &dockerRegistryRandomServiceEndpointProjectID

var dockerRegistryTestServiceEndpoint = serviceendpoint.ServiceEndpoint{ //todo change
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": "DH_TEST_username",
			"password": "",
			"email":    "DH_TEST_email",
			"registry": "https://index.docker.io/v1/",
		},
		Scheme: converter.String("UsernamePassword"),
	},
	Data: &map[string]string{
		"registrytype": "Others",
	},
	Id:          &dockerRegistryTestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("dockerregistry"),
	Url:         converter.String("https://hub.docker.com/"),
	Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: dockerRegistryTestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func TestServiceEndpointDockerRegistry_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointDockerRegistry().Schema, nil)
	resourceData.Set("project_id", (*dockerRegistryTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointDockerRegistry(resourceData, &dockerRegistryTestServiceEndpoint)

	serviceEndpointAfterRoundTrip, err := expandServiceEndpointDockerRegistry(resourceData)

	require.Equal(t, dockerRegistryTestServiceEndpoint, *serviceEndpointAfterRoundTrip)
	require.Equal(t, dockerRegistryTestServiceEndpointProjectID, (*serviceEndpointAfterRoundTrip.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
	require.Nil(t, err)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestServiceEndpointDockerRegistry_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointDockerRegistry()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*dockerRegistryTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointDockerRegistry(resourceData, &dockerRegistryTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &dockerRegistryTestServiceEndpoint}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointDockerRegistry_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointDockerRegistry()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*dockerRegistryTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointDockerRegistry(resourceData, &dockerRegistryTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: dockerRegistryTestServiceEndpoint.Id,
		Project:    converter.String(dockerRegistryTestServiceEndpointProjectID.String()),
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
func TestServiceEndpointDockerRegistry_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointDockerRegistry()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*dockerRegistryTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointDockerRegistry(resourceData, &dockerRegistryTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: dockerRegistryTestServiceEndpoint.Id,
		ProjectIds: &[]string{
			dockerRegistryTestServiceEndpointProjectID.String(),
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
func TestServiceEndpointDockerRegistry_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointDockerRegistry()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*dockerRegistryTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointDockerRegistry(resourceData, &dockerRegistryTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &dockerRegistryTestServiceEndpoint,
		EndpointId: dockerRegistryTestServiceEndpoint.Id,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

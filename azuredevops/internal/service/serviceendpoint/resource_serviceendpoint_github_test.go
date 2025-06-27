//go:build (all || resource_serviceendpoint_github) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_github
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
	ghTestServiceEndpointID          = uuid.New()
	ghRandomServiceEndpointProjectID = uuid.New()
	ghTestServiceEndpointProjectID   = &ghRandomServiceEndpointProjectID
)

var ghTestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"AccessToken": "UNIT_TEST_ACCESS_TOKEN",
		},
		Scheme: converter.String("Token"),
	},
	Id:          &ghTestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_NAME"),
	Owner:       converter.String("library"),
	Type:        converter.String("github"),
	Url:         converter.String("https://github.com"),
	Description: converter.String("UNIT_TEST_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: ghTestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_NAME"),
			Description: converter.String("UNIT_TEST_DESCRIPTION"),
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func TestServiceEndpointGitHub_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointGitHub().Schema, nil)
	resourceData.Set("project_id", (*ghTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	configureAuthPersonal(resourceData)
	flattenServiceEndpointGitHub(resourceData, &ghTestServiceEndpoint)

	serviceEndpointAfterRoundTrip, err := expandServiceEndpointGitHub(resourceData)

	require.Nil(t, err)
	require.Equal(t, ghTestServiceEndpoint, *serviceEndpointAfterRoundTrip)
	require.Equal(t, ghTestServiceEndpointProjectID, (*serviceEndpointAfterRoundTrip.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestServiceEndpointGitHub_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointGitHub()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*ghTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	configureAuthPersonal(resourceData)
	flattenServiceEndpointGitHub(resourceData, &ghTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &ghTestServiceEndpoint}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointGitHub_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointGitHub()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*ghTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointGitHub(resourceData, &ghTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: ghTestServiceEndpoint.Id,
		Project:    converter.String(ghTestServiceEndpointProjectID.String()),
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
func TestServiceEndpointGitHub_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointGitHub()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*ghTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointGitHub(resourceData, &ghTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: ghTestServiceEndpoint.Id,
		ProjectIds: &[]string{
			ghTestServiceEndpointProjectID.String(),
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
func TestServiceEndpointGitHub_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointGitHub()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*ghTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	configureAuthPersonal(resourceData)
	flattenServiceEndpointGitHub(resourceData, &ghTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &ghTestServiceEndpoint,
		EndpointId: ghTestServiceEndpoint.Id,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

func configureAuthPersonal(d *schema.ResourceData) {
	d.Set("auth_personal", &[]map[string]interface{}{
		{
			"personal_access_token": "UNIT_TEST_ACCESS_TOKEN",
		},
	})
}

//go:build (all || resource_serviceendpoint_externaltfs) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_externaltfs
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

var (
	externalTfsTestServiceEndpointID          = uuid.New()
	externalTfsRandomServiceEndpointProjectID = uuid.New()
	externalTfsTestServiceEndpointProjectID   = &externalTfsRandomServiceEndpointProjectID
)

var externalTfsTestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": "UNIT_TEST_ACCESS_TOKEN",
		},
		Scheme: converter.String("Token"),
	},
	Id:          &externalTfsTestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_NAME"),
	Owner:       converter.String("library"),
	Type:        converter.String("externaltfs"),
	Url:         converter.String("https://dev.azure.com/myorganization"),
	Description: converter.String("UNIT_TEST_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: externalTfsTestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_NAME"),
			Description: converter.String("UNIT_TEST_DESCRIPTION"),
		},
	},
}

func TestServiceEndpointExternalTFS_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointExternalTFS().Schema, nil)
	configureExternalTfsAuthPersonal(resourceData)
	flattenServiceEndpointExternalTFS(
		resourceData,
		&externalTfsTestServiceEndpoint,
		externalTfsTestServiceEndpointProjectID.String(),
	)

	serviceEndpointAfterRoundTrip, projectID, err := expandServiceEndpointExternalTFS(resourceData)

	require.Nil(t, err)
	require.Equal(t, externalTfsTestServiceEndpoint, *serviceEndpointAfterRoundTrip)
	require.Equal(t, externalTfsTestServiceEndpointProjectID.String(), projectID)
}

func TestServiceEndpointExternalTFS_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointExternalTFS()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureExternalTfsAuthPersonal(resourceData)
	flattenServiceEndpointExternalTFS(
		resourceData,
		&externalTfsTestServiceEndpoint,
		externalTfsTestServiceEndpointProjectID.String(),
	)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &externalTfsTestServiceEndpoint}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

func TestServiceEndpointExternalTFS_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointExternalTFS()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointExternalTFS(
		resourceData,
		&externalTfsTestServiceEndpoint,
		externalTfsTestServiceEndpointProjectID.String(),
	)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: externalTfsTestServiceEndpoint.Id,
		Project:    converter.String(externalTfsTestServiceEndpointProjectID.String()),
	}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
}

func TestServiceEndpointExternalTFS_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointExternalTFS()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointExternalTFS(
		resourceData,
		&externalTfsTestServiceEndpoint,
		externalTfsTestServiceEndpointProjectID.String(),
	)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: externalTfsTestServiceEndpoint.Id,
		ProjectIds: &[]string{
			externalTfsTestServiceEndpointProjectID.String(),
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

func TestServiceEndpointExternalTFS_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointExternalTFS()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureExternalTfsAuthPersonal(resourceData)
	flattenServiceEndpointExternalTFS(
		resourceData,
		&externalTfsTestServiceEndpoint,
		externalTfsTestServiceEndpointProjectID.String(),
	)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &externalTfsTestServiceEndpoint,
		EndpointId: externalTfsTestServiceEndpoint.Id,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

func configureExternalTfsAuthPersonal(d *schema.ResourceData) {
	d.Set("auth_personal", &[]map[string]interface{}{
		{
			personalAccessTokenExternalTFS: "UNIT_TEST_ACCESS_TOKEN",
		},
	})
}

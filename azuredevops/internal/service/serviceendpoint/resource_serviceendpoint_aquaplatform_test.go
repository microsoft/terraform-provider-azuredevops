//go:build (all || resource_serviceendpoint_aquaplatform) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_aquaplatform
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
	aquaPlatformTestServiceEndpointID          = uuid.New()
	aquaPlatformRandomServiceEndpointProjectID = uuid.New()
	aquaPlatformTestServiceEndpointProjectID   = &aquaPlatformRandomServiceEndpointProjectID
)

var aquaPlatformTestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"aquaKey":    "key",
			"aquaSecret": "secret",
		},
		Scheme: converter.String("None"),
	},
	Data: &map[string]string{
		"aquaPlatformUrl": "https://aqua.com/",
		"aquaAuthUrl":     "https://api.cloudsploit.com",
	},
	Id:          &aquaPlatformTestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Owner:       converter.String("library"),
	Type:        converter.String("aquaplatform"),
	Url:         converter.String("https://aqua.com/"),
	Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: aquaPlatformTestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

func TestServiceEndpointAquaPlatform_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointAquaPlatform().Schema, nil)
	resourceData.Set("project_id", (*aquaPlatformTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointAquaPlatform(resourceData, &aquaPlatformTestServiceEndpoint)

	// Since aqua_key and aqua_secret are not flattened, we need to set them manually
	resourceData.Set("aqua_key", "key")
	resourceData.Set("aqua_secret", "secret")

	serviceEndpointAfterRoundTrip := expandServiceEndpointAquaPlatform(resourceData)

	require.Equal(t, aquaPlatformTestServiceEndpoint, *serviceEndpointAfterRoundTrip)
	require.Equal(t, aquaPlatformTestServiceEndpointProjectID, (*serviceEndpointAfterRoundTrip.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
}

func TestServiceEndpointAquaPlatform_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAquaPlatform()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*aquaPlatformTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointAquaPlatform(resourceData, &aquaPlatformTestServiceEndpoint)
	resourceData.Set("aqua_key", "key")
	resourceData.Set("aqua_secret", "secret")

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &aquaPlatformTestServiceEndpoint}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

func TestServiceEndpointAquaPlatform_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAquaPlatform()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*aquaPlatformTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointAquaPlatform(resourceData, &aquaPlatformTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: aquaPlatformTestServiceEndpoint.Id,
		Project:    converter.String(aquaPlatformTestServiceEndpointProjectID.String()),
	}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
}

func TestServiceEndpointAquaPlatform_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAquaPlatform()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*aquaPlatformTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointAquaPlatform(resourceData, &aquaPlatformTestServiceEndpoint)
	resourceData.Set("aqua_key", "key")
	resourceData.Set("aqua_secret", "secret")

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: aquaPlatformTestServiceEndpoint.Id,
		ProjectIds: &[]string{
			aquaPlatformTestServiceEndpointProjectID.String(),
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

func TestServiceEndpointAquaPlatform_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAquaPlatform()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*aquaPlatformTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointAquaPlatform(resourceData, &aquaPlatformTestServiceEndpoint)
	resourceData.Set("aqua_key", "key")
	resourceData.Set("aqua_secret", "secret")

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &aquaPlatformTestServiceEndpoint,
		EndpointId: aquaPlatformTestServiceEndpoint.Id,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

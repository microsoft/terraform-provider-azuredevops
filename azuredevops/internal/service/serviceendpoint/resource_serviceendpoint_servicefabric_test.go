//go:build (all || resource_serviceendpoint_service_faric) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_service_faric
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

var serviceFabricTestServiceEndpointID = uuid.New()
var serviceFabricRandomServiceEndpointProjectID = uuid.New()
var serviceFabricTestServiceEndpointProjectID = &serviceFabricRandomServiceEndpointProjectID

var serviceFabricTestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"certLookup":           "Thumbprint",
			"servercertthumbprint": "THUMBPRINT_TEST",
			"certificate":          "CERTIFICATE_TEST",
			"certificatepassword":  "CERTIFICATE_PASSWORD_TEST",
		},
		Scheme: converter.String("Certificate"),
	},
	Id:          &serviceFabricTestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_NAME"),
	Owner:       converter.String("library"),
	Type:        converter.String("servicefabric"),
	Url:         converter.String("tcp://servicefabric.com"),
	Description: converter.String("UNIT_TEST_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: serviceFabricTestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_NAME"),
			Description: converter.String("UNIT_TEST_DESCRIPTION"),
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func TestServiceEndpointServiceFabric_FlattenExpand_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointServiceFabric().Schema, nil)
	resourceData.Set("project_id", (*serviceFabricTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	configureAuthServiceFabricCertificate(resourceData)
	flattenServiceEndpointServiceFabric(resourceData, &serviceFabricTestServiceEndpoint)

	serviceEndpointAfterRoundTrip, err := expandServiceEndpointServiceFabric(resourceData)

	require.Nil(t, err)
	require.Equal(t, serviceFabricTestServiceEndpoint, *serviceEndpointAfterRoundTrip)
	require.Equal(t, serviceFabricTestServiceEndpointProjectID, (*serviceEndpointAfterRoundTrip.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestServiceEndpointServiceFabric_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointServiceFabric()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*serviceFabricTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	configureAuthServiceFabricCertificate(resourceData)
	flattenServiceEndpointServiceFabric(resourceData, &serviceFabricTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &serviceFabricTestServiceEndpoint}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointServiceFabric_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointServiceFabric()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*serviceFabricTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	configureAuthServiceFabricCertificate(resourceData)
	flattenServiceEndpointServiceFabric(resourceData, &serviceFabricTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: serviceFabricTestServiceEndpoint.Id,
		Project:    converter.String(serviceFabricTestServiceEndpointProjectID.String()),
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
func TestServiceEndpointServiceFabric_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointServiceFabric()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*serviceFabricTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	configureAuthServiceFabricCertificate(resourceData)
	flattenServiceEndpointServiceFabric(resourceData, &serviceFabricTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: serviceFabricTestServiceEndpoint.Id,
		ProjectIds: &[]string{
			serviceFabricTestServiceEndpointProjectID.String(),
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
func TestServiceEndpointServiceFabric_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointServiceFabric()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*serviceFabricTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	configureAuthServiceFabricCertificate(resourceData)
	flattenServiceEndpointServiceFabric(resourceData, &serviceFabricTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &serviceFabricTestServiceEndpoint,
		EndpointId: serviceFabricTestServiceEndpoint.Id,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

func configureAuthServiceFabricCertificate(d *schema.ResourceData) {
	d.Set("certificate", &[]map[string]interface{}{
		{
			"server_certificate_lookup":     "Thumbprint",
			"server_certificate_thumbprint": "THUMBPRINT_TEST",
			"client_certificate":            "CERTIFICATE_TEST",
			"client_certificate_password":   "CERTIFICATE_PASSWORD_TEST",
		},
	})
}

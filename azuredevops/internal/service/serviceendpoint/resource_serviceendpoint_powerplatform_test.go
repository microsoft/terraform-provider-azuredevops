//go:build (all || resource_serviceendpoint_powerplatform) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_powerplatform
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
	powerplatformTestServiceEndpointID          = uuid.New()
	powerplatformRandomServiceEndpointProjectID = uuid.New()
	powerplatformTestServiceEndpointProjectID   = &powerplatformRandomServiceEndpointProjectID
)

// getTestServiceEndpointPowerPlatform returns a valid PowerPlatform service endpoint struct
// matching the "powerplatform-spn" type and "None" scheme structure.
func getTestServiceEndpointPowerPlatform() serviceendpoint.ServiceEndpoint {
	return serviceendpoint.ServiceEndpoint{
		Authorization: &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"tenantId":      "aba07645-051c-44b4-b806-c34d33f3dcd1", // fake tenant
				"applicationId": "e31eaaac-47da-4156-b433-9b0538c94b7e", // fake app ID
				"clientSecret":  "supersecretkey",                       // fake secret
			},
			Scheme: converter.String("None"),
		},
		Data: &map[string]string{}, // Empty data as per requirement
		Id:   &powerplatformTestServiceEndpointID,
		Name: converter.String("_POWERPLATFORM_UNIT_TEST_CONN_NAME"),
		Type: converter.String("powerplatform-spn"),
		Url:  converter.String("https://org.crm.dynamics.com/"),
		ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
			{
				ProjectReference: &serviceendpoint.ProjectReference{
					Id: powerplatformTestServiceEndpointProjectID,
				},
				Name:        converter.String("_POWERPLATFORM_UNIT_TEST_CONN_NAME"),
				Description: converter.String("_POWERPLATFORM_UNIT_TEST_CONN_DESCRIPTION"),
			},
		},
	}
}

var powerplatformTestServiceEndpoints = []serviceendpoint.ServiceEndpoint{
	getTestServiceEndpointPowerPlatform(),
}

// verifies that the flatten/expand round trip yields the same service endpoint
func TestServiceEndpointPowerPlatform_ExpandFlatten_Roundtrip(t *testing.T) {
	for _, resource := range powerplatformTestServiceEndpoints {
		resourceData := getResourceDataPowerPlatform(t, resource)

		flattenServiceEndpointPowerPlatform(resourceData, &resource)

		serviceEndpointAfterRoundTrip, err := expandServiceEndpointPowerPlatform(resourceData)
		require.Nil(t, err)

		require.Equal(t, *resource.Authorization.Parameters, *serviceEndpointAfterRoundTrip.Authorization.Parameters)
		require.Equal(t, resource.Url, serviceEndpointAfterRoundTrip.Url)
		require.Equal(t, resource.Type, serviceEndpointAfterRoundTrip.Type)
		require.Equal(t, powerplatformTestServiceEndpointProjectID, (*serviceEndpointAfterRoundTrip.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
	}
}

// verifies that if an error is produced on create, the error is not swallowed
func TestServiceEndpointPowerPlatform_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointPowerPlatform()
	for _, resource := range powerplatformTestServiceEndpoints {
		resourceData := getResourceDataPowerPlatform(t, resource)
		flattenServiceEndpointPowerPlatform(resourceData, &resource)

		buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
		clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

		expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &resource}

		buildClient.
			EXPECT().
			CreateServiceEndpoint(clients.Ctx, expectedArgs).
			Return(nil, errors.New("CreateServiceEndpoint() Failed")).
			Times(1)

		err := r.Create(resourceData, clients)
		require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
	}
}

// verifies that if validation is enabled and fails, the error is returned and endpoint deleted
func TestServiceEndpointPowerPlatform_CreateWithValidate_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointPowerPlatform()
	for _, resource := range powerplatformTestServiceEndpoints {
		resourceData := getResourceDataPowerPlatform(t, resource)
		flattenServiceEndpointPowerPlatform(resourceData, &resource)

		buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
		clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

		buildClient.
			EXPECT().
			CreateServiceEndpoint(clients.Ctx, serviceendpoint.CreateServiceEndpointArgs{Endpoint: &resource}).
			Return(&resource, nil).
			Times(1)

		returnedServiceEndpoint := resource
		returnedServiceEndpoint.IsReady = converter.Bool(true)
		buildClient.
			EXPECT().
			GetServiceEndpointDetails(clients.Ctx, serviceendpoint.GetServiceEndpointDetailsArgs{
				Project:    converter.String(powerplatformRandomServiceEndpointProjectID.String()),
				EndpointId: resource.Id,
			},
			).
			Return(&returnedServiceEndpoint, nil).
			Times(1)

		reqArgs := genExecuteServiceEndpointArgsPowerPlatform(&resource)
		buildClient.
			EXPECT().
			ExecuteServiceEndpointRequest(clients.Ctx, *reqArgs).
			Return(nil, errors.New("ExecuteServiceEndpointRequest() Failed")).
			Times(1)

		buildClient.
			EXPECT().
			DeleteServiceEndpoint(clients.Ctx, serviceendpoint.DeleteServiceEndpointArgs{
				ProjectIds: &[]string{powerplatformTestServiceEndpointProjectID.String()}, EndpointId: resource.Id,
			}).
			Return(nil).
			Times(1)

		err := r.Create(resourceData, clients)
		require.Contains(t, err.Error(), "ExecuteServiceEndpointRequest() Failed")
	}
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointPowerPlatform_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointPowerPlatform()
	for _, resource := range powerplatformTestServiceEndpoints {
		resourceData := getResourceDataPowerPlatform(t, resource)
		flattenServiceEndpointPowerPlatform(resourceData, &resource)

		buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
		clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

		expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
			EndpointId: resource.Id,
			Project:    converter.String(powerplatformTestServiceEndpointProjectID.String()),
		}

		buildClient.
			EXPECT().
			GetServiceEndpointDetails(clients.Ctx, expectedArgs).
			Return(nil, errors.New("GetServiceEndpoint() Failed")).
			Times(1)

		err := r.Read(resourceData, clients)
		require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
	}
}

// verifies that if an error is produced on an update, it is not swallowed
func TestServiceEndpointPowerPlatform_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointPowerPlatform()
	for _, resource := range powerplatformTestServiceEndpoints {
		resourceData := getResourceDataPowerPlatform(t, resource)
		flattenServiceEndpointPowerPlatform(resourceData, &resource)

		buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
		clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

		expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
			Endpoint:   &resource,
			EndpointId: resource.Id,
		}

		buildClient.
			EXPECT().
			UpdateServiceEndpoint(clients.Ctx, expectedArgs).
			Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
			Times(1)

		err := r.Update(resourceData, clients)
		require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
	}
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestServiceEndpointPowerPlatform_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointPowerPlatform()
	for _, resource := range powerplatformTestServiceEndpoints {
		resourceData := getResourceDataPowerPlatform(t, resource)
		flattenServiceEndpointPowerPlatform(resourceData, &resource)

		buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
		clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

		expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
			EndpointId: resource.Id,
			ProjectIds: &[]string{
				powerplatformTestServiceEndpointProjectID.String(),
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
}

// Helper to create ResourceData with the correct fields for PowerPlatform
func getResourceDataPowerPlatform(t *testing.T, resource serviceendpoint.ServiceEndpoint) *schema.ResourceData {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointPowerPlatform().Schema, nil)

	resourceData.Set("project_id", (*resource.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())

	if resource.Url != nil {
		resourceData.Set("url", *resource.Url)
	}

	params := *resource.Authorization.Parameters
	credentials := []interface{}{
		map[string]interface{}{
			"serviceprincipalid":  params["applicationId"],
			"serviceprincipalkey": params["clientSecret"],
			"tenantId":            params["tenantId"],
		},
	}
	resourceData.Set("credentials", credentials)

	return resourceData
}

// Helper to generate execution args for validation mocks
func genExecuteServiceEndpointArgsPowerPlatform(endpoint *serviceendpoint.ServiceEndpoint) *serviceendpoint.ExecuteServiceEndpointRequestArgs {
	return &serviceendpoint.ExecuteServiceEndpointRequestArgs{
		ServiceEndpointRequest: &serviceendpoint.ServiceEndpointRequest{
			DataSourceDetails: &serviceendpoint.DataSourceDetails{
				DataSourceName: converter.String("TestConnection"),
			},
			ResultTransformationDetails: &serviceendpoint.ResultTransformationDetails{},
			ServiceEndpointDetails: &serviceendpoint.ServiceEndpointDetails{
				Data:          endpoint.Data,
				Authorization: endpoint.Authorization,
				Url:           endpoint.Url,
				Type:          endpoint.Type,
			},
		},
		Project:    converter.String((*endpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String()),
		EndpointId: converter.String(endpoint.Id.String()),
	}
}

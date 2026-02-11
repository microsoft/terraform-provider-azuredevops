package serviceendpoint

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/testhelper"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var serviceEndPointWellFormed = serviceendpoint.ServiceEndpoint{
	Id:   testhelper.CreateUUID(),
	Name: converter.String("Endpoint1"),
	Url:  converter.String("https://generic.com"),
	Type: converter.String("Generic"),
	Data: &map[string]string{
		"someKey": "someValue",
	},
	Authorization: &serviceendpoint.EndpointAuthorization{
		Scheme: converter.String("UsernamePassword"),
		Parameters: &map[string]string{
			"username": "user",
		},
	},
}

var projectID = "123e4567-e89b-12d3-a456-426655440000"

func TestDataEndpointGeneric_Read_TestFindEndpointByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceEndpointClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{
		ServiceEndpointClient: serviceEndpointClient,
		Ctx:                   context.Background(),
	}

	expectedGetByNameArgs := serviceendpoint.GetServiceEndpointsByNamesArgs{
		Project:       converter.String(projectID),
		EndpointNames: &[]string{*serviceEndPointWellFormed.Name},
	}

	serviceEndpointClient.
		EXPECT().
		GetServiceEndpointsByNames(clients.Ctx, expectedGetByNameArgs).
		Return(&[]serviceendpoint.ServiceEndpoint{serviceEndPointWellFormed}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataServiceEndpointGeneric().Schema, nil)
	resourceData.Set("project_id", projectID)
	resourceData.Set("service_endpoint_name", serviceEndPointWellFormed.Name)

	err := dataSourceServiceEndpointGenericRead(resourceData, clients)

	require.Nil(t, err)
	require.Equal(t, *serviceEndPointWellFormed.Name, resourceData.Get("service_endpoint_name").(string))
	require.Equal(t, "someValue", resourceData.Get("data").(map[string]interface{})["someKey"])
	require.Equal(t, "user", resourceData.Get("authorization").(map[string]interface{})["username"])
	require.Equal(t, "UsernamePassword", resourceData.Get("authorization").(map[string]interface{})["scheme"])
}

func TestDataEndpointGeneric_Read_TestEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceEndpointClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{
		ServiceEndpointClient: serviceEndpointClient,
		Ctx:                   context.Background(),
	}

	expectedGetByNameArgs := serviceendpoint.GetServiceEndpointsByNamesArgs{
		Project:       converter.String(projectID),
		EndpointNames: &[]string{*serviceEndPointWellFormed.Name},
	}

	serviceEndpointClient.
		EXPECT().
		GetServiceEndpointsByNames(clients.Ctx, expectedGetByNameArgs).
		Return(&[]serviceendpoint.ServiceEndpoint{serviceEndPointWellFormed}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataServiceEndpointGeneric().Schema, nil)

	resourceData.Set("project_id", projectID)

	resourceData.Set("service_endpoint_name", *serviceEndPointWellFormed.Name)

	err := dataSourceServiceEndpointGenericRead(resourceData, clients)
	require.Nil(t, err)

	require.Equal(t, *serviceEndPointWellFormed.Name, resourceData.Get("service_endpoint_name").(string))
	require.Equal(t, "someValue", resourceData.Get("data").(map[string]interface{})["someKey"])
	require.Equal(t, "user", resourceData.Get("authorization").(map[string]interface{})["username"])
	require.Equal(t, "UsernamePassword", resourceData.Get("authorization").(map[string]interface{})["scheme"])
}

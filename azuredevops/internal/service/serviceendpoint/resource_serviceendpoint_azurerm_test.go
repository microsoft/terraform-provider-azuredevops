//go:build (all || resource_serviceendpoint_azurerm) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_azurerm
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

var azurermTestServiceEndpointAzureRMID = uuid.New()
var azurermRandomServiceEndpointAzureRMProjectID = uuid.New()
var azurermTestServiceEndpointAzureRMProjectID = &azurermRandomServiceEndpointAzureRMProjectID

func getManualAuthServiceEndpoint() serviceendpoint.ServiceEndpoint {
	return serviceendpoint.ServiceEndpoint{
		Authorization: &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"authenticationType":  "spnKey",
				"serviceprincipalid":  "e31eaaac-47da-4156-b433-9b0538c94b7e", //fake value
				"serviceprincipalkey": "serviceprincipalkey",                  //fake value
				"tenantid":            "aba07645-051c-44b4-b806-c34d33f3dcd1", //fake value
			},
			Scheme: converter.String("ServicePrincipal"),
		},
		Data: &map[string]string{
			"creationMode":     "Manual",
			"environment":      "AzureCloud",
			"scopeLevel":       "Subscription",
			"subscriptionId":   "42125daf-72fd-417c-9ea7-080690625ad3", //fake value
			"subscriptionName": "SUBSCRIPTION_TEST",
		},
		Id:          &azurermTestServiceEndpointAzureRMID,
		Name:        converter.String("_AZURERM_UNIT_TEST_CONN_NAME"),
		Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
		Type:        converter.String("azurerm"),
		Url:         converter.String("https://management.azure.com/"),
		Description: converter.String("_AZURERM_UNIT_TEST_CONN_DESCRIPTION"),
		ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
			{
				ProjectReference: &serviceendpoint.ProjectReference{
					Id: azurermTestServiceEndpointAzureRMProjectID,
				},
				Name:        converter.String("_AZURERM_UNIT_TEST_CONN_NAME"),
				Description: converter.String("_AZURERM_UNIT_TEST_CONN_DESCRIPTION"),
			},
		},
	}
}

var azurermTestServiceEndpointsAzureRM = []serviceendpoint.ServiceEndpoint{
	getManualAuthServiceEndpoint(),
	{
		Authorization: &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"authenticationType":  "spnKey",
				"serviceprincipalid":  "",
				"serviceprincipalkey": "",
				"tenantid":            "aba07645-051c-44b4-b806-c34d33f3dcd1", //fake value
			},
			Scheme: converter.String("ServicePrincipal"),
		},
		Data: &map[string]string{
			"creationMode":     "Automatic",
			"environment":      "AzureCloud",
			"scopeLevel":       "Subscription",
			"subscriptionId":   "42125daf-72fd-417c-9ea7-080690625ad3", //fake value
			"subscriptionName": "SUBSCRIPTION_TEST",
		},
		Id:          &azurermTestServiceEndpointAzureRMID,
		Name:        converter.String("_AZURERM_UNIT_TEST_CONN_NAME"),
		Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
		Type:        converter.String("azurerm"),
		Url:         converter.String("https://management.azure.com/"),
		Description: converter.String("_AZURERM_UNIT_TEST_CONN_DESCRIPTION"),
		ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
			{
				ProjectReference: &serviceendpoint.ProjectReference{
					Id: azurermTestServiceEndpointAzureRMProjectID,
				},
				Name:        converter.String("_AZURERM_UNIT_TEST_CONN_NAME"),
				Description: converter.String("_AZURERM_UNIT_TEST_CONN_DESCRIPTION"),
			},
		},
	},
	{
		Authorization: &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"authenticationType":  "spnKey",
				"serviceprincipalid":  "",
				"serviceprincipalkey": "",
				"tenantid":            "aba07645-051c-44b4-b806-c34d33f3dcd1", //fake value
				"scope":               "/subscriptions/42125daf-72fd-417c-9ea7-080690625ad3/resourcegroups/test",
			},
			Scheme: converter.String("ServicePrincipal"),
		},
		Data: &map[string]string{
			"creationMode":     "Automatic",
			"environment":      "AzureCloud",
			"scopeLevel":       "Subscription",
			"subscriptionId":   "42125daf-72fd-417c-9ea7-080690625ad3", //fake value
			"subscriptionName": "SUBSCRIPTION_TEST",
		},
		Id:          &azurermTestServiceEndpointAzureRMID,
		Name:        converter.String("_AZURERM_UNIT_TEST_CONN_NAME"),
		Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
		Type:        converter.String("azurerm"),
		Url:         converter.String("https://management.azure.com/"),
		Description: converter.String("_AZURERM_UNIT_TEST_CONN_DESCRIPTION"),
		ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
			{
				ProjectReference: &serviceendpoint.ProjectReference{
					Id: azurermTestServiceEndpointAzureRMProjectID,
				},
				Name:        converter.String("_AZURERM_UNIT_TEST_CONN_NAME"),
				Description: converter.String("_AZURERM_UNIT_TEST_CONN_DESCRIPTION"),
			},
		},
	},
	{
		Authorization: &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"tenantid": "aba07645-051c-44b4-b806-c34d33f3dcd1", //fake value
			},
			Scheme: converter.String("ManagedServiceIdentity"),
		},
		Data: &map[string]string{
			"environment":      "AzureCloud",
			"scopeLevel":       "Subscription",
			"subscriptionId":   "42125daf-72fd-417c-9ea7-080690625ad3", //fake value
			"subscriptionName": "SUBSCRIPTION_TEST",
		},
		Id:          &azurermTestServiceEndpointAzureRMID,
		Name:        converter.String("_AZURERM_UNIT_TEST_CONN_NAME"),
		Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
		Type:        converter.String("azurerm"),
		Url:         converter.String("https://management.azure.com/"),
		Description: converter.String("_AZURERM_UNIT_TEST_CONN_DESCRIPTION"),
		ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
			{
				ProjectReference: &serviceendpoint.ProjectReference{
					Id: azurermTestServiceEndpointAzureRMProjectID,
				},
				Name:        converter.String("_AZURERM_UNIT_TEST_CONN_NAME"),
				Description: converter.String("_AZURERM_UNIT_TEST_CONN_DESCRIPTION"),
			},
		},
	},
	{
		Authorization: &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"tenantid":           "aba07645-051c-44b4-b806-c34d33f3dcd1", //fake value
				"serviceprincipalid": "bba07645-051c-44b4-b806-c34d33f3dcd2", //fake value
			},
			Scheme: converter.String("WorkloadIdentityFederation"),
		},
		Data: &map[string]string{
			"environment":      "AzureCloud",
			"scopeLevel":       "Subscription",
			"subscriptionId":   "42125daf-72fd-417c-9ea7-080690625ad3", //fake value
			"subscriptionName": "SUBSCRIPTION_TEST",
			"creationMode":     "Manual",
		},
		Id:          &azurermTestServiceEndpointAzureRMID,
		Name:        converter.String("_AZURERM_UNIT_TEST_CONN_NAME"),
		Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
		Type:        converter.String("azurerm"),
		Url:         converter.String("https://management.azure.com/"),
		Description: converter.String("_AZURERM_UNIT_TEST_CONN_DESCRIPTION"),
		ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
			{
				ProjectReference: &serviceendpoint.ProjectReference{
					Id: azurermTestServiceEndpointAzureRMProjectID,
				},
				Name:        converter.String("_AZURERM_UNIT_TEST_CONN_NAME"),
				Description: converter.String("_AZURERM_UNIT_TEST_CONN_DESCRIPTION"),
			},
		},
	},
	{
		Authorization: &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"tenantid":           "aba07645-051c-44b4-b806-c34d33f3dcd1", //fake value
				"serviceprincipalid": "",
			},
			Scheme: converter.String("WorkloadIdentityFederation"),
		},
		Data: &map[string]string{
			"environment":      "AzureCloud",
			"scopeLevel":       "Subscription",
			"subscriptionId":   "42125daf-72fd-417c-9ea7-080690625ad3", //fake value
			"subscriptionName": "SUBSCRIPTION_TEST",
			"creationMode":     "Automatic",
		},
		Id:          &azurermTestServiceEndpointAzureRMID,
		Name:        converter.String("_AZURERM_UNIT_TEST_CONN_NAME"),
		Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
		Type:        converter.String("azurerm"),
		Url:         converter.String("https://management.azure.com/"),
		Description: converter.String("_AZURERM_UNIT_TEST_CONN_DESCRIPTION"),
		ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
			{
				ProjectReference: &serviceendpoint.ProjectReference{
					Id: azurermTestServiceEndpointAzureRMProjectID,
				},
				Name:        converter.String("_AZURERM_UNIT_TEST_CONN_NAME"),
				Description: converter.String("_AZURERM_UNIT_TEST_CONN_DESCRIPTION"),
			},
		},
	},
	{
		Authorization: &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"tenantid":           "aba07645-051c-44b4-b806-c34d33f3dcd1", //fake value
				"serviceprincipalid": "",
				"scope":              "/subscriptions/42125daf-72fd-417c-9ea7-080690625ad3/resourcegroups/test",
			},
			Scheme: converter.String("WorkloadIdentityFederation"),
		},
		Data: &map[string]string{
			"environment":      "AzureCloud",
			"scopeLevel":       "Subscription",
			"subscriptionId":   "42125daf-72fd-417c-9ea7-080690625ad3", //fake value
			"subscriptionName": "SUBSCRIPTION_TEST",
			"creationMode":     "Automatic",
		},
		Id:          &azurermTestServiceEndpointAzureRMID,
		Name:        converter.String("_AZURERM_UNIT_TEST_CONN_NAME"),
		Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
		Type:        converter.String("azurerm"),
		Url:         converter.String("https://management.azure.com/"),
		Description: converter.String("_AZURERM_UNIT_TEST_CONN_DESCRIPTION"),
		ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
			{
				ProjectReference: &serviceendpoint.ProjectReference{
					Id: azurermTestServiceEndpointAzureRMProjectID,
				},
				Name:        converter.String("_AZURERM_UNIT_TEST_CONN_NAME"),
				Description: converter.String("_AZURERM_UNIT_TEST_CONN_DESCRIPTION"),
			},
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func TestServiceEndpointAzureRM_ExpandFlatten_Roundtrip(t *testing.T) {
	for _, resource := range azurermTestServiceEndpointsAzureRM {
		resourceData := getResourceData(t, resource)
		resourceData.Set("project_id", (*resource.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
		if (*resource.Data)["creationMode"] == "Manual" {
			credentials := []interface{}{
				map[string]string{
					"serviceprincipalid":          (*resource.Authorization.Parameters)["serviceprincipalid"],
					"serviceprincipalkey":         (*resource.Authorization.Parameters)["serviceprincipalkey"],
					"serviceprincipalcertificate": (*resource.Authorization.Parameters)[""],
				},
			}
			resourceData.Set("credentials", credentials)
		}
		flattenServiceEndpointAzureRM(resourceData, &resource)
		serviceEndpointAfterRoundTrip, _ := expandServiceEndpointAzureRM(resourceData)

		require.Equal(t, resource, *serviceEndpointAfterRoundTrip)
		require.Equal(t, azurermTestServiceEndpointAzureRMProjectID, (*serviceEndpointAfterRoundTrip.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
	}
}

// verifies that if an error is produced on create, the error is not swallowed
func TestServiceEndpointAzureRM_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAzureRM()
	for _, resource := range azurermTestServiceEndpointsAzureRM {
		resourceData := getResourceData(t, resource)
		resourceData.Set("project_id", (*resource.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
		flattenServiceEndpointAzureRM(resourceData, &resource)

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

func TestServiceEndpointAzureRM_CreateWithValidate_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAzureRM()
	for _, resource := range azurermTestServiceEndpointsAzureRM {
		resourceData := getResourceData(t, resource)
		resourceData.Set("project_id", (*resource.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
		if (*resource.Data)["creationMode"] == "Manual" {
			credentials := []interface{}{
				map[string]string{
					"serviceprincipalid":          (*resource.Authorization.Parameters)["serviceprincipalid"],
					"serviceprincipalkey":         (*resource.Authorization.Parameters)["serviceprincipalkey"],
					"serviceprincipalcertificate": (*resource.Authorization.Parameters)[""],
				},
			}
			resourceData.Set("credentials", credentials)
		}
		flattenServiceEndpointAzureRM(resourceData, &resource)

		features := initializeFeaturesWithValidate(true)
		resourceData.Set("features", features)

		buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
		clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

		createArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &resource}

		buildClient.
			EXPECT().
			CreateServiceEndpoint(clients.Ctx, createArgs).
			Return(&resource, nil).
			Times(1)

		returnedServiceEndpoint := resource
		returnedServiceEndpoint.IsReady = converter.Bool(true)
		buildClient.
			EXPECT().
			GetServiceEndpointDetails(clients.Ctx, serviceendpoint.GetServiceEndpointDetailsArgs{
				Project:    converter.String(azurermRandomServiceEndpointAzureRMProjectID.String()),
				EndpointId: resource.Id,
			},
			).
			Return(&returnedServiceEndpoint, nil).
			Times(1)

		reqArgs := genExecuteServiceEndpointArgs(&resource)
		buildClient.
			EXPECT().
			ExecuteServiceEndpointRequest(clients.Ctx, *reqArgs).
			Return(nil, errors.New("ExecuteServiceEndpointRequest() Failed")).
			Times(1)

		buildClient.
			EXPECT().
			DeleteServiceEndpoint(clients.Ctx, serviceendpoint.DeleteServiceEndpointArgs{
				ProjectIds: &[]string{azurermTestServiceEndpointAzureRMProjectID.String()}, EndpointId: resource.Id}).
			Return(nil).
			Times(1)

		err := r.Create(resourceData, clients)
		require.Contains(t, err.Error(), "ExecuteServiceEndpointRequest() Failed")
	}
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointAzureRM_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAzureRM()
	for _, resource := range azurermTestServiceEndpointsAzureRM {
		resourceData := getResourceData(t, resource)
		resourceData.Set("project_id", (*resource.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
		flattenServiceEndpointAzureRM(resourceData, &resource)

		buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
		clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

		expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
			EndpointId: resource.Id,
			Project:    converter.String(azurermTestServiceEndpointAzureRMProjectID.String()),
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

// verifies that if an error is produced on a delete, it is not swallowed
func TestServiceEndpointAzureRM_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAzureRM()
	for _, resource := range azurermTestServiceEndpointsAzureRM {
		resourceData := getResourceData(t, resource)
		resourceData.Set("project_id", (*resource.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
		flattenServiceEndpointAzureRM(resourceData, &resource)

		buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
		clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

		expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
			EndpointId: resource.Id,
			ProjectIds: &[]string{
				azurermTestServiceEndpointAzureRMProjectID.String(),
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

// verifies that if an error is produced on an update, it is not swallowed
func TestServiceEndpointAzureRM_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAzureRM()
	for _, resource := range azurermTestServiceEndpointsAzureRM {
		resourceData := getResourceData(t, resource)
		resourceData.Set("project_id", (*resource.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
		flattenServiceEndpointAzureRM(resourceData, &resource)

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

func TestServiceEndpointAzureRM_UpdateWithValidate_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointAzureRM()
	for _, resource := range azurermTestServiceEndpointsAzureRM {
		resourceData := getResourceData(t, resource)
		resourceData.Set("project_id", (*resource.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
		if (*resource.Data)["creationMode"] == "Manual" {
			credentials := []interface{}{
				map[string]string{
					"serviceprincipalid":          (*resource.Authorization.Parameters)["serviceprincipalid"],
					"serviceprincipalkey":         (*resource.Authorization.Parameters)["serviceprincipalkey"],
					"serviceprincipalcertificate": (*resource.Authorization.Parameters)[""],
				},
			}
			resourceData.Set("credentials", credentials)
		}
		flattenServiceEndpointAzureRM(resourceData, &resource)

		features := initializeFeaturesWithValidate(true)
		resourceData.Set("features", features)

		buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
		clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

		reqArgs := genExecuteServiceEndpointArgs(&resource)
		buildClient.
			EXPECT().
			ExecuteServiceEndpointRequest(clients.Ctx, *reqArgs).
			Return(nil, errors.New("ExecuteServiceEndpointRequest() Failed")).
			Times(1)

		err := r.Update(resourceData, clients)
		require.Contains(t, err.Error(), "ExecuteServiceEndpointRequest() Failed")
	}
}

// This is a little different than most. The steps done, along with the motivation behind each, are as follows:
//	(1) The service endpoint is configured. The `serviceprincipalkey` is set to `""`, which matches
//		the Azure DevOps API behavior. The service will intentionally hide the value of
//		`serviceprincipalkey` because it is a secret value
//	(2) The resource is flattened/expanded
//	(3) The `serviceprincipalkey` field is inspected and asserted to equal `"null"`. This special
//		value, which is unfortunately not documented in the REST API, will be interpreted by the
//		Azure DevOps API as an indicator to "not update" the field. The resulting behavior is that
//		this Terraform Resource will be able to update the Service Endpoint without needing to
//		pass the password along in each request.
//func TestServiceEndpointAzureRM_ExpandHandlesMissingSpnKeyInAPIResponse(t *testing.T) {
//	// step (1)
//	endpoint := getManualAuthServiceEndpoint()
//	resourceData := getResourceData(t, endpoint)
//	(*endpoint.Authorization.Parameters)["serviceprincipalkey"] = ""
//
//	// step (2)
//	flattenServiceEndpointAzureRM(resourceData, &endpoint, azurermTestServiceEndpointAzureRMProjectID)
//	expandedEndpoint, _, _ := expandServiceEndpointAzureRM(resourceData)
//
//	// step (3)
//	spnKeyProperty := (*expandedEndpoint.Authorization.Parameters)["serviceprincipalkey"]
//	require.Equal(t, "null", spnKeyProperty)
//}

func getResourceData(t *testing.T, resource serviceendpoint.ServiceEndpoint) *schema.ResourceData {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointAzureRM().Schema, nil)
	if key := (*resource.Authorization.Parameters)["serviceprincipalkey"]; key != "" {
		resourceData.Set("credentials", []map[string]interface{}{{
			"serviceprincipalid":  (*resource.Authorization.Parameters)["serviceprincipalid"],
			"serviceprincipalkey": (*resource.Authorization.Parameters)["serviceprincipalkey"],
		}})
	}
	return resourceData
}

func genExecuteServiceEndpointArgs(endpoint *serviceendpoint.ServiceEndpoint) *serviceendpoint.ExecuteServiceEndpointRequestArgs {
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

func initializeFeaturesWithValidate(validate bool) []map[string]interface{} {
	var features []map[string]interface{}
	feature := make(map[string]interface{})
	feature["validate"] = validate
	features = append(features, feature)
	return features
}

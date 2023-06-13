//go:build (all || resource_serviceendpoint_gcp_terraform) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_gcp_terraform
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

var gcpForTerraformTestServiceEndpointID = uuid.New()
var gcpRandomServiceEndpointProjectID = uuid.New()
var gcpForTerraformTestServiceEndpointProjectID = &gcpRandomServiceEndpointProjectID

var gcpForTerraformTestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"Issuer":     "GCP_TEST_client_email",
			"Audience":   "GCP_TEST_token_uri",
			"Scope":      "GCP_TEST_scope",
			"PrivateKey": "",
		},
		Scheme: converter.String("JWT"),
	},
	Data: &map[string]string{
		"project": "GCP_TEST_project_id",
	},
	Id:    &gcpForTerraformTestServiceEndpointID,
	Name:  converter.String("UNIT_TEST_CONN_NAME"),
	Owner: converter.String("library"), // Supported values are "library", "agentcloud"
	Type:  converter.String("GoogleCloudServiceEndpoint"),
	Url:   converter.String("https://www.googleapis.com/"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: gcpForTerraformTestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func TestServiceEndpointGcp_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointGcp().Schema, nil)
	flattenServiceEndpointGcp(resourceData, &gcpForTerraformTestServiceEndpoint, gcpForTerraformTestServiceEndpointProjectID)

	serviceEndpointAfterRoundTrip, projectID, err := expandServiceEndpointGcp(resourceData)

	require.Equal(t, gcpForTerraformTestServiceEndpoint, *serviceEndpointAfterRoundTrip)
	require.Equal(t, gcpForTerraformTestServiceEndpointProjectID, projectID)
	require.Nil(t, err)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestServiceEndpointGcp_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointGcp()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointGcp(resourceData, &gcpForTerraformTestServiceEndpoint, gcpForTerraformTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &gcpForTerraformTestServiceEndpoint}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointGcp_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointGcp()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointGcp(resourceData, &gcpForTerraformTestServiceEndpoint, gcpForTerraformTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: gcpForTerraformTestServiceEndpoint.Id,
		Project:    converter.String(gcpForTerraformTestServiceEndpointProjectID.String()),
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
func TestServiceEndpointGcp_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointGcp()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointGcp(resourceData, &gcpForTerraformTestServiceEndpoint, gcpForTerraformTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: gcpForTerraformTestServiceEndpoint.Id,
		ProjectIds: &[]string{
			gcpForTerraformTestServiceEndpointProjectID.String(),
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
func TestServiceEndpointGcp_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointGcp()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointGcp(resourceData, &gcpForTerraformTestServiceEndpoint, gcpForTerraformTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &gcpForTerraformTestServiceEndpoint,
		EndpointId: gcpForTerraformTestServiceEndpoint.Id,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

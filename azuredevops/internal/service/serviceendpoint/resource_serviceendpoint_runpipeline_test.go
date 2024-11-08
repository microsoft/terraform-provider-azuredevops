//go:build (all || resource_serviceendpoint_rpipeline) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_rpipeline
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

// "rp" stands for RunPipeline
var rpTestServiceEndpointID = uuid.New()
var rpRandomServiceEndpointProjectID = uuid.New()
var rpTestServiceEndpointProjectID = &rpRandomServiceEndpointProjectID

var rpTestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": "UNIT_TEST_ACCESS_TOKEN",
		},
		Scheme: converter.String("Token"),
	},
	Data: &map[string]string{
		"releaseUrl": "https://vsrm.dev.azure.com/example",
	},
	Id:          &rpTestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_NAME"),
	Owner:       converter.String("library"),
	Type:        converter.String("azdoapi"),
	Url:         converter.String("https://dev.azure.com/example"),
	Description: converter.String("UNIT_TEST_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: rpTestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_NAME"),
			Description: converter.String("UNIT_TEST_DESCRIPTION"),
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func TestServiceEndpointRunPipeline_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointRunPipeline().Schema, nil)
	resourceData.Set("project_id", (*rpTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	rpConfigureExtraFields(resourceData)
	flattenServiceEndpointRunPipeline(resourceData, &rpTestServiceEndpoint)

	serviceEndpointAfterRoundTrip, err := expandServiceEndpointRunPipeline(resourceData)

	require.Nil(t, err)
	require.Equal(t, rpTestServiceEndpoint, *serviceEndpointAfterRoundTrip)
	require.Equal(t, rpTestServiceEndpointProjectID, (*serviceEndpointAfterRoundTrip.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestServiceEndpointRunPipeline_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointRunPipeline()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*rpTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	rpConfigureExtraFields(resourceData)
	flattenServiceEndpointRunPipeline(resourceData, &rpTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &rpTestServiceEndpoint}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointRunPipeline_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointRunPipeline()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*rpTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointRunPipeline(resourceData, &rpTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: rpTestServiceEndpoint.Id,
		Project:    converter.String(rpTestServiceEndpointProjectID.String()),
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
func TestServiceEndpointRunPipeline_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointRunPipeline()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*rpTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointRunPipeline(resourceData, &rpTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: rpTestServiceEndpoint.Id,
		ProjectIds: &[]string{
			rpTestServiceEndpointProjectID.String(),
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
func TestServiceEndpointRunPipeline_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointRunPipeline()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*rpTestServiceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	rpConfigureExtraFields(resourceData)
	flattenServiceEndpointRunPipeline(resourceData, &rpTestServiceEndpoint)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &rpTestServiceEndpoint,
		EndpointId: rpTestServiceEndpoint.Id,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

func rpConfigureExtraFields(d *schema.ResourceData) {
	d.Set("organization_name", "example")
	d.Set("auth_personal", &[]map[string]interface{}{
		{
			"personal_access_token": "UNIT_TEST_ACCESS_TOKEN",
		},
	})
}

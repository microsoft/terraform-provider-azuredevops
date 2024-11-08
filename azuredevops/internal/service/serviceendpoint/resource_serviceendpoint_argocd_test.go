//go:build (all || resource_serviceendpoint_argocd) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_argocd
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

var argocdTestServiceEndpointIDpassword = uuid.New()
var argocdRandomServiceEndpointProjectIDpassword = uuid.New()
var argocdTestServiceEndpointProjectIDpassword = &argocdRandomServiceEndpointProjectIDpassword

var argocdTestServiceEndpointPassword = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": "",
			"password": "",
		},
		Scheme: converter.String("UsernamePassword"),
	},
	Id:          &argocdTestServiceEndpointIDpassword,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("argocd"),
	Url:         converter.String("https://www.argocd.com"),
	Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: argocdTestServiceEndpointProjectIDpassword,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

var argocdTestServiceEndpointID = uuid.New()
var argocdRandomServiceEndpointProjectID = uuid.New()
var argocdTestServiceEndpointProjectID = &argocdRandomServiceEndpointProjectID

var argocdTestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": "",
		},
		Scheme: converter.String("Token"),
	},
	Id:          &argocdTestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("argocd"),
	Url:         converter.String("https://www.argocd.com"),
	Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: argocdTestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func testServiceEndpointArgoCD_ExpandFlatten_Roundtrip(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	for _, ep := range []*serviceendpoint.ServiceEndpoint{ep, ep} {

		resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointArgoCD().Schema, nil)
		resourceData.Set("project_id", (*(*ep.ServiceEndpointProjectReferences)[0].ProjectReference.Id).String())
		flattenServiceEndpointArgoCD(resourceData, ep)

		serviceEndpointAfterRoundTrip, err := expandServiceEndpointArgoCD(resourceData)
		require.Nil(t, err)
		require.Equal(t, *ep, *serviceEndpointAfterRoundTrip)
		require.Equal(t, id, (*serviceEndpointAfterRoundTrip.ServiceEndpointProjectReferences)[0].ProjectReference.Id)

	}
}
func TestServiceEndpointArgoCD_ExpandFlatten_RoundtripPassword(t *testing.T) {
	testServiceEndpointArgoCD_ExpandFlatten_Roundtrip(t, &argocdTestServiceEndpointPassword, argocdTestServiceEndpointProjectIDpassword)
}

func TestServiceEndpointArgoCD_ExpandFlatten_RoundtripToken(t *testing.T) {
	testServiceEndpointArgoCD_ExpandFlatten_Roundtrip(t, &argocdTestServiceEndpoint, argocdTestServiceEndpointProjectID)
}

// verifies that if an error is produced on create, the error is not swallowed
func testServiceEndpointArgoCD_Create_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointArgoCD()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*ep.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointArgoCD(resourceData, ep)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: ep}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}
func TestServiceEndpointArgoCD_Create_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointArgoCD_Create_DoesNotSwallowError(t, &argocdTestServiceEndpoint, argocdTestServiceEndpointProjectID)
}
func TestServiceEndpointArgoCD_Create_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointArgoCD_Create_DoesNotSwallowError(t, &argocdTestServiceEndpointPassword, argocdTestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on a read, it is not swallowed
func testServiceEndpointArgoCD_Read_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointArgoCD()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", id.String())
	flattenServiceEndpointArgoCD(resourceData, ep)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: ep.Id,
		Project:    converter.String(id.String()),
	}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
}
func TestServiceEndpointArgoCD_Read_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointArgoCD_Read_DoesNotSwallowError(t, &argocdTestServiceEndpoint, argocdTestServiceEndpointProjectID)
}
func TestServiceEndpointArgoCD_Read_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointArgoCD_Read_DoesNotSwallowError(t, &argocdTestServiceEndpointPassword, argocdTestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on a delete, it is not swallowed
func testServiceEndpointArgoCD_Delete_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointArgoCD()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", id.String())
	flattenServiceEndpointArgoCD(resourceData, ep)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: ep.Id,
		ProjectIds: &[]string{
			id.String(),
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

func TestServiceEndpointArgoCD_Delete_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointArgoCD_Delete_DoesNotSwallowError(t, &argocdTestServiceEndpoint, argocdTestServiceEndpointProjectID)
}

func TestServiceEndpointArgoCD_Delete_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointArgoCD_Delete_DoesNotSwallowError(t, &argocdTestServiceEndpointPassword, argocdTestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on an update, it is not swallowed
func testServiceEndpointArgoCD_Update_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointArgoCD()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", id.String())
	flattenServiceEndpointArgoCD(resourceData, ep)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   ep,
		EndpointId: ep.Id,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

func TestServiceEndpointArgoCD_Update_DoesNotSwallowErrorToken(t *testing.T) {
	testServiceEndpointArgoCD_Delete_DoesNotSwallowError(t, &argocdTestServiceEndpoint, argocdTestServiceEndpointProjectID)
}

func TestServiceEndpointArgoCD_Update_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointArgoCD_Delete_DoesNotSwallowError(t, &argocdTestServiceEndpointPassword, argocdTestServiceEndpointProjectIDpassword)
}

func TestServiceEndpointArgoCD_Update_DoesNotSwallowError(t *testing.T) {
	testServiceEndpointArgoCD_Update_DoesNotSwallowError(t, &argocdTestServiceEndpoint, argocdTestServiceEndpointProjectID)
}

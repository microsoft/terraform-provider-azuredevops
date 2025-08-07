//go:build (all || resource_serviceendpoint_jenkins) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_jenkins
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
	jenkinsTestServiceEndpointIDPassword          = uuid.New()
	jenkinsRandomServiceEndpointProjectIDPassword = uuid.New()
	jenkinsTestServiceEndpointProjectIDPassword   = &jenkinsRandomServiceEndpointProjectIDPassword
)

var jenkinsTestServiceEndpointPassword = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": "JENKINS_TEST_username",
			"password": "",
		},
		Scheme: converter.String("UsernamePassword"),
	},
	Data: &map[string]string{
		"AcceptUntrustedCerts": "false",
	},
	Id:          &jenkinsTestServiceEndpointIDPassword,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("Jenkins"),
	Url:         converter.String("https://www.jenkins.com"),
	Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: jenkinsTestServiceEndpointProjectIDPassword,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func testServiceEndpointJenkins_ExpandFlatten_Roundtrip(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	for _, se := range []*serviceendpoint.ServiceEndpoint{ep, ep} {
		resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointJenkins().Schema, nil)
		resourceData.Set("project_id", (*se.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
		flattenServiceEndpointJenkins(resourceData, se)

		serviceEndpointAfterRoundTrip, err := expandServiceEndpointJenkins(resourceData)

		require.Nil(t, err)
		require.Equal(t, *se, *serviceEndpointAfterRoundTrip)
		require.Equal(t, id, (*serviceEndpointAfterRoundTrip.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
	}
}

func TestServiceEndpointJenkins_ExpandFlatten_RoundtripPassword(t *testing.T) {
	testServiceEndpointJenkins_ExpandFlatten_Roundtrip(t, &jenkinsTestServiceEndpointPassword, jenkinsTestServiceEndpointProjectIDPassword)
}

func TestServiceEndpointJenkins_Create_DoesNotSwallowErrorPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointJenkins()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*jenkinsTestServiceEndpointPassword.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointJenkins(resourceData, &jenkinsTestServiceEndpointPassword)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	seJenkins, _ := expandServiceEndpointJenkins(resourceData)
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, serviceendpoint.CreateServiceEndpointArgs{Endpoint: seJenkins}).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on read, the error is not swallowed
func testServiceEndpointJenkins_Read_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointJenkins()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*ep.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointJenkins(resourceData, ep)

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

func TestServiceEndpointJenkins_Read_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointJenkins_Read_DoesNotSwallowError(t, &jenkinsTestServiceEndpointPassword, jenkinsTestServiceEndpointProjectIDPassword)
}

// verifies that if an error is produced on a delete, it is not swallowed
func testServiceEndpointJenkins_Delete_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointJenkins()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("project_id", (*ep.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	flattenServiceEndpointJenkins(resourceData, ep)

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

func TestServiceEndpointJenkins_Delete_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointJenkins_Delete_DoesNotSwallowError(t, &jenkinsTestServiceEndpointPassword, jenkinsTestServiceEndpointProjectIDPassword)
}

func TestServiceEndpointJenkins_Update_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointJenkins_Delete_DoesNotSwallowError(t, &jenkinsTestServiceEndpointPassword, jenkinsTestServiceEndpointProjectIDPassword)
}

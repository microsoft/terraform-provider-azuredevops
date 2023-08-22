//go:build (all || resource_serviceendpoint_jenkins) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_jenkins
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

var jenkinsTestServiceEndpointIDpassword = uuid.New()
var jenkinsRandomServiceEndpointProjectIDpassword = uuid.New()
var jenkinsTestServiceEndpointProjectIDpassword = &jenkinsRandomServiceEndpointProjectIDpassword

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
	Id:          &jenkinsTestServiceEndpointIDpassword,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("Jenkins"),
	Url:         converter.String("https://www.jenkins.com"),
	Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: jenkinsTestServiceEndpointProjectIDpassword,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("UNIT_TEST_CONN_DESCRIPTION"),
		},
	},
}

// verifies that the flatten/expand round trip yields the same service endpoint
func testServiceEndpointJenkins_ExpandFlatten_Roundtrip(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	for _, ep := range []*serviceendpoint.ServiceEndpoint{ep, ep} {
		resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointJenkins().Schema, nil)
		flattenServiceEndpointJenkins(resourceData, ep, id)

		serviceEndpointAfterRoundTrip, projectID, err := expandServiceEndpointJenkins(resourceData)

		require.Nil(t, err)
		require.Equal(t, *ep, *serviceEndpointAfterRoundTrip)
		require.Equal(t, id, projectID)
	}
}
func TestServiceEndpointJenkins_ExpandFlatten_RoundtripPassword(t *testing.T) {
	testServiceEndpointJenkins_ExpandFlatten_Roundtrip(t, &jenkinsTestServiceEndpointPassword, jenkinsTestServiceEndpointProjectIDpassword)
}

func TestServiceEndpointJenkins_Create_DoesNotSwallowErrorPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointJenkins()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointJenkins(resourceData, &jenkinsTestServiceEndpointPassword, jenkinsTestServiceEndpointProjectIDpassword)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	seJenkins, _, _ := expandServiceEndpointJenkins(resourceData)
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
	flattenServiceEndpointJenkins(resourceData, ep, id)

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
	testServiceEndpointJenkins_Read_DoesNotSwallowError(t, &jenkinsTestServiceEndpointPassword, jenkinsTestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on a delete, it is not swallowed
func testServiceEndpointJenkins_Delete_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointJenkins()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointJenkins(resourceData, ep, id)

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
	testServiceEndpointJenkins_Delete_DoesNotSwallowError(t, &jenkinsTestServiceEndpointPassword, jenkinsTestServiceEndpointProjectIDpassword)
}

// verifies that if an error is produced on a update, it is not swallowed
func testServiceEndpointJenkins_Update_DoesNotSwallowError(t *testing.T, ep *serviceendpoint.ServiceEndpoint, id *uuid.UUID) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointJenkins()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointJenkins(resourceData, ep, id)

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
func TestServiceEndpointJenkins_Update_DoesNotSwallowErrorPassword(t *testing.T) {
	testServiceEndpointJenkins_Delete_DoesNotSwallowError(t, &jenkinsTestServiceEndpointPassword, jenkinsTestServiceEndpointProjectIDpassword)
}

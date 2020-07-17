// +build all resource_resource_authorization
// +build !exclude_resource_authorization

package build

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"testing"

	"github.com/terraform-providers/terraform-provider-azuredevops/azdosdkmocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

var projectID = "projectid"
var definitionID = 666
var endpointId = uuid.New()

var resourceReferenceAuthorized = build.DefinitionResourceReference{
	Authorized: converter.Bool(true),
	Id:         converter.String(endpointId.String()),
	Name:       nil,
	Type:       converter.String("endpoint"),
}

var resourceReferenceNotAuthorized = build.DefinitionResourceReference{
	Authorized: converter.Bool(false),
	Id:         converter.String(endpointId.String()),
	Name:       nil,
	Type:       converter.String("endpoint"),
}

var projectResourcesArgsAuthorized = build.AuthorizeProjectResourcesArgs{
	Resources: &[]build.DefinitionResourceReference{resourceReferenceAuthorized},
	Project:   &projectID,
}

var definitionResourcesArgsAuthorized = build.AuthorizeDefinitionResourcesArgs{
	Resources:    &[]build.DefinitionResourceReference{resourceReferenceAuthorized},
	Project:      &projectID,
	DefinitionId: &definitionID,
}

func TestResourceAuthorization_FlattenExpand_RoundTrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceResourceAuthorization().Schema, nil)
	flattenAuthorizedResource(resourceData, &resourceReferenceAuthorized, projectID, definitionID)

	resourceReferenceAfterRoundtrip, projectIdAfterRoundtrip, definitionIDAfterRoundTrip := expandAuthorizedResource(resourceData)
	require.Equal(t, resourceReferenceAuthorized, *resourceReferenceAfterRoundtrip)
	require.Equal(t, projectID, projectIdAfterRoundtrip)
	require.Equal(t, definitionID, definitionIDAfterRoundTrip)
}

var tests = []struct {
	Name              string
	DefinitionID      int
	MockedFunction    func(*azdosdkmocks.MockBuildClientMockRecorder, *client.AggregatedClient) *gomock.Call
	FunctionUnderTest func(*client.AggregatedClient, *schema.Resource, *schema.ResourceData) error
}{
	{
		Name: "Create project resource authorizations",
		MockedFunction: func(mr *azdosdkmocks.MockBuildClientMockRecorder, clients *client.AggregatedClient) *gomock.Call {
			return mr.AuthorizeProjectResources(clients.Ctx, projectResourcesArgsAuthorized)
		},
		FunctionUnderTest: func(clients *client.AggregatedClient, r *schema.Resource, resourceData *schema.ResourceData) error {
			return r.Create(resourceData, clients)
		},
	},
	{
		Name:         "Create pipeline resource authorizations",
		DefinitionID: definitionID,
		MockedFunction: func(mr *azdosdkmocks.MockBuildClientMockRecorder, clients *client.AggregatedClient) *gomock.Call {
			return mr.AuthorizeDefinitionResources(clients.Ctx, definitionResourcesArgsAuthorized)
		},
		FunctionUnderTest: func(clients *client.AggregatedClient, r *schema.Resource, resourceData *schema.ResourceData) error {
			return r.Create(resourceData, clients)
		},
	},
	{
		Name: "Create project resource authorizations",
		MockedFunction: func(mr *azdosdkmocks.MockBuildClientMockRecorder, clients *client.AggregatedClient) *gomock.Call {
			return mr.AuthorizeProjectResources(clients.Ctx, projectResourcesArgsAuthorized)
		},
		FunctionUnderTest: func(clients *client.AggregatedClient, r *schema.Resource, resourceData *schema.ResourceData) error {
			return r.Update(resourceData, clients)
		},
	},
	{
		Name:         "Create pipeline resource authorizations",
		DefinitionID: definitionID,
		MockedFunction: func(mr *azdosdkmocks.MockBuildClientMockRecorder, clients *client.AggregatedClient) *gomock.Call {
			return mr.AuthorizeDefinitionResources(clients.Ctx, definitionResourcesArgsAuthorized)
		},
		FunctionUnderTest: func(clients *client.AggregatedClient, r *schema.Resource, resourceData *schema.ResourceData) error {
			return r.Update(resourceData, clients)
		},
	},
	{
		Name: "Read project resource authorizations",
		MockedFunction: func(mr *azdosdkmocks.MockBuildClientMockRecorder, clients *client.AggregatedClient) *gomock.Call {
			return mr.GetProjectResources(clients.Ctx, build.GetProjectResourcesArgs{
				Project: &projectID,
				Type:    resourceReferenceAuthorized.Type,
				Id:      resourceReferenceAuthorized.Id,
			})
		},
		FunctionUnderTest: func(clients *client.AggregatedClient, r *schema.Resource, resourceData *schema.ResourceData) error {
			return r.Read(resourceData, clients)
		},
	},
	{
		Name:         "Read pipeline resource authorizations",
		DefinitionID: definitionID,
		MockedFunction: func(mr *azdosdkmocks.MockBuildClientMockRecorder, clients *client.AggregatedClient) *gomock.Call {
			return mr.GetDefinitionResources(clients.Ctx, build.GetDefinitionResourcesArgs{
				Project:      &projectID,
				DefinitionId: &definitionID,
			})
		},
		FunctionUnderTest: func(clients *client.AggregatedClient, r *schema.Resource, resourceData *schema.ResourceData) error {
			return r.Read(resourceData, clients)
		},
	},
	{
		Name: "Delete project resource authorizations",
		MockedFunction: func(mr *azdosdkmocks.MockBuildClientMockRecorder, clients *client.AggregatedClient) *gomock.Call {
			return mr.AuthorizeProjectResources(clients.Ctx, build.AuthorizeProjectResourcesArgs{
				Resources: &[]build.DefinitionResourceReference{resourceReferenceNotAuthorized},
				Project:   &projectID,
			})
		},
		FunctionUnderTest: func(clients *client.AggregatedClient, r *schema.Resource, resourceData *schema.ResourceData) error {
			return r.Delete(resourceData, clients)
		},
	},
	{
		Name:         "Delete pipeline resource authorizations",
		DefinitionID: definitionID,
		MockedFunction: func(mr *azdosdkmocks.MockBuildClientMockRecorder, clients *client.AggregatedClient) *gomock.Call {
			return mr.AuthorizeDefinitionResources(clients.Ctx, build.AuthorizeDefinitionResourcesArgs{
				Resources:    &[]build.DefinitionResourceReference{resourceReferenceNotAuthorized},
				Project:      &projectID,
				DefinitionId: &definitionID,
			})
		},
		FunctionUnderTest: func(clients *client.AggregatedClient, r *schema.Resource, resourceData *schema.ResourceData) error {
			return r.Delete(resourceData, clients)
		},
	},
}

func TestResourceAuthorization_DoesNotSwallowError(t *testing.T) {
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := ResourceResourceAuthorization()
			resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
			flattenAuthorizedResource(resourceData, &resourceReferenceAuthorized, projectID, tc.DefinitionID)

			buildClient := azdosdkmocks.NewMockBuildClient(ctrl)
			clients := &client.AggregatedClient{BuildClient: buildClient, Ctx: context.Background()}

			tc.MockedFunction(buildClient.
				EXPECT(), clients).
				Return(nil, errors.New("ResourceAuthorization Failed")).
				Times(1)

			err := tc.FunctionUnderTest(clients, r, resourceData)
			require.Contains(t, err.Error(), "ResourceAuthorization Failed")
		})
	}
}

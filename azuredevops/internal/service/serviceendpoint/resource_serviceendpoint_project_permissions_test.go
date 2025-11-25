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
	permissionsTestEndpointID     = uuid.New()
	permissionsTestOwnerProjectID = uuid.New()
	permissionsTestTargetID1      = uuid.New()
	permissionsTestTargetID2      = uuid.New()
)

var permissionsBaseEndpoint = serviceendpoint.ServiceEndpoint{
	Id:   &permissionsTestEndpointID,
	Name: converter.String("TEST_ENDPOINT"),
	Type: converter.String("azurerm"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: &permissionsTestOwnerProjectID,
			},
			Name:        converter.String("TEST_ENDPOINT"),
			Description: converter.String("Owned by this project"),
		},
	},
}

// Helper: Sets up TF data WITH the required source project_id
// Uses raw map construction to support Create operations in tests
func getPermissionsResourceData(t *testing.T, targetProject uuid.UUID) *schema.ResourceData {
	r := ResourceServiceEndpointProjectPermissions()

	input := map[string]interface{}{
		"service_endpoint_id": permissionsTestEndpointID.String(),
		"project_id":          permissionsTestOwnerProjectID.String(),
		"project_reference": []interface{}{
			map[string]interface{}{
				"project_id":            targetProject.String(),
				"service_endpoint_name": "SHARED_NAME",
				"description":           "SHARED_DESC",
			},
		},
	}

	return schema.TestResourceDataRaw(t, r.Schema, input)
}

func TestServiceEndpointProjectPermissions_Create_AppendsReference(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointProjectPermissions()
	resourceData := getPermissionsResourceData(t, permissionsTestTargetID1)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	// 1. Expect Initial Get
	// We use gomock.Any() because pointer addresses for IDs will differ between test and prod code.
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, args serviceendpoint.GetServiceEndpointDetailsArgs) (*serviceendpoint.ServiceEndpoint, error) {
			// Validate Values manually
			require.Equal(t, permissionsTestEndpointID, *args.EndpointId, "EndpointID mismatch in Get 1")
			require.Equal(t, permissionsTestOwnerProjectID.String(), *args.Project, "ProjectID mismatch in Get 1")
			return &permissionsBaseEndpoint, nil
		}).
		Times(1)

	// 2. Expect Update
	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, args serviceendpoint.UpdateServiceEndpointArgs) (*serviceendpoint.ServiceEndpoint, error) {
			// Validate Values manually
			require.Equal(t, permissionsTestEndpointID, *args.EndpointId, "EndpointID mismatch in Update")

			refs := *args.Endpoint.ServiceEndpointProjectReferences
			require.Len(t, refs, 2, "Should have Owner + New Share")

			// Verify Owner is still index 0
			require.Equal(t, permissionsTestOwnerProjectID.String(), refs[0].ProjectReference.Id.String())

			// Verify New Share is index 1
			require.Equal(t, permissionsTestTargetID1.String(), refs[1].ProjectReference.Id.String())
			require.Equal(t, "SHARED_NAME", *refs[1].Name)

			return args.Endpoint, nil
		}).
		Times(1)

	// 3. Expect Final Read (after create)
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, args serviceendpoint.GetServiceEndpointDetailsArgs) (*serviceendpoint.ServiceEndpoint, error) {
			// Validate Values manually
			require.Equal(t, permissionsTestEndpointID, *args.EndpointId, "EndpointID mismatch in Get 2")
			require.Equal(t, permissionsTestOwnerProjectID.String(), *args.Project, "ProjectID mismatch in Get 2")
			return &permissionsBaseEndpoint, nil
		}).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Nil(t, err)
}

// Test: Delete removes ONLY the target reference, keeps Owner
func TestServiceEndpointProjectPermissions_Delete_RemovesOnlyTarget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointProjectPermissions()
	resourceData := getPermissionsResourceData(t, permissionsTestTargetID1)
	resourceData.SetId(permissionsTestEndpointID.String())

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	// Setup: Endpoint currently has Owner + Target1
	endpointWithShare := permissionsBaseEndpoint
	currentRefs := append(*endpointWithShare.ServiceEndpointProjectReferences, serviceendpoint.ServiceEndpointProjectReference{
		ProjectReference: &serviceendpoint.ProjectReference{Id: &permissionsTestTargetID1},
		Name:             converter.String("SHARED_NAME"),
	})
	endpointWithShare.ServiceEndpointProjectReferences = &currentRefs

	// 1. Expect Get
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, args serviceendpoint.GetServiceEndpointDetailsArgs) (*serviceendpoint.ServiceEndpoint, error) {
			require.Equal(t, permissionsTestEndpointID, *args.EndpointId)
			require.Equal(t, permissionsTestOwnerProjectID.String(), *args.Project)
			return &endpointWithShare, nil
		}).
		Times(1)

	// 2. Expect Update
	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, args serviceendpoint.UpdateServiceEndpointArgs) (*serviceendpoint.ServiceEndpoint, error) {

			refs := *args.Endpoint.ServiceEndpointProjectReferences
			require.Len(t, refs, 1, "Should only have Owner left")
			require.Equal(t, permissionsTestOwnerProjectID.String(), refs[0].ProjectReference.Id.String())

			return args.Endpoint, nil
		}).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Nil(t, err)
}

// Test: Read filters correctly
func TestServiceEndpointProjectPermissions_Read_FiltersExpectedProjects(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointProjectPermissions()
	resourceData := getPermissionsResourceData(t, permissionsTestTargetID1)
	resourceData.SetId(permissionsTestEndpointID.String())

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	// Endpoint has Owner + Target1 + Target2 (Target 2 is not in this TF resource)
	endpointMixed := permissionsBaseEndpoint
	refs := append(*endpointMixed.ServiceEndpointProjectReferences,
		serviceendpoint.ServiceEndpointProjectReference{
			ProjectReference: &serviceendpoint.ProjectReference{Id: &permissionsTestTargetID1},
			Name:             converter.String("SHARED_NAME"),
		},
		serviceendpoint.ServiceEndpointProjectReference{
			ProjectReference: &serviceendpoint.ProjectReference{Id: &permissionsTestTargetID2},
			Name:             converter.String("OTHER_SHARE"),
		},
	)
	endpointMixed.ServiceEndpointProjectReferences = &refs

	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, gomock.Any()).
		Return(&endpointMixed, nil).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Nil(t, err)

	// Verify that the state only contains Target1, not Owner and not Target2
	resultSet := resourceData.Get("project_reference").(*schema.Set)
	require.Equal(t, 1, resultSet.Len())

	list := resultSet.List()
	obj := list[0].(map[string]interface{})
	require.Equal(t, permissionsTestTargetID1.String(), obj["project_id"])
}

// Test: Error propagation on Create
func TestServiceEndpointProjectPermissions_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointProjectPermissions()
	resourceData := getPermissionsResourceData(t, permissionsTestTargetID1)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, gomock.Any()).
		Return(nil, errors.New("GetServiceEndpointDetails() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpointDetails() Failed")
}

// Test: Error propagation on Update
func TestServiceEndpointProjectPermissions_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointProjectPermissions()
	resourceData := getPermissionsResourceData(t, permissionsTestTargetID1)
	// Simulate Update by giving it an ID
	resourceData.SetId(permissionsTestEndpointID.String())

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	// Get succeeds
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, gomock.Any()).
		Return(&permissionsBaseEndpoint, nil).
		Times(1)

	// Update fails
	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, gomock.Any()).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

// Test: Error propagation on Delete
func TestServiceEndpointProjectPermissions_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointProjectPermissions()
	resourceData := getPermissionsResourceData(t, permissionsTestTargetID1)
	resourceData.SetId(permissionsTestEndpointID.String())

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, gomock.Any()).
		Return(nil, errors.New("GetServiceEndpointDetails() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpointDetails() Failed")
}

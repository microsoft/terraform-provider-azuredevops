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
		ShareServiceEndpoint(clients.Ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, args serviceendpoint.ShareServiceEndpointArgs) error {
			// Validate Values manually
			require.Equal(t, permissionsTestEndpointID, *args.EndpointId, "EndpointID mismatch in Update")

			refs := *args.EndpointProjectReferences
			require.Len(t, refs, 1, "Should have New Share")

			// Verify New Share is index 0
			require.Equal(t, permissionsTestTargetID1.String(), refs[0].ProjectReference.Id.String())
			require.Equal(t, "SHARED_NAME", *refs[0].Name)

			return nil
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

	err := r.CreateContext(clients.Ctx, resourceData, clients)
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
		DeleteServiceEndpoint(clients.Ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, args serviceendpoint.DeleteServiceEndpointArgs) error {
			require.Len(t, *args.ProjectIds, 1, "Should only delete 1 project")
			require.Equal(t, permissionsTestTargetID1.String(), (*args.ProjectIds)[0])

			return nil
		}).
		Times(1)

	err := r.DeleteContext(clients.Ctx, resourceData, clients)
	require.Nil(t, err)
}

// Test: Read filters correctly
func TestServiceEndpointProjectPermissions_Read_TracksAllSharedProjects(t *testing.T) {
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

	err := r.ReadContext(clients.Ctx, resourceData, clients)
	require.Nil(t, err)

	// Verify that the state contains Target1 and Target2, but not Owner
	resultSet := resourceData.Get("project_reference").([]interface{})
	require.Equal(t, 2, len(resultSet))

	pids := []string{
		resultSet[0].(map[string]interface{})["project_id"].(string),
		resultSet[1].(map[string]interface{})["project_id"].(string),
	}
	require.Contains(t, pids, permissionsTestTargetID1.String())
	require.Contains(t, pids, permissionsTestTargetID2.String())
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

	diag := r.CreateContext(clients.Ctx, resourceData, clients)
	require.True(t, diag.HasError(), "Expected diagnostics to have an error")
	require.Contains(t, diag[0].Summary, "GetServiceEndpointDetails() Failed")
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
		ShareServiceEndpoint(clients.Ctx, gomock.Any()).
		Return(errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	diags := r.UpdateContext(clients.Ctx, resourceData, clients)
	require.True(t, diags.HasError(), "Expected diagnostics to have an error")
	require.Contains(t, diags[0].Summary, "UpdateServiceEndpoint() Failed")
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

	diags := r.DeleteContext(clients.Ctx, resourceData, clients)
	require.True(t, diags.HasError(), "Expected diagnostics to have an error")
	require.Contains(t, diags[0].Summary, "GetServiceEndpointDetails() Failed")
}

//go:build (all || resource_securityrole_assignment) && !exclude_securityroles
// +build all resource_securityrole_assignment
// +build !exclude_securityroles

package securityroles

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/securityroles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var SecurityRoleAssignmentIdentityID = uuid.New()
var SecurityRoleAssignmentScope = "some:scope"
var SecurityRoleAssignmentResourceID = "123456789"
var SecurityRoleAssignmentRole = "Admin"

// verifies that if an error is produced on create, the error is not swallowed

func TestSecurityRoleAssignment_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceSecurityRoleAssignment()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"role_name":   SecurityRoleAssignmentRole,
		"identity_id": SecurityRoleAssignmentIdentityID.String(),
		"resource_id": SecurityRoleAssignmentResourceID,
		"scope":       SecurityRoleAssignmentScope,
	})

	securityrolesClient := azdosdkmocks.NewMockSecurityrolesClient(ctrl)
	clients := &client.AggregatedClient{SecurityRolesClient: securityrolesClient, Ctx: context.Background()}

	expectedArgs := securityroles.SetSecurityRoleAssignmentArgs{
		IdentityId: &SecurityRoleAssignmentIdentityID,
		Scope:      &SecurityRoleAssignmentScope,
		ResourceId: &SecurityRoleAssignmentResourceID,
		RoleName:   &SecurityRoleAssignmentRole,
	}
	securityrolesClient.
		EXPECT().
		SetSecurityRoleAssignment(clients.Ctx, &expectedArgs).
		Return(fmt.Errorf("invalid UUID length")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid UUID length")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestSecurityRoleAssignment_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceSecurityRoleAssignment()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"identity_id": SecurityRoleAssignmentIdentityID.String(),
		"resource_id": SecurityRoleAssignmentResourceID,
		"scope":       SecurityRoleAssignmentScope,
	})

	securityrolesClient := azdosdkmocks.NewMockSecurityrolesClient(ctrl)
	clients := &client.AggregatedClient{SecurityRolesClient: securityrolesClient, Ctx: context.Background()}

	expectedArgs := securityroles.GetSecurityRoleAssignmentArgs{
		Scope:      &SecurityRoleAssignmentScope,
		ResourceId: &SecurityRoleAssignmentResourceID,
		IdentityId: &SecurityRoleAssignmentIdentityID,
	}

	securityrolesClient.
		EXPECT().
		GetSecurityRoleAssignment(clients.Ctx, &expectedArgs).
		Return(nil, errors.New("invalid UUID length")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "invalid UUID length")
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestSecurityRoleAssignment_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceSecurityRoleAssignment()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"identity_id": SecurityRoleAssignmentIdentityID.String(),
		"resource_id": SecurityRoleAssignmentResourceID,
		"scope":       SecurityRoleAssignmentScope,
	})

	securityrolesClient := azdosdkmocks.NewMockSecurityrolesClient(ctrl)
	clients := &client.AggregatedClient{SecurityRolesClient: securityrolesClient, Ctx: context.Background()}

	expectedArgs := securityroles.DeleteSecurityRoleAssignmentArgs{
		Scope:      &SecurityRoleAssignmentScope,
		ResourceId: &SecurityRoleAssignmentResourceID,
		IdentityId: &SecurityRoleAssignmentIdentityID,
	}

	securityrolesClient.
		EXPECT().
		DeleteSecurityRoleAssignment(clients.Ctx, &expectedArgs).
		Return(errors.New("invalid UUID length")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "invalid UUID length")
}

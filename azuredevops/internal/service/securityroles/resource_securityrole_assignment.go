package securityroles

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/securityroles"
)

func ResourceSecurityRoleAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityRoleAssignmentCreateOrUpdate,
		Read:   resourceSecurityRoleAssignmentRead,
		Update: resourceSecurityRoleAssignmentCreateOrUpdate,
		Delete: resourceSecurityRoleAssignmentDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"scope": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
				Required:     true,
			},
			"resource_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
				Required:     true,
				ForceNew:     true,
			},
			"identity_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Required:     true,
			},
			"role_name": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
				Required:     true,
			},
		},
	}
}

func resourceSecurityRoleAssignmentCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	scope := d.Get("scope").(string)
	resourceId := d.Get("resource_id").(string)

	identityId, err := uuid.Parse(d.Get("identity_id").(string))
	if err != nil {
		return err
	}

	roleName := d.Get("role_name").(string)
	err = clients.SecurityRolesClient.SetSecurityRoleAssignment(clients.Ctx, &securityroles.SetSecurityRoleAssignmentArgs{
		Scope:      &scope,
		ResourceId: &resourceId,
		IdentityId: &identityId,
		RoleName:   &roleName,
	})

	if err != nil {
		return err
	}

	stateConf := &retry.StateChangeConf{
		ContinuousTargetOccurence: 2,
		Delay:                     5 * time.Second,
		MinTimeout:                10 * time.Second,
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Pending:                   []string{"syncing"},
		Target:                    []string{"succeed", "failed"},
		Refresh:                   getSecurityRoleAssignment(*clients, scope, roleName, resourceId, identityId),
	}
	_, err = stateConf.WaitForStateContext(clients.Ctx)
	if err != nil {
		return err
	}

	d.SetId("sra-" + uuid.New().String())
	return resourceSecurityRoleAssignmentRead(d, m)
}

func resourceSecurityRoleAssignmentRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	scope := d.Get("scope").(string)
	resourceId := d.Get("resource_id").(string)
	identityId, err := uuid.Parse(d.Get("identity_id").(string))
	if err != nil {
		return err
	}

	assignment, err := clients.SecurityRolesClient.GetSecurityRoleAssignment(clients.Ctx, &securityroles.GetSecurityRoleAssignmentArgs{
		Scope:      &scope,
		ResourceId: &resourceId,
		IdentityId: &identityId,
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(" reading group memberships during read: %+v", err)
	}

	if assignment != nil && (assignment.Identity == nil && assignment.Role == nil) {
		d.SetId("")
		return nil
	}

	if assignment.Role != nil {
		d.Set("scope", *assignment.Role.Scope)
		d.Set("role_name", *assignment.Role.Name)
	}
	if assignment.Identity != nil {
		d.Set("identity_id", *assignment.Identity.ID)
	}
	d.Set("resource_id", resourceId)

	return nil
}

func resourceSecurityRoleAssignmentDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	scope := d.Get("scope").(string)
	resourceId := d.Get("resource_id").(string)

	identityId, err := uuid.Parse(d.Get("identity_id").(string))
	if err != nil {
		return err
	}

	err = clients.SecurityRolesClient.DeleteSecurityRoleAssignment(clients.Ctx, &securityroles.DeleteSecurityRoleAssignmentArgs{
		Scope:      &scope,
		ResourceId: &resourceId,
		IdentityId: &identityId,
	})

	if err != nil {
		return err
	}

	return nil
}

func getSecurityRoleAssignment(clients client.AggregatedClient, scope, roleName, resourceId string, identityId uuid.UUID) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		assigns, err := clients.SecurityRolesClient.GetSecurityRoleAssignment(clients.Ctx, &securityroles.GetSecurityRoleAssignmentArgs{
			Scope:      &scope,
			ResourceId: &resourceId,
			IdentityId: &identityId,
		})

		if err != nil {
			return "", "failed", nil
		}

		if assigns != nil && (assigns.Identity == nil && assigns.Role == nil) {
			return "", "syncing", nil
		}

		if assigns != nil && assigns.Identity != nil && assigns.Identity.ID != nil &&
			!strings.EqualFold(*assigns.Identity.ID, identityId.String()) {
			return "", "syncing", nil
		}

		if assigns.Role != nil && assigns.Role.Name != nil &&
			!strings.EqualFold(*assigns.Role.Name, roleName) && !strings.EqualFold(*assigns.Role.Scope, scope) {
			return "", "syncing", nil
		}
		return "", "succeed", nil
	}
}

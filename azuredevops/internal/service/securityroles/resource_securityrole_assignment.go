package securityroles

import (
	"fmt"

	"github.com/google/uuid"
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

	if assignment != nil {
		if assignment.Role != nil {
			d.Set("scope", *assignment.Role.Scope)
			d.Set("role_name", *assignment.Role.Name)
		}
		if assignment.Identity != nil {
			d.Set("identity_id", *assignment.Identity.ID)
		}
		d.Set("resource_id", resourceId)
	}

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

	d.SetId("")
	return nil
}

package permissions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
)

func ResourceIterationPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceIterationPermissionsCreateOrUpdate,
		Read:   resourceIterationPermissionsRead,
		Update: resourceIterationPermissionsCreateOrUpdate,
		Delete: resourceIterationPermissionsDelete,
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Required:     true,
				ForceNew:     true,
			},
			"path": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Default:      "",
				ForceNew:     true,
				Optional:     true,
			},
		}),
	}
}

func createIterationToken(context context.Context, workitemtrackingClient workitemtracking.Client, d *schema.ResourceData) (*string, error) {
	projectID := d.Get("project_id").(string)
	path := d.Get("path").(string)
	aclToken, err := securityhelper.CreateClassificationNodeSecurityToken(context, workitemtrackingClient, workitemtracking.TreeStructureGroupValues.Iterations, projectID, path)
	if err != nil {
		return nil, err
	}
	return &aclToken, nil
}

func resourceIterationPermissionsCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.Iteration,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return err
	}

	aclToken, err := createIterationToken(clients.Ctx, clients.WorkItemTrackingClient, d)
	if err != nil {
		return err
	}

	err = securityhelper.SetPrincipalPermissions(d, sn, aclToken, nil, false)
	if err != nil {
		return err
	}

	return resourceIterationPermissionsRead(d, m)
}

func resourceIterationPermissionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.Iteration,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return err
	}

	aclToken, err := createIterationToken(clients.Ctx, clients.WorkItemTrackingClient, d)
	if err != nil {
		return err
	}

	principalPermissions, err := securityhelper.GetPrincipalPermissions(d, sn, aclToken)
	if err != nil {
		return err
	}

	d.Set("permissions", principalPermissions.Permissions)
	return nil
}

func resourceIterationPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.Iteration,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return err
	}

	aclToken, err := createIterationToken(clients.Ctx, clients.WorkItemTrackingClient, d)
	if err != nil {
		return err
	}

	err = securityhelper.SetPrincipalPermissions(d, sn, aclToken, &securityhelper.PermissionTypeValues.NotSet, true)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

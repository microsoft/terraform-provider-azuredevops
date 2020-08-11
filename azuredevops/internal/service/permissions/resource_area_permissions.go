package permissions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
)

func ResourceAreaPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceAreaPermissionsCreateOrUpdate,
		Read:   resourceAreaPermissionsRead,
		Update: resourceAreaPermissionsCreateOrUpdate,
		Delete: resourceAreaPermissionsDelete,
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
				ForceNew:     true,
				Optional:     true,
			},
		}),
	}
}

func resourceAreaPermissionsCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, aclToken, err := initializeAreaSecurityNamespaceAndToken(d, clients)
	if err != nil {
		return err
	}

	if err = securityhelper.SetPrincipalPermissions(d, sn, aclToken, nil, false); err != nil {
		return err
	}

	return resourceAreaPermissionsRead(d, m)
}

func resourceAreaPermissionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, aclToken, err := initializeAreaSecurityNamespaceAndToken(d, clients)
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

func resourceAreaPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, aclToken, err := initializeAreaSecurityNamespaceAndToken(d, clients)
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

func initializeAreaSecurityNamespaceAndToken(d *schema.ResourceData, clients *client.AggregatedClient) (*securityhelper.SecurityNamespace, *string, error) {
	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.CSS,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return nil, nil, err
	}

	aclToken, err := createAreaToken(clients.Ctx, clients.WorkItemTrackingClient, d)
	if err != nil {
		return nil, nil, err
	}

	return sn, aclToken, nil
}

func createAreaToken(context context.Context, workitemtrackingClient workitemtracking.Client, d *schema.ResourceData) (*string, error) {
	projectID := d.Get("project_id").(string)
	path := d.Get("path").(string)
	aclToken, err := securityhelper.CreateClassificationNodeSecurityToken(context, workitemtrackingClient, workitemtracking.TreeStructureGroupValues.Areas, projectID, path)
	if err != nil {
		return nil, err
	}
	return &aclToken, nil
}

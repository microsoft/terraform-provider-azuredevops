package permissions

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
)

// ResourceServiceHooksPermissions schema and implementation for servicehooks permission resource
func ResourceServiceHooksPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceHooksPermissionsCreateOrUpdate,
		Read:   resourceServiceHooksPermissionsRead,
		Update: resourceServiceHooksPermissionsCreateOrUpdate,
		Delete: resourceServiceHooksPermissionsDelete,
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Required:     true,
				ForceNew:     true,
			},
		}),
	}
}

func resourceServiceHooksPermissionsCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.ServiceHooks, createServiceHooksToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, nil, false); err != nil {
		return err
	}

	return resourceServiceHooksPermissionsRead(d, m)
}

func resourceServiceHooksPermissionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.ServiceHooks, createServiceHooksToken)
	if err != nil {
		return err
	}

	principalPermissions, err := securityhelper.GetPrincipalPermissions(d, sn)
	if err != nil {
		return err
	}
	if principalPermissions == nil {
		d.SetId("")
		log.Printf("[INFO] Permissions for ACL token %q not found. Removing from state", sn.GetToken())
		return nil
	}

	d.Set("permissions", principalPermissions.Permissions)
	return nil
}

func resourceServiceHooksPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.ServiceHooks, createServiceHooksToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, &securityhelper.PermissionTypeValues.NotSet, true); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func createServiceHooksToken(d *schema.ResourceData, clients *client.AggregatedClient) (string, error) {
	projectID := d.Get("project_id").(string)
	aclToken := "PublisherSecurity/" + projectID
	return aclToken, nil
}

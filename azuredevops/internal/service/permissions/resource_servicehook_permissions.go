package permissions

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
)

// ResourceServiceHookPermissions schema and implementation for servicehook permission resource
func ResourceServiceHookPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceHookPermissionsCreateOrUpdate,
		Read:   resourceServiceHookPermissionsRead,
		Update: resourceServiceHookPermissionsCreateOrUpdate,
		Delete: resourceServiceHookPermissionsDelete,
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Optional:     true,
				ForceNew:     true,
			},
		}),
	}
}

func resourceServiceHookPermissionsCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.ServiceHooks, createServiceHookToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, nil, false); err != nil {
		return err
	}

	return resourceServiceHookPermissionsRead(d, m)
}

func resourceServiceHookPermissionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.ServiceHooks, createServiceHookToken)
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

func resourceServiceHookPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.ServiceHooks, createServiceHookToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, &securityhelper.PermissionTypeValues.NotSet, true); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func createServiceHookToken(d *schema.ResourceData, clients *client.AggregatedClient) (string, error) {
	projectID, ok := d.GetOk("project_id")
	if !ok {
		return "PublisherSecurity", nil
	}
	return fmt.Sprintf("PublisherSecurity/%s", projectID.(string)), nil
}

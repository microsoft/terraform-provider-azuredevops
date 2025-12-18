package permissions

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
)

func ResourceWorkItemTrackingProcessPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceProcessPermissionsCreateOrUpdate,
		Read:   resourceProcessPermissionsRead,
		Update: resourceProcessPermissionsCreateOrUpdate,
		Delete: resourceProcessPermissionsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
			"parent_process_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Required:     true,
				ForceNew:     true,
				Description:  "The ID of the parent process.",
			},
			"process_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Required:     true,
				ForceNew:     true,
				Description:  "The ID of the process.",
			},
		}),
	}
}

func resourceProcessPermissionsCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.Process, createProcessToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, nil, false); err != nil {
		return err
	}

	return resourceProcessPermissionsRead(d, m)
}

func resourceProcessPermissionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.Process, createProcessToken)
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

func resourceProcessPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.Process, createProcessToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, &securityhelper.PermissionTypeValues.NotSet, true); err != nil {
		return err
	}
	return nil
}

func createProcessToken(d *schema.ResourceData, clients *client.AggregatedClient) (string, error) {
	parentProcessID, ok := d.GetOk("parent_process_id")
	if !ok {
		return "", fmt.Errorf("Failed to get 'parent_process_id' from schema")
	}

	processID, ok := d.GetOk("process_id")
	if !ok {
		return "", fmt.Errorf("Failed to get 'process_id' from schema")
	}

	aclToken := fmt.Sprintf("$PROCESS:%s:%s:", parentProcessID.(string), processID.(string))
	return aclToken, nil
}

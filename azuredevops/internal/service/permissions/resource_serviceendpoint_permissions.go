package permissions

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
)

// ResourceServiceEndpointPermissions schema and implementation for serviceendpoint permission resource
func ResourceServiceEndpointPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceEndpointPermissionsCreateOrUpdate,
		Read:   resourceServiceEndpointPermissionsRead,
		Update: resourceServiceEndpointPermissionsCreateOrUpdate,
		Delete: resourceServiceEndpointPermissionsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Required:     true,
				ForceNew:     true,
			},
			"serviceendpoint_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				ForceNew:     true,
				Optional:     true,
			},
		}),
	}
}

func resourceServiceEndpointPermissionsCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.ServiceEndpoints, createServiceEndpointToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, nil, false); err != nil {
		return err
	}

	return resourceServiceEndpointPermissionsRead(d, m)
}

func resourceServiceEndpointPermissionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.ServiceEndpoints, createServiceEndpointToken)
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

func resourceServiceEndpointPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.ServiceEndpoints, createServiceEndpointToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, &securityhelper.PermissionTypeValues.NotSet, true); err != nil {
		return err
	}

	return nil
}

func createServiceEndpointToken(d *schema.ResourceData, clients *client.AggregatedClient) (string, error) {
	projectID := d.Get("project_id").(string)
	// Token format for ALL service endpoints in a project: endpoints/ProjectID
	// Token format for a specific endpoint in a project: endpoints/ProjectID/ServiceEndpointID
	aclToken := "endpoints/" + projectID
	serviceEndpointID, serviceEndpointOk := d.GetOk("serviceendpoint_id")
	if serviceEndpointOk {
		aclToken += "/" + serviceEndpointID.(string)
	}
	return aclToken, nil
}

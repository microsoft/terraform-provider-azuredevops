// Implementation for the azuredevops_serviceendpoint_project_permissions resource
package serviceendpoint

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceServiceEndpointProjectPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceEndpointProjectPermissionsCreate,
		Read:   resourceServiceEndpointProjectPermissionsRead,
		Update: resourceServiceEndpointProjectPermissionsUpdate,
		Delete: resourceServiceEndpointProjectPermissionsDelete,

		Schema: map[string]*schema.Schema{
			"serviceendpoint_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_reference": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"service_endpoint_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceServiceEndpointProjectPermissionsCreate(d *schema.ResourceData, m interface{}) error {
	// Implementation logic for creating the resource
	return fmt.Errorf("Not implemented yet")
}

func resourceServiceEndpointProjectPermissionsRead(d *schema.ResourceData, m interface{}) error {
	// Implementation logic for reading the resource
	return fmt.Errorf("Not implemented yet")
}

func resourceServiceEndpointProjectPermissionsUpdate(d *schema.ResourceData, m interface{}) error {
	// Implementation logic for updating the resource
	return fmt.Errorf("Not implemented yet")
}

func resourceServiceEndpointProjectPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	// Implementation logic for deleting the resource
	return fmt.Errorf("Not implemented yet")
}

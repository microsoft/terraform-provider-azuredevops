package azuredevops

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)

func resourceAzureProjectPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzureProjectPermissionsCreate,
		Read:   resourceAzureProjectPermissionsRead,
		Update: resourceAzureProjectPermissionsUpdate,
		Delete: resourceAzureProjectPermissionsDelete,

		Schema: map[string]*schema.Schema{
			// add properties here
		},
	}
}

func resourceAzureProjectPermissionsCreate(d *schema.ResourceData, m interface{}) error {

	return resourceAzureProjectPermissionsRead(d, m)
}

func resourceAzureProjectPermissionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)

	return nil
}

func resourceAzureProjectPermissionsUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceAzureProjectPermissionsRead(d, m)
}

func resourceAzureProjectPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

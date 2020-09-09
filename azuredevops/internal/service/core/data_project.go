package core

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

// DataProject schema and implementation for project data source
func DataProject() *schema.Resource {
	return &schema.Resource{
		Read: dataProjectRead,
		Schema: map[string]*schema.Schema{
			"project_identifier": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_control": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"work_item_template": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"process_template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"features": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

// Introducing a read method here which is almost the same code a in resource_project.go
// but this follows the `A little copying is better than a little dependency.` GO proverb.
func dataProjectRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	identifier := d.Get("project_identifier").(string)

	project, err := projectRead(clients, identifier, identifier)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return fmt.Errorf("Project with name or ID %s does not exist", identifier)
		}
		return fmt.Errorf("Error looking up project with Name or ID %s, %w", identifier, err)
	}

	err = flattenProject(clients, d, project)
	if err != nil {
		return fmt.Errorf("Error flattening project: %v", err)
	}
	return nil
}

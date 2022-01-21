package core

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/core"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// DataProject schema and implementation for project data source
func DataProject() *schema.Resource {
	return &schema.Resource{
		Read: dataProjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				ConflictsWith: []string{
					"name",
				},
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

	name := d.Get("name").(string)
	id := d.Get("project_id").(string)

	if name == "" && id == "" {
		return fmt.Errorf("Either project_id or name must be set ")
	}

	identifier := id
	if identifier == "" {
		identifier = name
	}

	project, err := clients.CoreClient.GetProject(clients.Ctx, core.GetProjectArgs{
		ProjectId:           &identifier,
		IncludeCapabilities: converter.Bool(true),
		IncludeHistory:      converter.Bool(false),
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return fmt.Errorf("Project with name %s or ID %s does not exist ", name, id)
		}
		return fmt.Errorf("Error looking up project with Name %s or ID %s, %+v ", name, id, err)
	}

	err = flattenProject(clients, d, project)
	d.Set("project_id", project.Id.String())
	if err != nil {
		return fmt.Errorf("Error flattening project: %v", err)
	}
	return nil
}

package core

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
)

// DataProject schema and implementation for project data source
func DataProject() *schema.Resource {
	baseSchema := ResourceProject()
	for k, v := range baseSchema.Schema {
		if k == "project_name" {
			v.Required = false
			v.Optional = true
			baseSchema.Schema[k] = v
		} else {
			baseSchema.Schema[k] = &schema.Schema{
				Type:     v.Type,
				Computed: true,
			}
		}
	}

	baseSchema.Schema["project_id"] = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ValidateFunc:     validation.StringIsNotWhiteSpace,
		DiffSuppressFunc: suppress.CaseDifference,
		ConflictsWith: []string{
			"project_name",
		},
	}

	return &schema.Resource{
		Read:   dataProjectRead,
		Schema: baseSchema.Schema,
	}
}

// Introducing a read method here which is almost the same code a in resource_project.go
// but this follows the `A little copying is better than a little dependency.` GO proverb.
func dataProjectRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	name := d.Get("project_name").(string)
	id := d.Get("project_id").(string)

	if name == "" && id == "" {
		return fmt.Errorf("Either project_id or project_name must be set")
	}

	project, err := projectRead(clients, id, name)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return fmt.Errorf("Project with name %s does not exist", name)
		}
		return fmt.Errorf("Error looking up project with Name %s, %w", name, err)
	}

	err = flattenProject(clients, d, project)
	if err != nil {
		return fmt.Errorf("Error flattening project: %v", err)
	}
	return nil
}

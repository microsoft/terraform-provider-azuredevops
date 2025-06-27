package core

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// DataProject schema and implementation for project data source
func DataProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataProjectRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsNotWhiteSpace,
				ConflictsWith: []string{"project_id"},
				AtLeastOneOf:  []string{"name", "project_id"},
			},
			"project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsNotWhiteSpace,
				ConflictsWith: []string{"name"},
				AtLeastOneOf:  []string{"name", "project_id"},
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

func dataProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	name := d.Get("name").(string)
	id := d.Get("project_id").(string)

	identifier := id
	if identifier == "" {
		identifier = name
	}

	project, err := clients.CoreClient.GetProject(ctx, core.GetProjectArgs{
		ProjectId:           &identifier,
		IncludeCapabilities: converter.Bool(true),
		IncludeHistory:      converter.Bool(false),
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return diag.FromErr(fmt.Errorf("Project with name %s or ID %s does not exist ", name, id))
		}
		return diag.FromErr(fmt.Errorf("Error looking up project with Name %s or ID %s, %+v ", name, id, err))
	}

	d.SetId(project.Id.String())
	d.Set("project_id", project.Id.String())

	err = flattenProject(clients, d, project)
	if err != nil {
		return diag.FromErr(fmt.Errorf("flattening project: %v", err))
	}
	return nil
}

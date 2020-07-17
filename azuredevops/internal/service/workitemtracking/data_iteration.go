package azuredevops

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"
)

func dataIteration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIterationRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.UUID,
			},
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func dataSourceIterationRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	projectID := d.Get("project_id").(string)

	args := workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
		Depth:          converter.Int(999),
	}

	path, ok := d.GetOk("path")
	if ok {
		args.Path = converter.String(path.(string))
	}

	iteration, err := clients.WitClient.GetClassificationNode(clients.Ctx, args)
	if err != nil {
		return fmt.Errorf("Error getting Iteration with path %q: %+v", path, err)
	}

	d.SetId(iteration.Identifier.String())
	return nil
}

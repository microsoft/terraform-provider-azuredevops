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
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.UUID,
			},
		},
	}
}

func dataSourceIterationRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	path, projectID := d.Get("path").(string), d.Get("project_id").(string)

	getClassificationNodeArgs := workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
		Path:           &path,
		Depth:          converter.Int(999),
	}

	iteration, err := clients.WitClient.GetClassificationNode(clients.Ctx, getClassificationNodeArgs)
	if err != nil {
		return fmt.Errorf("Error getting Iteration with path %q: %+v", path, err)
	}

	d.SetId(iteration.Identifier.String())
	return nil
}

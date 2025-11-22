package workitemtracking

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/workitemtracking/utils"
)

func ResourceIteration() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateIteration,
		Read:   resourceReadIteration,
		Update: resourceUpdateIteration,
		Delete: resourceDeleteIteration,

		Schema: utils.CreateClassificationNodeResourceSchema(workitemtracking.TreeStructureGroupValues.Iterations),
	}
}

func resourceCreateIteration(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	return utils.CreateOrUpdateClassificationNode(clients, d, workitemtracking.TreeStructureGroupValues.Iterations)
}

func resourceReadIteration(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	return utils.ReadClassificationNode(clients, d, workitemtracking.TreeStructureGroupValues.Iterations)
}

func resourceUpdateIteration(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	return utils.CreateOrUpdateClassificationNode(clients, d, workitemtracking.TreeStructureGroupValues.Iterations)
}

func resourceDeleteIteration(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	return utils.DeleteClassificationNode(clients, d, workitemtracking.TreeStructureGroupValues.Iterations)
}

package workitemtracking

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/workitemtracking/utils"
)

// DataIteration schema for iteration data
func DataIteration() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceIterationRead,
		Schema: utils.CreateClassificationNodeSchema(map[string]*schema.Schema{}),
	}
}

func dataSourceIterationRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	return utils.ReadClassificationNode(clients, d, workitemtracking.TreeStructureGroupValues.Iterations)
}

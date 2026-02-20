package workitemtracking

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/workitemtracking/utils"
)

func ResourceArea() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateArea,
		Read:   resourceReadArea,
		Update: resourceUpdateArea,
		Delete: resourceDeleteArea,

		Schema: utils.CreateClassificationNodeResourceSchema(workitemtracking.TreeStructureGroupValues.Areas),
	}
}

func resourceCreateArea(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	return utils.CreateOrUpdateClassificationNode(clients, d, workitemtracking.TreeStructureGroupValues.Areas)
}

func resourceReadArea(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	return utils.ReadClassificationNode(clients, d, workitemtracking.TreeStructureGroupValues.Areas)
}

func resourceUpdateArea(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	return utils.CreateOrUpdateClassificationNode(clients, d, workitemtracking.TreeStructureGroupValues.Areas)
}

func resourceDeleteArea(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	return utils.DeleteClassificationNode(clients, d, workitemtracking.TreeStructureGroupValues.Areas)
}

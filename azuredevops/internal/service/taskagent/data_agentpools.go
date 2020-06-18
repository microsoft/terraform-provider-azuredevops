package taskagent

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
)

// DataAgentPools schema and implementation for agent pools data source
func DataAgentPools() *schema.Resource {
	baseSchema := ResourceAgentPool()

	// Now that the base schema's ID is not being used as the resource's ID, we can correctly
	// set it to be an integer.
	baseSchema.Schema["id"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}

	for k, v := range baseSchema.Schema {
		baseSchema.Schema[k] = &schema.Schema{
			Type:     v.Type,
			Computed: true,
		}
	}

	return &schema.Resource{
		Read: dataSourceAgentPoolsRead,

		Schema: map[string]*schema.Schema{
			"agent_pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: baseSchema.Schema,
				},
			},
		},
	}
}

func dataSourceAgentPoolsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	agentPools, err := getAgentPools(clients)
	if err != nil {
		return fmt.Errorf("Error finding agent pools. Error: %v", err)
	}
	log.Printf("[TRACE] plugin.terraform-provider-azuredevops: Read [%d] agent pools from current organization", len(*agentPools))

	err = d.Set("agent_pools", flattenAgentPoolReferences(agentPools))
	if err != nil {
		return fmt.Errorf("Error setting agent_pools field in state. Error: %v", err)
	}

	d.SetId(time.Now().UTC().String())
	return nil
}

func flattenAgentPoolReferences(input *[]taskagent.TaskAgentPool) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)

	for _, element := range *input {
		output := make(map[string]interface{})
		if element.Name != nil {
			output["name"] = *element.Name
		}

		if element.Id != nil {
			output["id"] = *element.Id
		}

		if element.PoolType != nil {
			output["pool_type"] = string(*element.PoolType)
		}

		if element.AutoProvision != nil {
			output["auto_provision"] = *element.AutoProvision
		}

		results = append(results, output)
	}

	return results
}

func getAgentPools(clients *client.AggregatedClient) (*[]taskagent.TaskAgentPool, error) {
	return clients.TaskAgentClient.GetAgentPools(clients.Ctx, taskagent.GetAgentPoolsArgs{})
}

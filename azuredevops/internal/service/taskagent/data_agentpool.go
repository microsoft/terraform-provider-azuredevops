package taskagent

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// DataAgentPool schema and implementation for agent pool data source
func DataAgentPool() *schema.Resource {
	baseSchema := ResourceAgentPool()
	for k, v := range baseSchema.Schema {
		if k != "name" {
			baseSchema.Schema[k] = &schema.Schema{
				Type:     v.Type,
				Computed: true,
			}
		}
	}

	return &schema.Resource{
		Read:   dataSourceAgentPoolRead,
		Schema: baseSchema.Schema,
	}
}

func dataSourceAgentPoolRead(d *schema.ResourceData, m interface{}) error {
	agentPoolName := d.Get("name").(string)
	clients := m.(*client.AggregatedClient)

	agentPool, err := getAgentPoolByName(clients, &agentPoolName)
	if err != nil {
		return fmt.Errorf("Error getting agent pool by name: %v", err)
	}

	flattenAzureAgentPool(d, agentPool)
	return nil
}

func getAgentPoolByName(clients *client.AggregatedClient, name *string) (*taskagent.TaskAgentPool, error) {
	agentPools, err := clients.TaskAgentClient.GetAgentPools(clients.Ctx, taskagent.GetAgentPoolsArgs{
		PoolName: name,
	})

	if err != nil {
		return nil, err
	}

	if len(*agentPools) > 1 {
		return nil, fmt.Errorf("Found multiple agent pools for name: %s. Agent pools found: %+v", *name, agentPools)
	}

	if len(*agentPools) == 0 {
		return nil, fmt.Errorf("Unable to find agent pool with name: %s", *name)
	}

	return &(*agentPools)[0], nil
}

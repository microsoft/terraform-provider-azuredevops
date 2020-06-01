package azuredevops

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
)

func dataAzureAgentPool() *schema.Resource {
	baseSchema := resourceAzureAgentPool()
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
	clients := m.(*config.AggregatedClient)

	agentPool, err := getAgentPoolByName(clients, &agentPoolName)
	if err != nil {
		return fmt.Errorf("Error getting agent pool by name: %v", err)
	}

	flattenAzureAgentPool(d, agentPool)
	return err
}

func getAgentPoolByName(clients *config.AggregatedClient, name *string) (*taskagent.TaskAgentPool, error) {
	agentPools, err := clients.TaskAgentClient.GetAgentPools(clients.Ctx, taskagent.GetAgentPoolsArgs{
		PoolName: name,
	})

	if err != nil {
		return nil, err
	}

	if len(*agentPools) > 1 {
		return nil, fmt.Errorf("Found multiple agent pools for name: %s. Agent pools found: %v", *name, agentPools)
	}

	if len(*agentPools) == 0 {
		return nil, fmt.Errorf("Unable to find agent pool with name: %s", *name)
	}

	return &(*agentPools)[0], nil
}

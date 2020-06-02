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

	agentPools, err := getAgentsPoolByName(clients, &agentPoolName)
	if err != nil {
		return fmt.Errorf("Error getting agent pool by name: %v", err)
	}

	// too many agent pools found
	if len(*agentPools) > 1 {
		return fmt.Errorf("Found multiple agent pools for name: %s. Agent pools found: %v", agentPoolName, agentPools)
	}

	// no agent pools found - handle gracefully
	if len(*agentPools) == 0 {
		d.SetId("")
		return nil
	}

	flattenAzureAgentPool(d, &(*agentPools)[0])
	return err
}

func getAgentsPoolByName(clients *config.AggregatedClient, name *string) (*[]taskagent.TaskAgentPool, error) {
	return clients.TaskAgentClient.GetAgentPools(clients.Ctx, taskagent.GetAgentPoolsArgs{
		PoolName: name,
	})
}

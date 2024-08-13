package taskagent

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// DataAgentPools schema and implementation for agent pools data source
func DataAgentPools() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAgentPoolsRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"agent_pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_provision": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto_update": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pool_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAgentPoolsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	agentPools, err := clients.TaskAgentClient.GetAgentPools(clients.Ctx, taskagent.GetAgentPoolsArgs{})
	if err != nil {
		return fmt.Errorf(" finding agent pools. Error: %v", err)
	}
	log.Printf("[TRACE] plugin.terraform-provider-azuredevops: Read [%d] agent pools from current organization", len(*agentPools))

	err = d.Set("agent_pools", flattenAgentPoolReferences(agentPools))
	if err != nil {
		return fmt.Errorf(" setting agent_pools field in state. Error: %v", err)
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

		if element.AutoUpdate != nil {
			output["auto_update"] = *element.AutoUpdate
		}

		results = append(results, output)
	}

	return results
}

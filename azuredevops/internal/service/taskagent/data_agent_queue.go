package taskagent

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
	"strconv"
)

// DataAgentQueue schema and implementation for agent queue source
func DataAgentQueue() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAgentQueueRead,
		Schema: map[string]*schema.Schema{
			projectID: {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.NoZeroValues,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceAgentQueueRead(d *schema.ResourceData, m interface{}) error {
	agentQueueName := d.Get("name").(string)
	projectID := d.Get("project_id").(string)
	clients := m.(*client.AggregatedClient)

	agentQueue, err := getAgentQueueByName(clients, &agentQueueName, &projectID)
	if err != nil {
		return fmt.Errorf("Error getting agent queue by name: %v", err)
	}

	flattenAzureAgentQueue(d, agentQueue)
	return nil
}

func flattenAzureAgentQueue(d *schema.ResourceData, agentQueue *taskagent.TaskAgentQueue) {
	d.SetId(strconv.Itoa(*agentQueue.Id))
	d.Set("name", converter.ToString(agentQueue.Name, ""))
	d.Set("pool_id", strconv.Itoa(*agentQueue.Pool.Id))
	//d.Set("project_id", *agentQueue.ProjectId)
}

func getAgentQueueByName(clients *client.AggregatedClient, name, projectID *string) (*taskagent.TaskAgentQueue, error) {
	agentQueues, err := clients.TaskAgentClient.GetAgentQueuesByNames(clients.Ctx, taskagent.GetAgentQueuesByNamesArgs{
		Project:    projectID,
		QueueNames: &[]string{*name},
	})

	if err != nil {
		return nil, err
	}

	if len(*agentQueues) > 1 {
		return nil, fmt.Errorf("Found multiple agent queue for name: %s. Agent queues found: %v", *name, agentQueues)
	}

	if len(*agentQueues) == 0 {
		return nil, fmt.Errorf("Unable to find agent queue with name: %s", *name)
	}

	return &(*agentQueues)[0], nil
}

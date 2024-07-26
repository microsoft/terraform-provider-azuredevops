package taskagent

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
)

// DataAgentQueue schema and implementation for agent queue source
func DataAgentQueue() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAgentQueueRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			projectID: {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.NoZeroValues,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			agentPoolID: {
				Type:     schema.TypeInt,
				Computed: true,
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
	d.Set("name", *agentQueue.Name)
	d.Set(agentPoolID, *agentQueue.Pool.Id)
	d.Set("project_id", agentQueue.ProjectId.String())
}

func getAgentQueueByName(clients *client.AggregatedClient, name, projectID *string) (*taskagent.TaskAgentQueue, error) {
	agentQueues, err := clients.TaskAgentClient.GetAgentQueues(clients.Ctx, taskagent.GetAgentQueuesArgs{
		Project:   projectID,
		QueueName: name,
	})

	if err != nil {
		return nil, err
	}

	if len(*agentQueues) > 1 {
		return nil, fmt.Errorf("Found multiple agent queues for name: %s. Agent queues found: %+v", *name, agentQueues)
	}

	if len(*agentQueues) == 0 {
		return nil, fmt.Errorf("Unable to find agent queues with name: %s", *name)
	}

	return &(*agentQueues)[0], nil
}

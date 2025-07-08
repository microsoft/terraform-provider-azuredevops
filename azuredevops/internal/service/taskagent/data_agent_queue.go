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
			"project_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.IsUUID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"agent_pool_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceAgentQueueRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	agentQueue, err := getAgentQueueByName(clients, d.Get("name").(string), d.Get("project_id").(string))
	if err != nil {
		return fmt.Errorf("getting agent queue by name: %v", err)
	}

	d.SetId(strconv.Itoa(*agentQueue.Id))
	flattenAzureAgentQueue(d, agentQueue)
	return nil
}

func flattenAzureAgentQueue(d *schema.ResourceData, agentQueue *taskagent.TaskAgentQueue) {
	if agentQueue.Name != nil {
		d.Set("name", *agentQueue.Name)
	}

	if agentQueue.Pool != nil && agentQueue.Pool.Id != nil {
		d.Set("agent_pool_id", *agentQueue.Pool.Id)
	}

	if agentQueue.ProjectId != nil {
		d.Set("project_id", agentQueue.ProjectId.String())
	}
}

func getAgentQueueByName(clients *client.AggregatedClient, name, projectID string) (*taskagent.TaskAgentQueue, error) {
	agentQueues, err := clients.TaskAgentClient.GetAgentQueues(clients.Ctx, taskagent.GetAgentQueuesArgs{
		Project:   &projectID,
		QueueName: &name,
	})
	if err != nil {
		return nil, err
	}

	if len(*agentQueues) > 1 {
		return nil, fmt.Errorf("Found multiple agent queues for name: %s. Agent queues found: %+v", name, agentQueues)
	}

	if len(*agentQueues) == 0 {
		return nil, fmt.Errorf("Unable to find agent queues with name: %s", name)
	}

	return &(*agentQueues)[0], nil
}

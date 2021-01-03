package taskagent

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

const (
	agentPoolID                      = "agent_pool_id"
	projectID                        = "project_id"
	invalidQueueIDErrorMessageFormat = "Queue ID was unexpectedly not a valid integer: %+v"
)

// ResourceAgentQueue schema and implementation for agent queue resource
func ResourceAgentQueue() *schema.Resource {
	// Note: there is no update API, so all fields will require a new resource
	return &schema.Resource{
		Create:   resourceAgentQueueCreate,
		Read:     resourceAgentQueueRead,
		Delete:   resourceAgentQueueDelete,
		Importer: tfhelper.ImportProjectQualifiedResourceInteger(),
		Schema: map[string]*schema.Schema{
			agentPoolID: {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			projectID: {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validation.NoZeroValues,
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceAgentQueueCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	queue, projectID, err := expandAgentQueue(d)
	if err != nil {
		return fmt.Errorf("Error expanding the agent queue resource from state: %+v", err)
	}

	referencedPool, err := azureAgentPoolRead(clients, *queue.Pool.Id)
	if err != nil {
		return fmt.Errorf("Error looking up referenced agent pool: %+v", err)
	}

	queue.Name = referencedPool.Name
	createdQueue, err := clients.TaskAgentClient.AddAgentQueue(clients.Ctx, taskagent.AddAgentQueueArgs{
		Queue:              queue,
		Project:            &projectID,
		AuthorizePipelines: converter.Bool(false),
	})

	if err != nil {
		return fmt.Errorf("Error creating agent queue: %+v", err)
	}

	d.SetId(strconv.Itoa(*createdQueue.Id))
	return resourceAgentQueueRead(d, m)
}

func expandAgentQueue(d *schema.ResourceData) (*taskagent.TaskAgentQueue, string, error) {
	queue := &taskagent.TaskAgentQueue{
		Pool: &taskagent.TaskAgentPoolReference{
			Id: converter.Int(d.Get(agentPoolID).(int)),
		},
	}

	if d.Id() != "" {
		id, err := converter.ASCIIToIntPtr(d.Id())
		if err != nil {
			return nil, "", fmt.Errorf(invalidQueueIDErrorMessageFormat, err)
		}
		queue.Id = id
	}

	return queue, d.Get(projectID).(string), nil
}

func resourceAgentQueueRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	queueID, err := converter.ASCIIToIntPtr(d.Id())
	if err != nil {
		return fmt.Errorf(invalidQueueIDErrorMessageFormat, err)
	}

	queue, err := clients.TaskAgentClient.GetAgentQueue(clients.Ctx, taskagent.GetAgentQueueArgs{
		QueueId: queueID,
		Project: converter.String(d.Get(projectID).(string)),
	})

	if utils.ResponseWasNotFound(err) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error reading the agent queue resource: %+v", err)
	}

	if queue.Pool != nil && queue.Pool.Id != nil {
		d.Set(agentPoolID, *queue.Pool.Id)
	}

	return nil
}

func resourceAgentQueueDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	queueID, err := converter.ASCIIToIntPtr(d.Id())
	if err != nil {
		return fmt.Errorf(invalidQueueIDErrorMessageFormat, err)
	}

	err = clients.TaskAgentClient.DeleteAgentQueue(clients.Ctx, taskagent.DeleteAgentQueueArgs{
		QueueId: queueID,
		Project: converter.String(d.Get(projectID).(string)),
	})

	if err != nil {
		return fmt.Errorf("Error deleting agent queue: %+v", err)
	}

	d.SetId("")
	return nil
}

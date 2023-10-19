package taskagent

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

const (
	agentQueueName                   = "name"
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
			agentQueueName: {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				AtLeastOneOf:  []string{agentPoolID, agentQueueName},
				ConflictsWith: []string{agentPoolID},
			},
			agentPoolID: {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				AtLeastOneOf:  []string{agentPoolID, agentQueueName},
				ConflictsWith: []string{agentQueueName},
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

	if queue.Pool != nil {
		referencedPool, err := clients.TaskAgentClient.GetAgentPool(clients.Ctx, taskagent.GetAgentPoolArgs{
			PoolId: queue.Pool.Id,
		})
		if err != nil {
			return fmt.Errorf("Error looking up referenced agent pool: %+v", err)
		}
		queue.Name = referencedPool.Name
	} else {
		queue.Name = converter.String(d.Get(agentQueueName).(string))
	}

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
	queue := &taskagent.TaskAgentQueue{}

	poolId := converter.Int(d.Get(agentPoolID).(int))
	if *poolId != 0 {
		queue.Pool = &taskagent.TaskAgentPoolReference{
			Id: poolId,
		}
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

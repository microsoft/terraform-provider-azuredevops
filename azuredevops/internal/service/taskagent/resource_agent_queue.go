package taskagent

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceAgentQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceAgentQueueCreate,
		Read:   resourceAgentQueueRead,
		Delete: resourceAgentQueueDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceInteger(),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringIsNotWhiteSpace,
				AtLeastOneOf:  []string{"agent_pool_id", "name"},
				ConflictsWith: []string{"agent_pool_id"},
			},
			"agent_pool_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.NoZeroValues,
				AtLeastOneOf:  []string{"agent_pool_id", "name"},
				ConflictsWith: []string{"name"},
			},
			"project_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validation.IsUUID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceAgentQueueCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	queue, projectID, err := expandAgentQueue(d)
	if err != nil {
		return fmt.Errorf(" expanding the agent queue resource from state: %+v", err)
	}

	if queue.Pool != nil {
		referencedPool, err := clients.TaskAgentClient.GetAgentPool(clients.Ctx, taskagent.GetAgentPoolArgs{
			PoolId: queue.Pool.Id,
		})
		if err != nil {
			return fmt.Errorf(" looking up referenced agent pool: %+v", err)
		}
		queue.Name = referencedPool.Name
	} else {
		queue.Name = converter.String(d.Get("name").(string))
	}

	createdQueue, err := clients.TaskAgentClient.AddAgentQueue(clients.Ctx, taskagent.AddAgentQueueArgs{
		Queue:              queue,
		Project:            &projectID,
		AuthorizePipelines: converter.Bool(false),
	})

	if err != nil {
		return fmt.Errorf(" creating agent queue: %+v", err)
	}

	d.SetId(strconv.Itoa(*createdQueue.Id))
	return resourceAgentQueueRead(d, m)
}

func resourceAgentQueueRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	queueID, err := converter.ASCIIToIntPtr(d.Id())
	if err != nil {
		return fmt.Errorf(" Queue ID was unexpectedly not a valid integer: %+v", err)
	}

	queue, err := clients.TaskAgentClient.GetAgentQueue(clients.Ctx, taskagent.GetAgentQueueArgs{
		QueueId: queueID,
		Project: converter.String(d.Get("project_id").(string)),
	})

	if utils.ResponseWasNotFound(err) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf(" reading the agent queue resource: %+v", err)
	}

	if queue.Pool != nil && queue.Pool.Id != nil {
		d.Set("agent_pool_id", *queue.Pool.Id)
	}

	if queue.Name != nil {
		d.Set("name", *queue.Name)
	}

	return nil
}

func resourceAgentQueueDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	queueID, err := converter.ASCIIToIntPtr(d.Id())
	if err != nil {
		return fmt.Errorf(" Queue ID was unexpectedly not a valid integer: %+v", err)
	}

	err = clients.TaskAgentClient.DeleteAgentQueue(clients.Ctx, taskagent.DeleteAgentQueueArgs{
		QueueId: queueID,
		Project: converter.String(d.Get("project_id").(string)),
	})

	if err != nil {
		return fmt.Errorf(" deleting agent queue: %+v", err)
	}

	return nil
}

func expandAgentQueue(d *schema.ResourceData) (*taskagent.TaskAgentQueue, string, error) {
	queue := &taskagent.TaskAgentQueue{}

	if v, ok := d.GetOk("agent_pool_id"); ok {
		queue.Pool = &taskagent.TaskAgentPoolReference{
			Id: converter.Int(v.(int)),
		}
	}

	if d.Id() != "" {
		id, err := converter.ASCIIToIntPtr(d.Id())
		if err != nil {
			return nil, "", fmt.Errorf(" Queue ID was unexpectedly not a valid integer: %+v", err)
		}
		queue.Id = id
	}

	return queue, d.Get("project_id").(string), nil
}

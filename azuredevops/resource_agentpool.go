package azuredevops

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"
)

func resourceAzureAgentPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzureAgentPoolCreate,
		Read:   resourceAzureAgentPoolRead,
		Update: resourceAzureAgentPoolUpdate,
		Delete: resourceAzureAgentPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"pool_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          taskagent.TaskAgentPoolTypeValues.Automation,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(taskagent.TaskAgentPoolTypeValues.Automation),
					string(taskagent.TaskAgentPoolTypeValues.Deployment),
				}, false),
			},
			"auto_provision": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAzureAgentPoolCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	agentPool, err := expandAgentPool(d, true)
	if err != nil {
		return fmt.Errorf("Error converting terraform data model to AzDO agentPool reference: %+v", err)
	}

	createdAgentPool, err := createAzureAgentPool(clients, agentPool)
	if err != nil {
		return fmt.Errorf("Error creating agent pool in Azure DevOps: %+v", err)
	}

	flattenAzureAgentPool(d, createdAgentPool)

	return resourceAzureAgentPoolRead(d, m)
}

func resourceAzureAgentPoolRead(d *schema.ResourceData, m interface{}) error {
	poolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error getting agent pool Id: %+v", err)
	}

	clients := m.(*config.AggregatedClient)
	agentPool, err := azureAgentPoolRead(clients, poolID)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error looking up agent pool with ID %d. Error: %v", poolID, err)
	}

	flattenAzureAgentPool(d, agentPool)
	return nil
}

func resourceAzureAgentPoolUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	agentPool, err := expandAgentPool(d, false)
	if err != nil {
		return fmt.Errorf("Error converting terraform data model to AzDO agent pool reference: %+v", err)
	}

	_, err = azureAgentPoolUpdate(clients, agentPool)
	if err != nil {
		return fmt.Errorf("Error updating agent pool in Azure DevOps: %+v", err)
	}

	return resourceAzureAgentPoolRead(d, m)
}

func resourceAzureAgentPoolDelete(d *schema.ResourceData, m interface{}) error {
	poolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error getting agent pool Id: %+v", err)
	}

	clients := m.(*config.AggregatedClient)
	return clients.TaskAgentClient.DeleteAgentPool(clients.Ctx, taskagent.DeleteAgentPoolArgs{
		PoolId: &poolID,
	})
}

func createAzureAgentPool(clients *config.AggregatedClient, agentPool *taskagent.TaskAgentPool) (*taskagent.TaskAgentPool, error) {
	args := taskagent.AddAgentPoolArgs{
		Pool: agentPool,
	}

	newTaskAgent, err := clients.TaskAgentClient.AddAgentPool(clients.Ctx, args)
	return newTaskAgent, err
}

func azureAgentPoolRead(clients *config.AggregatedClient, poolID int) (*taskagent.TaskAgentPool, error) {
	return clients.TaskAgentClient.GetAgentPool(clients.Ctx, taskagent.GetAgentPoolArgs{
		PoolId: &poolID,
	})
}

func azureAgentPoolUpdate(clients *config.AggregatedClient, agentPool *taskagent.TaskAgentPool) (*taskagent.TaskAgentPool, error) {
	return clients.TaskAgentClient.UpdateAgentPool(
		clients.Ctx,
		taskagent.UpdateAgentPoolArgs{
			PoolId: agentPool.Id,
			Pool: &taskagent.TaskAgentPool{
				Name:          agentPool.Name,
				PoolType:      agentPool.PoolType,
				AutoProvision: agentPool.AutoProvision,
			},
		})
}

func flattenAzureAgentPool(d *schema.ResourceData, agentPool *taskagent.TaskAgentPool) {
	d.SetId(strconv.Itoa(*agentPool.Id))
	d.Set("name", converter.ToString(agentPool.Name, ""))
	d.Set("pool_type", *agentPool.PoolType)
	d.Set("auto_provision", *agentPool.AutoProvision)
}

func expandAgentPool(d *schema.ResourceData, forCreate bool) (*taskagent.TaskAgentPool, error) {
	poolID, err := strconv.Atoi(d.Id())
	if !forCreate && err != nil {
		return nil, fmt.Errorf("Error getting agent pool Id: %+v", err)
	}

	poolType := taskagent.TaskAgentPoolType(d.Get("pool_type").(string))

	pool := &taskagent.TaskAgentPool{
		Id:            &poolID,
		Name:          converter.String(d.Get("name").(string)),
		PoolType:      &poolType,
		AutoProvision: converter.Bool(d.Get("auto_provision").(bool)),
	}

	return pool, nil
}

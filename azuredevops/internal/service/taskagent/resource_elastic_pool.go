package taskagent

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/elastic"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceAgentPoolVMSS() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzureAgentPoolVMSSCreate,
		Read:   resourceAzureAgentPoolVMSSRead,
		Update: resourceAzureAgentPoolVMSSUpdate,
		Delete: resourceAzureAgentPoolVMSSDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},

			"azure_resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"service_endpoint_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"service_endpoint_scope": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"desired_idle": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"max_capacity": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"recycle_after_each_use": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"agent_interactive_ui": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"time_to_live_minutes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"auto_provision": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"auto_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAzureAgentPoolVMSSCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	args := elastic.CreateElasticPoolArgs{
		ElasticPool: &elastic.ElasticPool{
			AgentInteractiveUI:  converter.ToPtr(d.Get("agent_interactive_ui").(bool)),
			AzureId:             converter.ToPtr(d.Get("azure_resource_id").(string)),
			DesiredIdle:         converter.ToPtr(d.Get("desired_idle").(int)),
			MaxCapacity:         converter.ToPtr(d.Get("max_capacity").(int)),
			TimeToLiveMinutes:   converter.ToPtr(d.Get("time_to_live_minutes").(int)),
			RecycleAfterEachUse: converter.ToPtr(d.Get("recycle_after_each_use").(bool)),
		},
		PoolName: converter.String(d.Get("name").(string)),
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectId, err := uuid.Parse(v.(string))
		if err != nil {
			return fmt.Errorf(" parse Project Id: %s. Error: %+v", v, err)
		}
		args.ProjectId = &projectId

	}

	seId := d.Get("service_endpoint_id").(string)
	seUUId, err := uuid.Parse(seId)
	if err != nil {
		return err
	}
	args.ElasticPool.ServiceEndpointId = &seUUId

	seScope := d.Get("service_endpoint_scope").(string)
	seScopeUUId, err := uuid.Parse(seScope)
	if err != nil {
		return err
	}
	args.ElasticPool.ServiceEndpointScope = &seScopeUUId

	desiredIdle := d.Get("desired_idle").(int)
	maxCapacity := d.Get("max_capacity").(int)
	if desiredIdle > maxCapacity {
		return fmt.Errorf(" `desired_idle` can not be greater than `max_capacity`. Valid range is from 0 to the maximum number of virtual machines in the scale set.")
	}

	elasticPool, err := clients.ElasticClient.CreateElasticPool(clients.Ctx, args)
	if err != nil {
		return fmt.Errorf(" creating Elastic Pool: %+v", err)
	}

	updateArgs := taskagent.UpdateAgentPoolArgs{
		PoolId: elasticPool.ElasticPool.PoolId,
		Pool: &taskagent.TaskAgentPool{
			Name:          args.PoolName,
			AutoProvision: converter.ToPtr(d.Get("auto_provision").(bool)),
			AutoUpdate:    converter.ToPtr(d.Get("auto_update").(bool)),
		},
	}
	_, err = clients.TaskAgentClient.UpdateAgentPool(clients.Ctx, updateArgs)
	if err := syncElasticPoolStatus(updateArgs, clients); err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf(" updating agent pool in Azure DevOps: %+v", err)
	}

	d.SetId(strconv.Itoa(*elasticPool.ElasticPool.PoolId))
	return resourceAzureAgentPoolVMSSRead(d, m)
}

func resourceAzureAgentPoolVMSSRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	poolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf(" parse Elastic Pool ID: %+v", err)
	}

	elasticPool, err := clients.ElasticClient.GetElasticPool(clients.Ctx, elastic.GetElasticPoolArgs{
		PoolId: &poolID,
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(" looking up Elastic Pool with ID %d. Error: %v", poolID, err)
	}

	agentPool, err := clients.TaskAgentClient.GetAgentPool(clients.Ctx, taskagent.GetAgentPoolArgs{
		PoolId: &poolID,
	})
	if err != nil {
		return fmt.Errorf(" looking up Agent Pool with ID %d. Error: %v", poolID, err)
	}

	d.Set("name", agentPool.Name)
	d.Set("azure_resource_id", elasticPool.AzureId)
	d.Set("service_endpoint_id", elasticPool.ServiceEndpointId.String())
	d.Set("service_endpoint_scope", elasticPool.ServiceEndpointScope.String())
	d.Set("desired_idle", elasticPool.DesiredIdle)
	d.Set("max_capacity", elasticPool.MaxCapacity)
	d.Set("recycle_after_each_use", elasticPool.RecycleAfterEachUse)
	d.Set("agent_interactive_ui", elasticPool.AgentInteractiveUI)
	d.Set("time_to_live_minutes", elasticPool.TimeToLiveMinutes)
	d.Set("project_id", d.Get("project_id").(string))

	d.Set("auto_provision", agentPool.AutoProvision)
	d.Set("auto_update", agentPool.AutoUpdate)

	return nil
}

func resourceAzureAgentPoolVMSSUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	poolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf(" getting Elastic Pool Id: %+v", err)
	}

	elasticPoolArgs := elastic.UpdateElasticPoolArgs{
		PoolId: &poolID,
		ElasticPoolSettings: &elastic.ElasticPoolSettings{
			AgentInteractiveUI:  converter.ToPtr(d.Get("agent_interactive_ui").(bool)),
			AzureId:             converter.ToPtr(d.Get("azure_resource_id").(string)),
			TimeToLiveMinutes:   converter.ToPtr(d.Get("time_to_live_minutes").(int)),
			RecycleAfterEachUse: converter.ToPtr(d.Get("recycle_after_each_use").(bool)),
		},
	}

	seId := d.Get("service_endpoint_id").(string)
	seUUId, err := uuid.Parse(seId)
	if err != nil {
		return err
	}
	elasticPoolArgs.ElasticPoolSettings.ServiceEndpointId = &seUUId

	seScope := d.Get("service_endpoint_scope").(string)
	seScopeUUId, err := uuid.Parse(seScope)
	if err != nil {
		return err
	}
	elasticPoolArgs.ElasticPoolSettings.ServiceEndpointScope = &seScopeUUId

	desiredIdle := d.Get("desired_idle").(int)
	maxCapacity := d.Get("max_capacity").(int)
	if desiredIdle > maxCapacity {
		return fmt.Errorf(" `desired_idle` can not be greater than `max_capacity`. Valid range is from 0 to the maximum number of virtual machines in the scale set.")
	}

	elasticPoolArgs.ElasticPoolSettings.DesiredIdle = &desiredIdle
	elasticPoolArgs.ElasticPoolSettings.MaxCapacity = &maxCapacity

	if _, err := clients.ElasticClient.UpdateElasticPool(clients.Ctx, elasticPoolArgs); err != nil {
		return fmt.Errorf(" updating Elastic Pool in Azure DevOps: %+v", err)
	}

	agentPoolArgs := taskagent.UpdateAgentPoolArgs{
		Pool: &taskagent.TaskAgentPool{
			Name:          converter.String(d.Get("name").(string)),
			AutoProvision: converter.Bool(d.Get("auto_provision").(bool)),
			AutoUpdate:    converter.Bool(d.Get("auto_update").(bool)),
		},
	}

	agentPoolArgs.PoolId = &poolID
	if _, err = clients.TaskAgentClient.UpdateAgentPool(clients.Ctx, agentPoolArgs); err != nil {
		return fmt.Errorf(" updating Elastic Pool in Azure DevOps: %+v", err)
	}

	if err := syncElasticPoolStatus(agentPoolArgs, clients); err != nil {
		return err
	}

	return resourceAzureAgentPoolVMSSRead(d, m)
}

func resourceAzureAgentPoolVMSSDelete(d *schema.ResourceData, m interface{}) error {
	poolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf(" parse agent pool ID: %+v", err)
	}

	clients := m.(*client.AggregatedClient)
	if err := clients.TaskAgentClient.DeleteAgentPool(clients.Ctx, taskagent.DeleteAgentPoolArgs{PoolId: &poolID}); err != nil {
		return err
	}

	//  waiting resource deleted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Waiting"},
		Target:  []string{"Synched"},
		Refresh: func() (interface{}, string, error) {
			state := "Waiting"
			agentPool, err := clients.TaskAgentClient.GetAgentPool(clients.Ctx, taskagent.GetAgentPoolArgs{PoolId: &poolID})
			if err != nil {
				if utils.ResponseWasNotFound(err) {
					state = "Synched"
				} else {
					return nil, "", fmt.Errorf(" looking up Agent Pool with ID: %+v", err)
				}
			}
			if agentPool == nil {
				state = "Synched"
			}
			return state, state, nil
		},
		Timeout:                   5 * time.Minute,
		MinTimeout:                3 * time.Second,
		Delay:                     1 * time.Second,
		ContinuousTargetOccurence: 2,
	}
	if _, err := stateConf.WaitForStateContext(clients.Ctx); err != nil {
		return fmt.Errorf(" waiting for Elastic Pool deleted. %v ", err)
	}
	return nil
}

func syncElasticPoolStatus(params taskagent.UpdateAgentPoolArgs, client *client.AggregatedClient) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Waiting"},
		Target:  []string{"Synched"},
		Refresh: func() (interface{}, string, error) {
			state := "Waiting"
			agentPool, err := client.TaskAgentClient.GetAgentPool(client.Ctx, taskagent.GetAgentPoolArgs{PoolId: params.PoolId})
			if err != nil {
				return nil, "", fmt.Errorf(" looking up Agent Pool with ID: %+v", err)
			}
			if *agentPool.AutoUpdate == *params.Pool.AutoUpdate &&
				*agentPool.AutoProvision == *params.Pool.AutoProvision {
				state = "Synched"
			}
			return state, state, nil
		},
		Timeout:                   5 * time.Minute,
		MinTimeout:                3 * time.Second,
		Delay:                     1 * time.Second,
		ContinuousTargetOccurence: 2,
	}
	if _, err := stateConf.WaitForStateContext(client.Ctx); err != nil {
		return fmt.Errorf(" waiting for Agent Pool ready. %v ", err)
	}
	return nil
}

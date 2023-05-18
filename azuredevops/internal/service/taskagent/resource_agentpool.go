package taskagent

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
)

// ResourceAgentPool schema and implementation for agent pool resource
func ResourceAgentPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzureAgentPoolCreate,
		Read:   resourceAzureAgentPoolRead,
		Update: resourceAzureAgentPoolUpdate,
		Delete: resourceAzureAgentPoolDelete,
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
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
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
			"auto_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceAzureAgentPoolCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	args := taskagent.AddAgentPoolArgs{
		Pool: &taskagent.TaskAgentPool{
			Name:          converter.String(d.Get("name").(string)),
			PoolType:      converter.ToPtr(taskagent.TaskAgentPoolType(d.Get("pool_type").(string))),
			AutoProvision: converter.Bool(d.Get("auto_provision").(bool)),
			AutoUpdate:    converter.Bool(d.Get("auto_update").(bool)),
		},
	}

	agentPool, err := clients.TaskAgentClient.AddAgentPool(clients.Ctx, args)
	if err != nil {
		return fmt.Errorf(" creating agent pool in Azure DevOps: %+v", err)
	}

	// auto update can only be set to true on creation
	if args.Pool.AutoUpdate != nil && !*args.Pool.AutoUpdate {
		updateArgs := taskagent.UpdateAgentPoolArgs{
			PoolId: agentPool.Id,
			Pool: &taskagent.TaskAgentPool{
				Name:          args.Pool.Name,
				PoolType:      args.Pool.PoolType,
				AutoProvision: args.Pool.AutoProvision,
				AutoUpdate:    args.Pool.AutoUpdate,
			},
		}
		agentPool, err = clients.TaskAgentClient.UpdateAgentPool(clients.Ctx, updateArgs)

		if err := syncStatus(updateArgs, clients); err != nil {
			return err
		}

		if err != nil {
			return fmt.Errorf(" updating agent pool in Azure DevOps: %+v", err)
		}
	}
	d.SetId(strconv.Itoa(*agentPool.Id))
	return resourceAzureAgentPoolRead(d, m)
}

func resourceAzureAgentPoolRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	poolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf(" parse agent pool ID: %+v", err)
	}

	agentPool, err := clients.TaskAgentClient.GetAgentPool(clients.Ctx, taskagent.GetAgentPoolArgs{PoolId: &poolID})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(" looking up Agent Pool with ID %d. Error: %v", poolID, err)
	}

	d.Set("name", agentPool.Name)
	d.Set("pool_type", agentPool.PoolType)
	d.Set("auto_provision", agentPool.AutoProvision)

	if agentPool.AutoUpdate != nil {
		d.Set("auto_update", agentPool.AutoUpdate)
	}
	return nil
}

func resourceAzureAgentPoolUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	parameter := taskagent.UpdateAgentPoolArgs{
		Pool: &taskagent.TaskAgentPool{
			Name:          converter.String(d.Get("name").(string)),
			PoolType:      converter.ToPtr(taskagent.TaskAgentPoolType(d.Get("pool_type").(string))),
			AutoProvision: converter.Bool(d.Get("auto_provision").(bool)),
			AutoUpdate:    converter.Bool(d.Get("auto_update").(bool)),
		},
	}

	poolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf(" getting agent pool Id: %+v", err)
	}
	parameter.PoolId = &poolID

	if _, err = clients.TaskAgentClient.UpdateAgentPool(clients.Ctx, parameter); err != nil {
		return fmt.Errorf(" updating agent pool in Azure DevOps: %+v", err)
	}

	if err := syncStatus(parameter, clients); err != nil {
		return err
	}

	return resourceAzureAgentPoolRead(d, m)
}

func resourceAzureAgentPoolDelete(d *schema.ResourceData, m interface{}) error {
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
		return fmt.Errorf(" waiting for Agent Pool deleted. %v ", err)
	}
	return nil
}

func syncStatus(params taskagent.UpdateAgentPoolArgs, client *client.AggregatedClient) error {
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
				*agentPool.AutoProvision == *params.Pool.AutoProvision &&
				*agentPool.PoolType == *params.Pool.PoolType {
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

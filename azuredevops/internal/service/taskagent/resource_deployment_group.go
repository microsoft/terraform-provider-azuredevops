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
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceDeploymentGroup schema and implementation for deployment group resource
func ResourceDeploymentGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeploymentGroupCreate,
		Read:   resourceDeploymentGroupRead,
		Update: resourceDeploymentGroupUpdate,
		Delete: resourceDeploymentGroupDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceInteger(),
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The ID of the project.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "The name of the deployment group.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The description of the deployment group.",
			},
			"pool_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the deployment pool in which deployment agents are registered. If not specified, a new pool will be created.",
			},
			"machine_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of deployment targets in the deployment group.",
			},
		},
	}
}

func resourceDeploymentGroupCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	createParams := &taskagent.DeploymentGroupCreateParameter{
		Name:        &name,
		Description: &description,
	}

	// If pool_id is specified, use it
	if v, ok := d.GetOk("pool_id"); ok {
		poolId := v.(int)
		createParams.PoolId = &poolId
	}

	deploymentGroup, err := clients.TaskAgentClient.AddDeploymentGroup(clients.Ctx, taskagent.AddDeploymentGroupArgs{
		Project:         &projectID,
		DeploymentGroup: createParams,
	})
	if err != nil {
		return fmt.Errorf("Error creating deployment group: %+v", err)
	}

	if deploymentGroup == nil || deploymentGroup.Id == nil {
		return fmt.Errorf("Error creating deployment group: response or ID is nil")
	}

	d.SetId(strconv.Itoa(*deploymentGroup.Id))

	return resourceDeploymentGroupRead(d, m)
}

func resourceDeploymentGroupRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	deploymentGroupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing deployment group ID: %+v", err)
	}

	deploymentGroup, err := clients.TaskAgentClient.GetDeploymentGroup(clients.Ctx, taskagent.GetDeploymentGroupArgs{
		Project:           &projectID,
		DeploymentGroupId: &deploymentGroupId,
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading deployment group: %+v", err)
	}

	d.Set("project_id", projectID)
	d.Set("name", converter.ToString(deploymentGroup.Name, ""))
	d.Set("description", converter.ToString(deploymentGroup.Description, ""))
	d.Set("machine_count", converter.ToInt(deploymentGroup.MachineCount, 0))

	if deploymentGroup.Pool != nil && deploymentGroup.Pool.Id != nil {
		d.Set("pool_id", *deploymentGroup.Pool.Id)
	}

	return nil
}

func resourceDeploymentGroupUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	deploymentGroupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing deployment group ID: %+v", err)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)

	updateParams := &taskagent.DeploymentGroupUpdateParameter{
		Name:        &name,
		Description: &description,
	}

	_, err = clients.TaskAgentClient.UpdateDeploymentGroup(clients.Ctx, taskagent.UpdateDeploymentGroupArgs{
		Project:           &projectID,
		DeploymentGroupId: &deploymentGroupId,
		DeploymentGroup:   updateParams,
	})
	if err != nil {
		return fmt.Errorf("Error updating deployment group: %+v", err)
	}

	return resourceDeploymentGroupRead(d, m)
}

func resourceDeploymentGroupDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	deploymentGroupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing deployment group ID: %+v", err)
	}

	err = clients.TaskAgentClient.DeleteDeploymentGroup(clients.Ctx, taskagent.DeleteDeploymentGroupArgs{
		Project:           &projectID,
		DeploymentGroupId: &deploymentGroupId,
	})
	if err != nil {
		return fmt.Errorf("Error deleting deployment group: %+v", err)
	}

	d.SetId("")
	return nil
}

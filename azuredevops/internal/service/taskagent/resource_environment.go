package taskagent

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/taskagent/validate"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceEnvironment schema and implementation for environment resource
func ResourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnvironmentCreate,
		Read:   resourceEnvironmentRead,
		Update: resourceEnvironmentUpdate,
		Delete: resourceEnvironmentDelete,
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
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.EnvironmentName,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceEnvironmentCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	environment, err := expandEnvironment(d)
	if err != nil {
		return fmt.Errorf(" Expanding the environment resource from state: %+v", err)
	}

	createdEnvironment, err := createEnvironment(clients, environment)
	if err != nil {
		return fmt.Errorf(" Creating environment in Azure DevOps: %+v", err)
	}

	flattenEnvironment(d, createdEnvironment)
	return resourceEnvironmentRead(d, m)
}

func resourceEnvironmentRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	environmentID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error getting environment Id: %+v", err)
	}

	environment, err := clients.TaskAgentClient.GetEnvironmentById(clients.Ctx, taskagent.GetEnvironmentByIdArgs{
		EnvironmentId: &environmentID,
		Project:       converter.String(d.Get("project_id").(string)),
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading the environment resource: %+v", err)
	}

	flattenEnvironment(d, environment)
	return nil
}

func resourceEnvironmentUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	environment, err := expandEnvironment(d)
	if err != nil {
		return fmt.Errorf("Error converting terraform data model to AzDO environment reference: %+v", err)
	}

	_, err = updateEnvironment(clients, environment)
	if err != nil {
		return fmt.Errorf("Error updating environment in Azure DevOps: %+v", err)
	}

	return resourceEnvironmentRead(d, m)
}

func resourceEnvironmentDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	environmentId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error getting environment id: %+v", err)
	}

	err = clients.TaskAgentClient.DeleteEnvironment(clients.Ctx, taskagent.DeleteEnvironmentArgs{
		Project:       converter.String(d.Get("project_id").(string)),
		EnvironmentId: &environmentId,
	})

	if err != nil {
		return fmt.Errorf("Error deleting environment: %+v", err)
	}

	d.SetId("")
	return nil
}

func createEnvironment(clients *client.AggregatedClient, environment *taskagent.EnvironmentInstance) (*taskagent.EnvironmentInstance, error) {
	return clients.TaskAgentClient.AddEnvironment(
		clients.Ctx,
		taskagent.AddEnvironmentArgs{
			Project: converter.String(environment.Project.Id.String()),
			EnvironmentCreateParameter: &taskagent.EnvironmentCreateParameter{
				Name:        environment.Name,
				Description: environment.Description,
			},
		})
}

func updateEnvironment(clients *client.AggregatedClient, environment *taskagent.EnvironmentInstance) (*taskagent.EnvironmentInstance, error) {
	return clients.TaskAgentClient.UpdateEnvironment(
		clients.Ctx,
		taskagent.UpdateEnvironmentArgs{
			Project:       converter.String(environment.Project.Id.String()),
			EnvironmentId: environment.Id,
			EnvironmentUpdateParameter: &taskagent.EnvironmentUpdateParameter{
				Name:        environment.Name,
				Description: environment.Description,
			},
		})
}

func expandEnvironment(d *schema.ResourceData) (*taskagent.EnvironmentInstance, error) {
	projectId, err := uuid.Parse(d.Get("project_id").(string))
	if err != nil {
		return nil, fmt.Errorf(" faild parse project ID to UUID: %s, %+v", "project_id", err)
	}
	environment := &taskagent.EnvironmentInstance{
		Name:        converter.String(d.Get("name").(string)),
		Description: converter.String(d.Get("description").(string)),
		Project: &taskagent.ProjectReference{
			Id: &projectId,
		},
	}
	// Look for the ID. This may not exist if we are within the context of a "create" operation,
	// so it is OK if it is missing.
	if d.Id() != "" {
		environmentId, err := strconv.Atoi(d.Id())
		if err != nil {
			return nil, fmt.Errorf("Error getting environment id: %+v", err)
		}
		environment.Id = &environmentId
	}
	return environment, nil
}

func flattenEnvironment(d *schema.ResourceData, environment *taskagent.EnvironmentInstance) {
	d.SetId(strconv.Itoa(*environment.Id))
	d.Set("project_id", environment.Project.Id.String())
	d.Set("name", *environment.Name)
	d.Set("description", converter.ToString(environment.Description, ""))
}

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
)

func DataEnvironment() *schema.Resource {
	return &schema.Resource{
		Read: dataEnvironmentRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"environment_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
				ConflictsWith: []string{
					"name",
				},
				AtLeastOneOf: []string{"environment_id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataEnvironmentRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	

	if name == "" {
		environment, err := clients.TaskAgentClient.GetEnvironmentById(clients.Ctx, taskagent.GetEnvironmentByIdArgs{
			EnvironmentId: converter.ToPtr(d.Get("environment_id").(int)),
			Project:       converter.String(d.Get("project_id").(string)),
		})

		if err != nil {
			if utils.ResponseWasNotFound(err) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error reading the environment resource: %+v", err)
		}

		d.SetId(strconv.Itoa(*environment.Id))
		d.Set("project_id", environment.Project.Id.String())
		d.Set("name", *environment.Name)
		d.Set("description", converter.ToString(environment.Description, ""))
		return nil
	} else {
		response, err := clients.TaskAgentClient.GetEnvironments(clients.Ctx, taskagent.GetEnvironmentsArgs{
			Name:    converter.String(d.Get("name").(string)),
			Project: converter.String(d.Get("project_id").(string)),
			Top:     converter.Int(1),
		})
		if err != nil {
			return err
		}

		if len(response.Value) == 0 {
			return fmt.Errorf("Unable to find environment with name: %s", name)
		}

		flattenEnvironment(d, &(response.Value)[0])
		return nil
	}
}

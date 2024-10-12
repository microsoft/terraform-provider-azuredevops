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
)

func DataEnvironment() *schema.Resource {
	return &schema.Resource{
		Read: dataEnvironmentRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
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

	var environment *taskagent.EnvironmentInstance
	var err error
	if envId, ok := d.GetOk("environment_id"); ok {
		environment, err = clients.TaskAgentClient.GetEnvironmentById(clients.Ctx, taskagent.GetEnvironmentByIdArgs{
			EnvironmentId: converter.ToPtr(envId.(int)),
			Project:       converter.String(d.Get("project_id").(string)),
		})
		if err != nil {
			if utils.ResponseWasNotFound(err) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error reading the environment resource: %+v", err)
		}
	} else {
		name := d.Get("name").(string)
		response, err := clients.TaskAgentClient.GetEnvironments(clients.Ctx, taskagent.GetEnvironmentsArgs{
			Name:    converter.String(name),
			Project: converter.String(d.Get("project_id").(string)),
			Top:     converter.Int(1),
		})
		if err != nil {
			return err
		}
		if len(response.Value) == 0 {
			return fmt.Errorf(" Unable to find environment with name: %s", name)
		}
		environment = &response.Value[0]
	}
	d.SetId(strconv.Itoa(*environment.Id))
	d.Set("project_id", environment.Project.Id.String())
	d.Set("name", *environment.Name)
	d.Set("description", converter.ToString(environment.Description, ""))
	return nil
}

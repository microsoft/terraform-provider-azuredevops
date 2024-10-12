package taskagent

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func DataVariableGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVariableGroupRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allow_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"variable": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_value": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"is_secret": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"content_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"expires": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: getVariableHash,
			},
			"key_vault": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVariableGroupRead(d *schema.ResourceData, m interface{}) error {
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)
	clients := m.(*client.AggregatedClient)

	variableGroups, err := clients.TaskAgentClient.GetVariableGroups(clients.Ctx, taskagent.GetVariableGroupsArgs{
		Project:   &projectID,
		GroupName: &name,
		Top:       converter.Int(1),
	})
	if err != nil {
		return err
	}

	if len(*variableGroups) == 0 {
		return fmt.Errorf(" Unable to find variable group with name: %s", name)
	}

	err = flattenVariableGroup(d, &(*variableGroups)[0], &projectID)
	if err != nil {
		return fmt.Errorf(" flattening variable group: %v", err)
	}

	return nil
}

func getVariableHash(v interface{}) int {
	return tfhelper.HashString(v.(map[string]interface{})["name"].(string))
}

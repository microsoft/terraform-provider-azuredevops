package taskagent

import (
	"fmt"

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
		Schema: map[string]*schema.Schema{
			vgProjectID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			vgName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			vgDescription: {
				Type:     schema.TypeString,
				Computed: true,
			},
			vgAllowAccess: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			vgVariable: {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						vgName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						vgValue: {
							Type:     schema.TypeString,
							Computed: true,
						},
						secretVgValue: {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						vgIsSecret: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						vgContentType: {
							Type:     schema.TypeString,
							Computed: true,
						},
						vgEnabled: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						vgExpires: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: getVariableHash,
			},
			vgKeyVault: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						vgName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						vgServiceEndpointID: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func getVariableHash(v interface{}) int {
	return tfhelper.HashString(v.(map[string]interface{})[vgName].(string))
}

func dataSourceVariableGroupRead(d *schema.ResourceData, m interface{}) error {
	projectID := d.Get(vgProjectID).(string)
	name := d.Get(vgName).(string)
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
		return fmt.Errorf("Unable to find variable group with name: %s", name)
	}

	err = flattenVariableGroup(d, &(*variableGroups)[0], &projectID)
	if err != nil {
		return fmt.Errorf(flatteningVariableGroupErrorMessageFormat, err)
	}

	return nil
}

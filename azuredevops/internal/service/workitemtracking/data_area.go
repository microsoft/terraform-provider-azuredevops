package workitemtracking

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func DataArea() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAreaRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func dataSourceAreaRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)

	args := workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Depth:          converter.Int(999),
	}

	path, ok := d.GetOk("path")
	if ok {
		args.Path = converter.String(path.(string))
	}

	area, err := clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, args)
	if err != nil {
		return fmt.Errorf("Error getting Area with path %q: %+v", path, err)
	}

	d.SetId(area.Identifier.String())
	return nil
}

package workitemtrackingprocess

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func DataProcess() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataProcess,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The ID of the process.",
			},
			"expand": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "none",
				ValidateFunc: validation.StringInSlice([]string{"none", "projects"}, false),
				Description:  "Specifies the expand option when getting the processes.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the process.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the process.",
			},
			"parent_process_type_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the parent process.",
			},
			"reference_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference name of process being created. If not specified, server will assign a unique reference name.",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is the process default?",
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is the process enabled?",
			},
			"customization_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the type of customization on this process. System Process is default process. Inherited Process is modified process that was System process before.",
			},
			"projects": {
				Type: schema.TypeSet,
				Set:  getProjectHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the project.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the project.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the project.",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Url of the project.",
						},
					},
				},
				Computed:    true,
				Description: "Returns associated projects when using the 'projects' expand option.",
			},
		},
	}
}

var getProcessExpandLevelMap = map[string]workitemtrackingprocess.GetProcessExpandLevel{
	"none":     workitemtrackingprocess.GetProcessExpandLevelValues.None,
	"projects": workitemtrackingprocess.GetProcessExpandLevelValues.Projects,
}

func readDataProcess(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	processId := d.Get("id").(string)

	expand := getProcessExpandLevelMap[d.Get("expand").(string)]

	getProcessArgs := workitemtrackingprocess.GetProcessByItsIdArgs{
		ProcessTypeId: converter.UUID(processId),
		Expand:        &expand,
	}
	process, err := clients.WorkItemTrackingProcessClient.GetProcessByItsId(ctx, getProcessArgs)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(" Getting process with id: %s. Error: %+v", processId, err)
	}

	d.Set("name", process.Name)
	d.Set("description", process.Description)
	d.Set("parent_process_type_id", process.ParentProcessTypeId.String())
	d.Set("reference_name", process.ReferenceName)
	d.Set("is_default", process.IsDefault)
	d.Set("is_enabled", process.IsEnabled)
	d.Set("customization_type", string(*process.CustomizationType))

	if process.Projects != nil {
		var projects []map[string]any
		for _, p := range *process.Projects {
			project := map[string]any{
				"id":          p.Id.String(),
				"description": p.Description,
				"name":        p.Name,
				"url":         p.Url,
			}
			projects = append(projects, project)
		}
		d.Set("projects", projects)
	}

	d.SetId(processId)

	return nil
}

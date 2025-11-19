package workitemtrackingprocess

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

func DataProcesses() *schema.Resource {
	return &schema.Resource{
		ReadContext: readProcesses,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the process",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the process",
			},
			"parent_process_type_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the parent process",
			},
			"reference_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference name of process being created. If not specified, server will assign a unique reference name",
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
				Description: "Indicates the type of customization on this process. System Process is default process. Inherited Process is modified process that was System process before",
			},
			"expand": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "none",
				ValidateFunc: validation.StringInSlice([]string{"none", "projects"}, false),
				Description:  "Specifies the expand option when getting the processes",
			},
			"projects": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the project",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the project",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the project",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Url of the project",
						},
					},
				},
				Computed:    true,
				Description: "Returns associated projects when using the 'projects' expand option",
			},
		},
	}
}

func readProcesses(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	expandInput := d.Get("expand").(string)
	expand := getProcessExpandLevelMap[expandInput]

	getListOfProcessesArgs := workitemtrackingprocess.GetListOfProcessesArgs{
		Expand: &expand,
	}
	retrievedProcesses, err := clients.WorkItemTrackingProcessClient.GetListOfProcesses(ctx, getListOfProcessesArgs)
	if err != nil {
		return diag.Errorf(" Getting list of processes: Error: %+v", err)
	}

	processes := make([]any, 0)
	for _, retrievedProcess := range *retrievedProcesses {
		process := make(map[string]any)
		if retrievedProcess.Name != nil {
			process["name"] = *retrievedProcess.Name
		}
		if retrievedProcess.Description != nil {
			process["description"] = *retrievedProcess.Description
		}
		if retrievedProcess.ParentProcessTypeId != nil {
			process["parent_process_type_id"] = retrievedProcess.ParentProcessTypeId.String()
		}
		if retrievedProcess.ReferenceName != nil {
			process["reference_name"] = *retrievedProcess.ReferenceName
		}
		if retrievedProcess.IsDefault != nil {
			process["is_default"] = *retrievedProcess.IsDefault
		}
		if retrievedProcess.IsEnabled != nil {
			process["is_enabled"] = *retrievedProcess.IsEnabled
		}
		if retrievedProcess.CustomizationType != nil {
			process["customization_type"] = string(*retrievedProcess.CustomizationType)
		}

		if retrievedProcess.Projects != nil {
			projects := make([]any, 0)
			for _, retrievedProject := range *retrievedProcess.Projects {
				project := make(map[string]any)
				if retrievedProject.Id != nil {
					project["id"] = retrievedProject.Id.String()
				}
				if retrievedProject.Description != nil {
					project["description"] = *retrievedProject.Description
				}
				if retrievedProject.Name != nil {
					project["name"] = *retrievedProject.Name
				}
				if retrievedProject.Url != nil {
					project["url"] = *retrievedProject.Url
				}
				projects = append(projects, project)
			}
		}
		processes = append(processes, process)
	}

	// Expand is the only input to the data source query
	d.SetId(expandInput)

	err = d.Set("processes", processes)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

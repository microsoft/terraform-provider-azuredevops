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

func ResourceProcess() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceProcess,
		ReadContext:   readResourceProcess,
		UpdateContext: updateResourceProcess,
		DeleteContext: deleteResourceProcess,
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
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Name of the process.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of the process.",
			},
			"parent_process_type_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "ID of the parent process.",
			},
			"reference_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Reference name of process being created. If not specified, server will assign a unique reference name.",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Is the process default?",
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Is the process enabled?",
			},
			"customization_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the type of customization on this process. System Process is default process. Inherited Process is modified process that was System process before.",
			},
		},
	}
}

func createResourceProcess(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	createProcessModel := &workitemtrackingprocess.CreateProcessModel{
		Name:                converter.String(d.Get("name").(string)),
		Description:         converter.String(d.Get("description").(string)),
		ParentProcessTypeId: converter.UUID(d.Get("parent_process_type_id").(string)),
	}
	if referenceName, ok := d.GetOk("reference_name"); ok {
		createProcessModel.ReferenceName = converter.String(referenceName.(string))
	}

	args := workitemtrackingprocess.CreateNewProcessArgs{
		CreateRequest: createProcessModel,
	}
	createdProcess, err := clients.WorkItemTrackingProcessClient.CreateNewProcess(ctx, args)
	if err != nil {
		return diag.Errorf(" Creating process. Error %+v", err)
	}

	d.SetId(createdProcess.TypeId.String())

	isDefault := d.Get("is_default").(bool)
	isEnabled := d.Get("is_enabled").(bool)
	if createdProcess.IsDefault == nil || *createdProcess.IsDefault != isDefault ||
		createdProcess.IsEnabled == nil || *createdProcess.IsEnabled != isEnabled {
		return updateResourceProcess(ctx, d, m)
	}

	return flattenProcess(d, createdProcess)
}

func readResourceProcess(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	processId := d.Id()

	getProcessArgs := workitemtrackingprocess.GetProcessByItsIdArgs{
		ProcessTypeId: converter.UUID(processId),
		Expand:        &workitemtrackingprocess.GetProcessExpandLevelValues.None,
	}
	process, err := clients.WorkItemTrackingProcessClient.GetProcessByItsId(ctx, getProcessArgs)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(" Getting process with id: %s. Error: %+v", processId, err)
	}

	return flattenProcess(d, process)
}

func updateResourceProcess(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	updateProcessModel := &workitemtrackingprocess.UpdateProcessModel{
		Name:        converter.String(d.Get("name").(string)),
		Description: converter.String(d.Get("description").(string)),
		IsDefault:   converter.Bool(d.Get("is_default").(bool)),
		IsEnabled:   converter.Bool(d.Get("is_enabled").(bool)),
	}

	args := workitemtrackingprocess.EditProcessArgs{
		ProcessTypeId: converter.UUID(d.Id()),
		UpdateRequest: updateProcessModel,
	}
	updatedProcess, err := clients.WorkItemTrackingProcessClient.EditProcess(ctx, args)
	if err != nil {
		return diag.Errorf(" Update process. Error %+v", err)
	}

	return flattenProcess(d, updatedProcess)
}

func deleteResourceProcess(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	args := workitemtrackingprocess.DeleteProcessByIdArgs{
		ProcessTypeId: converter.UUID(d.Id()),
	}

	err := clients.WorkItemTrackingProcessClient.DeleteProcessById(ctx, args)
	if err != nil {
		return diag.Errorf(" Delete process. Error %+v", err)
	}

	return nil
}

func flattenProcess(d *schema.ResourceData, process *workitemtrackingprocess.ProcessInfo) diag.Diagnostics {
	if process.Name != nil {
		d.Set("name", process.Name)
	}
	if process.Description != nil {
		d.Set("description", process.Description)
	}
	if process.ParentProcessTypeId != nil {
		d.Set("parent_process_type_id", process.ParentProcessTypeId.String())
	}
	if process.ReferenceName != nil {
		d.Set("reference_name", process.ReferenceName)
	}
	if process.IsDefault != nil {
		d.Set("is_default", process.IsDefault)
	}
	if process.IsEnabled != nil {
		d.Set("is_enabled", process.IsEnabled)
	}
	if process.CustomizationType != nil {
		d.Set("customization_type", string(*process.CustomizationType))
	}

	return nil
}

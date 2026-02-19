package workitemtrackingprocess

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceField() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFieldCreate,
		ReadContext:   resourceFieldRead,
		UpdateContext: resourceFieldUpdate,
		DeleteContext: resourceFieldDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importField,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"process_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Description:      "The ID of the process.",
			},
			"work_item_type_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The ID (reference name) of the work item type.",
			},
			"field_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The ID (reference name) of the field.",
			},
			"default_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The default value of the field.",
			},
			"read_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the field cannot be edited.",
			},
			"required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the field cannot be empty.",
			},
			// We set this to write-only to circumvent this bug: https://developercommunity.visualstudio.com/t/Custom-field-APIs-are-missing-attributes/11032086
			"allow_groups": {
				Type:        schema.TypeBool,
				Optional:    true,
				WriteOnly:   true,
				Description: "Allow setting field value to a group identity. Only applies to identity fields.",
			},
			"customization": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the type of customization on this work item. Possible values are `system`, `inherited`, or `custom`.",
			},
			"is_locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the field definition is locked for editing.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the field resource.",
			},
		},
	}
}

func resourceFieldCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_id").(string)
	referenceName := d.Get("field_id").(string)

	fieldRequest := &workitemtrackingprocess.AddProcessWorkItemTypeFieldRequest{
		ReferenceName: &referenceName,
		ReadOnly:      converter.Bool(d.Get("read_only").(bool)),
		Required:      converter.Bool(d.Get("required").(bool)),
	}

	if v, ok := d.GetOk("default_value"); ok {
		fieldRequest.DefaultValue = v.(string)
	}
	rawConfig := d.GetRawConfig().AsValueMap()
	if allowGroups := rawConfig["allow_groups"]; !allowGroups.IsNull() {
		fieldRequest.AllowGroups = converter.Bool(allowGroups.True())
	}

	args := workitemtrackingprocess.AddFieldToWorkItemTypeArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: &witRefName,
		Field:      fieldRequest,
	}

	createdField, err := clients.WorkItemTrackingProcessClient.AddFieldToWorkItemType(ctx, args)
	if err != nil {
		return diag.Errorf("adding field %s to work item type %s: %+v", referenceName, witRefName, err)
	}

	if createdField.ReferenceName == nil {
		return diag.Errorf("created field has no reference name")
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", processId, witRefName, *createdField.ReferenceName))
	return resourceFieldRead(ctx, d, m)
}

func resourceFieldRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_id").(string)
	fieldRefName := d.Get("field_id").(string)

	args := workitemtrackingprocess.GetWorkItemTypeFieldArgs{
		ProcessId:    converter.UUID(processId),
		WitRefName:   &witRefName,
		FieldRefName: &fieldRefName,
		Expand:       &workitemtrackingprocess.ProcessWorkItemTypeFieldsExpandLevelValues.All,
	}

	field, err := clients.WorkItemTrackingProcessClient.GetWorkItemTypeField(ctx, args)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("reading field %s: %+v", fieldRefName, err)
	}
	if field == nil {
		return diag.Errorf("field %s returned nil", fieldRefName)
	}

	if field.ReferenceName != nil {
		d.Set("field_id", *field.ReferenceName)
	}
	if field.DefaultValue != nil {
		d.Set("default_value", fmt.Sprintf("%v", field.DefaultValue))
	}
	if field.ReadOnly != nil {
		d.Set("read_only", *field.ReadOnly)
	} else {
		d.Set("read_only", false)
	}
	if field.Required != nil {
		d.Set("required", *field.Required)
	} else {
		d.Set("required", false)
	}
	if field.Customization != nil {
		d.Set("customization", string(*field.Customization))
	}
	if field.IsLocked != nil {
		d.Set("is_locked", *field.IsLocked)
	}
	if field.Url != nil {
		d.Set("url", *field.Url)
	}
	return nil
}

func resourceFieldUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_id").(string)
	fieldRefName := d.Get("field_id").(string)

	fieldUpdate := &workitemtrackingprocess.UpdateProcessWorkItemTypeFieldRequest{
		ReadOnly: converter.Bool(d.Get("read_only").(bool)),
		Required: converter.Bool(d.Get("required").(bool)),
	}

	if v, ok := d.GetOk("default_value"); ok {
		fieldUpdate.DefaultValue = v.(string)
	}
	rawConfig := d.GetRawConfig().AsValueMap()
	if allowGroups := rawConfig["allow_groups"]; !allowGroups.IsNull() {
		fieldUpdate.AllowGroups = converter.Bool(allowGroups.True())
	}

	args := workitemtrackingprocess.UpdateWorkItemTypeFieldArgs{
		ProcessId:    converter.UUID(processId),
		WitRefName:   &witRefName,
		FieldRefName: &fieldRefName,
		Field:        fieldUpdate,
	}

	_, err := clients.WorkItemTrackingProcessClient.UpdateWorkItemTypeField(ctx, args)
	if err != nil {
		return diag.Errorf("updating field %s: %+v", fieldRefName, err)
	}

	return resourceFieldRead(ctx, d, m)
}

func resourceFieldDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_id").(string)
	fieldRefName := d.Get("field_id").(string)

	args := workitemtrackingprocess.RemoveWorkItemTypeFieldArgs{
		ProcessId:    converter.UUID(processId),
		WitRefName:   &witRefName,
		FieldRefName: &fieldRefName,
	}

	err := clients.WorkItemTrackingProcessClient.RemoveWorkItemTypeField(ctx, args)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return nil
		}
		return diag.Errorf("removing field %s: %+v", fieldRefName, err)
	}

	return nil
}

func importField(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts, err := tfhelper.ParseImportedNameParts(d.Id(), "process_id/work_item_type_id/field_id", 3)
	if err != nil {
		return nil, err
	}
	d.Set("process_id", parts[0])
	d.Set("work_item_type_id", parts[1])
	d.Set("field_id", parts[2])
	return []*schema.ResourceData{d}, nil
}

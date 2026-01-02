package workitemtrackingprocess

import (
	"context"
	"encoding/json"
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
			"work_item_type_ref_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The reference name of the work item type.",
			},
			"reference_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The reference name of the field.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the field.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the field.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the field.",
			},
			"default_value_json": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The default value of the field as a JSON-encoded string.",
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
			"allow_groups": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow setting field value to a group identity. Only applies to identity fields.",
			},
			"allowed_values": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of field allowed values.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	witRefName := d.Get("work_item_type_ref_name").(string)
	referenceName := d.Get("reference_name").(string)

	fieldRequest := &workitemtrackingprocess.AddProcessWorkItemTypeFieldRequest{
		ReferenceName: &referenceName,
		ReadOnly:      converter.Bool(d.Get("read_only").(bool)),
		Required:      converter.Bool(d.Get("required").(bool)),
		AllowGroups:   converter.Bool(d.Get("allow_groups").(bool)),
	}

	if v, ok := d.GetOk("default_value_json"); ok {
		fieldRequest.DefaultValue = expandDefaultValue(v.(string))
	}

	if v, ok := d.GetOk("allowed_values"); ok {
		allowedValues := expandAllowedValues(v.([]interface{}))
		fieldRequest.AllowedValues = &allowedValues
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

	d.SetId(*createdField.ReferenceName)
	return resourceFieldRead(ctx, d, m)
}

func resourceFieldRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_ref_name").(string)
	fieldRefName := d.Id()

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

	flattenProcessField(d, field)
	return nil
}

func resourceFieldUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_ref_name").(string)
	fieldRefName := d.Id()

	fieldUpdate := &workitemtrackingprocess.UpdateProcessWorkItemTypeFieldRequest{
		ReadOnly:    converter.Bool(d.Get("read_only").(bool)),
		Required:    converter.Bool(d.Get("required").(bool)),
		AllowGroups: converter.Bool(d.Get("allow_groups").(bool)),
	}

	if v, ok := d.GetOk("default_value_json"); ok {
		fieldUpdate.DefaultValue = expandDefaultValue(v.(string))
	}

	if v, ok := d.GetOk("allowed_values"); ok {
		allowedValues := expandAllowedValues(v.([]interface{}))
		fieldUpdate.AllowedValues = &allowedValues
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
	witRefName := d.Get("work_item_type_ref_name").(string)
	fieldRefName := d.Id()

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

func expandAllowedValues(values []interface{}) []string {
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = v.(string)
	}
	return result
}

// expandDefaultValue parses a JSON string from Terraform into a value for the API.
func expandDefaultValue(jsonStr string) interface{} {
	if jsonStr == "" {
		return nil
	}
	var value interface{}
	if err := json.Unmarshal([]byte(jsonStr), &value); err != nil {
		// If it's not valid JSON, treat it as a plain string
		return jsonStr
	}
	return value
}

// flattenDefaultValue converts the API's default value to a JSON string for Terraform state.
func flattenDefaultValue(value interface{}) string {
	if value == nil {
		return ""
	}
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}
	return string(jsonBytes)
}

func importField(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts, err := tfhelper.ParseImportedNameParts(d.Id(), "process_id/work_item_type_ref_name/field_ref_name", 3)
	if err != nil {
		return nil, err
	}
	d.Set("process_id", parts[0])
	d.Set("work_item_type_ref_name", parts[1])
	d.SetId(parts[2])
	return []*schema.ResourceData{d}, nil
}

func flattenProcessField(d *schema.ResourceData, field *workitemtrackingprocess.ProcessWorkItemTypeField) {
	if field.Name != nil {
		d.Set("name", *field.Name)
	}
	if field.ReferenceName != nil {
		d.Set("reference_name", *field.ReferenceName)
	}
	if field.Type != nil {
		d.Set("type", string(*field.Type))
	}
	if field.Description != nil {
		d.Set("description", *field.Description)
	}
	if field.DefaultValue != nil {
		d.Set("default_value_json", flattenDefaultValue(field.DefaultValue))
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
	if field.AllowGroups != nil {
		d.Set("allow_groups", *field.AllowGroups)
	} else {
		d.Set("allow_groups", false)
	}
	if field.AllowedValues != nil {
		allowedValues := make([]string, len(*field.AllowedValues))
		for i, v := range *field.AllowedValues {
			if strVal, ok := v.(string); ok {
				allowedValues[i] = strVal
			}
		}
		d.Set("allowed_values", allowedValues)
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
}

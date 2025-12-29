package workitemtracking

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceField() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFieldCreate,
		ReadContext:   resourceFieldRead,
		UpdateContext: resourceFieldUpdate,
		DeleteContext: resourceFieldDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceFieldImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The friendly name of the field.",
			},
			"reference_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					validation.StringIsNotWhiteSpace,
					validation.StringLenBetween(1, 128),
					validation.StringDoesNotMatch(regexp.MustCompile(`[,;~:/\\*|?"&%$!+=()[\]{}<>\-]`), "cannot contain the following characters: ',;~:/\\*|?\"&%$!+=()[]{}<>-'"),
				)),
				Description: "The reference name of the field (e.g., Custom.MyField).",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					"string",
					"integer",
					"dateTime",
					"plainText",
					"html",
					"treePath",
					"history",
					"double",
					"guid",
					"boolean",
					"identity",
				}, false)),
				Description: "The type of the field. Possible values: `string`, `integer`, `dateTime`, `plainText`, `html`, `treePath`, `history`, `double`, `guid`, `boolean`, `identity`.",
			},
			"project_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Description:      "The ID of the project. If not specified, the field is created at the organization level.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The description of the field.",
			},
			"usage": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "workItem",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					"none",
					"workItem",
					"workItemLink",
					"tree",
					"workItemTypeExtension",
				}, false)),
				Description: "The usage of the field. Possible values: `none`, `workItem`, `workItemLink`, `tree`, `workItemTypeExtension`. Default: `workItem`.",
			},
			"read_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Indicates whether the field is read-only. Default: `false`.",
			},
			"can_sort_by": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     true,
				Description: "Indicates whether the field can be sorted in server queries. Default: `true`.",
			},
			"is_queryable": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     true,
				Description: "Indicates whether the field can be queried in the server. Default: `true`.",
			},
			"is_identity": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Indicates whether this field is an identity field. Default: `false`.",
			},
			"is_picklist": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Indicates whether this field is a picklist. Default: `false`.",
			},
			"is_picklist_suggested": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Indicates whether this field is a suggested picklist. Default: `false`.",
			},
			"picklist_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Description:      "The identifier of the picklist associated with this field, if applicable.",
			},
			"is_locked": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether this field is locked for editing. Default: `false`.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the field resource.",
			},
			"is_deleted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether this field is deleted. Default: `false`.",
			},
			"supported_operations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The supported operations on this field.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The friendly name of the operation.",
						},
						"reference_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reference name of the operation.",
						},
					},
				},
			},
		},
	}
}

func resourceFieldCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	field := expandField(d)

	var project *string
	if v, ok := d.GetOk("project_id"); ok {
		project = converter.String(v.(string))
	}

	args := workitemtracking.CreateWorkItemFieldArgs{
		WorkItemField: field,
		Project:       project,
	}

	createdField, err := clients.WorkItemTrackingClient.CreateWorkItemField(clients.Ctx, args)
	if err != nil {
		return diag.Errorf("creating field: %+v", err)
	}

	if createdField.ReferenceName == nil {
		return diag.Errorf("created field has no reference name")
	}

	d.SetId(*createdField.ReferenceName)
	return resourceFieldRead(ctx, d, m)
}

func resourceFieldRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	referenceName := d.Id()

	var project *string
	if v, ok := d.GetOk("project_id"); ok {
		project = converter.String(v.(string))
	}

	args := workitemtracking.GetWorkItemFieldArgs{
		FieldNameOrRefName: &referenceName,
		Project:            project,
	}

	field, err := clients.WorkItemTrackingClient.GetWorkItemField(clients.Ctx, args)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("reading field %s: %+v", referenceName, err)
	}

	flattenField(d, field)
	return nil
}

func resourceFieldUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChange("is_locked") && d.HasChange("is_deleted") {
		return diag.Errorf("cannot update is_locked and is_deleted at the same time")
	}

	clients := m.(*client.AggregatedClient)
	referenceName := d.Id()

	var project *string
	if v, ok := d.GetOk("project_id"); ok {
		project = converter.String(v.(string))
	}

	if d.HasChange("is_locked") {
		isLocked := d.Get("is_locked").(bool)
		args := workitemtracking.UpdateWorkItemFieldArgs{
			FieldNameOrRefName: &referenceName,
			Project:            project,
			Payload: &workitemtracking.FieldUpdate{
				IsLocked: &isLocked,
			},
		}

		_, err := clients.WorkItemTrackingClient.UpdateWorkItemField(clients.Ctx, args)
		if err != nil {
			return diag.Errorf("updating field %s: %+v", referenceName, err)
		}
	}

	if d.HasChange("is_deleted") {
		isDeleted := d.Get("is_deleted").(bool)
		args := workitemtracking.UpdateWorkItemFieldArgs{
			FieldNameOrRefName: &referenceName,
			Project:            project,
			Payload: &workitemtracking.FieldUpdate{
				IsDeleted: &isDeleted,
			},
		}

		_, err := clients.WorkItemTrackingClient.UpdateWorkItemField(clients.Ctx, args)
		if err != nil {
			return diag.Errorf("updating field %s: %+v", referenceName, err)
		}
	}

	return resourceFieldRead(ctx, d, m)
}

func resourceFieldDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	referenceName := d.Id()

	var project *string
	if v, ok := d.GetOk("project_id"); ok {
		project = converter.String(v.(string))
	}

	args := workitemtracking.DeleteWorkItemFieldArgs{
		FieldNameOrRefName: &referenceName,
		Project:            project,
	}

	err := clients.WorkItemTrackingClient.DeleteWorkItemField(clients.Ctx, args)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return nil
		}
		return diag.Errorf("deleting field %s: %+v", referenceName, err)
	}

	return nil
}

func resourceFieldImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	switch len(parts) {
	case 1:
		d.SetId(parts[0])
	case 2:
		d.Set("project_id", parts[0])
		d.SetId(parts[1])
	default:
		return nil, fmt.Errorf("invalid import ID format, expected: reference_name or project_id/reference_name")
	}

	return []*schema.ResourceData{d}, nil
}

func expandField(d *schema.ResourceData) *workitemtracking.WorkItemField2 {
	fieldType := workitemtracking.FieldType(d.Get("type").(string))

	field := &workitemtracking.WorkItemField2{
		Name:                converter.String(d.Get("name").(string)),
		ReferenceName:       converter.String(d.Get("reference_name").(string)),
		Type:                &fieldType,
		ReadOnly:            converter.Bool(d.Get("read_only").(bool)),
		CanSortBy:           converter.Bool(d.Get("can_sort_by").(bool)),
		IsDeleted:           converter.Bool(d.Get("is_deleted").(bool)),
		IsQueryable:         converter.Bool(d.Get("is_queryable").(bool)),
		IsIdentity:          converter.Bool(d.Get("is_identity").(bool)),
		IsPicklist:          converter.Bool(d.Get("is_picklist").(bool)),
		IsPicklistSuggested: converter.Bool(d.Get("is_picklist_suggested").(bool)),
		IsLocked:            converter.Bool(d.Get("is_locked").(bool)),
	}

	if v, ok := d.GetOk("usage"); ok {
		fieldUsage := workitemtracking.FieldUsage(v.(string))
		field.Usage = &fieldUsage
	}

	if v, ok := d.GetOk("description"); ok {
		field.Description = converter.String(v.(string))
	}

	if v, ok := d.GetOk("picklist_id"); ok {
		picklistId, _ := uuid.Parse(v.(string))
		field.PicklistId = &picklistId
	}

	return field
}

func flattenField(d *schema.ResourceData, field *workitemtracking.WorkItemField2) {
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
	if field.Usage != nil {
		d.Set("usage", string(*field.Usage))
	}
	if field.ReadOnly != nil {
		d.Set("read_only", *field.ReadOnly)
	}
	if field.CanSortBy != nil {
		d.Set("can_sort_by", *field.CanSortBy)
	}
	if field.IsQueryable != nil {
		d.Set("is_queryable", *field.IsQueryable)
	}
	if field.IsIdentity != nil {
		d.Set("is_identity", *field.IsIdentity)
	}
	if field.IsPicklist != nil {
		d.Set("is_picklist", *field.IsPicklist)
	}
	if field.IsPicklistSuggested != nil {
		d.Set("is_picklist_suggested", *field.IsPicklistSuggested)
	}
	if field.PicklistId != nil {
		d.Set("picklist_id", field.PicklistId.String())
	}
	if field.IsLocked != nil {
		d.Set("is_locked", *field.IsLocked)
	}
	if field.Url != nil {
		d.Set("url", *field.Url)
	}
	if field.IsDeleted != nil {
		d.Set("is_deleted", *field.IsDeleted)
	} else {
		d.Set("is_deleted", false)
	}

	if field.SupportedOperations != nil {
		operations := make([]map[string]interface{}, len(*field.SupportedOperations))
		for i, op := range *field.SupportedOperations {
			operation := map[string]interface{}{}
			if op.Name != nil {
				operation["name"] = *op.Name
			}
			if op.ReferenceName != nil {
				operation["reference_name"] = *op.ReferenceName
			}
			operations[i] = operation
		}
		d.Set("supported_operations", operations)
	}
}

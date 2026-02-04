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
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

var conditionTypes = []string{
	"when",
	"whenNot",
	"whenChanged",
	"whenNotChanged",
	"whenWas",
	"whenCurrentUserIsMemberOfGroup",
	"whenCurrentUserIsNotMemberOfGroup",
}

var actionTypes = []string{
	"makeRequired",
	"makeReadOnly",
	"setDefaultValue",
	"setDefaultFromClock",
	"setDefaultFromCurrentUser",
	"setDefaultFromField",
	"copyValue",
	"copyFromClock",
	"copyFromCurrentUser",
	"copyFromField",
	"setValueToEmpty",
	"copyFromServerClock",
	"copyFromServerCurrentUser",
	"hideTargetField",
	"disallowValue",
}

func ResourceRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRuleCreate,
		ReadContext:   resourceRuleRead,
		UpdateContext: resourceRuleUpdate,
		DeleteContext: resourceRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importRule,
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
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "Name of the rule.",
			},
			"is_disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates if the rule is disabled.",
			},
			"condition": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Description: "Set of conditions when the rule should be triggered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition_type": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(conditionTypes, false)),
							Description:      "Type of condition.",
						},
						"field": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Field reference name for the condition.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Value to match for the condition.",
						},
					},
				},
			},
			"action": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Description: "Set of actions to take when the rule is triggered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_type": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(actionTypes, false)),
							Description:      "Type of action.",
						},
						"target_field": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to act on.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Value to set on the target field.",
						},
					},
				},
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of the rule resource.",
			},
		},
	}
}

func resourceRuleCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_id").(string)

	ruleRequest := &workitemtrackingprocess.CreateProcessRuleRequest{
		Name:       converter.String(d.Get("name").(string)),
		IsDisabled: converter.Bool(d.Get("is_disabled").(bool)),
		Conditions: expandConditions(d.Get("condition").(*schema.Set).List()),
		Actions:    expandActions(d.Get("action").(*schema.Set).List()),
	}

	args := workitemtrackingprocess.AddProcessWorkItemTypeRuleArgs{
		ProcessId:         converter.UUID(processId),
		WitRefName:        &witRefName,
		ProcessRuleCreate: ruleRequest,
	}

	createdRule, err := clients.WorkItemTrackingProcessClient.AddProcessWorkItemTypeRule(ctx, args)
	if err != nil {
		return diag.Errorf(" Creating rule. Error: %+v", err)
	}

	if createdRule == nil {
		return diag.Errorf(" Created rule is nil")
	}

	if createdRule.Id == nil {
		return diag.Errorf(" Created rule has no ID")
	}

	d.SetId(createdRule.Id.String())
	return resourceRuleRead(ctx, d, m)
}

func resourceRuleRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_id").(string)
	ruleId := d.Id()

	args := workitemtrackingprocess.GetProcessWorkItemTypeRuleArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: &witRefName,
		RuleId:     converter.UUID(ruleId),
	}

	rule, err := clients.WorkItemTrackingProcessClient.GetProcessWorkItemTypeRule(ctx, args)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(" Reading rule %s. Error: %+v", ruleId, err)
	}

	if rule == nil {
		return diag.Errorf(" Rule %s is nil", ruleId)
	}

	if rule.Name != nil {
		d.Set("name", *rule.Name)
	}
	if rule.IsDisabled != nil {
		d.Set("is_disabled", *rule.IsDisabled)
	}
	if rule.Url != nil {
		d.Set("url", *rule.Url)
	}
	if rule.Conditions != nil {
		if err := d.Set("condition", flattenConditions(*rule.Conditions)); err != nil {
			return diag.FromErr(err)
		}
	}
	if rule.Actions != nil {
		if err := d.Set("action", flattenActions(*rule.Actions)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceRuleUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_id").(string)
	ruleId := d.Id()

	ruleUpdate := &workitemtrackingprocess.UpdateProcessRuleRequest{
		Id:         converter.UUID(ruleId),
		Name:       converter.String(d.Get("name").(string)),
		IsDisabled: converter.Bool(d.Get("is_disabled").(bool)),
		Conditions: expandConditions(d.Get("condition").(*schema.Set).List()),
		Actions:    expandActions(d.Get("action").(*schema.Set).List()),
	}

	args := workitemtrackingprocess.UpdateProcessWorkItemTypeRuleArgs{
		ProcessId:   converter.UUID(processId),
		WitRefName:  &witRefName,
		RuleId:      converter.UUID(ruleId),
		ProcessRule: ruleUpdate,
	}

	_, err := clients.WorkItemTrackingProcessClient.UpdateProcessWorkItemTypeRule(ctx, args)
	if err != nil {
		return diag.Errorf(" Updating rule %s. Error: %+v", ruleId, err)
	}

	return resourceRuleRead(ctx, d, m)
}

func resourceRuleDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_id").(string)
	ruleId := d.Id()

	args := workitemtrackingprocess.DeleteProcessWorkItemTypeRuleArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: &witRefName,
		RuleId:     converter.UUID(ruleId),
	}

	err := clients.WorkItemTrackingProcessClient.DeleteProcessWorkItemTypeRule(ctx, args)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return nil
		}
		return diag.Errorf(" Deleting rule %s. Error: %+v", ruleId, err)
	}

	return nil
}

func importRule(ctx context.Context, d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	parts, err := tfhelper.ParseImportedNameParts(d.Id(), "process_id/work_item_type_id/rule_id", 3)
	if err != nil {
		return nil, err
	}
	d.Set("process_id", parts[0])
	d.Set("work_item_type_id", parts[1])
	d.SetId(parts[2])
	return []*schema.ResourceData{d}, nil
}

func expandConditions(conditions []any) *[]workitemtrackingprocess.RuleCondition {
	expandedConditions := make([]workitemtrackingprocess.RuleCondition, len(conditions))
	for i, condition := range conditions {
		condition := condition.(map[string]any)
		expandedCondition := workitemtrackingprocess.RuleCondition{}

		if conditionType, ok := condition["condition_type"].(string); ok && conditionType != "" {
			ct := workitemtrackingprocess.RuleConditionType(conditionType)
			expandedCondition.ConditionType = &ct
		}
		if field, ok := condition["field"].(string); ok && field != "" {
			expandedCondition.Field = &field
		}
		if value, ok := condition["value"].(string); ok && value != "" {
			expandedCondition.Value = &value
		}

		expandedConditions[i] = expandedCondition
	}
	return &expandedConditions
}

func flattenConditions(conditions []workitemtrackingprocess.RuleCondition) []map[string]any {
	flattenedConditions := make([]map[string]any, len(conditions))
	for i, condition := range conditions {
		flattenedCondition := make(map[string]any)
		if condition.ConditionType != nil {
			flattenedCondition["condition_type"] = string(*condition.ConditionType)
		}
		if condition.Field != nil {
			flattenedCondition["field"] = *condition.Field
		}
		if condition.Value != nil {
			flattenedCondition["value"] = *condition.Value
		}
		flattenedConditions[i] = flattenedCondition
	}
	return flattenedConditions
}

func expandActions(actions []any) *[]workitemtrackingprocess.RuleAction {
	expandedActions := make([]workitemtrackingprocess.RuleAction, len(actions))
	for i, action := range actions {
		action := action.(map[string]any)
		expandedAction := workitemtrackingprocess.RuleAction{}

		if actionType, ok := action["action_type"].(string); ok && actionType != "" {
			at := workitemtrackingprocess.RuleActionType(actionType)
			expandedAction.ActionType = &at
		}
		if targetField, ok := action["target_field"].(string); ok && targetField != "" {
			expandedAction.TargetField = &targetField
		}
		if value, ok := action["value"].(string); ok && value != "" {
			expandedAction.Value = &value
		}

		expandedActions[i] = expandedAction
	}
	return &expandedActions
}

func flattenActions(actions []workitemtrackingprocess.RuleAction) []map[string]any {
	flattenedActions := make([]map[string]any, len(actions))
	for i, action := range actions {
		flattenedAction := make(map[string]any)
		if action.ActionType != nil {
			flattenedAction["action_type"] = string(*action.ActionType)
		}
		if action.TargetField != nil {
			flattenedAction["target_field"] = *action.TargetField
		}
		if action.Value != nil {
			flattenedAction["value"] = *action.Value
		}
		flattenedActions[i] = flattenedAction
	}
	return flattenedActions
}

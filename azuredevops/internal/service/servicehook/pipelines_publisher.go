package servicehook

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/servicehooks"
)

var (
	apiType2ResourceBlock = map[string]string{
		"ms.vss-pipelines.run-state-changed-event":   "run_state_changed_event",
		"ms.vss-pipelines.stage-state-changed-event": "stage_state_changed_event",
	}

	resourceBlock2ApiType = map[string]string{
		"run_state_changed_event":   "ms.vss-pipelines.run-state-changed-event",
		"stage_state_changed_event": "ms.vss-pipelines.stage-state-changed-event",
	}
)

func genPipelinesPublisherSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"stage_state_changed_event": {
			Type:          schema.TypeList,
			Optional:      true,
			MaxItems:      1,
			AtLeastOneOf:  []string{"stage_state_changed_event", "run_state_changed_event"},
			ConflictsWith: []string{"run_state_changed_event"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"pipeline_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The pipeline ID to be monitored. If not specified, all pipelines in the project will trigger the event",
					},
					"stage_name": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Which stage should generate an event. Only valid if published_event is `StageStateChanged`. If not specified, all stages will trigger the event",
					},
					"stage_state_filter": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"NotStarted", "Waiting", "Running", "Completed"}, false),
						Description:  "Which stage state should generate an event. Only valid if published_event is `StageStateChanged`. If not specified, all states will trigger the event",
					},
					"stage_result_filter": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"Canceled", "Failed", "Rejected", "Skipped", "Succeeded"}, false),
						Description:  "Which stage result should generate an event. Only valid if published_event is `StageStateChanged`. If not specified, all results will trigger the event",
					},
				},
			},
		},
		"run_state_changed_event": {
			Type:          schema.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"stage_state_changed_event"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"pipeline_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The pipeline ID to be monitored. If not specified, all pipelines in the project will trigger the event",
					},
					"run_state_filter": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"InProgress", "Canceling", "Completed"}, false),
						Description:  "Which run state should generate an event. Only valid if published_event is `RunStateChanged`. If not specified, all states will trigger the event",
					},
					"run_result_filter": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"Canceled", "Failed", "Succeeded"}, false),
						Description:  "Which run result should generate an event. Only valid if published_event is `RunStateChanged`. If not specified, all results will trigger the event",
					},
				},
			},
		},
	}
}

func expandPipelinesEventConfig(d *schema.ResourceData) (map[string]string, string) {
	eventConfig := make(map[string]string)
	var eventType string
	if inputsList, ok := d.GetOkExists("stage_state_changed_event"); ok && len(inputsList.([]interface{})) > 0 {
		eventType = "stage_state_changed_event"
		if inputs, ok := inputsList.([]interface{}); ok && inputs[0] != nil {
			eventConfig["pipelineId"] = inputs[0].(map[string]interface{})["pipeline_id"].(string)
			eventConfig["stageNameId"] = inputs[0].(map[string]interface{})["stage_name"].(string)
			eventConfig["stageStateId"] = inputs[0].(map[string]interface{})["stage_state_filter"].(string)
			eventConfig["stageResultId"] = inputs[0].(map[string]interface{})["stage_result_filter"].(string)
		}
	}
	if inputsList, ok := d.GetOkExists("run_state_changed_event"); ok && len(inputsList.([]interface{})) > 0 {
		eventType = "run_state_changed_event"
		if inputs, ok := inputsList.([]interface{}); ok && inputs[0] != nil {
			eventConfig["pipelineId"] = inputs[0].(map[string]interface{})["pipeline_id"].(string)
			eventConfig["runStateId"] = inputs[0].(map[string]interface{})["run_state_filter"].(string)
			eventConfig["runResultId"] = inputs[0].(map[string]interface{})["run_result_filter"].(string)
		}
	}
	eventConfig["projectId"] = d.Get("project_id").(string)
	return eventConfig, resourceBlock2ApiType[eventType]
}

func flattenPipelinesEventConfig(subscription *servicehooks.Subscription) (string, []interface{}) {
	eventType := apiType2ResourceBlock[*subscription.EventType]
	if isNilEventConfig(*subscription.PublisherInputs) {
		return eventType, []interface{}{nil}
	}
	event := *subscription.PublisherInputs
	eventConfig := make(map[string]interface{})
	eventConfig["pipeline_id"] = event["pipelineId"]
	switch eventType {
	case "stage_state_changed_event":
		eventConfig["stage_name"] = event["stageNameId"]
		eventConfig["stage_state_filter"] = event["stageStateId"]
		eventConfig["stage_result_filter"] = event["stageResultId"]
	case "run_state_changed_event":
		eventConfig["run_state_filter"] = event["runStateId"]
		eventConfig["run_result_filter"] = event["runResultId"]
	}

	return eventType, []interface{}{eventConfig}
}

func isNilEventConfig(eventConfig map[string]string) bool {
	for key := range eventConfig {
		if key != "projectId" && key != "tfsSubscriptionId" && eventConfig[key] != "" {
			return false
		}
	}
	return true
}

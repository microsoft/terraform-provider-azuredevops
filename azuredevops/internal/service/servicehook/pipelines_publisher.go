package servicehook

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type pipelineEvent string

const (
	stageStateChanged pipelineEvent = "StageStateChanged"
	runStateChanged   pipelineEvent = "RunStateChanged"
)

type pipelineEventType string

const (
	stageStateChangedEventType pipelineEventType = "ms.vss-pipelines.stage-state-changed-event"
	runStateChangedEventType   pipelineEventType = "ms.vss-pipelines.run-state-changed-event"
)

var pipelineEvent2apiType = map[pipelineEvent]pipelineEventType{
	stageStateChanged: stageStateChangedEventType,
	runStateChanged:   runStateChangedEventType,
}

var apiType2pipelineEvent = map[pipelineEventType]pipelineEvent{
	stageStateChangedEventType: stageStateChanged,
	runStateChangedEventType:   runStateChanged,
}

func genPipelinesPublisherSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"published_event": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"StageStateChanged", "RunStateChanged"}, false),
			Description:  "The trigger event",
		},
		"event_config": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"pipeline_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "",
						Description: "The pipeline ID to be monitored. If not specified, all pipelines in the project will trigger the event",
					},
					"stage_name": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "",
						Description: "Which stage should generate an event. Only valid if published_event is `StageStateChanged`. If not specified, all stages will trigger the event",
					},
					"stage_state_filter": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "",
						ValidateFunc: validation.StringInSlice([]string{"NotStarted", "Waiting", "Running", "Completed"}, false),
						Description:  "Which stage state should generate an event. Only valid if published_event is `StageStateChanged`. If not specified, all states will trigger the event",
					},
					"stage_result_filter": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "",
						ValidateFunc: validation.StringInSlice([]string{"Canceled", "Failed", "Rejected", "Skipped", "Succeeded"}, false),
						Description:  "Which stage result should generate an event. Only valid if published_event is `StageStateChanged`. If not specified, all results will trigger the event",
					},
					"run_state_filter": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "",
						ValidateFunc: validation.StringInSlice([]string{"InProgress", "Canceling", "Completed"}, false),
						Description:  "Which run state should generate an event. Only valid if published_event is `RunStateChanged`. If not specified, all states will trigger the event",
					},
					"run_result_filter": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "",
						ValidateFunc: validation.StringInSlice([]string{"Canceled", "Failed", "Succeeded"}, false),
						Description:  "Which run result should generate an event. Only valid if published_event is `RunStateChanged`. If not specified, all results will trigger the event",
					},
				},
			},
		},
	}
}

func expandPipelinesEventConfig(d *schema.ResourceData) (*map[string]string, *string) {
	publishedEvent := pipelineEvent(d.Get("published_event").(string))

	eventConfig := make(map[string]string)
	if v, ok := d.GetOk("event_config"); ok {
		inputs := v.([]interface{})[0].(map[string]interface{})
		eventConfig["pipelineId"] = inputs["pipeline_id"].(string)
		if publishedEvent == stageStateChanged {
			eventConfig["stageStateId"] = inputs["stage_state_filter"].(string)
			eventConfig["stageResultId"] = inputs["stage_result_filter"].(string)
			eventConfig["stageNameId"] = inputs["stage_name"].(string)
		}
		if publishedEvent == runStateChanged {
			eventConfig["runStateId"] = inputs["run_state_filter"].(string)
			eventConfig["runResultId"] = inputs["run_result_filter"].(string)
		}
	} else {
		eventConfig["pipelineId"] = ""
		if publishedEvent == stageStateChanged {
			eventConfig["stageStateId"] = ""
			eventConfig["stageResultId"] = ""
			eventConfig["stageNameId"] = ""
		}
		if publishedEvent == runStateChanged {
			eventConfig["runStateId"] = ""
			eventConfig["runResultId"] = ""
		}
	}
	eventConfig["projectId"] = d.Get("project_id").(string)

	eventType := string(pipelineEvent2apiType[publishedEvent])
	return &eventConfig, &eventType
}

func flattenPipelinesEventConfig(publishedEvent pipelineEvent, event *map[string]string) []interface{} {
	eventConfig := make(map[string]interface{})

	eventConfig["pipeline_id"] = (*event)["pipelineId"]

	if publishedEvent == stageStateChanged {
		eventConfig["stage_state_filter"] = (*event)["stageStateId"]
		eventConfig["stage_result_filter"] = (*event)["stageResultId"]
		eventConfig["stage_name"] = (*event)["stageNameId"]
	}

	if publishedEvent == runStateChanged {
		eventConfig["run_state_filter"] = (*event)["runStateId"]
		eventConfig["run_result_filter"] = (*event)["runResultId"]
	}

	return []interface{}{eventConfig}
}

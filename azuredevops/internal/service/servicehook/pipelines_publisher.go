package servicehook

import (
	"context"
	"fmt"
	"strings"

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
		"stage_state_changed_event": {
			Type:          schema.TypeList,
			Optional:      true,
			MaxItems:      1,
			Default:       nil,
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
			Default:       nil,
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

func validateEventConfigDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	publishedEvent := pipelineEvent(d.Get("published_event").(string))
	expectedResourceBlock := convertFromApiType2ResourceBlock(string(pipelineEvent2apiType[publishedEvent]))
	// iterate through all arguments to find a block
	changedKeys := d.GetChangedKeysPrefix("")
	for _, key := range changedKeys {
		if strings.Contains(key, "changed_event") && !strings.HasPrefix(key, expectedResourceBlock) {
			if _, ok := d.GetOk(key); ok {
				return fmt.Errorf("Only '%s' block is supported if published_event is '%s'", expectedResourceBlock, publishedEvent)
			}
		}
	}

	return nil
}

func expandPipelinesEventConfig(d *schema.ResourceData) (*map[string]string, *string) {
	publishedEvent := pipelineEvent(d.Get("published_event").(string))
	eventType := string(pipelineEvent2apiType[publishedEvent])
	convertedEventType := convertFromApiType2ResourceBlock(eventType)

	eventConfig := make(map[string]string)
	if _, ok := d.GetOk(convertedEventType); ok {
		inputs, ok := d.Get(convertedEventType).([]interface{})[0].(map[string]interface{})
		if ok {
			eventConfig["pipelineId"] = inputs["pipeline_id"].(string)
			switch publishedEvent {
			case stageStateChanged:
				eventConfig["stageNameId"] = inputs["stage_name"].(string)
				eventConfig["stageStateId"] = inputs["stage_state_filter"].(string)
				eventConfig["stageResultId"] = inputs["stage_result_filter"].(string)
			case runStateChanged:
				eventConfig["runStateId"] = inputs["run_state_filter"].(string)
				eventConfig["runResultId"] = inputs["run_result_filter"].(string)
			}
		}
	}
	eventConfig["projectId"] = d.Get("project_id").(string)

	return &eventConfig, &eventType
}

func flattenPipelinesEventConfig(publishedEvent pipelineEvent, event *map[string]string) []interface{} {
	if isNilEventConfig(*event) {
		return nil
	}
	eventConfig := make(map[string]interface{})
	eventConfig["pipeline_id"] = (*event)["pipelineId"]
	switch publishedEvent {
	case stageStateChanged:
		eventConfig["stage_name"] = (*event)["stageNameId"]
		eventConfig["stage_state_filter"] = (*event)["stageStateId"]
		eventConfig["stage_result_filter"] = (*event)["stageResultId"]
	case runStateChanged:
		eventConfig["run_state_filter"] = (*event)["runStateId"]
		eventConfig["run_result_filter"] = (*event)["runResultId"]
	}

	return []interface{}{eventConfig}
}

func convertFromApiType2ResourceBlock(eventType string) string {
	eventTypeSplited := strings.Split(eventType, ".")
	eventTypeSplited[len(eventTypeSplited)-1] = strings.Replace(eventTypeSplited[len(eventTypeSplited)-1], "-", "_", -1)
	return eventTypeSplited[len(eventTypeSplited)-1]
}

func isNilEventConfig(eventConfig map[string]string) bool {
	for key := range eventConfig {
		if key != "projectId" && key != "tfsSubscriptionId" && eventConfig[key] != "" {
			return false
		}
	}
	return true
}

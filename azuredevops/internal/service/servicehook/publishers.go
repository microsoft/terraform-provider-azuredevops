package servicehook

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var eventTypeMap = map[string]map[string]string{
	"pipelines": {
		"stage_state_changed": "ms.vss-pipelines.stage-state-changed-event",
		"run_state_changed":   "ms.vss-pipelines.run-state-changed-event",
	},
	"distributedtasks": {},
	"rm":               {},
}

var publisherResourceVersionMap = map[string]string{
	"pipelines":        "5.1-preview.1",
	"distributedtasks": "1.0-preview.1",
	"rm":               "3.0-preview.1",
}

func genPublisherSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"publisher": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"pipelines", "distributedtasks", "rm"}, false),
					},
					"stage_state_changed": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Default:  nil,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"pipeline_id": {
									Type:     schema.TypeString,
									Optional: true,
									Default:  "",
								},
								"stage_name": {
									Type:     schema.TypeString,
									Optional: true,
									Default:  "",
								},
								"state_filter": {
									Type:         schema.TypeString,
									Optional:     true,
									Default:      "",
									ValidateFunc: validation.StringInSlice([]string{"NotStarted", "Waiting", "Running", "Completed"}, false),
								},
								"result_filter": {
									Type:         schema.TypeString,
									Optional:     true,
									Default:      "",
									ValidateFunc: validation.StringInSlice([]string{"Canceled", "Failed", "Rejected", "Skipped", "Succeeded"}, false),
								},
							},
						},
					},
					"run_state_changed": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Default:  nil,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"pipeline_id": {
									Type:     schema.TypeString,
									Optional: true,
									Default:  "",
								},
								"state_filter": {
									Type:         schema.TypeString,
									Optional:     true,
									Default:      "",
									ValidateFunc: validation.StringInSlice([]string{"InProgress", "Canceling", "Completed"}, false),
								},
								"result_filter": {
									Type:         schema.TypeString,
									Optional:     true,
									Default:      "",
									ValidateFunc: validation.StringInSlice([]string{"Canceled", "Failed", "Succeeded"}, false),
								},
							},
						},
					},
				},
			},
		},
	}
}

func getEventType(publisher map[string]interface{}) *string {
	publisherName := publisher["name"].(string)
	for k, v := range publisher {
		if eventInputs, ok := v.([]interface{}); ok && len(eventInputs) > 0 {
			eventType := eventTypeMap[publisherName][k]
			return &eventType
		}
	}
	return nil
}

func expandPublisherInputs(projectId string, publisher []interface{}) *map[string]string {
	publisherMap := publisher[0].(map[string]interface{})

	publisherInputs := make(map[string]string)

	publisherName := publisherMap["name"].(string)
	switch publisherName {
	case "pipelines":
		// For simplicity, I'm only handling the 'stage_state_changed' event here.
		if stageStateChangedList, ok := publisherMap["stage_state_changed"].([]interface{}); ok {
			stageStateChanged := stageStateChangedList[0].(map[string]interface{})
			publisherInputs["pipelineId"] = stageStateChanged["pipeline_id"].(string)
			publisherInputs["stageNameId"] = stageStateChanged["stage_name"].(string)
			publisherInputs["stageStateId"] = stageStateChanged["state_filter"].(string)
			publisherInputs["stageResultId"] = stageStateChanged["result_filter"].(string)
			publisherInputs["projectId"] = projectId
		}
		// case "distributedtasks": ...
		// case "rm": ...
	}

	return &publisherInputs
}

func flattenPublisherInputs(publisherId string, publisherInputs map[string]string) []interface{} {
	// Create an empty map for the publisher to be returned
	publisher := make(map[string]interface{})

	// Set the publisher name
	publisher["name"] = publisherId

	// Based on the publisher name, determine which set of configurations to fill
	switch publisherId {
	case "pipelines":
		// Again, for simplicity, only handling the 'stage_state_changed' event.
		stageStateChanged := make(map[string]interface{})
		stageStateChanged["pipeline_id"] = publisherInputs["pipelineId"]
		stageStateChanged["stage_name"] = publisherInputs["stageNameId"]
		stageStateChanged["state_filter"] = publisherInputs["stageStateId"]
		stageStateChanged["result_filter"] = publisherInputs["stageResultId"]
		publisher["stage_state_changed"] = []interface{}{stageStateChanged}
		// case "distributedtasks": ...
		// case "rm": ...
	}

	return []interface{}{publisher}
}

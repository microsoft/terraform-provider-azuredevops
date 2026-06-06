package utils

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/audit"
)

// ResourceAuditStreamSchema returns the schema for the audit stream resource
func ResourceAuditStreamSchema(outer map[string]*schema.Schema) map[string]*schema.Schema {
	if outer == nil {
		outer = make(map[string]*schema.Schema)
	}
	baseSchema := map[string]*schema.Schema{
		"display_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The human-readable name for the audit stream.",
		},
		"consumer_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The type of the consumer (e.g., 'splunk', 'azureMonitorLogs').",
		},
		"status": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "enabled",
			Description: "The status of the stream. Valid values are 'enabled' and 'disabledByUser'.",
		},
		"consumer_inputs": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "A list of key-value pairs of consumer inputs.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:      schema.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},
		"created_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The time when the stream was created.",
		},
		"updated_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The time when the stream was last updated.",
		},
		"status_reason": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The reason for the current status.",
		},
	}
	for key, elem := range baseSchema {
		outer[key] = elem
	}

	return outer
}

// DataAuditStreamSchema returns the schema for the audit stream data source
func DataAuditStreamSchema(outer map[string]*schema.Schema) map[string]*schema.Schema {
	if outer == nil {
		outer = make(map[string]*schema.Schema)
	}
	baseSchema := map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The unique ID of the audit stream.",
		},
		"display_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The human-readable name for the audit stream.",
		},
		"consumer_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The type of the consumer.",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the stream.",
		},
		"consumer_inputs": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "A list of key-value pairs of consumer inputs.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"value": {
						Type:      schema.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},
		"created_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The time when the stream was created.",
		},
		"updated_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The time when the stream was last updated.",
		},
		"status_reason": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The reason for the current status.",
		},
	}
	for key, elem := range baseSchema {
		outer[key] = elem
	}

	return outer
}

// DataAuditStreamsSchema returns the schema for the audit streams data source
func DataAuditStreamsSchema(outer map[string]*schema.Schema) map[string]*schema.Schema {
	if outer == nil {
		outer = make(map[string]*schema.Schema)
	}
	baseSchema := map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The unique ID of the data source.",
		},
		"streams": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "A list of all configured audit streams in the organization.",
			Elem: &schema.Resource{
				Schema: DataAuditStreamSchemaComputed(make(map[string]*schema.Schema)),
			},
		},
	}
	for key, elem := range baseSchema {
		outer[key] = elem
	}

	return outer
}

// DataAuditStreamSchemaComputed returns the schema for the audit stream data source with all fields as computed
func DataAuditStreamSchemaComputed(outer map[string]*schema.Schema) map[string]*schema.Schema {
	if outer == nil {
		outer = make(map[string]*schema.Schema)
	}
	baseSchema := map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The unique ID of the audit stream.",
		},
		"display_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The human-readable name for the audit stream.",
		},
		"consumer_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The type of the consumer.",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the stream.",
		},
		"consumer_inputs": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "A list of key-value pairs of consumer inputs.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"value": {
						Type:      schema.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},
		"created_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The time when the stream was created.",
		},
		"updated_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The time when the stream was last updated.",
		},
		"status_reason": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The reason for the current status.",
		},
	}
	for key, elem := range baseSchema {
		outer[key] = elem
	}

	return outer
}

// ExpandAuditStreamStatus expands a string to an AuditStreamStatus
func ExpandAuditStreamStatus(statusStr string) *audit.AuditStreamStatus {
	var status audit.AuditStreamStatus

	switch statusStr {
	case "enabled":
		status = audit.AuditStreamStatusValues.Enabled
	case "disabledByUser":
		status = audit.AuditStreamStatusValues.DisabledByUser
	case "disabledBySystem":
		status = audit.AuditStreamStatusValues.DisabledBySystem
	case "deleted":
		status = audit.AuditStreamStatusValues.Deleted
	case "backfilling":
		status = audit.AuditStreamStatusValues.Backfilling
	default:
		status = audit.AuditStreamStatusValues.Enabled
	}

	return &status
}

// ExpandAuditStream expands a ResourceData to an AuditStream
func ExpandAuditStream(d *schema.ResourceData) audit.AuditStream {
	displayName := d.Get("display_name").(string)
	consumerType := d.Get("consumer_type").(string)
	status := ExpandAuditStreamStatus(d.Get("status").(string))

	return audit.AuditStream{
		DisplayName:    &displayName,
		ConsumerType:   &consumerType,
		Status:         status,
		ConsumerInputs: ExpandConsumerInputs(d),
	}
}

// ExpandConsumerInputs expands a ResourceData to a map of consumer inputs
func ExpandConsumerInputs(d *schema.ResourceData) *map[string]string {
	v, ok := d.GetOk("consumer_inputs")
	if !ok || v == nil {
		return nil
	}

	tfInputs := v.(*schema.Set).List()
	apiInputs := make(map[string]string, len(tfInputs))

	for _, input := range tfInputs {
		inputMap := input.(map[string]interface{})
		key := inputMap["key"].(string)
		value := inputMap["value"].(string)
		apiInputs[key] = value
	}

	return &apiInputs
}

// FlattenConsumerInputs flattens a map of consumer inputs to a schema.Set
func FlattenConsumerInputs(inputs *map[string]string) *schema.Set {
	if inputs == nil {
		return nil
	}

	inputSet := schema.NewSet(schema.HashResource(ResourceAuditStreamSchema(nil)["consumer_inputs"].Elem.(*schema.Resource)), []interface{}{})

	for key, value := range *inputs {
		inputMap := map[string]interface{}{
			"key":   key,
			"value": value,
		}
		inputSet.Add(inputMap)
	}
	return inputSet
}

// FlattenAuditStream flattens an AuditStream to a ResourceData
func FlattenAuditStream(d *schema.ResourceData, stream *audit.AuditStream) error {
	if stream == nil {
		return nil
	}

	d.Set("display_name", *stream.DisplayName)
	d.Set("consumer_type", *stream.ConsumerType)
	d.Set("status", string(*stream.Status))

	if stream.CreatedTime != nil {
		d.Set("created_time", stream.CreatedTime.String())
	}

	if stream.UpdatedTime != nil {
		d.Set("updated_time", stream.UpdatedTime.String())
	}

	if stream.StatusReason != nil {
		d.Set("status_reason", *stream.StatusReason)
	}

	d.Set("consumer_inputs", FlattenConsumerInputs(stream.ConsumerInputs))

	return nil
}

// FlattenSingleAuditStream flattens an AuditStream to a map
func FlattenSingleAuditStream(m map[string]interface{}, stream *audit.AuditStream) error {
	if stream == nil {
		return nil
	}

	m["id"] = strconv.Itoa(*stream.Id)
	m["display_name"] = *stream.DisplayName
	m["consumer_type"] = *stream.ConsumerType
	m["status"] = string(*stream.Status)

	if stream.CreatedTime != nil {
		m["created_time"] = stream.CreatedTime.String()
	}
	if stream.UpdatedTime != nil {
		m["updated_time"] = stream.UpdatedTime.String()
	}

	if stream.StatusReason != nil {
		m["status_reason"] = *stream.StatusReason
	}

	if inputs := FlattenConsumerInputs(stream.ConsumerInputs); inputs != nil {
		m["consumer_inputs"] = inputs.List()
	} else {
		m["consumer_inputs"] = []any{}
	}

	return nil
}

package utils

import (
	"strconv" // Import toegevoegd voor FlattenSingleAuditStream

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/audit"
)

func ResourceAuditStreamSchema(outer map[string]*schema.Schema) map[string]*schema.Schema {
	baseSchema := map[string]*schema.Schema{
		"display_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "De leesbare naam voor de auditstroom. (Vertaalt naar REST API veld 'displayName').",
		},
		"consumer_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Het type externe service (bijv. 'splunk', 'azureMonitorLogs'). (Vertaalt naar REST API veld 'consumerType').",
		},
		"status": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "enabled",
			Description: "De gewenste status van de stroom ('enabled' of 'disabledByUser').",
		},
		"consumer_inputs": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "Een lijst van key/value paren met de benodigde invoerparameters voor de consument.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "De sleutel van de invoerparameter (bijv. 'url', 'token').",
					},
					"value": {
						Type:        schema.TypeString,
						Required:    true,
						Sensitive:   true,
						Description: "De bijbehorende waarde of secret.",
					},
				},
			},
		},
		"created_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Het tijdstip waarop de stroom is aangemaakt.",
		},
		"updated_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Het tijdstip waarop de stroom voor het laatst is bijgewerkt.",
		},
		"status_reason": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "De reden voor de huidige status, indien van toepassing.",
		},
	}
	for key, elem := range baseSchema {
		outer[key] = elem
	}

	return outer
}

func DataAuditStreamSchema(outer map[string]*schema.Schema) map[string]*schema.Schema {
	baseSchema := map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "De unieke ID van de gevonden auditstroom.",
		},
		"display_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "De leesbare naam om de auditstroom op te zoeken.",
		},
		"consumer_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Het type externe service van de gevonden stream.",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "De status van de gevonden stream.",
		},
		"consumer_inputs": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "De key/value paren met de invoerparameters van de consument.",
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
			Description: "Het tijdstip waarop de stroom is aangemaakt.",
		},
		"updated_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Het tijdstip waarop de stroom voor het laatst is bijgewerkt.",
		},
		"status_reason": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "De reden voor de huidige status.",
		},
	}
	for key, elem := range baseSchema {
		outer[key] = elem
	}

	return outer
}

func DataAuditStreamsSchema(outer map[string]*schema.Schema) map[string]*schema.Schema {
	baseSchema := map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "De unieke ID van de data source (wordt ingesteld op een constante waarde).",
		},

		"streams": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Een lijst van alle geconfigureerde audit streams in de organisatie.",
			Elem: &schema.Resource{
				Schema: DataAuditStreamSchema(outer),
			},
		},
	}
	for key, elem := range baseSchema {
		outer[key] = elem
	}

	return outer
}

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

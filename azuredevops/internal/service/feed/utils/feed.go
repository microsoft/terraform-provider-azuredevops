package utils

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/feed"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func FlattenUpstreamSources(sources []feed.UpstreamSource) []any {
	var result []any
	for _, source := range sources {
		s := make(map[string]any)

		if source.Id != nil {
			s["id"] = source.Id.String()
		}
		if source.Name != nil {
			s["name"] = *source.Name
		}
		if source.Protocol != nil {
			s["protocol"] = *source.Protocol
		}
		if source.Location != nil {
			s["location"] = *source.Location
		}
		if source.DisplayLocation != nil {
			s["display_location"] = *source.DisplayLocation
		}
		if source.UpstreamSourceType != nil {
			s["upstream_source_type"] = string(*source.UpstreamSourceType)
		}
		if source.Status != nil {
			s["status"] = string(*source.Status)
		}
		if source.ServiceEndpointId != nil {
			s["service_endpoint_id"] = source.ServiceEndpointId.String()
		}
		if source.ServiceEndpointProjectId != nil {
			s["service_endpoint_project_id"] = source.ServiceEndpointProjectId.String()
		}
		if source.InternalUpstreamCollectionId != nil {
			s["internal_upstream_collection_id"] = source.InternalUpstreamCollectionId.String()
		}
		if source.InternalUpstreamFeedId != nil {
			s["internal_upstream_feed_id"] = source.InternalUpstreamFeedId.String()
		}
		if source.InternalUpstreamViewId != nil {
			s["internal_upstream_view_id"] = source.InternalUpstreamViewId.String()
		}

		result = append(result, s)
	}
	return result
}

func CommonFeedFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"url": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"badges_enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"hide_deleted_package_versions": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"upstream_enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"upstream_sources": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: upstreamSourceSchema(),
			},
		},
	}
}

func upstreamSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"protocol": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"location": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"display_location": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"upstream_source_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func ResourceUpstreamSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"protocol": {
			Type:     schema.TypeString,
			Required: true,
		},
		"location": {
			Type:     schema.TypeString,
			Required: true,
		},
		"upstream_source_type": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "public",
			ValidateFunc: validation.StringInSlice([]string{
				"public", "internal",
			}, false),
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"display_location": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"service_endpoint_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
			Description:  "The ID of the Service Endpoint used for authentication.",
		},
		"service_endpoint_project_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
			Description:  "The ID of the Project where the Service Endpoint exists.",
		},
		"internal_upstream_collection_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
			Description:  "The ID of the collection where the internal feed resides.",
		},
		"internal_upstream_feed_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
			Description:  "The ID of the internal feed.",
		},
		"internal_upstream_view_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
			Description:  "The ID of the view in the internal upstream feed.",
		},
	}
}

func ExpandUpstreamSources(input []any) *[]feed.UpstreamSource {
	if len(input) == 0 {
		return nil
	}

	results := make([]feed.UpstreamSource, 0, len(input))

	for _, item := range input {
		data := item.(map[string]any)

		source := feed.UpstreamSource{
			Name:     converter.String(data["name"].(string)),
			Protocol: converter.String(data["protocol"].(string)),
			Location: converter.String(data["location"].(string)),
		}

		if v, ok := data["upstream_source_type"].(string); ok && v != "" {
			t := feed.UpstreamSourceType(v)
			source.UpstreamSourceType = &t
		}

		if v, ok := data["service_endpoint_id"].(string); ok && v != "" {
			if id, err := uuid.Parse(v); err == nil {
				source.ServiceEndpointId = &id
			}
		}

		if v, ok := data["service_endpoint_project_id"].(string); ok && v != "" {
			if id, err := uuid.Parse(v); err == nil {
				source.ServiceEndpointProjectId = &id
			}
		}

		if v, ok := data["internal_upstream_collection_id"].(string); ok && v != "" {
			if id, err := uuid.Parse(v); err == nil {
				source.InternalUpstreamCollectionId = &id
			}
		}

		if v, ok := data["internal_upstream_feed_id"].(string); ok && v != "" {
			if id, err := uuid.Parse(v); err == nil {
				source.InternalUpstreamFeedId = &id
			}
		}

		if v, ok := data["internal_upstream_view_id"].(string); ok && v != "" {
			if id, err := uuid.Parse(v); err == nil {
				source.InternalUpstreamViewId = &id
			}
		}

		results = append(results, source)
	}

	return &results
}

func ExpandFeedFeatures(input []interface{}) map[string]interface{} {
	if len(input) == 0 || input[0] == nil {
		return map[string]interface{}{}
	}
	return input[0].(map[string]interface{})
}

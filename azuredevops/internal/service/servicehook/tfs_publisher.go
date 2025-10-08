package servicehook

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/servicehooks"
)

// TFS Event Type mappings based on Microsoft Azure DevOps documentation
var tfsResourceBlock2ApiType = map[string]string{
	// Build and Release Events
	"build_completed": "build.complete",

	// Code Events - Git
	"git_push":                         "git.push",
	"git_pull_request_created":         "git.pullrequest.created",
	"git_pull_request_merge_attempted": "git.pullrequest.merged",
	"git_pull_request_updated":         "git.pullrequest.updated",
	"git_pull_request_commented":       "ms.vss-code.git-pullrequest-comment-event",
	"repository_created":               "git.repo.created",
	"repository_deleted":               "git.repo.deleted",
	"repository_forked":                "git.repo.forked",
	"repository_renamed":               "git.repo.renamed",
	"repository_status_changed":        "git.repo.statuschanged",

	// Code Events - TFVC
	"tfvc_checkin": "tfvc.checkin",

	// Work Item Events
	"work_item_created":   "workitem.created",
	"work_item_deleted":   "workitem.deleted",
	"work_item_restored":  "workitem.restored",
	"work_item_updated":   "workitem.updated",
	"work_item_commented": "workitem.commented",

	// Service Connection Events
	"service_connection_created": "ms.vss-endpoint.endpoint-created",
	"service_connection_updated": "ms.vss-endpoint.endpoint-updated",
}

var tfsApiType2ResourceBlock = map[string]string{
	// Build and Release Events
	"build.complete": "build_completed",

	// Code Events - Git
	"git.push":                                  "git_push",
	"git.pullrequest.created":                   "git_pull_request_created",
	"git.pullrequest.merged":                    "git_pull_request_merge_attempted",
	"git.pullrequest.updated":                   "git_pull_request_updated",
	"ms.vss-code.git-pullrequest-comment-event": "git_pull_request_commented",
	"git.repo.created":                          "repository_created",
	"git.repo.deleted":                          "repository_deleted",
	"git.repo.forked":                           "repository_forked",
	"git.repo.renamed":                          "repository_renamed",
	"git.repo.statuschanged":                    "repository_status_changed",

	// Code Events - TFVC
	"tfvc.checkin": "tfvc_checkin",

	// Work Item Events
	"workitem.created":   "work_item_created",
	"workitem.deleted":   "work_item_deleted",
	"workitem.restored":  "work_item_restored",
	"workitem.updated":   "work_item_updated",
	"workitem.commented": "work_item_commented",

	// Service Connection Events
	"ms.vss-endpoint.endpoint-created": "service_connection_created",
	"ms.vss-endpoint.endpoint-updated": "service_connection_updated",
}

// TfsPublisherSchema represents the publisher schema for TFS events
func genTfsPublisherSchema() map[string]*schema.Schema {
	eventTypes := []string{
		"build_completed",
		"git_pull_request_commented",
		"git_pull_request_created",
		"git_pull_request_merge_attempted",
		"git_pull_request_updated",
		"git_push",
		"repository_created",
		"repository_deleted",
		"repository_forked",
		"repository_renamed",
		"repository_status_changed",
		"service_connection_created",
		"service_connection_updated",
		"tfvc_checkin",
		"work_item_commented",
		"work_item_created",
		"work_item_deleted",
		"work_item_restored",
		"work_item_updated",
	}

	return map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The project ID that will be used for the TFS event subscription",
		},

		// Build and Release Events
		"build_completed": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"definition_name": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for completed builds for a specific pipeline.",
					},
					"build_status": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"Succeeded", "PartiallySucceeded", "Failed", "Stopped"}, false),
						Description:  "Include only events for completed builds that have a specific completion status.",
					},
				},
			},
		},

		// Git Events
		"git_pull_request_commented": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"repository_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for pull requests in a specific repository (repository ID). If not specified, all repositories in the project will trigger the event.",
					},
					"branch": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for pull requests in a specific branch.",
					},
				},
			},
		},

		"git_pull_request_created": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"repository_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for pull requests in a specific repository (repository ID). If not specified, all repositories in the project will trigger the event.",
					},
					"pull_request_created_by": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for pull requests created by users in a specific group.",
					},
					"pull_request_reviewers_contains": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for pull requests with reviewers in a specific group.",
					},
					"branch": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for pull requests in a specific branch.",
					},
				},
			},
		},

		"git_pull_request_merge_attempted": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"repository_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for pull requests in a specific repository (repository ID). If not specified, all repositories in the project will trigger the event.",
					},
					"pull_request_created_by": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for pull requests created by users in a specific group.",
					},
					"pull_request_reviewers_contains": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for pull requests with reviewers in a specific group.",
					},
					"branch": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for pull requests in a specific branch.",
					},
					"merge_result": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"Succeeded", "Unsuccessful", "Conflicts", "Failure", "RejectedByPolicy"}, false),
						Description:  "Include only events for pull requests with a specific merge result.",
					},
				},
			},
		},

		"git_pull_request_updated": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"notification_type": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"PushNotification", "ReviewersUpdateNotification", "StatusUpdateNotification", "ReviewerVoteNotification"}, false),
						Description:  "Include only events for pull requests with a specific change.",
					},
					"repository_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for pull requests in a specific repository (repository ID). If not specified, all repositories in the project will trigger the event.",
					},
					"pull_request_created_by": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for pull requests created by users in a specific group.",
					},
					"pull_request_reviewers_contains": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for pull requests with reviewers in a specific group.",
					},
					"branch": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for pull requests in a specific branch.",
					},
				},
			},
		},

		"git_push": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"branch": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for code pushes to a specific branch.",
					},
					"pushed_by": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for code pushes by users in a specific group.",
					},
					"repository_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for code pushes to a specific repository (repository ID). If not specified, all repositories in the project will trigger the event.",
					},
				},
			},
		},

		// Repository Events
		"repository_created": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"project_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for repositories created in a specific project.",
					},
				},
			},
		},

		"repository_deleted": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"repository_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for repositories with a specific repository ID.",
					},
				},
			},
		},

		"repository_forked": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"repository_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for repositories with a specific repository ID.",
					},
				},
			},
		},

		"repository_renamed": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"repository_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for repositories with a specific repository ID.",
					},
				},
			},
		},

		"repository_status_changed": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"repository_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for repositories with a specific repository ID.",
					},
				},
			},
		},

		// Service Connection Events
		"service_connection_created": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"project_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for service connections created in a specific project.",
					},
				},
			},
		},

		"service_connection_updated": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"project_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for service connections updated in a specific project.",
					},
				},
			},
		},

		// TFVC Events
		"tfvc_checkin": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"path": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Include only events for check-ins that change files under a specific path.",
					},
				},
			},
		},

		// Work Item Events
		"work_item_commented": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"area_path": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items under a specific area path.",
					},
					"comment_pattern": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items with a comment that contains a specific string.",
					},
					"work_item_type": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items of a specific type.",
					},
					"tag": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items that contain a specific tag.",
					},
				},
			},
		},

		"work_item_created": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"area_path": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items under a specific area path.",
					},
					"work_item_type": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items of a specific type.",
					},
					"links_changed": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Include only events for work items with one or more links added or removed.",
					},
					"tag": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items that contain a specific tag.",
					},
				},
			},
		},

		"work_item_deleted": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"area_path": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items under a specific area path.",
					},
					"work_item_type": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items of a specific type.",
					},
					"tag": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items that contain a specific tag.",
					},
				},
			},
		},

		"work_item_restored": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"area_path": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items under a specific area path.",
					},
					"work_item_type": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items of a specific type.",
					},
					"tag": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items that contain a specific tag.",
					},
				},
			},
		},

		"work_item_updated": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: eventTypes,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"area_path": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items under a specific area path.",
					},
					"changed_fields": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items with a change in a specific field.",
					},
					"work_item_type": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items of a specific type.",
					},
					"links_changed": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Include only events for work items with one or more links added or removed.",
					},
					"tag": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Include only events for work items that contain a specific tag.",
					},
				},
			},
		},
	}
}

// expandTfsEventConfig expands the TFS event configuration from schema
func expandTfsEventConfig(d *schema.ResourceData) (map[string]string, string) {
	eventConfig := make(map[string]string)

	eventTypes := []string{
		"build_completed",
		"git_pull_request_commented",
		"git_pull_request_created",
		"git_pull_request_merge_attempted",
		"git_pull_request_updated",
		"git_push",
		"repository_created",
		"repository_deleted",
		"repository_forked",
		"repository_renamed",
		"repository_status_changed",
		"service_connection_created",
		"service_connection_updated",
		"tfvc_checkin",
		"work_item_commented",
		"work_item_created",
		"work_item_deleted",
		"work_item_restored",
		"work_item_updated",
	}

	// eventSelector selects the specified event config.
	// Returns: event type and the raw event config
	eventSelector := func(d *schema.ResourceData) (string, map[string]interface{}) {
		for _, evt := range eventTypes {
			if l := d.Get(evt); len(l.([]interface{})) != 0 {
				if rawConfig, ok := l.([]interface{})[0].(map[string]interface{}); ok {
					return evt, rawConfig
				} else {
					// This can happen when the event config's nested fields are all optional and absent in the config
					return evt, map[string]interface{}{}
				}
			}
		}
		return "", nil
	}

	eventType, rawConfig := eventSelector(d)

	// Helper function to extract string from rawConfig and add to eventConfig if non-empty
	addField := func(schemaKey, apiKey string) {
		if val, ok := rawConfig[schemaKey].(string); ok && val != "" {
			eventConfig[apiKey] = val
		}
	}

	// Map schema fields to API fields based on event type
	if rawConfig != nil {
		switch eventType {
		case "build_completed":
			addField("build_status", "buildStatus")
			addField("definition_name", "definitionName")

		case "git_push":
			addField("repository_id", "repository")
			addField("branch", "branch")
			addField("pushed_by", "pushedBy")

		case "git_pull_request_created":
			addField("repository_id", "repository")
			addField("branch", "branch")
			addField("pull_request_created_by", "pullrequestCreatedBy")
			addField("pull_request_reviewers_contains", "pullrequestReviewersContains")

		case "git_pull_request_updated":
			addField("repository_id", "repository")
			addField("branch", "branch")
			addField("notification_type", "notificationType")
			addField("pull_request_created_by", "pullrequestCreatedBy")
			addField("pull_request_reviewers_contains", "pullrequestReviewersContains")

		case "git_pull_request_merge_attempted":
			addField("repository_id", "repository")
			addField("branch", "branch")
			addField("pull_request_created_by", "pullrequestCreatedBy")
			addField("pull_request_reviewers_contains", "pullrequestReviewersContains")
			addField("merge_result", "mergeResult")

		case "git_pull_request_commented":
			addField("repository_id", "repository")
			addField("branch", "branch")

		case "repository_created":
			addField("project_id", "projectId")

		case "repository_deleted", "repository_forked", "repository_renamed", "repository_status_changed":
			addField("repository_id", "repository")

		case "tfvc_checkin":
			addField("path", "path")

		case "work_item_created", "work_item_deleted", "work_item_restored":
			addField("work_item_type", "workItemType")
			addField("area_path", "areaPath")
			addField("tag", "tag")

		case "work_item_updated":
			addField("work_item_type", "workItemType")
			addField("area_path", "areaPath")
			addField("tag", "tag")
			addField("changed_fields", "changedFields")

		case "work_item_commented":
			addField("work_item_type", "workItemType")
			addField("area_path", "areaPath")
			addField("tag", "tag")
			addField("comment_pattern", "commentPattern")

		case "service_connection_created", "service_connection_updated":
			addField("project", "project")
		}
	}

	eventConfig["projectId"] = d.Get("project_id").(string)
	return eventConfig, tfsResourceBlock2ApiType[eventType]
}

// flattenTfsEventConfig flattens the TFS event configuration to schema
func flattenTfsEventConfig(subscription *servicehooks.Subscription) (string, []interface{}) {
	if subscription == nil {
		return "", []interface{}{nil}
	}
	eventType := tfsApiType2ResourceBlock[*subscription.EventType]
	if isNilTfsEventConfig(*subscription.PublisherInputs) {
		return eventType, []interface{}{nil}
	}

	event := *subscription.PublisherInputs
	eventConfig := make(map[string]interface{})

	// Set filters based on event type
	switch eventType {
	case "build_completed":
		if val, exists := event["buildStatus"]; exists {
			eventConfig["build_status"] = val
		}
		if val, exists := event["definitionName"]; exists {
			eventConfig["definition_name"] = val
		}

	case "git_push":
		if val, exists := event["repository"]; exists {
			eventConfig["repository_id"] = val
		}
		if val, exists := event["branch"]; exists {
			eventConfig["branch"] = val
		}
		if val, exists := event["pushedBy"]; exists {
			eventConfig["pushed_by"] = val
		}

	case "git_pull_request_created":
		if val, exists := event["repository"]; exists {
			eventConfig["repository_id"] = val
		}
		if val, exists := event["branch"]; exists {
			eventConfig["branch"] = val
		}
		if val, exists := event["pullrequestCreatedBy"]; exists {
			eventConfig["pull_request_created_by"] = val
		}
		if val, exists := event["pullrequestReviewersContains"]; exists {
			eventConfig["pull_request_reviewers_contains"] = val
		}

	case "git_pull_request_updated":
		if val, exists := event["repository"]; exists {
			eventConfig["repository_id"] = val
		}
		if val, exists := event["branch"]; exists {
			eventConfig["branch"] = val
		}
		if val, exists := event["notificationType"]; exists {
			eventConfig["notification_type"] = val
		}
		if val, exists := event["pullrequestCreatedBy"]; exists {
			eventConfig["pull_request_created_by"] = val
		}
		if val, exists := event["pullrequestReviewersContains"]; exists {
			eventConfig["pull_request_reviewers_contains"] = val
		}

	case "git_pull_request_merge_attempted":
		if val, exists := event["repository"]; exists {
			eventConfig["repository_id"] = val
		}
		if val, exists := event["branch"]; exists {
			eventConfig["branch"] = val
		}
		if val, exists := event["pullrequestCreatedBy"]; exists {
			eventConfig["pull_request_created_by"] = val
		}
		if val, exists := event["pullrequestReviewersContains"]; exists {
			eventConfig["pull_request_reviewers_contains"] = val
		}
		if val, exists := event["mergeResult"]; exists {
			eventConfig["merge_result"] = val
		}

	case "git_pull_request_commented":
		if val, exists := event["repository"]; exists {
			eventConfig["repository_id"] = val
		}
		if val, exists := event["branch"]; exists {
			eventConfig["branch"] = val
		}

	case "repository_created":
		if val, exists := event["projectId"]; exists {
			eventConfig["project_id"] = val
		}

	case "repository_deleted", "repository_forked", "repository_renamed", "repository_status_changed":
		if val, exists := event["repository"]; exists {
			eventConfig["repository_id"] = val
		}

	case "tfvc_checkin":
		if val, exists := event["path"]; exists {
			eventConfig["path"] = val
		}

	case "work_item_created", "work_item_deleted", "work_item_restored":
		if val, exists := event["workItemType"]; exists {
			eventConfig["work_item_type"] = val
		}
		if val, exists := event["areaPath"]; exists {
			eventConfig["area_path"] = val
		}
		if val, exists := event["tag"]; exists {
			eventConfig["tag"] = val
		}

	case "work_item_updated":
		if val, exists := event["workItemType"]; exists {
			eventConfig["work_item_type"] = val
		}
		if val, exists := event["areaPath"]; exists {
			eventConfig["area_path"] = val
		}
		if val, exists := event["tag"]; exists {
			eventConfig["tag"] = val
		}
		if val, exists := event["changedFields"]; exists {
			eventConfig["changed_fields"] = val
		}

	case "work_item_commented":
		if val, exists := event["workItemType"]; exists {
			eventConfig["work_item_type"] = val
		}
		if val, exists := event["areaPath"]; exists {
			eventConfig["area_path"] = val
		}
		if val, exists := event["tag"]; exists {
			eventConfig["tag"] = val
		}
		if val, exists := event["commentPattern"]; exists {
			eventConfig["comment_pattern"] = val
		}

	case "service_connection_created", "service_connection_updated":
		if val, exists := event["project"]; exists {
			eventConfig["project"] = val
		}
	}

	return eventType, []interface{}{eventConfig}
}

// isNilTfsEventConfig checks if TFS event config is empty
func isNilTfsEventConfig(eventConfig map[string]string) bool {
	for key := range eventConfig {
		if key != "projectId" && key != "tfsSubscriptionId" && eventConfig[key] != "" {
			return false
		}
	}
	return true
}

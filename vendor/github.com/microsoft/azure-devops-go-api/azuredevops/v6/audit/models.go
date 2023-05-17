// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package audit

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
)

// Defines all the categories an AuditAction can be
type AuditActionCategory string

type auditActionCategoryValuesType struct {
	Unknown AuditActionCategory
	Modify  AuditActionCategory
	Remove  AuditActionCategory
	Create  AuditActionCategory
	Access  AuditActionCategory
	Execute AuditActionCategory
}

var AuditActionCategoryValues = auditActionCategoryValuesType{
	// The category is not known
	Unknown: "unknown",
	// An artifact has been Modified
	Modify: "modify",
	// An artifact has been Removed
	Remove: "remove",
	// An artifact has been Created
	Create: "create",
	// An artifact has been Accessed
	Access: "access",
	// An artifact has been Executed
	Execute: "execute",
}

type AuditActionInfo struct {
	// The action id for the event, i.e Git.CreateRepo, Project.RenameProject
	ActionId *string `json:"actionId,omitempty"`
	// Area of Azure DevOps the action occurred
	Area *string `json:"area,omitempty"`
	// Type of action executed
	Category *AuditActionCategory `json:"category,omitempty"`
}

// The object returned when the audit log is queried. It contains the log and the information needed to query more audit entries.
type AuditLogQueryResult struct {
	// The continuation token to pass to get the next set of results
	ContinuationToken *string `json:"continuationToken,omitempty"`
	// The list of audit log entries
	DecoratedAuditLogEntries *[]DecoratedAuditLogEntry `json:"decoratedAuditLogEntries,omitempty"`
	// True when there are more matching results to be fetched, false otherwise.
	HasMore *bool `json:"hasMore,omitempty"`
}

// This class represents an audit stream
type AuditStream struct {
	// Inputs used to communicate with external service. Inputs could be url, a connection string, a token, etc.
	ConsumerInputs *map[string]string `json:"consumerInputs,omitempty"`
	// Type of the consumer, i.e. splunk, azureEventHub, etc.
	ConsumerType *string `json:"consumerType,omitempty"`
	// The time when the stream was created
	CreatedTime *azuredevops.Time `json:"createdTime,omitempty"`
	// Used to identify individual streams
	DisplayName *string `json:"displayName,omitempty"`
	// Unique stream identifier
	Id *int `json:"id,omitempty"`
	// Status of the stream, Enabled, Disabled
	Status *AuditStreamStatus `json:"status,omitempty"`
	// Reason for the current stream status, i.e. Disabled by the system, Invalid credentials, etc.
	StatusReason *string `json:"statusReason,omitempty"`
	// The time when the stream was last updated
	UpdatedTime *azuredevops.Time `json:"updatedTime,omitempty"`
}

// Represents the status of a stream
type AuditStreamStatus string

type auditStreamStatusValuesType struct {
	Unknown          AuditStreamStatus
	Enabled          AuditStreamStatus
	DisabledByUser   AuditStreamStatus
	DisabledBySystem AuditStreamStatus
	Deleted          AuditStreamStatus
	Backfilling      AuditStreamStatus
}

var AuditStreamStatusValues = auditStreamStatusValuesType{
	// The state has not been set, The stream is new
	Unknown: "unknown",
	// The stream is enabled and can deliver events
	Enabled: "enabled",
	// The stream has been disabled by a user
	DisabledByUser: "disabledByUser",
	// The stream has been disabled by the system
	DisabledBySystem: "disabledBySystem",
	// The stream has been marked for deletion
	Deleted: "deleted",
	// The stream is delivering old events
	Backfilling: "backfilling",
}

type DecoratedAuditLogEntry struct {
	// The action id for the event, i.e Git.CreateRepo, Project.RenameProject
	ActionId *string `json:"actionId,omitempty"`
	// ActivityId
	ActivityId *uuid.UUID `json:"activityId,omitempty"`
	// The Actor's CUID
	ActorCUID *uuid.UUID `json:"actorCUID,omitempty"`
	// DisplayName of the user who initiated the action
	ActorDisplayName *string `json:"actorDisplayName,omitempty"`
	// URL of Actor's Profile image
	ActorImageUrl *string `json:"actorImageUrl,omitempty"`
	// The Actor's UPN
	ActorUPN *string `json:"actorUPN,omitempty"`
	// The Actor's User Id
	ActorUserId *uuid.UUID `json:"actorUserId,omitempty"`
	// Area of Azure DevOps the action occurred
	Area *string `json:"area,omitempty"`
	// Type of authentication used by the actor
	AuthenticationMechanism *string `json:"authenticationMechanism,omitempty"`
	// Type of action executed
	Category *AuditActionCategory `json:"category,omitempty"`
	// DisplayName of the category
	CategoryDisplayName *string `json:"categoryDisplayName,omitempty"`
	// This allows related audit entries to be grouped together. Generally this occurs when a single action causes a cascade of audit entries. For example, project creation.
	CorrelationId *uuid.UUID `json:"correlationId,omitempty"`
	// External data such as CUIDs, item names, etc.
	Data *map[string]interface{} `json:"data,omitempty"`
	// Decorated details
	Details *string `json:"details,omitempty"`
	// EventId - Needs to be unique per service
	Id *string `json:"id,omitempty"`
	// IP Address where the event was originated
	IpAddress *string `json:"ipAddress,omitempty"`
	// When specified, the id of the project this event is associated to
	ProjectId *uuid.UUID `json:"projectId,omitempty"`
	// When specified, the name of the project this event is associated to
	ProjectName *string `json:"projectName,omitempty"`
	// DisplayName of the scope
	ScopeDisplayName *string `json:"scopeDisplayName,omitempty"`
	// The organization Id (Organization is the only scope currently supported)
	ScopeId *uuid.UUID `json:"scopeId,omitempty"`
	// The type of the scope (Organization is only scope currently supported)
	ScopeType *string `json:"scopeType,omitempty"`
	// The time when the event occurred in UTC
	Timestamp *azuredevops.Time `json:"timestamp,omitempty"`
	// The user agent from the request
	UserAgent *string `json:"userAgent,omitempty"`
}

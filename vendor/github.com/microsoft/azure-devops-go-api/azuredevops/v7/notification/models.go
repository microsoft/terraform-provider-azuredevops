// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package notification

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/forminput"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
)

type ActorFilter struct {
	EventType  *string                `json:"eventType,omitempty"`
	Criteria   *ExpressionFilterModel `json:"criteria,omitempty"`
	Type       *string                `json:"type,omitempty"`
	Exclusions *[]string              `json:"exclusions,omitempty"`
	Inclusions *[]string              `json:"inclusions,omitempty"`
}

type ActorNotificationReason struct {
	NotificationReasonType *NotificationReasonType `json:"notificationReasonType,omitempty"`
	TargetIdentities       *[]webapi.IdentityRef   `json:"targetIdentities,omitempty"`
	MatchedRoles           *[]string               `json:"matchedRoles,omitempty"`
}

// Artifact filter options. Used in "follow" subscriptions.
type ArtifactFilter struct {
	EventType    *string `json:"eventType,omitempty"`
	ArtifactId   *string `json:"artifactId,omitempty"`
	ArtifactType *string `json:"artifactType,omitempty"`
	ArtifactUri  *string `json:"artifactUri,omitempty"`
	Type         *string `json:"type,omitempty"`
}

type BaseSubscriptionFilter struct {
	EventType *string `json:"eventType,omitempty"`
	Type      *string `json:"type,omitempty"`
}

type BatchNotificationOperation struct {
	NotificationOperation       *NotificationOperation        `json:"notificationOperation,omitempty"`
	NotificationQueryConditions *[]NotificationQueryCondition `json:"notificationQueryConditions,omitempty"`
}

type BlockFilter struct {
	EventType  *string                `json:"eventType,omitempty"`
	Criteria   *ExpressionFilterModel `json:"criteria,omitempty"`
	Type       *string                `json:"type,omitempty"`
	Exclusions *[]string              `json:"exclusions,omitempty"`
	Inclusions *[]string              `json:"inclusions,omitempty"`
}

type BlockSubscriptionChannel struct {
	Type *string `json:"type,omitempty"`
}

// Default delivery preference for group subscribers. Indicates how the subscriber should be notified.
type DefaultGroupDeliveryPreference string

type defaultGroupDeliveryPreferenceValuesType struct {
	NoDelivery DefaultGroupDeliveryPreference
	EachMember DefaultGroupDeliveryPreference
}

var DefaultGroupDeliveryPreferenceValues = defaultGroupDeliveryPreferenceValuesType{
	NoDelivery: "noDelivery",
	EachMember: "eachMember",
}

type DiagnosticIdentity struct {
	DisplayName  *string    `json:"displayName,omitempty"`
	EmailAddress *string    `json:"emailAddress,omitempty"`
	Id           *uuid.UUID `json:"id,omitempty"`
}

type DiagnosticNotification struct {
	EventId        *int                                `json:"eventId,omitempty"`
	EventType      *string                             `json:"eventType,omitempty"`
	Id             *int                                `json:"id,omitempty"`
	Messages       *[]NotificationDiagnosticLogMessage `json:"messages,omitempty"`
	Recipients     *map[uuid.UUID]DiagnosticRecipient  `json:"recipients,omitempty"`
	Result         *string                             `json:"result,omitempty"`
	Stats          *map[string]int                     `json:"stats,omitempty"`
	SubscriptionId *string                             `json:"subscriptionId,omitempty"`
}

type DiagnosticRecipient struct {
	Recipient *DiagnosticIdentity `json:"recipient,omitempty"`
	Status    *string             `json:"status,omitempty"`
}

type EmailHtmlSubscriptionChannel struct {
	Address          *string `json:"address,omitempty"`
	UseCustomAddress *bool   `json:"useCustomAddress,omitempty"`
	Type             *string `json:"type,omitempty"`
}

type EmailPlaintextSubscriptionChannel struct {
	Address          *string `json:"address,omitempty"`
	UseCustomAddress *bool   `json:"useCustomAddress,omitempty"`
	Type             *string `json:"type,omitempty"`
}

// Describes the subscription evaluation operation status.
type EvaluationOperationStatus string

type evaluationOperationStatusValuesType struct {
	NotSet     EvaluationOperationStatus
	Queued     EvaluationOperationStatus
	InProgress EvaluationOperationStatus
	Cancelled  EvaluationOperationStatus
	Succeeded  EvaluationOperationStatus
	Failed     EvaluationOperationStatus
	TimedOut   EvaluationOperationStatus
	NotFound   EvaluationOperationStatus
}

var EvaluationOperationStatusValues = evaluationOperationStatusValuesType{
	// The operation object does not have the status set.
	NotSet: "notSet",
	// The operation has been queued.
	Queued: "queued",
	// The operation is in progress.
	InProgress: "inProgress",
	// The operation was cancelled by the user.
	Cancelled: "cancelled",
	// The operation completed successfully.
	Succeeded: "succeeded",
	// The operation completed with a failure.
	Failed: "failed",
	// The operation timed out.
	TimedOut: "timedOut",
	// The operation could not be found.
	NotFound: "notFound",
}

type EventBacklogStatus struct {
	CaptureTime             *azuredevops.Time `json:"captureTime,omitempty"`
	JobId                   *uuid.UUID        `json:"jobId,omitempty"`
	LastEventBatchStartTime *azuredevops.Time `json:"lastEventBatchStartTime,omitempty"`
	LastEventProcessedTime  *azuredevops.Time `json:"lastEventProcessedTime,omitempty"`
	LastJobBatchStartTime   *azuredevops.Time `json:"lastJobBatchStartTime,omitempty"`
	LastJobProcessedTime    *azuredevops.Time `json:"lastJobProcessedTime,omitempty"`
	OldestPendingEventTime  *azuredevops.Time `json:"oldestPendingEventTime,omitempty"`
	Publisher               *string           `json:"publisher,omitempty"`
	UnprocessedEvents       *int              `json:"unprocessedEvents,omitempty"`
}

type EventBatch struct {
	EndTime             interface{}     `json:"endTime,omitempty"`
	EventCounts         *map[string]int `json:"eventCounts,omitempty"`
	EventIds            *string         `json:"eventIds,omitempty"`
	NotificationCounts  *map[string]int `json:"notificationCounts,omitempty"`
	PreProcessEndTime   interface{}     `json:"preProcessEndTime,omitempty"`
	PreProcessStartTime interface{}     `json:"preProcessStartTime,omitempty"`
	ProcessEndTime      interface{}     `json:"processEndTime,omitempty"`
	ProcessStartTime    interface{}     `json:"processStartTime,omitempty"`
	StartTime           interface{}     `json:"startTime,omitempty"`
	SubscriptionCounts  *map[string]int `json:"subscriptionCounts,omitempty"`
}

type EventProcessingLog struct {
	// Identifier used for correlating to other diagnostics that may have been recorded elsewhere.
	ActivityId  *uuid.UUID        `json:"activityId,omitempty"`
	Description *string           `json:"description,omitempty"`
	EndTime     *azuredevops.Time `json:"endTime,omitempty"`
	Errors      *int              `json:"errors,omitempty"`
	// Unique instance identifier.
	Id         *uuid.UUID                          `json:"id,omitempty"`
	LogType    *string                             `json:"logType,omitempty"`
	Messages   *[]NotificationDiagnosticLogMessage `json:"messages,omitempty"`
	Properties *map[string]string                  `json:"properties,omitempty"`
	// This identifier depends on the logType.  For notification jobs, this will be the job Id. For subscription tracing, this will be a special root Guid with the subscription Id encoded.
	Source         *uuid.UUID                 `json:"source,omitempty"`
	StartTime      *azuredevops.Time          `json:"startTime,omitempty"`
	Warnings       *int                       `json:"warnings,omitempty"`
	Result         *string                    `json:"result,omitempty"`
	Stats          *map[string]map[string]int `json:"stats,omitempty"`
	Batches        *[]EventBatch              `json:"batches,omitempty"`
	MatcherResults *[]MatcherResult           `json:"matcherResults,omitempty"`
}

// [Flags] Set of flags used to determine which set of information is retrieved when querying for event publishers
type EventPublisherQueryFlags string

type eventPublisherQueryFlagsValuesType struct {
	None                  EventPublisherQueryFlags
	IncludeRemoteServices EventPublisherQueryFlags
}

var EventPublisherQueryFlagsValues = eventPublisherQueryFlagsValuesType{
	None: "none",
	// Include event types from the remote services too
	IncludeRemoteServices: "includeRemoteServices",
}

// Encapsulates events result properties. It defines the total number of events used and the number of matched events.
type EventsEvaluationResult struct {
	// Count of events evaluated.
	Count *int `json:"count,omitempty"`
	// Count of matched events.
	MatchedCount *int `json:"matchedCount,omitempty"`
}

// A transform request specify the properties of a notification event to be transformed.
type EventTransformRequest struct {
	// Event payload.
	EventPayload *string `json:"eventPayload,omitempty"`
	// Event type.
	EventType *string `json:"eventType,omitempty"`
	// System inputs.
	SystemInputs *map[string]string `json:"systemInputs,omitempty"`
}

// Result of transforming a notification event.
type EventTransformResult struct {
	// Transformed html content.
	Content *string `json:"content,omitempty"`
	// Calculated data.
	Data interface{} `json:"data,omitempty"`
	// Calculated system inputs.
	SystemInputs *map[string]string `json:"systemInputs,omitempty"`
}

// [Flags] Set of flags used to determine which set of information is retrieved when querying for eventtypes
type EventTypeQueryFlags string

type eventTypeQueryFlagsValuesType struct {
	None          EventTypeQueryFlags
	IncludeFields EventTypeQueryFlags
}

var EventTypeQueryFlagsValues = eventTypeQueryFlagsValuesType{
	None: "none",
	// IncludeFields will include all fields and their types
	IncludeFields: "includeFields",
}

type ExpressionFilter struct {
	EventType *string                `json:"eventType,omitempty"`
	Criteria  *ExpressionFilterModel `json:"criteria,omitempty"`
	Type      *string                `json:"type,omitempty"`
}

// Subscription Filter Clause represents a single clause in a subscription filter e.g. If the subscription has the following criteria "Project Name = [Current Project] AND Assigned To = [Me] it will be represented as two Filter Clauses Clause 1: Index = 1, Logical Operator: NULL  , FieldName = 'Project Name', Operator = '=', Value = '[Current Project]' Clause 2: Index = 2, Logical Operator: 'AND' , FieldName = 'Assigned To' , Operator = '=', Value = '[Me]'
type ExpressionFilterClause struct {
	FieldName *string `json:"fieldName,omitempty"`
	// The order in which this clause appeared in the filter query
	Index *int `json:"index,omitempty"`
	// Logical Operator 'AND', 'OR' or NULL (only for the first clause in the filter)
	LogicalOperator *string `json:"logicalOperator,omitempty"`
	Operator        *string `json:"operator,omitempty"`
	Value           *string `json:"value,omitempty"`
}

// Represents a hierarchy of SubscritionFilterClauses that have been grouped together through either adding a group in the WebUI or using parethesis in the Subscription condition string
type ExpressionFilterGroup struct {
	// The index of the last FilterClause in this group
	End *int `json:"end,omitempty"`
	// Level of the group, since groups can be nested for each nested group the level will increase by 1
	Level *int `json:"level,omitempty"`
	// The index of the first FilterClause in this group
	Start *int `json:"start,omitempty"`
}

type ExpressionFilterModel struct {
	// Flat list of clauses in this subscription
	Clauses *[]ExpressionFilterClause `json:"clauses,omitempty"`
	// Grouping of clauses in the subscription
	Groups *[]ExpressionFilterGroup `json:"groups,omitempty"`
	// Max depth of the Subscription tree
	MaxGroupLevel *int `json:"maxGroupLevel,omitempty"`
}

type FieldInputValues struct {
	// The default value to use for this input
	DefaultValue *string `json:"defaultValue,omitempty"`
	// Errors encountered while computing dynamic values.
	Error *forminput.InputValuesError `json:"error,omitempty"`
	// The id of the input
	InputId *string `json:"inputId,omitempty"`
	// Should this input be disabled
	IsDisabled *bool `json:"isDisabled,omitempty"`
	// Should the value be restricted to one of the values in the PossibleValues (True) or are the values in PossibleValues just a suggestion (False)
	IsLimitedToPossibleValues *bool `json:"isLimitedToPossibleValues,omitempty"`
	// Should this input be made read-only
	IsReadOnly *bool `json:"isReadOnly,omitempty"`
	// Possible values that this input can take
	PossibleValues *[]forminput.InputValue `json:"possibleValues,omitempty"`
	Operators      *[]byte                 `json:"operators,omitempty"`
}

type FieldValuesQuery struct {
	CurrentValues *map[string]string `json:"currentValues,omitempty"`
	// Subscription containing information about the publisher/consumer and the current input values
	Resource    interface{}         `json:"resource,omitempty"`
	InputValues *[]FieldInputValues `json:"inputValues,omitempty"`
	Scope       *string             `json:"scope,omitempty"`
}

type GeneratedNotification struct {
	Recipients *[]DiagnosticIdentity `json:"recipients,omitempty"`
}

type GroupSubscriptionChannel struct {
	Address          *string `json:"address,omitempty"`
	UseCustomAddress *bool   `json:"useCustomAddress,omitempty"`
	Type             *string `json:"type,omitempty"`
}

// Abstraction interface for the diagnostic log.  Primarily for deserialization.
type INotificationDiagnosticLog struct {
	// Identifier used for correlating to other diagnostics that may have been recorded elsewhere.
	ActivityId *uuid.UUID `json:"activityId,omitempty"`
	// Description of what subscription or notification job is being logged.
	Description *string `json:"description,omitempty"`
	// Time the log ended.
	EndTime *azuredevops.Time `json:"endTime,omitempty"`
	// Unique instance identifier.
	Id *uuid.UUID `json:"id,omitempty"`
	// Type of information being logged.
	LogType *string `json:"logType,omitempty"`
	// List of log messages.
	Messages *[]NotificationDiagnosticLogMessage `json:"messages,omitempty"`
	// Dictionary of log properties and settings for the job.
	Properties *map[string]string `json:"properties,omitempty"`
	// This identifier depends on the logType.  For notification jobs, this will be the job Id. For subscription tracing, this will be a special root Guid with the subscription Id encoded.
	Source *uuid.UUID `json:"source,omitempty"`
	// Time the log started.
	StartTime *azuredevops.Time `json:"startTime,omitempty"`
}

type ISubscriptionFilter struct {
	EventType *string `json:"eventType,omitempty"`
	Type      *string `json:"type,omitempty"`
}

type ISubscriptionChannel struct {
	Type *string `json:"type,omitempty"`
}

type MatcherResult struct {
	Matcher *string                    `json:"matcher,omitempty"`
	Stats   *map[string]map[string]int `json:"stats,omitempty"`
}

type MessageQueueSubscriptionChannel struct {
	Type *string `json:"type,omitempty"`
}

type NotificationAdminSettings struct {
	// The default group delivery preference for groups in this collection
	DefaultGroupDeliveryPreference *DefaultGroupDeliveryPreference `json:"defaultGroupDeliveryPreference,omitempty"`
}

type NotificationAdminSettingsUpdateParameters struct {
	DefaultGroupDeliveryPreference *DefaultGroupDeliveryPreference `json:"defaultGroupDeliveryPreference,omitempty"`
}

type NotificationBacklogStatus struct {
	CaptureTime                    *azuredevops.Time `json:"captureTime,omitempty"`
	Channel                        *string           `json:"channel,omitempty"`
	JobId                          *uuid.UUID        `json:"jobId,omitempty"`
	LastJobBatchStartTime          *azuredevops.Time `json:"lastJobBatchStartTime,omitempty"`
	LastJobProcessedTime           *azuredevops.Time `json:"lastJobProcessedTime,omitempty"`
	LastNotificationBatchStartTime *azuredevops.Time `json:"lastNotificationBatchStartTime,omitempty"`
	LastNotificationProcessedTime  *azuredevops.Time `json:"lastNotificationProcessedTime,omitempty"`
	OldestPendingNotificationTime  *azuredevops.Time `json:"oldestPendingNotificationTime,omitempty"`
	Publisher                      *string           `json:"publisher,omitempty"`
	// Null status is unprocessed
	Status                   *string `json:"status,omitempty"`
	UnprocessedNotifications *int    `json:"unprocessedNotifications,omitempty"`
}

type NotificationBatch struct {
	EndTime                  interface{}               `json:"endTime,omitempty"`
	NotificationCount        *int                      `json:"notificationCount,omitempty"`
	NotificationIds          *string                   `json:"notificationIds,omitempty"`
	ProblematicNotifications *[]DiagnosticNotification `json:"problematicNotifications,omitempty"`
	StartTime                interface{}               `json:"startTime,omitempty"`
}

type NotificationDeliveryLog struct {
	// Identifier used for correlating to other diagnostics that may have been recorded elsewhere.
	ActivityId  *uuid.UUID        `json:"activityId,omitempty"`
	Description *string           `json:"description,omitempty"`
	EndTime     *azuredevops.Time `json:"endTime,omitempty"`
	Errors      *int              `json:"errors,omitempty"`
	// Unique instance identifier.
	Id         *uuid.UUID                          `json:"id,omitempty"`
	LogType    *string                             `json:"logType,omitempty"`
	Messages   *[]NotificationDiagnosticLogMessage `json:"messages,omitempty"`
	Properties *map[string]string                  `json:"properties,omitempty"`
	// This identifier depends on the logType.  For notification jobs, this will be the job Id. For subscription tracing, this will be a special root Guid with the subscription Id encoded.
	Source    *uuid.UUID                 `json:"source,omitempty"`
	StartTime *azuredevops.Time          `json:"startTime,omitempty"`
	Warnings  *int                       `json:"warnings,omitempty"`
	Result    *string                    `json:"result,omitempty"`
	Stats     *map[string]map[string]int `json:"stats,omitempty"`
	Batches   *[]NotificationBatch       `json:"batches,omitempty"`
}

// Abstract base class for all of the diagnostic logs.
type NotificationDiagnosticLog struct {
	// Identifier used for correlating to other diagnostics that may have been recorded elsewhere.
	ActivityId  *uuid.UUID        `json:"activityId,omitempty"`
	Description *string           `json:"description,omitempty"`
	EndTime     *azuredevops.Time `json:"endTime,omitempty"`
	Errors      *int              `json:"errors,omitempty"`
	// Unique instance identifier.
	Id         *uuid.UUID                          `json:"id,omitempty"`
	LogType    *string                             `json:"logType,omitempty"`
	Messages   *[]NotificationDiagnosticLogMessage `json:"messages,omitempty"`
	Properties *map[string]string                  `json:"properties,omitempty"`
	// This identifier depends on the logType.  For notification jobs, this will be the job Id. For subscription tracing, this will be a special root Guid with the subscription Id encoded.
	Source    *uuid.UUID        `json:"source,omitempty"`
	StartTime *azuredevops.Time `json:"startTime,omitempty"`
	Warnings  *int              `json:"warnings,omitempty"`
}

type NotificationDiagnosticLogMessage struct {
	// Corresponds to .Net TraceLevel enumeration
	Level   *int        `json:"level,omitempty"`
	Message *string     `json:"message,omitempty"`
	Time    interface{} `json:"time,omitempty"`
}

type NotificationEventBacklogStatus struct {
	EventBacklogStatus        *[]EventBacklogStatus        `json:"eventBacklogStatus,omitempty"`
	NotificationBacklogStatus *[]NotificationBacklogStatus `json:"notificationBacklogStatus,omitempty"`
}

// Encapsulates the properties of a filterable field. A filterable field is a field in an event that can used to filter notifications for a certain event type.
type NotificationEventField struct {
	// Gets or sets the type of this field.
	FieldType *NotificationEventFieldType `json:"fieldType,omitempty"`
	// Gets or sets the unique identifier of this field.
	Id *string `json:"id,omitempty"`
	// Gets or sets the name of this field.
	Name *string `json:"name,omitempty"`
	// Gets or sets the path to the field in the event object. This path can be either Json Path or XPath, depending on if the event will be serialized into Json or XML
	Path *string `json:"path,omitempty"`
	// Gets or sets the scopes that this field supports. If not specified then the event type scopes apply.
	SupportedScopes *[]string `json:"supportedScopes,omitempty"`
}

// Encapsulates the properties of a field type. It includes a unique id for the operator and a localized string for display name
type NotificationEventFieldOperator struct {
	// Gets or sets the display name of an operator
	DisplayName *string `json:"displayName,omitempty"`
	// Gets or sets the id of an operator
	Id *string `json:"id,omitempty"`
}

// Encapsulates the properties of a field type. It describes the data type of a field, the operators it support and how to populate it in the UI
type NotificationEventFieldType struct {
	// Gets or sets the unique identifier of this field type.
	Id                  *string               `json:"id,omitempty"`
	OperatorConstraints *[]OperatorConstraint `json:"operatorConstraints,omitempty"`
	// Gets or sets the list of operators that this type supports.
	Operators             *[]NotificationEventFieldOperator `json:"operators,omitempty"`
	SubscriptionFieldType *SubscriptionFieldType            `json:"subscriptionFieldType,omitempty"`
	// Gets or sets the value definition of this field like the getValuesMethod and template to display in the UI
	Value *ValueDefinition `json:"value,omitempty"`
}

// Encapsulates the properties of a notification event publisher.
type NotificationEventPublisher struct {
	Id                         *string                 `json:"id,omitempty"`
	SubscriptionManagementInfo *SubscriptionManagement `json:"subscriptionManagementInfo,omitempty"`
	Url                        *string                 `json:"url,omitempty"`
}

// Encapsulates the properties of an event role.  An event Role is used for role based subscription for example for a buildCompletedEvent, one role is request by field
type NotificationEventRole struct {
	// Gets or sets an Id for that role, this id is used by the event.
	Id *string `json:"id,omitempty"`
	// Gets or sets the Name for that role, this name is used for UI display.
	Name *string `json:"name,omitempty"`
	// Gets or sets whether this role can be a group or just an individual user
	SupportsGroups *bool `json:"supportsGroups,omitempty"`
}

// Encapsulates the properties of an event type. It defines the fields, that can be used for filtering, for that event type.
type NotificationEventType struct {
	Category *NotificationEventTypeCategory `json:"category,omitempty"`
	// Gets or sets the color representing this event type. Example: rgb(128,245,211) or #fafafa
	Color                      *string                            `json:"color,omitempty"`
	CustomSubscriptionsAllowed *bool                              `json:"customSubscriptionsAllowed,omitempty"`
	EventPublisher             *NotificationEventPublisher        `json:"eventPublisher,omitempty"`
	Fields                     *map[string]NotificationEventField `json:"fields,omitempty"`
	HasInitiator               *bool                              `json:"hasInitiator,omitempty"`
	// Gets or sets the icon representing this event type. Can be a URL or a CSS class. Example: css://some-css-class
	Icon *string `json:"icon,omitempty"`
	// Gets or sets the unique identifier of this event definition.
	Id *string `json:"id,omitempty"`
	// Gets or sets the name of this event definition.
	Name  *string                  `json:"name,omitempty"`
	Roles *[]NotificationEventRole `json:"roles,omitempty"`
	// Gets or sets the scopes that this event type supports
	SupportedScopes *[]string `json:"supportedScopes,omitempty"`
	// Gets or sets the rest end point to get this event type details (fields, fields types)
	Url *string `json:"url,omitempty"`
}

// Encapsulates the properties of a category. A category will be used by the UI to group event types
type NotificationEventTypeCategory struct {
	// Gets or sets the unique identifier of this category.
	Id *string `json:"id,omitempty"`
	// Gets or sets the friendly name of this category.
	Name *string `json:"name,omitempty"`
}

type NotificationJobDiagnosticLog struct {
	// Identifier used for correlating to other diagnostics that may have been recorded elsewhere.
	ActivityId  *uuid.UUID        `json:"activityId,omitempty"`
	Description *string           `json:"description,omitempty"`
	EndTime     *azuredevops.Time `json:"endTime,omitempty"`
	Errors      *int              `json:"errors,omitempty"`
	// Unique instance identifier.
	Id         *uuid.UUID                          `json:"id,omitempty"`
	LogType    *string                             `json:"logType,omitempty"`
	Messages   *[]NotificationDiagnosticLogMessage `json:"messages,omitempty"`
	Properties *map[string]string                  `json:"properties,omitempty"`
	// This identifier depends on the logType.  For notification jobs, this will be the job Id. For subscription tracing, this will be a special root Guid with the subscription Id encoded.
	Source    *uuid.UUID                 `json:"source,omitempty"`
	StartTime *azuredevops.Time          `json:"startTime,omitempty"`
	Warnings  *int                       `json:"warnings,omitempty"`
	Result    *string                    `json:"result,omitempty"`
	Stats     *map[string]map[string]int `json:"stats,omitempty"`
}

type NotificationOperation string

type notificationOperationValuesType struct {
	None               NotificationOperation
	SuspendUnprocessed NotificationOperation
}

var NotificationOperationValues = notificationOperationValuesType{
	None:               "none",
	SuspendUnprocessed: "suspendUnprocessed",
}

type NotificationQueryCondition struct {
	EventInitiator *uuid.UUID `json:"eventInitiator,omitempty"`
	EventType      *string    `json:"eventType,omitempty"`
	Subscriber     *uuid.UUID `json:"subscriber,omitempty"`
	SubscriptionId *string    `json:"subscriptionId,omitempty"`
}

type NotificationReason struct {
	NotificationReasonType *NotificationReasonType `json:"notificationReasonType,omitempty"`
	TargetIdentities       *[]webapi.IdentityRef   `json:"targetIdentities,omitempty"`
}

type NotificationReasonType string

type notificationReasonTypeValuesType struct {
	Unknown                 NotificationReasonType
	Follows                 NotificationReasonType
	Personal                NotificationReasonType
	PersonalAlias           NotificationReasonType
	DirectMember            NotificationReasonType
	IndirectMember          NotificationReasonType
	GroupAlias              NotificationReasonType
	SubscriptionAlias       NotificationReasonType
	SingleRole              NotificationReasonType
	DirectMemberGroupRole   NotificationReasonType
	InDirectMemberGroupRole NotificationReasonType
	AliasMemberGroupRole    NotificationReasonType
}

var NotificationReasonTypeValues = notificationReasonTypeValuesType{
	Unknown:                 "unknown",
	Follows:                 "follows",
	Personal:                "personal",
	PersonalAlias:           "personalAlias",
	DirectMember:            "directMember",
	IndirectMember:          "indirectMember",
	GroupAlias:              "groupAlias",
	SubscriptionAlias:       "subscriptionAlias",
	SingleRole:              "singleRole",
	DirectMemberGroupRole:   "directMemberGroupRole",
	InDirectMemberGroupRole: "inDirectMemberGroupRole",
	AliasMemberGroupRole:    "aliasMemberGroupRole",
}

// Encapsulates notifications result properties. It defines the number of notifications and the recipients of notifications.
type NotificationsEvaluationResult struct {
	// Count of generated notifications
	Count *int `json:"count,omitempty"`
}

type NotificationStatistic struct {
	Date     *azuredevops.Time          `json:"date,omitempty"`
	HitCount *int                       `json:"hitCount,omitempty"`
	Path     *string                    `json:"path,omitempty"`
	Type     *NotificationStatisticType `json:"type,omitempty"`
	User     *webapi.IdentityRef        `json:"user,omitempty"`
}

type NotificationStatisticsQuery struct {
	Conditions *[]NotificationStatisticsQueryConditions `json:"conditions,omitempty"`
}

type NotificationStatisticsQueryConditions struct {
	EndDate         *azuredevops.Time          `json:"endDate,omitempty"`
	HitCountMinimum *int                       `json:"hitCountMinimum,omitempty"`
	Path            *string                    `json:"path,omitempty"`
	StartDate       *azuredevops.Time          `json:"startDate,omitempty"`
	Type            *NotificationStatisticType `json:"type,omitempty"`
	User            *webapi.IdentityRef        `json:"user,omitempty"`
}

type NotificationStatisticType string

type notificationStatisticTypeValuesType struct {
	NotificationBySubscription                             NotificationStatisticType
	EventsByEventType                                      NotificationStatisticType
	NotificationByEventType                                NotificationStatisticType
	EventsByEventTypePerUser                               NotificationStatisticType
	NotificationByEventTypePerUser                         NotificationStatisticType
	Events                                                 NotificationStatisticType
	Notifications                                          NotificationStatisticType
	NotificationFailureBySubscription                      NotificationStatisticType
	UnprocessedRangeStart                                  NotificationStatisticType
	UnprocessedEventsByPublisher                           NotificationStatisticType
	UnprocessedEventDelayByPublisher                       NotificationStatisticType
	UnprocessedNotificationsByChannelByPublisher           NotificationStatisticType
	UnprocessedNotificationDelayByChannelByPublisher       NotificationStatisticType
	DelayRangeStart                                        NotificationStatisticType
	TotalPipelineTime                                      NotificationStatisticType
	NotificationPipelineTime                               NotificationStatisticType
	EventPipelineTime                                      NotificationStatisticType
	HourlyRangeStart                                       NotificationStatisticType
	HourlyNotificationBySubscription                       NotificationStatisticType
	HourlyEventsByEventTypePerUser                         NotificationStatisticType
	HourlyEvents                                           NotificationStatisticType
	HourlyNotifications                                    NotificationStatisticType
	HourlyUnprocessedEventsByPublisher                     NotificationStatisticType
	HourlyUnprocessedEventDelayByPublisher                 NotificationStatisticType
	HourlyUnprocessedNotificationsByChannelByPublisher     NotificationStatisticType
	HourlyUnprocessedNotificationDelayByChannelByPublisher NotificationStatisticType
	HourlyTotalPipelineTime                                NotificationStatisticType
	HourlyNotificationPipelineTime                         NotificationStatisticType
	HourlyEventPipelineTime                                NotificationStatisticType
}

var NotificationStatisticTypeValues = notificationStatisticTypeValuesType{
	NotificationBySubscription:                             "notificationBySubscription",
	EventsByEventType:                                      "eventsByEventType",
	NotificationByEventType:                                "notificationByEventType",
	EventsByEventTypePerUser:                               "eventsByEventTypePerUser",
	NotificationByEventTypePerUser:                         "notificationByEventTypePerUser",
	Events:                                                 "events",
	Notifications:                                          "notifications",
	NotificationFailureBySubscription:                      "notificationFailureBySubscription",
	UnprocessedRangeStart:                                  "unprocessedRangeStart",
	UnprocessedEventsByPublisher:                           "unprocessedEventsByPublisher",
	UnprocessedEventDelayByPublisher:                       "unprocessedEventDelayByPublisher",
	UnprocessedNotificationsByChannelByPublisher:           "unprocessedNotificationsByChannelByPublisher",
	UnprocessedNotificationDelayByChannelByPublisher:       "unprocessedNotificationDelayByChannelByPublisher",
	DelayRangeStart:                                        "delayRangeStart",
	TotalPipelineTime:                                      "totalPipelineTime",
	NotificationPipelineTime:                               "notificationPipelineTime",
	EventPipelineTime:                                      "eventPipelineTime",
	HourlyRangeStart:                                       "hourlyRangeStart",
	HourlyNotificationBySubscription:                       "hourlyNotificationBySubscription",
	HourlyEventsByEventTypePerUser:                         "hourlyEventsByEventTypePerUser",
	HourlyEvents:                                           "hourlyEvents",
	HourlyNotifications:                                    "hourlyNotifications",
	HourlyUnprocessedEventsByPublisher:                     "hourlyUnprocessedEventsByPublisher",
	HourlyUnprocessedEventDelayByPublisher:                 "hourlyUnprocessedEventDelayByPublisher",
	HourlyUnprocessedNotificationsByChannelByPublisher:     "hourlyUnprocessedNotificationsByChannelByPublisher",
	HourlyUnprocessedNotificationDelayByChannelByPublisher: "hourlyUnprocessedNotificationDelayByChannelByPublisher",
	HourlyTotalPipelineTime:                                "hourlyTotalPipelineTime",
	HourlyNotificationPipelineTime:                         "hourlyNotificationPipelineTime",
	HourlyEventPipelineTime:                                "hourlyEventPipelineTime",
}

// A subscriber is a user or group that has the potential to receive notifications.
type NotificationSubscriber struct {
	// Indicates how the subscriber should be notified by default.
	DeliveryPreference *NotificationSubscriberDeliveryPreference `json:"deliveryPreference,omitempty"`
	Flags              *SubscriberFlags                          `json:"flags,omitempty"`
	// Identifier of the subscriber.
	Id *uuid.UUID `json:"id,omitempty"`
	// Preferred email address of the subscriber. A null or empty value indicates no preferred email address has been set.
	PreferredEmailAddress *string `json:"preferredEmailAddress,omitempty"`
}

// Delivery preference for a subscriber. Indicates how the subscriber should be notified.
type NotificationSubscriberDeliveryPreference string

type notificationSubscriberDeliveryPreferenceValuesType struct {
	NoDelivery            NotificationSubscriberDeliveryPreference
	PreferredEmailAddress NotificationSubscriberDeliveryPreference
	EachMember            NotificationSubscriberDeliveryPreference
	UseDefault            NotificationSubscriberDeliveryPreference
}

var NotificationSubscriberDeliveryPreferenceValues = notificationSubscriberDeliveryPreferenceValuesType{
	// Do not send notifications by default. Note: notifications can still be delivered to this subscriber, for example via a custom subscription.
	NoDelivery: "noDelivery",
	// Deliver notifications to the subscriber's preferred email address.
	PreferredEmailAddress: "preferredEmailAddress",
	// Deliver notifications to each member of the group representing the subscriber. Only applicable when the subscriber is a group.
	EachMember: "eachMember",
	// Use default
	UseDefault: "useDefault",
}

// Updates to a subscriber. Typically used to change (or set) a preferred email address or default delivery preference.
type NotificationSubscriberUpdateParameters struct {
	// New delivery preference for the subscriber (indicates how the subscriber should be notified).
	DeliveryPreference *NotificationSubscriberDeliveryPreference `json:"deliveryPreference,omitempty"`
	// New preferred email address for the subscriber. Specify an empty string to clear the current address.
	PreferredEmailAddress *string `json:"preferredEmailAddress,omitempty"`
}

// A subscription defines criteria for matching events and how the subscription's subscriber should be notified about those events.
type NotificationSubscription struct {
	// Links to related resources, APIs, and views for the subscription.
	Links interface{} `json:"_links,omitempty"`
	// Admin-managed settings for the subscription. Only applies when the subscriber is a group.
	AdminSettings *SubscriptionAdminSettings `json:"adminSettings,omitempty"`
	// Description of the subscription. Typically describes filter criteria which helps identity the subscription.
	Description *string `json:"description,omitempty"`
	// Diagnostics for this subscription.
	Diagnostics *SubscriptionDiagnostics `json:"diagnostics,omitempty"`
	// Any extra properties like detailed description for different contexts, user/group contexts
	ExtendedProperties *map[string]string `json:"extendedProperties,omitempty"`
	// Matching criteria for the subscription. ExpressionFilter
	Filter *ISubscriptionFilter `json:"filter,omitempty"`
	// Read-only indicators that further describe the subscription.
	Flags *SubscriptionFlags `json:"flags,omitempty"`
	// Channel for delivering notifications triggered by the subscription.
	Channel *ISubscriptionChannel `json:"channel,omitempty"`
	// Subscription identifier.
	Id *string `json:"id,omitempty"`
	// User that last modified (or created) the subscription.
	LastModifiedBy *webapi.IdentityRef `json:"lastModifiedBy,omitempty"`
	// Date when the subscription was last modified. If the subscription has not been updated since it was created, this value will indicate when the subscription was created.
	ModifiedDate *azuredevops.Time `json:"modifiedDate,omitempty"`
	// The permissions the user have for this subscriptions.
	Permissions *SubscriptionPermissions `json:"permissions,omitempty"`
	// The container in which events must be published from in order to be matched by the subscription. If empty, the scope is the current host (typically an account or project collection). For example, a subscription scoped to project A will not produce notifications for events published from project B.
	Scope *SubscriptionScope `json:"scope,omitempty"`
	// Status of the subscription. Typically indicates whether the subscription is enabled or not.
	Status *SubscriptionStatus `json:"status,omitempty"`
	// Message that provides more details about the status of the subscription.
	StatusMessage *string `json:"statusMessage,omitempty"`
	// User or group that will receive notifications for events matching the subscription's filter criteria.
	Subscriber *webapi.IdentityRef `json:"subscriber,omitempty"`
	// REST API URL of the subscription.
	Url *string `json:"url,omitempty"`
	// User-managed settings for the subscription. Only applies when the subscriber is a group. Typically used to indicate whether the calling user is opted in or out of a group subscription.
	UserSettings *SubscriptionUserSettings `json:"userSettings,omitempty"`
}

// Parameters for creating a new subscription. A subscription defines criteria for matching events and how the subscription's subscriber should be notified about those events.
type NotificationSubscriptionCreateParameters struct {
	// Brief description for the new subscription. Typically describes filter criteria which helps identity the subscription.
	Description *string `json:"description,omitempty"`
	// Matching criteria for the new subscription. ExpressionFilter
	Filter *ISubscriptionFilter `json:"filter,omitempty"`
	// Channel for delivering notifications triggered by the new subscription.
	Channel *ISubscriptionChannel `json:"channel,omitempty"`
	// The container in which events must be published from in order to be matched by the new subscription. If not specified, defaults to the current host (typically an account or project collection). For example, a subscription scoped to project A will not produce notifications for events published from project B.
	Scope *SubscriptionScope `json:"scope,omitempty"`
	// User or group that will receive notifications for events matching the subscription's filter criteria. If not specified, defaults to the calling user.
	Subscriber *webapi.IdentityRef `json:"subscriber,omitempty"`
}

type NotificationSubscriptionTemplate struct {
	Description                  *string                   `json:"description,omitempty"`
	Filter                       *ISubscriptionFilter      `json:"filter,omitempty"`
	Id                           *string                   `json:"id,omitempty"`
	NotificationEventInformation *NotificationEventType    `json:"notificationEventInformation,omitempty"`
	Type                         *SubscriptionTemplateType `json:"type,omitempty"`
}

// Parameters for updating an existing subscription. A subscription defines criteria for matching events and how the subscription's subscriber should be notified about those events. Note: only the fields to be updated should be set.
type NotificationSubscriptionUpdateParameters struct {
	// Admin-managed settings for the subscription. Only applies to subscriptions where the subscriber is a group.
	AdminSettings *SubscriptionAdminSettings `json:"adminSettings,omitempty"`
	// Updated description for the subscription. Typically describes filter criteria which helps identity the subscription.
	Description *string `json:"description,omitempty"`
	// Matching criteria for the subscription. ExpressionFilter
	Filter *ISubscriptionFilter `json:"filter,omitempty"`
	// Channel for delivering notifications triggered by the subscription.
	Channel *ISubscriptionChannel `json:"channel,omitempty"`
	// The container in which events must be published from in order to be matched by the new subscription. If not specified, defaults to the current host (typically the current account or project collection). For example, a subscription scoped to project A will not produce notifications for events published from project B.
	Scope *SubscriptionScope `json:"scope,omitempty"`
	// Updated status for the subscription. Typically used to enable or disable a subscription.
	Status *SubscriptionStatus `json:"status,omitempty"`
	// Optional message that provides more details about the updated status.
	StatusMessage *string `json:"statusMessage,omitempty"`
	// User-managed settings for the subscription. Only applies to subscriptions where the subscriber is a group. Typically used to opt-in or opt-out a user from a group subscription.
	UserSettings *SubscriptionUserSettings `json:"userSettings,omitempty"`
}

// Encapsulates the properties of an operator constraint. An operator constraint defines if some operator is available only for specific scope like a project scope.
type OperatorConstraint struct {
	Operator *string `json:"operator,omitempty"`
	// Gets or sets the list of scopes that this type supports.
	SupportedScopes *[]string `json:"supportedScopes,omitempty"`
}

type ProcessedEvent struct {
	// All of the users that were associated with this event and their role.
	Actors             *[]webapi.EventActor  `json:"actors,omitempty"`
	AllowedChannels    *string               `json:"allowedChannels,omitempty"`
	ArtifactUri        *string               `json:"artifactUri,omitempty"`
	DeliveryIdentities *ProcessingIdentities `json:"deliveryIdentities,omitempty"`
	// Evaluations for each user
	Evaluations *map[uuid.UUID]SubscriptionEvaluation `json:"evaluations,omitempty"`
	EventId     *int                                  `json:"eventId,omitempty"`
	// Which members were excluded from evaluation (only applies to ActorMatcher subscriptions)
	Exclusions *[]webapi.EventActor `json:"exclusions,omitempty"`
	// Which members were included for evaluation (only applies to ActorMatcher subscriptions)
	Inclusions    *[]webapi.EventActor     `json:"inclusions,omitempty"`
	Notifications *[]GeneratedNotification `json:"notifications,omitempty"`
}

type ProcessingDiagnosticIdentity struct {
	DisplayName        *string    `json:"displayName,omitempty"`
	EmailAddress       *string    `json:"emailAddress,omitempty"`
	Id                 *uuid.UUID `json:"id,omitempty"`
	DeliveryPreference *string    `json:"deliveryPreference,omitempty"`
	IsActive           *bool      `json:"isActive,omitempty"`
	IsGroup            *bool      `json:"isGroup,omitempty"`
	Message            *string    `json:"message,omitempty"`
}

type ProcessingIdentities struct {
	ExcludedIdentities *map[uuid.UUID]ProcessingDiagnosticIdentity `json:"excludedIdentities,omitempty"`
	IncludedIdentities *map[uuid.UUID]ProcessingDiagnosticIdentity `json:"includedIdentities,omitempty"`
	Messages           *[]NotificationDiagnosticLogMessage         `json:"messages,omitempty"`
	MissingIdentities  *[]uuid.UUID                                `json:"missingIdentities,omitempty"`
	Properties         *map[string]string                          `json:"properties,omitempty"`
}

type RoleBasedFilter struct {
	EventType  *string                `json:"eventType,omitempty"`
	Criteria   *ExpressionFilterModel `json:"criteria,omitempty"`
	Type       *string                `json:"type,omitempty"`
	Exclusions *[]string              `json:"exclusions,omitempty"`
	Inclusions *[]string              `json:"inclusions,omitempty"`
}

type ServiceBusSubscriptionChannel struct {
	Type *string `json:"type,omitempty"`
}

type ServiceHooksSubscriptionChannel struct {
	Type *string `json:"type,omitempty"`
}

type SoapSubscriptionChannel struct {
	Address          *string `json:"address,omitempty"`
	UseCustomAddress *bool   `json:"useCustomAddress,omitempty"`
	Type             *string `json:"type,omitempty"`
}

// [Flags]
type SubscriberFlags string

type subscriberFlagsValuesType struct {
	None                                  SubscriberFlags
	DeliveryPreferencesEditable           SubscriberFlags
	SupportsPreferredEmailAddressDelivery SubscriberFlags
	SupportsEachMemberDelivery            SubscriberFlags
	SupportsNoDelivery                    SubscriberFlags
	IsUser                                SubscriberFlags
	IsGroup                               SubscriberFlags
	IsTeam                                SubscriberFlags
}

var SubscriberFlagsValues = subscriberFlagsValuesType{
	None: "none",
	// Subscriber's delivery preferences could be updated
	DeliveryPreferencesEditable: "deliveryPreferencesEditable",
	// Subscriber's delivery preferences supports email delivery
	SupportsPreferredEmailAddressDelivery: "supportsPreferredEmailAddressDelivery",
	// Subscriber's delivery preferences supports individual members delivery(group expansion)
	SupportsEachMemberDelivery: "supportsEachMemberDelivery",
	// Subscriber's delivery preferences supports no delivery
	SupportsNoDelivery: "supportsNoDelivery",
	// Subscriber is a user
	IsUser: "isUser",
	// Subscriber is a group
	IsGroup: "isGroup",
	// Subscriber is a team
	IsTeam: "isTeam",
}

// Admin-managed settings for a group subscription.
type SubscriptionAdminSettings struct {
	// If true, members of the group subscribed to the associated subscription cannot opt (choose not to get notified)
	BlockUserOptOut *bool `json:"blockUserOptOut,omitempty"`
}

// Contains all the diagnostics settings for a subscription.
type SubscriptionDiagnostics struct {
	// Diagnostics settings for retaining delivery results.  Used for Service Hooks subscriptions.
	DeliveryResults *SubscriptionTracing `json:"deliveryResults,omitempty"`
	// Diagnostics settings for troubleshooting notification delivery.
	DeliveryTracing *SubscriptionTracing `json:"deliveryTracing,omitempty"`
	// Diagnostics settings for troubleshooting event matching.
	EvaluationTracing *SubscriptionTracing `json:"evaluationTracing,omitempty"`
}

type SubscriptionEvaluation struct {
	Clauses *[]SubscriptionEvaluationClause `json:"clauses,omitempty"`
	User    *DiagnosticIdentity             `json:"user,omitempty"`
}

type SubscriptionEvaluationClause struct {
	Clause *string `json:"clause,omitempty"`
	Order  *int    `json:"order,omitempty"`
	Result *bool   `json:"result,omitempty"`
}

// Encapsulates the properties of a SubscriptionEvaluationRequest. It defines the subscription to be evaluated and time interval for events used in evaluation.
type SubscriptionEvaluationRequest struct {
	// The min created date for the events used for matching in UTC. Use all events created since this date
	MinEventsCreatedDate *azuredevops.Time `json:"minEventsCreatedDate,omitempty"`
	// User or group that will receive notifications for events matching the subscription's filter criteria. If not specified, defaults to the calling user.
	SubscriptionCreateParameters *NotificationSubscriptionCreateParameters `json:"subscriptionCreateParameters,omitempty"`
}

// Encapsulates the subscription evaluation results. It defines the Date Interval that was used, number of events evaluated and events and notifications results
type SubscriptionEvaluationResult struct {
	// Subscription evaluation job status
	EvaluationJobStatus *EvaluationOperationStatus `json:"evaluationJobStatus,omitempty"`
	// Subscription evaluation events results.
	Events *EventsEvaluationResult `json:"events,omitempty"`
	// The requestId which is the subscription evaluation jobId
	Id *uuid.UUID `json:"id,omitempty"`
	// Subscription evaluation  notification results.
	Notifications *NotificationsEvaluationResult `json:"notifications,omitempty"`
}

// Encapsulates the subscription evaluation settings needed for the UI
type SubscriptionEvaluationSettings struct {
	// Indicates whether subscription evaluation before saving is enabled or not
	Enabled *bool `json:"enabled,omitempty"`
	// Time interval to check on subscription evaluation job in seconds
	Interval *int `json:"interval,omitempty"`
	// Threshold on the number of notifications for considering a subscription too noisy
	Threshold *int `json:"threshold,omitempty"`
	// Time out for the subscription evaluation check in seconds
	TimeOut *int `json:"timeOut,omitempty"`
}

type SubscriptionFieldType string

type subscriptionFieldTypeValuesType struct {
	String          SubscriptionFieldType
	Integer         SubscriptionFieldType
	DateTime        SubscriptionFieldType
	PlainText       SubscriptionFieldType
	Html            SubscriptionFieldType
	TreePath        SubscriptionFieldType
	History         SubscriptionFieldType
	Double          SubscriptionFieldType
	Guid            SubscriptionFieldType
	Boolean         SubscriptionFieldType
	Identity        SubscriptionFieldType
	PicklistInteger SubscriptionFieldType
	PicklistString  SubscriptionFieldType
	PicklistDouble  SubscriptionFieldType
	TeamProject     SubscriptionFieldType
}

var SubscriptionFieldTypeValues = subscriptionFieldTypeValuesType{
	String:          "string",
	Integer:         "integer",
	DateTime:        "dateTime",
	PlainText:       "plainText",
	Html:            "html",
	TreePath:        "treePath",
	History:         "history",
	Double:          "double",
	Guid:            "guid",
	Boolean:         "boolean",
	Identity:        "identity",
	PicklistInteger: "picklistInteger",
	PicklistString:  "picklistString",
	PicklistDouble:  "picklistDouble",
	TeamProject:     "teamProject",
}

// [Flags] Read-only indicators that further describe the subscription.
type SubscriptionFlags string

type subscriptionFlagsValuesType struct {
	None                    SubscriptionFlags
	GroupSubscription       SubscriptionFlags
	ContributedSubscription SubscriptionFlags
	CanOptOut               SubscriptionFlags
	TeamSubscription        SubscriptionFlags
	OneActorMatches         SubscriptionFlags
}

var SubscriptionFlagsValues = subscriptionFlagsValuesType{
	// None
	None: "none",
	// Subscription's subscriber is a group, not a user
	GroupSubscription: "groupSubscription",
	// Subscription is contributed and not persisted. This means certain fields of the subscription, like Filter, are read-only.
	ContributedSubscription: "contributedSubscription",
	// A user that is member of the subscription's subscriber group can opt in/out of the subscription.
	CanOptOut: "canOptOut",
	// If the subscriber is a group, is it a team.
	TeamSubscription: "teamSubscription",
	// For role based subscriptions, there is an expectation that there will always be at least one actor that matches
	OneActorMatches: "oneActorMatches",
}

type SubscriptionChannelWithAddress struct {
	Address          *string `json:"address,omitempty"`
	Type             *string `json:"type,omitempty"`
	UseCustomAddress *bool   `json:"useCustomAddress,omitempty"`
}

// Encapsulates the properties needed to manage subscriptions, opt in and out of subscriptions.
type SubscriptionManagement struct {
	ServiceInstanceType *uuid.UUID `json:"serviceInstanceType,omitempty"`
	Url                 *string    `json:"url,omitempty"`
}

// [Flags] The permissions that a user has to a certain subscription
type SubscriptionPermissions string

type subscriptionPermissionsValuesType struct {
	None   SubscriptionPermissions
	View   SubscriptionPermissions
	Edit   SubscriptionPermissions
	Delete SubscriptionPermissions
}

var SubscriptionPermissionsValues = subscriptionPermissionsValuesType{
	// None
	None: "none",
	// full view of description, filters, etc. Not limited.
	View: "view",
	// update subscription
	Edit: "edit",
	// delete subscription
	Delete: "delete",
}

// Notification subscriptions query input.
type SubscriptionQuery struct {
	// One or more conditions to query on. If more than 2 conditions are specified, the combined results of each condition is returned (i.e. conditions are logically OR'ed).
	Conditions *[]SubscriptionQueryCondition `json:"conditions,omitempty"`
	// Flags the refine the types of subscriptions that will be returned from the query.
	QueryFlags *SubscriptionQueryFlags `json:"queryFlags,omitempty"`
}

// Conditions a subscription must match to qualify for the query result set. Not all fields are required. A subscription must match all conditions specified in order to qualify for the result set.
type SubscriptionQueryCondition struct {
	// Filter conditions that matching subscriptions must have. Typically only the filter's type and event type are used for matching.
	Filter *ISubscriptionFilter `json:"filter,omitempty"`
	// Flags to specify the type subscriptions to query for.
	Flags *SubscriptionFlags `json:"flags,omitempty"`
	// Scope that matching subscriptions must have.
	Scope *string `json:"scope,omitempty"`
	// ID of the subscriber (user or group) that matching subscriptions must be subscribed to.
	SubscriberId *uuid.UUID `json:"subscriberId,omitempty"`
	// ID of the subscription to query for.
	SubscriptionId *string `json:"subscriptionId,omitempty"`
}

// [Flags] Flags that influence the result set of a subscription query.
type SubscriptionQueryFlags string

type subscriptionQueryFlagsValuesType struct {
	None                         SubscriptionQueryFlags
	IncludeInvalidSubscriptions  SubscriptionQueryFlags
	IncludeDeletedSubscriptions  SubscriptionQueryFlags
	IncludeFilterDetails         SubscriptionQueryFlags
	AlwaysReturnBasicInformation SubscriptionQueryFlags
	IncludeSystemSubscriptions   SubscriptionQueryFlags
}

var SubscriptionQueryFlagsValues = subscriptionQueryFlagsValuesType{
	None: "none",
	// Include subscriptions with invalid subscribers.
	IncludeInvalidSubscriptions: "includeInvalidSubscriptions",
	// Include subscriptions marked for deletion.
	IncludeDeletedSubscriptions: "includeDeletedSubscriptions",
	// Include the full filter details with each subscription.
	IncludeFilterDetails: "includeFilterDetails",
	// For a subscription the caller does not have permission to view, return basic (non-confidential) information.
	AlwaysReturnBasicInformation: "alwaysReturnBasicInformation",
	// Include system subscriptions.
	IncludeSystemSubscriptions: "includeSystemSubscriptions",
}

// A resource, typically an account or project, in which events are published from.
type SubscriptionScope struct {
	// Required: This is the identity of the scope for the type.
	Id *uuid.UUID `json:"id,omitempty"`
	// Optional: The display name of the scope
	Name *string `json:"name,omitempty"`
	// Required: The event specific type of a scope.
	Type *string `json:"type,omitempty"`
}

// Subscription status values. A value greater than or equal to zero indicates the subscription is enabled. A negative value indicates the subscription is disabled.
type SubscriptionStatus string

type subscriptionStatusValuesType struct {
	JailedByNotificationsVolume      SubscriptionStatus
	PendingDeletion                  SubscriptionStatus
	DisabledArgumentException        SubscriptionStatus
	DisabledProjectInvalid           SubscriptionStatus
	DisabledMissingPermissions       SubscriptionStatus
	DisabledFromProbation            SubscriptionStatus
	DisabledInactiveIdentity         SubscriptionStatus
	DisabledMessageQueueNotSupported SubscriptionStatus
	DisabledMissingIdentity          SubscriptionStatus
	DisabledInvalidRoleExpression    SubscriptionStatus
	DisabledInvalidPathClause        SubscriptionStatus
	DisabledAsDuplicateOfDefault     SubscriptionStatus
	DisabledByAdmin                  SubscriptionStatus
	Disabled                         SubscriptionStatus
	Enabled                          SubscriptionStatus
	EnabledOnProbation               SubscriptionStatus
}

var SubscriptionStatusValues = subscriptionStatusValuesType{
	// Subscription is disabled because it generated a high volume of notifications.
	JailedByNotificationsVolume: "jailedByNotificationsVolume",
	// Subscription is disabled and will be deleted.
	PendingDeletion: "pendingDeletion",
	// Subscription is disabled because of an Argument Exception while processing the subscription
	DisabledArgumentException: "disabledArgumentException",
	// Subscription is disabled because the project is invalid
	DisabledProjectInvalid: "disabledProjectInvalid",
	// Subscription is disabled because the identity does not have the appropriate permissions
	DisabledMissingPermissions: "disabledMissingPermissions",
	// Subscription is disabled service due to failures.
	DisabledFromProbation: "disabledFromProbation",
	// Subscription is disabled because the identity is no longer active
	DisabledInactiveIdentity: "disabledInactiveIdentity",
	// Subscription is disabled because message queue is not supported.
	DisabledMessageQueueNotSupported: "disabledMessageQueueNotSupported",
	// Subscription is disabled because its subscriber is unknown.
	DisabledMissingIdentity: "disabledMissingIdentity",
	// Subscription is disabled because it has an invalid role expression.
	DisabledInvalidRoleExpression: "disabledInvalidRoleExpression",
	// Subscription is disabled because it has an invalid filter expression.
	DisabledInvalidPathClause: "disabledInvalidPathClause",
	// Subscription is disabled because it is a duplicate of a default subscription.
	DisabledAsDuplicateOfDefault: "disabledAsDuplicateOfDefault",
	// Subscription is disabled by an administrator, not the subscription's subscriber.
	DisabledByAdmin: "disabledByAdmin",
	// Subscription is disabled, typically by the owner of the subscription, and will not produce any notifications.
	Disabled: "disabled",
	// Subscription is active.
	Enabled: "enabled",
	// Subscription is active, but is on probation due to failed deliveries or other issues with the subscription.
	EnabledOnProbation: "enabledOnProbation",
}

// [Flags] Set of flags used to determine which set of templates is retrieved when querying for subscription templates
type SubscriptionTemplateQueryFlags string

type subscriptionTemplateQueryFlagsValuesType struct {
	None                        SubscriptionTemplateQueryFlags
	IncludeUser                 SubscriptionTemplateQueryFlags
	IncludeGroup                SubscriptionTemplateQueryFlags
	IncludeUserAndGroup         SubscriptionTemplateQueryFlags
	IncludeEventTypeInformation SubscriptionTemplateQueryFlags
}

var SubscriptionTemplateQueryFlagsValues = subscriptionTemplateQueryFlagsValuesType{
	None: "none",
	// Include user templates
	IncludeUser: "includeUser",
	// Include group templates
	IncludeGroup: "includeGroup",
	// Include user and group templates
	IncludeUserAndGroup: "includeUserAndGroup",
	// Include the event type details like the fields and operators
	IncludeEventTypeInformation: "includeEventTypeInformation",
}

type SubscriptionTemplateType string

type subscriptionTemplateTypeValuesType struct {
	User SubscriptionTemplateType
	Team SubscriptionTemplateType
	Both SubscriptionTemplateType
	None SubscriptionTemplateType
}

var SubscriptionTemplateTypeValues = subscriptionTemplateTypeValuesType{
	User: "user",
	Team: "team",
	Both: "both",
	None: "none",
}

type SubscriptionTraceDiagnosticLog struct {
	// Identifier used for correlating to other diagnostics that may have been recorded elsewhere.
	ActivityId  *uuid.UUID        `json:"activityId,omitempty"`
	Description *string           `json:"description,omitempty"`
	EndTime     *azuredevops.Time `json:"endTime,omitempty"`
	Errors      *int              `json:"errors,omitempty"`
	// Unique instance identifier.
	Id         *uuid.UUID                          `json:"id,omitempty"`
	LogType    *string                             `json:"logType,omitempty"`
	Messages   *[]NotificationDiagnosticLogMessage `json:"messages,omitempty"`
	Properties *map[string]string                  `json:"properties,omitempty"`
	// This identifier depends on the logType.  For notification jobs, this will be the job Id. For subscription tracing, this will be a special root Guid with the subscription Id encoded.
	Source    *uuid.UUID        `json:"source,omitempty"`
	StartTime *azuredevops.Time `json:"startTime,omitempty"`
	Warnings  *int              `json:"warnings,omitempty"`
	// Indicates the job Id that processed or delivered this subscription
	JobId *uuid.UUID `json:"jobId,omitempty"`
	// Indicates unique instance identifier for the job that processed or delivered this subscription
	JobInstanceId  *uuid.UUID `json:"jobInstanceId,omitempty"`
	SubscriptionId *string    `json:"subscriptionId,omitempty"`
}

type SubscriptionTraceEventProcessingLog struct {
	// Identifier used for correlating to other diagnostics that may have been recorded elsewhere.
	ActivityId  *uuid.UUID        `json:"activityId,omitempty"`
	Description *string           `json:"description,omitempty"`
	EndTime     *azuredevops.Time `json:"endTime,omitempty"`
	Errors      *int              `json:"errors,omitempty"`
	// Unique instance identifier.
	Id         *uuid.UUID                          `json:"id,omitempty"`
	LogType    *string                             `json:"logType,omitempty"`
	Messages   *[]NotificationDiagnosticLogMessage `json:"messages,omitempty"`
	Properties *map[string]string                  `json:"properties,omitempty"`
	// This identifier depends on the logType.  For notification jobs, this will be the job Id. For subscription tracing, this will be a special root Guid with the subscription Id encoded.
	Source    *uuid.UUID        `json:"source,omitempty"`
	StartTime *azuredevops.Time `json:"startTime,omitempty"`
	Warnings  *int              `json:"warnings,omitempty"`
	// Indicates the job Id that processed or delivered this subscription
	JobId *uuid.UUID `json:"jobId,omitempty"`
	// Indicates unique instance identifier for the job that processed or delivered this subscription
	JobInstanceId        *uuid.UUID            `json:"jobInstanceId,omitempty"`
	SubscriptionId       *string               `json:"subscriptionId,omitempty"`
	EvaluationIdentities *ProcessingIdentities `json:"evaluationIdentities,omitempty"`
	Channel              *string               `json:"channel,omitempty"`
	// Which members opted out from receiving notifications from this subscription
	OptedOut        *[]DiagnosticIdentity   `json:"optedOut,omitempty"`
	ProcessedEvents *map[int]ProcessedEvent `json:"processedEvents,omitempty"`
}

type SubscriptionTraceNotificationDeliveryLog struct {
	// Identifier used for correlating to other diagnostics that may have been recorded elsewhere.
	ActivityId  *uuid.UUID        `json:"activityId,omitempty"`
	Description *string           `json:"description,omitempty"`
	EndTime     *azuredevops.Time `json:"endTime,omitempty"`
	Errors      *int              `json:"errors,omitempty"`
	// Unique instance identifier.
	Id         *uuid.UUID                          `json:"id,omitempty"`
	LogType    *string                             `json:"logType,omitempty"`
	Messages   *[]NotificationDiagnosticLogMessage `json:"messages,omitempty"`
	Properties *map[string]string                  `json:"properties,omitempty"`
	// This identifier depends on the logType.  For notification jobs, this will be the job Id. For subscription tracing, this will be a special root Guid with the subscription Id encoded.
	Source    *uuid.UUID        `json:"source,omitempty"`
	StartTime *azuredevops.Time `json:"startTime,omitempty"`
	Warnings  *int              `json:"warnings,omitempty"`
	// Indicates the job Id that processed or delivered this subscription
	JobId *uuid.UUID `json:"jobId,omitempty"`
	// Indicates unique instance identifier for the job that processed or delivered this subscription
	JobInstanceId  *uuid.UUID                `json:"jobInstanceId,omitempty"`
	SubscriptionId *string                   `json:"subscriptionId,omitempty"`
	Notifications  *[]DiagnosticNotification `json:"notifications,omitempty"`
}

// Data controlling a single diagnostic setting for a subscription.
type SubscriptionTracing struct {
	// Indicates whether the diagnostic tracing is enabled or not.
	Enabled *bool `json:"enabled,omitempty"`
	// Trace until the specified end date.
	EndDate *azuredevops.Time `json:"endDate,omitempty"`
	// The maximum number of result details to trace.
	MaxTracedEntries *int `json:"maxTracedEntries,omitempty"`
	// The date and time tracing started.
	StartDate *azuredevops.Time `json:"startDate,omitempty"`
	// Trace until remaining count reaches 0.
	TracedEntries *int `json:"tracedEntries,omitempty"`
}

// User-managed settings for a group subscription.
type SubscriptionUserSettings struct {
	// Indicates whether the user will receive notifications for the associated group subscription.
	OptedOut *bool `json:"optedOut,omitempty"`
}

type UnsupportedFilter struct {
	EventType *string `json:"eventType,omitempty"`
	Type      *string `json:"type,omitempty"`
}

type UnsupportedSubscriptionChannel struct {
	Type *string `json:"type,omitempty"`
}

// Parameters to update diagnostics settings for a subscription.
type UpdateSubscripitonDiagnosticsParameters struct {
	// Diagnostics settings for retaining delivery results.  Used for Service Hooks subscriptions.
	DeliveryResults *UpdateSubscripitonTracingParameters `json:"deliveryResults,omitempty"`
	// Diagnostics settings for troubleshooting notification delivery.
	DeliveryTracing *UpdateSubscripitonTracingParameters `json:"deliveryTracing,omitempty"`
	// Diagnostics settings for troubleshooting event matching.
	EvaluationTracing *UpdateSubscripitonTracingParameters `json:"evaluationTracing,omitempty"`
}

// Parameters to update a specific diagnostic setting.
type UpdateSubscripitonTracingParameters struct {
	// Indicates whether to enable to disable the diagnostic tracing.
	Enabled *bool `json:"enabled,omitempty"`
}

type UserSubscriptionChannel struct {
	Address          *string `json:"address,omitempty"`
	UseCustomAddress *bool   `json:"useCustomAddress,omitempty"`
	Type             *string `json:"type,omitempty"`
}

type UserSystemSubscriptionChannel struct {
	Address          *string `json:"address,omitempty"`
	UseCustomAddress *bool   `json:"useCustomAddress,omitempty"`
	Type             *string `json:"type,omitempty"`
}

// Encapsulates the properties of a field value definition. It has the information needed to retrieve the list of possible values for a certain field and how to handle that field values in the UI. This information includes what type of object this value represents, which property to use for UI display and which property to use for saving the subscription
type ValueDefinition struct {
	// Gets or sets the data source.
	DataSource *[]forminput.InputValue `json:"dataSource,omitempty"`
	// Gets or sets the rest end point.
	EndPoint *string `json:"endPoint,omitempty"`
	// Gets or sets the result template.
	ResultTemplate *string `json:"resultTemplate,omitempty"`
}

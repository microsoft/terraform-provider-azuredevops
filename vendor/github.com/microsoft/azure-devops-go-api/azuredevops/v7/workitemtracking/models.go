// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package workitemtracking

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
)

type AccountMyWorkResult struct {
	// True, when length of WorkItemDetails is same as the limit
	QuerySizeLimitExceeded *bool `json:"querySizeLimitExceeded,omitempty"`
	// WorkItem Details
	WorkItemDetails *[]AccountWorkWorkItemModel `json:"workItemDetails,omitempty"`
}

// Represents Work Item Recent Activity
type AccountRecentActivityWorkItemModel struct {
	// Date of the last Activity by the user
	ActivityDate *azuredevops.Time `json:"activityDate,omitempty"`
	// Type of the activity
	ActivityType *WorkItemRecentActivityType `json:"activityType,omitempty"`
	// Last changed date of the work item
	ChangedDate *azuredevops.Time `json:"changedDate,omitempty"`
	// Work Item Id
	Id *int `json:"id,omitempty"`
	// TeamFoundationId of the user this activity belongs to
	IdentityId *uuid.UUID `json:"identityId,omitempty"`
	// State of the work item
	State *string `json:"state,omitempty"`
	// Team project the work item belongs to
	TeamProject *string `json:"teamProject,omitempty"`
	// Title of the work item
	Title *string `json:"title,omitempty"`
	// Type of Work Item
	WorkItemType *string `json:"workItemType,omitempty"`
	// Assigned To
	AssignedTo *string `json:"assignedTo,omitempty"`
}

// Represents Work Item Recent Activity
type AccountRecentActivityWorkItemModel2 struct {
	// Date of the last Activity by the user
	ActivityDate *azuredevops.Time `json:"activityDate,omitempty"`
	// Type of the activity
	ActivityType *WorkItemRecentActivityType `json:"activityType,omitempty"`
	// Last changed date of the work item
	ChangedDate *azuredevops.Time `json:"changedDate,omitempty"`
	// Work Item Id
	Id *int `json:"id,omitempty"`
	// TeamFoundationId of the user this activity belongs to
	IdentityId *uuid.UUID `json:"identityId,omitempty"`
	// State of the work item
	State *string `json:"state,omitempty"`
	// Team project the work item belongs to
	TeamProject *string `json:"teamProject,omitempty"`
	// Title of the work item
	Title *string `json:"title,omitempty"`
	// Type of Work Item
	WorkItemType *string `json:"workItemType,omitempty"`
	// Assigned To
	AssignedTo *webapi.IdentityRef `json:"assignedTo,omitempty"`
}

// Represents Work Item Recent Activity
type AccountRecentActivityWorkItemModelBase struct {
	// Date of the last Activity by the user
	ActivityDate *azuredevops.Time `json:"activityDate,omitempty"`
	// Type of the activity
	ActivityType *WorkItemRecentActivityType `json:"activityType,omitempty"`
	// Last changed date of the work item
	ChangedDate *azuredevops.Time `json:"changedDate,omitempty"`
	// Work Item Id
	Id *int `json:"id,omitempty"`
	// TeamFoundationId of the user this activity belongs to
	IdentityId *uuid.UUID `json:"identityId,omitempty"`
	// State of the work item
	State *string `json:"state,omitempty"`
	// Team project the work item belongs to
	TeamProject *string `json:"teamProject,omitempty"`
	// Title of the work item
	Title *string `json:"title,omitempty"`
	// Type of Work Item
	WorkItemType *string `json:"workItemType,omitempty"`
}

// Represents Recent Mention Work Item
type AccountRecentMentionWorkItemModel struct {
	// Assigned To
	AssignedTo *string `json:"assignedTo,omitempty"`
	// Work Item Id
	Id *int `json:"id,omitempty"`
	// Latest date that the user were mentioned
	MentionedDateField *azuredevops.Time `json:"mentionedDateField,omitempty"`
	// State of the work item
	State *string `json:"state,omitempty"`
	// Team project the work item belongs to
	TeamProject *string `json:"teamProject,omitempty"`
	// Title of the work item
	Title *string `json:"title,omitempty"`
	// Type of Work Item
	WorkItemType *string `json:"workItemType,omitempty"`
}

type AccountWorkWorkItemModel struct {
	AssignedTo   *string           `json:"assignedTo,omitempty"`
	ChangedDate  *azuredevops.Time `json:"changedDate,omitempty"`
	Id           *int              `json:"id,omitempty"`
	State        *string           `json:"state,omitempty"`
	TeamProject  *string           `json:"teamProject,omitempty"`
	Title        *string           `json:"title,omitempty"`
	WorkItemType *string           `json:"workItemType,omitempty"`
}

// Contains criteria for querying work items based on artifact URI.
type ArtifactUriQuery struct {
	// List of artifact URIs to use for querying work items.
	ArtifactUris *[]string `json:"artifactUris,omitempty"`
}

// Defines result of artifact URI query on work items. Contains mapping of work item IDs to artifact URI.
type ArtifactUriQueryResult struct {
	// A Dictionary that maps a list of work item references to the given list of artifact URI.
	ArtifactUrisQueryResult *map[string][]WorkItemReference `json:"artifactUrisQueryResult,omitempty"`
}

type AttachmentReference struct {
	Id  *uuid.UUID `json:"id,omitempty"`
	Url *string    `json:"url,omitempty"`
}

// Flag to control error policy in a batch classification nodes get request.
type ClassificationNodesErrorPolicy string

type classificationNodesErrorPolicyValuesType struct {
	Fail ClassificationNodesErrorPolicy
	Omit ClassificationNodesErrorPolicy
}

var ClassificationNodesErrorPolicyValues = classificationNodesErrorPolicyValuesType{
	Fail: "fail",
	Omit: "omit",
}

// Comment on a Work Item.
type Comment struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// IdentityRef of the creator of the comment.
	CreatedBy *webapi.IdentityRef `json:"createdBy,omitempty"`
	// The creation date of the comment.
	CreatedDate *azuredevops.Time `json:"createdDate,omitempty"`
	// Effective Date/time value for adding the comment. Can be optionally different from CreatedDate.
	CreatedOnBehalfDate *azuredevops.Time `json:"createdOnBehalfDate,omitempty"`
	// Identity on whose behalf this comment has been added. Can be optionally different from CreatedBy.
	CreatedOnBehalfOf *webapi.IdentityRef `json:"createdOnBehalfOf,omitempty"`
	// Represents the possible types for the comment format.
	Format *CommentFormat `json:"format,omitempty"`
	// The id assigned to the comment.
	Id *int `json:"id,omitempty"`
	// Indicates if the comment has been deleted.
	IsDeleted *bool `json:"isDeleted,omitempty"`
	// The mentions of the comment.
	Mentions *[]CommentMention `json:"mentions,omitempty"`
	// IdentityRef of the user who last modified the comment.
	ModifiedBy *webapi.IdentityRef `json:"modifiedBy,omitempty"`
	// The last modification date of the comment.
	ModifiedDate *azuredevops.Time `json:"modifiedDate,omitempty"`
	// The reactions of the comment.
	Reactions *[]CommentReaction `json:"reactions,omitempty"`
	// The text of the comment in HTML format.
	RenderedText *string `json:"renderedText,omitempty"`
	// The text of the comment.
	Text *string `json:"text,omitempty"`
	// The current version of the comment.
	Version *int `json:"version,omitempty"`
	// The id of the work item this comment belongs to.
	WorkItemId *int `json:"workItemId,omitempty"`
}

// Represents a request to create a work item comment.
type CommentCreate struct {
	// The text of the comment.
	Text *string `json:"text,omitempty"`
}

// [Flags] Specifies the additional data retrieval options for work item comments.
type CommentExpandOptions string

type commentExpandOptionsValuesType struct {
	None             CommentExpandOptions
	Reactions        CommentExpandOptions
	RenderedText     CommentExpandOptions
	RenderedTextOnly CommentExpandOptions
	All              CommentExpandOptions
}

var CommentExpandOptionsValues = commentExpandOptionsValuesType{
	None: "none",
	// Include comment reactions.
	Reactions: "reactions",
	// Include the rendered text (html) in addition to MD text.
	RenderedText: "renderedText",
	// If specified, then ONLY rendered text (html) will be returned, w/o markdown. Supposed to be used internally from data provides for optimization purposes.
	RenderedTextOnly: "renderedTextOnly",
	All:              "all",
}

// Represents the possible types for the comment format. Should be in sync with WorkItemCommentFormat.cs
type CommentFormat string

type commentFormatValuesType struct {
	Markdown CommentFormat
	Html     CommentFormat
}

var CommentFormatValues = commentFormatValuesType{
	Markdown: "markdown",
	Html:     "html",
}

// Represents a list of work item comments.
type CommentList struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// List of comments in the current batch.
	Comments *[]Comment `json:"comments,omitempty"`
	// A string token that can be used to retrieving next page of comments if available. Otherwise null.
	ContinuationToken *string `json:"continuationToken,omitempty"`
	// The count of comments in the current batch.
	Count *int `json:"count,omitempty"`
	// Uri to the next page of comments if it is available. Otherwise null.
	NextPage *string `json:"nextPage,omitempty"`
	// Total count of comments on a work item.
	TotalCount *int `json:"totalCount,omitempty"`
}

type CommentMention struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// The artifact portion of the parsed text. (i.e. the work item's id)
	ArtifactId *string `json:"artifactId,omitempty"`
	// The type the parser assigned to the mention. (i.e. person, work item, etc)
	ArtifactType *string `json:"artifactType,omitempty"`
	// The comment id of the mention.
	CommentId *int `json:"commentId,omitempty"`
	// The resolved target of the mention. An example of this could be a user's tfid
	TargetId *string `json:"targetId,omitempty"`
}

// Contains information about work item comment reaction for a particular reaction type.
type CommentReaction struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// The id of the comment this reaction belongs to.
	CommentId *int `json:"commentId,omitempty"`
	// Total number of reactions for the CommentReactionType.
	Count *int `json:"count,omitempty"`
	// Flag to indicate if the current user has engaged on this particular EngagementType (e.g. if they liked the associated comment).
	IsCurrentUserEngaged *bool `json:"isCurrentUserEngaged,omitempty"`
	// Type of the reaction.
	Type *CommentReactionType `json:"type,omitempty"`
}

// Represents different reaction types for a work item comment.
type CommentReactionType string

type commentReactionTypeValuesType struct {
	Like     CommentReactionType
	Dislike  CommentReactionType
	Heart    CommentReactionType
	Hooray   CommentReactionType
	Smile    CommentReactionType
	Confused CommentReactionType
}

var CommentReactionTypeValues = commentReactionTypeValuesType{
	Like:     "like",
	Dislike:  "dislike",
	Heart:    "heart",
	Hooray:   "hooray",
	Smile:    "smile",
	Confused: "confused",
}

type CommentSortOrder string

type commentSortOrderValuesType struct {
	Asc  CommentSortOrder
	Desc CommentSortOrder
}

var CommentSortOrderValues = commentSortOrderValuesType{
	// The results will be sorted in Ascending order.
	Asc: "asc",
	// The results will be sorted in Descending order.
	Desc: "desc",
}

// Represents a request to update a work item comment.
type CommentUpdate struct {
	// The updated text of the comment.
	Text *string `json:"text,omitempty"`
}

// Represents a specific version of a comment on a work item.
type CommentVersion struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// IdentityRef of the creator of the comment.
	CreatedBy *webapi.IdentityRef `json:"createdBy,omitempty"`
	// The creation date of the comment.
	CreatedDate *azuredevops.Time `json:"createdDate,omitempty"`
	// Effective Date/time value for adding the comment. Can be optionally different from CreatedDate.
	CreatedOnBehalfDate *azuredevops.Time `json:"createdOnBehalfDate,omitempty"`
	// Identity on whose behalf this comment has been added. Can be optionally different from CreatedBy.
	CreatedOnBehalfOf *webapi.IdentityRef `json:"createdOnBehalfOf,omitempty"`
	// The id assigned to the comment.
	Id *int `json:"id,omitempty"`
	// Indicates if the comment has been deleted at this version.
	IsDeleted *bool `json:"isDeleted,omitempty"`
	// IdentityRef of the user who modified the comment at this version.
	ModifiedBy *webapi.IdentityRef `json:"modifiedBy,omitempty"`
	// The modification date of the comment for this version.
	ModifiedDate *azuredevops.Time `json:"modifiedDate,omitempty"`
	// The rendered content of the comment at this version.
	RenderedText *string `json:"renderedText,omitempty"`
	// The text of the comment at this version.
	Text *string `json:"text,omitempty"`
	// The version number.
	Version *int `json:"version,omitempty"`
}

type EmailRecipients struct {
	// Plaintext email addresses.
	EmailAddresses *[]string `json:"emailAddresses,omitempty"`
	// TfIds
	TfIds *[]uuid.UUID `json:"tfIds,omitempty"`
	// Unresolved entity ids
	UnresolvedEntityIds *[]uuid.UUID `json:"unresolvedEntityIds,omitempty"`
}

type ExternalDeployment struct {
	ArtifactId         *uuid.UUID           `json:"artifactId,omitempty"`
	CreatedBy          *uuid.UUID           `json:"createdBy,omitempty"`
	Description        *string              `json:"description,omitempty"`
	DisplayName        *string              `json:"displayName,omitempty"`
	Environment        *ExternalEnvironment `json:"environment,omitempty"`
	Group              *string              `json:"group,omitempty"`
	Pipeline           *ExternalPipeline    `json:"pipeline,omitempty"`
	RelatedWorkItemIds *[]int               `json:"relatedWorkItemIds,omitempty"`
	RunId              *int                 `json:"runId,omitempty"`
	SequenceNumber     *int                 `json:"sequenceNumber,omitempty"`
	Status             *string              `json:"status,omitempty"`
	StatusDate         *azuredevops.Time    `json:"statusDate,omitempty"`
	Url                *string              `json:"url,omitempty"`
}

type ExternalEnvironment struct {
	DisplayName *string `json:"displayName,omitempty"`
	Id          *int    `json:"id,omitempty"`
	Type        *string `json:"type,omitempty"`
}

type ExternalPipeline struct {
	DisplayName *string `json:"displayName,omitempty"`
	Id          *int    `json:"id,omitempty"`
	Url         *string `json:"url,omitempty"`
}

// Describes a list of dependent fields for a rule.
type FieldDependentRule struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// The dependent fields.
	DependentFields *[]WorkItemFieldReference `json:"dependentFields,omitempty"`
}

// Enum for field types.
type FieldType string

type fieldTypeValuesType struct {
	String          FieldType
	Integer         FieldType
	DateTime        FieldType
	PlainText       FieldType
	Html            FieldType
	TreePath        FieldType
	History         FieldType
	Double          FieldType
	Guid            FieldType
	Boolean         FieldType
	Identity        FieldType
	PicklistString  FieldType
	PicklistInteger FieldType
	PicklistDouble  FieldType
}

var FieldTypeValues = fieldTypeValuesType{
	// String field type.
	String: "string",
	// Integer field type.
	Integer: "integer",
	// Datetime field type.
	DateTime: "dateTime",
	// Plain text field type.
	PlainText: "plainText",
	// HTML (Multiline) field type.
	Html: "html",
	// Treepath field type.
	TreePath: "treePath",
	// History field type.
	History: "history",
	// Double field type.
	Double: "double",
	// Guid field type.
	Guid: "guid",
	// Boolean field type.
	Boolean: "boolean",
	// Identity field type.
	Identity: "identity",
	// String picklist field type. When creating a string picklist field from REST API, use "String" FieldType.
	PicklistString: "picklistString",
	// Integer picklist field type. When creating a integer picklist field from REST API, use "Integer" FieldType.
	PicklistInteger: "picklistInteger",
	// Double picklist field type. When creating a double picklist field from REST API, use "Double" FieldType.
	PicklistDouble: "picklistDouble",
}

// Describes an update request for a work item field.
type FieldUpdate struct {
	// Indicates whether the user wants to restore the field.
	IsDeleted *bool `json:"isDeleted,omitempty"`
	// Indicates whether the user wants to lock the field.
	IsLocked *bool `json:"isLocked,omitempty"`
}

// Enum for field usages.
type FieldUsage string

type fieldUsageValuesType struct {
	None                  FieldUsage
	WorkItem              FieldUsage
	WorkItemLink          FieldUsage
	Tree                  FieldUsage
	WorkItemTypeExtension FieldUsage
}

var FieldUsageValues = fieldUsageValuesType{
	// Empty usage.
	None: "none",
	// Work item field usage.
	WorkItem: "workItem",
	// Work item link field usage.
	WorkItemLink: "workItemLink",
	// Treenode field usage.
	Tree: "tree",
	// Work Item Type Extension usage.
	WorkItemTypeExtension: "workItemTypeExtension",
}

// Flag to expand types of fields.
type GetFieldsExpand string

type getFieldsExpandValuesType struct {
	None            GetFieldsExpand
	ExtensionFields GetFieldsExpand
	IncludeDeleted  GetFieldsExpand
}

var GetFieldsExpandValues = getFieldsExpandValuesType{
	// Default behavior.
	None: "none",
	// Adds extension fields to the response.
	ExtensionFields: "extensionFields",
	// Includes fields that have been deleted.
	IncludeDeleted: "includeDeleted",
}

// Describes Github connection.
type GitHubConnectionModel struct {
	// Github connection authorization type (f. e. PAT, OAuth)
	AuthorizationType *string `json:"authorizationType,omitempty"`
	// Github connection created by
	CreatedBy *webapi.IdentityRef `json:"createdBy,omitempty"`
	// Github connection id
	Id *uuid.UUID `json:"id,omitempty"`
	// Whether current Github connection is valid or not
	IsConnectionValid *bool `json:"isConnectionValid,omitempty"`
	// Github connection name (should contain organization/user name)
	Name *string `json:"name,omitempty"`
}

// Describes Github connection's repo.
type GitHubConnectionRepoModel struct {
	// Error message
	ErrorMessage *string `json:"errorMessage,omitempty"`
	// Repository web url
	GitHubRepositoryUrl *string `json:"gitHubRepositoryUrl,omitempty"`
}

// Describes Github connection's repo bulk request
type GitHubConnectionReposBatchRequest struct {
	// Requested repos urls
	GitHubRepositoryUrls *[]GitHubConnectionRepoModel `json:"gitHubRepositoryUrls,omitempty"`
	// Operation type (f. e. add, remove)
	OperationType *string `json:"operationType,omitempty"`
}

// Describes a reference to an identity.
type IdentityReference struct {
	// This field contains zero or more interesting links about the graph subject. These links may be invoked to obtain additional relationships or more detailed information about this graph subject.
	Links interface{} `json:"_links,omitempty"`
	// The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the same graph subject across both Accounts and Organizations.
	Descriptor *string `json:"descriptor,omitempty"`
	// This is the non-unique display name of the graph subject. To change this field, you must alter its value in the source provider.
	DisplayName *string `json:"displayName,omitempty"`
	// This url is the full route to the source resource of this graph subject.
	Url *string `json:"url,omitempty"`
	// Deprecated - Can be retrieved by querying the Graph user referenced in the "self" entry of the IdentityRef "_links" dictionary
	DirectoryAlias *string `json:"directoryAlias,omitempty"`
	// Deprecated - Available in the "avatar" entry of the IdentityRef "_links" dictionary
	ImageUrl *string `json:"imageUrl,omitempty"`
	// Deprecated - Can be retrieved by querying the Graph membership state referenced in the "membershipState" entry of the GraphUser "_links" dictionary
	Inactive *bool `json:"inactive,omitempty"`
	// Deprecated - Can be inferred from the subject type of the descriptor (Descriptor.IsAadUserType/Descriptor.IsAadGroupType)
	IsAadIdentity *bool `json:"isAadIdentity,omitempty"`
	// Deprecated - Can be inferred from the subject type of the descriptor (Descriptor.IsGroupType)
	IsContainer       *bool `json:"isContainer,omitempty"`
	IsDeletedInOrigin *bool `json:"isDeletedInOrigin,omitempty"`
	// Deprecated - not in use in most preexisting implementations of ToIdentityRef
	ProfileUrl *string `json:"profileUrl,omitempty"`
	// Deprecated - use Domain+PrincipalName instead
	UniqueName *string    `json:"uniqueName,omitempty"`
	Id         *uuid.UUID `json:"id,omitempty"`
	// Legacy back-compat property. This has been the WIT specific value from Constants. Will be hidden (but exists) on the client unless they are targeting the newest version
	Name *string `json:"name,omitempty"`
}

// Link description.
type Link struct {
	// Collection of link attributes.
	Attributes *map[string]interface{} `json:"attributes,omitempty"`
	// Relation type.
	Rel *string `json:"rel,omitempty"`
	// Link url.
	Url *string `json:"url,omitempty"`
}

// The link query mode which determines the behavior of the query.
type LinkQueryMode string

type linkQueryModeValuesType struct {
	WorkItems                    LinkQueryMode
	LinksOneHopMustContain       LinkQueryMode
	LinksOneHopMayContain        LinkQueryMode
	LinksOneHopDoesNotContain    LinkQueryMode
	LinksRecursiveMustContain    LinkQueryMode
	LinksRecursiveMayContain     LinkQueryMode
	LinksRecursiveDoesNotContain LinkQueryMode
}

var LinkQueryModeValues = linkQueryModeValuesType{
	// Returns flat list of work items.
	WorkItems: "workItems",
	// Returns work items where the source, target, and link criteria are all satisfied.
	LinksOneHopMustContain: "linksOneHopMustContain",
	// Returns work items that satisfy the source and link criteria, even if no linked work item satisfies the target criteria.
	LinksOneHopMayContain: "linksOneHopMayContain",
	// Returns work items that satisfy the source, only if no linked work item satisfies the link and target criteria.
	LinksOneHopDoesNotContain: "linksOneHopDoesNotContain",
	LinksRecursiveMustContain: "linksRecursiveMustContain",
	// Returns work items a hierarchy of work items that by default satisfy the source
	LinksRecursiveMayContain:     "linksRecursiveMayContain",
	LinksRecursiveDoesNotContain: "linksRecursiveDoesNotContain",
}

type LogicalOperation string

type logicalOperationValuesType struct {
	None LogicalOperation
	And  LogicalOperation
	Or   LogicalOperation
}

var LogicalOperationValues = logicalOperationValuesType{
	None: "none",
	And:  "and",
	Or:   "or",
}

type MailMessage struct {
	// The mail body in HTML format.
	Body *string `json:"body,omitempty"`
	// CC recipients.
	Cc *EmailRecipients `json:"cc,omitempty"`
	// The in-reply-to header value
	InReplyTo *string `json:"inReplyTo,omitempty"`
	// The Message Id value
	MessageId *string `json:"messageId,omitempty"`
	// Reply To recipients.
	ReplyTo *EmailRecipients `json:"replyTo,omitempty"`
	// The mail subject.
	Subject *string `json:"subject,omitempty"`
	// To recipients
	To *EmailRecipients `json:"to,omitempty"`
}

// Stores process ID.
type ProcessIdModel struct {
	// The ID of the process.
	TypeId *uuid.UUID `json:"typeId,omitempty"`
}

// Stores project ID and its process ID.
type ProcessMigrationResultModel struct {
	// The ID of the process.
	ProcessId *uuid.UUID `json:"processId,omitempty"`
	// The ID of the project.
	ProjectId *uuid.UUID `json:"projectId,omitempty"`
}

// Project work item type state colors
type ProjectWorkItemStateColors struct {
	// Project name
	ProjectName *string `json:"projectName,omitempty"`
	// State colors for all work item type in a project
	WorkItemTypeStateColors *[]WorkItemTypeStateColors `json:"workItemTypeStateColors,omitempty"`
}

// Enumerates the possible provisioning actions that can be triggered on process template update.
type ProvisioningActionType string

type provisioningActionTypeValuesType struct {
	Import   ProvisioningActionType
	Validate ProvisioningActionType
}

var ProvisioningActionTypeValues = provisioningActionTypeValuesType{
	Import:   "import",
	Validate: "validate",
}

// Result of an update work item type XML update operation.
type ProvisioningResult struct {
	// Details about of the provisioning import events.
	ProvisioningImportEvents *[]string `json:"provisioningImportEvents,omitempty"`
}

// Describes a request to get a list of queries
type QueryBatchGetRequest struct {
	// The expand parameters for queries. Possible options are { None, Wiql, Clauses, All, Minimal }
	Expand *QueryExpand `json:"$expand,omitempty"`
	// The flag to control error policy in a query batch request. Possible options are { Fail, Omit }.
	ErrorPolicy *QueryErrorPolicy `json:"errorPolicy,omitempty"`
	// The requested query ids
	Ids *[]uuid.UUID `json:"ids,omitempty"`
}

// Enum to control error policy in a query batch request.
type QueryErrorPolicy string

type queryErrorPolicyValuesType struct {
	Fail QueryErrorPolicy
	Omit QueryErrorPolicy
}

var QueryErrorPolicyValues = queryErrorPolicyValuesType{
	Fail: "fail",
	Omit: "omit",
}

// Determines which set of additional query properties to display
type QueryExpand string

type queryExpandValuesType struct {
	None    QueryExpand
	Wiql    QueryExpand
	Clauses QueryExpand
	All     QueryExpand
	Minimal QueryExpand
}

var QueryExpandValues = queryExpandValuesType{
	// Expands Columns, Links and ChangeInfo
	None: "none",
	// Expands Columns, Links,  ChangeInfo and WIQL text
	Wiql: "wiql",
	// Expands Columns, Links, ChangeInfo, WIQL text and clauses
	Clauses: "clauses",
	// Expands all properties
	All: "all",
	// Displays minimal properties and the WIQL text
	Minimal: "minimal",
}

// Represents an item in the work item query hierarchy. This can be either a query or a folder.
type QueryHierarchyItem struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// The clauses for a flat query.
	Clauses *WorkItemQueryClause `json:"clauses,omitempty"`
	// The columns of the query.
	Columns *[]WorkItemFieldReference `json:"columns,omitempty"`
	// The identity who created the query item.
	CreatedBy *IdentityReference `json:"createdBy,omitempty"`
	// When the query item was created.
	CreatedDate *azuredevops.Time `json:"createdDate,omitempty"`
	// The link query mode.
	FilterOptions *LinkQueryMode `json:"filterOptions,omitempty"`
	// If this is a query folder, indicates if it contains any children.
	HasChildren *bool `json:"hasChildren,omitempty"`
	// The child query items inside a query folder.
	Children *[]QueryHierarchyItem `json:"children,omitempty"`
	// The id of the query item.
	Id *uuid.UUID `json:"id,omitempty"`
	// Indicates if this query item is deleted. Setting this to false on a deleted query item will undelete it. Undeleting a query or folder will not bring back the permission changes that were previously applied to it.
	IsDeleted *bool `json:"isDeleted,omitempty"`
	// Indicates if this is a query folder or a query.
	IsFolder *bool `json:"isFolder,omitempty"`
	// Indicates if the WIQL of this query is invalid. This could be due to invalid syntax or a no longer valid area/iteration path.
	IsInvalidSyntax *bool `json:"isInvalidSyntax,omitempty"`
	// Indicates if this query item is public or private.
	IsPublic *bool `json:"isPublic,omitempty"`
	// The identity who last ran the query.
	LastExecutedBy *IdentityReference `json:"lastExecutedBy,omitempty"`
	// When the query was last run.
	LastExecutedDate *azuredevops.Time `json:"lastExecutedDate,omitempty"`
	// The identity who last modified the query item.
	LastModifiedBy *IdentityReference `json:"lastModifiedBy,omitempty"`
	// When the query item was last modified.
	LastModifiedDate *azuredevops.Time `json:"lastModifiedDate,omitempty"`
	// The link query clause.
	LinkClauses *WorkItemQueryClause `json:"linkClauses,omitempty"`
	// The name of the query item.
	Name *string `json:"name,omitempty"`
	// The path of the query item.
	Path *string `json:"path,omitempty"`
	// The recursion option for use in a tree query.
	QueryRecursionOption *QueryRecursionOption `json:"queryRecursionOption,omitempty"`
	// The type of query.
	QueryType *QueryType `json:"queryType,omitempty"`
	// The sort columns of the query.
	SortColumns *[]WorkItemQuerySortColumn `json:"sortColumns,omitempty"`
	// The source clauses in a tree or one-hop link query.
	SourceClauses *WorkItemQueryClause `json:"sourceClauses,omitempty"`
	// The target clauses in a tree or one-hop link query.
	TargetClauses *WorkItemQueryClause `json:"targetClauses,omitempty"`
	// The WIQL text of the query
	Wiql *string `json:"wiql,omitempty"`
}

type QueryHierarchyItemsResult struct {
	// The count of items.
	Count *int `json:"count,omitempty"`
	// Indicates if the max return limit was hit but there are still more items
	HasMore *bool `json:"hasMore,omitempty"`
	// The list of items
	Value *[]QueryHierarchyItem `json:"value,omitempty"`
}

type QueryOption string

type queryOptionValuesType struct {
	Doing    QueryOption
	Done     QueryOption
	Followed QueryOption
}

var QueryOptionValues = queryOptionValuesType{
	Doing:    "doing",
	Done:     "done",
	Followed: "followed",
}

// Determines whether a tree query matches parents or children first.
type QueryRecursionOption string

type queryRecursionOptionValuesType struct {
	ParentFirst QueryRecursionOption
	ChildFirst  QueryRecursionOption
}

var QueryRecursionOptionValues = queryRecursionOptionValuesType{
	// Returns work items that satisfy the source, even if no linked work item satisfies the target and link criteria.
	ParentFirst: "parentFirst",
	// Returns work items that satisfy the target criteria, even if no work item satisfies the source and link criteria.
	ChildFirst: "childFirst",
}

// The query result type
type QueryResultType string

type queryResultTypeValuesType struct {
	WorkItem     QueryResultType
	WorkItemLink QueryResultType
}

var QueryResultTypeValues = queryResultTypeValuesType{
	// A list of work items (for flat queries).
	WorkItem: "workItem",
	// A list of work item links (for OneHop and Tree queries).
	WorkItemLink: "workItemLink",
}

// The type of query.
type QueryType string

type queryTypeValuesType struct {
	Flat   QueryType
	Tree   QueryType
	OneHop QueryType
}

var QueryTypeValues = queryTypeValuesType{
	// Gets a flat list of work items.
	Flat: "flat",
	// Gets a tree of work items showing their link hierarchy.
	Tree: "tree",
	// Gets a list of work items and their direct links.
	OneHop: "oneHop",
}

// The reporting revision expand level.
type ReportingRevisionsExpand string

type reportingRevisionsExpandValuesType struct {
	None   ReportingRevisionsExpand
	Fields ReportingRevisionsExpand
}

var ReportingRevisionsExpandValues = reportingRevisionsExpandValuesType{
	// Default behavior.
	None: "none",
	// Add fields to the response.
	Fields: "fields",
}

type ReportingWorkItemLinksBatch struct {
	// ContinuationToken acts as a waterMark. Used while querying large results.
	ContinuationToken *string `json:"continuationToken,omitempty"`
	// Returns 'true' if it's last batch, 'false' otherwise.
	IsLastBatch *bool `json:"isLastBatch,omitempty"`
	// The next link for the work item.
	NextLink *string `json:"nextLink,omitempty"`
	// Values such as rel, sourceId, TargetId, ChangedDate, isActive.
	Values *[]interface{} `json:"values,omitempty"`
}

type ReportingWorkItemRevisionsBatch struct {
	// ContinuationToken acts as a waterMark. Used while querying large results.
	ContinuationToken *string `json:"continuationToken,omitempty"`
	// Returns 'true' if it's last batch, 'false' otherwise.
	IsLastBatch *bool `json:"isLastBatch,omitempty"`
	// The next link for the work item.
	NextLink *string `json:"nextLink,omitempty"`
	// Values such as rel, sourceId, TargetId, ChangedDate, isActive.
	Values *[]interface{} `json:"values,omitempty"`
}

// The class represents the reporting work item revision filer.
type ReportingWorkItemRevisionsFilter struct {
	// A list of fields to return in work item revisions. Omit this parameter to get all reportable fields.
	Fields *[]string `json:"fields,omitempty"`
	// Include deleted work item in the result.
	IncludeDeleted *bool `json:"includeDeleted,omitempty"`
	// Return an identity reference instead of a string value for identity fields.
	IncludeIdentityRef *bool `json:"includeIdentityRef,omitempty"`
	// Include only the latest version of a work item, skipping over all previous revisions of the work item.
	IncludeLatestOnly *bool `json:"includeLatestOnly,omitempty"`
	// Include tag reference instead of string value for System.Tags field
	IncludeTagRef *bool `json:"includeTagRef,omitempty"`
	// A list of types to filter the results to specific work item types. Omit this parameter to get work item revisions of all work item types.
	Types *[]string `json:"types,omitempty"`
}

type SendMailBody struct {
	Fields        *[]string    `json:"fields,omitempty"`
	Ids           *[]int       `json:"ids,omitempty"`
	Message       *MailMessage `json:"message,omitempty"`
	PersistenceId *uuid.UUID   `json:"persistenceId,omitempty"`
	ProjectId     *string      `json:"projectId,omitempty"`
	SortFields    *[]string    `json:"sortFields,omitempty"`
	TempQueryId   *string      `json:"tempQueryId,omitempty"`
	Wiql          *string      `json:"wiql,omitempty"`
}

// The class describes reporting work item revision batch.
type StreamedBatch struct {
	// ContinuationToken acts as a waterMark. Used while querying large results.
	ContinuationToken *string `json:"continuationToken,omitempty"`
	// Returns 'true' if it's last batch, 'false' otherwise.
	IsLastBatch *bool `json:"isLastBatch,omitempty"`
	// The next link for the work item.
	NextLink *string `json:"nextLink,omitempty"`
	// Values such as rel, sourceId, TargetId, ChangedDate, isActive.
	Values *[]interface{} `json:"values,omitempty"`
}

// Enumerates types of supported xml templates used for customization.
type TemplateType string

type templateTypeValuesType struct {
	WorkItemType   TemplateType
	GlobalWorkflow TemplateType
}

var TemplateTypeValues = templateTypeValuesType{
	WorkItemType:   "workItemType",
	GlobalWorkflow: "globalWorkflow",
}

// Describes a request to create a temporary query
type TemporaryQueryRequestModel struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// The WIQL text of the temporary query
	Wiql *string `json:"wiql,omitempty"`
}

// The result of a temporary query creation.
type TemporaryQueryResponseModel struct {
	// The id of the temporary query item.
	Id *uuid.UUID `json:"id,omitempty"`
}

// Types of tree node structures.
type TreeNodeStructureType string

type treeNodeStructureTypeValuesType struct {
	Area      TreeNodeStructureType
	Iteration TreeNodeStructureType
}

var TreeNodeStructureTypeValues = treeNodeStructureTypeValuesType{
	// Area type.
	Area: "area",
	// Iteration type.
	Iteration: "iteration",
}

// Types of tree structures groups.
type TreeStructureGroup string

type treeStructureGroupValuesType struct {
	Areas      TreeStructureGroup
	Iterations TreeStructureGroup
}

var TreeStructureGroupValues = treeStructureGroupValuesType{
	Areas:      "areas",
	Iterations: "iterations",
}

// Describes an update request for a work item field.
type UpdateWorkItemField struct {
	// Indicates whether the user wants to restore the field.
	IsDeleted *bool `json:"isDeleted,omitempty"`
}

// A WIQL query
type Wiql struct {
	// The text of the WIQL query
	Query *string `json:"query,omitempty"`
}

// A work artifact link describes an outbound artifact link type.
type WorkArtifactLink struct {
	// Target artifact type.
	ArtifactType *string `json:"artifactType,omitempty"`
	// Outbound link type.
	LinkType *string `json:"linkType,omitempty"`
	// Target tool type.
	ToolType *string `json:"toolType,omitempty"`
}

// Describes a work item.
type WorkItem struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// Reference to a specific version of the comment added/edited/deleted in this revision.
	CommentVersionRef *WorkItemCommentVersionRef `json:"commentVersionRef,omitempty"`
	// Map of field and values for the work item.
	Fields *map[string]interface{} `json:"fields,omitempty"`
	// The work item ID.
	Id *int `json:"id,omitempty"`
	// Relations of the work item.
	Relations *[]WorkItemRelation `json:"relations,omitempty"`
	// Revision number of the work item.
	Rev *int `json:"rev,omitempty"`
}

// Describes a request to get a set of work items
type WorkItemBatchGetRequest struct {
	// The expand parameters for work item attributes. Possible options are { None, Relations, Fields, Links, All }
	Expand *WorkItemExpand `json:"$expand,omitempty"`
	// AsOf UTC date time string
	AsOf *azuredevops.Time `json:"asOf,omitempty"`
	// The flag to control error policy in a bulk get work items request. Possible options are {Fail, Omit}.
	ErrorPolicy *WorkItemErrorPolicy `json:"errorPolicy,omitempty"`
	// The requested fields
	Fields *[]string `json:"fields,omitempty"`
	// The requested work item ids
	Ids *[]int `json:"ids,omitempty"`
}

// Defines a classification node for work item tracking.
type WorkItemClassificationNode struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// Dictionary that has node attributes like start/finish date for iteration nodes.
	Attributes *map[string]interface{} `json:"attributes,omitempty"`
	// Flag that indicates if the classification node has any child nodes.
	HasChildren *bool `json:"hasChildren,omitempty"`
	// List of child nodes fetched.
	Children *[]WorkItemClassificationNode `json:"children,omitempty"`
	// Integer ID of the classification node.
	Id *int `json:"id,omitempty"`
	// GUID ID of the classification node.
	Identifier *uuid.UUID `json:"identifier,omitempty"`
	// Name of the classification node.
	Name *string `json:"name,omitempty"`
	// Path of the classification node.
	Path *string `json:"path,omitempty"`
	// Node structure type.
	StructureType *TreeNodeStructureType `json:"structureType,omitempty"`
}

// Comment on Work Item
type WorkItemComment struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// Represents the possible types for the comment format.
	Format *CommentFormat `json:"format,omitempty"`
	// The text of the comment in HTML format.
	RenderedText *string `json:"renderedText,omitempty"`
	// Identity of user who added the comment.
	RevisedBy *IdentityReference `json:"revisedBy,omitempty"`
	// The date of comment.
	RevisedDate *azuredevops.Time `json:"revisedDate,omitempty"`
	// The work item revision number.
	Revision *int `json:"revision,omitempty"`
	// The text of the comment.
	Text *string `json:"text,omitempty"`
}

// Collection of comments.
type WorkItemComments struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// Comments collection.
	Comments *[]WorkItemComment `json:"comments,omitempty"`
	// The count of comments.
	Count *int `json:"count,omitempty"`
	// Count of comments from the revision.
	FromRevisionCount *int `json:"fromRevisionCount,omitempty"`
	// Total count of comments.
	TotalCount *int `json:"totalCount,omitempty"`
}

// Represents the reference to a specific version of a comment on a Work Item.
type WorkItemCommentVersionRef struct {
	Url *string `json:"url,omitempty"`
	// The id assigned to the comment.
	CommentId *int `json:"commentId,omitempty"`
	// [Internal] The work item revision where this comment was originally added.
	CreatedInRevision *int `json:"createdInRevision,omitempty"`
	// [Internal] Specifies whether comment was deleted.
	IsDeleted *bool `json:"isDeleted,omitempty"`
	// [Internal] The text of the comment.
	Text *string `json:"text,omitempty"`
	// The version number.
	Version *int `json:"version,omitempty"`
}

// Full deleted work item object. Includes the work item itself.
type WorkItemDelete struct {
	// The HTTP status code for work item operation in a batch request.
	Code *int `json:"code,omitempty"`
	// The user who deleted the work item type.
	DeletedBy *string `json:"deletedBy,omitempty"`
	// The work item deletion date.
	DeletedDate *string `json:"deletedDate,omitempty"`
	// Work item ID.
	Id *int `json:"id,omitempty"`
	// The exception message for work item operation in a batch request.
	Message *string `json:"message,omitempty"`
	// Name or title of the work item.
	Name *string `json:"name,omitempty"`
	// Parent project of the deleted work item.
	Project *string `json:"project,omitempty"`
	// Type of work item.
	Type *string `json:"type,omitempty"`
	// REST API URL of the resource
	Url *string `json:"url,omitempty"`
	// The work item object that was deleted.
	Resource *WorkItem `json:"resource,omitempty"`
}

// Describes response to delete a set of work items.
type WorkItemDeleteBatch struct {
	// List of results for each work item
	Results *[]WorkItemDelete `json:"results,omitempty"`
}

// Describes a request to delete a set of work items
type WorkItemDeleteBatchRequest struct {
	// Optional parameter, if set to true, the work item is deleted permanently. Please note: the destroy action is PERMANENT and cannot be undone.
	Destroy *bool `json:"destroy,omitempty"`
	// The requested work item ids
	Ids *[]int `json:"ids,omitempty"`
	// Optional parameter, if set to true, notifications will be disabled.
	SkipNotifications *bool `json:"skipNotifications,omitempty"`
}

// Reference to a deleted work item.
type WorkItemDeleteReference struct {
	// The HTTP status code for work item operation in a batch request.
	Code *int `json:"code,omitempty"`
	// The user who deleted the work item type.
	DeletedBy *string `json:"deletedBy,omitempty"`
	// The work item deletion date.
	DeletedDate *string `json:"deletedDate,omitempty"`
	// Work item ID.
	Id *int `json:"id,omitempty"`
	// The exception message for work item operation in a batch request.
	Message *string `json:"message,omitempty"`
	// Name or title of the work item.
	Name *string `json:"name,omitempty"`
	// Parent project of the deleted work item.
	Project *string `json:"project,omitempty"`
	// Type of work item.
	Type *string `json:"type,omitempty"`
	// REST API URL of the resource
	Url *string `json:"url,omitempty"`
}

// Shallow Reference to a deleted work item.
type WorkItemDeleteShallowReference struct {
	// Work item ID.
	Id *int `json:"id,omitempty"`
	// REST API URL of the resource
	Url *string `json:"url,omitempty"`
}

// Describes an update request for a deleted work item.
type WorkItemDeleteUpdate struct {
	// Sets a value indicating whether this work item is deleted.
	IsDeleted *bool `json:"isDeleted,omitempty"`
}

// Enum to control error policy in a bulk get work items request.
type WorkItemErrorPolicy string

type workItemErrorPolicyValuesType struct {
	Fail WorkItemErrorPolicy
	Omit WorkItemErrorPolicy
}

var WorkItemErrorPolicyValues = workItemErrorPolicyValuesType{
	// Fail work error policy.
	Fail: "fail",
	// Omit work error policy.
	Omit: "omit",
}

// Flag to control payload properties from get work item command.
type WorkItemExpand string

type workItemExpandValuesType struct {
	None      WorkItemExpand
	Relations WorkItemExpand
	Fields    WorkItemExpand
	Links     WorkItemExpand
	All       WorkItemExpand
}

var WorkItemExpandValues = workItemExpandValuesType{
	// Default behavior.
	None: "none",
	// Relations work item expand.
	Relations: "relations",
	// Fields work item expand.
	Fields: "fields",
	// Links work item expand.
	Links: "links",
	// Expands all.
	All: "all",
}

// Describes a field on a work item and it's properties specific to that work item type.
type WorkItemField struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// Indicates whether the field is sortable in server queries.
	CanSortBy *bool `json:"canSortBy,omitempty"`
	// The description of the field.
	Description *string `json:"description,omitempty"`
	// Indicates whether this field is deleted.
	IsDeleted *bool `json:"isDeleted,omitempty"`
	// Indicates whether this field is an identity field.
	IsIdentity *bool `json:"isIdentity,omitempty"`
	// Indicates whether this instance is picklist.
	IsPicklist *bool `json:"isPicklist,omitempty"`
	// Indicates whether this instance is a suggested picklist .
	IsPicklistSuggested *bool `json:"isPicklistSuggested,omitempty"`
	// Indicates whether the field can be queried in the server.
	IsQueryable *bool `json:"isQueryable,omitempty"`
	// The name of the field.
	Name *string `json:"name,omitempty"`
	// If this field is picklist, the identifier of the picklist associated, otherwise null
	PicklistId *uuid.UUID `json:"picklistId,omitempty"`
	// Indicates whether the field is [read only].
	ReadOnly *bool `json:"readOnly,omitempty"`
	// The reference name of the field.
	ReferenceName *string `json:"referenceName,omitempty"`
	// The supported operations on this field.
	SupportedOperations *[]WorkItemFieldOperation `json:"supportedOperations,omitempty"`
	// The type of the field.
	Type *FieldType `json:"type,omitempty"`
	// The usage of the field.
	Usage *FieldUsage `json:"usage,omitempty"`
}

// Describes a field on a work item and it's properties specific to that work item type.
type WorkItemField2 struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// Indicates whether the field is sortable in server queries.
	CanSortBy *bool `json:"canSortBy,omitempty"`
	// The description of the field.
	Description *string `json:"description,omitempty"`
	// Indicates whether this field is deleted.
	IsDeleted *bool `json:"isDeleted,omitempty"`
	// Indicates whether this field is an identity field.
	IsIdentity *bool `json:"isIdentity,omitempty"`
	// Indicates whether this instance is picklist.
	IsPicklist *bool `json:"isPicklist,omitempty"`
	// Indicates whether this instance is a suggested picklist .
	IsPicklistSuggested *bool `json:"isPicklistSuggested,omitempty"`
	// Indicates whether the field can be queried in the server.
	IsQueryable *bool `json:"isQueryable,omitempty"`
	// The name of the field.
	Name *string `json:"name,omitempty"`
	// If this field is picklist, the identifier of the picklist associated, otherwise null
	PicklistId *uuid.UUID `json:"picklistId,omitempty"`
	// Indicates whether the field is [read only].
	ReadOnly *bool `json:"readOnly,omitempty"`
	// The reference name of the field.
	ReferenceName *string `json:"referenceName,omitempty"`
	// The supported operations on this field.
	SupportedOperations *[]WorkItemFieldOperation `json:"supportedOperations,omitempty"`
	// The type of the field.
	Type *FieldType `json:"type,omitempty"`
	// The usage of the field.
	Usage *FieldUsage `json:"usage,omitempty"`
	// Indicates whether this field is marked as locked for editing.
	IsLocked *bool `json:"isLocked,omitempty"`
}

// Describes the list of allowed values of the field.
type WorkItemFieldAllowedValues struct {
	// The list of field allowed values.
	AllowedValues *[]string `json:"allowedValues,omitempty"`
	// Name of the field.
	FieldName *string `json:"fieldName,omitempty"`
}

// Describes a work item field operation.
type WorkItemFieldOperation struct {
	// Friendly name of the operation.
	Name *string `json:"name,omitempty"`
	// Reference name of the operation.
	ReferenceName *string `json:"referenceName,omitempty"`
}

// Reference to a field in a work item
type WorkItemFieldReference struct {
	// The friendly name of the field.
	Name *string `json:"name,omitempty"`
	// The reference name of the field.
	ReferenceName *string `json:"referenceName,omitempty"`
	// The REST URL of the resource.
	Url *string `json:"url,omitempty"`
}

// Describes an update to a work item field.
type WorkItemFieldUpdate struct {
	// The new value of the field.
	NewValue interface{} `json:"newValue,omitempty"`
	// The old value of the field.
	OldValue interface{} `json:"oldValue,omitempty"`
}

type WorkItemHistory struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links       interface{}        `json:"_links,omitempty"`
	Rev         *int               `json:"rev,omitempty"`
	RevisedBy   *IdentityReference `json:"revisedBy,omitempty"`
	RevisedDate *azuredevops.Time  `json:"revisedDate,omitempty"`
	Value       *string            `json:"value,omitempty"`
}

// Reference to a work item icon.
type WorkItemIcon struct {
	// The identifier of the icon.
	Id *string `json:"id,omitempty"`
	// The REST URL of the resource.
	Url *string `json:"url,omitempty"`
}

// A link between two work items.
type WorkItemLink struct {
	// The type of link.
	Rel *string `json:"rel,omitempty"`
	// The source work item.
	Source *WorkItemReference `json:"source,omitempty"`
	// The target work item.
	Target *WorkItemReference `json:"target,omitempty"`
}

// Describes the next state for a work item.
type WorkItemNextStateOnTransition struct {
	// Error code if there is no next state transition possible.
	ErrorCode *string `json:"errorCode,omitempty"`
	// Work item ID.
	Id *int `json:"id,omitempty"`
	// Error message if there is no next state transition possible.
	Message *string `json:"message,omitempty"`
	// Name of the next state on transition.
	StateOnTransition *string `json:"stateOnTransition,omitempty"`
}

// Represents a clause in a work item query. This shows the structure of a work item query.
type WorkItemQueryClause struct {
	// Child clauses if the current clause is a logical operator
	Clauses *[]WorkItemQueryClause `json:"clauses,omitempty"`
	// Field associated with condition
	Field *WorkItemFieldReference `json:"field,omitempty"`
	// Right side of the condition when a field to field comparison
	FieldValue *WorkItemFieldReference `json:"fieldValue,omitempty"`
	// Determines if this is a field to field comparison
	IsFieldValue *bool `json:"isFieldValue,omitempty"`
	// Logical operator separating the condition clause
	LogicalOperator *LogicalOperation `json:"logicalOperator,omitempty"`
	// The field operator
	Operator *WorkItemFieldOperation `json:"operator,omitempty"`
	// Right side of the condition when a field to value comparison
	Value *string `json:"value,omitempty"`
}

// The result of a work item query.
type WorkItemQueryResult struct {
	// The date the query was run in the context of.
	AsOf *azuredevops.Time `json:"asOf,omitempty"`
	// The columns of the query.
	Columns *[]WorkItemFieldReference `json:"columns,omitempty"`
	// The result type
	QueryResultType *QueryResultType `json:"queryResultType,omitempty"`
	// The type of the query
	QueryType *QueryType `json:"queryType,omitempty"`
	// The sort columns of the query.
	SortColumns *[]WorkItemQuerySortColumn `json:"sortColumns,omitempty"`
	// The work item links returned by the query.
	WorkItemRelations *[]WorkItemLink `json:"workItemRelations,omitempty"`
	// The work items returned by the query.
	WorkItems *[]WorkItemReference `json:"workItems,omitempty"`
}

// A sort column.
type WorkItemQuerySortColumn struct {
	// The direction to sort by.
	Descending *bool `json:"descending,omitempty"`
	// A work item field.
	Field *WorkItemFieldReference `json:"field,omitempty"`
}

// Type of the activity
type WorkItemRecentActivityType string

type workItemRecentActivityTypeValuesType struct {
	Visited  WorkItemRecentActivityType
	Edited   WorkItemRecentActivityType
	Deleted  WorkItemRecentActivityType
	Restored WorkItemRecentActivityType
}

var WorkItemRecentActivityTypeValues = workItemRecentActivityTypeValuesType{
	Visited:  "visited",
	Edited:   "edited",
	Deleted:  "deleted",
	Restored: "restored",
}

// Contains reference to a work item.
type WorkItemReference struct {
	// Work item ID.
	Id *int `json:"id,omitempty"`
	// REST API URL of the resource
	Url *string `json:"url,omitempty"`
}

type WorkItemRelation struct {
	// Collection of link attributes.
	Attributes *map[string]interface{} `json:"attributes,omitempty"`
	// Relation type.
	Rel *string `json:"rel,omitempty"`
	// Link url.
	Url *string `json:"url,omitempty"`
}

// Represents the work item type relation type.
type WorkItemRelationType struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// The name.
	Name *string `json:"name,omitempty"`
	// The reference name.
	ReferenceName *string `json:"referenceName,omitempty"`
	// The collection of relation type attributes.
	Attributes *map[string]interface{} `json:"attributes,omitempty"`
}

// Describes updates to a work item's relations.
type WorkItemRelationUpdates struct {
	// List of newly added relations.
	Added *[]WorkItemRelation `json:"added,omitempty"`
	// List of removed relations.
	Removed *[]WorkItemRelation `json:"removed,omitempty"`
	// List of updated relations.
	Updated *[]WorkItemRelation `json:"updated,omitempty"`
}

// Work item type state name, color and state category
type WorkItemStateColor struct {
	// Category of state
	Category *string `json:"category,omitempty"`
	// Color value
	Color *string `json:"color,omitempty"`
	// Work item type state name
	Name *string `json:"name,omitempty"`
}

// Describes a state transition in a work item.
type WorkItemStateTransition struct {
	// Gets a list of actions needed to transition to that state.
	Actions *[]string `json:"actions,omitempty"`
	// Name of the next state.
	To *string `json:"to,omitempty"`
}

type WorkItemTagDefinition struct {
	Id          *uuid.UUID        `json:"id,omitempty"`
	LastUpdated *azuredevops.Time `json:"lastUpdated,omitempty"`
	Name        *string           `json:"name,omitempty"`
	Url         *string           `json:"url,omitempty"`
}

// Describes a work item template.
type WorkItemTemplate struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// The description of the work item template.
	Description *string `json:"description,omitempty"`
	// The identifier of the work item template.
	Id *uuid.UUID `json:"id,omitempty"`
	// The name of the work item template.
	Name *string `json:"name,omitempty"`
	// The name of the work item type.
	WorkItemTypeName *string `json:"workItemTypeName,omitempty"`
	// Mapping of field and its templated value.
	Fields *map[string]string `json:"fields,omitempty"`
}

// Describes a shallow reference to a work item template.
type WorkItemTemplateReference struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// The description of the work item template.
	Description *string `json:"description,omitempty"`
	// The identifier of the work item template.
	Id *uuid.UUID `json:"id,omitempty"`
	// The name of the work item template.
	Name *string `json:"name,omitempty"`
	// The name of the work item type.
	WorkItemTypeName *string `json:"workItemTypeName,omitempty"`
}

type WorkItemTrackingReference struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// The name.
	Name *string `json:"name,omitempty"`
	// The reference name.
	ReferenceName *string `json:"referenceName,omitempty"`
}

// Base class for WIT REST resources.
type WorkItemTrackingResource struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
}

// Base class for work item tracking resource references.
type WorkItemTrackingResourceReference struct {
	Url *string `json:"url,omitempty"`
}

// Describes a work item type.
type WorkItemType struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// The color.
	Color *string `json:"color,omitempty"`
	// The description of the work item type.
	Description *string `json:"description,omitempty"`
	// The fields that exist on the work item type.
	FieldInstances *[]WorkItemTypeFieldInstance `json:"fieldInstances,omitempty"`
	// The fields that exist on the work item type.
	Fields *[]WorkItemTypeFieldInstance `json:"fields,omitempty"`
	// The icon of the work item type.
	Icon *WorkItemIcon `json:"icon,omitempty"`
	// True if work item type is disabled
	IsDisabled *bool `json:"isDisabled,omitempty"`
	// Gets the name of the work item type.
	Name *string `json:"name,omitempty"`
	// The reference name of the work item type.
	ReferenceName *string `json:"referenceName,omitempty"`
	// Gets state information for the work item type.
	States *[]WorkItemStateColor `json:"states,omitempty"`
	// Gets the various state transition mappings in the work item type.
	Transitions *map[string][]WorkItemStateTransition `json:"transitions,omitempty"`
	// The XML form.
	XmlForm *string `json:"xmlForm,omitempty"`
}

// Describes a work item type category.
type WorkItemTypeCategory struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// Gets or sets the default type of the work item.
	DefaultWorkItemType *WorkItemTypeReference `json:"defaultWorkItemType,omitempty"`
	// The name of the category.
	Name *string `json:"name,omitempty"`
	// The reference name of the category.
	ReferenceName *string `json:"referenceName,omitempty"`
	// The work item types that belong to the category.
	WorkItemTypes *[]WorkItemTypeReference `json:"workItemTypes,omitempty"`
}

// Describes a work item type's colors.
type WorkItemTypeColor struct {
	// Gets or sets the color of the primary.
	PrimaryColor *string `json:"primaryColor,omitempty"`
	// Gets or sets the color of the secondary.
	SecondaryColor *string `json:"secondaryColor,omitempty"`
	// The name of the work item type.
	WorkItemTypeName *string `json:"workItemTypeName,omitempty"`
}

// Describes work item type name, its icon and color.
type WorkItemTypeColorAndIcon struct {
	// The color of the work item type in hex format.
	Color *string `json:"color,omitempty"`
	// The work item type icon.
	Icon *string `json:"icon,omitempty"`
	// Indicates if the work item is disabled in the process.
	IsDisabled *bool `json:"isDisabled,omitempty"`
	// The name of the work item type.
	WorkItemTypeName *string `json:"workItemTypeName,omitempty"`
}

// Field instance of a work item type.
type WorkItemTypeFieldInstance struct {
	// The friendly name of the field.
	Name *string `json:"name,omitempty"`
	// The reference name of the field.
	ReferenceName *string `json:"referenceName,omitempty"`
	// The REST URL of the resource.
	Url *string `json:"url,omitempty"`
	// Indicates whether field value is always required.
	AlwaysRequired *bool `json:"alwaysRequired,omitempty"`
	// The list of dependent fields.
	DependentFields *[]WorkItemFieldReference `json:"dependentFields,omitempty"`
	// Gets the help text for the field.
	HelpText *string `json:"helpText,omitempty"`
	// The list of field allowed values.
	AllowedValues *[]string `json:"allowedValues,omitempty"`
	// Represents the default value of the field.
	DefaultValue *string `json:"defaultValue,omitempty"`
}

// Base field instance for workItemType fields.
type WorkItemTypeFieldInstanceBase struct {
	// The friendly name of the field.
	Name *string `json:"name,omitempty"`
	// The reference name of the field.
	ReferenceName *string `json:"referenceName,omitempty"`
	// The REST URL of the resource.
	Url *string `json:"url,omitempty"`
	// Indicates whether field value is always required.
	AlwaysRequired *bool `json:"alwaysRequired,omitempty"`
	// The list of dependent fields.
	DependentFields *[]WorkItemFieldReference `json:"dependentFields,omitempty"`
	// Gets the help text for the field.
	HelpText *string `json:"helpText,omitempty"`
}

// Expand options for the work item field(s) request.
type WorkItemTypeFieldsExpandLevel string

type workItemTypeFieldsExpandLevelValuesType struct {
	None            WorkItemTypeFieldsExpandLevel
	AllowedValues   WorkItemTypeFieldsExpandLevel
	DependentFields WorkItemTypeFieldsExpandLevel
	All             WorkItemTypeFieldsExpandLevel
}

var WorkItemTypeFieldsExpandLevelValues = workItemTypeFieldsExpandLevelValuesType{
	// Includes only basic properties of the field.
	None: "none",
	// Includes allowed values for the field.
	AllowedValues: "allowedValues",
	// Includes dependent fields of the field.
	DependentFields: "dependentFields",
	// Includes allowed values and dependent fields of the field.
	All: "all",
}

// Field Instance of a workItemype with detailed references.
type WorkItemTypeFieldWithReferences struct {
	// The friendly name of the field.
	Name *string `json:"name,omitempty"`
	// The reference name of the field.
	ReferenceName *string `json:"referenceName,omitempty"`
	// The REST URL of the resource.
	Url *string `json:"url,omitempty"`
	// Indicates whether field value is always required.
	AlwaysRequired *bool `json:"alwaysRequired,omitempty"`
	// The list of dependent fields.
	DependentFields *[]WorkItemFieldReference `json:"dependentFields,omitempty"`
	// Gets the help text for the field.
	HelpText *string `json:"helpText,omitempty"`
	// The list of field allowed values.
	AllowedValues *[]interface{} `json:"allowedValues,omitempty"`
	// Represents the default value of the field.
	DefaultValue interface{} `json:"defaultValue,omitempty"`
}

// Reference to a work item type.
type WorkItemTypeReference struct {
	Url *string `json:"url,omitempty"`
	// Name of the work item type.
	Name *string `json:"name,omitempty"`
}

// State colors for a work item type
type WorkItemTypeStateColors struct {
	// Work item type state colors
	StateColors *[]WorkItemStateColor `json:"stateColors,omitempty"`
	// Work item type name
	WorkItemTypeName *string `json:"workItemTypeName,omitempty"`
}

// Describes a work item type template.
type WorkItemTypeTemplate struct {
	// XML template in string format.
	Template *string `json:"template,omitempty"`
}

// Describes a update work item type template request body.
type WorkItemTypeTemplateUpdateModel struct {
	// Describes the type of the action for the update request.
	ActionType *ProvisioningActionType `json:"actionType,omitempty"`
	// Methodology to which the template belongs, eg. Agile, Scrum, CMMI.
	Methodology *string `json:"methodology,omitempty"`
	// String representation of the work item type template.
	Template *string `json:"template,omitempty"`
	// The type of the template described in the request body.
	TemplateType *TemplateType `json:"templateType,omitempty"`
}

// Describes an update to a work item.
type WorkItemUpdate struct {
	Url *string `json:"url,omitempty"`
	// Link references to related REST resources.
	Links interface{} `json:"_links,omitempty"`
	// List of updates to fields.
	Fields *map[string]WorkItemFieldUpdate `json:"fields,omitempty"`
	// ID of update.
	Id *int `json:"id,omitempty"`
	// List of updates to relations.
	Relations *WorkItemRelationUpdates `json:"relations,omitempty"`
	// The revision number of work item update.
	Rev *int `json:"rev,omitempty"`
	// Identity for the work item update.
	RevisedBy *IdentityReference `json:"revisedBy,omitempty"`
	// The work item updates revision date.
	RevisedDate *azuredevops.Time `json:"revisedDate,omitempty"`
	// The work item ID.
	WorkItemId *int `json:"workItemId,omitempty"`
}

// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package work

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
)

type Activity struct {
	CapacityPerDay *float32 `json:"capacityPerDay,omitempty"`
	Name           *string  `json:"name,omitempty"`
}

type attribute struct {
}

type BacklogColumn struct {
	ColumnFieldReference *workitemtracking.WorkItemFieldReference `json:"columnFieldReference,omitempty"`
	Width                *int                                     `json:"width,omitempty"`
}

type BacklogConfiguration struct {
	// Behavior/type field mapping
	BacklogFields *BacklogFields `json:"backlogFields,omitempty"`
	// Bugs behavior
	BugsBehavior *BugsBehavior `json:"bugsBehavior,omitempty"`
	// Hidden Backlog
	HiddenBacklogs *[]string `json:"hiddenBacklogs,omitempty"`
	// Is BugsBehavior Configured in the process
	IsBugsBehaviorConfigured *bool `json:"isBugsBehaviorConfigured,omitempty"`
	// Portfolio backlog descriptors
	PortfolioBacklogs *[]BacklogLevelConfiguration `json:"portfolioBacklogs,omitempty"`
	// Requirement backlog
	RequirementBacklog *BacklogLevelConfiguration `json:"requirementBacklog,omitempty"`
	// Task backlog
	TaskBacklog *BacklogLevelConfiguration `json:"taskBacklog,omitempty"`
	Url         *string                    `json:"url,omitempty"`
	// Mapped states for work item types
	WorkItemTypeMappedStates *[]WorkItemTypeStateInfo `json:"workItemTypeMappedStates,omitempty"`
}

type BacklogFields struct {
	// Field Type (e.g. Order, Activity) to Field Reference Name map
	TypeFields *map[string]string `json:"typeFields,omitempty"`
}

// Contract representing a backlog level
type BacklogLevel struct {
	// Reference name of the corresponding WIT category
	CategoryReferenceName *string `json:"categoryReferenceName,omitempty"`
	// Plural name for the backlog level
	PluralName *string `json:"pluralName,omitempty"`
	// Collection of work item states that are included in the plan. The server will filter to only these work item types.
	WorkItemStates *[]string `json:"workItemStates,omitempty"`
	// Collection of valid workitem type names for the given backlog level
	WorkItemTypes *[]string `json:"workItemTypes,omitempty"`
}

type BacklogLevelConfiguration struct {
	// List of fields to include in Add Panel
	AddPanelFields *[]workitemtracking.WorkItemFieldReference `json:"addPanelFields,omitempty"`
	// Color for the backlog level
	Color *string `json:"color,omitempty"`
	// Default list of columns for the backlog
	ColumnFields *[]BacklogColumn `json:"columnFields,omitempty"`
	// Default Work Item Type for the backlog
	DefaultWorkItemType *workitemtracking.WorkItemTypeReference `json:"defaultWorkItemType,omitempty"`
	// Backlog Id (for Legacy Backlog Level from process config it can be categoryref name)
	Id *string `json:"id,omitempty"`
	// Indicates whether the backlog level is hidden
	IsHidden *bool `json:"isHidden,omitempty"`
	// Backlog Name
	Name *string `json:"name,omitempty"`
	// Backlog Rank (Taskbacklog is 0)
	Rank *int `json:"rank,omitempty"`
	// The type of this backlog level
	Type *BacklogType `json:"type,omitempty"`
	// Max number of work items to show in the given backlog
	WorkItemCountLimit *int `json:"workItemCountLimit,omitempty"`
	// Work Item types participating in this backlog as known by the project/Process, can be overridden by team settings for bugs
	WorkItemTypes *[]workitemtracking.WorkItemTypeReference `json:"workItemTypes,omitempty"`
}

// Represents work items in a backlog level
type BacklogLevelWorkItems struct {
	// A list of work items within a backlog level
	WorkItems *[]workitemtracking.WorkItemLink `json:"workItems,omitempty"`
}

// Definition of the type of backlog level
type BacklogType string

type backlogTypeValuesType struct {
	Portfolio   BacklogType
	Requirement BacklogType
	Task        BacklogType
}

var BacklogTypeValues = backlogTypeValuesType{
	// Portfolio backlog level
	Portfolio: "portfolio",
	// Requirement backlog level
	Requirement: "requirement",
	// Task backlog level
	Task: "task",
}

type Board struct {
	// Id of the resource
	Id *uuid.UUID `json:"id,omitempty"`
	// Name of the resource
	Name *string `json:"name,omitempty"`
	// Full http link to the resource
	Url             *string                         `json:"url,omitempty"`
	Links           interface{}                     `json:"_links,omitempty"`
	AllowedMappings *map[string]map[string][]string `json:"allowedMappings,omitempty"`
	CanEdit         *bool                           `json:"canEdit,omitempty"`
	Columns         *[]BoardColumn                  `json:"columns,omitempty"`
	Fields          *BoardFields                    `json:"fields,omitempty"`
	IsValid         *bool                           `json:"isValid,omitempty"`
	Revision        *int                            `json:"revision,omitempty"`
	Rows            *[]BoardRow                     `json:"rows,omitempty"`
}

// Represents a board badge.
type BoardBadge struct {
	// The ID of the board represented by this badge.
	BoardId *uuid.UUID `json:"boardId,omitempty"`
	// A link to the SVG resource.
	ImageUrl *string `json:"imageUrl,omitempty"`
}

// Determines what columns to include on the board badge
type BoardBadgeColumnOptions string

type boardBadgeColumnOptionsValuesType struct {
	InProgressColumns BoardBadgeColumnOptions
	AllColumns        BoardBadgeColumnOptions
	CustomColumns     BoardBadgeColumnOptions
}

var BoardBadgeColumnOptionsValues = boardBadgeColumnOptionsValuesType{
	// Only include In Progress columns
	InProgressColumns: "inProgressColumns",
	// Include all columns
	AllColumns: "allColumns",
	// Include a custom set of columns
	CustomColumns: "customColumns",
}

type BoardCardRuleSettings struct {
	Links interface{}        `json:"_links,omitempty"`
	Rules *map[string][]Rule `json:"rules,omitempty"`
	Url   *string            `json:"url,omitempty"`
}

type BoardCardSettings struct {
	Cards *map[string][]FieldSetting `json:"cards,omitempty"`
}

type BoardColumn struct {
	ColumnType    *BoardColumnType   `json:"columnType,omitempty"`
	Description   *string            `json:"description,omitempty"`
	Id            *uuid.UUID         `json:"id,omitempty"`
	IsSplit       *bool              `json:"isSplit,omitempty"`
	ItemLimit     *int               `json:"itemLimit,omitempty"`
	Name          *string            `json:"name,omitempty"`
	StateMappings *map[string]string `json:"stateMappings,omitempty"`
}

type BoardColumnType string

type boardColumnTypeValuesType struct {
	Incoming   BoardColumnType
	InProgress BoardColumnType
	Outgoing   BoardColumnType
}

var BoardColumnTypeValues = boardColumnTypeValuesType{
	Incoming:   "incoming",
	InProgress: "inProgress",
	Outgoing:   "outgoing",
}

type BoardFields struct {
	ColumnField *FieldReference `json:"columnField,omitempty"`
	DoneField   *FieldReference `json:"doneField,omitempty"`
	RowField    *FieldReference `json:"rowField,omitempty"`
}

type BoardChart struct {
	// Name of the resource
	Name *string `json:"name,omitempty"`
	// Full http link to the resource
	Url *string `json:"url,omitempty"`
	// The links for the resource
	Links interface{} `json:"_links,omitempty"`
	// The settings for the resource
	Settings *map[string]interface{} `json:"settings,omitempty"`
}

type BoardChartReference struct {
	// Name of the resource
	Name *string `json:"name,omitempty"`
	// Full http link to the resource
	Url *string `json:"url,omitempty"`
}

type BoardReference struct {
	// Id of the resource
	Id *uuid.UUID `json:"id,omitempty"`
	// Name of the resource
	Name *string `json:"name,omitempty"`
	// Full http link to the resource
	Url *string `json:"url,omitempty"`
}

type BoardRow struct {
	Color *string    `json:"color,omitempty"`
	Id    *uuid.UUID `json:"id,omitempty"`
	Name  *string    `json:"name,omitempty"`
}

type BoardSuggestedValue struct {
	Name *string `json:"name,omitempty"`
}

type BoardUserSettings struct {
	AutoRefreshState *bool `json:"autoRefreshState,omitempty"`
}

// The behavior of the work item types that are in the work item category specified in the BugWorkItems section in the Process Configuration
type BugsBehavior string

type bugsBehaviorValuesType struct {
	Off            BugsBehavior
	AsRequirements BugsBehavior
	AsTasks        BugsBehavior
}

var BugsBehaviorValues = bugsBehaviorValuesType{
	Off:            "off",
	AsRequirements: "asRequirements",
	AsTasks:        "asTasks",
}

type CapacityContractBase struct {
	// Collection of links relevant to resource
	Links interface{} `json:"_links,omitempty"`
	// Full http link to the resource
	Url *string `json:"url,omitempty"`
	// Collection of capacities associated with the team member
	Activities *[]Activity `json:"activities,omitempty"`
	// The days off associated with the team member
	DaysOff *[]DateRange `json:"daysOff,omitempty"`
}

// Expected data from PATCH
type CapacityPatch struct {
	Activities *[]Activity  `json:"activities,omitempty"`
	DaysOff    *[]DateRange `json:"daysOff,omitempty"`
}

// Card settings, such as fields and rules
type CardFieldSettings struct {
	// A collection of field information of additional fields on cards. The index in the collection signifies the order of the field among the additional fields. Currently unused. Should be used with User Story 691539: Card setting: additional fields
	AdditionalFields *[]FieldInfo `json:"additionalFields,omitempty"`
	// Display format for the assigned to field
	AssignedToDisplayFormat *IdentityDisplayFormat `json:"assignedToDisplayFormat,omitempty"`
	// A collection of field information of rendered core fields on cards.
	CoreFields *[]FieldInfo `json:"coreFields,omitempty"`
	// Flag indicating whether to show assigned to field on cards. When true, AssignedToDisplayFormat will determine how the field will be displayed
	ShowAssignedTo *bool `json:"showAssignedTo,omitempty"`
	// Flag indicating whether to show empty fields on cards
	ShowEmptyFields *bool `json:"showEmptyFields,omitempty"`
	// Flag indicating whether to show child rollup on cards
	ShowChildRollup *bool `json:"showChildRollup,omitempty"`
	// Flag indicating whether to show ID on cards
	ShowId *bool `json:"showId,omitempty"`
	// Flag indicating whether to show parent field on cards
	ShowParent *bool `json:"showParent,omitempty"`
	// Flag indicating whether to show state field on cards
	ShowState *bool `json:"showState,omitempty"`
	// Flag indicating whether to show tags on cards
	ShowTags *bool `json:"showTags,omitempty"`
}

// Card settings, such as fields and rules
type CardSettings struct {
	// A collection of settings related to rendering of fields on cards
	Fields *CardFieldSettings `json:"fields,omitempty"`
}

// Details about a given backlog category
type CategoryConfiguration struct {
	// Name
	Name *string `json:"name,omitempty"`
	// Category Reference Name
	ReferenceName *string `json:"referenceName,omitempty"`
	// Work item types for the backlog category
	WorkItemTypes *[]workitemtracking.WorkItemTypeReference `json:"workItemTypes,omitempty"`
}

type CreatePlan struct {
	// Description of the plan
	Description *string `json:"description,omitempty"`
	// Name of the plan to create.
	Name *string `json:"name,omitempty"`
	// Plan properties.
	Properties interface{} `json:"properties,omitempty"`
	// Type of plan to create.
	Type *PlanType `json:"type,omitempty"`
}

type DateRange struct {
	// End of the date range.
	End *azuredevops.Time `json:"end,omitempty"`
	// Start of the date range.
	Start *azuredevops.Time `json:"start,omitempty"`
}

// Data contract for Data of Delivery View
type DeliveryViewData struct {
	Id       *uuid.UUID `json:"id,omitempty"`
	Revision *int       `json:"revision,omitempty"`
	// Filter criteria status of the timeline
	CriteriaStatus *TimelineCriteriaStatus `json:"criteriaStatus,omitempty"`
	// The end date of the delivery view data
	EndDate *azuredevops.Time `json:"endDate,omitempty"`
	// Work item child id to parent id map
	ChildIdToParentIdMap *map[int]int `json:"childIdToParentIdMap,omitempty"`
	// Max number of teams that can be configured for a delivery plan
	MaxExpandedTeams *int `json:"maxExpandedTeams,omitempty"`
	// Mapping between parent id, title and all the child work item ids
	ParentItemMaps *[]ParentChildWIMap `json:"parentItemMaps,omitempty"`
	// The start date for the delivery view data
	StartDate *azuredevops.Time `json:"startDate,omitempty"`
	// All the team data
	Teams *[]TimelineTeamData `json:"teams,omitempty"`
	// List of all work item ids that have a dependency but not a violation
	WorkItemDependencies *[]int `json:"workItemDependencies,omitempty"`
	// List of all work item ids that have a violation
	WorkItemViolations *[]int `json:"workItemViolations,omitempty"`
}

// Collection of properties, specific to the DeliveryTimelineView
type DeliveryViewPropertyCollection struct {
	// Card settings
	CardSettings *CardSettings `json:"cardSettings,omitempty"`
	// Field criteria
	Criteria *[]FilterClause `json:"criteria,omitempty"`
	// Markers. Will be missing/null if there are no markers.
	Markers *[]Marker `json:"markers,omitempty"`
	// Card style settings
	StyleSettings *[]Rule `json:"styleSettings,omitempty"`
	// tag style settings
	TagStyleSettings *[]Rule `json:"tagStyleSettings,omitempty"`
	// Team backlog mappings
	TeamBacklogMappings *[]TeamBacklogMapping `json:"teamBacklogMappings,omitempty"`
}

// Object bag storing the set of permissions relevant to this plan
type FieldInfo struct {
	// The additional field display name
	DisplayName *string `json:"displayName,omitempty"`
	// The additional field type
	FieldType *FieldType `json:"fieldType,omitempty"`
	// Indicates if the field definition is for an identity field.
	IsIdentity *bool `json:"isIdentity,omitempty"`
	// The additional field reference name
	ReferenceName *string `json:"referenceName,omitempty"`
}

// An abstracted reference to a field
type FieldReference struct {
	// fieldRefName for the field
	ReferenceName *string `json:"referenceName,omitempty"`
	// Full http link to more information about the field
	Url *string `json:"url,omitempty"`
}

type FieldSetting struct {
}

type FieldType string

type fieldTypeValuesType struct {
	String    FieldType
	PlainText FieldType
	Integer   FieldType
	DateTime  FieldType
	TreePath  FieldType
	Boolean   FieldType
	Double    FieldType
}

var FieldTypeValues = fieldTypeValuesType{
	String:    "string",
	PlainText: "plainText",
	Integer:   "integer",
	DateTime:  "dateTime",
	TreePath:  "treePath",
	Boolean:   "boolean",
	Double:    "double",
}

type FilterClause struct {
	FieldName       *string `json:"fieldName,omitempty"`
	Index           *int    `json:"index,omitempty"`
	LogicalOperator *string `json:"logicalOperator,omitempty"`
	Operator        *string `json:"operator,omitempty"`
	Value           *string `json:"value,omitempty"`
}

type FilterGroup struct {
	End   *int `json:"end,omitempty"`
	Level *int `json:"level,omitempty"`
	Start *int `json:"start,omitempty"`
}

// Enum for the various modes of identity picker
type IdentityDisplayFormat string

type identityDisplayFormatValuesType struct {
	AvatarOnly        IdentityDisplayFormat
	FullName          IdentityDisplayFormat
	AvatarAndFullName IdentityDisplayFormat
}

var IdentityDisplayFormatValues = identityDisplayFormatValuesType{
	// Display avatar only
	AvatarOnly: "avatarOnly",
	// Display Full name only
	FullName: "fullName",
	// Display Avatar and Full name
	AvatarAndFullName: "avatarAndFullName",
}

type ITaskboardColumnMapping struct {
	State        *string `json:"state,omitempty"`
	WorkItemType *string `json:"workItemType,omitempty"`
}

// Capacity and teams for all teams in an iteration
type IterationCapacity struct {
	Teams                        *[]TeamCapacityTotals `json:"teams,omitempty"`
	TotalIterationCapacityPerDay *float64              `json:"totalIterationCapacityPerDay,omitempty"`
	TotalIterationDaysOff        *int                  `json:"totalIterationDaysOff,omitempty"`
}

// Represents work items in an iteration backlog
type IterationWorkItems struct {
	// Collection of links relevant to resource
	Links interface{} `json:"_links,omitempty"`
	// Full http link to the resource
	Url *string `json:"url,omitempty"`
	// Work item relations
	WorkItemRelations *[]workitemtracking.WorkItemLink `json:"workItemRelations,omitempty"`
}

// Client serialization contract for Delivery Timeline Markers.
type Marker struct {
	// Color associated with the marker.
	Color *string `json:"color,omitempty"`
	// Where the marker should be displayed on the timeline.
	Date *azuredevops.Time `json:"date,omitempty"`
	// Label/title for the marker.
	Label *string `json:"label,omitempty"`
}

type Member struct {
	DisplayName *string    `json:"displayName,omitempty"`
	Id          *uuid.UUID `json:"id,omitempty"`
	ImageUrl    *string    `json:"imageUrl,omitempty"`
	UniqueName  *string    `json:"uniqueName,omitempty"`
	Url         *string    `json:"url,omitempty"`
}

type ParentChildWIMap struct {
	ChildWorkItemIds *[]int  `json:"childWorkItemIds,omitempty"`
	Id               *int    `json:"id,omitempty"`
	Title            *string `json:"title,omitempty"`
	WorkItemTypeName *string `json:"workItemTypeName,omitempty"`
}

// Data contract for the plan definition
type Plan struct {
	// Identity that created this plan. Defaults to null for records before upgrading to ScaledAgileViewComponent4.
	CreatedByIdentity *webapi.IdentityRef `json:"createdByIdentity,omitempty"`
	// Date when the plan was created
	CreatedDate *azuredevops.Time `json:"createdDate,omitempty"`
	// Description of the plan
	Description *string `json:"description,omitempty"`
	// Id of the plan
	Id *uuid.UUID `json:"id,omitempty"`
	// Date when the plan was last accessed. Default is null.
	LastAccessed *azuredevops.Time `json:"lastAccessed,omitempty"`
	// Identity that last modified this plan. Defaults to null for records before upgrading to ScaledAgileViewComponent4.
	ModifiedByIdentity *webapi.IdentityRef `json:"modifiedByIdentity,omitempty"`
	// Date when the plan was last modified. Default to CreatedDate when the plan is first created.
	ModifiedDate *azuredevops.Time `json:"modifiedDate,omitempty"`
	// Name of the plan
	Name *string `json:"name,omitempty"`
	// The PlanPropertyCollection instance associated with the plan. These are dependent on the type of the plan. For example, DeliveryTimelineView, it would be of type DeliveryViewPropertyCollection.
	Properties interface{} `json:"properties,omitempty"`
	// Revision of the plan. Used to safeguard users from overwriting each other's changes.
	Revision *int `json:"revision,omitempty"`
	// Type of the plan
	Type *PlanType `json:"type,omitempty"`
	// The resource url to locate the plan via rest api
	Url *string `json:"url,omitempty"`
	// Bit flag indicating set of permissions a user has to the plan.
	UserPermissions *PlanUserPermissions `json:"userPermissions,omitempty"`
}

// Metadata about a plan definition that is stored in favorites service
type PlanMetadata struct {
	// Identity of the creator of the plan
	CreatedByIdentity *webapi.IdentityRef `json:"createdByIdentity,omitempty"`
	// Description of plan
	Description *string `json:"description,omitempty"`
	// Last modified date of the plan
	ModifiedDate *azuredevops.Time `json:"modifiedDate,omitempty"`
	// Bit flag indicating set of permissions a user has to the plan.
	UserPermissions *PlanUserPermissions `json:"userPermissions,omitempty"`
}

// Enum for the various types of plans
type PlanType string

type planTypeValuesType struct {
	DeliveryTimelineView PlanType
}

var PlanTypeValues = planTypeValuesType{
	DeliveryTimelineView: "deliveryTimelineView",
}

// [Flags] Flag for permissions a user can have for this plan.
type PlanUserPermissions string

type planUserPermissionsValuesType struct {
	None           PlanUserPermissions
	View           PlanUserPermissions
	Edit           PlanUserPermissions
	Delete         PlanUserPermissions
	Manage         PlanUserPermissions
	AllPermissions PlanUserPermissions
}

var PlanUserPermissionsValues = planUserPermissionsValuesType{
	// None
	None: "none",
	// Permission to view this plan.
	View: "view",
	// Permission to update this plan.
	Edit: "edit",
	// Permission to delete this plan.
	Delete: "delete",
	// Permission to manage this plan.
	Manage: "manage",
	// Full control permission for this plan.
	AllPermissions: "allPermissions",
}

// Base class for plan view data contracts. Anything common goes here.
type PlanViewData struct {
	Id       *uuid.UUID `json:"id,omitempty"`
	Revision *int       `json:"revision,omitempty"`
}

// Represents a single pre-defined query.
type PredefinedQuery struct {
	// Whether or not the query returned the complete set of data or if the data was truncated.
	HasMore *bool `json:"hasMore,omitempty"`
	// Id of the query
	Id *string `json:"id,omitempty"`
	// Localized name of the query
	Name *string `json:"name,omitempty"`
	// The results of the query.  This will be a set of WorkItem objects with only the 'id' set.  The client is responsible for paging in the data as needed.
	Results *[]workitemtracking.WorkItem `json:"results,omitempty"`
	// REST API Url to use to retrieve results for this query
	Url *string `json:"url,omitempty"`
	// Url to use to display a page in the browser with the results of this query
	WebUrl *string `json:"webUrl,omitempty"`
}

// Process Configurations for the project
type ProcessConfiguration struct {
	// Details about bug work items
	BugWorkItems *CategoryConfiguration `json:"bugWorkItems,omitempty"`
	// Details about portfolio backlogs
	PortfolioBacklogs *[]CategoryConfiguration `json:"portfolioBacklogs,omitempty"`
	// Details of requirement backlog
	RequirementBacklog *CategoryConfiguration `json:"requirementBacklog,omitempty"`
	// Details of task backlog
	TaskBacklog *CategoryConfiguration `json:"taskBacklog,omitempty"`
	// Type fields for the process configuration
	TypeFields *map[string]workitemtracking.WorkItemFieldReference `json:"typeFields,omitempty"`
	Url        *string                                             `json:"url,omitempty"`
}

// Represents a reorder request for one or more work items.
type ReorderOperation struct {
	// IDs of the work items to be reordered.  Must be valid WorkItem Ids.
	Ids *[]int `json:"ids,omitempty"`
	// IterationPath for reorder operation. This is only used when we reorder from the Iteration Backlog
	IterationPath *string `json:"iterationPath,omitempty"`
	// ID of the work item that should be after the reordered items. Can use 0 to specify the end of the list.
	NextId *int `json:"nextId,omitempty"`
	// Parent ID for all of the work items involved in this operation. Can use 0 to indicate the items don't have a parent.
	ParentId *int `json:"parentId,omitempty"`
	// ID of the work item that should be before the reordered items. Can use 0 to specify the beginning of the list.
	PreviousId *int `json:"previousId,omitempty"`
}

// Represents a reorder result for a work item.
type ReorderResult struct {
	// The ID of the work item that was reordered.
	Id *int `json:"id,omitempty"`
	// The updated order value of the work item that was reordered.
	Order *float64 `json:"order,omitempty"`
}

type Rule struct {
	Clauses   *[]FilterClause `json:"clauses,omitempty"`
	Filter    *string         `json:"filter,omitempty"`
	IsEnabled *string         `json:"isEnabled,omitempty"`
	Name      *string         `json:"name,omitempty"`
	Settings  *attribute      `json:"settings,omitempty"`
}

// Represents the taskbord column
type TaskboardColumn struct {
	// Column ID
	Id *uuid.UUID `json:"id,omitempty"`
	// Work item type states mapped to this column to support auto state update when column is updated.
	Mappings *[]ITaskboardColumnMapping `json:"mappings,omitempty"`
	// Column name
	Name *string `json:"name,omitempty"`
	// Column position relative to other columns in the same board
	Order *int `json:"order,omitempty"`
}

// Represents the state to column mapping per work item type This allows auto state update when the column changes
type TaskboardColumnMapping struct {
	// State of the work item type mapped to the column
	State *string `json:"state,omitempty"`
	// Work Item Type name who's state is mapped to the column
	WorkItemType *string `json:"workItemType,omitempty"`
}

type TaskboardColumns struct {
	Columns *[]TaskboardColumn `json:"columns,omitempty"`
	// Are the columns cutomized for this team
	IsCustomized *bool `json:"isCustomized,omitempty"`
	// Specifies if the referenced WIT and State is valid
	IsValid *bool `json:"isValid,omitempty"`
	// Details of validation failure if the state to column mapping is invalid
	ValidationMesssage *string `json:"validationMesssage,omitempty"`
}

// Column value of a work item in the taskboard
type TaskboardWorkItemColumn struct {
	// Work item column value in the taskboard
	Column *string `json:"column,omitempty"`
	// Work item column id in the taskboard
	ColumnId *uuid.UUID `json:"columnId,omitempty"`
	// Work Item state value
	State *string `json:"state,omitempty"`
	// Work item id
	WorkItemId *int `json:"workItemId,omitempty"`
}

// Mapping of teams to the corresponding work item category
type TeamBacklogMapping struct {
	CategoryReferenceName *string    `json:"categoryReferenceName,omitempty"`
	TeamId                *uuid.UUID `json:"teamId,omitempty"`
}

// Represents team member capacity with totals aggregated
type TeamCapacity struct {
	TeamMembers         *[]TeamMemberCapacityIdentityRef `json:"teamMembers,omitempty"`
	TotalCapacityPerDay *float64                         `json:"totalCapacityPerDay,omitempty"`
	TotalDaysOff        *int                             `json:"totalDaysOff,omitempty"`
}

// Team information with total capacity and days off
type TeamCapacityTotals struct {
	TeamCapacityPerDay *float64   `json:"teamCapacityPerDay,omitempty"`
	TeamId             *uuid.UUID `json:"teamId,omitempty"`
	TeamTotalDaysOff   *int       `json:"teamTotalDaysOff,omitempty"`
}

// Represents a single TeamFieldValue
type TeamFieldValue struct {
	IncludeChildren *bool   `json:"includeChildren,omitempty"`
	Value           *string `json:"value,omitempty"`
}

// Essentially a collection of team field values
type TeamFieldValues struct {
	// Collection of links relevant to resource
	Links interface{} `json:"_links,omitempty"`
	// Full http link to the resource
	Url *string `json:"url,omitempty"`
	// The default team field value
	DefaultValue *string `json:"defaultValue,omitempty"`
	// Shallow ref to the field being used as a team field
	Field *FieldReference `json:"field,omitempty"`
	// Collection of all valid team field values
	Values *[]TeamFieldValue `json:"values,omitempty"`
}

// Expected data from PATCH
type TeamFieldValuesPatch struct {
	DefaultValue *string           `json:"defaultValue,omitempty"`
	Values       *[]TeamFieldValue `json:"values,omitempty"`
}

type TeamIterationAttributes struct {
	// Finish date of the iteration. Date-only, correct unadjusted at midnight in UTC.
	FinishDate *azuredevops.Time `json:"finishDate,omitempty"`
	// Start date of the iteration. Date-only, correct unadjusted at midnight in UTC.
	StartDate *azuredevops.Time `json:"startDate,omitempty"`
	// Time frame of the iteration, such as past, current or future.
	TimeFrame *TimeFrame `json:"timeFrame,omitempty"`
}

// Represents capacity for a specific team member
type TeamMemberCapacity struct {
	// Collection of links relevant to resource
	Links interface{} `json:"_links,omitempty"`
	// Full http link to the resource
	Url *string `json:"url,omitempty"`
	// Collection of capacities associated with the team member
	Activities *[]Activity `json:"activities,omitempty"`
	// The days off associated with the team member
	DaysOff *[]DateRange `json:"daysOff,omitempty"`
	// Shallow Ref to the associated team member
	TeamMember *Member `json:"teamMember,omitempty"`
}

// Represents capacity for a specific team member
type TeamMemberCapacityIdentityRef struct {
	// Collection of links relevant to resource
	Links interface{} `json:"_links,omitempty"`
	// Full http link to the resource
	Url *string `json:"url,omitempty"`
	// Collection of capacities associated with the team member
	Activities *[]Activity `json:"activities,omitempty"`
	// The days off associated with the team member
	DaysOff *[]DateRange `json:"daysOff,omitempty"`
	// Identity ref of the associated team member
	TeamMember *webapi.IdentityRef `json:"teamMember,omitempty"`
}

// Data contract for TeamSettings
type TeamSetting struct {
	// Collection of links relevant to resource
	Links interface{} `json:"_links,omitempty"`
	// Full http link to the resource
	Url *string `json:"url,omitempty"`
	// Backlog Iteration
	BacklogIteration *TeamSettingsIteration `json:"backlogIteration,omitempty"`
	// Information about categories that are visible on the backlog.
	BacklogVisibilities *map[string]bool `json:"backlogVisibilities,omitempty"`
	// BugsBehavior (Off, AsTasks, AsRequirements, ...)
	BugsBehavior *BugsBehavior `json:"bugsBehavior,omitempty"`
	// Default Iteration, the iteration used when creating a new work item on the queries page.
	DefaultIteration *TeamSettingsIteration `json:"defaultIteration,omitempty"`
	// Default Iteration macro (if any)
	DefaultIterationMacro *string `json:"defaultIterationMacro,omitempty"`
	// Days that the team is working
	WorkingDays *[]string `json:"workingDays,omitempty"`
}

// Base class for TeamSettings data contracts. Anything common goes here.
type TeamSettingsDataContractBase struct {
	// Collection of links relevant to resource
	Links interface{} `json:"_links,omitempty"`
	// Full http link to the resource
	Url *string `json:"url,omitempty"`
}

type TeamSettingsDaysOff struct {
	// Collection of links relevant to resource
	Links interface{} `json:"_links,omitempty"`
	// Full http link to the resource
	Url     *string      `json:"url,omitempty"`
	DaysOff *[]DateRange `json:"daysOff,omitempty"`
}

type TeamSettingsDaysOffPatch struct {
	DaysOff *[]DateRange `json:"daysOff,omitempty"`
}

// Represents a shallow ref for a single iteration.
type TeamSettingsIteration struct {
	// Collection of links relevant to resource
	Links interface{} `json:"_links,omitempty"`
	// Full http link to the resource
	Url *string `json:"url,omitempty"`
	// Attributes of the iteration such as start and end date.
	Attributes *TeamIterationAttributes `json:"attributes,omitempty"`
	// Id of the iteration.
	Id *uuid.UUID `json:"id,omitempty"`
	// Name of the iteration.
	Name *string `json:"name,omitempty"`
	// Relative path of the iteration.
	Path *string `json:"path,omitempty"`
}

// Data contract for what we expect to receive when PATCH
type TeamSettingsPatch struct {
	BacklogIteration      *uuid.UUID       `json:"backlogIteration,omitempty"`
	BacklogVisibilities   *map[string]bool `json:"backlogVisibilities,omitempty"`
	BugsBehavior          *BugsBehavior    `json:"bugsBehavior,omitempty"`
	DefaultIteration      *uuid.UUID       `json:"defaultIteration,omitempty"`
	DefaultIterationMacro *string          `json:"defaultIterationMacro,omitempty"`
	WorkingDays           *[]string        `json:"workingDays,omitempty"`
}

type TimeFrame string

type timeFrameValuesType struct {
	Past    TimeFrame
	Current TimeFrame
	Future  TimeFrame
}

var TimeFrameValues = timeFrameValuesType{
	Past:    "past",
	Current: "current",
	Future:  "future",
}

type TimelineCriteriaStatus struct {
	Message *string                     `json:"message,omitempty"`
	Type    *TimelineCriteriaStatusCode `json:"type,omitempty"`
}

type TimelineCriteriaStatusCode string

type timelineCriteriaStatusCodeValuesType struct {
	Ok                  TimelineCriteriaStatusCode
	InvalidFilterClause TimelineCriteriaStatusCode
	Unknown             TimelineCriteriaStatusCode
}

var TimelineCriteriaStatusCodeValues = timelineCriteriaStatusCodeValuesType{
	// No error - filter is good.
	Ok: "ok",
	// One of the filter clause is invalid.
	InvalidFilterClause: "invalidFilterClause",
	// Unknown error.
	Unknown: "unknown",
}

type TimelineIterationStatus struct {
	Message *string                      `json:"message,omitempty"`
	Type    *TimelineIterationStatusCode `json:"type,omitempty"`
}

type TimelineIterationStatusCode string

type timelineIterationStatusCodeValuesType struct {
	Ok            TimelineIterationStatusCode
	IsOverlapping TimelineIterationStatusCode
}

var TimelineIterationStatusCodeValues = timelineIterationStatusCodeValuesType{
	// No error - iteration data is good.
	Ok: "ok",
	// This iteration overlaps with another iteration, no data is returned for this iteration.
	IsOverlapping: "isOverlapping",
}

type TimelineTeamData struct {
	// Backlog matching the mapped backlog associated with this team.
	Backlog *BacklogLevel `json:"backlog,omitempty"`
	// The field reference names of the work item data
	FieldReferenceNames *[]string `json:"fieldReferenceNames,omitempty"`
	// The id of the team
	Id *uuid.UUID `json:"id,omitempty"`
	// Was iteration and work item data retrieved for this team. <remarks> Teams with IsExpanded false have not had their iteration, work item, and field related data queried and will never contain this data. If true then these items are queried and, if there are items in the queried range, there will be data. </remarks>
	IsExpanded *bool `json:"isExpanded,omitempty"`
	// The iteration data, including the work items, in the queried date range.
	Iterations *[]TimelineTeamIteration `json:"iterations,omitempty"`
	// The name of the team
	Name *string `json:"name,omitempty"`
	// The order by field name of this team
	OrderByField *string `json:"orderByField,omitempty"`
	// The field reference names of the partially paged work items, such as ID, WorkItemType
	PartiallyPagedFieldReferenceNames *[]string        `json:"partiallyPagedFieldReferenceNames,omitempty"`
	PartiallyPagedWorkItems           *[][]interface{} `json:"partiallyPagedWorkItems,omitempty"`
	// The project id the team belongs team
	ProjectId *uuid.UUID `json:"projectId,omitempty"`
	// Work item types for which we will collect roll up data on the client side
	RollupWorkItemTypes *[]string `json:"rollupWorkItemTypes,omitempty"`
	// Status for this team.
	Status *TimelineTeamStatus `json:"status,omitempty"`
	// The team field default value
	TeamFieldDefaultValue *string `json:"teamFieldDefaultValue,omitempty"`
	// The team field name of this team
	TeamFieldName *string `json:"teamFieldName,omitempty"`
	// The team field values
	TeamFieldValues *[]TeamFieldValue `json:"teamFieldValues,omitempty"`
	// Work items associated with the team that are not under any of the team's iterations
	WorkItems *[][]interface{} `json:"workItems,omitempty"`
	// Colors for the work item types.
	WorkItemTypeColors *[]WorkItemColor `json:"workItemTypeColors,omitempty"`
}

type TimelineTeamIteration struct {
	// The iteration CSS Node Id
	CssNodeId *string `json:"cssNodeId,omitempty"`
	// The end date of the iteration
	FinishDate *azuredevops.Time `json:"finishDate,omitempty"`
	// The iteration name
	Name *string `json:"name,omitempty"`
	// All the partially paged workitems in this iteration.
	PartiallyPagedWorkItems *[][]interface{} `json:"partiallyPagedWorkItems,omitempty"`
	// The iteration path
	Path *string `json:"path,omitempty"`
	// The start date of the iteration
	StartDate *azuredevops.Time `json:"startDate,omitempty"`
	// The status of this iteration
	Status *TimelineIterationStatus `json:"status,omitempty"`
	// The work items that have been paged in this iteration
	WorkItems *[][]interface{} `json:"workItems,omitempty"`
}

type TimelineTeamStatus struct {
	Message *string                 `json:"message,omitempty"`
	Type    *TimelineTeamStatusCode `json:"type,omitempty"`
}

type TimelineTeamStatusCode string

type timelineTeamStatusCodeValuesType struct {
	Ok                        TimelineTeamStatusCode
	DoesntExistOrAccessDenied TimelineTeamStatusCode
	MaxTeamsExceeded          TimelineTeamStatusCode
	MaxTeamFieldsExceeded     TimelineTeamStatusCode
	BacklogInError            TimelineTeamStatusCode
	MissingTeamFieldValue     TimelineTeamStatusCode
	NoIterationsExist         TimelineTeamStatusCode
}

var TimelineTeamStatusCodeValues = timelineTeamStatusCodeValuesType{
	// No error - all data for team is good.
	Ok: "ok",
	// Team does not exist or access is denied.
	DoesntExistOrAccessDenied: "doesntExistOrAccessDenied",
	// Maximum number of teams was exceeded. No team data will be returned for this team.
	MaxTeamsExceeded: "maxTeamsExceeded",
	// Maximum number of team fields (ie Area paths) have been exceeded. No team data will be returned for this team.
	MaxTeamFieldsExceeded: "maxTeamFieldsExceeded",
	// Backlog does not exist or is missing crucial information.
	BacklogInError: "backlogInError",
	// Team field value is not set for this team. No team data will be returned for this team
	MissingTeamFieldValue: "missingTeamFieldValue",
	// Team does not have a single iteration with date range.
	NoIterationsExist: "noIterationsExist",
}

type UpdatePlan struct {
	// Description of the plan
	Description *string `json:"description,omitempty"`
	// Name of the plan to create.
	Name *string `json:"name,omitempty"`
	// Plan properties.
	Properties interface{} `json:"properties,omitempty"`
	// Revision of the plan that was updated - the value used here should match the one the server gave the client in the Plan.
	Revision *int `json:"revision,omitempty"`
	// Type of the plan
	Type *PlanType `json:"type,omitempty"`
}

type UpdateTaskboardColumn struct {
	// Column ID, keep it null for new column
	Id *uuid.UUID `json:"id,omitempty"`
	// Work item type states mapped to this column to support auto state update when column is updated.
	Mappings *[]TaskboardColumnMapping `json:"mappings,omitempty"`
	// Column name is required
	Name *string `json:"name,omitempty"`
	// Column position relative to other columns in the same board
	Order *int `json:"order,omitempty"`
}

type UpdateTaskboardWorkItemColumn struct {
	NewColumn *string `json:"newColumn,omitempty"`
}

// Work item color and icon.
type WorkItemColor struct {
	Icon             *string `json:"icon,omitempty"`
	PrimaryColor     *string `json:"primaryColor,omitempty"`
	WorkItemTypeName *string `json:"workItemTypeName,omitempty"`
}

type WorkItemTypeStateInfo struct {
	// State name to state category map
	States *map[string]string `json:"states,omitempty"`
	// Work Item type name
	WorkItemTypeName *string `json:"workItemTypeName,omitempty"`
}

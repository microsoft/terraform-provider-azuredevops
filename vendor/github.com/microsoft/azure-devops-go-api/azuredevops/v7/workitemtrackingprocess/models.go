// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package workitemtrackingprocess

import (
	"github.com/google/uuid"
)

// Class that describes a request to add a field in a work item type.
type AddProcessWorkItemTypeFieldRequest struct {
	// The list of field allowed values.
	AllowedValues *[]string `json:"allowedValues,omitempty"`
	// Allow setting field value to a group identity. Only applies to identity fields.
	AllowGroups *bool `json:"allowGroups,omitempty"`
	// The default value of the field.
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	// If true the field cannot be edited.
	ReadOnly *bool `json:"readOnly,omitempty"`
	// Reference name of the field.
	ReferenceName *string `json:"referenceName,omitempty"`
	// If true the field cannot be empty.
	Required *bool `json:"required,omitempty"`
}

// Represent a control in the form.
type Control struct {
	// Contribution for the control.
	Contribution *WitContribution `json:"contribution,omitempty"`
	// Type of the control.
	ControlType *string `json:"controlType,omitempty"`
	// Height of the control, for html controls.
	Height *int `json:"height,omitempty"`
	// The id for the layout node.
	Id *string `json:"id,omitempty"`
	// A value indicating whether this layout node has been inherited. from a parent layout.  This is expected to only be only set by the combiner.
	Inherited *bool `json:"inherited,omitempty"`
	// A value indicating if the layout node is contribution or not.
	IsContribution *bool `json:"isContribution,omitempty"`
	// Label for the field.
	Label *string `json:"label,omitempty"`
	// Inner text of the control.
	Metadata *string `json:"metadata,omitempty"`
	// Order in which the control should appear in its group.
	Order *int `json:"order,omitempty"`
	// A value indicating whether this layout node has been overridden . by a child layout.
	Overridden *bool `json:"overridden,omitempty"`
	// A value indicating if the control is readonly.
	ReadOnly *bool `json:"readOnly,omitempty"`
	// A value indicating if the control should be hidden or not.
	Visible *bool `json:"visible,omitempty"`
	// Watermark text for the textbox.
	Watermark *string `json:"watermark,omitempty"`
}

// Describes a process being created.
type CreateProcessModel struct {
	// Description of the process
	Description *string `json:"description,omitempty"`
	// Name of the process
	Name *string `json:"name,omitempty"`
	// The ID of the parent process
	ParentProcessTypeId *uuid.UUID `json:"parentProcessTypeId,omitempty"`
	// Reference name of process being created. If not specified, server will assign a unique reference name
	ReferenceName *string `json:"referenceName,omitempty"`
}

// Request object/class for creating a rule on a work item type.
type CreateProcessRuleRequest struct {
	// List of actions to take when the rule is triggered.
	Actions *[]RuleAction `json:"actions,omitempty"`
	// List of conditions when the rule should be triggered.
	Conditions *[]RuleCondition `json:"conditions,omitempty"`
	// Indicates if the rule is disabled.
	IsDisabled *bool `json:"isDisabled,omitempty"`
	// Name for the rule.
	Name *string `json:"name,omitempty"`
}

// Class for create work item type request
type CreateProcessWorkItemTypeRequest struct {
	// Color hexadecimal code to represent the work item type
	Color *string `json:"color,omitempty"`
	// Description of the work item type
	Description *string `json:"description,omitempty"`
	// Icon to represent the work item type
	Icon *string `json:"icon,omitempty"`
	// Parent work item type for work item type
	InheritsFrom *string `json:"inheritsFrom,omitempty"`
	// True if the work item type need to be disabled
	IsDisabled *bool `json:"isDisabled,omitempty"`
	// Name of work item type
	Name *string `json:"name,omitempty"`
}

// Indicates the customization-type. Customization-type is System if is system generated or by default. Customization-type is Inherited if the existing workitemtype of inherited process is customized. Customization-type is Custom if the newly created workitemtype is customized.
type CustomizationType string

type customizationTypeValuesType struct {
	System    CustomizationType
	Inherited CustomizationType
	Custom    CustomizationType
}

var CustomizationTypeValues = customizationTypeValuesType{
	// Customization-type is System if is system generated workitemtype.
	System: "system",
	// Customization-type is Inherited if the existing workitemtype of inherited process is customized.
	Inherited: "inherited",
	// Customization-type is Custom if the newly created workitemtype is customized.
	Custom: "custom",
}

// Represents the extensions part of the layout
type Extension struct {
	// Id of the extension
	Id *string `json:"id,omitempty"`
}

type FieldModel struct {
	Description *string    `json:"description,omitempty"`
	Id          *string    `json:"id,omitempty"`
	IsIdentity  *bool      `json:"isIdentity,omitempty"`
	IsLocked    *bool      `json:"isLocked,omitempty"`
	Name        *string    `json:"name,omitempty"`
	Type        *FieldType `json:"type,omitempty"`
	Url         *string    `json:"url,omitempty"`
}

type FieldRuleModel struct {
	Actions      *[]RuleActionModel    `json:"actions,omitempty"`
	Conditions   *[]RuleConditionModel `json:"conditions,omitempty"`
	FriendlyName *string               `json:"friendlyName,omitempty"`
	Id           *uuid.UUID            `json:"id,omitempty"`
	IsDisabled   *bool                 `json:"isDisabled,omitempty"`
	IsSystem     *bool                 `json:"isSystem,omitempty"`
}

// Enum for the type of a field.
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
	PicklistInteger FieldType
	PicklistString  FieldType
	PicklistDouble  FieldType
}

var FieldTypeValues = fieldTypeValuesType{
	// String field type.
	String: "string",
	// Integer field type.
	Integer: "integer",
	// DateTime field type.
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
	// Integer picklist field type.
	PicklistInteger: "picklistInteger",
	// String picklist field type.
	PicklistString: "picklistString",
	// Double picklist field type.
	PicklistDouble: "picklistDouble",
}

// Describes the layout of a work item type
type FormLayout struct {
	// Gets and sets extensions list.
	Extensions *[]Extension `json:"extensions,omitempty"`
	// Top level tabs of the layout.
	Pages *[]Page `json:"pages,omitempty"`
	// Headers controls of the layout.
	SystemControls *[]Control `json:"systemControls,omitempty"`
}

// Expand options to fetch fields for behaviors API.
type GetBehaviorsExpand string

type getBehaviorsExpandValuesType struct {
	None           GetBehaviorsExpand
	Fields         GetBehaviorsExpand
	CombinedFields GetBehaviorsExpand
}

var GetBehaviorsExpandValues = getBehaviorsExpandValuesType{
	// Default none option.
	None: "none",
	// This option returns fields associated with a behavior.
	Fields: "fields",
	// This option returns fields associated with this behavior and all behaviors from which it inherits.
	CombinedFields: "combinedFields",
}

// [Flags] The expand level of returned processes.
type GetProcessExpandLevel string

type getProcessExpandLevelValuesType struct {
	None     GetProcessExpandLevel
	Projects GetProcessExpandLevel
}

var GetProcessExpandLevelValues = getProcessExpandLevelValuesType{
	// No expand level.
	None: "none",
	// Projects expand level.
	Projects: "projects",
}

// [Flags] Flag to define what properties to return in get work item type response.
type GetWorkItemTypeExpand string

type getWorkItemTypeExpandValuesType struct {
	None      GetWorkItemTypeExpand
	States    GetWorkItemTypeExpand
	Behaviors GetWorkItemTypeExpand
	Layout    GetWorkItemTypeExpand
}

var GetWorkItemTypeExpandValues = getWorkItemTypeExpandValuesType{
	// Returns no properties in get work item type response.
	None: "none",
	// Returns states property in get work item type response.
	States: "states",
	// Returns behaviors property in get work item type response.
	Behaviors: "behaviors",
	// Returns layout property in get work item type response.
	Layout: "layout",
}

// Represent a group in the form that holds controls in it.
type Group struct {
	// Contribution for the group.
	Contribution *WitContribution `json:"contribution,omitempty"`
	// Controls to be put in the group.
	Controls *[]Control `json:"controls,omitempty"`
	// The height for the contribution.
	Height *int `json:"height,omitempty"`
	// The id for the layout node.
	Id *string `json:"id,omitempty"`
	// A value indicating whether this layout node has been inherited from a parent layout.  This is expected to only be only set by the combiner.
	Inherited *bool `json:"inherited,omitempty"`
	// A value indicating if the layout node is contribution are not.
	IsContribution *bool `json:"isContribution,omitempty"`
	// Label for the group.
	Label *string `json:"label,omitempty"`
	// Order in which the group should appear in the section.
	Order *int `json:"order,omitempty"`
	// A value indicating whether this layout node has been overridden by a child layout.
	Overridden *bool `json:"overridden,omitempty"`
	// A value indicating if the group should be hidden or not.
	Visible *bool `json:"visible,omitempty"`
}

// Class that describes the work item state is hidden.
type HideStateModel struct {
	// Returns 'true', if workitem state is hidden, 'false' otherwise.
	Hidden *bool `json:"hidden,omitempty"`
}

// Describes a page in the work item form layout
type Page struct {
	// Contribution for the page.
	Contribution *WitContribution `json:"contribution,omitempty"`
	// The id for the layout node.
	Id *string `json:"id,omitempty"`
	// A value indicating whether this layout node has been inherited from a parent layout.  This is expected to only be only set by the combiner.
	Inherited *bool `json:"inherited,omitempty"`
	// A value indicating if the layout node is contribution are not.
	IsContribution *bool `json:"isContribution,omitempty"`
	// The label for the page.
	Label *string `json:"label,omitempty"`
	// A value indicating whether any user operations are permitted on this page and the contents of this page
	Locked *bool `json:"locked,omitempty"`
	// Order in which the page should appear in the layout.
	Order *int `json:"order,omitempty"`
	// A value indicating whether this layout node has been overridden by a child layout.
	Overridden *bool `json:"overridden,omitempty"`
	// The icon for the page.
	PageType *PageType `json:"pageType,omitempty"`
	// The sections of the page.
	Sections *[]Section `json:"sections,omitempty"`
	// A value indicating if the page should be hidden or not.
	Visible *bool `json:"visible,omitempty"`
}

// Enum for the types of pages in the work item form layout
type PageType string

type pageTypeValuesType struct {
	Custom      PageType
	History     PageType
	Links       PageType
	Attachments PageType
}

var PageTypeValues = pageTypeValuesType{
	// Custom page type.
	Custom: "custom",
	// History page type.
	History: "history",
	// Link page type.
	Links: "links",
	// Attachment page type.
	Attachments: "attachments",
}

// Picklist.
type PickList struct {
	// ID of the picklist
	Id *uuid.UUID `json:"id,omitempty"`
	// Indicates whether items outside of suggested list are allowed
	IsSuggested *bool `json:"isSuggested,omitempty"`
	// Name of the picklist
	Name *string `json:"name,omitempty"`
	// DataType of picklist
	Type *string `json:"type,omitempty"`
	// Url of the picklist
	Url *string `json:"url,omitempty"`
	// A list of PicklistItemModel.
	Items *[]string `json:"items,omitempty"`
}

// Metadata for picklist.
type PickListMetadata struct {
	// ID of the picklist
	Id *uuid.UUID `json:"id,omitempty"`
	// Indicates whether items outside of suggested list are allowed
	IsSuggested *bool `json:"isSuggested,omitempty"`
	// Name of the picklist
	Name *string `json:"name,omitempty"`
	// DataType of picklist
	Type *string `json:"type,omitempty"`
	// Url of the picklist
	Url *string `json:"url,omitempty"`
}

// Process Behavior Model.
type ProcessBehavior struct {
	// Color.
	Color *string `json:"color,omitempty"`
	// Indicates the type of customization on this work item. System behaviors are inherited from parent process but not modified. Inherited behaviors are modified behaviors that were inherited from parent process. Custom behaviors are behaviors created by user in current process.
	Customization *CustomizationType `json:"customization,omitempty"`
	// . Description
	Description *string `json:"description,omitempty"`
	// Process Behavior Fields.
	Fields *[]ProcessBehaviorField `json:"fields,omitempty"`
	// Parent behavior reference.
	Inherits *ProcessBehaviorReference `json:"inherits,omitempty"`
	// Behavior Name.
	Name *string `json:"name,omitempty"`
	// Rank of the behavior
	Rank *int `json:"rank,omitempty"`
	// Behavior Id
	ReferenceName *string `json:"referenceName,omitempty"`
	// Url of the behavior.
	Url *string `json:"url,omitempty"`
}

// Process Behavior Create Payload.
type ProcessBehaviorCreateRequest struct {
	// Color.
	Color *string `json:"color,omitempty"`
	// Parent behavior id.
	Inherits *string `json:"inherits,omitempty"`
	// Name of the behavior.
	Name *string `json:"name,omitempty"`
	// ReferenceName is optional, if not specified will be auto-generated.
	ReferenceName *string `json:"referenceName,omitempty"`
}

// Process Behavior Field.
type ProcessBehaviorField struct {
	// Name of the field.
	Name *string `json:"name,omitempty"`
	// Reference name of the field.
	ReferenceName *string `json:"referenceName,omitempty"`
	// Url to field.
	Url *string `json:"url,omitempty"`
}

// Process behavior Reference.
type ProcessBehaviorReference struct {
	// Id of a Behavior.
	BehaviorRefName *string `json:"behaviorRefName,omitempty"`
	// Url to behavior.
	Url *string `json:"url,omitempty"`
}

// Process Behavior Replace Payload.
type ProcessBehaviorUpdateRequest struct {
	// Color.
	Color *string `json:"color,omitempty"`
	// Behavior Name.
	Name *string `json:"name,omitempty"`
}

type ProcessClass string

type processClassValuesType struct {
	System  ProcessClass
	Derived ProcessClass
	Custom  ProcessClass
}

var ProcessClassValues = processClassValuesType{
	System:  "system",
	Derived: "derived",
	Custom:  "custom",
}

// Process.
type ProcessInfo struct {
	// Indicates the type of customization on this process. System Process is default process. Inherited Process is modified process that was System process before.
	CustomizationType *CustomizationType `json:"customizationType,omitempty"`
	// Description of the process.
	Description *string `json:"description,omitempty"`
	// Is the process default.
	IsDefault *bool `json:"isDefault,omitempty"`
	// Is the process enabled.
	IsEnabled *bool `json:"isEnabled,omitempty"`
	// Name of the process.
	Name *string `json:"name,omitempty"`
	// ID of the parent process.
	ParentProcessTypeId *uuid.UUID `json:"parentProcessTypeId,omitempty"`
	// Projects in this process to which the user is subscribed to.
	Projects *[]ProjectReference `json:"projects,omitempty"`
	// Reference name of the process.
	ReferenceName *string `json:"referenceName,omitempty"`
	// The ID of the process.
	TypeId *uuid.UUID `json:"typeId,omitempty"`
}

type ProcessModel struct {
	// Description of the process
	Description *string `json:"description,omitempty"`
	// Name of the process
	Name *string `json:"name,omitempty"`
	// Projects in this process
	Projects *[]ProjectReference `json:"projects,omitempty"`
	// Properties of the process
	Properties *ProcessProperties `json:"properties,omitempty"`
	// Reference name of the process
	ReferenceName *string `json:"referenceName,omitempty"`
	// The ID of the process
	TypeId *uuid.UUID `json:"typeId,omitempty"`
}

// Properties of the process.
type ProcessProperties struct {
	// Class of the process.
	Class *ProcessClass `json:"class,omitempty"`
	// Is the process default process.
	IsDefault *bool `json:"isDefault,omitempty"`
	// Is the process enabled.
	IsEnabled *bool `json:"isEnabled,omitempty"`
	// ID of the parent process.
	ParentProcessTypeId *uuid.UUID `json:"parentProcessTypeId,omitempty"`
	// Version of the process.
	Version *string `json:"version,omitempty"`
}

// Process Rule Response.
type ProcessRule struct {
	// List of actions to take when the rule is triggered.
	Actions *[]RuleAction `json:"actions,omitempty"`
	// List of conditions when the rule should be triggered.
	Conditions *[]RuleCondition `json:"conditions,omitempty"`
	// Indicates if the rule is disabled.
	IsDisabled *bool `json:"isDisabled,omitempty"`
	// Name for the rule.
	Name *string `json:"name,omitempty"`
	// Indicates if the rule is system generated or created by user.
	CustomizationType *CustomizationType `json:"customizationType,omitempty"`
	// Id to uniquely identify the rule.
	Id *uuid.UUID `json:"id,omitempty"`
	// Resource Url.
	Url *string `json:"url,omitempty"`
}

// Class that describes a work item type object
type ProcessWorkItemType struct {
	Behaviors *[]WorkItemTypeBehavior `json:"behaviors,omitempty"`
	// Color hexadecimal code to represent the work item type
	Color *string `json:"color,omitempty"`
	// Indicates the type of customization on this work item System work item types are inherited from parent process but not modified Inherited work item types are modified work item that were inherited from parent process Custom work item types are work item types that were created in the current process
	Customization *CustomizationType `json:"customization,omitempty"`
	// Description of the work item type
	Description *string `json:"description,omitempty"`
	// Icon to represent the work item typ
	Icon *string `json:"icon,omitempty"`
	// Reference name of the parent work item type
	Inherits *string `json:"inherits,omitempty"`
	// Indicates if a work item type is disabled
	IsDisabled *bool       `json:"isDisabled,omitempty"`
	Layout     *FormLayout `json:"layout,omitempty"`
	// Name of the work item type
	Name *string `json:"name,omitempty"`
	// Reference name of work item type
	ReferenceName *string                     `json:"referenceName,omitempty"`
	States        *[]WorkItemStateResultModel `json:"states,omitempty"`
	// Url of the work item type
	Url *string `json:"url,omitempty"`
}

// Class that describes a field in a work item type and its properties.
type ProcessWorkItemTypeField struct {
	// The list of field allowed values.
	AllowedValues *[]interface{} `json:"allowedValues,omitempty"`
	// Allow setting field value to a group identity. Only applies to identity fields.
	AllowGroups *bool `json:"allowGroups,omitempty"`
	// Indicates the type of customization on this work item.
	Customization *CustomizationType `json:"customization,omitempty"`
	// The default value of the field.
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	// Description of the field.
	Description *string `json:"description,omitempty"`
	// Information about field definition being locked for editing
	IsLocked *bool `json:"isLocked,omitempty"`
	// Name of the field.
	Name *string `json:"name,omitempty"`
	// If true the field cannot be edited.
	ReadOnly *bool `json:"readOnly,omitempty"`
	// Reference name of the field.
	ReferenceName *string `json:"referenceName,omitempty"`
	// If true the field cannot be empty.
	Required *bool `json:"required,omitempty"`
	// Type of the field.
	Type *FieldType `json:"type,omitempty"`
	// Resource URL of the field.
	Url *string `json:"url,omitempty"`
}

// Expand options for the work item field(s) request.
type ProcessWorkItemTypeFieldsExpandLevel string

type processWorkItemTypeFieldsExpandLevelValuesType struct {
	None          ProcessWorkItemTypeFieldsExpandLevel
	AllowedValues ProcessWorkItemTypeFieldsExpandLevel
	All           ProcessWorkItemTypeFieldsExpandLevel
}

var ProcessWorkItemTypeFieldsExpandLevelValues = processWorkItemTypeFieldsExpandLevelValuesType{
	// Includes only basic properties of the field.
	None: "none",
	// Includes allowed values for the field.
	AllowedValues: "allowedValues",
	// Includes allowed values and dependent fields of the field.
	All: "all",
}

// Defines the project reference class.
type ProjectReference struct {
	// Description of the project
	Description *string `json:"description,omitempty"`
	// The ID of the project
	Id *uuid.UUID `json:"id,omitempty"`
	// Name of the project
	Name *string `json:"name,omitempty"`
	// Url of the project
	Url *string `json:"url,omitempty"`
}

// Action to take when the rule is triggered.
type RuleAction struct {
	// Type of action to take when the rule is triggered.
	ActionType *RuleActionType `json:"actionType,omitempty"`
	// Field on which the action should be taken.
	TargetField *string `json:"targetField,omitempty"`
	// Value to apply on target field, once the action is taken.
	Value *string `json:"value,omitempty"`
}

// Action to take when the rule is triggered.
type RuleActionModel struct {
	ActionType  *string `json:"actionType,omitempty"`
	TargetField *string `json:"targetField,omitempty"`
	Value       *string `json:"value,omitempty"`
}

// Type of action to take when the rule is triggered.
type RuleActionType string

type ruleActionTypeValuesType struct {
	MakeRequired              RuleActionType
	MakeReadOnly              RuleActionType
	SetDefaultValue           RuleActionType
	SetDefaultFromClock       RuleActionType
	SetDefaultFromCurrentUser RuleActionType
	SetDefaultFromField       RuleActionType
	CopyValue                 RuleActionType
	CopyFromClock             RuleActionType
	CopyFromCurrentUser       RuleActionType
	CopyFromField             RuleActionType
	SetValueToEmpty           RuleActionType
	CopyFromServerClock       RuleActionType
	CopyFromServerCurrentUser RuleActionType
	HideTargetField           RuleActionType
	DisallowValue             RuleActionType
}

var RuleActionTypeValues = ruleActionTypeValuesType{
	// Make the target field required. Example : {"actionType":"$makeRequired","targetField":"Microsoft.VSTS.Common.Activity","value":""}
	MakeRequired: "makeRequired",
	// Make the target field read-only. Example : {"actionType":"$makeReadOnly","targetField":"Microsoft.VSTS.Common.Activity","value":""}
	MakeReadOnly: "makeReadOnly",
	// Set a default value on the target field. This is used if the user creates a integer/string field and sets a default value of this field.
	SetDefaultValue: "setDefaultValue",
	// Set the default value on the target field from server clock. This is used if user creates the field like Date/Time and uses default value.
	SetDefaultFromClock: "setDefaultFromClock",
	// Set the default current user value on the target field. This is used if the user creates the field of type identity and uses default value.
	SetDefaultFromCurrentUser: "setDefaultFromCurrentUser",
	// Set the default value on from existing field to the target field.  This used wants to set a existing field value to the current field.
	SetDefaultFromField: "setDefaultFromField",
	// Set the value of target field to given value. Example : {actionType: "$copyValue", targetField: "ScrumInherited.mypicklist", value: "samplevalue"}
	CopyValue: "copyValue",
	// Set the value from clock.
	CopyFromClock: "copyFromClock",
	// Set the current user to the target field. Example : {"actionType":"$copyFromCurrentUser","targetField":"System.AssignedTo","value":""}.
	CopyFromCurrentUser: "copyFromCurrentUser",
	// Copy the value from a specified field and set to target field. Example : {actionType: "$copyFromField", targetField: "System.AssignedTo", value:"System.ChangedBy"}. Here, value is copied from "System.ChangedBy" and set to "System.AssingedTo" field.
	CopyFromField: "copyFromField",
	// Set the value of the target field to empty.
	SetValueToEmpty: "setValueToEmpty",
	// Use the current time to set the value of the target field. Example : {actionType: "$copyFromServerClock", targetField: "System.CreatedDate", value: ""}
	CopyFromServerClock: "copyFromServerClock",
	// Use the current user to set the value of the target field.
	CopyFromServerCurrentUser: "copyFromServerCurrentUser",
	// Hides target field from the form. This is a server side only action.
	HideTargetField: "hideTargetField",
	// Disallows a field from being set to a specific value.
	DisallowValue: "disallowValue",
}

// Defines a condition on a field when the rule should be triggered.
type RuleCondition struct {
	// Type of condition. $When. This condition limits the execution of its children to cases when another field has a particular value, i.e. when the Is value of the referenced field is equal to the given literal value. $WhenNot.This condition limits the execution of its children to cases when another field does not have a particular value, i.e.when the Is value of the referenced field is not equal to the given literal value. $WhenChanged.This condition limits the execution of its children to cases when another field has changed, i.e.when the Is value of the referenced field is not equal to the Was value of that field. $WhenNotChanged.This condition limits the execution of its children to cases when another field has not changed, i.e.when the Is value of the referenced field is equal to the Was value of that field.
	ConditionType *RuleConditionType `json:"conditionType,omitempty"`
	// Field that defines condition.
	Field *string `json:"field,omitempty"`
	// Value of field to define the condition for rule.
	Value *string `json:"value,omitempty"`
}

type RuleConditionModel struct {
	ConditionType *string `json:"conditionType,omitempty"`
	Field         *string `json:"field,omitempty"`
	Value         *string `json:"value,omitempty"`
}

// Type of rule condition.
type RuleConditionType string

type ruleConditionTypeValuesType struct {
	When                              RuleConditionType
	WhenNot                           RuleConditionType
	WhenChanged                       RuleConditionType
	WhenNotChanged                    RuleConditionType
	WhenWas                           RuleConditionType
	WhenStateChangedTo                RuleConditionType
	WhenStateChangedFromAndTo         RuleConditionType
	WhenWorkItemIsCreated             RuleConditionType
	WhenValueIsDefined                RuleConditionType
	WhenValueIsNotDefined             RuleConditionType
	WhenCurrentUserIsMemberOfGroup    RuleConditionType
	WhenCurrentUserIsNotMemberOfGroup RuleConditionType
}

var RuleConditionTypeValues = ruleConditionTypeValuesType{
	// $When. This condition limits the execution of its children to cases when another field has a particular value, i.e. when the Is value of the referenced field is equal to the given literal value.
	When: "when",
	// $WhenNot.This condition limits the execution of its children to cases when another field does not have a particular value, i.e.when the Is value of the referenced field is not equal to the given literal value.
	WhenNot: "whenNot",
	// $WhenChanged.This condition limits the execution of its children to cases when another field has changed, i.e.when the Is value of the referenced field is not equal to the Was value of that field.
	WhenChanged: "whenChanged",
	// $WhenNotChanged.This condition limits the execution of its children to cases when another field has not changed, i.e.when the Is value of the referenced field is equal to the Was value of that field.
	WhenNotChanged:                    "whenNotChanged",
	WhenWas:                           "whenWas",
	WhenStateChangedTo:                "whenStateChangedTo",
	WhenStateChangedFromAndTo:         "whenStateChangedFromAndTo",
	WhenWorkItemIsCreated:             "whenWorkItemIsCreated",
	WhenValueIsDefined:                "whenValueIsDefined",
	WhenValueIsNotDefined:             "whenValueIsNotDefined",
	WhenCurrentUserIsMemberOfGroup:    "whenCurrentUserIsMemberOfGroup",
	WhenCurrentUserIsNotMemberOfGroup: "whenCurrentUserIsNotMemberOfGroup",
}

// Defines a section of the work item form layout
type Section struct {
	// List of child groups in this section
	Groups *[]Group `json:"groups,omitempty"`
	// The id for the layout node.
	Id *string `json:"id,omitempty"`
	// A value indicating whether this layout node has been overridden by a child layout.
	Overridden *bool `json:"overridden,omitempty"`
}

// Describes a request to update a process
type UpdateProcessModel struct {
	// New description of the process
	Description *string `json:"description,omitempty"`
	// If true new projects will use this process by default
	IsDefault *bool `json:"isDefault,omitempty"`
	// If false the process will be disabled and cannot be used to create projects
	IsEnabled *bool `json:"isEnabled,omitempty"`
	// New name of the process
	Name *string `json:"name,omitempty"`
}

// Request class/object to update the rule.
type UpdateProcessRuleRequest struct {
	// List of actions to take when the rule is triggered.
	Actions *[]RuleAction `json:"actions,omitempty"`
	// List of conditions when the rule should be triggered.
	Conditions *[]RuleCondition `json:"conditions,omitempty"`
	// Indicates if the rule is disabled.
	IsDisabled *bool `json:"isDisabled,omitempty"`
	// Name for the rule.
	Name *string `json:"name,omitempty"`
	// Id to uniquely identify the rule.
	Id *uuid.UUID `json:"id,omitempty"`
}

// Class to describe a request that updates a field's properties in a work item type.
type UpdateProcessWorkItemTypeFieldRequest struct {
	// The list of field allowed values.
	AllowedValues *[]string `json:"allowedValues,omitempty"`
	// Allow setting field value to a group identity. Only applies to identity fields.
	AllowGroups *bool `json:"allowGroups,omitempty"`
	// The default value of the field.
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	// If true the field cannot be edited.
	ReadOnly *bool `json:"readOnly,omitempty"`
	// The default value of the field.
	Required *bool `json:"required,omitempty"`
}

// Class for update request on a work item type
type UpdateProcessWorkItemTypeRequest struct {
	// Color of the work item type
	Color *string `json:"color,omitempty"`
	// Description of the work item type
	Description *string `json:"description,omitempty"`
	// Icon of the work item type
	Icon *string `json:"icon,omitempty"`
	// If set will disable the work item type
	IsDisabled *bool `json:"isDisabled,omitempty"`
}

// Properties of a work item form contribution
type WitContribution struct {
	// The id for the contribution.
	ContributionId *string `json:"contributionId,omitempty"`
	// The height for the contribution.
	Height *int `json:"height,omitempty"`
	// A dictionary holding key value pairs for contribution inputs.
	Inputs *map[string]interface{} `json:"inputs,omitempty"`
	// A value indicating if the contribution should be show on deleted workItem.
	ShowOnDeletedWorkItem *bool `json:"showOnDeletedWorkItem,omitempty"`
}

type WorkItemBehavior struct {
	Abstract    *bool                      `json:"abstract,omitempty"`
	Color       *string                    `json:"color,omitempty"`
	Description *string                    `json:"description,omitempty"`
	Fields      *[]WorkItemBehaviorField   `json:"fields,omitempty"`
	Id          *string                    `json:"id,omitempty"`
	Inherits    *WorkItemBehaviorReference `json:"inherits,omitempty"`
	Name        *string                    `json:"name,omitempty"`
	Overriden   *bool                      `json:"overriden,omitempty"`
	Rank        *int                       `json:"rank,omitempty"`
	Url         *string                    `json:"url,omitempty"`
}

type WorkItemBehaviorField struct {
	BehaviorFieldId *string `json:"behaviorFieldId,omitempty"`
	Id              *string `json:"id,omitempty"`
	Url             *string `json:"url,omitempty"`
}

// Reference to the behavior of a work item type.
type WorkItemBehaviorReference struct {
	// The ID of the reference behavior.
	Id *string `json:"id,omitempty"`
	// The url of the reference behavior.
	Url *string `json:"url,omitempty"`
}

// Class That represents a work item state input.
type WorkItemStateInputModel struct {
	// Color of the state
	Color *string `json:"color,omitempty"`
	// Name of the state
	Name *string `json:"name,omitempty"`
	// Order in which state should appear
	Order *int `json:"order,omitempty"`
	// Category of the state
	StateCategory *string `json:"stateCategory,omitempty"`
}

// Class that represents a work item state result.
type WorkItemStateResultModel struct {
	// Work item state color.
	Color *string `json:"color,omitempty"`
	// Work item state customization type.
	CustomizationType *CustomizationType `json:"customizationType,omitempty"`
	// If the Work item state is hidden.
	Hidden *bool `json:"hidden,omitempty"`
	// Id of the Workitemstate.
	Id *uuid.UUID `json:"id,omitempty"`
	// Work item state name.
	Name *string `json:"name,omitempty"`
	// Work item state order.
	Order *int `json:"order,omitempty"`
	// Work item state statecategory.
	StateCategory *string `json:"stateCategory,omitempty"`
	// Work item state url.
	Url *string `json:"url,omitempty"`
}

// Association between a work item type and it's behavior
type WorkItemTypeBehavior struct {
	// Reference to the behavior of a work item type
	Behavior *WorkItemBehaviorReference `json:"behavior,omitempty"`
	// If true the work item type is the default work item type in the behavior
	IsDefault *bool `json:"isDefault,omitempty"`
	// If true the work item type is the default work item type in the parent behavior
	IsLegacyDefault *bool `json:"isLegacyDefault,omitempty"`
	// URL of the work item type behavior
	Url *string `json:"url,omitempty"`
}

type WorkItemTypeClass string

type workItemTypeClassValuesType struct {
	System  WorkItemTypeClass
	Derived WorkItemTypeClass
	Custom  WorkItemTypeClass
}

var WorkItemTypeClassValues = workItemTypeClassValuesType{
	System:  "system",
	Derived: "derived",
	Custom:  "custom",
}

type WorkItemTypeModel struct {
	Behaviors   *[]WorkItemTypeBehavior `json:"behaviors,omitempty"`
	Class       *WorkItemTypeClass      `json:"class,omitempty"`
	Color       *string                 `json:"color,omitempty"`
	Description *string                 `json:"description,omitempty"`
	Icon        *string                 `json:"icon,omitempty"`
	Id          *string                 `json:"id,omitempty"`
	// Parent WIT Id/Internal ReferenceName that it inherits from
	Inherits   *string                     `json:"inherits,omitempty"`
	IsDisabled *bool                       `json:"isDisabled,omitempty"`
	Layout     *FormLayout                 `json:"layout,omitempty"`
	Name       *string                     `json:"name,omitempty"`
	States     *[]WorkItemStateResultModel `json:"states,omitempty"`
	Url        *string                     `json:"url,omitempty"`
}

// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package featuremanagement

import (
	"github.com/google/uuid"
)

// A feature that can be enabled or disabled
type ContributedFeature struct {
	// Named links describing the feature
	Links interface{} `json:"_links,omitempty"`
	// If true, the feature is enabled unless overridden at some scope
	DefaultState *bool `json:"defaultState,omitempty"`
	// Rules for setting the default value if not specified by any setting/scope. Evaluated in order until a rule returns an Enabled or Disabled state (not Undefined)
	DefaultValueRules *[]ContributedFeatureValueRule `json:"defaultValueRules,omitempty"`
	// The description of the feature
	Description *string `json:"description,omitempty"`
	// Extra properties for the feature
	FeatureProperties *map[string]interface{} `json:"featureProperties,omitempty"`
	// Handler for listening to setter calls on feature value. These listeners are only invoked after a successful set has occurred
	FeatureStateChangedListeners *[]ContributedFeatureListener `json:"featureStateChangedListeners,omitempty"`
	// The full contribution id of the feature
	Id *string `json:"id,omitempty"`
	// If this is set to true, then the id for this feature will be added to the list of claims for the request.
	IncludeAsClaim *bool `json:"includeAsClaim,omitempty"`
	// The friendly name of the feature
	Name *string `json:"name,omitempty"`
	// Suggested order to display feature in.
	Order *int `json:"order,omitempty"`
	// Rules for overriding a feature value. These rules are run before explicit user/host state values are checked. They are evaluated in order until a rule returns an Enabled or Disabled state (not Undefined)
	OverrideRules *[]ContributedFeatureValueRule `json:"overrideRules,omitempty"`
	// The scopes/levels at which settings can set the enabled/disabled state of this feature
	Scopes *[]ContributedFeatureSettingScope `json:"scopes,omitempty"`
	// The service instance id of the service that owns this feature
	ServiceInstanceType *uuid.UUID `json:"serviceInstanceType,omitempty"`
	// Tags associated with the feature.
	Tags *[]string `json:"tags,omitempty"`
}

// The current state of a feature within a given scope
type ContributedFeatureEnabledValue string

type contributedFeatureEnabledValueValuesType struct {
	Undefined ContributedFeatureEnabledValue
	Disabled  ContributedFeatureEnabledValue
	Enabled   ContributedFeatureEnabledValue
}

var ContributedFeatureEnabledValueValues = contributedFeatureEnabledValueValuesType{
	// The state of the feature is not set for the specified scope
	Undefined: "undefined",
	// The feature is disabled at the specified scope
	Disabled: "disabled",
	// The feature is enabled at the specified scope
	Enabled: "enabled",
}

type ContributedFeatureHandlerSettings struct {
	// Name of the handler to run
	Name *string `json:"name,omitempty"`
	// Properties to feed to the handler
	Properties *map[string]interface{} `json:"properties,omitempty"`
}

// An identifier and properties used to pass into a handler for a listener or plugin
type ContributedFeatureListener struct {
	// Name of the handler to run
	Name *string `json:"name,omitempty"`
	// Properties to feed to the handler
	Properties *map[string]interface{} `json:"properties,omitempty"`
}

// The scope to which a feature setting applies
type ContributedFeatureSettingScope struct {
	// The name of the settings scope to use when reading/writing the setting
	SettingScope *string `json:"settingScope,omitempty"`
	// Whether this is a user-scope or this is a host-wide (all users) setting
	UserScoped *bool `json:"userScoped,omitempty"`
}

// A contributed feature/state pair
type ContributedFeatureState struct {
	// The full contribution id of the feature
	FeatureId *string `json:"featureId,omitempty"`
	// True if the effective state was set by an override rule (indicating that the state cannot be managed by the end user)
	Overridden *bool `json:"overridden,omitempty"`
	// Reason that the state was set (by a plugin/rule).
	Reason *string `json:"reason,omitempty"`
	// The scope at which this state applies
	Scope *ContributedFeatureSettingScope `json:"scope,omitempty"`
	// The current state of this feature
	State *ContributedFeatureEnabledValue `json:"state,omitempty"`
}

// A query for the effective contributed feature states for a list of feature ids
type ContributedFeatureStateQuery struct {
	// The list of feature ids to query
	FeatureIds *[]string `json:"featureIds,omitempty"`
	// The query result containing the current feature states for each of the queried feature ids
	FeatureStates *map[string]ContributedFeatureState `json:"featureStates,omitempty"`
	// A dictionary of scope values (project name, etc.) to use in the query (if querying across scopes)
	ScopeValues *map[string]string `json:"scopeValues,omitempty"`
}

// A rule for dynamically getting the enabled/disabled state of a feature
type ContributedFeatureValueRule struct {
	// Name of the handler to run
	Name *string `json:"name,omitempty"`
	// Properties to feed to the handler
	Properties *map[string]interface{} `json:"properties,omitempty"`
}

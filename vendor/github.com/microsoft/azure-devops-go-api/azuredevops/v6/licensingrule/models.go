// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package licensingrule

import (
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/licensing"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/operations"
)

type ApplicationStatus struct {
	Extensions  *[]ExtensionApplicationStatus `json:"extensions,omitempty"`
	IsTruncated *bool                         `json:"isTruncated,omitempty"`
	Licenses    *[]LicenseApplicationStatus   `json:"licenses,omitempty"`
	Option      *RuleOption                   `json:"option,omitempty"`
	Status      *operations.OperationStatus   `json:"status,omitempty"`
}

type ExtensionApplicationStatus struct {
	Assigned              *int    `json:"assigned,omitempty"`
	Failed                *int    `json:"failed,omitempty"`
	InsufficientResources *int    `json:"insufficientResources,omitempty"`
	ExtensionId           *string `json:"extensionId,omitempty"`
	Incompatible          *int    `json:"incompatible,omitempty"`
	Unassigned            *int    `json:"unassigned,omitempty"`
}

// Represents an Extension Rule
type ExtensionRule struct {
	// Extension Id
	ExtensionId *string `json:"extensionId,omitempty"`
	// Status of the group rule (applied, missing licenses, etc)
	Status *GroupLicensingRuleStatus `json:"status,omitempty"`
}

// Batching of subjects to lookup using the Graph API
type GraphSubjectLookup struct {
	LookupKeys *[]GraphSubjectLookupKey `json:"lookupKeys,omitempty"`
}

type GraphSubjectLookupKey struct {
	Descriptor *string `json:"descriptor,omitempty"`
}

// Represents a GroupLicensingRule
type GroupLicensingRule struct {
	// Extension Rules
	ExtensionRules *[]ExtensionRule `json:"extensionRules,omitempty"`
	// License Rule
	LicenseRule *LicenseRule `json:"licenseRule,omitempty"`
	// SubjectDescriptor for the rule
	SubjectDescriptor *string `json:"subjectDescriptor,omitempty"`
}

type GroupLicensingRuleStatus string

type groupLicensingRuleStatusValuesType struct {
	ApplyPending  GroupLicensingRuleStatus
	Applied       GroupLicensingRuleStatus
	Incompatible  GroupLicensingRuleStatus
	UnableToApply GroupLicensingRuleStatus
}

var GroupLicensingRuleStatusValues = groupLicensingRuleStatusValuesType{
	// Rule is created or updated, but apply is pending
	ApplyPending: "applyPending",
	// Rule is applied
	Applied: "applied",
	// The group rule was incompatible
	Incompatible: "incompatible",
	// Rule failed to apply unexpectedly and should be retried
	UnableToApply: "unableToApply",
}

// Represents an GroupLicensingRuleUpdate Model
type GroupLicensingRuleUpdate struct {
	// Extensions to Add
	ExtensionsToAdd *[]string `json:"extensionsToAdd,omitempty"`
	// Extensions to Remove
	ExtensionsToRemove *[]string `json:"extensionsToRemove,omitempty"`
	// New License
	License *licensing.License `json:"license,omitempty"`
	// SubjectDescriptor for the rule
	SubjectDescriptor *string `json:"subjectDescriptor,omitempty"`
}

type LicenseApplicationStatus struct {
	Assigned              *int                          `json:"assigned,omitempty"`
	Failed                *int                          `json:"failed,omitempty"`
	InsufficientResources *int                          `json:"insufficientResources,omitempty"`
	AccountUserLicense    *licensing.AccountUserLicense `json:"accountUserLicense,omitempty"`
	License               *licensing.License            `json:"license,omitempty"`
}

// Represents a License Rule
type LicenseRule struct {
	// The last time the rule was executed (regardless of whether any changes were made)
	LastExecuted *azuredevops.Time `json:"lastExecuted,omitempty"`
	// Lasted updated timestamp of the licensing rule
	LastUpdated *azuredevops.Time `json:"lastUpdated,omitempty"`
	// License
	License *licensing.License `json:"license,omitempty"`
	// Status of the group rule (applied, missing licenses, etc)
	Status *GroupLicensingRuleStatus `json:"status,omitempty"`
}

type LicensingApplicationStatus struct {
	Assigned              *int `json:"assigned,omitempty"`
	Failed                *int `json:"failed,omitempty"`
	InsufficientResources *int `json:"insufficientResources,omitempty"`
}

type RuleOption string

type ruleOptionValuesType struct {
	ApplyGroupRule     RuleOption
	TestApplyGroupRule RuleOption
}

var RuleOptionValues = ruleOptionValuesType{
	ApplyGroupRule:     "applyGroupRule",
	TestApplyGroupRule: "testApplyGroupRule",
}

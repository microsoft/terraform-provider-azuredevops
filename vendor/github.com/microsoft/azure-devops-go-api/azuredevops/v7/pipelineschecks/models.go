// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package pipelineschecks

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/pipelinesapproval"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/pipelinestaskcheck"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
)

type ApprovalCheckConfiguration struct {
	// Check configuration id.
	Id *int `json:"id,omitempty"`
	// Resource on which check get configured.
	Resource *Resource `json:"resource,omitempty"`
	// Check configuration type
	Type *CheckType `json:"type,omitempty"`
	// The URL from which one can fetch the configured check.
	Url *string `json:"url,omitempty"`
	// Reference links.
	Links interface{} `json:"_links,omitempty"`
	// Identity of person who configured check.
	CreatedBy *webapi.IdentityRef `json:"createdBy,omitempty"`
	// Time when check got configured.
	CreatedOn *azuredevops.Time `json:"createdOn,omitempty"`
	// Issue connected to check configuration.
	Issue *CheckIssue `json:"issue,omitempty"`
	// Identity of person who modified the configured check.
	ModifiedBy *webapi.IdentityRef `json:"modifiedBy,omitempty"`
	// Time when configured check was modified.
	ModifiedOn *azuredevops.Time `json:"modifiedOn,omitempty"`
	// Timeout in minutes for the check.
	Timeout *int `json:"timeout,omitempty"`
	// Settings for the approval check configuration.
	Settings *pipelinesapproval.ApprovalConfigSettings `json:"settings,omitempty"`
}

type GenericCheckConfiguration struct {
	// Check configuration id.
	Id *int `json:"id,omitempty"`
	// Resource on which check get configured.
	Resource *Resource `json:"resource,omitempty"`
	// Check configuration type
	Type *CheckType `json:"type,omitempty"`
	// The URL from which one can fetch the configured check.
	Url *string `json:"url,omitempty"`
	// Reference links.
	Links interface{} `json:"_links,omitempty"`
	// Identity of person who configured check.
	CreatedBy *webapi.IdentityRef `json:"createdBy,omitempty"`
	// Time when check got configured.
	CreatedOn *azuredevops.Time `json:"createdOn,omitempty"`
	// Issue connected to check configuration.
	Issue *CheckIssue `json:"issue,omitempty"`
	// Identity of person who modified the configured check.
	ModifiedBy *webapi.IdentityRef `json:"modifiedBy,omitempty"`
	// Time when configured check was modified.
	ModifiedOn *azuredevops.Time `json:"modifiedOn,omitempty"`
	// Timeout in minutes for the check.
	Timeout *int `json:"timeout,omitempty"`
	// Settings for the generic check configuration.
	Settings interface{} `json:"settings,omitempty"`
}

type CheckConfiguration struct {
	// Check configuration id.
	Id *int `json:"id,omitempty"`
	// Resource on which check get configured.
	Resource *Resource `json:"resource,omitempty"`
	// Check configuration type
	Type *CheckType `json:"type,omitempty"`
	// The URL from which one can fetch the configured check.
	Url *string `json:"url,omitempty"`
	// Reference links.
	Links interface{} `json:"_links,omitempty"`
	// Identity of person who configured check.
	CreatedBy *webapi.IdentityRef `json:"createdBy,omitempty"`
	// Time when check got configured.
	CreatedOn *azuredevops.Time `json:"createdOn,omitempty"`
	// Issue connected to check configuration.
	Issue *CheckIssue `json:"issue,omitempty"`
	// Identity of person who modified the configured check.
	ModifiedBy *webapi.IdentityRef `json:"modifiedBy,omitempty"`
	// Time when configured check was modified.
	ModifiedOn *azuredevops.Time `json:"modifiedOn,omitempty"`
	// Timeout in minutes for the check.
	Timeout *int `json:"timeout,omitempty"`
}

type CheckConfigurationData struct {
	// Definition Ref Id of the particular check.
	DefinitionRefId *uuid.UUID `json:"definitionRefId,omitempty"`
	// Check configuration of the check.
	CheckConfiguration *CheckConfiguration `json:"checkConfiguration,omitempty"`
}

// [Flags]
type CheckConfigurationExpandParameter string

type checkConfigurationExpandParameterValuesType struct {
	None     CheckConfigurationExpandParameter
	Settings CheckConfigurationExpandParameter
}

var CheckConfigurationExpandParameterValues = checkConfigurationExpandParameterValuesType{
	None:     "none",
	Settings: "settings",
}

type CheckConfigurationRef struct {
	// Check configuration id.
	Id *int `json:"id,omitempty"`
	// Resource on which check get configured.
	Resource *Resource `json:"resource,omitempty"`
	// Check configuration type
	Type *CheckType `json:"type,omitempty"`
	// The URL from which one can fetch the configured check.
	Url *string `json:"url,omitempty"`
}

type CheckData struct {
	// List of default check settings
	DefaultCheckSettings *map[string]string `json:"defaultCheckSettings,omitempty"`
	// List of check configuration data
	CheckConfigurationDataList *[]CheckConfigurationData `json:"checkConfigurationDataList,omitempty"`
	// List of check definitions
	CheckDefinitions *[]CheckDefinitionData `json:"checkDefinitions,omitempty"`
	// List of time zones.
	TimeZoneList *[]TimeZone `json:"timeZoneList,omitempty"`
}

type CheckDefinitionData struct {
	// Flag to allow multiple configurations of a particular check on a resource.
	AllowMultipleConfigurations *bool `json:"allowMultipleConfigurations,omitempty"`
	// Check DefinitionRef Id
	DefinitionRefId *uuid.UUID `json:"definitionRefId,omitempty"`
	// Description about the check
	Description *string `json:"description,omitempty"`
	// Details about the check
	CheckDefinition interface{} `json:"checkDefinition,omitempty"`
	// Icon for the check
	Icon *CheckIcon `json:"icon,omitempty"`
	// Name of the check
	Name *string `json:"name,omitempty"`
	// Check UI contribution Dependencies
	UiContributionDependencies *[]string `json:"uiContributionDependencies,omitempty"`
	// Check UI contribution Type
	UiContributionType *string `json:"uiContributionType,omitempty"`
}

type CheckIcon struct {
	// Asset Location of the icon
	AssetLocation *string `json:"assetLocation,omitempty"`
	// Name of the icon
	Name *string `json:"name,omitempty"`
	// Url of the icon
	Url *string `json:"url,omitempty"`
}

// An issue (error, warning) associated with a check configuration.
type CheckIssue struct {
	// A more detailed description of issue.
	DetailedMessage *string `json:"detailedMessage,omitempty"`
	// A description of issue.
	Message *string `json:"message,omitempty"`
	// The type (error, warning) of the issue.
	Type *CheckIssueType `json:"type,omitempty"`
}

// The type of issue based on severity.
type CheckIssueType string

type checkIssueTypeValuesType struct {
	Error   CheckIssueType
	Warning CheckIssueType
}

var CheckIssueTypeValues = checkIssueTypeValuesType{
	Error:   "error",
	Warning: "warning",
}

type CheckRun struct {
	ResultMessage         *string                `json:"resultMessage,omitempty"`
	Status                *CheckRunStatus        `json:"status,omitempty"`
	CompletedDate         *azuredevops.Time      `json:"completedDate,omitempty"`
	CreatedDate           *azuredevops.Time      `json:"createdDate,omitempty"`
	CheckConfigurationRef *CheckConfigurationRef `json:"checkConfigurationRef,omitempty"`
	Id                    *uuid.UUID             `json:"id,omitempty"`
}

type CheckRunResult struct {
	ResultMessage *string         `json:"resultMessage,omitempty"`
	Status        *CheckRunStatus `json:"status,omitempty"`
}

// [Flags]
type CheckRunStatus string

type checkRunStatusValuesType struct {
	None      CheckRunStatus
	Queued    CheckRunStatus
	Running   CheckRunStatus
	Approved  CheckRunStatus
	Rejected  CheckRunStatus
	Canceled  CheckRunStatus
	TimedOut  CheckRunStatus
	Failed    CheckRunStatus
	Completed CheckRunStatus
	All       CheckRunStatus
}

var CheckRunStatusValues = checkRunStatusValuesType{
	None:      "none",
	Queued:    "queued",
	Running:   "running",
	Approved:  "approved",
	Rejected:  "rejected",
	Canceled:  "canceled",
	TimedOut:  "timedOut",
	Failed:    "failed",
	Completed: "completed",
	All:       "all",
}

type CheckSuite struct {
	// Evaluation context for the check suite request
	Context interface{} `json:"context,omitempty"`
	// Unique suite id generated by the pipeline orchestrator for the pipeline check runs request on the list of resources Pipeline orchestrator will used this identifier to map the check requests on a stage
	Id *uuid.UUID `json:"id,omitempty"`
	// Reference links.
	Links interface{} `json:"_links,omitempty"`
	// Completed date of the given check suite request
	CompletedDate *azuredevops.Time `json:"completedDate,omitempty"`
	// List of check runs associated with the given check suite request.
	CheckRuns *[]CheckRun `json:"checkRuns,omitempty"`
	// Optional message for the given check suite request
	Message *string `json:"message,omitempty"`
	// Overall check runs status for the given suite request. This is check suite status
	Status *CheckRunStatus `json:"status,omitempty"`
}

// [Flags]
type CheckSuiteExpandParameter string

type checkSuiteExpandParameterValuesType struct {
	None      CheckSuiteExpandParameter
	Resources CheckSuiteExpandParameter
}

var CheckSuiteExpandParameterValues = checkSuiteExpandParameterValuesType{
	None:      "none",
	Resources: "resources",
}

type CheckSuiteRef struct {
	// Evaluation context for the check suite request
	Context interface{} `json:"context,omitempty"`
	// Unique suite id generated by the pipeline orchestrator for the pipeline check runs request on the list of resources Pipeline orchestrator will used this identifier to map the check requests on a stage
	Id *uuid.UUID `json:"id,omitempty"`
}

type CheckSuiteRequest struct {
	Context   interface{} `json:"context,omitempty"`
	Id        *uuid.UUID  `json:"id,omitempty"`
	Resources *[]Resource `json:"resources,omitempty"`
}

type CheckType struct {
	// Gets or sets check type id.
	Id *uuid.UUID `json:"id,omitempty"`
	// Name of the check type.
	Name *string `json:"name,omitempty"`
}

type Resource struct {
	// Id of the resource.
	Id *string `json:"id,omitempty"`
	// Name of the resource.
	Name *string `json:"name,omitempty"`
	// Type of the resource.
	Type *string `json:"type,omitempty"`
}

type TaskCheckConfiguration struct {
	// Check configuration id.
	Id *int `json:"id,omitempty"`
	// Resource on which check get configured.
	Resource *Resource `json:"resource,omitempty"`
	// Check configuration type
	Type *CheckType `json:"type,omitempty"`
	// The URL from which one can fetch the configured check.
	Url *string `json:"url,omitempty"`
	// Reference links.
	Links interface{} `json:"_links,omitempty"`
	// Identity of person who configured check.
	CreatedBy *webapi.IdentityRef `json:"createdBy,omitempty"`
	// Time when check got configured.
	CreatedOn *azuredevops.Time `json:"createdOn,omitempty"`
	// Issue connected to check configuration.
	Issue *CheckIssue `json:"issue,omitempty"`
	// Identity of person who modified the configured check.
	ModifiedBy *webapi.IdentityRef `json:"modifiedBy,omitempty"`
	// Time when configured check was modified.
	ModifiedOn *azuredevops.Time `json:"modifiedOn,omitempty"`
	// Timeout in minutes for the check.
	Timeout *int `json:"timeout,omitempty"`
	// Settings for the task check configuration.
	Settings *pipelinestaskcheck.TaskCheckConfig `json:"settings,omitempty"`
}

type TimeZone struct {
	// Display name of the time zone.
	DisplayName *string `json:"displayName,omitempty"`
	// Id of the time zone.
	Id *string `json:"id,omitempty"`
}

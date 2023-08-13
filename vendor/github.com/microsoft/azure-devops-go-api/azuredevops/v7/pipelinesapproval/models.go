// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package pipelinesapproval

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
)

type Approval struct {
	// /// Gets the links to access the approval object.
	Links interface{} `json:"_links,omitempty"`
	// Identities which are not allowed to approve.
	BlockedApprovers *[]webapi.IdentityRef `json:"blockedApprovers,omitempty"`
	// Date on which approval got created.
	CreatedOn *azuredevops.Time `json:"createdOn,omitempty"`
	// Order in which approvers will be actionable.
	ExecutionOrder *ApprovalExecutionOrder `json:"executionOrder,omitempty"`
	// Unique identifier of the approval.
	Id *uuid.UUID `json:"id,omitempty"`
	// Instructions for the approvers.
	Instructions *string `json:"instructions,omitempty"`
	// Date on which approval was last modified.
	LastModifiedOn *azuredevops.Time `json:"lastModifiedOn,omitempty"`
	// Minimum number of approvers that should approve for the entire approval to be considered approved.
	MinRequiredApprovers *int `json:"minRequiredApprovers,omitempty"`
	// Current user permissions for approval object.
	Permissions *ApprovalPermissions `json:"permissions,omitempty"`
	// Overall status of the approval.
	Status *ApprovalStatus `json:"status,omitempty"`
	// List of steps associated with the approval.
	Steps *[]ApprovalStep `json:"steps,omitempty"`
}

type ApprovalCompletedNotificationEvent struct {
	Approval  *Approval  `json:"approval,omitempty"`
	ProjectId *uuid.UUID `json:"projectId,omitempty"`
}

// Config to create a new approval.
type ApprovalConfig struct {
	// Ordered list of approvers.
	Approvers *[]webapi.IdentityRef `json:"approvers,omitempty"`
	// Identities which are not allowed to approve.
	BlockedApprovers *[]webapi.IdentityRef `json:"blockedApprovers,omitempty"`
	// Order in which approvers will be actionable.
	ExecutionOrder *ApprovalExecutionOrder `json:"executionOrder,omitempty"`
	// Instructions for the approver.
	Instructions *string `json:"instructions,omitempty"`
	// Minimum number of approvers that should approve for the entire approval to be considered approved. Defaults to all.
	MinRequiredApprovers *int `json:"minRequiredApprovers,omitempty"`
}

// Config to create a new approval.
type ApprovalConfigSettings struct {
	// Ordered list of approvers.
	Approvers *[]webapi.IdentityRef `json:"approvers,omitempty"`
	// Identities which are not allowed to approve.
	BlockedApprovers *[]webapi.IdentityRef `json:"blockedApprovers,omitempty"`
	// Order in which approvers will be actionable.
	ExecutionOrder *ApprovalExecutionOrder `json:"executionOrder,omitempty"`
	// Instructions for the approver.
	Instructions *string `json:"instructions,omitempty"`
	// Minimum number of approvers that should approve for the entire approval to be considered approved. Defaults to all.
	MinRequiredApprovers *int `json:"minRequiredApprovers,omitempty"`
	// Determines whether check requester can approve the check.
	RequesterCannotBeApprover *bool `json:"requesterCannotBeApprover,omitempty"`
}

// [Flags]
type ApprovalDetailsExpandParameter string

type approvalDetailsExpandParameterValuesType struct {
	None        ApprovalDetailsExpandParameter
	Steps       ApprovalDetailsExpandParameter
	Permissions ApprovalDetailsExpandParameter
}

var ApprovalDetailsExpandParameterValues = approvalDetailsExpandParameterValuesType{
	None:        "none",
	Steps:       "steps",
	Permissions: "permissions",
}

type ApprovalExecutionOrder string

type approvalExecutionOrderValuesType struct {
	AnyOrder   ApprovalExecutionOrder
	InSequence ApprovalExecutionOrder
}

var ApprovalExecutionOrderValues = approvalExecutionOrderValuesType{
	// Indicates that the approvers can approve in any order.
	AnyOrder: "anyOrder",
	// Indicates that the approvers can only approve in a sequential order(Order in which they were assigned).
	InSequence: "inSequence",
}

// Data for notification base class for approval events.
type ApprovalNotificationEventBase struct {
	Approval  *Approval  `json:"approval,omitempty"`
	ProjectId *uuid.UUID `json:"projectId,omitempty"`
}

// [Flags]
type ApprovalPermissions string

type approvalPermissionsValuesType struct {
	None          ApprovalPermissions
	View          ApprovalPermissions
	Update        ApprovalPermissions
	Reassign      ApprovalPermissions
	ResourceAdmin ApprovalPermissions
	QueueBuild    ApprovalPermissions
}

var ApprovalPermissionsValues = approvalPermissionsValuesType{
	None:          "none",
	View:          "view",
	Update:        "update",
	Reassign:      "reassign",
	ResourceAdmin: "resourceAdmin",
	QueueBuild:    "queueBuild",
}

// Request to create a new approval.
type ApprovalRequest struct {
	// Unique identifier with which the approval is to be registered.
	ApprovalId *uuid.UUID `json:"approvalId,omitempty"`
	// Configuration of the approval request.
	Config *ApprovalConfig `json:"config,omitempty"`
}

type ApprovalsQueryParameters struct {
	// Query approvals based on list of approval IDs.
	ApprovalIds *[]uuid.UUID `json:"approvalIds,omitempty"`
}

// [Flags] Status of an approval as a whole or of an individual step.
type ApprovalStatus string

type approvalStatusValuesType struct {
	Undefined   ApprovalStatus
	Uninitiated ApprovalStatus
	Pending     ApprovalStatus
	Approved    ApprovalStatus
	Rejected    ApprovalStatus
	Skipped     ApprovalStatus
	Canceled    ApprovalStatus
	TimedOut    ApprovalStatus
	Failed      ApprovalStatus
	Completed   ApprovalStatus
	All         ApprovalStatus
}

var ApprovalStatusValues = approvalStatusValuesType{
	Undefined: "undefined",
	// Indicates the approval is Uninitiated. Used in case of in sequence order of execution where given approver is not yet actionable.
	Uninitiated: "uninitiated",
	// Indicates the approval is Pending.
	Pending: "pending",
	// Indicates the approval is Approved.
	Approved: "approved",
	// Indicates the approval is Rejected.
	Rejected: "rejected",
	// Indicates the approval is Skipped.
	Skipped: "skipped",
	// Indicates the approval is Canceled.
	Canceled: "canceled",
	// Indicates the approval is Timed out.
	TimedOut:  "timedOut",
	Failed:    "failed",
	Completed: "completed",
	All:       "all",
}

// Data for a single approval step.
type ApprovalStep struct {
	// Identity who approved.
	ActualApprover *webapi.IdentityRef `json:"actualApprover,omitempty"`
	// Identity who should approve.
	AssignedApprover *webapi.IdentityRef `json:"assignedApprover,omitempty"`
	// Comment associated with this step.
	Comment *string `json:"comment,omitempty"`
	// History of the approval step
	History *[]ApprovalStepHistory `json:"history,omitempty"`
	// Timestamp at which this step was initiated.
	InitiatedOn *azuredevops.Time `json:"initiatedOn,omitempty"`
	// Identity by which this step was last modified.
	LastModifiedBy *webapi.IdentityRef `json:"lastModifiedBy,omitempty"`
	// Timestamp at which this step was last modified.
	LastModifiedOn *azuredevops.Time `json:"lastModifiedOn,omitempty"`
	// Order in which the approvers are allowed to approve.
	Order *int `json:"order,omitempty"`
	// Current user permissions for step.
	Permissions *ApprovalPermissions `json:"permissions,omitempty"`
	// Current status of this step.
	Status *ApprovalStatus `json:"status,omitempty"`
}

// Data for a single approval step history.
type ApprovalStepHistory struct {
	// Identity who was assigned this approval
	AssignedTo *webapi.IdentityRef `json:"assignedTo,omitempty"`
	// Comment associated with this step history.
	Comment *string `json:"comment,omitempty"`
	// Identity by which this step history was created.
	CreatedBy *webapi.IdentityRef `json:"createdBy,omitempty"`
	// Timestamp at which this step history was created.
	CreatedOn *azuredevops.Time `json:"createdOn,omitempty"`
}

// Data to update an approval object or its individual step.
type ApprovalUpdateParameters struct {
	// ID of the approval to be updated.
	ApprovalId *uuid.UUID `json:"approvalId,omitempty"`
	// Current approver.
	AssignedApprover *webapi.IdentityRef `json:"assignedApprover,omitempty"`
	// Gets or sets comment.
	Comment *string `json:"comment,omitempty"`
	// Reassigned Approver.
	ReassignTo *webapi.IdentityRef `json:"reassignTo,omitempty"`
	// Gets or sets status.
	Status *ApprovalStatus `json:"status,omitempty"`
}

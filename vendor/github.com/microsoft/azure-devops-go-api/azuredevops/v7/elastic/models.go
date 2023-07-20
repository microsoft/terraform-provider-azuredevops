// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package elastic

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
)

// [Flags]
type ElasticAgentState string

type elasticAgentStateValuesType struct {
	None     ElasticAgentState
	Enabled  ElasticAgentState
	Online   ElasticAgentState
	Assigned ElasticAgentState
}

var ElasticAgentStateValues = elasticAgentStateValuesType{
	None:     "none",
	Enabled:  "enabled",
	Online:   "online",
	Assigned: "assigned",
}

type ElasticComputeState string

type elasticComputeStateValuesType struct {
	None      ElasticComputeState
	Healthy   ElasticComputeState
	Creating  ElasticComputeState
	Deleting  ElasticComputeState
	Failed    ElasticComputeState
	Stopped   ElasticComputeState
	Reimaging ElasticComputeState
}

var ElasticComputeStateValues = elasticComputeStateValuesType{
	None:      "none",
	Healthy:   "healthy",
	Creating:  "creating",
	Deleting:  "deleting",
	Failed:    "failed",
	Stopped:   "stopped",
	Reimaging: "reimaging",
}

// Data and settings for an elastic node
type ElasticNode struct {
	// Distributed Task's Agent Id
	AgentId *int `json:"agentId,omitempty"`
	// Summary of the state of the agent
	AgentState *ElasticAgentState `json:"agentState,omitempty"`
	// Compute Id.  VMSS's InstanceId
	ComputeId *string `json:"computeId,omitempty"`
	// State of the compute host
	ComputeState *ElasticComputeState `json:"computeState,omitempty"`
	// Users can force state changes to specific states (ToReimage, ToDelete, Save)
	DesiredState *ElasticNodeState `json:"desiredState,omitempty"`
	// Unique identifier since the agent and/or VM may be null
	Id *int `json:"id,omitempty"`
	// Computer name. Used to match a scaleset VM with an agent
	Name *string `json:"name,omitempty"`
	// Pool Id that this node belongs to
	PoolId *int `json:"poolId,omitempty"`
	// Last job RequestId assigned to this agent
	RequestId *uint64 `json:"requestId,omitempty"`
	// State of the ElasticNode
	State *ElasticNodeState `json:"state,omitempty"`
	// Last state change. Only updated by SQL.
	StateChangedOn *azuredevops.Time `json:"stateChangedOn,omitempty"`
}

// Class used for updating an elastic node where only certain members are populated
type ElasticNodeSettings struct {
	// State of the ElasticNode
	State *ElasticNodeState `json:"state,omitempty"`
}

type ElasticNodeState string

type elasticNodeStateValuesType struct {
	None                         ElasticNodeState
	New                          ElasticNodeState
	CreatingCompute              ElasticNodeState
	StartingAgent                ElasticNodeState
	Idle                         ElasticNodeState
	Assigned                     ElasticNodeState
	Offline                      ElasticNodeState
	PendingReimage               ElasticNodeState
	PendingDelete                ElasticNodeState
	Saved                        ElasticNodeState
	DeletingCompute              ElasticNodeState
	Deleted                      ElasticNodeState
	Lost                         ElasticNodeState
	ReimagingCompute             ElasticNodeState
	RestartingAgent              ElasticNodeState
	FailedToStartPendingDelete   ElasticNodeState
	FailedToRestartPendingDelete ElasticNodeState
	FailedVMPendingDelete        ElasticNodeState
	AssignedPendingDelete        ElasticNodeState
	RetryDelete                  ElasticNodeState
}

var ElasticNodeStateValues = elasticNodeStateValuesType{
	None:                         "none",
	New:                          "new",
	CreatingCompute:              "creatingCompute",
	StartingAgent:                "startingAgent",
	Idle:                         "idle",
	Assigned:                     "assigned",
	Offline:                      "offline",
	PendingReimage:               "pendingReimage",
	PendingDelete:                "pendingDelete",
	Saved:                        "saved",
	DeletingCompute:              "deletingCompute",
	Deleted:                      "deleted",
	Lost:                         "lost",
	ReimagingCompute:             "reimagingCompute",
	RestartingAgent:              "restartingAgent",
	FailedToStartPendingDelete:   "failedToStartPendingDelete",
	FailedToRestartPendingDelete: "failedToRestartPendingDelete",
	FailedVMPendingDelete:        "failedVMPendingDelete",
	AssignedPendingDelete:        "assignedPendingDelete",
	RetryDelete:                  "retryDelete",
}

// Data and settings for an elastic pool
type ElasticPool struct {
	// Set whether agents should be configured to run with interactive UI
	AgentInteractiveUI *bool `json:"agentInteractiveUI,omitempty"`
	// Azure string representing to location of the resource
	AzureId *string `json:"azureId,omitempty"`
	// Number of agents to have ready waiting for jobs
	DesiredIdle *int `json:"desiredIdle,omitempty"`
	// The desired size of the pool
	DesiredSize *int `json:"desiredSize,omitempty"`
	// Maximum number of nodes that will exist in the elastic pool
	MaxCapacity *int `json:"maxCapacity,omitempty"`
	// Keep nodes in the pool on failure for investigation
	MaxSavedNodeCount *int `json:"maxSavedNodeCount,omitempty"`
	// Timestamp the pool was first detected to be offline
	OfflineSince *azuredevops.Time `json:"offlineSince,omitempty"`
	// Operating system type of the nodes in the pool
	OrchestrationType *OrchestrationType `json:"orchestrationType,omitempty"`
	// Operating system type of the nodes in the pool
	OsType *OperatingSystemType `json:"osType,omitempty"`
	// Id of the associated TaskAgentPool
	PoolId *int `json:"poolId,omitempty"`
	// Discard node after each job completes
	RecycleAfterEachUse *bool `json:"recycleAfterEachUse,omitempty"`
	// Id of the Service Endpoint used to connect to Azure
	ServiceEndpointId *uuid.UUID `json:"serviceEndpointId,omitempty"`
	// Scope the Service Endpoint belongs to
	ServiceEndpointScope *uuid.UUID `json:"serviceEndpointScope,omitempty"`
	// The number of sizing attempts executed while trying to achieve a desired size
	SizingAttempts *int `json:"sizingAttempts,omitempty"`
	// State of the pool
	State *ElasticPoolState `json:"state,omitempty"`
	// The minimum time in minutes to keep idle agents alive
	TimeToLiveMinutes *int `json:"timeToLiveMinutes,omitempty"`
}

// Returned result from creating a new elastic pool
type ElasticPoolCreationResult struct {
	// Created agent pool
	AgentPool *TaskAgentPool `json:"agentPool,omitempty"`
	// Created agent queue
	AgentQueue *TaskAgentQueue `json:"agentQueue,omitempty"`
	// Created elastic pool
	ElasticPool *ElasticPool `json:"elasticPool,omitempty"`
}

// Log data for an Elastic Pool
type ElasticPoolLog struct {
	// Log Id
	Id *uint64 `json:"id,omitempty"`
	// E.g. error, warning, info
	Level *LogLevel `json:"level,omitempty"`
	// Log contents
	Message *string `json:"message,omitempty"`
	// Operation that triggered the message being logged
	Operation *OperationType `json:"operation,omitempty"`
	// Id of the associated TaskAgentPool
	PoolId *int `json:"poolId,omitempty"`
	// Datetime that the log occurred
	Timestamp *azuredevops.Time `json:"timestamp,omitempty"`
}

// Class used for updating an elastic pool where only certain members are populated
type ElasticPoolSettings struct {
	// Set whether agents should be configured to run with interactive UI
	AgentInteractiveUI *bool `json:"agentInteractiveUI,omitempty"`
	// Azure string representing to location of the resource
	AzureId *string `json:"azureId,omitempty"`
	// Number of machines to have ready waiting for jobs
	DesiredIdle *int `json:"desiredIdle,omitempty"`
	// Maximum number of machines that will exist in the elastic pool
	MaxCapacity *int `json:"maxCapacity,omitempty"`
	// Keep machines in the pool on failure for investigation
	MaxSavedNodeCount *int `json:"maxSavedNodeCount,omitempty"`
	// Operating system type of the machines in the pool
	OrchestrationType *OrchestrationType `json:"orchestrationType,omitempty"`
	// Operating system type of the machines in the pool
	OsType *OperatingSystemType `json:"osType,omitempty"`
	// Discard machines after each job completes
	RecycleAfterEachUse *bool `json:"recycleAfterEachUse,omitempty"`
	// Id of the Service Endpoint used to connect to Azure
	ServiceEndpointId *uuid.UUID `json:"serviceEndpointId,omitempty"`
	// Scope the Service Endpoint belongs to
	ServiceEndpointScope *uuid.UUID `json:"serviceEndpointScope,omitempty"`
	// The minimum time in minutes to keep idle agents alive
	TimeToLiveMinutes *int `json:"timeToLiveMinutes,omitempty"`
}

type ElasticPoolState string

type elasticPoolStateValuesType struct {
	Online    ElasticPoolState
	Offline   ElasticPoolState
	Unhealthy ElasticPoolState
	New       ElasticPoolState
}

var ElasticPoolStateValues = elasticPoolStateValuesType{
	// Online and healthy
	Online:    "online",
	Offline:   "offline",
	Unhealthy: "unhealthy",
	New:       "new",
}

type LogLevel string

type logLevelValuesType struct {
	Error   LogLevel
	Warning LogLevel
	Info    LogLevel
}

var LogLevelValues = logLevelValuesType{
	Error:   "error",
	Warning: "warning",
	Info:    "info",
}

type OperatingSystemType string

type operatingSystemTypeValuesType struct {
	Windows OperatingSystemType
	Linux   OperatingSystemType
}

var OperatingSystemTypeValues = operatingSystemTypeValuesType{
	Windows: "windows",
	Linux:   "linux",
}

type OperationType string

type operationTypeValuesType struct {
	ConfigurationJob OperationType
	SizingJob        OperationType
	IncreaseCapacity OperationType
	Reimage          OperationType
	DeleteVMs        OperationType
}

var OperationTypeValues = operationTypeValuesType{
	ConfigurationJob: "configurationJob",
	SizingJob:        "sizingJob",
	IncreaseCapacity: "increaseCapacity",
	Reimage:          "reimage",
	DeleteVMs:        "deleteVMs",
}

type OrchestrationType string

type orchestrationTypeValuesType struct {
	Uniform  OrchestrationType
	Flexible OrchestrationType
}

var OrchestrationTypeValues = orchestrationTypeValuesType{
	Uniform:  "uniform",
	Flexible: "flexible",
}

// An organization-level grouping of agents.
type TaskAgentPool struct {
	Id *int `json:"id,omitempty"`
	// Gets or sets a value indicating whether or not this pool is managed by the service.
	IsHosted *bool `json:"isHosted,omitempty"`
	// Determines whether the pool is legacy.
	IsLegacy *bool   `json:"isLegacy,omitempty"`
	Name     *string `json:"name,omitempty"`
	// Additional pool settings and details
	Options *TaskAgentPoolOptions `json:"options,omitempty"`
	// Gets or sets the type of the pool
	PoolType *TaskAgentPoolType `json:"poolType,omitempty"`
	Scope    *uuid.UUID         `json:"scope,omitempty"`
	// Gets the current size of the pool.
	Size *int `json:"size,omitempty"`
	// The ID of the associated agent cloud.
	AgentCloudId *int `json:"agentCloudId,omitempty"`
	// Whether or not a queue should be automatically provisioned for each project collection.
	AutoProvision *bool `json:"autoProvision,omitempty"`
	// Whether or not the pool should autosize itself based on the Agent Cloud Provider settings.
	AutoSize *bool `json:"autoSize,omitempty"`
	// Whether or not agents in this pool are allowed to automatically update
	AutoUpdate *bool `json:"autoUpdate,omitempty"`
	// Creator of the pool. The creator of the pool is automatically added into the administrators group for the pool on creation.
	CreatedBy *webapi.IdentityRef `json:"createdBy,omitempty"`
	// The date/time of the pool creation.
	CreatedOn *azuredevops.Time `json:"createdOn,omitempty"`
	// Owner or administrator of the pool.
	Owner      *webapi.IdentityRef `json:"owner,omitempty"`
	Properties interface{}         `json:"properties,omitempty"`
	// Target parallelism - Only applies to agent pools that are backed by pool providers. It will be null for regular pools.
	TargetSize *int `json:"targetSize,omitempty"`
}

// [Flags] Additional settings and descriptors for a TaskAgentPool
type TaskAgentPoolOptions string

type taskAgentPoolOptionsValuesType struct {
	None                      TaskAgentPoolOptions
	ElasticPool               TaskAgentPoolOptions
	SingleUseAgents           TaskAgentPoolOptions
	PreserveAgentOnJobFailure TaskAgentPoolOptions
}

var TaskAgentPoolOptionsValues = taskAgentPoolOptionsValuesType{
	None: "none",
	// TaskAgentPool backed by the Elastic pool service
	ElasticPool: "elasticPool",
	// Set to true if agents are re-imaged after each TaskAgentJobRequest
	SingleUseAgents: "singleUseAgents",
	// Set to true if agents are held for investigation after a TaskAgentJobRequest failure
	PreserveAgentOnJobFailure: "preserveAgentOnJobFailure",
}

type TaskAgentPoolReference struct {
	Id *int `json:"id,omitempty"`
	// Gets or sets a value indicating whether or not this pool is managed by the service.
	IsHosted *bool `json:"isHosted,omitempty"`
	// Determines whether the pool is legacy.
	IsLegacy *bool   `json:"isLegacy,omitempty"`
	Name     *string `json:"name,omitempty"`
	// Additional pool settings and details
	Options *TaskAgentPoolOptions `json:"options,omitempty"`
	// Gets or sets the type of the pool
	PoolType *TaskAgentPoolType `json:"poolType,omitempty"`
	Scope    *uuid.UUID         `json:"scope,omitempty"`
	// Gets the current size of the pool.
	Size *int `json:"size,omitempty"`
}

// The type of agent pool.
type TaskAgentPoolType string

type taskAgentPoolTypeValuesType struct {
	Automation TaskAgentPoolType
	Deployment TaskAgentPoolType
}

var TaskAgentPoolTypeValues = taskAgentPoolTypeValuesType{
	// A typical pool of task agents
	Automation: "automation",
	// A deployment pool
	Deployment: "deployment",
}

// An agent queue.
type TaskAgentQueue struct {
	// ID of the queue
	Id *int `json:"id,omitempty"`
	// Name of the queue
	Name *string `json:"name,omitempty"`
	// Pool reference for this queue
	Pool *TaskAgentPoolReference `json:"pool,omitempty"`
	// Project ID
	ProjectId *uuid.UUID `json:"projectId,omitempty"`
}

// The result of an operation tracked by a timeline record.
type TaskResult string

type taskResultValuesType struct {
	Succeeded           TaskResult
	SucceededWithIssues TaskResult
	Failed              TaskResult
	Canceled            TaskResult
	Skipped             TaskResult
	Abandoned           TaskResult
}

var TaskResultValues = taskResultValuesType{
	Succeeded:           "succeeded",
	SucceededWithIssues: "succeededWithIssues",
	Failed:              "failed",
	Canceled:            "canceled",
	Skipped:             "skipped",
	Abandoned:           "abandoned",
}

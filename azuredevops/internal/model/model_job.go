package model

type PipelineJobType string

type pipelineJobTypeValuesType struct {
	AgentJob     PipelineJobType
	AgentlessJob PipelineJobType
}

var PipelineJobTypeValues = pipelineJobTypeValuesType{
	AgentJob:     "AgentJob",
	AgentlessJob: "AgentlessJob",
}

var PipelineJobTypeTypeValueMap = map[string]int{
	string(PipelineJobTypeValues.AgentJob):     1,
	string(PipelineJobTypeValues.AgentlessJob): 2,
}

var PipelineJobTypeValueTypeMap = map[int]string{
	1: string(PipelineJobTypeValues.AgentJob),
	2: string(PipelineJobTypeValues.AgentlessJob),
}

type JobExecutionOptionsType string

type JobExecutionOptionsTypeValuesType struct {
	None               JobExecutionOptionsType
	MultiConfiguration JobExecutionOptionsType
	MultiAgent         JobExecutionOptionsType
}

var JobExecutionOptionsTypeValues = JobExecutionOptionsTypeValuesType{
	None:               "None",                // 0
	MultiConfiguration: "Multi-Configuration", // 1
	MultiAgent:         "Multi-Agent",         // 2
}

var JobExecutionOptionsTypValueTypeMap = map[int]string{
	0: string(JobExecutionOptionsTypeValues.None),
	1: string(JobExecutionOptionsTypeValues.MultiConfiguration),
	2: string(JobExecutionOptionsTypeValues.MultiAgent),
}

type JobExecutionOptions struct {
	Multipliers     *[]string `json:"multipliers,omitempty"`
	MaxConcurrency  *int      `json:"maxConcurrency,omitempty"`
	ContinueOnError *bool     `json:"continueOnError,omitempty"`
	Type            *int      `json:"type,omitempty"`
}

type JobDependency struct {
	Scope *string `json:"scope,omitempty"`
	Event *string `json:"event,omitempty"`
}

type JobTarget struct {
	Demands                      *[]string            `json:"demands,omitempty"`
	ExecutionOptions             *JobExecutionOptions `json:"executionOptions,omitempty"`
	Type                         *int                 `json:"type,omitempty"`
	AllowScriptsAuthAccessOption *bool                `json:"allowScriptsAuthAccessOption,omitempty"`
}

type PipelineJob struct {
	Name                      *string          `json:"name,omitempty"`
	RefName                   *string          `json:"refName,omitempty"`
	Condition                 *string          `json:"condition,omitempty"`
	Dependencies              *[]JobDependency `json:"dependencies,omitempty"`
	Target                    *JobTarget       `json:"target,omitempty"`
	JobTimeoutInMinutes       *int             `json:"jobTimeoutInMinutes,omitempty"`
	JobCancelTimeoutInMinutes *int             `json:"jobCancelTimeoutInMinutes,omitempty"`
	JobAuthorizationScope     *string          `json:"JobAuthorizationScope,omitempty"`
}

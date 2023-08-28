// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package pipelines

import (
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
)

// Artifacts are collections of files produced by a pipeline. Use artifacts to share files between stages in a pipeline or between different pipelines.
type Artifact struct {
	// The name of the artifact.
	Name *string `json:"name,omitempty"`
	// Signed url for downloading this artifact
	SignedContent *webapi.SignedUrl `json:"signedContent,omitempty"`
	// Self-referential url
	Url *string `json:"url,omitempty"`
}

type BuildResourceParameters struct {
	Version *string `json:"version,omitempty"`
}

type ConfigurationType string

type configurationTypeValuesType struct {
	Unknown            ConfigurationType
	Yaml               ConfigurationType
	DesignerJson       ConfigurationType
	JustInTime         ConfigurationType
	DesignerHyphenJson ConfigurationType
}

var ConfigurationTypeValues = configurationTypeValuesType{
	// Unknown type.
	Unknown: "unknown",
	// YAML.
	Yaml: "yaml",
	// Designer JSON.
	DesignerJson: "designerJson",
	// Just-in-time.
	JustInTime: "justInTime",
	// Designer-JSON.
	DesignerHyphenJson: "designerHyphenJson",
}

type Container struct {
	Environment     *map[string]string `json:"environment,omitempty"`
	Image           *string            `json:"image,omitempty"`
	MapDockerSocket *bool              `json:"mapDockerSocket,omitempty"`
	Options         *string            `json:"options,omitempty"`
	Ports           *[]string          `json:"ports,omitempty"`
	Volumes         *[]string          `json:"volumes,omitempty"`
}

type ContainerResource struct {
	Container *Container `json:"container,omitempty"`
}

type ContainerResourceParameters struct {
	Version *string `json:"version,omitempty"`
}

// Configuration parameters of the pipeline.
type CreatePipelineConfigurationParameters struct {
	// Type of configuration.
	Type *ConfigurationType `json:"type,omitempty"`
}

// Parameters to create a pipeline.
type CreatePipelineParameters struct {
	// Configuration parameters of the pipeline.
	Configuration *CreatePipelineConfigurationParameters `json:"configuration,omitempty"`
	// Folder of the pipeline.
	Folder *string `json:"folder,omitempty"`
	// Name of the pipeline.
	Name *string `json:"name,omitempty"`
}

// [Flags] Expansion options for GetArtifact and ListArtifacts.
type GetArtifactExpandOptions string

type getArtifactExpandOptionsValuesType struct {
	None          GetArtifactExpandOptions
	SignedContent GetArtifactExpandOptions
}

var GetArtifactExpandOptionsValues = getArtifactExpandOptionsValuesType{
	// No expansion.
	None: "none",
	// Include signed content.
	SignedContent: "signedContent",
}

// [Flags] $expand options for GetLog and ListLogs.
type GetLogExpandOptions string

type getLogExpandOptionsValuesType struct {
	None          GetLogExpandOptions
	SignedContent GetLogExpandOptions
}

var GetLogExpandOptionsValues = getLogExpandOptionsValuesType{
	None:          "none",
	SignedContent: "signedContent",
}

// Log for a pipeline.
type Log struct {
	// The date and time the log was created.
	CreatedOn *azuredevops.Time `json:"createdOn,omitempty"`
	// The ID of the log.
	Id *int `json:"id,omitempty"`
	// The date and time the log was last changed.
	LastChangedOn *azuredevops.Time `json:"lastChangedOn,omitempty"`
	// The number of lines in the log.
	LineCount     *uint64           `json:"lineCount,omitempty"`
	SignedContent *webapi.SignedUrl `json:"signedContent,omitempty"`
	Url           *string           `json:"url,omitempty"`
}

// A collection of logs.
type LogCollection struct {
	// The list of logs.
	Logs          *[]Log            `json:"logs,omitempty"`
	SignedContent *webapi.SignedUrl `json:"signedContent,omitempty"`
	// URL of the log.
	Url *string `json:"url,omitempty"`
}

type PackageResourceParameters struct {
	Version *string `json:"version,omitempty"`
}

// Definition of a pipeline.
type Pipeline struct {
	// Pipeline folder
	Folder *string `json:"folder,omitempty"`
	// Pipeline ID
	Id *int `json:"id,omitempty"`
	// Pipeline name
	Name *string `json:"name,omitempty"`
	// Revision number
	Revision      *int                   `json:"revision,omitempty"`
	Links         interface{}            `json:"_links,omitempty"`
	Configuration *PipelineConfiguration `json:"configuration,omitempty"`
	// URL of the pipeline
	Url *string `json:"url,omitempty"`
}

type PipelineBase struct {
	// Pipeline folder
	Folder *string `json:"folder,omitempty"`
	// Pipeline ID
	Id *int `json:"id,omitempty"`
	// Pipeline name
	Name *string `json:"name,omitempty"`
	// Revision number
	Revision *int `json:"revision,omitempty"`
}

type PipelineConfiguration struct {
	Type *ConfigurationType `json:"type,omitempty"`
}

// A reference to a Pipeline.
type PipelineReference struct {
	// Pipeline folder
	Folder *string `json:"folder,omitempty"`
	// Pipeline ID
	Id *int `json:"id,omitempty"`
	// Pipeline name
	Name *string `json:"name,omitempty"`
	// Revision number
	Revision *int    `json:"revision,omitempty"`
	Url      *string `json:"url,omitempty"`
}

type PipelineResource struct {
	Pipeline *PipelineReference `json:"pipeline,omitempty"`
	Version  *string            `json:"version,omitempty"`
}

type PipelineResourceParameters struct {
	Version *string `json:"version,omitempty"`
}

type PreviewRun struct {
	FinalYaml *string `json:"finalYaml,omitempty"`
}

type Repository struct {
	Type *RepositoryType `json:"type,omitempty"`
}

type RepositoryResource struct {
	RefName    *string     `json:"refName,omitempty"`
	Repository *Repository `json:"repository,omitempty"`
	Version    *string     `json:"version,omitempty"`
}

type RepositoryResourceParameters struct {
	RefName *string `json:"refName,omitempty"`
	// This is the security token to use when connecting to the repository.
	Token *string `json:"token,omitempty"`
	// Optional. This is the type of the token given. If not provided, a type of "Bearer" is assumed. Note: Use "Basic" for a PAT token.
	TokenType *string `json:"tokenType,omitempty"`
	Version   *string `json:"version,omitempty"`
}

type RepositoryType string

type repositoryTypeValuesType struct {
	Unknown                 RepositoryType
	GitHub                  RepositoryType
	AzureReposGit           RepositoryType
	GitHubEnterprise        RepositoryType
	AzureReposGitHyphenated RepositoryType
}

var RepositoryTypeValues = repositoryTypeValuesType{
	Unknown:                 "unknown",
	GitHub:                  "gitHub",
	AzureReposGit:           "azureReposGit",
	GitHubEnterprise:        "gitHubEnterprise",
	AzureReposGitHyphenated: "azureReposGitHyphenated",
}

type Run struct {
	Id                 *int                    `json:"id,omitempty"`
	Name               *string                 `json:"name,omitempty"`
	Links              interface{}             `json:"_links,omitempty"`
	CreatedDate        *azuredevops.Time       `json:"createdDate,omitempty"`
	FinalYaml          *string                 `json:"finalYaml,omitempty"`
	FinishedDate       *azuredevops.Time       `json:"finishedDate,omitempty"`
	Pipeline           *PipelineReference      `json:"pipeline,omitempty"`
	Resources          *RunResources           `json:"resources,omitempty"`
	Result             *RunResult              `json:"result,omitempty"`
	State              *RunState               `json:"state,omitempty"`
	TemplateParameters *map[string]interface{} `json:"templateParameters,omitempty"`
	Url                *string                 `json:"url,omitempty"`
	Variables          *map[string]Variable    `json:"variables,omitempty"`
}

// Settings which influence pipeline runs.
type RunPipelineParameters struct {
	// If true, don't actually create a new run. Instead, return the final YAML document after parsing templates.
	PreviewRun *bool `json:"previewRun,omitempty"`
	// The resources the run requires.
	Resources          *RunResourcesParameters `json:"resources,omitempty"`
	StagesToSkip       *[]string               `json:"stagesToSkip,omitempty"`
	TemplateParameters *map[string]string      `json:"templateParameters,omitempty"`
	Variables          *map[string]Variable    `json:"variables,omitempty"`
	// If you use the preview run option, you may optionally supply different YAML. This allows you to preview the final YAML document without committing a changed file.
	YamlOverride *string `json:"yamlOverride,omitempty"`
}

type RunReference struct {
	Id   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type RunResources struct {
	Containers   *map[string]ContainerResource  `json:"containers,omitempty"`
	Pipelines    *map[string]PipelineResource   `json:"pipelines,omitempty"`
	Repositories *map[string]RepositoryResource `json:"repositories,omitempty"`
}

type RunResourcesParameters struct {
	Builds       *map[string]BuildResourceParameters      `json:"builds,omitempty"`
	Containers   *map[string]ContainerResourceParameters  `json:"containers,omitempty"`
	Packages     *map[string]PackageResourceParameters    `json:"packages,omitempty"`
	Pipelines    *map[string]PipelineResourceParameters   `json:"pipelines,omitempty"`
	Repositories *map[string]RepositoryResourceParameters `json:"repositories,omitempty"`
}

// This is not a Flags enum because we don't want to set multiple results on a build. However, when adding values, please stick to powers of 2 as if it were a Flags enum. This will make it easier to query multiple results.
type RunResult string

type runResultValuesType struct {
	Unknown   RunResult
	Succeeded RunResult
	Failed    RunResult
	Canceled  RunResult
}

var RunResultValues = runResultValuesType{
	Unknown:   "unknown",
	Succeeded: "succeeded",
	Failed:    "failed",
	Canceled:  "canceled",
}

// This is not a Flags enum because we don't want to set multiple states on a build. However, when adding values, please stick to powers of 2 as if it were a Flags enum. This will make it easier to query multiple states.
type RunState string

type runStateValuesType struct {
	Unknown    RunState
	InProgress RunState
	Canceling  RunState
	Completed  RunState
}

var RunStateValues = runStateValuesType{
	Unknown:    "unknown",
	InProgress: "inProgress",
	Canceling:  "canceling",
	Completed:  "completed",
}

type SignalRConnection struct {
	SignedContent *webapi.SignedUrl `json:"signedContent,omitempty"`
}

type Variable struct {
	IsSecret *bool   `json:"isSecret,omitempty"`
	Value    *string `json:"value,omitempty"`
}

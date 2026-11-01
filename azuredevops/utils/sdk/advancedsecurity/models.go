// The Azure DevOps Go SDK does not contain an Advanced Security client.
// This file cannot be under "internal", because azdosdkmocks/advancedsecurity_sdk_mock.go depends on it.

package advancedsecurity

import "github.com/google/uuid"

// RepoEnablementSettings is the Advanced Security enablement of a repository.
type RepoEnablementSettings struct {
	ProjectId                *uuid.UUID                `json:"projectId,omitempty"`
	RepositoryId             *uuid.UUID                `json:"repositoryId,omitempty"`
	SecretProtectionFeatures *SecretProtectionFeatures `json:"secretProtectionFeatures,omitempty"`
	CodeSecurityFeatures     *CodeSecurityFeatures     `json:"codeSecurityFeatures,omitempty"`
}

// SecretProtectionFeatures are the Secret Protection settings of a repository.
type SecretProtectionFeatures struct {
	SecretProtectionEnabled *bool `json:"secretProtectionEnabled,omitempty"`
	BlockPushes             *bool `json:"blockPushes,omitempty"`
}

// CodeSecurityFeatures are the Code Security settings of a repository.
type CodeSecurityFeatures struct {
	CodeSecurityEnabled                *bool `json:"codeSecurityEnabled,omitempty"`
	DependencyScanningInjectionEnabled *bool `json:"dependencyScanningInjectionEnabled,omitempty"`
	CodeQLEnabled                      *bool `json:"codeQLEnabled,omitempty"`
	AutofixEnabled                     *bool `json:"autofixEnabled,omitempty"`
}

// GetRepoEnablementArgs are the arguments for GetRepoEnablement.
type GetRepoEnablementArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Repository ID or repository name
	Repository *string
	// (optional) Also return whether pushes containing secrets are blocked
	IncludeAllProperties *bool
}

// UpdateRepoEnablementArgs are the arguments for UpdateRepoEnablement.
type UpdateRepoEnablementArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Repository ID or repository name
	Repository *string
	// (required) The enablement to apply; nil features and fields are left unchanged
	Settings *RepoEnablementSettings
}

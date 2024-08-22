package tokens

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
)

// PatTokenCreateRequest encapsulates the request parameters for creating a new personal access token (PAT)
type PatTokenCreateRequest struct {
	// True, if this personal access token (PAT) is for all of the user's accessible organizations. False, if otherwise (e.g. if the token is for a specific organization)
	AllOrgs *bool `json:"allOrgs,omitempty"`
	// The token name
	DisplayName *string `json:"displayName,omitempty"`
	// The token scopes for accessing Azure DevOps resources
	Scope *string `json:"scope,omitempty"`
	// The token expiration date. If the "Enforce maximum personal access token lifespan" policy is enabled and the provided token expiration date is past the maximum allowed lifespan, it will return back a PAT with a validTo date equal to the current date + maximum allowed lifespan.
	ValidTo *azuredevops.Time `json:"validTo,omitempty"`
}

// PatTokenResult contains the resulting personal access token (PAT) and the error (if any) that occurred during the operation
type PatTokenResult struct {
	// Represents a personal access token (PAT) used to access Azure DevOps resources
	PatToken PatToken `json:"patToken,omitempty"`
	// The error (if any) that occurred
	PatTokenError *PatTokenError `json:"patTokenError,omitempty"`
}

// PatToken represents a personal access token (PAT) used to access Azure DevOps resources
type PatToken struct {
	// Unique guid identifier
	AuthorizationId *uuid.UUID `json:"authorizationId,omitempty"`
	// The token name
	DisplayName *string `json:"displayName,omitempty"`
	// The token scopes for accessing Azure DevOps resources
	Scope *string `json:"scope,omitempty"`
	// The organizations for which the token is valid; null if the token applies to all of the user's accessible organizations
	TargetAccounts *[]uuid.UUID `json:"targetAccounts,omitempty"`
	// The unique token string generated at creation
	Token *string `json:"token,omitempty"`
	// The token creation date
	ValidFrom *azuredevops.Time `json:"validFrom,omitempty"`
	// The token expiration date
	ValidTo *azuredevops.Time `json:"validTo,omitempty"`
}

// Enumeration of possible errors returned when creating a personal access token (PAT)
type PatTokenError string

type buildPatTokenErrorType struct {
	None                        PatTokenError
	DisplayNameRequired         PatTokenError
	InvalideDisplayName         PatTokenError
	InvalidValidTo              PatTokenError
	InvalidScope                PatTokenError
	UserIdRequired              PatTokenError
	InvalidUserId               PatTokenError
	InvalidUserType             PatTokenError
	AccessDenied                PatTokenError
	FailedToIssueAccessToken    PatTokenError
	InvalidClient               PatTokenError
	InvalidClientType           PatTokenError
	InvalidClientId             PatTokenError
	InvalideTargetAccounts      PatTokenError
	HostAuthorizationNotFound   PatTokenError
	AuthorizationNotFound       PatTokenError
	FailedToUpdateAccessToken   PatTokenError
	SourceNotSupported          PatTokenError
	InvalidSourceIP             PatTokenError
	InvalideSource              PatTokenError
	DuplicateHash               PatTokenError
	SshPolicyDisabled           PatTokenError
	InvalidToken                PatTokenError
	TokenNotFound               PatTokenError
	InvalidAuthorizationId      PatTokenError
	FailedToReadTenantPolicy    PatTokenError
	GlobalPatPolicyViolation    PatTokenError
	FullScopePatPolicyViolation PatTokenError
	PatLifespanPolicyViolation  PatTokenError
	InvalidTokenType            PatTokenError
	InvalidAudience             PatTokenError
	InvalidSubject              PatTokenError
	DeploymentHostNotSupported  PatTokenError
}

var PatTokeErrorValues = buildPatTokenErrorType{
	None:                        "none",
	DisplayNameRequired:         "displayNameRequired",
	InvalideDisplayName:         "invalidDisplayName",
	InvalidValidTo:              "invalidValidTo",
	InvalidScope:                "invalidScope",
	UserIdRequired:              "userIdRequired",
	InvalidUserId:               "invalidUserId",
	InvalidUserType:             "invalidUserType",
	AccessDenied:                "accessDenied",
	FailedToIssueAccessToken:    "failedToIssueAccessToken",
	InvalidClient:               "invalidClient",
	InvalidClientType:           "invalidClientType",
	InvalidClientId:             "invalidClientId",
	InvalideTargetAccounts:      "invalidTargetAccounts",
	HostAuthorizationNotFound:   "hostAuthorizationNotFound",
	AuthorizationNotFound:       "authorizationNotFound",
	FailedToUpdateAccessToken:   "failedToUpdateAccessToken",
	SourceNotSupported:          "sourceNotSupported",
	InvalidSourceIP:             "invalidSourceIP",
	InvalideSource:              "invalidSource",
	DuplicateHash:               "duplicateHash",
	SshPolicyDisabled:           "sshPolicyDisabled",
	InvalidToken:                "invalidToken",
	TokenNotFound:               "tokenNotFound",
	InvalidAuthorizationId:      "invalidAuthorizationId",
	FailedToReadTenantPolicy:    "failedToReadTenantPolicy",
	GlobalPatPolicyViolation:    "globalPatPolicyViolation",
	FullScopePatPolicyViolation: "fullScopePatPolicyViolation",
	PatLifespanPolicyViolation:  "patLifespanPolicyViolation",
	InvalidTokenType:            "invalidTokenType",
	InvalidAudience:             "invalidAudience",
	InvalidSubject:              "invalidSubject",
	DeploymentHostNotSupported:  "deploymentHostNotSupported",
}

// PatTokenUpdateRequest encapsulates the request parameters for updating a personal access token (PAT)
type PatTokenUpdateRequest struct {
	// (Optional) True if this personal access token (PAT) is for all of the user's accessible organizations. False if otherwise (e.g. if the token is for a specific organization)
	AllOrgs *bool `json:"allOrgs,omitempty"`
	// The authorizationId identifying a single, unique personal access token (PAT)
	AuthorizationId *uuid.UUID `json:"authorizationId,omitempty"`
	// (Optional) The token name
	DisplayName *string `json:"displayName,omitempty"`
	// (Optional) The token scopes for accessing Azure DevOps resources
	Scope *string `json:"scope,omitempty"`
	// (Optional) The token expiration date. If the \"Enforce maximum personal access token lifespan\" policy is enabled and the provided token expiration date is past the maximum allowed lifespan, it will return back a PAT with a validTo date equal to the date when the PAT was intially created + maximum allowed lifespan.
	ValidTo *azuredevops.Time `json:"validTo,omitempty"`
}

// Enumerates display filter options for Personal Access Tokens (PATs)
type DisplayFilterOption string

type buildDisplayFilterOptionType struct {
	Active  DisplayFilterOption
	Revoked DisplayFilterOption
	Expired DisplayFilterOption
	All     DisplayFilterOption
}

var DisplayFilterOptionValues = buildDisplayFilterOptionType{
	Active:  "active",
	Revoked: "revoked",
	Expired: "expired",
	All:     "all",
}

// Enumerates sort by options for Personal Access Tokens (PATs)
type SortByOption string

type buildSortByOptionType struct {
	DisplayName SortByOption
	DisplayDate SortByOption
	Status      SortByOption
}

var SortByOptionValues = buildSortByOptionType{
	DisplayName: "displayName",
	DisplayDate: "displayDate",
	Status:      "status",
}

// PagedPatResults returned by the List method; contains a list of personal access tokens (PATs) and the continuation token to get the next page of results
type PagedPatResults struct {
	// "Used to access the next page of results in successive API calls to list personal access tokens (PATs)
	ContinuationToken *string `json:"continuationToken,omitempty"`
	// The list of personal access tokens (PATs)
	PatTokens *[]PatToken `json:"patTokens,omitempty"`
}

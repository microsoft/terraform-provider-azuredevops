// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package serviceendpoint

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/forminput"
	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
)

type AadLoginPromptOption string

type aadLoginPromptOptionValuesType struct {
	NoOption          AadLoginPromptOption
	Login             AadLoginPromptOption
	SelectAccount     AadLoginPromptOption
	FreshLogin        AadLoginPromptOption
	FreshLoginWithMfa AadLoginPromptOption
}

var AadLoginPromptOptionValues = aadLoginPromptOptionValuesType{
	// Do not provide a prompt option
	NoOption: "noOption",
	// Force the user to login again.
	Login: "login",
	// Force the user to select which account they are logging in with instead of automatically picking the user up from the session state. NOTE: This does not work for switching between the variants of a dual-homed user.
	SelectAccount: "selectAccount",
	// Force the user to login again. <remarks> Ignore current authentication state and force the user to authenticate again. This option should be used instead of Login. </remarks>
	FreshLogin: "freshLogin",
	// Force the user to login again with mfa. <remarks> Ignore current authentication state and force the user to authenticate again. This option should be used instead of Login, if MFA is required. </remarks>
	FreshLoginWithMfa: "freshLoginWithMfa",
}

type AadOauthTokenRequest struct {
	Refresh  *bool   `json:"refresh,omitempty"`
	Resource *string `json:"resource,omitempty"`
	TenantId *string `json:"tenantId,omitempty"`
	Token    *string `json:"token,omitempty"`
}

type AadOauthTokenResult struct {
	AccessToken       *string `json:"accessToken,omitempty"`
	RefreshTokenCache *string `json:"refreshTokenCache,omitempty"`
}

type AccessTokenRequestType string

type accessTokenRequestTypeValuesType struct {
	None   AccessTokenRequestType
	Oauth  AccessTokenRequestType
	Direct AccessTokenRequestType
}

var AccessTokenRequestTypeValues = accessTokenRequestTypeValuesType{
	None:   "none",
	Oauth:  "oauth",
	Direct: "direct",
}

type AuthConfiguration struct {
	// Gets or sets the ClientId
	ClientId *string `json:"clientId,omitempty"`
	// Gets or sets the ClientSecret
	ClientSecret *string `json:"clientSecret,omitempty"`
	// Gets or sets the identity who created the config.
	CreatedBy *webapi.IdentityRef `json:"createdBy,omitempty"`
	// Gets or sets the time when config was created.
	CreatedOn *azuredevops.Time `json:"createdOn,omitempty"`
	// Gets or sets the type of the endpoint.
	EndpointType *string `json:"endpointType,omitempty"`
	// Gets or sets the unique identifier of this field
	Id *uuid.UUID `json:"id,omitempty"`
	// Gets or sets the identity who modified the config.
	ModifiedBy *webapi.IdentityRef `json:"modifiedBy,omitempty"`
	// Gets or sets the time when variable group was modified
	ModifiedOn *azuredevops.Time `json:"modifiedOn,omitempty"`
	// Gets or sets the name
	Name *string `json:"name,omitempty"`
	// Gets or sets the Url
	Url *string `json:"url,omitempty"`
	// Gets or sets parameters contained in configuration object.
	Parameters *map[string]Parameter `json:"parameters,omitempty"`
}

// Specifies the authentication scheme to be used for authentication.
type AuthenticationSchemeReference struct {
	// Gets or sets the key and value of the fields used for authentication.
	Inputs *map[string]string `json:"inputs,omitempty"`
	// Gets or sets the type of authentication scheme of an endpoint.
	Type *string `json:"type,omitempty"`
}

// Represents the header of the REST request.
type AuthorizationHeader struct {
	// Gets or sets the name of authorization header.
	Name *string `json:"name,omitempty"`
	// Gets or sets the value of authorization header.
	Value *string `json:"value,omitempty"`
}

type AzureKeyVaultPermission struct {
	Provisioned      *bool   `json:"provisioned,omitempty"`
	ResourceProvider *string `json:"resourceProvider,omitempty"`
	ResourceGroup    *string `json:"resourceGroup,omitempty"`
	Vault            *string `json:"vault,omitempty"`
}

// Azure Management Group
type AzureManagementGroup struct {
	// Display name of azure management group
	DisplayName *string `json:"displayName,omitempty"`
	// Id of azure management group
	Id *string `json:"id,omitempty"`
	// Azure management group name
	Name *string `json:"name,omitempty"`
	// Id of tenant from which azure management group belogs
	TenantId *string `json:"tenantId,omitempty"`
}

// Azure management group query result
type AzureManagementGroupQueryResult struct {
	// Error message in case of an exception
	ErrorMessage *string `json:"errorMessage,omitempty"`
	// List of azure management groups
	Value *[]AzureManagementGroup `json:"value,omitempty"`
}

type AzureMLWorkspace struct {
	Id       *string `json:"id,omitempty"`
	Location *string `json:"location,omitempty"`
	Name     *string `json:"name,omitempty"`
}

type AzurePermission struct {
	Provisioned      *bool   `json:"provisioned,omitempty"`
	ResourceProvider *string `json:"resourceProvider,omitempty"`
}

type AzureResourcePermission struct {
	Provisioned      *bool   `json:"provisioned,omitempty"`
	ResourceProvider *string `json:"resourceProvider,omitempty"`
	ResourceGroup    *string `json:"resourceGroup,omitempty"`
}

type AzureRoleAssignmentPermission struct {
	Provisioned      *bool      `json:"provisioned,omitempty"`
	ResourceProvider *string    `json:"resourceProvider,omitempty"`
	RoleAssignmentId *uuid.UUID `json:"roleAssignmentId,omitempty"`
}

type AzureSpnOperationStatus struct {
	State         *string `json:"state,omitempty"`
	StatusMessage *string `json:"statusMessage,omitempty"`
}

type AzureSubscription struct {
	DisplayName            *string `json:"displayName,omitempty"`
	SubscriptionId         *string `json:"subscriptionId,omitempty"`
	SubscriptionTenantId   *string `json:"subscriptionTenantId,omitempty"`
	SubscriptionTenantName *string `json:"subscriptionTenantName,omitempty"`
}

type AzureSubscriptionQueryResult struct {
	ErrorMessage *string              `json:"errorMessage,omitempty"`
	Value        *[]AzureSubscription `json:"value,omitempty"`
}

// Specifies the client certificate to be used for the endpoint request.
type ClientCertificate struct {
	// Gets or sets the value of client certificate.
	Value *string `json:"value,omitempty"`
}

// Specifies the data sources for this endpoint.
type DataSource struct {
	// Gets or sets the authentication scheme for the endpoint request.
	AuthenticationScheme *AuthenticationSchemeReference `json:"authenticationScheme,omitempty"`
	// Gets or sets the pagination format supported by this data source(ContinuationToken/SkipTop).
	CallbackContextTemplate *string `json:"callbackContextTemplate,omitempty"`
	// Gets or sets the template to check if subsequent call is needed.
	CallbackRequiredTemplate *string `json:"callbackRequiredTemplate,omitempty"`
	// Gets or sets the endpoint url of the data source.
	EndpointUrl *string `json:"endpointUrl,omitempty"`
	// Gets or sets the authorization headers of the request.
	Headers *[]AuthorizationHeader `json:"headers,omitempty"`
	// Gets or sets the initial value of the query params.
	InitialContextTemplate *string `json:"initialContextTemplate,omitempty"`
	// Gets or sets the name of the data source.
	Name *string `json:"name,omitempty"`
	// Gets or sets the request content of the endpoint request.
	RequestContent *string `json:"requestContent,omitempty"`
	// Gets or sets the request method of the endpoint request.
	RequestVerb *string `json:"requestVerb,omitempty"`
	// Gets or sets the resource url of the endpoint request.
	ResourceUrl *string `json:"resourceUrl,omitempty"`
	// Gets or sets the result selector to filter the response of the endpoint request.
	ResultSelector *string `json:"resultSelector,omitempty"`
}

// Represents the data source binding of the endpoint.
type DataSourceBinding struct {
}

// Represents details of the service endpoint data source.
type DataSourceDetails struct {
	// Gets or sets the data source name.
	DataSourceName *string `json:"dataSourceName,omitempty"`
	// Gets or sets the data source url.
	DataSourceUrl *string `json:"dataSourceUrl,omitempty"`
	// Gets or sets the request headers.
	Headers *[]AuthorizationHeader `json:"headers,omitempty"`
	// Gets or sets the initialization context used for the initial call to the data source
	InitialContextTemplate *string `json:"initialContextTemplate,omitempty"`
	// Gets the parameters of data source.
	Parameters *map[string]string `json:"parameters,omitempty"`
	// Gets or sets the data source request content.
	RequestContent *string `json:"requestContent,omitempty"`
	// Gets or sets the data source request verb. Get/Post are the only implemented types
	RequestVerb *string `json:"requestVerb,omitempty"`
	// Gets or sets the resource url of data source.
	ResourceUrl *string `json:"resourceUrl,omitempty"`
	// Gets or sets the result selector.
	ResultSelector *string `json:"resultSelector,omitempty"`
}

// Represents the details of the input on which a given input is dependent.
type DependencyBinding struct {
	// Gets or sets the value of the field on which url is dependent.
	Key *string `json:"key,omitempty"`
	// Gets or sets the corresponding value of url.
	Value *string `json:"value,omitempty"`
}

// Represents the dependency data for the endpoint inputs.
type DependencyData struct {
	// Gets or sets the category of dependency data.
	Input *string `json:"input,omitempty"`
	// Gets or sets the key-value pair to specify properties and their values.
	Map *[]azuredevops.KeyValuePair `json:"map,omitempty"`
}

// Represents the inputs on which any given input is dependent.
type DependsOn struct {
	// Gets or sets the ID of the field on which URL's value is dependent.
	Input *string `json:"input,omitempty"`
	// Gets or sets key-value pair containing other's field value and corresponding url value.
	Map *[]DependencyBinding `json:"map,omitempty"`
}

// Represents the authorization used for service endpoint.
type EndpointAuthorization struct {
	// Gets or sets the parameters for the selected authorization scheme.
	Parameters *map[string]string `json:"parameters,omitempty"`
	// Gets or sets the scheme used for service endpoint authentication.
	Scheme *string `json:"scheme,omitempty"`
}

type EndpointOperationStatus struct {
	State         *string `json:"state,omitempty"`
	StatusMessage *string `json:"statusMessage,omitempty"`
}

// Represents url of the service endpoint.
type EndpointUrl struct {
	// Gets or sets the dependency bindings.
	DependsOn *DependsOn `json:"dependsOn,omitempty"`
	// Gets or sets the display name of service endpoint url.
	DisplayName *string `json:"displayName,omitempty"`
	// Gets or sets the help text of service endpoint url.
	HelpText *string `json:"helpText,omitempty"`
	// Gets or sets the visibility of service endpoint url.
	IsVisible *string `json:"isVisible,omitempty"`
	// Gets or sets the value of service endpoint url.
	Value *string `json:"value,omitempty"`
}

// Specifies the public url of the help documentation.
type HelpLink struct {
	// Gets or sets the help text.
	Text *string `json:"text,omitempty"`
	// Gets or sets the public url of the help documentation.
	Url *string `json:"url,omitempty"`
}

type OAuth2TokenResult struct {
	AccessToken      *string `json:"accessToken,omitempty"`
	Error            *string `json:"error,omitempty"`
	ErrorDescription *string `json:"errorDescription,omitempty"`
	ExpiresIn        *string `json:"expiresIn,omitempty"`
	IssuedAt         *string `json:"issuedAt,omitempty"`
	RefreshToken     *string `json:"refreshToken,omitempty"`
	Scope            *string `json:"scope,omitempty"`
}

type OAuthConfiguration struct {
	// Gets or sets the ClientId
	ClientId *string `json:"clientId,omitempty"`
	// Gets or sets the ClientSecret
	ClientSecret *string `json:"clientSecret,omitempty"`
	// Gets or sets the identity who created the config.
	CreatedBy *webapi.IdentityRef `json:"createdBy,omitempty"`
	// Gets or sets the time when config was created.
	CreatedOn *azuredevops.Time `json:"createdOn,omitempty"`
	// Gets or sets the type of the endpoint.
	EndpointType *string `json:"endpointType,omitempty"`
	// Gets or sets the unique identifier of this field
	Id *uuid.UUID `json:"id,omitempty"`
	// Gets or sets the identity who modified the config.
	ModifiedBy *webapi.IdentityRef `json:"modifiedBy,omitempty"`
	// Gets or sets the time when variable group was modified
	ModifiedOn *azuredevops.Time `json:"modifiedOn,omitempty"`
	// Gets or sets the name
	Name *string `json:"name,omitempty"`
	// Gets or sets the Url
	Url *string `json:"url,omitempty"`
}

// [Flags]
type OAuthConfigurationActionFilter string

type oAuthConfigurationActionFilterValuesType struct {
	None   OAuthConfigurationActionFilter
	Manage OAuthConfigurationActionFilter
	Use    OAuthConfigurationActionFilter
}

var OAuthConfigurationActionFilterValues = oAuthConfigurationActionFilterValuesType{
	None:   "none",
	Manage: "manage",
	Use:    "use",
}

type OAuthConfigurationParams struct {
	// Gets or sets the ClientId
	ClientId *string `json:"clientId,omitempty"`
	// Gets or sets the ClientSecret
	ClientSecret *string `json:"clientSecret,omitempty"`
	// Gets or sets the type of the endpoint.
	EndpointType *string `json:"endpointType,omitempty"`
	// Gets or sets the name
	Name *string `json:"name,omitempty"`
	// Gets or sets the Url
	Url *string `json:"url,omitempty"`
}

type OAuthEndpointStatus struct {
	State         *string `json:"state,omitempty"`
	StatusMessage *string `json:"statusMessage,omitempty"`
}

type Parameter struct {
	IsSecret *bool   `json:"isSecret,omitempty"`
	Value    *string `json:"value,omitempty"`
}

type ProjectReference struct {
	Id   *uuid.UUID `json:"id,omitempty"`
	Name *string    `json:"name,omitempty"`
}

// Represents template to transform the result data.
type ResultTransformationDetails struct {
	// Gets or sets the template for callback parameters
	CallbackContextTemplate *string `json:"callbackContextTemplate,omitempty"`
	// Gets or sets the template to decide whether to callback or not
	CallbackRequiredTemplate *string `json:"callbackRequiredTemplate,omitempty"`
	// Gets or sets the template for result transformation.
	ResultTemplate *string `json:"resultTemplate,omitempty"`
}

// Represents an endpoint which may be used by an orchestration job.
type ServiceEndpoint struct {
	// Gets or sets the identity reference for the administrators group of the service endpoint.
	AdministratorsGroup *webapi.IdentityRef `json:"administratorsGroup,omitempty"`
	// Gets or sets the authorization data for talking to the endpoint.
	Authorization *EndpointAuthorization `json:"authorization,omitempty"`
	// Gets or sets the identity reference for the user who created the Service endpoint.
	CreatedBy *webapi.IdentityRef `json:"createdBy,omitempty"`
	Data      *map[string]string  `json:"data,omitempty"`
	// Gets or sets the description of endpoint.
	Description *string `json:"description,omitempty"`
	// This is a deprecated field.
	GroupScopeId *uuid.UUID `json:"groupScopeId,omitempty"`
	// Gets or sets the identifier of this endpoint.
	Id *uuid.UUID `json:"id,omitempty"`
	// EndPoint state indicator
	IsReady *bool `json:"isReady,omitempty"`
	// Indicates whether service endpoint is shared with other projects or not.
	IsShared *bool `json:"isShared,omitempty"`
	// Gets or sets the friendly name of the endpoint.
	Name *string `json:"name,omitempty"`
	// Error message during creation/deletion of endpoint
	OperationStatus interface{} `json:"operationStatus,omitempty"`
	// Owner of the endpoint Supported values are "library", "agentcloud"
	Owner *string `json:"owner,omitempty"`
	// Gets or sets the identity reference for the readers group of the service endpoint.
	ReadersGroup *webapi.IdentityRef `json:"readersGroup,omitempty"`
	// Gets or sets the type of the endpoint.
	Type *string `json:"type,omitempty"`
	// Gets or sets the url of the endpoint.
	Url *string `json:"url,omitempty"`
}

// [Flags]
type ServiceEndpointActionFilter string

type serviceEndpointActionFilterValuesType struct {
	None   ServiceEndpointActionFilter
	Manage ServiceEndpointActionFilter
	Use    ServiceEndpointActionFilter
}

var ServiceEndpointActionFilterValues = serviceEndpointActionFilterValuesType{
	None:   "none",
	Manage: "manage",
	Use:    "use",
}

// Represents the authentication scheme used to authenticate the endpoint.
type ServiceEndpointAuthenticationScheme struct {
	// Gets or sets the authorization headers of service endpoint authentication scheme.
	AuthorizationHeaders *[]AuthorizationHeader `json:"authorizationHeaders,omitempty"`
	// Gets or sets the Authorization url required to authenticate using OAuth2
	AuthorizationUrl *string `json:"authorizationUrl,omitempty"`
	// Gets or sets the certificates of service endpoint authentication scheme.
	ClientCertificates *[]ClientCertificate `json:"clientCertificates,omitempty"`
	// Gets or sets the data source bindings of the endpoint.
	DataSourceBindings *[]DataSourceBinding `json:"dataSourceBindings,omitempty"`
	// Gets or sets the display name for the service endpoint authentication scheme.
	DisplayName *string `json:"displayName,omitempty"`
	// Gets or sets the input descriptors for the service endpoint authentication scheme.
	InputDescriptors *[]forminput.InputDescriptor `json:"inputDescriptors,omitempty"`
	// Gets or sets the scheme for service endpoint authentication.
	Scheme *string `json:"scheme,omitempty"`
}

// Represents details of the service endpoint.
type ServiceEndpointDetails struct {
	// Gets or sets the authorization of service endpoint.
	Authorization *EndpointAuthorization `json:"authorization,omitempty"`
	// Gets or sets the data of service endpoint.
	Data *map[string]string `json:"data,omitempty"`
	// Gets or sets the type of service endpoint.
	Type *string `json:"type,omitempty"`
	// Gets or sets the connection url of service endpoint.
	Url *string `json:"url,omitempty"`
}

// Represents service endpoint execution data.
type ServiceEndpointExecutionData struct {
	// Gets the definition of service endpoint execution owner.
	Definition *ServiceEndpointExecutionOwner `json:"definition,omitempty"`
	// Gets the finish time of service endpoint execution.
	FinishTime *azuredevops.Time `json:"finishTime,omitempty"`
	// Gets the Id of service endpoint execution data.
	Id *uint64 `json:"id,omitempty"`
	// Gets the owner of service endpoint execution data.
	Owner *ServiceEndpointExecutionOwner `json:"owner,omitempty"`
	// Gets the plan type of service endpoint execution data.
	PlanType *string `json:"planType,omitempty"`
	// Gets the result of service endpoint execution.
	Result *ServiceEndpointExecutionResult `json:"result,omitempty"`
	// Gets the start time of service endpoint execution.
	StartTime *azuredevops.Time `json:"startTime,omitempty"`
}

// Represents execution owner of the service endpoint.
type ServiceEndpointExecutionOwner struct {
	Links interface{} `json:"_links,omitempty"`
	// Gets or sets the Id of service endpoint execution owner.
	Id *int `json:"id,omitempty"`
	// Gets or sets the name of service endpoint execution owner.
	Name *string `json:"name,omitempty"`
}

// Represents the details of service endpoint execution.
type ServiceEndpointExecutionRecord struct {
	// Gets the execution data of service endpoint execution.
	Data *ServiceEndpointExecutionData `json:"data,omitempty"`
	// Gets the Id of service endpoint.
	EndpointId *uuid.UUID `json:"endpointId,omitempty"`
}

type ServiceEndpointExecutionRecordsInput struct {
	Data        *ServiceEndpointExecutionData `json:"data,omitempty"`
	EndpointIds *[]uuid.UUID                  `json:"endpointIds,omitempty"`
}

type ServiceEndpointExecutionResult string

type serviceEndpointExecutionResultValuesType struct {
	Succeeded           ServiceEndpointExecutionResult
	SucceededWithIssues ServiceEndpointExecutionResult
	Failed              ServiceEndpointExecutionResult
	Canceled            ServiceEndpointExecutionResult
	Skipped             ServiceEndpointExecutionResult
	Abandoned           ServiceEndpointExecutionResult
}

var ServiceEndpointExecutionResultValues = serviceEndpointExecutionResultValuesType{
	// "Service endpoint request succeeded.
	Succeeded: "succeeded",
	// "Service endpoint request succeeded but with some issues.
	SucceededWithIssues: "succeededWithIssues",
	// "Service endpoint request failed.
	Failed: "failed",
	// "Service endpoint request was cancelled.
	Canceled: "canceled",
	// "Service endpoint request was skipped.
	Skipped: "skipped",
	// "Service endpoint request was abandoned.
	Abandoned: "abandoned",
}

type ServiceEndpointOAuthConfigurationReference struct {
	ConfigurationId          *uuid.UUID `json:"configurationId,omitempty"`
	ServiceEndpointId        *uuid.UUID `json:"serviceEndpointId,omitempty"`
	ServiceEndpointProjectId *uuid.UUID `json:"serviceEndpointProjectId,omitempty"`
}

type ServiceEndpointRequest struct {
	// Gets or sets the data source details for the service endpoint request.
	DataSourceDetails *DataSourceDetails `json:"dataSourceDetails,omitempty"`
	// Gets or sets the result transformation details for the service endpoint request.
	ResultTransformationDetails *ResultTransformationDetails `json:"resultTransformationDetails,omitempty"`
	// Gets or sets the service endpoint details for the service endpoint request.
	ServiceEndpointDetails *ServiceEndpointDetails `json:"serviceEndpointDetails,omitempty"`
}

// Represents result of the service endpoint request.
type ServiceEndpointRequestResult struct {
	// Gets or sets the parameters used to make subsequent calls to the data source
	CallbackContextParameters *map[string]string `json:"callbackContextParameters,omitempty"`
	// Gets or sets the flat that decides if another call to the data source is to be made
	CallbackRequired *bool `json:"callbackRequired,omitempty"`
	// Gets or sets the error message of the service endpoint request result.
	ErrorMessage *string `json:"errorMessage,omitempty"`
	// Gets or sets the result of service endpoint request.
	Result interface{} `json:"result,omitempty"`
	// Gets or sets the status code of the service endpoint request result.
	StatusCode *string `json:"statusCode,omitempty"`
}

// Represents type of the service endpoint.
type ServiceEndpointType struct {
	// Authentication scheme of service endpoint type.
	AuthenticationSchemes *[]ServiceEndpointAuthenticationScheme `json:"authenticationSchemes,omitempty"`
	// Data sources of service endpoint type.
	DataSources *[]DataSource `json:"dataSources,omitempty"`
	// Dependency data of service endpoint type.
	DependencyData *[]DependencyData `json:"dependencyData,omitempty"`
	// Gets or sets the description of service endpoint type.
	Description *string `json:"description,omitempty"`
	// Gets or sets the display name of service endpoint type.
	DisplayName *string `json:"displayName,omitempty"`
	// Gets or sets the endpoint url of service endpoint type.
	EndpointUrl *EndpointUrl `json:"endpointUrl,omitempty"`
	// Gets or sets the help link of service endpoint type.
	HelpLink *HelpLink `json:"helpLink,omitempty"`
	// Gets or sets the help text shown at the endpoint create dialog.
	HelpMarkDown *string `json:"helpMarkDown,omitempty"`
	// Gets or sets the icon url of service endpoint type.
	IconUrl *string `json:"iconUrl,omitempty"`
	// Input descriptor of service endpoint type.
	InputDescriptors *[]forminput.InputDescriptor `json:"inputDescriptors,omitempty"`
	// Gets or sets the name of service endpoint type.
	Name *string `json:"name,omitempty"`
	// Trusted hosts of a service endpoint type.
	TrustedHosts *[]string `json:"trustedHosts,omitempty"`
	// Gets or sets the ui contribution id of service endpoint type.
	UiContributionId *string `json:"uiContributionId,omitempty"`
}

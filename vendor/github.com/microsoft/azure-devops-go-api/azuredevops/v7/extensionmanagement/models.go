// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package extensionmanagement

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/gallery"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
)

// How the acquisition is assigned
type AcquisitionAssignmentType string

type acquisitionAssignmentTypeValuesType struct {
	None AcquisitionAssignmentType
	Me   AcquisitionAssignmentType
	All  AcquisitionAssignmentType
}

var AcquisitionAssignmentTypeValues = acquisitionAssignmentTypeValuesType{
	None: "none",
	// Just assign for me
	Me: "me",
	// Assign for all users in the account
	All: "all",
}

type AcquisitionOperation struct {
	// State of the AcquisitionOperation for the current user
	OperationState *AcquisitionOperationState `json:"operationState,omitempty"`
	// AcquisitionOperationType: install, request, buy, etc...
	OperationType *AcquisitionOperationType `json:"operationType,omitempty"`
	// Optional reason to justify current state. Typically used with Disallow state.
	Reason *string `json:"reason,omitempty"`
	// List of reasons indicating why the operation is not allowed.
	Reasons *[]AcquisitionOperationDisallowReason `json:"reasons,omitempty"`
}

type AcquisitionOperationDisallowReason struct {
	// User-friendly message clarifying the reason for disallowance
	Message *string `json:"message,omitempty"`
	// Type of reason for disallowance - AlreadyInstalled, UnresolvedDemand, etc.
	Type *string `json:"type,omitempty"`
}

type AcquisitionOperationState string

type acquisitionOperationStateValuesType struct {
	Disallow  AcquisitionOperationState
	Allow     AcquisitionOperationState
	Completed AcquisitionOperationState
}

var AcquisitionOperationStateValues = acquisitionOperationStateValuesType{
	// Not allowed to use this AcquisitionOperation
	Disallow: "disallow",
	// Allowed to use this AcquisitionOperation
	Allow: "allow",
	// Operation has already been completed and is no longer available
	Completed: "completed",
}

// Set of different types of operations that can be requested.
type AcquisitionOperationType string

type acquisitionOperationTypeValuesType struct {
	Get             AcquisitionOperationType
	Install         AcquisitionOperationType
	Buy             AcquisitionOperationType
	Try             AcquisitionOperationType
	Request         AcquisitionOperationType
	None            AcquisitionOperationType
	PurchaseRequest AcquisitionOperationType
}

var AcquisitionOperationTypeValues = acquisitionOperationTypeValuesType{
	// Not yet used
	Get: "get",
	// Install this extension into the host provided
	Install: "install",
	// Buy licenses for this extension and install into the host provided
	Buy: "buy",
	// Try this extension
	Try: "try",
	// Request this extension for installation
	Request: "request",
	// No action found
	None: "none",
	// Request admins for purchasing extension
	PurchaseRequest: "purchaseRequest",
}

// Market item acquisition options (install, buy, etc) for an installation target.
type AcquisitionOptions struct {
	// Default Operation for the ItemId in this target
	DefaultOperation *AcquisitionOperation `json:"defaultOperation,omitempty"`
	// The item id that this options refer to
	ItemId *string `json:"itemId,omitempty"`
	// Operations allowed for the ItemId in this target
	Operations *[]AcquisitionOperation `json:"operations,omitempty"`
	// Additional properties which can be added to the request.
	Properties interface{} `json:"properties,omitempty"`
	// The target that this options refer to
	Target *string `json:"target,omitempty"`
}

// Representation of a ContributionNode that can be used for serialized to clients.
type ClientContribution struct {
	// Description of the contribution/type
	Description *string `json:"description,omitempty"`
	// Fully qualified identifier of the contribution/type
	Id *string `json:"id,omitempty"`
	// Includes is a set of contributions that should have this contribution included in their targets list.
	Includes *[]string `json:"includes,omitempty"`
	// Properties/attributes of this contribution
	Properties interface{} `json:"properties,omitempty"`
	// The ids of the contribution(s) that this contribution targets. (parent contributions)
	Targets *[]string `json:"targets,omitempty"`
	// Id of the Contribution Type
	Type *string `json:"type,omitempty"`
}

// Representation of a ContributionNode that can be used for serialized to clients.
type ClientContributionNode struct {
	// Contribution associated with this node.
	Contribution *ClientContribution `json:"contribution,omitempty"`
	// List of ids for contributions which are children to the current contribution.
	Children *[]string `json:"children,omitempty"`
	// List of ids for contributions which are parents to the current contribution.
	Parents *[]string `json:"parents,omitempty"`
}

type ClientContributionProviderDetails struct {
	// Friendly name for the provider.
	DisplayName *string `json:"displayName,omitempty"`
	// Unique identifier for this provider. The provider name can be used to cache the contribution data and refer back to it when looking for changes
	Name *string `json:"name,omitempty"`
	// Properties associated with the provider
	Properties *map[string]string `json:"properties,omitempty"`
	// Version of contributions associated with this contribution provider.
	Version *string `json:"version,omitempty"`
}

// A client data provider are the details needed to make the data provider request from the client.
type ClientDataProviderQuery struct {
	// Contextual information to pass to the data providers
	Context *DataProviderContext `json:"context,omitempty"`
	// The contribution ids of the data providers to resolve
	ContributionIds *[]string `json:"contributionIds,omitempty"`
	// The Id of the service instance type that should be communicated with in order to resolve the data providers from the client given the query values.
	QueryServiceInstanceType *uuid.UUID `json:"queryServiceInstanceType,omitempty"`
}

// An individual contribution made by an extension
type Contribution struct {
	// Description of the contribution/type
	Description *string `json:"description,omitempty"`
	// Fully qualified identifier of the contribution/type
	Id *string `json:"id,omitempty"`
	// VisibleTo can be used to restrict whom can reference a given contribution/type. This value should be a list of publishers or extensions access is restricted too.  Examples: "ms" - Means only the "ms" publisher can reference this. "ms.vss-web" - Means only the "vss-web" extension from the "ms" publisher can reference this.
	VisibleTo *[]string `json:"visibleTo,omitempty"`
	// List of constraints (filters) that should be applied to the availability of this contribution
	Constraints *[]ContributionConstraint `json:"constraints,omitempty"`
	// Includes is a set of contributions that should have this contribution included in their targets list.
	Includes *[]string `json:"includes,omitempty"`
	// Properties/attributes of this contribution
	Properties interface{} `json:"properties,omitempty"`
	// List of demanded claims in order for the user to see this contribution (like anonymous, public, member...).
	RestrictedTo *[]string `json:"restrictedTo,omitempty"`
	// The ids of the contribution(s) that this contribution targets. (parent contributions)
	Targets *[]string `json:"targets,omitempty"`
	// Id of the Contribution Type
	Type *string `json:"type,omitempty"`
}

// Base class shared by contributions and contribution types
type ContributionBase struct {
	// Description of the contribution/type
	Description *string `json:"description,omitempty"`
	// Fully qualified identifier of the contribution/type
	Id *string `json:"id,omitempty"`
	// VisibleTo can be used to restrict whom can reference a given contribution/type. This value should be a list of publishers or extensions access is restricted too.  Examples: "ms" - Means only the "ms" publisher can reference this. "ms.vss-web" - Means only the "vss-web" extension from the "ms" publisher can reference this.
	VisibleTo *[]string `json:"visibleTo,omitempty"`
}

// Specifies a constraint that can be used to dynamically include/exclude a given contribution
type ContributionConstraint struct {
	// An optional property that can be specified to group constraints together. All constraints within a group are AND'd together (all must be evaluate to True in order for the contribution to be included). Different groups of constraints are OR'd (only one group needs to evaluate to True for the contribution to be included).
	Group *int `json:"group,omitempty"`
	// Fully qualified identifier of a shared constraint
	Id *string `json:"id,omitempty"`
	// If true, negate the result of the filter (include the contribution if the applied filter returns false instead of true)
	Inverse *bool `json:"inverse,omitempty"`
	// Name of the IContributionFilter plugin
	Name *string `json:"name,omitempty"`
	// Properties that are fed to the contribution filter class
	Properties interface{} `json:"properties,omitempty"`
	// Constraints can be optionally be applied to one or more of the relationships defined in the contribution. If no relationships are defined then all relationships are associated with the constraint. This means the default behaviour will eliminate the contribution from the tree completely if the constraint is applied.
	Relationships *[]string `json:"relationships,omitempty"`
}

// Represents different ways of including contributions based on licensing
type ContributionLicensingBehaviorType string

type contributionLicensingBehaviorTypeValuesType struct {
	OnlyIfLicensed   ContributionLicensingBehaviorType
	OnlyIfUnlicensed ContributionLicensingBehaviorType
	AlwaysInclude    ContributionLicensingBehaviorType
}

var ContributionLicensingBehaviorTypeValues = contributionLicensingBehaviorTypeValuesType{
	// Default value - only include the contribution if the user is licensed for the extension
	OnlyIfLicensed: "onlyIfLicensed",
	// Only include the contribution if the user is NOT licensed for the extension
	OnlyIfUnlicensed: "onlyIfUnlicensed",
	// Always include the contribution regardless of whether or not the user is licensed for the extension
	AlwaysInclude: "alwaysInclude",
}

// A query that can be issued for contribution nodes
type ContributionNodeQuery struct {
	// The contribution ids of the nodes to find.
	ContributionIds *[]string `json:"contributionIds,omitempty"`
	// Contextual information that can be leveraged by contribution constraints
	DataProviderContext *DataProviderContext `json:"dataProviderContext,omitempty"`
	// Indicator if contribution provider details should be included in the result.
	IncludeProviderDetails *bool `json:"includeProviderDetails,omitempty"`
	// Query options tpo be used when fetching ContributionNodes
	QueryOptions *ContributionQueryOptions `json:"queryOptions,omitempty"`
}

// Result of a contribution node query.  Wraps the resulting contribution nodes and provider details.
type ContributionNodeQueryResult struct {
	// Map of contribution ids to corresponding node.
	Nodes *map[string]ClientContributionNode `json:"nodes,omitempty"`
	// Map of provider ids to the corresponding provider details object.
	ProviderDetails *map[string]ClientContributionProviderDetails `json:"providerDetails,omitempty"`
}

// Description about a property of a contribution type
type ContributionPropertyDescription struct {
	// Description of the property
	Description *string `json:"description,omitempty"`
	// Name of the property
	Name *string `json:"name,omitempty"`
	// True if this property is required
	Required *bool `json:"required,omitempty"`
	// The type of value used for this property
	Type *ContributionPropertyType `json:"type,omitempty"`
}

// [Flags] The type of value used for a property
type ContributionPropertyType string

type contributionPropertyTypeValuesType struct {
	Unknown    ContributionPropertyType
	String     ContributionPropertyType
	Uri        ContributionPropertyType
	Guid       ContributionPropertyType
	Boolean    ContributionPropertyType
	Integer    ContributionPropertyType
	Double     ContributionPropertyType
	DateTime   ContributionPropertyType
	Dictionary ContributionPropertyType
	Array      ContributionPropertyType
	Object     ContributionPropertyType
}

var ContributionPropertyTypeValues = contributionPropertyTypeValuesType{
	// Contribution type is unknown (value may be anything)
	Unknown: "unknown",
	// Value is a string
	String: "string",
	// Value is a Uri
	Uri: "uri",
	// Value is a GUID
	Guid: "guid",
	// Value is True or False
	Boolean: "boolean",
	// Value is an integer
	Integer: "integer",
	// Value is a double
	Double: "double",
	// Value is a DateTime object
	DateTime: "dateTime",
	// Value is a generic Dictionary/JObject/property bag
	Dictionary: "dictionary",
	// Value is an array
	Array: "array",
	// Value is an arbitrary/custom object
	Object: "object",
}

type ContributionProviderDetails struct {
	// Friendly name for the provider.
	DisplayName *string `json:"displayName,omitempty"`
	// Unique identifier for this provider. The provider name can be used to cache the contribution data and refer back to it when looking for changes
	Name *string `json:"name,omitempty"`
	// Properties associated with the provider
	Properties *map[string]string `json:"properties,omitempty"`
	// Version of contributions associated with this contribution provider.
	Version *string `json:"version,omitempty"`
}

// [Flags] Options that control the contributions to include in a query
type ContributionQueryOptions string

type contributionQueryOptionsValuesType struct {
	None              ContributionQueryOptions
	IncludeSelf       ContributionQueryOptions
	IncludeChildren   ContributionQueryOptions
	IncludeSubTree    ContributionQueryOptions
	IncludeAll        ContributionQueryOptions
	IgnoreConstraints ContributionQueryOptions
}

var ContributionQueryOptionsValues = contributionQueryOptionsValuesType{
	None: "none",
	// Include the direct contributions that have the ids queried.
	IncludeSelf: "includeSelf",
	// Include the contributions that directly target the contributions queried.
	IncludeChildren: "includeChildren",
	// Include the contributions from the entire sub-tree targeting the contributions queried.
	IncludeSubTree: "includeSubTree",
	// Include the contribution being queried as well as all contributions that target them recursively.
	IncludeAll: "includeAll",
	// Some callers may want the entire tree back without constraint evaluation being performed.
	IgnoreConstraints: "ignoreConstraints",
}

// A contribution type, given by a json schema
type ContributionType struct {
	// Description of the contribution/type
	Description *string `json:"description,omitempty"`
	// Fully qualified identifier of the contribution/type
	Id *string `json:"id,omitempty"`
	// VisibleTo can be used to restrict whom can reference a given contribution/type. This value should be a list of publishers or extensions access is restricted too.  Examples: "ms" - Means only the "ms" publisher can reference this. "ms.vss-web" - Means only the "vss-web" extension from the "ms" publisher can reference this.
	VisibleTo *[]string `json:"visibleTo,omitempty"`
	// Controls whether or not contributions of this type have the type indexed for queries. This allows clients to find all extensions that have a contribution of this type.  NOTE: Only TrustedPartners are allowed to specify indexed contribution types.
	Indexed *bool `json:"indexed,omitempty"`
	// Friendly name of the contribution/type
	Name *string `json:"name,omitempty"`
	// Describes the allowed properties for this contribution type
	Properties *map[string]ContributionPropertyDescription `json:"properties,omitempty"`
}

// Contextual information that data providers can examine when populating their data
type DataProviderContext struct {
	// Generic property bag that contains context-specific properties that data providers can use when populating their data dictionary
	Properties *map[string]interface{} `json:"properties,omitempty"`
}

type DataProviderExceptionDetails struct {
	// The type of the exception that was thrown.
	ExceptionType *string `json:"exceptionType,omitempty"`
	// Message that is associated with the exception.
	Message *string `json:"message,omitempty"`
	// The StackTrace from the exception turned into a string.
	StackTrace *string `json:"stackTrace,omitempty"`
}

// A query that can be issued for data provider data
type DataProviderQuery struct {
	// Contextual information to pass to the data providers
	Context *DataProviderContext `json:"context,omitempty"`
	// The contribution ids of the data providers to resolve
	ContributionIds *[]string `json:"contributionIds,omitempty"`
}

// Result structure from calls to GetDataProviderData
type DataProviderResult struct {
	// This is the set of data providers that were requested, but either they were defined as client providers, or as remote providers that failed and may be retried by the client.
	ClientProviders *map[string]ClientDataProviderQuery `json:"clientProviders,omitempty"`
	// Property bag of data keyed off of the data provider contribution id
	Data *map[string]interface{} `json:"data,omitempty"`
	// Set of exceptions that occurred resolving the data providers.
	Exceptions *map[string]DataProviderExceptionDetails `json:"exceptions,omitempty"`
	// List of data providers resolved in the data-provider query
	ResolvedProviders *[]ResolvedDataProvider `json:"resolvedProviders,omitempty"`
	// Scope name applied to this data provider result.
	ScopeName *string `json:"scopeName,omitempty"`
	// Scope value applied to this data provider result.
	ScopeValue *string `json:"scopeValue,omitempty"`
	// Property bag of shared data that was contributed to by any of the individual data providers
	SharedData *map[string]interface{} `json:"sharedData,omitempty"`
}

// Data bag that any data provider can contribute to. This shared dictionary is returned in the data provider result.
type DataProviderSharedData struct {
}

// Contract for handling the extension acquisition process
type ExtensionAcquisitionRequest struct {
	// How the item is being assigned
	AssignmentType *AcquisitionAssignmentType `json:"assignmentType,omitempty"`
	// The id of the subscription used for purchase
	BillingId *string `json:"billingId,omitempty"`
	// The marketplace id (publisherName.extensionName) for the item
	ItemId *string `json:"itemId,omitempty"`
	// The type of operation, such as install, request, purchase
	OperationType *AcquisitionOperationType `json:"operationType,omitempty"`
	// Additional properties which can be added to the request.
	Properties interface{} `json:"properties,omitempty"`
	// How many licenses should be purchased
	Quantity *int `json:"quantity,omitempty"`
}

// Audit log for an extension
type ExtensionAuditLog struct {
	// Collection of audit log entries
	Entries *[]ExtensionAuditLogEntry `json:"entries,omitempty"`
	// Extension that the change was made for
	ExtensionName *string `json:"extensionName,omitempty"`
	// Publisher that the extension is part of
	PublisherName *string `json:"publisherName,omitempty"`
}

// An audit log entry for an extension
type ExtensionAuditLogEntry struct {
	// Change that was made to extension
	AuditAction *string `json:"auditAction,omitempty"`
	// Date at which the change was made
	AuditDate *azuredevops.Time `json:"auditDate,omitempty"`
	// Extra information about the change
	Comment *string `json:"comment,omitempty"`
	// Represents the user who made the change
	UpdatedBy *webapi.IdentityRef `json:"updatedBy,omitempty"`
}

type ExtensionAuthorization struct {
	Id     *uuid.UUID `json:"id,omitempty"`
	Scopes *[]string  `json:"scopes,omitempty"`
}

// Represents a single collection for extension data documents
type ExtensionDataCollection struct {
	// The name of the collection
	CollectionName *string `json:"collectionName,omitempty"`
	// A list of documents belonging to the collection
	Documents *[]interface{} `json:"documents,omitempty"`
	// The type of the collection's scope, such as Default or User
	ScopeType *string `json:"scopeType,omitempty"`
	// The value of the collection's scope, such as Current or Me
	ScopeValue *string `json:"scopeValue,omitempty"`
}

// Represents a query to receive a set of extension data collections
type ExtensionDataCollectionQuery struct {
	// A list of collections to query
	Collections *[]ExtensionDataCollection `json:"collections,omitempty"`
}

type ExtensionEvent struct {
	// The extension which has been updated
	Extension *gallery.PublishedExtension `json:"extension,omitempty"`
	// The current version of the extension that was updated
	ExtensionVersion *string `json:"extensionVersion,omitempty"`
	// Name of the collection for which the extension was requested
	Host *ExtensionHost `json:"host,omitempty"`
	// Gallery host url
	Links *ExtensionEventUrls `json:"links,omitempty"`
	// Represents the user who initiated the update
	ModifiedBy *webapi.IdentityRef `json:"modifiedBy,omitempty"`
	// The type of update that was made
	UpdateType *ExtensionUpdateType `json:"updateType,omitempty"`
}

// Base class for an event callback for an extension
type ExtensionEventCallback struct {
	// The uri of the endpoint that is hit when an event occurs
	Uri *string `json:"uri,omitempty"`
}

// Collection of event callbacks - endpoints called when particular extension events occur.
type ExtensionEventCallbackCollection struct {
	// Optional.  Defines an endpoint that gets called via a POST request to notify that an extension disable has occurred.
	PostDisable *ExtensionEventCallback `json:"postDisable,omitempty"`
	// Optional.  Defines an endpoint that gets called via a POST request to notify that an extension enable has occurred.
	PostEnable *ExtensionEventCallback `json:"postEnable,omitempty"`
	// Optional.  Defines an endpoint that gets called via a POST request to notify that an extension install has completed.
	PostInstall *ExtensionEventCallback `json:"postInstall,omitempty"`
	// Optional.  Defines an endpoint that gets called via a POST request to notify that an extension uninstall has occurred.
	PostUninstall *ExtensionEventCallback `json:"postUninstall,omitempty"`
	// Optional.  Defines an endpoint that gets called via a POST request to notify that an extension update has occurred.
	PostUpdate *ExtensionEventCallback `json:"postUpdate,omitempty"`
	// Optional.  Defines an endpoint that gets called via a POST request to notify that an extension install is about to occur.  Response indicates whether to proceed or abort.
	PreInstall *ExtensionEventCallback `json:"preInstall,omitempty"`
	// For multi-version extensions, defines an endpoint that gets called via an OPTIONS request to determine the particular version of the extension to be used
	VersionCheck *ExtensionEventCallback `json:"versionCheck,omitempty"`
}

type ExtensionEventUrls struct {
	// Url of the extension icon
	ExtensionIcon *string `json:"extensionIcon,omitempty"`
	// Link to view the extension details page
	ExtensionPage *string `json:"extensionPage,omitempty"`
	// Url of the extension management page
	ManageExtensionsPage *string `json:"manageExtensionsPage,omitempty"`
}

// [Flags] Set of flags applied to extensions that are relevant to contribution consumers
type ExtensionFlags string

type extensionFlagsValuesType struct {
	BuiltIn ExtensionFlags
	Trusted ExtensionFlags
}

var ExtensionFlagsValues = extensionFlagsValuesType{
	// A built-in extension is installed for all VSTS accounts by default
	BuiltIn: "builtIn",
	// The extension comes from a fully-trusted publisher
	Trusted: "trusted",
}

type ExtensionHost struct {
	Id   *uuid.UUID `json:"id,omitempty"`
	Name *string    `json:"name,omitempty"`
}

// How an extension should handle including contributions based on licensing
type ExtensionLicensing struct {
	// A list of contributions which deviate from the default licensing behavior
	Overrides *[]LicensingOverride `json:"overrides,omitempty"`
}

// Base class for extension properties which are shared by the extension manifest and the extension model
type ExtensionManifest struct {
	// Uri used as base for other relative uri's defined in extension
	BaseUri *string `json:"baseUri,omitempty"`
	// List of shared constraints defined by this extension
	Constraints *[]ContributionConstraint `json:"constraints,omitempty"`
	// List of contributions made by this extension
	Contributions *[]Contribution `json:"contributions,omitempty"`
	// List of contribution types defined by this extension
	ContributionTypes *[]ContributionType `json:"contributionTypes,omitempty"`
	// List of explicit demands required by this extension
	Demands *[]string `json:"demands,omitempty"`
	// Collection of endpoints that get called when particular extension events occur
	EventCallbacks *ExtensionEventCallbackCollection `json:"eventCallbacks,omitempty"`
	// Secondary location that can be used as base for other relative uri's defined in extension
	FallbackBaseUri *string `json:"fallbackBaseUri,omitempty"`
	// Language Culture Name set by the Gallery
	Language *string `json:"language,omitempty"`
	// How this extension behaves with respect to licensing
	Licensing *ExtensionLicensing `json:"licensing,omitempty"`
	// Version of the extension manifest format/content
	ManifestVersion *float64 `json:"manifestVersion,omitempty"`
	// Default user claims applied to all contributions (except the ones which have been specified restrictedTo explicitly) to control the visibility of a contribution.
	RestrictedTo *[]string `json:"restrictedTo,omitempty"`
	// List of all oauth scopes required by this extension
	Scopes *[]string `json:"scopes,omitempty"`
	// The ServiceInstanceType(Guid) of the VSTS service that must be available to an account in order for the extension to be installed
	ServiceInstanceType *uuid.UUID `json:"serviceInstanceType,omitempty"`
}

// A request for an extension (to be installed or have a license assigned)
type ExtensionRequest struct {
	// Required message supplied if the request is rejected
	RejectMessage *string `json:"rejectMessage,omitempty"`
	// Date at which the request was made
	RequestDate *azuredevops.Time `json:"requestDate,omitempty"`
	// Represents the user who made the request
	RequestedBy *webapi.IdentityRef `json:"requestedBy,omitempty"`
	// Optional message supplied by the requester justifying the request
	RequestMessage *string `json:"requestMessage,omitempty"`
	// Represents the state of the request
	RequestState *ExtensionRequestState `json:"requestState,omitempty"`
	// Date at which the request was resolved
	ResolveDate *azuredevops.Time `json:"resolveDate,omitempty"`
	// Represents the user who resolved the request
	ResolvedBy *webapi.IdentityRef `json:"resolvedBy,omitempty"`
}

type ExtensionRequestEvent struct {
	// The extension which has been requested
	Extension *gallery.PublishedExtension `json:"extension,omitempty"`
	// Information about the host for which this extension is requested
	Host *ExtensionHost `json:"host,omitempty"`
	// Name of the collection for which the extension was requested
	HostName *string `json:"hostName,omitempty"`
	// Gallery host url
	Links *ExtensionRequestUrls `json:"links,omitempty"`
	// The extension request object
	Request *ExtensionRequest `json:"request,omitempty"`
	// The type of update that was made
	UpdateType *ExtensionRequestUpdateType `json:"updateType,omitempty"`
}

type ExtensionRequestsEvent struct {
	// The extension which has been requested
	Extension *gallery.PublishedExtension `json:"extension,omitempty"`
	// Information about the host for which this extension is requested
	Host *ExtensionHost `json:"host,omitempty"`
	// Gallery host url
	Links *ExtensionRequestUrls `json:"links,omitempty"`
	// The extension request object
	Requests *[]ExtensionRequest `json:"requests,omitempty"`
	// The type of update that was made
	UpdateType *ExtensionRequestUpdateType `json:"updateType,omitempty"`
}

// Represents the state of an extension request
type ExtensionRequestState string

type extensionRequestStateValuesType struct {
	Open     ExtensionRequestState
	Accepted ExtensionRequestState
	Rejected ExtensionRequestState
}

var ExtensionRequestStateValues = extensionRequestStateValuesType{
	// The request has been opened, but not yet responded to
	Open: "open",
	// The request was accepted (extension installed or license assigned)
	Accepted: "accepted",
	// The request was rejected (extension not installed or license not assigned)
	Rejected: "rejected",
}

type ExtensionRequestUpdateType string

type extensionRequestUpdateTypeValuesType struct {
	Created  ExtensionRequestUpdateType
	Approved ExtensionRequestUpdateType
	Rejected ExtensionRequestUpdateType
	Deleted  ExtensionRequestUpdateType
}

var ExtensionRequestUpdateTypeValues = extensionRequestUpdateTypeValuesType{
	Created:  "created",
	Approved: "approved",
	Rejected: "rejected",
	Deleted:  "deleted",
}

type ExtensionRequestUrls struct {
	// Url of the extension icon
	ExtensionIcon *string `json:"extensionIcon,omitempty"`
	// Link to view the extension details page
	ExtensionPage *string `json:"extensionPage,omitempty"`
	// Link to view the extension request
	RequestPage *string `json:"requestPage,omitempty"`
}

// The state of an extension
type ExtensionState struct {
	// States of an installed extension
	Flags *ExtensionStateFlags `json:"flags,omitempty"`
	// List of installation issues
	InstallationIssues *[]InstalledExtensionStateIssue `json:"installationIssues,omitempty"`
	// The time at which this installation was last updated
	LastUpdated   *azuredevops.Time `json:"lastUpdated,omitempty"`
	ExtensionName *string           `json:"extensionName,omitempty"`
	// The time at which the version was last checked
	LastVersionCheck *azuredevops.Time `json:"lastVersionCheck,omitempty"`
	PublisherName    *string           `json:"publisherName,omitempty"`
	Version          *string           `json:"version,omitempty"`
}

// [Flags] States of an extension Note:  If you add value to this enum, you need to do 2 other things.  First add the back compat enum in value src\Vssf\Sdk\Server\Contributions\InstalledExtensionMessage.cs.  Second, you can not send the new value on the message bus.  You need to remove it from the message bus event prior to being sent.
type ExtensionStateFlags string

type extensionStateFlagsValuesType struct {
	None                 ExtensionStateFlags
	Disabled             ExtensionStateFlags
	BuiltIn              ExtensionStateFlags
	MultiVersion         ExtensionStateFlags
	UnInstalled          ExtensionStateFlags
	VersionCheckError    ExtensionStateFlags
	Trusted              ExtensionStateFlags
	Error                ExtensionStateFlags
	NeedsReauthorization ExtensionStateFlags
	AutoUpgradeError     ExtensionStateFlags
	Warning              ExtensionStateFlags
}

var ExtensionStateFlagsValues = extensionStateFlagsValuesType{
	// No flags set
	None: "none",
	// Extension is disabled
	Disabled: "disabled",
	// Extension is a built in
	BuiltIn: "builtIn",
	// Extension has multiple versions
	MultiVersion: "multiVersion",
	// Extension is not installed.  This is for builtin extensions only and can not otherwise be set.
	UnInstalled: "unInstalled",
	// Error performing version check
	VersionCheckError: "versionCheckError",
	// Trusted extensions are ones that are given special capabilities. These tend to come from Microsoft and can't be published by the general public.  Note: BuiltIn extensions are always trusted.
	Trusted: "trusted",
	// Extension is currently in an error state
	Error: "error",
	// Extension scopes have changed and the extension requires re-authorization
	NeedsReauthorization: "needsReauthorization",
	// Error performing auto-upgrade. For example, if the new version has demands not supported the extension cannot be auto-upgraded.
	AutoUpgradeError: "autoUpgradeError",
	// Extension is currently in a warning state, that can cause a degraded experience. The degraded experience can be caused for example by some installation issues detected such as implicit demands not supported.
	Warning: "warning",
}

type ExtensionUpdateType string

type extensionUpdateTypeValuesType struct {
	Installed      ExtensionUpdateType
	Uninstalled    ExtensionUpdateType
	Enabled        ExtensionUpdateType
	Disabled       ExtensionUpdateType
	VersionUpdated ExtensionUpdateType
	ActionRequired ExtensionUpdateType
	ActionResolved ExtensionUpdateType
}

var ExtensionUpdateTypeValues = extensionUpdateTypeValuesType{
	Installed:      "installed",
	Uninstalled:    "uninstalled",
	Enabled:        "enabled",
	Disabled:       "disabled",
	VersionUpdated: "versionUpdated",
	ActionRequired: "actionRequired",
	ActionResolved: "actionResolved",
}

type ExtensionUrls struct {
	// Url of the extension icon
	ExtensionIcon *string `json:"extensionIcon,omitempty"`
	// Link to view the extension details page
	ExtensionPage *string `json:"extensionPage,omitempty"`
}

// Represents a VSTS extension along with its installation state
type InstalledExtension struct {
	// Uri used as base for other relative uri's defined in extension
	BaseUri *string `json:"baseUri,omitempty"`
	// List of shared constraints defined by this extension
	Constraints *[]ContributionConstraint `json:"constraints,omitempty"`
	// List of contributions made by this extension
	Contributions *[]Contribution `json:"contributions,omitempty"`
	// List of contribution types defined by this extension
	ContributionTypes *[]ContributionType `json:"contributionTypes,omitempty"`
	// List of explicit demands required by this extension
	Demands *[]string `json:"demands,omitempty"`
	// Collection of endpoints that get called when particular extension events occur
	EventCallbacks *ExtensionEventCallbackCollection `json:"eventCallbacks,omitempty"`
	// Secondary location that can be used as base for other relative uri's defined in extension
	FallbackBaseUri *string `json:"fallbackBaseUri,omitempty"`
	// Language Culture Name set by the Gallery
	Language *string `json:"language,omitempty"`
	// How this extension behaves with respect to licensing
	Licensing *ExtensionLicensing `json:"licensing,omitempty"`
	// Version of the extension manifest format/content
	ManifestVersion *float64 `json:"manifestVersion,omitempty"`
	// Default user claims applied to all contributions (except the ones which have been specified restrictedTo explicitly) to control the visibility of a contribution.
	RestrictedTo *[]string `json:"restrictedTo,omitempty"`
	// List of all oauth scopes required by this extension
	Scopes *[]string `json:"scopes,omitempty"`
	// The ServiceInstanceType(Guid) of the VSTS service that must be available to an account in order for the extension to be installed
	ServiceInstanceType *uuid.UUID `json:"serviceInstanceType,omitempty"`
	// The friendly extension id for this extension - unique for a given publisher.
	ExtensionId *string `json:"extensionId,omitempty"`
	// The display name of the extension.
	ExtensionName *string `json:"extensionName,omitempty"`
	// This is the set of files available from the extension.
	Files *[]gallery.ExtensionFile `json:"files,omitempty"`
	// Extension flags relevant to contribution consumers
	Flags *ExtensionFlags `json:"flags,omitempty"`
	// Information about this particular installation of the extension
	InstallState *InstalledExtensionState `json:"installState,omitempty"`
	// This represents the date/time the extensions was last updated in the gallery. This doesnt mean this version was updated the value represents changes to any and all versions of the extension.
	LastPublished *azuredevops.Time `json:"lastPublished,omitempty"`
	// Unique id of the publisher of this extension
	PublisherId *string `json:"publisherId,omitempty"`
	// The display name of the publisher
	PublisherName *string `json:"publisherName,omitempty"`
	// Unique id for this extension (the same id is used for all versions of a single extension)
	RegistrationId *uuid.UUID `json:"registrationId,omitempty"`
	// Version of this extension
	Version *string `json:"version,omitempty"`
}

type InstalledExtensionQuery struct {
	AssetTypes *[]string                      `json:"assetTypes,omitempty"`
	Monikers   *[]gallery.ExtensionIdentifier `json:"monikers,omitempty"`
}

// The state of an installed extension
type InstalledExtensionState struct {
	// States of an installed extension
	Flags *ExtensionStateFlags `json:"flags,omitempty"`
	// List of installation issues
	InstallationIssues *[]InstalledExtensionStateIssue `json:"installationIssues,omitempty"`
	// The time at which this installation was last updated
	LastUpdated *azuredevops.Time `json:"lastUpdated,omitempty"`
}

// Represents an installation issue
type InstalledExtensionStateIssue struct {
	// The error message
	Message *string `json:"message,omitempty"`
	// Source of the installation issue, for example  "Demands"
	Source *string `json:"source,omitempty"`
	// Installation issue type (Warning, Error)
	Type *InstalledExtensionStateIssueType `json:"type,omitempty"`
}

// Installation issue type (Warning, Error)
type InstalledExtensionStateIssueType string

type installedExtensionStateIssueTypeValuesType struct {
	Warning InstalledExtensionStateIssueType
	Error   InstalledExtensionStateIssueType
}

var InstalledExtensionStateIssueTypeValues = installedExtensionStateIssueTypeValuesType{
	// Represents an installation warning, for example an implicit demand not supported
	Warning: "warning",
	// Represents an installation error, for example an explicit demand not supported
	Error: "error",
}

// Maps a contribution to a licensing behavior
type LicensingOverride struct {
	// How the inclusion of this contribution should change based on licensing
	Behavior *ContributionLicensingBehaviorType `json:"behavior,omitempty"`
	// Fully qualified contribution id which we want to define licensing behavior for
	Id *string `json:"id,omitempty"`
}

// A request for an extension (to be installed or have a license assigned)
type RequestedExtension struct {
	// The unique name of the extension
	ExtensionName *string `json:"extensionName,omitempty"`
	// A list of each request for the extension
	ExtensionRequests *[]ExtensionRequest `json:"extensionRequests,omitempty"`
	// DisplayName of the publisher that owns the extension being published.
	PublisherDisplayName *string `json:"publisherDisplayName,omitempty"`
	// Represents the Publisher of the requested extension
	PublisherName *string `json:"publisherName,omitempty"`
	// The total number of requests for an extension
	RequestCount *int `json:"requestCount,omitempty"`
}

// Entry for a specific data provider's resulting data
type ResolvedDataProvider struct {
	// The total time the data provider took to resolve its data (in milliseconds)
	Duration *float32 `json:"duration,omitempty"`
	Error    *string  `json:"error,omitempty"`
	Id       *string  `json:"id,omitempty"`
}

type Scope struct {
	Description *string `json:"description,omitempty"`
	Title       *string `json:"title,omitempty"`
	Value       *string `json:"value,omitempty"`
}

// Information about the extension
type SupportedExtension struct {
	// Unique Identifier for this extension
	Extension *string `json:"extension,omitempty"`
	// Unique Identifier for this publisher
	Publisher *string `json:"publisher,omitempty"`
	// Supported version for this extension
	Version *string `json:"version,omitempty"`
}

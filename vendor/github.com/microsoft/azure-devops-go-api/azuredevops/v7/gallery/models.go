// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package gallery

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
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
	// The target that this options refer to
	Target *string `json:"target,omitempty"`
}

type Answers struct {
	// Gets or sets the vs marketplace extension name
	VsMarketplaceExtensionName *string `json:"vsMarketplaceExtensionName,omitempty"`
	// Gets or sets the vs marketplace publisher name
	VsMarketplacePublisherName *string `json:"vsMarketplacePublisherName,omitempty"`
}

type AssetDetails struct {
	// Gets or sets the Answers, which contains vs marketplace extension name and publisher name
	Answers *Answers `json:"answers,omitempty"`
	// Gets or sets the VS publisher Id
	PublisherNaturalIdentifier *string `json:"publisherNaturalIdentifier,omitempty"`
}

type AzurePublisher struct {
	AzurePublisherId *string `json:"azurePublisherId,omitempty"`
	PublisherName    *string `json:"publisherName,omitempty"`
}

type AzureRestApiRequestModel struct {
	// Gets or sets the Asset details
	AssetDetails *AssetDetails `json:"assetDetails,omitempty"`
	// Gets or sets the asset id
	AssetId *string `json:"assetId,omitempty"`
	// Gets or sets the asset version
	AssetVersion *uint64 `json:"assetVersion,omitempty"`
	// Gets or sets the customer support email
	CustomerSupportEmail *string `json:"customerSupportEmail,omitempty"`
	// Gets or sets the integration contact email
	IntegrationContactEmail *string `json:"integrationContactEmail,omitempty"`
	// Gets or sets the asset version
	Operation *string `json:"operation,omitempty"`
	// Gets or sets the plan identifier if any.
	PlanId *string `json:"planId,omitempty"`
	// Gets or sets the publisher id
	PublisherId *string `json:"publisherId,omitempty"`
	// Gets or sets the resource type
	Type *string `json:"type,omitempty"`
}

type AzureRestApiResponseModel struct {
	// Gets or sets the Asset details
	AssetDetails *AssetDetails `json:"assetDetails,omitempty"`
	// Gets or sets the asset id
	AssetId *string `json:"assetId,omitempty"`
	// Gets or sets the asset version
	AssetVersion *uint64 `json:"assetVersion,omitempty"`
	// Gets or sets the customer support email
	CustomerSupportEmail *string `json:"customerSupportEmail,omitempty"`
	// Gets or sets the integration contact email
	IntegrationContactEmail *string `json:"integrationContactEmail,omitempty"`
	// Gets or sets the asset version
	Operation *string `json:"operation,omitempty"`
	// Gets or sets the plan identifier if any.
	PlanId *string `json:"planId,omitempty"`
	// Gets or sets the publisher id
	PublisherId *string `json:"publisherId,omitempty"`
	// Gets or sets the resource type
	Type *string `json:"type,omitempty"`
	// Gets or sets the Asset operation status
	OperationStatus *RestApiResponseStatusModel `json:"operationStatus,omitempty"`
}

// This is the set of categories in response to the get category query
type CategoriesResult struct {
	Categories *[]ExtensionCategory `json:"categories,omitempty"`
}

// Definition of one title of a category
type CategoryLanguageTitle struct {
	// The language for which the title is applicable
	Lang *string `json:"lang,omitempty"`
	// The language culture id of the lang parameter
	Lcid *int `json:"lcid,omitempty"`
	// Actual title to be shown on the UI
	Title *string `json:"title,omitempty"`
}

// The structure of a Concern Rather than defining a separate data structure having same fields as QnAItem, we are inheriting from the QnAItem.
type Concern struct {
	// Time when the review was first created
	CreatedDate *azuredevops.Time `json:"createdDate,omitempty"`
	// Unique identifier of a QnA item
	Id *uint64 `json:"id,omitempty"`
	// Get status of item
	Status *QnAItemStatus `json:"status,omitempty"`
	// Text description of the QnA item
	Text *string `json:"text,omitempty"`
	// Time when the review was edited/updated
	UpdatedDate *azuredevops.Time `json:"updatedDate,omitempty"`
	// User details for the item.
	User *UserIdentityRef `json:"user,omitempty"`
	// Category of the concern
	Category *ConcernCategory `json:"category,omitempty"`
}

type ConcernCategory string

type concernCategoryValuesType struct {
	General ConcernCategory
	Abusive ConcernCategory
	Spam    ConcernCategory
}

var ConcernCategoryValues = concernCategoryValuesType{
	General: "general",
	Abusive: "abusive",
	Spam:    "spam",
}

// Stores Last Contact Date
type CustomerLastContact struct {
	// account for which customer was last contacted
	Account *string `json:"account,omitempty"`
	// Date on which the customer was last contacted
	LastContactDate *azuredevops.Time `json:"lastContactDate,omitempty"`
}

// An entity representing the data required to create a Customer Support Request.
type CustomerSupportRequest struct {
	// Display name of extension in concern
	DisplayName *string `json:"displayName,omitempty"`
	// Email of user making the support request
	EmailId *string `json:"emailId,omitempty"`
	// Extension name
	ExtensionName *string `json:"extensionName,omitempty"`
	// Link to the extension details page
	ExtensionURL *string `json:"extensionURL,omitempty"`
	// User-provided support request message.
	Message *string `json:"message,omitempty"`
	// Publisher name
	PublisherName *string `json:"publisherName,omitempty"`
	// Reason for support request
	Reason         *string `json:"reason,omitempty"`
	ReCaptchaToken *string `json:"reCaptchaToken,omitempty"`
	// VSID of the user making the support request
	ReporterVSID *string `json:"reporterVSID,omitempty"`
	// Review under concern
	Review *Review `json:"review,omitempty"`
	// The UI source through which the request was made
	SourceLink *string `json:"sourceLink,omitempty"`
}

type DraftPatchOperation string

type draftPatchOperationValuesType struct {
	Publish DraftPatchOperation
	Cancel  DraftPatchOperation
}

var DraftPatchOperationValues = draftPatchOperationValuesType{
	Publish: "publish",
	Cancel:  "cancel",
}

type DraftStateType string

type draftStateTypeValuesType struct {
	Unpublished DraftStateType
	Published   DraftStateType
	Cancelled   DraftStateType
	Error       DraftStateType
}

var DraftStateTypeValues = draftStateTypeValuesType{
	Unpublished: "unpublished",
	Published:   "published",
	Cancelled:   "cancelled",
	Error:       "error",
}

type EventCounts struct {
	// Average rating on the day for extension
	AverageRating *float32 `json:"averageRating,omitempty"`
	// Number of times the extension was bought in hosted scenario (applies only to VSTS extensions)
	BuyCount *int `json:"buyCount,omitempty"`
	// Number of times the extension was bought in connected scenario (applies only to VSTS extensions)
	ConnectedBuyCount *int `json:"connectedBuyCount,omitempty"`
	// Number of times the extension was installed in connected scenario (applies only to VSTS extensions)
	ConnectedInstallCount *int `json:"connectedInstallCount,omitempty"`
	// Number of times the extension was installed
	InstallCount *uint64 `json:"installCount,omitempty"`
	// Number of times the extension was installed as a trial (applies only to VSTS extensions)
	TryCount *int `json:"tryCount,omitempty"`
	// Number of times the extension was uninstalled (applies only to VSTS extensions)
	UninstallCount *int `json:"uninstallCount,omitempty"`
	// Number of times the extension was downloaded (applies to VSTS extensions and VSCode marketplace click installs)
	WebDownloadCount *uint64 `json:"webDownloadCount,omitempty"`
	// Number of detail page views
	WebPageViews *uint64 `json:"webPageViews,omitempty"`
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
	// A list of target guids where the item should be acquired (installed, requested, etc.), such as account id
	Targets *[]string `json:"targets,omitempty"`
}

type ExtensionBadge struct {
	Description *string `json:"description,omitempty"`
	ImgUri      *string `json:"imgUri,omitempty"`
	Link        *string `json:"link,omitempty"`
}

type ExtensionCategory struct {
	// The name of the products with which this category is associated to.
	AssociatedProducts *[]string `json:"associatedProducts,omitempty"`
	CategoryId         *int      `json:"categoryId,omitempty"`
	// This is the internal name for a category
	CategoryName *string `json:"categoryName,omitempty"`
	// This parameter is obsolete. Refer to LanguageTitles for language specific titles
	Language *string `json:"language,omitempty"`
	// The list of all the titles of this category in various languages
	LanguageTitles *[]CategoryLanguageTitle `json:"languageTitles,omitempty"`
	// This is the internal name of the parent if this is associated with a parent
	ParentCategoryName *string `json:"parentCategoryName,omitempty"`
}

type ExtensionDailyStat struct {
	// Stores the event counts
	Counts *EventCounts `json:"counts,omitempty"`
	// Generic key/value pair to store extended statistics. Used for sending paid extension stats like Upgrade, Downgrade, Cancel trend etc.
	ExtendedStats *map[string]interface{} `json:"extendedStats,omitempty"`
	// Timestamp of this data point
	StatisticDate *azuredevops.Time `json:"statisticDate,omitempty"`
	// Version of the extension
	Version *string `json:"version,omitempty"`
}

type ExtensionDailyStats struct {
	// List of extension statistics data points
	DailyStats *[]ExtensionDailyStat `json:"dailyStats,omitempty"`
	// Id of the extension, this will never be sent back to the client. For internal use only.
	ExtensionId *uuid.UUID `json:"extensionId,omitempty"`
	// Name of the extension
	ExtensionName *string `json:"extensionName,omitempty"`
	// Name of the publisher
	PublisherName *string `json:"publisherName,omitempty"`
	// Count of stats
	StatCount *int `json:"statCount,omitempty"`
}

type ExtensionDeploymentTechnology string

type extensionDeploymentTechnologyValuesType struct {
	Exe          ExtensionDeploymentTechnology
	Msi          ExtensionDeploymentTechnology
	Vsix         ExtensionDeploymentTechnology
	ReferralLink ExtensionDeploymentTechnology
}

var ExtensionDeploymentTechnologyValues = extensionDeploymentTechnologyValuesType{
	Exe:          "exe",
	Msi:          "msi",
	Vsix:         "vsix",
	ReferralLink: "referralLink",
}

type ExtensionDraft struct {
	Assets             *[]ExtensionDraftAsset      `json:"assets,omitempty"`
	CreatedDate        *azuredevops.Time           `json:"createdDate,omitempty"`
	DraftState         *DraftStateType             `json:"draftState,omitempty"`
	ExtensionName      *string                     `json:"extensionName,omitempty"`
	Id                 *uuid.UUID                  `json:"id,omitempty"`
	LastUpdated        *azuredevops.Time           `json:"lastUpdated,omitempty"`
	Payload            *ExtensionPayload           `json:"payload,omitempty"`
	Product            *string                     `json:"product,omitempty"`
	PublisherName      *string                     `json:"publisherName,omitempty"`
	ValidationErrors   *[]azuredevops.KeyValuePair `json:"validationErrors,omitempty"`
	ValidationWarnings *[]azuredevops.KeyValuePair `json:"validationWarnings,omitempty"`
}

type ExtensionDraftAsset struct {
	AssetType *string `json:"assetType,omitempty"`
	Language  *string `json:"language,omitempty"`
	Source    *string `json:"source,omitempty"`
}

type ExtensionDraftPatch struct {
	ExtensionData  *UnpackagedExtensionData `json:"extensionData,omitempty"`
	Operation      *DraftPatchOperation     `json:"operation,omitempty"`
	ReCaptchaToken *string                  `json:"reCaptchaToken,omitempty"`
}

// Stores details of each event
type ExtensionEvent struct {
	// Id which identifies each data point uniquely
	Id         *uint64     `json:"id,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
	// Timestamp of when the event occurred
	StatisticDate *azuredevops.Time `json:"statisticDate,omitempty"`
	// Version of the extension
	Version *string `json:"version,omitempty"`
}

// Container object for all extension events. Stores all install and uninstall events related to an extension. The events container is generic so can store data of any type of event. New event types can be added without altering the contract.
type ExtensionEvents struct {
	// Generic container for events data. The dictionary key denotes the type of event and the list contains properties related to that event
	Events *map[string][]ExtensionEvent `json:"events,omitempty"`
	// Id of the extension, this will never be sent back to the client. This field will mainly be used when EMS calls into Gallery REST API to update install/uninstall events for various extensions in one go.
	ExtensionId *uuid.UUID `json:"extensionId,omitempty"`
	// Name of the extension
	ExtensionName *string `json:"extensionName,omitempty"`
	// Name of the publisher
	PublisherName *string `json:"publisherName,omitempty"`
}

type ExtensionFile struct {
	AssetType *string `json:"assetType,omitempty"`
	Language  *string `json:"language,omitempty"`
	Source    *string `json:"source,omitempty"`
}

// The FilterResult is the set of extensions that matched a particular query filter.
type ExtensionFilterResult struct {
	// This is the set of applications that matched the query filter supplied.
	Extensions *[]PublishedExtension `json:"extensions,omitempty"`
	// The PagingToken is returned from a request when more records exist that match the result than were requested or could be returned. A follow-up query with this paging token can be used to retrieve more results.
	PagingToken *string `json:"pagingToken,omitempty"`
	// This is the additional optional metadata for the given result. E.g. Total count of results which is useful in case of paged results
	ResultMetadata *[]ExtensionFilterResultMetadata `json:"resultMetadata,omitempty"`
}

// ExtensionFilterResultMetadata is one set of metadata for the result e.g. Total count. There can be multiple metadata items for one metadata.
type ExtensionFilterResultMetadata struct {
	// The metadata items for the category
	MetadataItems *[]MetadataItem `json:"metadataItems,omitempty"`
	// Defines the category of metadata items
	MetadataType *string `json:"metadataType,omitempty"`
}

// Represents the component pieces of an extensions fully qualified name, along with the fully qualified name.
type ExtensionIdentifier struct {
	// The ExtensionName component part of the fully qualified ExtensionIdentifier
	ExtensionName *string `json:"extensionName,omitempty"`
	// The PublisherName component part of the fully qualified ExtensionIdentifier
	PublisherName *string `json:"publisherName,omitempty"`
}

// Type of event
type ExtensionLifecycleEventType string

type extensionLifecycleEventTypeValuesType struct {
	Uninstall   ExtensionLifecycleEventType
	Install     ExtensionLifecycleEventType
	Review      ExtensionLifecycleEventType
	Acquisition ExtensionLifecycleEventType
	Sales       ExtensionLifecycleEventType
	Other       ExtensionLifecycleEventType
}

var ExtensionLifecycleEventTypeValues = extensionLifecycleEventTypeValuesType{
	Uninstall:   "uninstall",
	Install:     "install",
	Review:      "review",
	Acquisition: "acquisition",
	Sales:       "sales",
	Other:       "other",
}

// Package that will be used to create or update a published extension
type ExtensionPackage struct {
	// Base 64 encoded extension package
	ExtensionManifest *string `json:"extensionManifest,omitempty"`
}

type ExtensionPayload struct {
	Description         *string                        `json:"description,omitempty"`
	DisplayName         *string                        `json:"displayName,omitempty"`
	FileName            *string                        `json:"fileName,omitempty"`
	InstallationTargets *[]InstallationTarget          `json:"installationTargets,omitempty"`
	IsPreview           *bool                          `json:"isPreview,omitempty"`
	IsSignedByMicrosoft *bool                          `json:"isSignedByMicrosoft,omitempty"`
	IsValid             *bool                          `json:"isValid,omitempty"`
	Metadata            *[]azuredevops.KeyValuePair    `json:"metadata,omitempty"`
	Type                *ExtensionDeploymentTechnology `json:"type,omitempty"`
}

// Policy with a set of permissions on extension operations
type ExtensionPolicy struct {
	// Permissions on 'Install' operation
	Install *ExtensionPolicyFlags `json:"install,omitempty"`
	// Permission on 'Request' operation
	Request *ExtensionPolicyFlags `json:"request,omitempty"`
}

// [Flags] Set of flags that can be associated with a given permission over an extension
type ExtensionPolicyFlags string

type extensionPolicyFlagsValuesType struct {
	None       ExtensionPolicyFlags
	Private    ExtensionPolicyFlags
	Public     ExtensionPolicyFlags
	Preview    ExtensionPolicyFlags
	Released   ExtensionPolicyFlags
	FirstParty ExtensionPolicyFlags
	All        ExtensionPolicyFlags
}

var ExtensionPolicyFlagsValues = extensionPolicyFlagsValuesType{
	// No permission
	None: "none",
	// Permission on private extensions
	Private: "private",
	// Permission on public extensions
	Public: "public",
	// Permission in extensions that are in preview
	Preview: "preview",
	// Permission in released extensions
	Released: "released",
	// Permission in 1st party extensions
	FirstParty: "firstParty",
	// Mask that defines all permissions
	All: "all",
}

// An ExtensionQuery is used to search the gallery for a set of extensions that match one of many filter values.
type ExtensionQuery struct {
	// When retrieving extensions with a query; frequently the caller only needs a small subset of the assets. The caller may specify a list of asset types that should be returned if the extension contains it. All other assets will not be returned.
	AssetTypes *[]string `json:"assetTypes,omitempty"`
	// Each filter is a unique query and will have matching set of extensions returned from the request. Each result will have the same index in the resulting array that the filter had in the incoming query.
	Filters *[]QueryFilter `json:"filters,omitempty"`
	// The Flags are used to determine which set of information the caller would like returned for the matched extensions.
	Flags *ExtensionQueryFlags `json:"flags,omitempty"`
}

// Type of extension filters that are supported in the queries.
type ExtensionQueryFilterType string

type extensionQueryFilterTypeValuesType struct {
	Tag                            ExtensionQueryFilterType
	DisplayName                    ExtensionQueryFilterType
	Private                        ExtensionQueryFilterType
	Id                             ExtensionQueryFilterType
	Category                       ExtensionQueryFilterType
	ContributionType               ExtensionQueryFilterType
	Name                           ExtensionQueryFilterType
	InstallationTarget             ExtensionQueryFilterType
	Featured                       ExtensionQueryFilterType
	SearchText                     ExtensionQueryFilterType
	FeaturedInCategory             ExtensionQueryFilterType
	ExcludeWithFlags               ExtensionQueryFilterType
	IncludeWithFlags               ExtensionQueryFilterType
	Lcid                           ExtensionQueryFilterType
	InstallationTargetVersion      ExtensionQueryFilterType
	InstallationTargetVersionRange ExtensionQueryFilterType
	VsixMetadata                   ExtensionQueryFilterType
	PublisherName                  ExtensionQueryFilterType
	PublisherDisplayName           ExtensionQueryFilterType
	IncludeWithPublisherFlags      ExtensionQueryFilterType
	OrganizationSharedWith         ExtensionQueryFilterType
	ProductArchitecture            ExtensionQueryFilterType
	TargetPlatform                 ExtensionQueryFilterType
	ExtensionName                  ExtensionQueryFilterType
}

var ExtensionQueryFilterTypeValues = extensionQueryFilterTypeValuesType{
	// The values are used as tags. All tags are treated as "OR" conditions with each other. There may be some value put on the number of matched tags from the query.
	Tag: "tag",
	// The Values are an ExtensionName or fragment that is used to match other extension names.
	DisplayName: "displayName",
	// The Filter is one or more tokens that define what scope to return private extensions for.
	Private: "private",
	// Retrieve a set of extensions based on their id's. The values should be the extension id's encoded as strings.
	Id: "id",
	// The category is unlike other filters. It is AND'd with the other filters instead of being a separate query.
	Category: "category",
	// Certain contribution types may be indexed to allow for query by type. User defined types can't be indexed at the moment.
	ContributionType: "contributionType",
	// Retrieve an set extension based on the name based identifier. This differs from the internal id (which is being deprecated).
	Name: "name",
	// The InstallationTarget for an extension defines the target consumer for the extension. This may be something like VS, VSOnline, or VSCode
	InstallationTarget: "installationTarget",
	// Query for featured extensions, no value is allowed when using the query type.
	Featured: "featured",
	// The SearchText provided by the user to search for extensions
	SearchText: "searchText",
	// Query for extensions that are featured in their own category, The filterValue for this is name of category of extensions.
	FeaturedInCategory: "featuredInCategory",
	// When retrieving extensions from a query, exclude the extensions which are having the given flags. The value specified for this filter should be a string representing the integer values of the flags to be excluded. In case of multiple flags to be specified, a logical OR of the interger values should be given as value for this filter This should be at most one filter of this type. This only acts as a restrictive filter after. In case of having a particular flag in both IncludeWithFlags and ExcludeWithFlags, excludeFlags will remove the included extensions giving empty result for that flag.
	ExcludeWithFlags: "excludeWithFlags",
	// When retrieving extensions from a query, include the extensions which are having the given flags. The value specified for this filter should be a string representing the integer values of the flags to be included. In case of multiple flags to be specified, a logical OR of the integer values should be given as value for this filter This should be at most one filter of this type. This only acts as a restrictive filter after. In case of having a particular flag in both IncludeWithFlags and ExcludeWithFlags, excludeFlags will remove the included extensions giving empty result for that flag. In case of multiple flags given in IncludeWithFlags in ORed fashion, extensions having any of the given flags will be included.
	IncludeWithFlags: "includeWithFlags",
	// Filter the extensions based on the LCID values applicable. Any extensions which are not having any LCID values will also be filtered. This is currently only supported for VS extensions.
	Lcid: "lcid",
	// Filter to provide the version of the installation target. This filter will be used along with InstallationTarget filter. The value should be a valid version string. Currently supported only if search text is provided.
	InstallationTargetVersion: "installationTargetVersion",
	// Filter type for specifying a range of installation target version. The filter will be used along with InstallationTarget filter. The value should be a pair of well formed version values separated by hyphen(-). Currently supported only if search text is provided.
	InstallationTargetVersionRange: "installationTargetVersionRange",
	// Filter type for specifying metadata key and value to be used for filtering.
	VsixMetadata: "vsixMetadata",
	// Filter to get extensions published by a publisher having supplied internal name
	PublisherName: "publisherName",
	// Filter to get extensions published by all publishers having supplied display name
	PublisherDisplayName: "publisherDisplayName",
	// When retrieving extensions from a query, include the extensions which have a publisher having the given flags. The value specified for this filter should be a string representing the integer values of the flags to be included. In case of multiple flags to be specified, a logical OR of the integer values should be given as value for this filter There should be at most one filter of this type. This only acts as a restrictive filter after. In case of multiple flags given in IncludeWithFlags in ORed fashion, extensions having any of the given flags will be included.
	IncludeWithPublisherFlags: "includeWithPublisherFlags",
	// Filter to get extensions shared with particular organization
	OrganizationSharedWith: "organizationSharedWith",
	// Filter to get VS IDE extensions by Product Architecture
	ProductArchitecture: "productArchitecture",
	// Filter to get VS Code extensions by target platform.
	TargetPlatform: "targetPlatform",
	// Retrieve an extension based on the extensionName.
	ExtensionName: "extensionName",
}

// [Flags] Set of flags used to determine which set of information is retrieved when reading published extensions
type ExtensionQueryFlags string

type extensionQueryFlagsValuesType struct {
	None                          ExtensionQueryFlags
	IncludeVersions               ExtensionQueryFlags
	IncludeFiles                  ExtensionQueryFlags
	IncludeCategoryAndTags        ExtensionQueryFlags
	IncludeSharedAccounts         ExtensionQueryFlags
	IncludeVersionProperties      ExtensionQueryFlags
	ExcludeNonValidated           ExtensionQueryFlags
	IncludeInstallationTargets    ExtensionQueryFlags
	IncludeAssetUri               ExtensionQueryFlags
	IncludeStatistics             ExtensionQueryFlags
	IncludeLatestVersionOnly      ExtensionQueryFlags
	UseFallbackAssetUri           ExtensionQueryFlags
	IncludeMetadata               ExtensionQueryFlags
	IncludeMinimalPayloadForVsIde ExtensionQueryFlags
	IncludeLcids                  ExtensionQueryFlags
	IncludeSharedOrganizations    ExtensionQueryFlags
	IncludeNameConflictInfo       ExtensionQueryFlags
	AllAttributes                 ExtensionQueryFlags
}

var ExtensionQueryFlagsValues = extensionQueryFlagsValuesType{
	// None is used to retrieve only the basic extension details.
	None: "none",
	// IncludeVersions will return version information for extensions returned
	IncludeVersions: "includeVersions",
	// IncludeFiles will return information about which files were found within the extension that were stored independent of the manifest. When asking for files, versions will be included as well since files are returned as a property of the versions.  These files can be retrieved using the path to the file without requiring the entire manifest be downloaded.
	IncludeFiles: "includeFiles",
	// Include the Categories and Tags that were added to the extension definition.
	IncludeCategoryAndTags: "includeCategoryAndTags",
	// Include the details about which accounts the extension has been shared with if the extension is a private extension.
	IncludeSharedAccounts: "includeSharedAccounts",
	// Include properties associated with versions of the extension
	IncludeVersionProperties: "includeVersionProperties",
	// Excluding non-validated extensions will remove any extension versions that either are in the process of being validated or have failed validation.
	ExcludeNonValidated: "excludeNonValidated",
	// Include the set of installation targets the extension has requested.
	IncludeInstallationTargets: "includeInstallationTargets",
	// Include the base uri for assets of this extension
	IncludeAssetUri: "includeAssetUri",
	// Include the statistics associated with this extension
	IncludeStatistics: "includeStatistics",
	// When retrieving versions from a query, only include the latest version of the extensions that matched. This is useful when the caller doesn't need all the published versions. It will save a significant size in the returned payload.
	IncludeLatestVersionOnly: "includeLatestVersionOnly",
	// This flag switches the asset uri to use GetAssetByName instead of CDN When this is used, values of base asset uri and base asset uri fallback are switched When this is used, source of asset files are pointed to Gallery service always even if CDN is available
	UseFallbackAssetUri: "useFallbackAssetUri",
	// This flag is used to get all the metadata values associated with the extension. This is not applicable to VSTS or VSCode extensions and usage is only internal.
	IncludeMetadata: "includeMetadata",
	// This flag is used to indicate to return very small data for extension required by VS IDE. This flag is only compatible when querying is done by VS IDE
	IncludeMinimalPayloadForVsIde: "includeMinimalPayloadForVsIde",
	// This flag is used to get Lcid values associated with the extension. This is not applicable to VSTS or VSCode extensions and usage is only internal
	IncludeLcids: "includeLcids",
	// Include the details about which organizations the extension has been shared with if the extension is a private extension.
	IncludeSharedOrganizations: "includeSharedOrganizations",
	// Include the details if an extension is in conflict list or not Currently being used for VSCode extensions.
	IncludeNameConflictInfo: "includeNameConflictInfo",
	// AllAttributes is designed to be a mask that defines all sub-elements of the extension should be returned.  NOTE: This is not actually All flags. This is now locked to the set defined since changing this enum would be a breaking change and would change the behavior of anyone using it. Try not to use this value when making calls to the service, instead be explicit about the options required.
	AllAttributes: "allAttributes",
}

// This is the set of extensions that matched a supplied query through the filters given.
type ExtensionQueryResult struct {
	// For each filter supplied in the query, a filter result will be returned in the query result.
	Results *[]ExtensionFilterResult `json:"results,omitempty"`
}

type ExtensionShare struct {
	Id    *string `json:"id,omitempty"`
	IsOrg *bool   `json:"isOrg,omitempty"`
	Name  *string `json:"name,omitempty"`
	Type  *string `json:"type,omitempty"`
}

type ExtensionStatistic struct {
	StatisticName *string  `json:"statisticName,omitempty"`
	Value         *float64 `json:"value,omitempty"`
}

type ExtensionStatisticOperation string

type extensionStatisticOperationValuesType struct {
	None      ExtensionStatisticOperation
	Set       ExtensionStatisticOperation
	Increment ExtensionStatisticOperation
	Decrement ExtensionStatisticOperation
	Delete    ExtensionStatisticOperation
}

var ExtensionStatisticOperationValues = extensionStatisticOperationValuesType{
	None:      "none",
	Set:       "set",
	Increment: "increment",
	Decrement: "decrement",
	Delete:    "delete",
}

type ExtensionStatisticUpdate struct {
	ExtensionName *string                      `json:"extensionName,omitempty"`
	Operation     *ExtensionStatisticOperation `json:"operation,omitempty"`
	PublisherName *string                      `json:"publisherName,omitempty"`
	Statistic     *ExtensionStatistic          `json:"statistic,omitempty"`
}

// Stats aggregation type
type ExtensionStatsAggregateType string

type extensionStatsAggregateTypeValuesType struct {
	Daily ExtensionStatsAggregateType
}

var ExtensionStatsAggregateTypeValues = extensionStatsAggregateTypeValuesType{
	Daily: "daily",
}

type ExtensionVersion struct {
	AssetUri                *string                     `json:"assetUri,omitempty"`
	Badges                  *[]ExtensionBadge           `json:"badges,omitempty"`
	FallbackAssetUri        *string                     `json:"fallbackAssetUri,omitempty"`
	Files                   *[]ExtensionFile            `json:"files,omitempty"`
	Flags                   *ExtensionVersionFlags      `json:"flags,omitempty"`
	LastUpdated             *azuredevops.Time           `json:"lastUpdated,omitempty"`
	Properties              *[]azuredevops.KeyValuePair `json:"properties,omitempty"`
	TargetPlatform          *string                     `json:"targetPlatform,omitempty"`
	ValidationResultMessage *string                     `json:"validationResultMessage,omitempty"`
	Version                 *string                     `json:"version,omitempty"`
	VersionDescription      *string                     `json:"versionDescription,omitempty"`
}

// [Flags] Set of flags that can be associated with a given extension version. These flags apply to a specific version of the extension.
type ExtensionVersionFlags string

type extensionVersionFlagsValuesType struct {
	None      ExtensionVersionFlags
	Validated ExtensionVersionFlags
}

var ExtensionVersionFlagsValues = extensionVersionFlagsValuesType{
	// No flags exist for this version.
	None: "none",
	// The Validated flag for a version means the extension version has passed validation and can be used..
	Validated: "validated",
}

// One condition in a QueryFilter.
type FilterCriteria struct {
	FilterType *int `json:"filterType,omitempty"`
	// The value used in the match based on the filter type.
	Value *string `json:"value,omitempty"`
}

type InstallationTarget struct {
	ExtensionVersion    *string `json:"extensionVersion,omitempty"`
	ProductArchitecture *string `json:"productArchitecture,omitempty"`
	Target              *string `json:"target,omitempty"`
	TargetPlatform      *string `json:"targetPlatform,omitempty"`
	TargetVersion       *string `json:"targetVersion,omitempty"`
}

// MetadataItem is one value of metadata under a given category of metadata
type MetadataItem struct {
	// The count of the metadata item
	Count *int `json:"count,omitempty"`
	// The name of the metadata item
	Name *string `json:"name,omitempty"`
}

// Information needed for sending mail notification
type NotificationsData struct {
	// Notification data needed
	Data *map[string]interface{} `json:"data,omitempty"`
	// List of users who should get the notification
	Identities *map[string]interface{} `json:"identities,omitempty"`
	// Type of Mail Notification.Can be Qna , review or CustomerContact
	Type *NotificationTemplateType `json:"type,omitempty"`
}

// Type of event
type NotificationTemplateType string

type notificationTemplateTypeValuesType struct {
	ReviewNotification                NotificationTemplateType
	QnaNotification                   NotificationTemplateType
	CustomerContactNotification       NotificationTemplateType
	PublisherMemberUpdateNotification NotificationTemplateType
}

var NotificationTemplateTypeValues = notificationTemplateTypeValuesType{
	// Template type for Review Notification.
	ReviewNotification: "reviewNotification",
	// Template type for Qna Notification.
	QnaNotification: "qnaNotification",
	// Template type for Customer Contact Notification.
	CustomerContactNotification: "customerContactNotification",
	// Template type for Publisher Member Notification.
	PublisherMemberUpdateNotification: "publisherMemberUpdateNotification",
}

// PagingDirection is used to define which set direction to move the returned result set based on a previous query.
type PagingDirection string

type pagingDirectionValuesType struct {
	Backward PagingDirection
	Forward  PagingDirection
}

var PagingDirectionValues = pagingDirectionValuesType{
	// Backward will return results from earlier in the resultset.
	Backward: "backward",
	// Forward will return results from later in the resultset.
	Forward: "forward",
}

// This is the set of categories in response to the get category query
type ProductCategoriesResult struct {
	Categories *[]ProductCategory `json:"categories,omitempty"`
}

// This is the interface object to be used by Root Categories and Category Tree APIs for Visual Studio Ide.
type ProductCategory struct {
	// Indicator whether this is a leaf or there are children under this category
	HasChildren *bool              `json:"hasChildren,omitempty"`
	Children    *[]ProductCategory `json:"children,omitempty"`
	// Individual Guid of the Category
	Id *uuid.UUID `json:"id,omitempty"`
	// Category Title in the requested language
	Title *string `json:"title,omitempty"`
}

type PublishedExtension struct {
	Categories          *[]string                      `json:"categories,omitempty"`
	DeploymentType      *ExtensionDeploymentTechnology `json:"deploymentType,omitempty"`
	DisplayName         *string                        `json:"displayName,omitempty"`
	ExtensionId         *uuid.UUID                     `json:"extensionId,omitempty"`
	ExtensionName       *string                        `json:"extensionName,omitempty"`
	Flags               *PublishedExtensionFlags       `json:"flags,omitempty"`
	InstallationTargets *[]InstallationTarget          `json:"installationTargets,omitempty"`
	LastUpdated         *azuredevops.Time              `json:"lastUpdated,omitempty"`
	LongDescription     *string                        `json:"longDescription,omitempty"`
	// Check if Extension is in conflict list or not. Taking as String and not as boolean because we don't want end customer to see this flag and by making it Boolean it is coming as false for all the cases.
	PresentInConflictList *string `json:"presentInConflictList,omitempty"`
	// Date on which the extension was first uploaded.
	PublishedDate *azuredevops.Time `json:"publishedDate,omitempty"`
	Publisher     *PublisherFacts   `json:"publisher,omitempty"`
	// Date on which the extension first went public.
	ReleaseDate      *azuredevops.Time     `json:"releaseDate,omitempty"`
	SharedWith       *[]ExtensionShare     `json:"sharedWith,omitempty"`
	ShortDescription *string               `json:"shortDescription,omitempty"`
	Statistics       *[]ExtensionStatistic `json:"statistics,omitempty"`
	Tags             *[]string             `json:"tags,omitempty"`
	Versions         *[]ExtensionVersion   `json:"versions,omitempty"`
}

// [Flags] Set of flags that can be associated with a given extension. These flags apply to all versions of the extension and not to a specific version.
type PublishedExtensionFlags string

type publishedExtensionFlagsValuesType struct {
	None         PublishedExtensionFlags
	Disabled     PublishedExtensionFlags
	BuiltIn      PublishedExtensionFlags
	Validated    PublishedExtensionFlags
	Trusted      PublishedExtensionFlags
	Paid         PublishedExtensionFlags
	Public       PublishedExtensionFlags
	MultiVersion PublishedExtensionFlags
	System       PublishedExtensionFlags
	Preview      PublishedExtensionFlags
	Unpublished  PublishedExtensionFlags
	Trial        PublishedExtensionFlags
	Locked       PublishedExtensionFlags
	Hidden       PublishedExtensionFlags
}

var PublishedExtensionFlagsValues = publishedExtensionFlagsValuesType{
	// No flags exist for this extension.
	None: "none",
	// The Disabled flag for an extension means the extension can't be changed and won't be used by consumers. The disabled flag is managed by the service and can't be supplied by the Extension Developers.
	Disabled: "disabled",
	// BuiltIn Extension are available to all Tenants. An explicit registration is not required. This attribute is reserved and can't be supplied by Extension Developers.  BuiltIn extensions are by definition Public. There is no need to set the public flag for extensions marked BuiltIn.
	BuiltIn: "builtIn",
	// This extension has been validated by the service. The extension meets the requirements specified. This attribute is reserved and can't be supplied by the Extension Developers. Validation is a process that ensures that all contributions are well formed. They meet the requirements defined by the contribution type they are extending. Note this attribute will be updated asynchronously as the extension is validated by the developer of the contribution type. There will be restricted access to the extension while this process is performed.
	Validated: "validated",
	// Trusted extensions are ones that are given special capabilities. These tend to come from Microsoft and can't be published by the general public.  Note: BuiltIn extensions are always trusted.
	Trusted: "trusted",
	// The Paid flag indicates that the commerce can be enabled for this extension. Publisher needs to setup Offer/Pricing plan in Azure. If Paid flag is set and a corresponding Offer is not available, the extension will automatically be marked as Preview. If the publisher intends to make the extension Paid in the future, it is mandatory to set the Preview flag. This is currently available only for VSTS extensions only.
	Paid: "paid",
	// This extension registration is public, making its visibility open to the public. This means all tenants have the ability to install this extension. Without this flag the extension will be private and will need to be shared with the tenants that can install it.
	Public: "public",
	// This extension has multiple versions active at one time and version discovery should be done using the defined "Version Discovery" protocol to determine the version available to a specific user or tenant.  @TODO: Link to Version Discovery Protocol.
	MultiVersion: "multiVersion",
	// The system flag is reserved, and cant be used by publishers.
	System: "system",
	// The Preview flag indicates that the extension is still under preview (not yet of "release" quality). These extensions may be decorated differently in the gallery and may have different policies applied to them.
	Preview: "preview",
	// The Unpublished flag indicates that the extension can't be installed/downloaded. Users who have installed such an extension can continue to use the extension.
	Unpublished: "unpublished",
	// The Trial flag indicates that the extension is in Trial version. The flag is right now being used only with respect to Visual Studio extensions.
	Trial: "trial",
	// The Locked flag indicates that extension has been locked from Marketplace. Further updates/acquisitions are not allowed on the extension until this is present. This should be used along with making the extension private/unpublished.
	Locked: "locked",
	// This flag is set for extensions we want to hide from Marketplace home and search pages. This will be used to override the exposure of builtIn flags.
	Hidden: "hidden",
}

type Publisher struct {
	DisplayName        *string               `json:"displayName,omitempty"`
	EmailAddress       *[]string             `json:"emailAddress,omitempty"`
	Extensions         *[]PublishedExtension `json:"extensions,omitempty"`
	Flags              *PublisherFlags       `json:"flags,omitempty"`
	LastUpdated        *azuredevops.Time     `json:"lastUpdated,omitempty"`
	LongDescription    *string               `json:"longDescription,omitempty"`
	PublisherId        *uuid.UUID            `json:"publisherId,omitempty"`
	PublisherName      *string               `json:"publisherName,omitempty"`
	ShortDescription   *string               `json:"shortDescription,omitempty"`
	State              *PublisherState       `json:"state,omitempty"`
	Links              interface{}           `json:"_links,omitempty"`
	Domain             *string               `json:"domain,omitempty"`
	IsDnsTokenVerified *bool                 `json:"isDnsTokenVerified,omitempty"`
	IsDomainVerified   *bool                 `json:"isDomainVerified,omitempty"`
	ReCaptchaToken     *string               `json:"reCaptchaToken,omitempty"`
}

// Keeping base class separate since publisher DB model class and publisher contract class share these common properties
type PublisherBase struct {
	DisplayName      *string               `json:"displayName,omitempty"`
	EmailAddress     *[]string             `json:"emailAddress,omitempty"`
	Extensions       *[]PublishedExtension `json:"extensions,omitempty"`
	Flags            *PublisherFlags       `json:"flags,omitempty"`
	LastUpdated      *azuredevops.Time     `json:"lastUpdated,omitempty"`
	LongDescription  *string               `json:"longDescription,omitempty"`
	PublisherId      *uuid.UUID            `json:"publisherId,omitempty"`
	PublisherName    *string               `json:"publisherName,omitempty"`
	ShortDescription *string               `json:"shortDescription,omitempty"`
	State            *PublisherState       `json:"state,omitempty"`
}

// High-level information about the publisher, like id's and names
type PublisherFacts struct {
	DisplayName      *string         `json:"displayName,omitempty"`
	Domain           *string         `json:"domain,omitempty"`
	Flags            *PublisherFlags `json:"flags,omitempty"`
	IsDomainVerified *bool           `json:"isDomainVerified,omitempty"`
	PublisherId      *uuid.UUID      `json:"publisherId,omitempty"`
	PublisherName    *string         `json:"publisherName,omitempty"`
}

// The FilterResult is the set of publishers that matched a particular query filter.
type PublisherFilterResult struct {
	// This is the set of applications that matched the query filter supplied.
	Publishers *[]Publisher `json:"publishers,omitempty"`
}

// [Flags]
type PublisherFlags string

type publisherFlagsValuesType struct {
	UnChanged    PublisherFlags
	None         PublisherFlags
	Disabled     PublisherFlags
	Verified     PublisherFlags
	Certified    PublisherFlags
	ServiceFlags PublisherFlags
}

var PublisherFlagsValues = publisherFlagsValuesType{
	// This should never be returned, it is used to represent a publisher who's flags haven't changed during update calls.
	UnChanged: "unChanged",
	// No flags exist for this publisher.
	None: "none",
	// The Disabled flag for a publisher means the publisher can't be changed and won't be used by consumers, this extends to extensions owned by the publisher as well. The disabled flag is managed by the service and can't be supplied by the Extension Developers.
	Disabled: "disabled",
	// A verified publisher is one that Microsoft has done some review of and ensured the publisher meets a set of requirements. The requirements to become a verified publisher are not listed here.  They can be found in public documentation (TBD).
	Verified: "verified",
	// A Certified publisher is one that is Microsoft verified and in addition meets a set of requirements for its published extensions. The requirements to become a certified publisher are not listed here.  They can be found in public documentation (TBD).
	Certified: "certified",
	// This is the set of flags that can't be supplied by the developer and is managed by the service itself.
	ServiceFlags: "serviceFlags",
}

// [Flags]
type PublisherPermissions string

type publisherPermissionsValuesType struct {
	Read              PublisherPermissions
	UpdateExtension   PublisherPermissions
	CreatePublisher   PublisherPermissions
	PublishExtension  PublisherPermissions
	Admin             PublisherPermissions
	TrustedPartner    PublisherPermissions
	PrivateRead       PublisherPermissions
	DeleteExtension   PublisherPermissions
	EditSettings      PublisherPermissions
	ViewPermissions   PublisherPermissions
	ManagePermissions PublisherPermissions
	DeletePublisher   PublisherPermissions
}

var PublisherPermissionsValues = publisherPermissionsValuesType{
	// This gives the bearer the rights to read Publishers and Extensions.
	Read: "read",
	// This gives the bearer the rights to update, delete, and share Extensions (but not the ability to create them).
	UpdateExtension: "updateExtension",
	// This gives the bearer the rights to create new Publishers at the root of the namespace.
	CreatePublisher: "createPublisher",
	// This gives the bearer the rights to create new Extensions within a publisher.
	PublishExtension: "publishExtension",
	// Admin gives the bearer the rights to manage restricted attributes of Publishers and Extensions.
	Admin: "admin",
	// TrustedPartner gives the bearer the rights to publish a extensions with restricted capabilities.
	TrustedPartner: "trustedPartner",
	// PrivateRead is another form of read designed to allow higher privilege accessors the ability to read private extensions.
	PrivateRead: "privateRead",
	// This gives the bearer the rights to delete any extension.
	DeleteExtension: "deleteExtension",
	// This gives the bearer the rights edit the publisher settings.
	EditSettings: "editSettings",
	// This gives the bearer the rights to see all permissions on the publisher.
	ViewPermissions: "viewPermissions",
	// This gives the bearer the rights to assign permissions on the publisher.
	ManagePermissions: "managePermissions",
	// This gives the bearer the rights to delete the publisher.
	DeletePublisher: "deletePublisher",
}

// An PublisherQuery is used to search the gallery for a set of publishers that match one of many filter values.
type PublisherQuery struct {
	// Each filter is a unique query and will have matching set of publishers returned from the request. Each result will have the same index in the resulting array that the filter had in the incoming query.
	Filters *[]QueryFilter `json:"filters,omitempty"`
	// The Flags are used to determine which set of information the caller would like returned for the matched publishers.
	Flags *PublisherQueryFlags `json:"flags,omitempty"`
}

// [Flags] Set of flags used to define the attributes requested when a publisher is returned. Some API's allow the caller to specify the level of detail needed.
type PublisherQueryFlags string

type publisherQueryFlagsValuesType struct {
	None                PublisherQueryFlags
	IncludeExtensions   PublisherQueryFlags
	IncludeEmailAddress PublisherQueryFlags
}

var PublisherQueryFlagsValues = publisherQueryFlagsValuesType{
	// None is used to retrieve only the basic publisher details.
	None: "none",
	// Is used to include a list of basic extension details for all extensions published by the requested publisher.
	IncludeExtensions: "includeExtensions",
	// Is used to include email address of all the users who are marked as owners for the publisher
	IncludeEmailAddress: "includeEmailAddress",
}

// This is the set of publishers that matched a supplied query through the filters given.
type PublisherQueryResult struct {
	// For each filter supplied in the query, a filter result will be returned in the query result.
	Results *[]PublisherFilterResult `json:"results,omitempty"`
}

// Access definition for a RoleAssignment.
type PublisherRoleAccess string

type publisherRoleAccessValuesType struct {
	Assigned  PublisherRoleAccess
	Inherited PublisherRoleAccess
}

var PublisherRoleAccessValues = publisherRoleAccessValuesType{
	// Access has been explicitly set.
	Assigned: "assigned",
	// Access has been inherited from a higher scope.
	Inherited: "inherited",
}

type PublisherRoleAssignment struct {
	// Designates the role as explicitly assigned or inherited.
	Access *PublisherRoleAccess `json:"access,omitempty"`
	// User friendly description of access assignment.
	AccessDisplayName *string `json:"accessDisplayName,omitempty"`
	// The user to whom the role is assigned.
	Identity *webapi.IdentityRef `json:"identity,omitempty"`
	// The role assigned to the user.
	Role *PublisherSecurityRole `json:"role,omitempty"`
}

type PublisherSecurityRole struct {
	// Permissions the role is allowed.
	AllowPermissions *int `json:"allowPermissions,omitempty"`
	// Permissions the role is denied.
	DenyPermissions *int `json:"denyPermissions,omitempty"`
	// Description of user access defined by the role
	Description *string `json:"description,omitempty"`
	// User friendly name of the role.
	DisplayName *string `json:"displayName,omitempty"`
	// Globally unique identifier for the role.
	Identifier *string `json:"identifier,omitempty"`
	// Unique name of the role in the scope.
	Name *string `json:"name,omitempty"`
	// Returns the id of the ParentScope.
	Scope *string `json:"scope,omitempty"`
}

// [Flags]
type PublisherState string

type publisherStateValuesType struct {
	None                  PublisherState
	VerificationPending   PublisherState
	CertificationPending  PublisherState
	CertificationRejected PublisherState
	CertificationRevoked  PublisherState
}

var PublisherStateValues = publisherStateValuesType{
	// No state exists for this publisher.
	None: "none",
	// This state indicates that publisher has applied for Marketplace verification (via UI) and still not been certified. This state would be reset once the publisher is verified.
	VerificationPending: "verificationPending",
	// This state indicates that publisher has applied for Marketplace certification (via UI) and still not been certified. This state would be reset once the publisher is certified.
	CertificationPending: "certificationPending",
	// This state indicates that publisher had applied for Marketplace certification (via UI) but his/her certification got rejected. This state would be reset if and when the publisher is certified.
	CertificationRejected: "certificationRejected",
	// This state indicates that publisher was certified on the Marketplace, but his/her certification got revoked. This state would never be reset, even after publisher gets re-certified. It would indicate that the publisher certification was revoked at least once.
	CertificationRevoked: "certificationRevoked",
}

type PublisherUserRoleAssignmentRef struct {
	// The name of the role assigned.
	RoleName *string `json:"roleName,omitempty"`
	// Identifier of the user given the role assignment.
	UniqueName *string `json:"uniqueName,omitempty"`
	// Unique id of the user given the role assignment.
	UserId *uuid.UUID `json:"userId,omitempty"`
}

// The core structure of a QnA item
type QnAItem struct {
	// Time when the review was first created
	CreatedDate *azuredevops.Time `json:"createdDate,omitempty"`
	// Unique identifier of a QnA item
	Id *uint64 `json:"id,omitempty"`
	// Get status of item
	Status *QnAItemStatus `json:"status,omitempty"`
	// Text description of the QnA item
	Text *string `json:"text,omitempty"`
	// Time when the review was edited/updated
	UpdatedDate *azuredevops.Time `json:"updatedDate,omitempty"`
	// User details for the item.
	User *UserIdentityRef `json:"user,omitempty"`
}

// [Flags] Denotes the status of the QnA Item
type QnAItemStatus string

type qnAItemStatusValuesType struct {
	None             QnAItemStatus
	UserEditable     QnAItemStatus
	PublisherCreated QnAItemStatus
}

var QnAItemStatusValues = qnAItemStatusValuesType{
	None: "none",
	// The UserEditable flag indicates whether the item is editable by the logged in user.
	UserEditable: "userEditable",
	// The PublisherCreated flag indicates whether the item has been created by extension publisher.
	PublisherCreated: "publisherCreated",
}

// A filter used to define a set of extensions to return during a query.
type QueryFilter struct {
	// The filter values define the set of values in this query. They are applied based on the QueryFilterType.
	Criteria *[]FilterCriteria `json:"criteria,omitempty"`
	// The PagingDirection is applied to a paging token if one exists. If not the direction is ignored, and Forward from the start of the resultset is used. Direction should be left out of the request unless a paging token is used to help prevent future issues.
	Direction *PagingDirection `json:"direction,omitempty"`
	// The page number requested by the user. If not provided 1 is assumed by default.
	PageNumber *int `json:"pageNumber,omitempty"`
	// The page size defines the number of results the caller wants for this filter. The count can't exceed the overall query size limits.
	PageSize *int `json:"pageSize,omitempty"`
	// The paging token is a distinct type of filter and the other filter fields are ignored. The paging token represents the continuation of a previously executed query. The information about where in the result and what fields are being filtered are embedded in the token.
	PagingToken *string `json:"pagingToken,omitempty"`
	// Defines the type of sorting to be applied on the results. The page slice is cut of the sorted results only.
	SortBy *int `json:"sortBy,omitempty"`
	// Defines the order of sorting, 1 for Ascending, 2 for Descending, else default ordering based on the SortBy value
	SortOrder *int `json:"sortOrder,omitempty"`
}

// The structure of the question / thread
type Question struct {
	// Time when the review was first created
	CreatedDate *azuredevops.Time `json:"createdDate,omitempty"`
	// Unique identifier of a QnA item
	Id *uint64 `json:"id,omitempty"`
	// Get status of item
	Status *QnAItemStatus `json:"status,omitempty"`
	// Text description of the QnA item
	Text *string `json:"text,omitempty"`
	// Time when the review was edited/updated
	UpdatedDate *azuredevops.Time `json:"updatedDate,omitempty"`
	// User details for the item.
	User           *UserIdentityRef `json:"user,omitempty"`
	ReCaptchaToken *string          `json:"reCaptchaToken,omitempty"`
	// List of answers in for the question / thread
	Responses *[]Response `json:"responses,omitempty"`
}

type QuestionsResult struct {
	// Flag indicating if there are more QnA threads to be shown (for paging)
	HasMoreQuestions *bool `json:"hasMoreQuestions,omitempty"`
	// List of the QnA threads
	Questions *[]Question `json:"questions,omitempty"`
}

type RatingCountPerRating struct {
	// Rating value
	Rating *byte `json:"rating,omitempty"`
	// Count of total ratings
	RatingCount *uint64 `json:"ratingCount,omitempty"`
}

// The structure of a response
type Response struct {
	// Time when the review was first created
	CreatedDate *azuredevops.Time `json:"createdDate,omitempty"`
	// Unique identifier of a QnA item
	Id *uint64 `json:"id,omitempty"`
	// Get status of item
	Status *QnAItemStatus `json:"status,omitempty"`
	// Text description of the QnA item
	Text *string `json:"text,omitempty"`
	// Time when the review was edited/updated
	UpdatedDate *azuredevops.Time `json:"updatedDate,omitempty"`
	// User details for the item.
	User           *UserIdentityRef `json:"user,omitempty"`
	ReCaptchaToken *string          `json:"reCaptchaToken,omitempty"`
}

// The status of a REST Api response status.
type RestApiResponseStatus string

type restApiResponseStatusValuesType struct {
	Completed  RestApiResponseStatus
	Failed     RestApiResponseStatus
	Inprogress RestApiResponseStatus
	Skipped    RestApiResponseStatus
}

var RestApiResponseStatusValues = restApiResponseStatusValuesType{
	// The operation is completed.
	Completed: "completed",
	// The operation is failed.
	Failed: "failed",
	// The operation is in progress.
	Inprogress: "inprogress",
	// The operation is in skipped.
	Skipped: "skipped",
}

// REST Api Response
type RestApiResponseStatusModel struct {
	// Gets or sets the operation details
	OperationDetails interface{} `json:"operationDetails,omitempty"`
	// Gets or sets the operation id
	OperationId *string `json:"operationId,omitempty"`
	// Gets or sets the completed status percentage
	PercentageCompleted *int `json:"percentageCompleted,omitempty"`
	// Gets or sets the status
	Status *RestApiResponseStatus `json:"status,omitempty"`
	// Gets or sets the status message
	StatusMessage *string `json:"statusMessage,omitempty"`
}

type Review struct {
	// Admin Reply, if any, for this review
	AdminReply *ReviewReply `json:"adminReply,omitempty"`
	// Unique identifier of a review item
	Id *uint64 `json:"id,omitempty"`
	// Flag for soft deletion
	IsDeleted *bool `json:"isDeleted,omitempty"`
	IsIgnored *bool `json:"isIgnored,omitempty"`
	// Version of the product for which review was submitted
	ProductVersion *string `json:"productVersion,omitempty"`
	// Rating provided by the user
	Rating         *byte   `json:"rating,omitempty"`
	ReCaptchaToken *string `json:"reCaptchaToken,omitempty"`
	// Reply, if any, for this review
	Reply *ReviewReply `json:"reply,omitempty"`
	// Text description of the review
	Text *string `json:"text,omitempty"`
	// Title of the review
	Title *string `json:"title,omitempty"`
	// Time when the review was edited/updated
	UpdatedDate *azuredevops.Time `json:"updatedDate,omitempty"`
	// Name of the user
	UserDisplayName *string `json:"userDisplayName,omitempty"`
	// Id of the user who submitted the review
	UserId *uuid.UUID `json:"userId,omitempty"`
}

// Type of operation
type ReviewEventOperation string

type reviewEventOperationValuesType struct {
	Create ReviewEventOperation
	Update ReviewEventOperation
	Delete ReviewEventOperation
}

var ReviewEventOperationValues = reviewEventOperationValuesType{
	Create: "create",
	Update: "update",
	Delete: "delete",
}

// Properties associated with Review event
type ReviewEventProperties struct {
	// Operation performed on Event - Create\Update
	EventOperation *ReviewEventOperation `json:"eventOperation,omitempty"`
	// Flag to see if reply is admin reply
	IsAdminReply *bool `json:"isAdminReply,omitempty"`
	// Flag to record if the review is ignored
	IsIgnored *bool `json:"isIgnored,omitempty"`
	// Rating at the time of event
	Rating *int `json:"rating,omitempty"`
	// Reply update date
	ReplyDate *azuredevops.Time `json:"replyDate,omitempty"`
	// Publisher reply text or admin reply text
	ReplyText *string `json:"replyText,omitempty"`
	// User who responded to the review
	ReplyUserId *uuid.UUID `json:"replyUserId,omitempty"`
	// Review Event Type - Review
	ResourceType *ReviewResourceType `json:"resourceType,omitempty"`
	// Review update date
	ReviewDate *azuredevops.Time `json:"reviewDate,omitempty"`
	// ReviewId of the review  on which the operation is performed
	ReviewId *uint64 `json:"reviewId,omitempty"`
	// Text in Review Text
	ReviewText *string `json:"reviewText,omitempty"`
	// User display name at the time of review
	UserDisplayName *string `json:"userDisplayName,omitempty"`
	// User who gave review
	UserId *uuid.UUID `json:"userId,omitempty"`
}

// [Flags] Options to GetReviews query
type ReviewFilterOptions string

type reviewFilterOptionsValuesType struct {
	None                 ReviewFilterOptions
	FilterEmptyReviews   ReviewFilterOptions
	FilterEmptyUserNames ReviewFilterOptions
}

var ReviewFilterOptionsValues = reviewFilterOptionsValuesType{
	// No filtering, all reviews are returned (default option)
	None: "none",
	// Filter out review items with empty review text
	FilterEmptyReviews: "filterEmptyReviews",
	// Filter out review items with empty usernames
	FilterEmptyUserNames: "filterEmptyUserNames",
}

type ReviewPatch struct {
	// Denotes the patch operation type
	Operation *ReviewPatchOperation `json:"operation,omitempty"`
	// Use when patch operation is FlagReview
	ReportedConcern *UserReportedConcern `json:"reportedConcern,omitempty"`
	// Use when patch operation is EditReview
	ReviewItem *Review `json:"reviewItem,omitempty"`
}

// Denotes the patch operation type
type ReviewPatchOperation string

type reviewPatchOperationValuesType struct {
	FlagReview             ReviewPatchOperation
	UpdateReview           ReviewPatchOperation
	ReplyToReview          ReviewPatchOperation
	AdminResponseForReview ReviewPatchOperation
	DeleteAdminReply       ReviewPatchOperation
	DeletePublisherReply   ReviewPatchOperation
}

var ReviewPatchOperationValues = reviewPatchOperationValuesType{
	// Flag a review
	FlagReview: "flagReview",
	// Update an existing review
	UpdateReview: "updateReview",
	// Submit a reply for a review
	ReplyToReview: "replyToReview",
	// Submit an admin response
	AdminResponseForReview: "adminResponseForReview",
	// Delete an Admin Reply
	DeleteAdminReply: "deleteAdminReply",
	// Delete Publisher Reply
	DeletePublisherReply: "deletePublisherReply",
}

type ReviewReply struct {
	// Id of the reply
	Id *uint64 `json:"id,omitempty"`
	// Flag for soft deletion
	IsDeleted *bool `json:"isDeleted,omitempty"`
	// Version of the product when the reply was submitted or updated
	ProductVersion *string `json:"productVersion,omitempty"`
	// Content of the reply
	ReplyText *string `json:"replyText,omitempty"`
	// Id of the review, to which this reply belongs
	ReviewId *uint64 `json:"reviewId,omitempty"`
	// Title of the reply
	Title *string `json:"title,omitempty"`
	// Date the reply was submitted or updated
	UpdatedDate *azuredevops.Time `json:"updatedDate,omitempty"`
	// Id of the user who left the reply
	UserId *uuid.UUID `json:"userId,omitempty"`
}

// Type of event
type ReviewResourceType string

type reviewResourceTypeValuesType struct {
	Review         ReviewResourceType
	PublisherReply ReviewResourceType
	AdminReply     ReviewResourceType
}

var ReviewResourceTypeValues = reviewResourceTypeValuesType{
	Review:         "review",
	PublisherReply: "publisherReply",
	AdminReply:     "adminReply",
}

type ReviewsResult struct {
	// Flag indicating if there are more reviews to be shown (for paging)
	HasMoreReviews *bool `json:"hasMoreReviews,omitempty"`
	// List of reviews
	Reviews *[]Review `json:"reviews,omitempty"`
	// Count of total review items
	TotalReviewCount *uint64 `json:"totalReviewCount,omitempty"`
}

type ReviewSummary struct {
	// Average Rating
	AverageRating *float32 `json:"averageRating,omitempty"`
	// Count of total ratings
	RatingCount *uint64 `json:"ratingCount,omitempty"`
	// Split of count across rating
	RatingSplit *[]RatingCountPerRating `json:"ratingSplit,omitempty"`
}

// Defines the sort order that can be defined for Extensions query
type SortByType string

type sortByTypeValuesType struct {
	Relevance       SortByType
	LastUpdatedDate SortByType
	Title           SortByType
	Publisher       SortByType
	InstallCount    SortByType
	PublishedDate   SortByType
	AverageRating   SortByType
	TrendingDaily   SortByType
	TrendingWeekly  SortByType
	TrendingMonthly SortByType
	ReleaseDate     SortByType
	Author          SortByType
	WeightedRating  SortByType
}

var SortByTypeValues = sortByTypeValuesType{
	// The results will be sorted by relevance in case search query is given, if no search query resutls will be provided as is
	Relevance: "relevance",
	// The results will be sorted as per Last Updated date of the extensions with recently updated at the top
	LastUpdatedDate: "lastUpdatedDate",
	// Results will be sorted Alphabetically as per the title of the extension
	Title: "title",
	// Results will be sorted Alphabetically as per Publisher title
	Publisher: "publisher",
	// Results will be sorted by Install Count
	InstallCount: "installCount",
	// The results will be sorted as per Published date of the extensions
	PublishedDate: "publishedDate",
	// The results will be sorted as per Average ratings of the extensions
	AverageRating: "averageRating",
	// The results will be sorted as per Trending Daily Score of the extensions
	TrendingDaily: "trendingDaily",
	// The results will be sorted as per Trending weekly Score of the extensions
	TrendingWeekly: "trendingWeekly",
	// The results will be sorted as per Trending monthly Score of the extensions
	TrendingMonthly: "trendingMonthly",
	// The results will be sorted as per ReleaseDate of the extensions (date on which the extension first went public)
	ReleaseDate: "releaseDate",
	// The results will be sorted as per Author defined in the VSix/Metadata. If not defined, publisher name is used This is specifically needed by VS IDE, other (new and old) clients are not encouraged to use this
	Author: "author",
	// The results will be sorted as per Weighted Rating of the extension.
	WeightedRating: "weightedRating",
}

// Defines the sort order that can be defined for Extensions query
type SortOrderType string

type sortOrderTypeValuesType struct {
	Default    SortOrderType
	Ascending  SortOrderType
	Descending SortOrderType
}

var SortOrderTypeValues = sortOrderTypeValuesType{
	// Results will be sorted in the default order as per the sorting type defined. The default varies for each type, e.g. for Relevance, default is Descending, for Title default is Ascending etc.
	Default: "default",
	// The results will be sorted in Ascending order
	Ascending: "ascending",
	// The results will be sorted in Descending order
	Descending: "descending",
}

type UnpackagedExtensionData struct {
	Categories            *[]string             `json:"categories,omitempty"`
	Description           *string               `json:"description,omitempty"`
	DisplayName           *string               `json:"displayName,omitempty"`
	DraftId               *uuid.UUID            `json:"draftId,omitempty"`
	ExtensionName         *string               `json:"extensionName,omitempty"`
	InstallationTargets   *[]InstallationTarget `json:"installationTargets,omitempty"`
	IsConvertedToMarkdown *bool                 `json:"isConvertedToMarkdown,omitempty"`
	IsPreview             *bool                 `json:"isPreview,omitempty"`
	PricingCategory       *string               `json:"pricingCategory,omitempty"`
	Product               *string               `json:"product,omitempty"`
	PublisherName         *string               `json:"publisherName,omitempty"`
	QnAEnabled            *bool                 `json:"qnAEnabled,omitempty"`
	ReferralUrl           *string               `json:"referralUrl,omitempty"`
	RepositoryUrl         *string               `json:"repositoryUrl,omitempty"`
	Tags                  *[]string             `json:"tags,omitempty"`
	Version               *string               `json:"version,omitempty"`
	VsixId                *string               `json:"vsixId,omitempty"`
}

// Represents the extension policy applied to a given user
type UserExtensionPolicy struct {
	// User display name that this policy refers to
	DisplayName *string `json:"displayName,omitempty"`
	// The extension policy applied to the user
	Permissions *ExtensionPolicy `json:"permissions,omitempty"`
	// User id that this policy refers to
	UserId *string `json:"userId,omitempty"`
}

// Identity reference with name and guid
type UserIdentityRef struct {
	// User display name
	DisplayName *string `json:"displayName,omitempty"`
	// User VSID
	Id *uuid.UUID `json:"id,omitempty"`
}

type UserReportedConcern struct {
	// Category of the concern
	Category *ConcernCategory `json:"category,omitempty"`
	// User comment associated with the report
	ConcernText *string `json:"concernText,omitempty"`
	// Id of the review which was reported
	ReviewId *uint64 `json:"reviewId,omitempty"`
	// Date the report was submitted
	SubmittedDate *azuredevops.Time `json:"submittedDate,omitempty"`
	// Id of the user who reported a review
	UserId *uuid.UUID `json:"userId,omitempty"`
}

type VSCodeWebExtensionStatisicsType string

type vsCodeWebExtensionStatisicsTypeValuesType struct {
	Install   VSCodeWebExtensionStatisicsType
	Update    VSCodeWebExtensionStatisicsType
	Uninstall VSCodeWebExtensionStatisicsType
}

var VSCodeWebExtensionStatisicsTypeValues = vsCodeWebExtensionStatisicsTypeValuesType{
	Install:   "install",
	Update:    "update",
	Uninstall: "uninstall",
}

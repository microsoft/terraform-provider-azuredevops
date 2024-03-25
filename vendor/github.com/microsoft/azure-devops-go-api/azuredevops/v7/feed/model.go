// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package feed

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
)

type BuildPackage struct {
	// Display name of the feed.
	FeedName *string `json:"feedName,omitempty"`
	// Package version description.
	PackageDescription *string `json:"packageDescription,omitempty"`
	// Display name of the package.
	PackageName *string `json:"packageName,omitempty"`
	// Version of the package.
	PackageVersion *string `json:"packageVersion,omitempty"`
	// TFS project id.
	ProjectId *uuid.UUID `json:"projectId,omitempty"`
	// Type of the package.
	ProtocolType *string `json:"protocolType,omitempty"`
}

// A container for artifacts.
type Feed struct {
	// Supported capabilities of a feed.
	Capabilities *FeedCapabilities `json:"capabilities,omitempty"`
	// This will either be the feed GUID or the feed GUID and view GUID depending on how the feed was accessed.
	FullyQualifiedId *string `json:"fullyQualifiedId,omitempty"`
	// Full name of the view, in feed@view format.
	FullyQualifiedName *string `json:"fullyQualifiedName,omitempty"`
	// A GUID that uniquely identifies this feed.
	Id *uuid.UUID `json:"id,omitempty"`
	// If set, all packages in the feed are immutable.  It is important to note that feed views are immutable; therefore, this flag will always be set for views.
	IsReadOnly *bool `json:"isReadOnly,omitempty"`
	// A name for the feed. feed names must follow these rules: <list type="bullet"><item><description> Must not exceed 64 characters </description></item><item><description> Must not contain whitespaces </description></item><item><description> Must not start with an underscore or a period </description></item><item><description> Must not end with a period </description></item><item><description> Must not contain any of the following illegal characters: <![CDATA[ @, ~, ;, {, }, \, +, =, <, >, |, /, \\, ?, :, &, $, *, \", #, [, ] ]]></description></item></list>
	Name *string `json:"name,omitempty"`
	// The project that this feed is associated with.
	Project *ProjectReference `json:"project,omitempty"`
	// This should always be true. Setting to false will override all sources in UpstreamSources.
	UpstreamEnabled *bool `json:"upstreamEnabled,omitempty"`
	// A list of sources that this feed will fetch packages from.  An empty list indicates that this feed will not search any additional sources for packages.
	UpstreamSources *[]UpstreamSource `json:"upstreamSources,omitempty"`
	// Definition of the view.
	View *FeedView `json:"view,omitempty"`
	// View Id.
	ViewId *uuid.UUID `json:"viewId,omitempty"`
	// View name.
	ViewName *string `json:"viewName,omitempty"`
	// Related REST links.
	Links interface{} `json:"_links,omitempty"`
	// If set, this feed supports generation of package badges.
	BadgesEnabled *bool `json:"badgesEnabled,omitempty"`
	// The view that the feed administrator has indicated is the default experience for readers.
	DefaultViewId *uuid.UUID `json:"defaultViewId,omitempty"`
	// The date that this feed was deleted.
	DeletedDate *azuredevops.Time `json:"deletedDate,omitempty"`
	// A description for the feed.  Descriptions must not exceed 255 characters.
	Description *string `json:"description,omitempty"`
	// If set, the feed will hide all deleted/unpublished versions
	HideDeletedPackageVersions *bool `json:"hideDeletedPackageVersions,omitempty"`
	// The date that this feed was permanently deleted.
	PermanentDeletedDate *azuredevops.Time `json:"permanentDeletedDate,omitempty"`
	// Explicit permissions for the feed.
	Permissions *[]FeedPermission `json:"permissions,omitempty"`
	// The date that this feed is scheduled to be permanently deleted.
	ScheduledPermanentDeleteDate *azuredevops.Time `json:"scheduledPermanentDeleteDate,omitempty"`
	// If set, time that the UpstreamEnabled property was changed. Will be null if UpstreamEnabled was never changed after Feed creation.
	UpstreamEnabledChangedDate *azuredevops.Time `json:"upstreamEnabledChangedDate,omitempty"`
	// The URL of the base feed in GUID form.
	Url *string `json:"url,omitempty"`
}

type FeedBatchOperation string

type feedBatchOperationValuesType struct {
	SaveCachedPackages FeedBatchOperation
}

var FeedBatchOperationValues = feedBatchOperationValuesType{
	SaveCachedPackages: "saveCachedPackages",
}

// [Flags] Capabilities are used to track features that are available to individual feeds. In general, newly created feeds should be given all available capabilities. These flags track breaking changes in behaviour to feeds, or changes that require user reaction.
type FeedCapabilities string

type feedCapabilitiesValuesType struct {
	None                FeedCapabilities
	UpstreamV2          FeedCapabilities
	UnderMaintenance    FeedCapabilities
	DefaultCapabilities FeedCapabilities
}

var FeedCapabilitiesValues = feedCapabilitiesValuesType{
	// No flags exist for this feed
	None: "none",
	// This feed can serve packages from upstream sources Upstream packages must be manually promoted to views
	UpstreamV2: "upstreamV2",
	// This feed is currently under maintenance and may have reduced functionality
	UnderMaintenance: "underMaintenance",
	// The capabilities given to a newly created feed
	DefaultCapabilities: "defaultCapabilities",
}

// An object that contains all of the settings for a specific feed.
type FeedCore struct {
	// Supported capabilities of a feed.
	Capabilities *FeedCapabilities `json:"capabilities,omitempty"`
	// This will either be the feed GUID or the feed GUID and view GUID depending on how the feed was accessed.
	FullyQualifiedId *string `json:"fullyQualifiedId,omitempty"`
	// Full name of the view, in feed@view format.
	FullyQualifiedName *string `json:"fullyQualifiedName,omitempty"`
	// A GUID that uniquely identifies this feed.
	Id *uuid.UUID `json:"id,omitempty"`
	// If set, all packages in the feed are immutable.  It is important to note that feed views are immutable; therefore, this flag will always be set for views.
	IsReadOnly *bool `json:"isReadOnly,omitempty"`
	// A name for the feed. feed names must follow these rules: <list type="bullet"><item><description> Must not exceed 64 characters </description></item><item><description> Must not contain whitespaces </description></item><item><description> Must not start with an underscore or a period </description></item><item><description> Must not end with a period </description></item><item><description> Must not contain any of the following illegal characters: <![CDATA[ @, ~, ;, {, }, \, +, =, <, >, |, /, \\, ?, :, &, $, *, \", #, [, ] ]]></description></item></list>
	Name *string `json:"name,omitempty"`
	// The project that this feed is associated with.
	Project *ProjectReference `json:"project,omitempty"`
	// This should always be true. Setting to false will override all sources in UpstreamSources.
	UpstreamEnabled *bool `json:"upstreamEnabled,omitempty"`
	// A list of sources that this feed will fetch packages from.  An empty list indicates that this feed will not search any additional sources for packages.
	UpstreamSources *[]UpstreamSource `json:"upstreamSources,omitempty"`
	// Definition of the view.
	View *FeedView `json:"view,omitempty"`
	// View Id.
	ViewId *uuid.UUID `json:"viewId,omitempty"`
	// View name.
	ViewName *string `json:"viewName,omitempty"`
}

// A container that encapsulates the state of the feed after a create, update, or delete.
type FeedChange struct {
	// The state of the feed after a after a create, update, or delete operation completed.
	Feed *Feed `json:"feed,omitempty"`
	// A token that identifies the next change in the log of changes.
	FeedContinuationToken *uint64 `json:"feedContinuationToken,omitempty"`
	// The type of operation.
	ChangeType *ChangeType `json:"changeType,omitempty"`
	// A token that identifies the latest package change for this feed.  This can be used to quickly determine if there have been any changes to packages in a specific feed.
	LatestPackageContinuationToken *uint64 `json:"latestPackageContinuationToken,omitempty"`
}

// A result set containing the feed changes for the range that was requested.
type FeedChangesResponse struct {
	Links interface{} `json:"_links,omitempty"`
	// The number of changes in this set.
	Count *int `json:"count,omitempty"`
	// A container that encapsulates the state of the feed after a create, update, or delete.
	FeedChanges *[]FeedChange `json:"feedChanges,omitempty"`
	// When iterating through the log of changes this value indicates the value that should be used for the next continuation token.
	NextFeedContinuationToken *uint64 `json:"nextFeedContinuationToken,omitempty"`
}

type FeedIdsResult struct {
	Id          *uuid.UUID `json:"id,omitempty"`
	Name        *string    `json:"name,omitempty"`
	ProjectId   *uuid.UUID `json:"projectId,omitempty"`
	ProjectName *string    `json:"projectName,omitempty"`
}

// Permissions for a feed.
type FeedPermission struct {
	// Display name for the identity.
	DisplayName *string `json:"displayName,omitempty"`
	// Identity associated with this role.
	IdentityDescriptor *string `json:"identityDescriptor,omitempty"`
	// Id of the identity associated with this role.
	IdentityId *uuid.UUID `json:"identityId,omitempty"`
	// Boolean indicating whether the role is inherited or set directly.
	IsInheritedRole *bool `json:"isInheritedRole,omitempty"`
	// The role for this identity on a feed.
	Role *FeedRole `json:"role,omitempty"`
}

// Retention policy settings.
type FeedRetentionPolicy struct {
	// This attribute is deprecated and is not honoured by retention
	AgeLimitInDays *int `json:"ageLimitInDays,omitempty"`
	// Maximum versions to preserve per package and package type.
	CountLimit *int `json:"countLimit,omitempty"`
	// Number of days to preserve a package version after its latest download.
	DaysToKeepRecentlyDownloadedPackages *int `json:"daysToKeepRecentlyDownloadedPackages,omitempty"`
}

type FeedRole string

type feedRoleValuesType struct {
	Custom        FeedRole
	None          FeedRole
	Reader        FeedRole
	Contributor   FeedRole
	Administrator FeedRole
	Collaborator  FeedRole
}

var FeedRoleValues = feedRoleValuesType{
	// Unsupported.
	Custom: "custom",
	// Unsupported.
	None: "none",
	// Readers can only read packages and view settings.
	Reader: "reader",
	// Contributors can do anything to packages in the feed including adding new packages, but they may not modify feed settings.
	Contributor: "contributor",
	// Administrators have total control over the feed.
	Administrator: "administrator",
	// Collaborators have the same permissions as readers, but can also ingest packages from configured upstream sources.
	Collaborator: "collaborator",
}

// Update a feed definition with these new values.
type FeedUpdate struct {
	// If set, the feed will allow upload of packages that exist on the upstream
	AllowUpstreamNameConflict *bool `json:"allowUpstreamNameConflict,omitempty"`
	// If set, this feed supports generation of package badges.
	BadgesEnabled *bool `json:"badgesEnabled,omitempty"`
	// The view that the feed administrator has indicated is the default experience for readers.
	DefaultViewId *uuid.UUID `json:"defaultViewId,omitempty"`
	// A description for the feed.  Descriptions must not exceed 255 characters.
	Description *string `json:"description,omitempty"`
	// If set, feed will hide all deleted/unpublished versions
	HideDeletedPackageVersions *bool `json:"hideDeletedPackageVersions,omitempty"`
	// A GUID that uniquely identifies this feed.
	Id *uuid.UUID `json:"id,omitempty"`
	// A name for the feed. feed names must follow these rules: <list type="bullet"><item><description> Must not exceed 64 characters </description></item><item><description> Must not contain whitespaces </description></item><item><description> Must not start with an underscore or a period </description></item><item><description> Must not end with a period </description></item><item><description> Must not contain any of the following illegal characters: <![CDATA[ @, ~, ;, {, }, \, +, =, <, >, |, /, \\, ?, :, &, $, *, \", #, [, ] ]]></description></item></list>
	Name *string `json:"name,omitempty"`
	// If set, the feed can proxy packages from an upstream feed
	UpstreamEnabled *bool `json:"upstreamEnabled,omitempty"`
	// A list of sources that this feed will fetch packages from.  An empty list indicates that this feed will not search any additional sources for packages.
	UpstreamSources *[]UpstreamSource `json:"upstreamSources,omitempty"`
}

// A view on top of a feed.
type FeedView struct {
	// Related REST links.
	Links interface{} `json:"_links,omitempty"`
	// Id of the view.
	Id *uuid.UUID `json:"id,omitempty"`
	// Name of the view.
	Name *string `json:"name,omitempty"`
	// Type of view.
	Type *FeedViewType `json:"type,omitempty"`
	// Url of the view.
	Url *string `json:"url,omitempty"`
	// Visibility status of the view.
	Visibility *FeedVisibility `json:"visibility,omitempty"`
}

// The type of view, often used to control capabilities and exposure to options such as promote.  Implicit views are internally created only.
type FeedViewType string

type feedViewTypeValuesType struct {
	None     FeedViewType
	Release  FeedViewType
	Implicit FeedViewType
}

var FeedViewTypeValues = feedViewTypeValuesType{
	// Default, unspecified view type.
	None: "none",
	// View used as a promotion destination to classify released artifacts.
	Release: "release",
	// Internal view type that is automatically created and managed by the system.
	Implicit: "implicit",
}

// Feed visibility controls the scope in which a certain feed is accessible by a particular user
type FeedVisibility string

type feedVisibilityValuesType struct {
	Private      FeedVisibility
	Collection   FeedVisibility
	Organization FeedVisibility
	AadTenant    FeedVisibility
}

var FeedVisibilityValues = feedVisibilityValuesType{
	// Only accessible by the permissions explicitly set by the feed administrator.
	Private: "private",
	// Feed is accessible by all the valid users present in the organization where the feed resides (for example across organization 'myorg' at 'dev.azure.com/myorg')
	Collection: "collection",
	// Feed is accessible by all the valid users present in the enterprise where the feed resides. Note that legacy naming and back compat leaves the name of this value out of sync with its new meaning.
	Organization: "organization",
	// Feed is accessible by all the valid users present in the Azure Active Directory tenant.
	AadTenant: "aadTenant",
}

// Permissions for feed service-wide operations such as the creation of new feeds.
type GlobalPermission struct {
	// Identity of the user with the provided Role.
	IdentityDescriptor *string `json:"identityDescriptor,omitempty"`
	// IdentityId corresponding to the IdentityDescriptor
	IdentityId *uuid.UUID `json:"identityId,omitempty"`
	// Role associated with the Identity.
	Role *GlobalRole `json:"role,omitempty"`
}

type GlobalRole string

type globalRoleValuesType struct {
	Custom        GlobalRole
	None          GlobalRole
	FeedCreator   GlobalRole
	Administrator GlobalRole
}

var GlobalRoleValues = globalRoleValuesType{
	// Invalid default value.
	Custom: "custom",
	// Explicit no permissions.
	None: "none",
	// Ability to create new feeds.
	FeedCreator: "feedCreator",
	// Read and manage any feed
	Administrator: "administrator",
}

// Type of operation last performed.
type ChangeType string

type changeTypeValuesType struct {
	AddOrUpdate     ChangeType
	Delete          ChangeType
	PermanentDelete ChangeType
}

var ChangeTypeValues = changeTypeValuesType{
	// A package version was added or updated.
	AddOrUpdate: "addOrUpdate",
	// A package version was deleted.
	Delete: "delete",
	// A feed was permanently deleted. This is not used for package version.
	PermanentDelete: "permanentDelete",
}

// Core data about any package, including its id and version information and basic state.
type MinimalPackageVersion struct {
	// Upstream source this package was ingested from.
	DirectUpstreamSourceId *uuid.UUID `json:"directUpstreamSourceId,omitempty"`
	// Id for the package.
	Id *uuid.UUID `json:"id,omitempty"`
	// [Obsolete] Used for legacy scenarios and may be removed in future versions.
	IsCachedVersion *bool `json:"isCachedVersion,omitempty"`
	// True if this package has been deleted.
	IsDeleted *bool `json:"isDeleted,omitempty"`
	// True if this is the latest version of the package by package type sort order.
	IsLatest *bool `json:"isLatest,omitempty"`
	// (NuGet and Cargo Only) True if this package is listed.
	IsListed *bool `json:"isListed,omitempty"`
	// Normalized version using normalization rules specific to a package type.
	NormalizedVersion *string `json:"normalizedVersion,omitempty"`
	// Package description.
	PackageDescription *string `json:"packageDescription,omitempty"`
	// UTC Date the package was published to the service.
	PublishDate *azuredevops.Time `json:"publishDate,omitempty"`
	// Internal storage id.
	StorageId *string `json:"storageId,omitempty"`
	// Display version.
	Version *string `json:"version,omitempty"`
	// List of views containing this package version.
	Views *[]FeedView `json:"views,omitempty"`
}

// A package, which is a container for one or more package versions.
type Package struct {
	// Related REST links.
	Links interface{} `json:"_links,omitempty"`
	// Id of the package.
	Id *uuid.UUID `json:"id,omitempty"`
	// Used for legacy scenarios and may be removed in future versions.
	IsCached *bool `json:"isCached,omitempty"`
	// The display name of the package.
	Name *string `json:"name,omitempty"`
	// The normalized name representing the identity of this package within its package type.
	NormalizedName *string `json:"normalizedName,omitempty"`
	// Type of the package.
	ProtocolType *string `json:"protocolType,omitempty"`
	// [Obsolete] - this field is unused and will be removed in a future release.
	StarCount *int `json:"starCount,omitempty"`
	// Url for this package.
	Url *string `json:"url,omitempty"`
	// All versions for this package within its feed.
	Versions *[]MinimalPackageVersion `json:"versions,omitempty"`
}

// A dependency on another package version.
type PackageDependency struct {
	// Dependency package group (an optional classification within some package types).
	Group *string `json:"group,omitempty"`
	// Dependency package name.
	PackageName *string `json:"packageName,omitempty"`
	// Dependency package version range.
	VersionRange *string `json:"versionRange,omitempty"`
}

// A package file for a specific package version, only relevant to package types that contain multiple files per version.
type PackageFile struct {
	// Hierarchical representation of files.
	Children *[]PackageFile `json:"children,omitempty"`
	// File name.
	Name *string `json:"name,omitempty"`
	// Extended data unique to a specific package type.
	ProtocolMetadata *ProtocolMetadata `json:"protocolMetadata,omitempty"`
}

// A single change to a feed's packages.
type PackageChange struct {
	// Package that was changed.
	Package *Package `json:"package,omitempty"`
	// Change that was performed on a package version.
	PackageVersionChange *PackageVersionChange `json:"packageVersionChange,omitempty"`
}

// A set of change operations to a feed's packages.
type PackageChangesResponse struct {
	// Related REST links.
	Links interface{} `json:"_links,omitempty"`
	// Number of changes in this batch.
	Count *int `json:"count,omitempty"`
	// Token that should be used in future calls for this feed to retrieve new changes.
	NextPackageContinuationToken *uint64 `json:"nextPackageContinuationToken,omitempty"`
	// List of changes.
	PackageChanges *[]PackageChange `json:"packageChanges,omitempty"`
}

// All metrics for a certain package id
type PackageMetrics struct {
	// Total count of downloads per package id.
	DownloadCount *float64 `json:"downloadCount,omitempty"`
	// Number of downloads per unique user per package id.
	DownloadUniqueUsers *float64 `json:"downloadUniqueUsers,omitempty"`
	// UTC date and time when package was last downloaded.
	LastDownloaded *azuredevops.Time `json:"lastDownloaded,omitempty"`
	// Package id.
	PackageId *uuid.UUID `json:"packageId,omitempty"`
}

// Query to get package metrics
type PackageMetricsQuery struct {
	// List of package ids
	PackageIds *[]uuid.UUID `json:"packageIds,omitempty"`
}

// A specific version of a package.
type PackageVersion struct {
	// Upstream source this package was ingested from.
	DirectUpstreamSourceId *uuid.UUID `json:"directUpstreamSourceId,omitempty"`
	// Id for the package.
	Id *uuid.UUID `json:"id,omitempty"`
	// [Obsolete] Used for legacy scenarios and may be removed in future versions.
	IsCachedVersion *bool `json:"isCachedVersion,omitempty"`
	// True if this package has been deleted.
	IsDeleted *bool `json:"isDeleted,omitempty"`
	// True if this is the latest version of the package by package type sort order.
	IsLatest *bool `json:"isLatest,omitempty"`
	// (NuGet and Cargo Only) True if this package is listed.
	IsListed *bool `json:"isListed,omitempty"`
	// Normalized version using normalization rules specific to a package type.
	NormalizedVersion *string `json:"normalizedVersion,omitempty"`
	// Package description.
	PackageDescription *string `json:"packageDescription,omitempty"`
	// UTC Date the package was published to the service.
	PublishDate *azuredevops.Time `json:"publishDate,omitempty"`
	// Internal storage id.
	StorageId *string `json:"storageId,omitempty"`
	// Display version.
	Version *string `json:"version,omitempty"`
	// List of views containing this package version.
	Views *[]FeedView `json:"views,omitempty"`
	// Related links
	Links interface{} `json:"_links,omitempty"`
	// Package version author.
	Author *string `json:"author,omitempty"`
	// UTC date that this package version was deleted.
	DeletedDate *azuredevops.Time `json:"deletedDate,omitempty"`
	// List of dependencies for this package version.
	Dependencies *[]PackageDependency `json:"dependencies,omitempty"`
	// Package version description.
	Description *string `json:"description,omitempty"`
	// Files associated with this package version, only relevant for multi-file package types.
	Files *[]PackageFile `json:"files,omitempty"`
	// Other versions of this package.
	OtherVersions *[]MinimalPackageVersion `json:"otherVersions,omitempty"`
	// Extended data specific to a package type.
	ProtocolMetadata *ProtocolMetadata `json:"protocolMetadata,omitempty"`
	// List of upstream sources through which a package version moved to land in this feed.
	SourceChain *[]UpstreamSource `json:"sourceChain,omitempty"`
	// Package version summary.
	Summary *string `json:"summary,omitempty"`
	// Package version tags.
	Tags *[]string `json:"tags,omitempty"`
	// Package version url.
	Url *string `json:"url,omitempty"`
}

// A change to a single package version.
type PackageVersionChange struct {
	// Token marker for this change, allowing the caller to send this value back to the service and receive changes beyond this one.
	ContinuationToken *uint64 `json:"continuationToken,omitempty"`
	// The type of change that was performed.
	ChangeType *ChangeType `json:"changeType,omitempty"`
	// Package version that was changed.
	PackageVersion *PackageVersion `json:"packageVersion,omitempty"`
}

// All metrics for a certain package version id
type PackageVersionMetrics struct {
	// Total count of downloads per package version id.
	DownloadCount *float64 `json:"downloadCount,omitempty"`
	// Number of downloads per unique user per package version id.
	DownloadUniqueUsers *float64 `json:"downloadUniqueUsers,omitempty"`
	// UTC date and time when package version was last downloaded.
	LastDownloaded *azuredevops.Time `json:"lastDownloaded,omitempty"`
	// Package id.
	PackageId *uuid.UUID `json:"packageId,omitempty"`
	// Package version id.
	PackageVersionId *uuid.UUID `json:"packageVersionId,omitempty"`
}

// Query to get package version metrics
type PackageVersionMetricsQuery struct {
	// List of package version ids
	PackageVersionIds *[]uuid.UUID `json:"packageVersionIds,omitempty"`
}

// Provenance for a published package version
type PackageVersionProvenance struct {
	// Name or Id of the feed.
	FeedId *uuid.UUID `json:"feedId,omitempty"`
	// Id of the package (GUID Id, not name).
	PackageId *uuid.UUID `json:"packageId,omitempty"`
	// Id of the package version (GUID Id, not name).
	PackageVersionId *uuid.UUID `json:"packageVersionId,omitempty"`
	// Provenance information for this package version.
	Provenance *Provenance `json:"provenance,omitempty"`
}

type ProjectReference struct {
	// Gets or sets id of the project.
	Id *uuid.UUID `json:"id,omitempty"`
	// Gets or sets name of the project.
	Name *string `json:"name,omitempty"`
	// Gets or sets visibility of the project.
	Visibility *string `json:"visibility,omitempty"`
}

// Extended metadata for a specific package type.
type ProtocolMetadata struct {
	// Extended metadata for a specific package type, formatted to the associated schema version definition.
	Data interface{} `json:"data,omitempty"`
	// Schema version.
	SchemaVersion *int `json:"schemaVersion,omitempty"`
}

// Data about the origin of a published package
type Provenance struct {
	// Other provenance data.
	Data *map[string]string `json:"data,omitempty"`
	// Type of provenance source, for example "InternalBuild", "InternalRelease"
	ProvenanceSource *string `json:"provenanceSource,omitempty"`
	// Identity of user that published the package
	PublisherUserIdentity *uuid.UUID `json:"publisherUserIdentity,omitempty"`
	// HTTP User-Agent used when pushing the package.
	UserAgent *string `json:"userAgent,omitempty"`
}

// A single package version within the recycle bin.
type RecycleBinPackageVersion struct {
	// Upstream source this package was ingested from.
	DirectUpstreamSourceId *uuid.UUID `json:"directUpstreamSourceId,omitempty"`
	// Id for the package.
	Id *uuid.UUID `json:"id,omitempty"`
	// [Obsolete] Used for legacy scenarios and may be removed in future versions.
	IsCachedVersion *bool `json:"isCachedVersion,omitempty"`
	// True if this package has been deleted.
	IsDeleted *bool `json:"isDeleted,omitempty"`
	// True if this is the latest version of the package by package type sort order.
	IsLatest *bool `json:"isLatest,omitempty"`
	// (NuGet and Cargo Only) True if this package is listed.
	IsListed *bool `json:"isListed,omitempty"`
	// Normalized version using normalization rules specific to a package type.
	NormalizedVersion *string `json:"normalizedVersion,omitempty"`
	// Package description.
	PackageDescription *string `json:"packageDescription,omitempty"`
	// UTC Date the package was published to the service.
	PublishDate *azuredevops.Time `json:"publishDate,omitempty"`
	// Internal storage id.
	StorageId *string `json:"storageId,omitempty"`
	// Display version.
	Version *string `json:"version,omitempty"`
	// List of views containing this package version.
	Views *[]FeedView `json:"views,omitempty"`
	// Related links
	Links interface{} `json:"_links,omitempty"`
	// Package version author.
	Author *string `json:"author,omitempty"`
	// UTC date that this package version was deleted.
	DeletedDate *azuredevops.Time `json:"deletedDate,omitempty"`
	// List of dependencies for this package version.
	Dependencies *[]PackageDependency `json:"dependencies,omitempty"`
	// Package version description.
	Description *string `json:"description,omitempty"`
	// Files associated with this package version, only relevant for multi-file package types.
	Files *[]PackageFile `json:"files,omitempty"`
	// Other versions of this package.
	OtherVersions *[]MinimalPackageVersion `json:"otherVersions,omitempty"`
	// Extended data specific to a package type.
	ProtocolMetadata *ProtocolMetadata `json:"protocolMetadata,omitempty"`
	// List of upstream sources through which a package version moved to land in this feed.
	SourceChain *[]UpstreamSource `json:"sourceChain,omitempty"`
	// Package version summary.
	Summary *string `json:"summary,omitempty"`
	// Package version tags.
	Tags *[]string `json:"tags,omitempty"`
	// Package version url.
	Url *string `json:"url,omitempty"`
	// UTC date on which the package will automatically be removed from the recycle bin and permanently deleted.
	ScheduledPermanentDeleteDate *azuredevops.Time `json:"scheduledPermanentDeleteDate,omitempty"`
}

// Upstream source definition, including its Identity, package type, and other associated information.
type UpstreamSource struct {
	// UTC date that this upstream was deleted.
	DeletedDate *azuredevops.Time `json:"deletedDate,omitempty"`
	// Locator for connecting to the upstream source in a user friendly format, that may potentially change over time
	DisplayLocation *string `json:"displayLocation,omitempty"`
	// Identity of the upstream source.
	Id *uuid.UUID `json:"id,omitempty"`
	// For an internal upstream type, track the Azure DevOps organization that contains it.
	InternalUpstreamCollectionId *uuid.UUID `json:"internalUpstreamCollectionId,omitempty"`
	// For an internal upstream type, track the feed id being referenced.
	InternalUpstreamFeedId *uuid.UUID `json:"internalUpstreamFeedId,omitempty"`
	// For an internal upstream type, track the project of the feed being referenced.
	InternalUpstreamProjectId *uuid.UUID `json:"internalUpstreamProjectId,omitempty"`
	// For an internal upstream type, track the view of the feed being referenced.
	InternalUpstreamViewId *uuid.UUID `json:"internalUpstreamViewId,omitempty"`
	// Consistent locator for connecting to the upstream source.
	Location *string `json:"location,omitempty"`
	// Display name.
	Name *string `json:"name,omitempty"`
	// Package type associated with the upstream source.
	Protocol *string `json:"protocol,omitempty"`
	// The identity of the service endpoint that holds credentials to use when accessing the upstream.
	ServiceEndpointId *uuid.UUID `json:"serviceEndpointId,omitempty"`
	// Specifies the projectId of the Service Endpoint.
	ServiceEndpointProjectId *uuid.UUID `json:"serviceEndpointProjectId,omitempty"`
	// Specifies the status of the upstream.
	Status *UpstreamStatus `json:"status,omitempty"`
	// Provides a human-readable reason for the status of the upstream.
	StatusDetails *[]UpstreamStatusDetail `json:"statusDetails,omitempty"`
	// Source type, such as Public or Internal.
	UpstreamSourceType *UpstreamSourceType `json:"upstreamSourceType,omitempty"`
}

// Type of an upstream source, such as Public or Internal.
type UpstreamSourceType string

type upstreamSourceTypeValuesType struct {
	Public   UpstreamSourceType
	Internal UpstreamSourceType
}

var UpstreamSourceTypeValues = upstreamSourceTypeValuesType{
	// Publicly available source.
	Public: "public",
	// Azure DevOps upstream source.
	Internal: "internal",
}

// Status of the upstream, such as Ok or Disabled.
type UpstreamStatus string

type upstreamStatusValuesType struct {
	Ok       UpstreamStatus
	Disabled UpstreamStatus
}

var UpstreamStatusValues = upstreamStatusValuesType{
	// Upstream source is ok.
	Ok: "ok",
	// Upstream source is disabled.
	Disabled: "disabled",
}

type UpstreamStatusDetail struct {
	// Provides a human-readable reason for the status of the upstream.
	Reason *string `json:"reason,omitempty"`
}
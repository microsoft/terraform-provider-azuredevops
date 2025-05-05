// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package dashboard

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
)

// Copy options of a Dashboard.
type CopyDashboardOptions struct {
	// Dashboard Scope. Can be either Project or Project_Team
	CopyDashboardScope *DashboardScope `json:"copyDashboardScope,omitempty"`
	// When this flag is set to true,option to select the folder to copy Queries of copy dashboard will appear.
	CopyQueriesFlag *bool `json:"copyQueriesFlag,omitempty"`
	// Description of the dashboard
	Description *string `json:"description,omitempty"`
	// Name of the dashboard
	Name *string `json:"name,omitempty"`
	// ID of the project. Provided by service at creation time.
	ProjectId *uuid.UUID `json:"projectId,omitempty"`
	// Path to which the queries should be copied of copy dashboard
	QueryFolderPath *uuid.UUID `json:"queryFolderPath,omitempty"`
	// Refresh interval of dashboard
	RefreshInterval *int `json:"refreshInterval,omitempty"`
	// ID of the team. Provided by service at creation time
	TeamId *uuid.UUID `json:"teamId,omitempty"`
}

type CopyDashboardResponse struct {
	// Copied Dashboard
	CopiedDashboard *Dashboard `json:"copiedDashboard,omitempty"`
	// Copy Dashboard options
	CopyDashboardOptions *CopyDashboardOptions `json:"copyDashboardOptions,omitempty"`
}

// Model of a Dashboard.
type Dashboard struct {
	Links interface{} `json:"_links,omitempty"`
	// Entity to which the dashboard is scoped.
	DashboardScope *DashboardScope `json:"dashboardScope,omitempty"`
	// Description of the dashboard.
	Description *string `json:"description,omitempty"`
	// Server defined version tracking value, used for edit collision detection.
	ETag *string `json:"eTag,omitempty"`
	// ID of the group for a dashboard. For team-scoped dashboards, this is the unique identifier for the team associated with the dashboard. For project-scoped dashboards this property is empty.
	GroupId *uuid.UUID `json:"groupId,omitempty"`
	// ID of the Dashboard. Provided by service at creation time.
	Id *uuid.UUID `json:"id,omitempty"`
	// Dashboard Last Accessed Date.
	LastAccessedDate *azuredevops.Time `json:"lastAccessedDate,omitempty"`
	// Id of the person who modified Dashboard.
	ModifiedBy *uuid.UUID `json:"modifiedBy,omitempty"`
	// Dashboard's last modified date.
	ModifiedDate *azuredevops.Time `json:"modifiedDate,omitempty"`
	// Name of the Dashboard.
	Name *string `json:"name,omitempty"`
	// ID of the owner for a dashboard. For team-scoped dashboards, this is the unique identifier for the team associated with the dashboard. For project-scoped dashboards, this is the unique identifier for the user identity associated with the dashboard.
	OwnerId *uuid.UUID `json:"ownerId,omitempty"`
	// Position of the dashboard, within a dashboard group. If unset at creation time, position is decided by the service.
	Position *int `json:"position,omitempty"`
	// Interval for client to automatically refresh the dashboard. Expressed in minutes.
	RefreshInterval *int    `json:"refreshInterval,omitempty"`
	Url             *string `json:"url,omitempty"`
	// The set of Widgets on the dashboard.
	Widgets *[]Widget `json:"widgets,omitempty"`
}

// Describes a list of dashboards associated to an owner. Currently, teams own dashboard groups.
type DashboardGroup struct {
	Links interface{} `json:"_links,omitempty"`
	// A list of Dashboards held by the Dashboard Group
	DashboardEntries *[]DashboardGroupEntry `json:"dashboardEntries,omitempty"`
	// Deprecated: The old permission model describing the level of permissions for the current team. Pre-M125.
	Permission *GroupMemberPermission `json:"permission,omitempty"`
	// A permissions bit mask describing the security permissions of the current team for dashboards. When this permission is the value None, use GroupMemberPermission. Permissions are evaluated based on the presence of a value other than None, else the GroupMemberPermission will be saved.
	TeamDashboardPermission *TeamDashboardPermission `json:"teamDashboardPermission,omitempty"`
	Url                     *string                  `json:"url,omitempty"`
}

// Dashboard group entry, wrapping around Dashboard (needed?)
type DashboardGroupEntry struct {
	Links interface{} `json:"_links,omitempty"`
	// Entity to which the dashboard is scoped.
	DashboardScope *DashboardScope `json:"dashboardScope,omitempty"`
	// Description of the dashboard.
	Description *string `json:"description,omitempty"`
	// Server defined version tracking value, used for edit collision detection.
	ETag *string `json:"eTag,omitempty"`
	// ID of the group for a dashboard. For team-scoped dashboards, this is the unique identifier for the team associated with the dashboard. For project-scoped dashboards this property is empty.
	GroupId *uuid.UUID `json:"groupId,omitempty"`
	// ID of the Dashboard. Provided by service at creation time.
	Id *uuid.UUID `json:"id,omitempty"`
	// Dashboard Last Accessed Date.
	LastAccessedDate *azuredevops.Time `json:"lastAccessedDate,omitempty"`
	// Id of the person who modified Dashboard.
	ModifiedBy *uuid.UUID `json:"modifiedBy,omitempty"`
	// Dashboard's last modified date.
	ModifiedDate *azuredevops.Time `json:"modifiedDate,omitempty"`
	// Name of the Dashboard.
	Name *string `json:"name,omitempty"`
	// ID of the owner for a dashboard. For team-scoped dashboards, this is the unique identifier for the team associated with the dashboard. For project-scoped dashboards, this is the unique identifier for the user identity associated with the dashboard.
	OwnerId *uuid.UUID `json:"ownerId,omitempty"`
	// Position of the dashboard, within a dashboard group. If unset at creation time, position is decided by the service.
	Position *int `json:"position,omitempty"`
	// Interval for client to automatically refresh the dashboard. Expressed in minutes.
	RefreshInterval *int    `json:"refreshInterval,omitempty"`
	Url             *string `json:"url,omitempty"`
	// The set of Widgets on the dashboard.
	Widgets *[]Widget `json:"widgets,omitempty"`
}

// Response from RestAPI when saving and editing DashboardGroupEntry
type DashboardGroupEntryResponse struct {
	Links interface{} `json:"_links,omitempty"`
	// Entity to which the dashboard is scoped.
	DashboardScope *DashboardScope `json:"dashboardScope,omitempty"`
	// Description of the dashboard.
	Description *string `json:"description,omitempty"`
	// Server defined version tracking value, used for edit collision detection.
	ETag *string `json:"eTag,omitempty"`
	// ID of the group for a dashboard. For team-scoped dashboards, this is the unique identifier for the team associated with the dashboard. For project-scoped dashboards this property is empty.
	GroupId *uuid.UUID `json:"groupId,omitempty"`
	// ID of the Dashboard. Provided by service at creation time.
	Id *uuid.UUID `json:"id,omitempty"`
	// Dashboard Last Accessed Date.
	LastAccessedDate *azuredevops.Time `json:"lastAccessedDate,omitempty"`
	// Id of the person who modified Dashboard.
	ModifiedBy *uuid.UUID `json:"modifiedBy,omitempty"`
	// Dashboard's last modified date.
	ModifiedDate *azuredevops.Time `json:"modifiedDate,omitempty"`
	// Name of the Dashboard.
	Name *string `json:"name,omitempty"`
	// ID of the owner for a dashboard. For team-scoped dashboards, this is the unique identifier for the team associated with the dashboard. For project-scoped dashboards, this is the unique identifier for the user identity associated with the dashboard.
	OwnerId *uuid.UUID `json:"ownerId,omitempty"`
	// Position of the dashboard, within a dashboard group. If unset at creation time, position is decided by the service.
	Position *int `json:"position,omitempty"`
	// Interval for client to automatically refresh the dashboard. Expressed in minutes.
	RefreshInterval *int    `json:"refreshInterval,omitempty"`
	Url             *string `json:"url,omitempty"`
	// The set of Widgets on the dashboard.
	Widgets *[]Widget `json:"widgets,omitempty"`
}

type DashboardResponse struct {
	Links interface{} `json:"_links,omitempty"`
	// Entity to which the dashboard is scoped.
	DashboardScope *DashboardScope `json:"dashboardScope,omitempty"`
	// Description of the dashboard.
	Description *string `json:"description,omitempty"`
	// Server defined version tracking value, used for edit collision detection.
	ETag *string `json:"eTag,omitempty"`
	// ID of the group for a dashboard. For team-scoped dashboards, this is the unique identifier for the team associated with the dashboard. For project-scoped dashboards this property is empty.
	GroupId *uuid.UUID `json:"groupId,omitempty"`
	// ID of the Dashboard. Provided by service at creation time.
	Id *uuid.UUID `json:"id,omitempty"`
	// Dashboard Last Accessed Date.
	LastAccessedDate *azuredevops.Time `json:"lastAccessedDate,omitempty"`
	// Id of the person who modified Dashboard.
	ModifiedBy *uuid.UUID `json:"modifiedBy,omitempty"`
	// Dashboard's last modified date.
	ModifiedDate *azuredevops.Time `json:"modifiedDate,omitempty"`
	// Name of the Dashboard.
	Name *string `json:"name,omitempty"`
	// ID of the owner for a dashboard. For team-scoped dashboards, this is the unique identifier for the team associated with the dashboard. For project-scoped dashboards, this is the unique identifier for the user identity associated with the dashboard.
	OwnerId *uuid.UUID `json:"ownerId,omitempty"`
	// Position of the dashboard, within a dashboard group. If unset at creation time, position is decided by the service.
	Position *int `json:"position,omitempty"`
	// Interval for client to automatically refresh the dashboard. Expressed in minutes.
	RefreshInterval *int    `json:"refreshInterval,omitempty"`
	Url             *string `json:"url,omitempty"`
	// The set of Widgets on the dashboard.
	Widgets *[]Widget `json:"widgets,omitempty"`
}

// identifies the scope of dashboard storage and permissions.
type DashboardScope string

type dashboardScopeValuesType struct {
	Collection_User DashboardScope
	Project_Team    DashboardScope
	Project         DashboardScope
}

var DashboardScopeValues = dashboardScopeValuesType{
	// [DEPRECATED] Dashboard is scoped to the collection user.
	Collection_User: "collection_User",
	// Dashboard is scoped to the team.
	Project_Team: "project_Team",
	// Dashboard is scoped to the project.
	Project: "project",
}

type GroupMemberPermission string

type groupMemberPermissionValuesType struct {
	None              GroupMemberPermission
	Edit              GroupMemberPermission
	Manage            GroupMemberPermission
	ManagePermissions GroupMemberPermission
}

var GroupMemberPermissionValues = groupMemberPermissionValuesType{
	None:              "none",
	Edit:              "edit",
	Manage:            "manage",
	ManagePermissions: "managePermissions",
}

// Lightbox configuration
type LightboxOptions struct {
	// Height of desired lightbox, in pixels
	Height *int `json:"height,omitempty"`
	// True to allow lightbox resizing, false to disallow lightbox resizing, defaults to false.
	Resizable *bool `json:"resizable,omitempty"`
	// Width of desired lightbox, in pixels
	Width *int `json:"width,omitempty"`
}

// versioning for an artifact as described at: http://semver.org/, of the form major.minor.patch.
type SemanticVersion struct {
	// Major version when you make incompatible API changes
	Major *int `json:"major,omitempty"`
	// Minor version when you add functionality in a backwards-compatible manner
	Minor *int `json:"minor,omitempty"`
	// Patch version when you make backwards-compatible bug fixes
	Patch *int `json:"patch,omitempty"`
}

// [Flags]
type TeamDashboardPermission string

type teamDashboardPermissionValuesType struct {
	None              TeamDashboardPermission
	Read              TeamDashboardPermission
	Create            TeamDashboardPermission
	Edit              TeamDashboardPermission
	Delete            TeamDashboardPermission
	ManagePermissions TeamDashboardPermission
}

var TeamDashboardPermissionValues = teamDashboardPermissionValuesType{
	None:              "none",
	Read:              "read",
	Create:            "create",
	Edit:              "edit",
	Delete:            "delete",
	ManagePermissions: "managePermissions",
}

// Widget data
type Widget struct {
	Links interface{} `json:"_links,omitempty"`
	// Refers to the allowed sizes for the widget. This gets populated when user wants to configure the widget
	AllowedSizes *[]WidgetSize `json:"allowedSizes,omitempty"`
	// Read-Only Property from Dashboard Service. Indicates if settings are blocked for the current user.
	AreSettingsBlockedForUser *bool `json:"areSettingsBlockedForUser,omitempty"`
	// Refers to unique identifier of a feature artifact. Used for pinning+unpinning a specific artifact.
	ArtifactId                          *string `json:"artifactId,omitempty"`
	ConfigurationContributionId         *string `json:"configurationContributionId,omitempty"`
	ConfigurationContributionRelativeId *string `json:"configurationContributionRelativeId,omitempty"`
	ContentUri                          *string `json:"contentUri,omitempty"`
	// The id of the underlying contribution defining the supplied Widget Configuration.
	ContributionId *string `json:"contributionId,omitempty"`
	// Optional partial dashboard content, to support exchanging dashboard-level version ETag for widget-level APIs
	Dashboard          *Dashboard       `json:"dashboard,omitempty"`
	ETag               *string          `json:"eTag,omitempty"`
	Id                 *uuid.UUID       `json:"id,omitempty"`
	IsEnabled          *bool            `json:"isEnabled,omitempty"`
	IsNameConfigurable *bool            `json:"isNameConfigurable,omitempty"`
	LightboxOptions    *LightboxOptions `json:"lightboxOptions,omitempty"`
	LoadingImageUrl    *string          `json:"loadingImageUrl,omitempty"`
	Name               *string          `json:"name,omitempty"`
	Position           *WidgetPosition  `json:"position,omitempty"`
	Settings           *string          `json:"settings,omitempty"`
	SettingsVersion    *SemanticVersion `json:"settingsVersion,omitempty"`
	Size               *WidgetSize      `json:"size,omitempty"`
	TypeId             *string          `json:"typeId,omitempty"`
	Url                *string          `json:"url,omitempty"`
}

// Contribution based information describing Dashboard Widgets.
type WidgetMetadata struct {
	// Sizes supported by the Widget.
	AllowedSizes *[]WidgetSize `json:"allowedSizes,omitempty"`
	// Opt-in boolean that indicates if the widget requires the Analytics Service to function. Widgets requiring the analytics service are hidden from the catalog if the Analytics Service is not available.
	AnalyticsServiceRequired *bool `json:"analyticsServiceRequired,omitempty"`
	// Resource for an icon in the widget catalog.
	CatalogIconUrl *string `json:"catalogIconUrl,omitempty"`
	// Opt-in URL string pointing at widget information. Defaults to extension marketplace URL if omitted
	CatalogInfoUrl *string `json:"catalogInfoUrl,omitempty"`
	// The id of the underlying contribution defining the supplied Widget custom configuration UI. Null if custom configuration UI is not available.
	ConfigurationContributionId *string `json:"configurationContributionId,omitempty"`
	// The relative id of the underlying contribution defining the supplied Widget custom configuration UI. Null if custom configuration UI is not available.
	ConfigurationContributionRelativeId *string `json:"configurationContributionRelativeId,omitempty"`
	// Indicates if the widget requires configuration before being added to dashboard.
	ConfigurationRequired *bool `json:"configurationRequired,omitempty"`
	// Uri for the widget content to be loaded from .
	ContentUri *string `json:"contentUri,omitempty"`
	// The id of the underlying contribution defining the supplied Widget.
	ContributionId *string `json:"contributionId,omitempty"`
	// Optional default settings to be copied into widget settings.
	DefaultSettings *string `json:"defaultSettings,omitempty"`
	// Summary information describing the widget.
	Description *string `json:"description,omitempty"`
	// Widgets can be disabled by the app store.  We'll need to gracefully handle for: - persistence (Allow) - Requests (Tag as disabled, and provide context)
	IsEnabled *bool `json:"isEnabled,omitempty"`
	// Opt-out boolean that indicates if the widget supports widget name/title configuration. Widgets ignoring the name should set it to false in the manifest.
	IsNameConfigurable *bool `json:"isNameConfigurable,omitempty"`
	// Opt-out boolean indicating if the widget is hidden from the catalog. Commonly, this is used to allow developers to disable creation of a deprecated widget. A widget must have a functional default state, or have a configuration experience, in order to be visible from the catalog.
	IsVisibleFromCatalog *bool `json:"isVisibleFromCatalog,omitempty"`
	// Keywords associated with this widget, non-filterable and invisible
	Keywords *[]string `json:"keywords,omitempty"`
	// Opt-in properties for customizing widget presentation in a "lightbox" dialog.
	LightboxOptions *LightboxOptions `json:"lightboxOptions,omitempty"`
	// Resource for a loading placeholder image on dashboard
	LoadingImageUrl *string `json:"loadingImageUrl,omitempty"`
	// User facing name of the widget type. Each widget must use a unique value here.
	Name *string `json:"name,omitempty"`
	// Publisher Name of this kind of widget.
	PublisherName *string `json:"publisherName,omitempty"`
	// Data contract required for the widget to function and to work in its container.
	SupportedScopes *[]WidgetScope `json:"supportedScopes,omitempty"`
	// Tags associated with this widget, visible on each widget and filterable.
	Tags *[]string `json:"tags,omitempty"`
	// Contribution target IDs
	Targets *[]string `json:"targets,omitempty"`
	// Deprecated: locally unique developer-facing id of this kind of widget. ContributionId provides a globally unique identifier for widget types.
	TypeId *string `json:"typeId,omitempty"`
}

type WidgetMetadataResponse struct {
	Uri            *string         `json:"uri,omitempty"`
	WidgetMetadata *WidgetMetadata `json:"widgetMetadata,omitempty"`
}

type WidgetPosition struct {
	Column *int `json:"column,omitempty"`
	Row    *int `json:"row,omitempty"`
}

// Response from RestAPI when saving and editing Widget
type WidgetResponse struct {
	Links interface{} `json:"_links,omitempty"`
	// Refers to the allowed sizes for the widget. This gets populated when user wants to configure the widget
	AllowedSizes *[]WidgetSize `json:"allowedSizes,omitempty"`
	// Read-Only Property from Dashboard Service. Indicates if settings are blocked for the current user.
	AreSettingsBlockedForUser *bool `json:"areSettingsBlockedForUser,omitempty"`
	// Refers to unique identifier of a feature artifact. Used for pinning+unpinning a specific artifact.
	ArtifactId                          *string `json:"artifactId,omitempty"`
	ConfigurationContributionId         *string `json:"configurationContributionId,omitempty"`
	ConfigurationContributionRelativeId *string `json:"configurationContributionRelativeId,omitempty"`
	ContentUri                          *string `json:"contentUri,omitempty"`
	// The id of the underlying contribution defining the supplied Widget Configuration.
	ContributionId *string `json:"contributionId,omitempty"`
	// Optional partial dashboard content, to support exchanging dashboard-level version ETag for widget-level APIs
	Dashboard          *Dashboard       `json:"dashboard,omitempty"`
	ETag               *string          `json:"eTag,omitempty"`
	Id                 *uuid.UUID       `json:"id,omitempty"`
	IsEnabled          *bool            `json:"isEnabled,omitempty"`
	IsNameConfigurable *bool            `json:"isNameConfigurable,omitempty"`
	LightboxOptions    *LightboxOptions `json:"lightboxOptions,omitempty"`
	LoadingImageUrl    *string          `json:"loadingImageUrl,omitempty"`
	Name               *string          `json:"name,omitempty"`
	Position           *WidgetPosition  `json:"position,omitempty"`
	Settings           *string          `json:"settings,omitempty"`
	SettingsVersion    *SemanticVersion `json:"settingsVersion,omitempty"`
	Size               *WidgetSize      `json:"size,omitempty"`
	TypeId             *string          `json:"typeId,omitempty"`
	Url                *string          `json:"url,omitempty"`
}

// data contract required for the widget to function in a webaccess area or page.
type WidgetScope string

type widgetScopeValuesType struct {
	Collection_User WidgetScope
	Project_Team    WidgetScope
}

var WidgetScopeValues = widgetScopeValuesType{
	Collection_User: "collection_User",
	Project_Team:    "project_Team",
}

type WidgetSize struct {
	// The Width of the widget, expressed in dashboard grid columns.
	ColumnSpan *int `json:"columnSpan,omitempty"`
	// The height of the widget, expressed in dashboard grid rows.
	RowSpan *int `json:"rowSpan,omitempty"`
}

// Wrapper class to support HTTP header generation using CreateResponse, ClientHeaderParameter and ClientResponseType in WidgetV2Controller
type WidgetsVersionedList struct {
	ETag    *[]string `json:"eTag,omitempty"`
	Widgets *[]Widget `json:"widgets,omitempty"`
}

type WidgetTypesResponse struct {
	Links       interface{}       `json:"_links,omitempty"`
	Uri         *string           `json:"uri,omitempty"`
	WidgetTypes *[]WidgetMetadata `json:"widgetTypes,omitempty"`
}

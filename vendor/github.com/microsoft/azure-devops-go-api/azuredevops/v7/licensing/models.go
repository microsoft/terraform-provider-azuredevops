// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package licensing

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/accounts"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/commerce"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
)

// License assigned to a user
type AccessLevel struct {
	// Type of Account License (e.g. Express, Stakeholder etc.)
	AccountLicenseType *AccountLicenseType `json:"accountLicenseType,omitempty"`
	// Assignment Source of the License (e.g. Group, Unknown etc.
	AssignmentSource *AssignmentSource `json:"assignmentSource,omitempty"`
	// Display name of the License
	LicenseDisplayName *string `json:"licenseDisplayName,omitempty"`
	// Licensing Source (e.g. Account. MSDN etc.)
	LicensingSource *LicensingSource `json:"licensingSource,omitempty"`
	// Type of MSDN License (e.g. Visual Studio Professional, Visual Studio Enterprise etc.)
	MsdnLicenseType *MsdnLicenseType `json:"msdnLicenseType,omitempty"`
	// User status in the account
	Status *accounts.AccountUserStatus `json:"status,omitempty"`
	// Status message.
	StatusMessage *string `json:"statusMessage,omitempty"`
}

// Represents a license granted to a user in an account
type AccountEntitlement struct {
	// Gets or sets the id of the account to which the license belongs
	AccountId *uuid.UUID `json:"accountId,omitempty"`
	// Gets or sets the date the license was assigned
	AssignmentDate *azuredevops.Time `json:"assignmentDate,omitempty"`
	// Assignment Source
	AssignmentSource *AssignmentSource `json:"assignmentSource,omitempty"`
	// Gets or sets the creation date of the user in this account
	DateCreated *azuredevops.Time `json:"dateCreated,omitempty"`
	// Gets or sets the date of the user last sign-in to this account
	LastAccessedDate *azuredevops.Time `json:"lastAccessedDate,omitempty"`
	License          *License          `json:"license,omitempty"`
	// Licensing origin
	Origin *LicensingOrigin `json:"origin,omitempty"`
	// The computed rights of this user in the account.
	Rights *AccountRights `json:"rights,omitempty"`
	// The status of the user in the account
	Status *accounts.AccountUserStatus `json:"status,omitempty"`
	// Identity information of the user to which the license belongs
	User *webapi.IdentityRef `json:"user,omitempty"`
	// Gets the id of the user to which the license belongs
	UserId *uuid.UUID `json:"userId,omitempty"`
}

// Model for updating an AccountEntitlement for a user, used for the Web API
type AccountEntitlementUpdateModel struct {
	// Gets or sets the license for the entitlement
	License *License `json:"license,omitempty"`
}

// Represents an Account license
type AccountLicense struct {
	// Gets the source of the license
	Source *LicensingSource `json:"source,omitempty"`
	// Gets the license type for the license
	License *AccountLicenseType `json:"license,omitempty"`
}

type AccountLicenseExtensionUsage struct {
	ExtensionId            *string                               `json:"extensionId,omitempty"`
	ExtensionName          *string                               `json:"extensionName,omitempty"`
	IncludedQuantity       *int                                  `json:"includedQuantity,omitempty"`
	IsTrial                *bool                                 `json:"isTrial,omitempty"`
	MinimumLicenseRequired *commerce.MinimumRequiredServiceLevel `json:"minimumLicenseRequired,omitempty"`
	MsdnUsedCount          *int                                  `json:"msdnUsedCount,omitempty"`
	ProvisionedCount       *int                                  `json:"provisionedCount,omitempty"`
	RemainingTrialDays     *int                                  `json:"remainingTrialDays,omitempty"`
	TrialExpiryDate        *azuredevops.Time                     `json:"trialExpiryDate,omitempty"`
	UsedCount              *int                                  `json:"usedCount,omitempty"`
}

type AccountLicenseType string

type accountLicenseTypeValuesType struct {
	None         AccountLicenseType
	EarlyAdopter AccountLicenseType
	Express      AccountLicenseType
	Professional AccountLicenseType
	Advanced     AccountLicenseType
	Stakeholder  AccountLicenseType
}

var AccountLicenseTypeValues = accountLicenseTypeValuesType{
	None:         "none",
	EarlyAdopter: "earlyAdopter",
	Express:      "express",
	Professional: "professional",
	Advanced:     "advanced",
	Stakeholder:  "stakeholder",
}

type AccountRights struct {
	Level  *VisualStudioOnlineServiceLevel `json:"level,omitempty"`
	Reason *string                         `json:"reason,omitempty"`
}

type AccountUserLicense struct {
	License *int             `json:"license,omitempty"`
	Source  *LicensingSource `json:"source,omitempty"`
}

type AssignmentSource string

type assignmentSourceValuesType struct {
	None      AssignmentSource
	Unknown   AssignmentSource
	GroupRule AssignmentSource
}

var AssignmentSourceValues = assignmentSourceValuesType{
	None:      "none",
	Unknown:   "unknown",
	GroupRule: "groupRule",
}

type AutoLicense struct {
	// Gets the source of the license
	Source *LicensingSource `json:"source,omitempty"`
}

type ClientRightsContainer struct {
	CertificateBytes *[]byte `json:"certificateBytes,omitempty"`
	Token            *string `json:"token,omitempty"`
}

// Model for assigning an extension to users, used for the Web API
type ExtensionAssignment struct {
	// Gets or sets the extension ID to assign.
	ExtensionGalleryId *string `json:"extensionGalleryId,omitempty"`
	// Set to true if this a auto assignment scenario.
	IsAutoAssignment *bool `json:"isAutoAssignment,omitempty"`
	// Gets or sets the licensing source.
	LicensingSource *LicensingSource `json:"licensingSource,omitempty"`
	// Gets or sets the user IDs to assign the extension to.
	UserIds *[]uuid.UUID `json:"userIds,omitempty"`
}

// Model for assigning an extension to users, used for the Web API
type ExtensionSource struct {
	// Assignment Source
	AssignmentSource *AssignmentSource `json:"assignmentSource,omitempty"`
	// extension Identifier
	ExtensionGalleryId *string `json:"extensionGalleryId,omitempty"`
	// The licensing source of the extension. Account, Msdn, etc.
	LicensingSource *LicensingSource `json:"licensingSource,omitempty"`
}

// The base class for a specific license source and license
type License struct {
	// Gets the source of the license
	Source *LicensingSource `json:"source,omitempty"`
}

type LicensingOrigin string

type licensingOriginValuesType struct {
	None                     LicensingOrigin
	OnDemandPrivateProject   LicensingOrigin
	OnDemandPublicProject    LicensingOrigin
	UserHubInvitation        LicensingOrigin
	PrivateProjectInvitation LicensingOrigin
	PublicProjectInvitation  LicensingOrigin
}

var LicensingOriginValues = licensingOriginValuesType{
	None:                     "none",
	OnDemandPrivateProject:   "onDemandPrivateProject",
	OnDemandPublicProject:    "onDemandPublicProject",
	UserHubInvitation:        "userHubInvitation",
	PrivateProjectInvitation: "privateProjectInvitation",
	PublicProjectInvitation:  "publicProjectInvitation",
}

// [Flags]
type LicensingSettingsSelectProperty string

type licensingSettingsSelectPropertyValuesType struct {
	DefaultAccessLevel LicensingSettingsSelectProperty
	AccessLevelOptions LicensingSettingsSelectProperty
	All                LicensingSettingsSelectProperty
}

var LicensingSettingsSelectPropertyValues = licensingSettingsSelectPropertyValuesType{
	DefaultAccessLevel: "defaultAccessLevel",
	AccessLevelOptions: "accessLevelOptions",
	All:                "all",
}

type LicensingSource string

type licensingSourceValuesType struct {
	None    LicensingSource
	Account LicensingSource
	Msdn    LicensingSource
	Profile LicensingSource
	Auto    LicensingSource
	Trial   LicensingSource
}

var LicensingSourceValues = licensingSourceValuesType{
	None:    "none",
	Account: "account",
	Msdn:    "msdn",
	Profile: "profile",
	Auto:    "auto",
	Trial:   "trial",
}

// Represents an Msdn license
type MsdnLicense struct {
	// Gets the source of the license
	Source *LicensingSource `json:"source,omitempty"`
	// Gets the license type for the license
	License *MsdnLicenseType `json:"license,omitempty"`
}

type MsdnLicenseType string

type msdnLicenseTypeValuesType struct {
	None             MsdnLicenseType
	Eligible         MsdnLicenseType
	Professional     MsdnLicenseType
	Platforms        MsdnLicenseType
	TestProfessional MsdnLicenseType
	Premium          MsdnLicenseType
	Ultimate         MsdnLicenseType
	Enterprise       MsdnLicenseType
}

var MsdnLicenseTypeValues = msdnLicenseTypeValuesType{
	None:             "none",
	Eligible:         "eligible",
	Professional:     "professional",
	Platforms:        "platforms",
	TestProfessional: "testProfessional",
	Premium:          "premium",
	Ultimate:         "ultimate",
	Enterprise:       "enterprise",
}

type NoLicense struct {
	// Gets the source of the license
	Source *LicensingSource `json:"source,omitempty"`
}

type VisualStudioOnlineServiceLevel string

type visualStudioOnlineServiceLevelValuesType struct {
	None         VisualStudioOnlineServiceLevel
	Express      VisualStudioOnlineServiceLevel
	Advanced     VisualStudioOnlineServiceLevel
	AdvancedPlus VisualStudioOnlineServiceLevel
	Stakeholder  VisualStudioOnlineServiceLevel
}

var VisualStudioOnlineServiceLevelValues = visualStudioOnlineServiceLevelValuesType{
	// No service rights. The user cannot access the account
	None: "none",
	// Default or minimum service level
	Express: "express",
	// Premium service level - either by purchasing on the Azure portal or by purchasing the appropriate MSDN subscription
	Advanced: "advanced",
	// Only available to a specific set of MSDN Subscribers
	AdvancedPlus: "advancedPlus",
	// Stakeholder service level
	Stakeholder: "stakeholder",
}

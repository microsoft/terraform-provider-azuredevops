// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package commerce

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
)

// The subscription account namespace. Denotes the 'category' of the account.
type AccountProviderNamespace string

type accountProviderNamespaceValuesType struct {
	VisualStudioOnline AccountProviderNamespace
	AppInsights        AccountProviderNamespace
	Marketplace        AccountProviderNamespace
	OnPremise          AccountProviderNamespace
}

var AccountProviderNamespaceValues = accountProviderNamespaceValuesType{
	VisualStudioOnline: "visualStudioOnline",
	AppInsights:        "appInsights",
	Marketplace:        "marketplace",
	OnPremise:          "onPremise",
}

// Encapsulates Azure specific plan structure, using a publisher defined publisher name, offer name, and plan name. These are all specified by the publisher and can vary from other meta data we store about the extension internally therefore need to be tracked separately for purposes of interacting with Azure.
type AzureOfferPlanDefinition struct {
	// Determines whether or not this plan is visible to all users
	IsPublic *bool `json:"isPublic,omitempty"`
	// The meter id which identifies the offer meter this plan is associated with
	MeterId *int `json:"meterId,omitempty"`
	// The offer / product name as defined by the publisher in Azure
	OfferId *string `json:"offerId,omitempty"`
	// The offer / product name as defined by the publisher in Azure
	OfferName *string `json:"offerName,omitempty"`
	// The id of the plan, which is usually in the format "{publisher}:{offer}:{plan}"
	PlanId *string `json:"planId,omitempty"`
	// The plan name as defined by the publisher in Azure
	PlanName *string `json:"planName,omitempty"`
	// The version string which optionally identifies the version of the plan
	PlanVersion *string `json:"planVersion,omitempty"`
	// The publisher of the plan as defined by the publisher in Azure
	Publisher *string `json:"publisher,omitempty"`
	// get/set publisher name
	PublisherName *string `json:"publisherName,omitempty"`
	// The number of users associated with the plan as defined in Azure
	Quantity *int `json:"quantity,omitempty"`
}

// These are known offer types to VSTS.
type AzureOfferType string

type azureOfferTypeValuesType struct {
	None        AzureOfferType
	Standard    AzureOfferType
	Ea          AzureOfferType
	Msdn        AzureOfferType
	Csp         AzureOfferType
	Unsupported AzureOfferType
}

var AzureOfferTypeValues = azureOfferTypeValuesType{
	None:        "none",
	Standard:    "standard",
	Ea:          "ea",
	Msdn:        "msdn",
	Csp:         "csp",
	Unsupported: "unsupported",
}

// Represents an azure region, used by ibiza for linking accounts
type AzureRegion struct {
	// Display Name of the azure region. Ex: North Central US.
	DisplayName *string `json:"displayName,omitempty"`
	// Unique Identifier
	Id *string `json:"id,omitempty"`
	// Region code of the azure region. Ex: NCUS.
	RegionCode *string `json:"regionCode,omitempty"`
}

// The responsible entity/method for billing.
type BillingProvider string

type billingProviderValuesType struct {
	SelfManaged       BillingProvider
	AzureStoreManaged BillingProvider
}

var BillingProviderValues = billingProviderValuesType{
	SelfManaged:       "selfManaged",
	AzureStoreManaged: "azureStoreManaged",
}

type ConnectedServer struct {
	// Hosted AccountId associated with the connected server NOTE: As of S112, this is now the CollectionId. Not changed as this is exposed to client code.
	AccountId *uuid.UUID `json:"accountId,omitempty"`
	// Hosted AccountName associated with the connected server NOTE: As of S112, this is now the collection name. Not changed as this is exposed to client code.
	AccountName *string `json:"accountName,omitempty"`
	// Object used to create credentials to call from OnPrem to hosted service.
	Authorization *ConnectedServerAuthorization `json:"authorization,omitempty"`
	// OnPrem server id associated with the connected server
	ServerId *uuid.UUID `json:"serverId,omitempty"`
	// OnPrem server associated with the connected server
	ServerName *string `json:"serverName,omitempty"`
	// SpsUrl of the hosted account that the onrepm server has been connected to.
	SpsUrl *string `json:"spsUrl,omitempty"`
	// The id of the subscription used for purchase
	SubscriptionId *uuid.UUID `json:"subscriptionId,omitempty"`
	// OnPrem target host associated with the connected server.  Typically the collection host id
	TargetId *uuid.UUID `json:"targetId,omitempty"`
	// OnPrem target associated with the connected server.
	TargetName *string `json:"targetName,omitempty"`
}

// Provides data necessary for authorizing the connecter server using OAuth 2.0 authentication flows.
type ConnectedServerAuthorization struct {
	// Gets or sets the endpoint used to obtain access tokens from the configured token service.
	AuthorizationUrl *string `json:"authorizationUrl,omitempty"`
	// Gets or sets the client identifier for this agent.
	ClientId *uuid.UUID `json:"clientId,omitempty"`
	// Gets or sets the public key used to verify the identity of this connected server.
	PublicKey *string `json:"publicKey,omitempty"`
}

type IAzureSubscription struct {
	AnniversaryDay *int                      `json:"anniversaryDay,omitempty"`
	Created        *azuredevops.Time         `json:"created,omitempty"`
	Id             *uuid.UUID                `json:"id,omitempty"`
	LastUpdated    *azuredevops.Time         `json:"lastUpdated,omitempty"`
	Namespace      *AccountProviderNamespace `json:"namespace,omitempty"`
	OfferType      *AzureOfferType           `json:"offerType,omitempty"`
	Source         *SubscriptionSource       `json:"source,omitempty"`
	Status         *SubscriptionStatus       `json:"status,omitempty"`
}

type ICommerceEvent struct {
	// Billed quantity (prorated) passed to Azure commerce
	BilledQuantity *float64   `json:"billedQuantity,omitempty"`
	CollectionId   *uuid.UUID `json:"collectionId,omitempty"`
	CollectionName *string    `json:"collectionName,omitempty"`
	// Quantity for current billing cycle
	CommittedQuantity *int `json:"committedQuantity,omitempty"`
	// Quantity for next billing cycle
	CurrentQuantity *int              `json:"currentQuantity,omitempty"`
	EffectiveDate   *azuredevops.Time `json:"effectiveDate,omitempty"`
	// Onpremise or hosted
	Environment              *string           `json:"environment,omitempty"`
	EventId                  *string           `json:"eventId,omitempty"`
	EventName                *string           `json:"eventName,omitempty"`
	EventSource              *string           `json:"eventSource,omitempty"`
	EventTime                *azuredevops.Time `json:"eventTime,omitempty"`
	GalleryId                *string           `json:"galleryId,omitempty"`
	IncludedQuantity         *int              `json:"includedQuantity,omitempty"`
	MaxQuantity              *int              `json:"maxQuantity,omitempty"`
	MeterName                *string           `json:"meterName,omitempty"`
	OrganizationId           *uuid.UUID        `json:"organizationId,omitempty"`
	OrganizationName         *string           `json:"organizationName,omitempty"`
	PreviousIncludedQuantity *int              `json:"previousIncludedQuantity,omitempty"`
	PreviousMaxQuantity      *int              `json:"previousMaxQuantity,omitempty"`
	// Previous quantity in case of upgrade/downgrade
	PreviousQuantity *int              `json:"previousQuantity,omitempty"`
	RenewalGroup     *string           `json:"renewalGroup,omitempty"`
	ServiceIdentity  *uuid.UUID        `json:"serviceIdentity,omitempty"`
	SubscriptionId   *uuid.UUID        `json:"subscriptionId,omitempty"`
	TrialEndDate     *azuredevops.Time `json:"trialEndDate,omitempty"`
	TrialStartDate   *azuredevops.Time `json:"trialStartDate,omitempty"`
	UserIdentity     *uuid.UUID        `json:"userIdentity,omitempty"`
	Version          *string           `json:"version,omitempty"`
}

// Encapsulates the state of offer meter definitions and purchases
type ICommercePackage struct {
	Configuration      *map[string]string   `json:"configuration,omitempty"`
	OfferMeters        *[]OfferMeter        `json:"offerMeters,omitempty"`
	OfferSubscriptions *[]OfferSubscription `json:"offerSubscriptions,omitempty"`
}

// Information about a resource associated with a subscription.
type IOfferSubscription struct {
	// Indicates whether users get auto assigned this license type duing first access.
	AutoAssignOnAccess *bool `json:"autoAssignOnAccess,omitempty"`
	// The azure subscription id
	AzureSubscriptionId *uuid.UUID `json:"azureSubscriptionId,omitempty"`
	// The azure subscription name
	AzureSubscriptionName *string `json:"azureSubscriptionName,omitempty"`
	// The azure subscription state
	AzureSubscriptionState *SubscriptionStatus `json:"azureSubscriptionState,omitempty"`
	// Quantity committed by the user, when resources is commitment based.
	CommittedQuantity *int `json:"committedQuantity,omitempty"`
	// A enumeration value indicating why the resource was disabled.
	DisabledReason *ResourceStatusReason `json:"disabledReason,omitempty"`
	// Uri pointing to user action on a disabled resource. It is based on DisabledReason value.
	DisabledResourceActionLink *string `json:"disabledResourceActionLink,omitempty"`
	// Quantity included for free.
	IncludedQuantity *int `json:"includedQuantity,omitempty"`
	// Returns true if paid billing is enabled on the resource. Returns false for non-azure subscriptions, disabled azure subscriptions or explicitly disabled by user
	IsPaidBillingEnabled *bool `json:"isPaidBillingEnabled,omitempty"`
	// Gets or sets a value indicating whether this instance is in preview.
	IsPreview *bool `json:"isPreview,omitempty"`
	// Gets the value indicating whether the puchase is canceled.
	IsPurchaseCanceled *bool `json:"isPurchaseCanceled,omitempty"`
	// Gets the value indicating whether current meter was purchased while the meter is still in trial
	IsPurchasedDuringTrial *bool `json:"isPurchasedDuringTrial,omitempty"`
	// Gets or sets a value indicating whether this instance is trial or preview.
	IsTrialOrPreview *bool `json:"isTrialOrPreview,omitempty"`
	// Returns true if resource is can be used otherwise returns false. DisabledReason can be used to identify why resource is disabled.
	IsUseable *bool `json:"isUseable,omitempty"`
	// Returns an integer representing the maximum quantity that can be billed for this resource. Any usage submitted over this number is automatically excluded from being sent to azure.
	MaximumQuantity *int `json:"maximumQuantity,omitempty"`
	// Gets the name of this resource.
	OfferMeter *OfferMeter `json:"offerMeter,omitempty"`
	// Gets the renewal group.
	RenewalGroup *ResourceRenewalGroup `json:"renewalGroup,omitempty"`
	// Returns a Date of UTC kind indicating when the next reset of quantities is going to happen. On this day at UTC 2:00 AM is when the reset will occur.
	ResetDate *azuredevops.Time `json:"resetDate,omitempty"`
	// Gets or sets the start date for this resource. First install date in any state.
	StartDate *azuredevops.Time `json:"startDate,omitempty"`
	// Gets or sets the trial expiry date.
	TrialExpiryDate *azuredevops.Time `json:"trialExpiryDate,omitempty"`
}

// The subscription account. Add Sub Type and Owner email later.
type ISubscriptionAccount struct {
	// Gets or sets the account host type.
	AccountHostType *int `json:"accountHostType,omitempty"`
	// Gets or sets the account identifier. Usually a guid.
	AccountId *uuid.UUID `json:"accountId,omitempty"`
	// Gets or sets the name of the account.
	AccountName *string `json:"accountName,omitempty"`
	// Gets or sets the account tenantId.
	AccountTenantId *uuid.UUID `json:"accountTenantId,omitempty"`
	// get or set purchase Error Reason
	FailedPurchaseReason *PurchaseErrorReason `json:"failedPurchaseReason,omitempty"`
	// Gets or sets the geo location.
	GeoLocation *string `json:"geoLocation,omitempty"`
	// Gets or sets a value indicating whether the calling user identity owns or is a PCA of the account.
	IsAccountOwner *bool `json:"isAccountOwner,omitempty"`
	// Gets or set the flag to enable purchase via subscription.
	IsEligibleForPurchase *bool `json:"isEligibleForPurchase,omitempty"`
	// get or set IsPrepaidFundSubscription
	IsPrepaidFundSubscription *bool `json:"isPrepaidFundSubscription,omitempty"`
	// get or set IsPricingPricingAvailable
	IsPricingAvailable *bool `json:"isPricingAvailable,omitempty"`
	// Gets or sets the subscription locale
	Locale *string `json:"locale,omitempty"`
	// Gets or sets the Offer Type of this subscription. A value of null means, this value has not been evaluated.
	OfferType *AzureOfferType `json:"offerType,omitempty"`
	// Gets or sets the subscription address country display name
	RegionDisplayName *string `json:"regionDisplayName,omitempty"`
	// Gets or sets the resource group.
	ResourceGroupName *string `json:"resourceGroupName,omitempty"`
	// Gets or sets the azure resource name.
	ResourceName *string `json:"resourceName,omitempty"`
	// A dictionary of service urls, mapping the service owner to the service owner url
	ServiceUrls *map[uuid.UUID]string `json:"serviceUrls,omitempty"`
	// Gets or sets the subscription identifier.
	SubscriptionId *uuid.UUID `json:"subscriptionId,omitempty"`
	// Gets or sets the azure subscription name
	SubscriptionName *string `json:"subscriptionName,omitempty"`
	// get or set object id of subscruption admin
	SubscriptionObjectId *uuid.UUID `json:"subscriptionObjectId,omitempty"`
	// get or set subscription offer code
	SubscriptionOfferCode *string `json:"subscriptionOfferCode,omitempty"`
	// Gets or sets the subscription status.
	SubscriptionStatus *SubscriptionStatus `json:"subscriptionStatus,omitempty"`
	// get or set tenant id of subscription
	SubscriptionTenantId *uuid.UUID `json:"subscriptionTenantId,omitempty"`
}

// Information about a resource associated with a subscription.
type ISubscriptionResource struct {
	// Quantity committed by the user, when resources is commitment based.
	CommittedQuantity *int `json:"committedQuantity,omitempty"`
	// A enumeration value indicating why the resource was disabled.
	DisabledReason *ResourceStatusReason `json:"disabledReason,omitempty"`
	// Uri pointing to user action on a disabled resource. It is based on DisabledReason value.
	DisabledResourceActionLink *string `json:"disabledResourceActionLink,omitempty"`
	// Quantity included for free.
	IncludedQuantity *int `json:"includedQuantity,omitempty"`
	// Returns true if paid billing is enabled on the resource. Returns false for non-azure subscriptions, disabled azure subscriptions or explicitly disabled by user
	IsPaidBillingEnabled *bool `json:"isPaidBillingEnabled,omitempty"`
	// Returns true if resource is can be used otherwise returns false. DisabledReason can be used to identify why resource is disabled.
	IsUseable *bool `json:"isUseable,omitempty"`
	// Returns an integer representing the maximum quantity that can be billed for this resource. Any usage submitted over this number is automatically excluded from being sent to azure.
	MaximumQuantity *int `json:"maximumQuantity,omitempty"`
	// Gets the name of this resource.
	Name *ResourceName `json:"name,omitempty"`
	// Returns a Date of UTC kind indicating when the next reset of quantities is going to happen. On this day at UTC 2:00 AM is when the reset will occur.
	ResetDate *azuredevops.Time `json:"resetDate,omitempty"`
}

// Represents the aggregated usage of a resource over a time span
type IUsageEventAggregate struct {
	// Gets or sets end time of the aggregated value, exclusive
	EndTime *azuredevops.Time `json:"endTime,omitempty"`
	// Gets or sets resource that the aggregated value represents
	Resource *ResourceName `json:"resource,omitempty"`
	// Gets or sets start time of the aggregated value, inclusive
	StartTime *azuredevops.Time `json:"startTime,omitempty"`
	// Gets or sets quantity of the resource used from start time to end time
	Value *int `json:"value,omitempty"`
}

// The meter billing state.
type MeterBillingState string

type meterBillingStateValuesType struct {
	Free MeterBillingState
	Paid MeterBillingState
}

var MeterBillingStateValues = meterBillingStateValuesType{
	Free: "free",
	Paid: "paid",
}

// Defines meter categories.
type MeterCategory string

type meterCategoryValuesType struct {
	Legacy    MeterCategory
	Bundle    MeterCategory
	Extension MeterCategory
}

var MeterCategoryValues = meterCategoryValuesType{
	Legacy:    "legacy",
	Bundle:    "bundle",
	Extension: "extension",
}

// Describes the Renewal frequncy of a Meter.
type MeterRenewalFrequecy string

type meterRenewalFrequecyValuesType struct {
	None     MeterRenewalFrequecy
	Monthly  MeterRenewalFrequecy
	Annually MeterRenewalFrequecy
}

var MeterRenewalFrequecyValues = meterRenewalFrequecyValuesType{
	None:     "none",
	Monthly:  "monthly",
	Annually: "annually",
}

// The meter state.
type MeterState string

type meterStateValuesType struct {
	Registered MeterState
	Active     MeterState
	Retired    MeterState
	Deleted    MeterState
}

var MeterStateValues = meterStateValuesType{
	Registered: "registered",
	Active:     "active",
	Retired:    "retired",
	Deleted:    "deleted",
}

type MinimumRequiredServiceLevel string

type minimumRequiredServiceLevelValuesType struct {
	None         MinimumRequiredServiceLevel
	Express      MinimumRequiredServiceLevel
	Advanced     MinimumRequiredServiceLevel
	AdvancedPlus MinimumRequiredServiceLevel
	Stakeholder  MinimumRequiredServiceLevel
}

var MinimumRequiredServiceLevelValues = minimumRequiredServiceLevelValuesType{
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

type OfferMeter struct {
	// Gets or sets the value of absolute maximum quantity for the resource
	AbsoluteMaximumQuantity *int `json:"absoluteMaximumQuantity,omitempty"`
	// Gets or sets the user assignment model.
	AssignmentModel *OfferMeterAssignmentModel `json:"assignmentModel,omitempty"`
	// Indicates whether users get auto assigned this license type duing first access.
	AutoAssignOnAccess *bool `json:"autoAssignOnAccess,omitempty"`
	// Gets or sets the responsible entity/method for billing. Determines how this meter is handled in the backend.
	BillingEntity *BillingProvider `json:"billingEntity,omitempty"`
	// Gets or sets the billing mode of the resource
	BillingMode *ResourceBillingMode `json:"billingMode,omitempty"`
	// Gets or sets the billing start date. If TrialDays + PreviewGraceDays > then, on 'BillingStartDate' it starts the preview Grace and/or trial period.
	BillingStartDate *azuredevops.Time `json:"billingStartDate,omitempty"`
	// Gets or sets the state of the billing.
	BillingState *MeterBillingState `json:"billingState,omitempty"`
	// Category.
	Category *MeterCategory `json:"category,omitempty"`
	// Quantity committed by the user, when resources is commitment based.
	CommittedQuantity *int `json:"committedQuantity,omitempty"`
	// Quantity used by the user, when resources is pay as you go or commitment based.
	CurrentQuantity *int `json:"currentQuantity,omitempty"`
	// Gets or sets the map of named quantity varied plans, plans can be purchased that vary only in the number of users included. Null if this offer meter does not support named fixed quantity plans.
	FixedQuantityPlans *[]AzureOfferPlanDefinition `json:"fixedQuantityPlans,omitempty"`
	// Gets or sets Gallery Id.
	GalleryId *string `json:"galleryId,omitempty"`
	// Gets or sets the Min license level the offer is free for.
	IncludedInLicenseLevel *MinimumRequiredServiceLevel `json:"includedInLicenseLevel,omitempty"`
	// Quantity included for free.
	IncludedQuantity *int `json:"includedQuantity,omitempty"`
	// Flag to identify whether the meter is First Party or Third Party based on BillingEntity If the BillingEntity is SelfManaged, the Meter is First Party otherwise its a Third Party Meter
	IsFirstParty *bool `json:"isFirstParty,omitempty"`
	// Gets or sets the value of maximum quantity for the resource
	MaximumQuantity *int `json:"maximumQuantity,omitempty"`
	// Meter Id.
	MeterId *int `json:"meterId,omitempty"`
	// Gets or sets the minimum required access level for the meter.
	MinimumRequiredAccessLevel *MinimumRequiredServiceLevel `json:"minimumRequiredAccessLevel,omitempty"`
	// Name of the resource
	Name *string `json:"name,omitempty"`
	// Gets or sets the offer scope.
	OfferScope *OfferScope `json:"offerScope,omitempty"`
	// Gets or sets the identifier representing this meter in commerce platform
	PlatformMeterId *uuid.UUID `json:"platformMeterId,omitempty"`
	// Gets or sets the preview grace days.
	PreviewGraceDays *byte `json:"previewGraceDays,omitempty"`
	// Gets or sets the Renewak Frequency.
	RenewalFrequency *MeterRenewalFrequecy `json:"renewalFrequency,omitempty"`
	// Gets or sets the status.
	Status *MeterState `json:"status,omitempty"`
	// Gets or sets the trial cycles.
	TrialCycles *int `json:"trialCycles,omitempty"`
	// Gets or sets the trial days.
	TrialDays *byte `json:"trialDays,omitempty"`
	// Measuring unit for this meter.
	Unit *string `json:"unit,omitempty"`
}

// The offer meter assignment model.
type OfferMeterAssignmentModel string

type offerMeterAssignmentModelValuesType struct {
	Explicit OfferMeterAssignmentModel
	Implicit OfferMeterAssignmentModel
}

var OfferMeterAssignmentModelValues = offerMeterAssignmentModelValuesType{
	// Users need to be explicitly assigned.
	Explicit: "explicit",
	// Users will be added automatically. All-or-nothing model.
	Implicit: "implicit",
}

type OfferMeterPrice struct {
	// Currency code
	CurrencyCode *string `json:"currencyCode,omitempty"`
	// The meter Name which identifies the offer meter this plan is associated with
	MeterName *string `json:"meterName,omitempty"`
	// The Name of the plan, which is usually in the format "{publisher}:{offer}:{plan}"
	PlanName *string `json:"planName,omitempty"`
	// Plan Price
	Price *float64 `json:"price,omitempty"`
	// Plan Quantity
	Quantity *float64 `json:"quantity,omitempty"`
	// Region price is for
	Region *string `json:"region,omitempty"`
}

// The offer scope.
type OfferScope string

type offerScopeValuesType struct {
	Account     OfferScope
	User        OfferScope
	UserAccount OfferScope
}

var OfferScopeValues = offerScopeValuesType{
	Account:     "account",
	User:        "user",
	UserAccount: "userAccount",
}

// Information about a resource associated with a subscription.
type OfferSubscription struct {
	// Indicates whether users get auto assigned this license type duing first access.
	AutoAssignOnAccess *bool `json:"autoAssignOnAccess,omitempty"`
	// The azure subscription id
	AzureSubscriptionId *uuid.UUID `json:"azureSubscriptionId,omitempty"`
	// The azure subscription name
	AzureSubscriptionName *string `json:"azureSubscriptionName,omitempty"`
	// The azure subscription state
	AzureSubscriptionState *SubscriptionStatus `json:"azureSubscriptionState,omitempty"`
	// Quantity committed by the user, when resources is commitment based.
	CommittedQuantity *int `json:"committedQuantity,omitempty"`
	// A enumeration value indicating why the resource was disabled.
	DisabledReason *ResourceStatusReason `json:"disabledReason,omitempty"`
	// Uri pointing to user action on a disabled resource. It is based on DisabledReason value.
	DisabledResourceActionLink *string `json:"disabledResourceActionLink,omitempty"`
	// Quantity included for free.
	IncludedQuantity *int `json:"includedQuantity,omitempty"`
	// Returns true if paid billing is enabled on the resource. Returns false for non-azure subscriptions, disabled azure subscriptions or explicitly disabled by user
	IsPaidBillingEnabled *bool `json:"isPaidBillingEnabled,omitempty"`
	// Gets or sets a value indicating whether this instance is in preview.
	IsPreview *bool `json:"isPreview,omitempty"`
	// Gets the value indicating whether the puchase is canceled.
	IsPurchaseCanceled *bool `json:"isPurchaseCanceled,omitempty"`
	// Gets the value indicating whether current meter was purchased while the meter is still in trial
	IsPurchasedDuringTrial *bool `json:"isPurchasedDuringTrial,omitempty"`
	// Gets or sets a value indicating whether this instance is trial or preview.
	IsTrialOrPreview *bool `json:"isTrialOrPreview,omitempty"`
	// Returns true if resource is can be used otherwise returns false. DisabledReason can be used to identify why resource is disabled.
	IsUseable *bool `json:"isUseable,omitempty"`
	// Returns an integer representing the maximum quantity that can be billed for this resource. Any usage submitted over this number is automatically excluded from being sent to azure.
	MaximumQuantity *int `json:"maximumQuantity,omitempty"`
	// Gets or sets the name of this resource.
	OfferMeter *OfferMeter `json:"offerMeter,omitempty"`
	// The unique identifier of this offer subscription
	OfferSubscriptionId *uuid.UUID `json:"offerSubscriptionId,omitempty"`
	// Gets the renewal group.
	RenewalGroup *ResourceRenewalGroup `json:"renewalGroup,omitempty"`
	// Returns a Date of UTC kind indicating when the next reset of quantities is going to happen. On this day at UTC 2:00 AM is when the reset will occur.
	ResetDate *azuredevops.Time `json:"resetDate,omitempty"`
	// Gets or sets the start date for this resource. First install date in any state.
	StartDate *azuredevops.Time `json:"startDate,omitempty"`
	// Gets or sets the trial expiry date.
	TrialExpiryDate *azuredevops.Time `json:"trialExpiryDate,omitempty"`
}

// The Purchasable offer meter.
type PurchasableOfferMeter struct {
	// Currency code for meter pricing
	CurrencyCode *string `json:"currencyCode,omitempty"`
	// Gets or sets the estimated renewal date.
	EstimatedRenewalDate *azuredevops.Time `json:"estimatedRenewalDate,omitempty"`
	// Locale for azure subscription
	LocaleCode *string `json:"localeCode,omitempty"`
	// Gets or sets the meter pricing (GraduatedPrice)
	MeterPricing *[]azuredevops.KeyValuePair `json:"meterPricing,omitempty"`
	// Gets or sets the offer meter definition.
	OfferMeterDefinition *OfferMeter `json:"offerMeterDefinition,omitempty"`
}

type PurchaseErrorReason string

type purchaseErrorReasonValuesType struct {
	None                           PurchaseErrorReason
	MonetaryLimitSet               PurchaseErrorReason
	InvalidOfferCode               PurchaseErrorReason
	NotAdminOrCoAdmin              PurchaseErrorReason
	InvalidRegionPurchase          PurchaseErrorReason
	PaymentInstrumentNotCreditCard PurchaseErrorReason
	InvalidOfferRegion             PurchaseErrorReason
	UnsupportedSubscription        PurchaseErrorReason
	DisabledSubscription           PurchaseErrorReason
	InvalidUser                    PurchaseErrorReason
	NotSubscriptionUser            PurchaseErrorReason
	UnsupportedSubscriptionCsp     PurchaseErrorReason
	TemporarySpendingLimit         PurchaseErrorReason
	AzureServiceError              PurchaseErrorReason
}

var PurchaseErrorReasonValues = purchaseErrorReasonValuesType{
	None:                           "none",
	MonetaryLimitSet:               "monetaryLimitSet",
	InvalidOfferCode:               "invalidOfferCode",
	NotAdminOrCoAdmin:              "notAdminOrCoAdmin",
	InvalidRegionPurchase:          "invalidRegionPurchase",
	PaymentInstrumentNotCreditCard: "paymentInstrumentNotCreditCard",
	InvalidOfferRegion:             "invalidOfferRegion",
	UnsupportedSubscription:        "unsupportedSubscription",
	DisabledSubscription:           "disabledSubscription",
	InvalidUser:                    "invalidUser",
	NotSubscriptionUser:            "notSubscriptionUser",
	UnsupportedSubscriptionCsp:     "unsupportedSubscriptionCsp",
	TemporarySpendingLimit:         "temporarySpendingLimit",
	AzureServiceError:              "azureServiceError",
}

// Represents a purchase request for requesting purchase by a user who does not have authorization to purchase.
type PurchaseRequest struct {
	// Name of the offer meter
	OfferMeterName *string `json:"offerMeterName,omitempty"`
	// Quantity for purchase
	Quantity *int `json:"quantity,omitempty"`
	// Reason for the purchase request
	Reason *string `json:"reason,omitempty"`
	// Response for this purchase request by the approver
	Response *PurchaseRequestResponse `json:"response,omitempty"`
}

// Type of purchase request response
type PurchaseRequestResponse string

type purchaseRequestResponseValuesType struct {
	None     PurchaseRequestResponse
	Approved PurchaseRequestResponse
	Denied   PurchaseRequestResponse
}

var PurchaseRequestResponseValues = purchaseRequestResponseValuesType{
	None:     "none",
	Approved: "approved",
	Denied:   "denied",
}

// The resource billing mode.
type ResourceBillingMode string

type resourceBillingModeValuesType struct {
	Committment ResourceBillingMode
	PayAsYouGo  ResourceBillingMode
}

var ResourceBillingModeValues = resourceBillingModeValuesType{
	Committment: "committment",
	PayAsYouGo:  "payAsYouGo",
}

// Various metered resources in VSTS
type ResourceName string

type resourceNameValuesType struct {
	StandardLicense             ResourceName
	AdvancedLicense             ResourceName
	ProfessionalLicense         ResourceName
	Build                       ResourceName
	LoadTest                    ResourceName
	PremiumBuildAgent           ResourceName
	PrivateOtherBuildAgent      ResourceName
	PrivateAzureBuildAgent      ResourceName
	Artifacts                   ResourceName
	MsHostedCICDforMacOS        ResourceName
	MsHostedCICDforWindowsLinux ResourceName
}

var ResourceNameValues = resourceNameValuesType{
	StandardLicense:             "standardLicense",
	AdvancedLicense:             "advancedLicense",
	ProfessionalLicense:         "professionalLicense",
	Build:                       "build",
	LoadTest:                    "loadTest",
	PremiumBuildAgent:           "premiumBuildAgent",
	PrivateOtherBuildAgent:      "privateOtherBuildAgent",
	PrivateAzureBuildAgent:      "privateAzureBuildAgent",
	Artifacts:                   "artifacts",
	MsHostedCICDforMacOS:        "msHostedCICDforMacOS",
	MsHostedCICDforWindowsLinux: "msHostedCICDforWindowsLinux",
}

// The resource renewal group.
type ResourceRenewalGroup string

type resourceRenewalGroupValuesType struct {
	Monthly ResourceRenewalGroup
	Jan     ResourceRenewalGroup
	Feb     ResourceRenewalGroup
	Mar     ResourceRenewalGroup
	Apr     ResourceRenewalGroup
	May     ResourceRenewalGroup
	Jun     ResourceRenewalGroup
	Jul     ResourceRenewalGroup
	Aug     ResourceRenewalGroup
	Sep     ResourceRenewalGroup
	Oct     ResourceRenewalGroup
	Nov     ResourceRenewalGroup
	Dec     ResourceRenewalGroup
}

var ResourceRenewalGroupValues = resourceRenewalGroupValuesType{
	Monthly: "monthly",
	Jan:     "jan",
	Feb:     "feb",
	Mar:     "mar",
	Apr:     "apr",
	May:     "may",
	Jun:     "jun",
	Jul:     "jul",
	Aug:     "aug",
	Sep:     "sep",
	Oct:     "oct",
	Nov:     "nov",
	Dec:     "dec",
}

// [Flags] Reason for disabled resource.
type ResourceStatusReason string

type resourceStatusReasonValuesType struct {
	None                   ResourceStatusReason
	NoAzureSubscription    ResourceStatusReason
	NoIncludedQuantityLeft ResourceStatusReason
	SubscriptionDisabled   ResourceStatusReason
	PaidBillingDisabled    ResourceStatusReason
	MaximumQuantityReached ResourceStatusReason
}

var ResourceStatusReasonValues = resourceStatusReasonValuesType{
	None:                   "none",
	NoAzureSubscription:    "noAzureSubscription",
	NoIncludedQuantityLeft: "noIncludedQuantityLeft",
	SubscriptionDisabled:   "subscriptionDisabled",
	PaidBillingDisabled:    "paidBillingDisabled",
	MaximumQuantityReached: "maximumQuantityReached",
}

// The subscription account. Add Sub Type and Owner email later.
type SubscriptionAccount struct {
	// Gets or sets the account host type.
	AccountHostType *int `json:"accountHostType,omitempty"`
	// Gets or sets the account identifier. Usually a guid.
	AccountId *uuid.UUID `json:"accountId,omitempty"`
	// Gets or sets the name of the account.
	AccountName *string `json:"accountName,omitempty"`
	// Gets or sets the account tenantId.
	AccountTenantId *uuid.UUID `json:"accountTenantId,omitempty"`
	// Purchase Error Reason
	FailedPurchaseReason *PurchaseErrorReason `json:"failedPurchaseReason,omitempty"`
	// Gets or sets the geo location.
	GeoLocation *string `json:"geoLocation,omitempty"`
	// Gets or sets a value indicating whether the calling user identity owns or is a PCA of the account.
	IsAccountOwner *bool `json:"isAccountOwner,omitempty"`
	// Gets or set the flag to enable purchase via subscription.
	IsEligibleForPurchase *bool `json:"isEligibleForPurchase,omitempty"`
	// get or set IsPrepaidFundSubscription
	IsPrepaidFundSubscription *bool `json:"isPrepaidFundSubscription,omitempty"`
	// get or set IsPricingPricingAvailable
	IsPricingAvailable *bool `json:"isPricingAvailable,omitempty"`
	// Gets or sets the subscription address country code
	Locale *string `json:"locale,omitempty"`
	// Gets or sets the Offer Type of this subscription.
	OfferType *AzureOfferType `json:"offerType,omitempty"`
	// Gets or sets the subscription address country display name
	RegionDisplayName *string `json:"regionDisplayName,omitempty"`
	// Gets or sets the resource group.
	ResourceGroupName *string `json:"resourceGroupName,omitempty"`
	// Gets or sets the azure resource name.
	ResourceName *string `json:"resourceName,omitempty"`
	// A dictionary of service urls, mapping the service owner to the service owner url
	ServiceUrls *map[uuid.UUID]string `json:"serviceUrls,omitempty"`
	// Gets or sets the subscription identifier.
	SubscriptionId *uuid.UUID `json:"subscriptionId,omitempty"`
	// Gets or sets the azure subscription name
	SubscriptionName *string `json:"subscriptionName,omitempty"`
	// object id of subscription admin
	SubscriptionObjectId *uuid.UUID `json:"subscriptionObjectId,omitempty"`
	// get or set subscription offer code
	SubscriptionOfferCode *string `json:"subscriptionOfferCode,omitempty"`
	// Gets or sets the subscription status.
	SubscriptionStatus *SubscriptionStatus `json:"subscriptionStatus,omitempty"`
	// tenant id of subscription
	SubscriptionTenantId *uuid.UUID `json:"subscriptionTenantId,omitempty"`
}

// Information about a resource associated with a subscription.
type SubscriptionResource struct {
	// Quantity committed by the user, when resources is commitment based.
	CommittedQuantity *int `json:"committedQuantity,omitempty"`
	// A enumeration value indicating why the resource was disabled.
	DisabledReason *ResourceStatusReason `json:"disabledReason,omitempty"`
	// Uri pointing to user action on a disabled resource. It is based on DisabledReason value.
	DisabledResourceActionLink *string `json:"disabledResourceActionLink,omitempty"`
	// Quantity included for free.
	IncludedQuantity *int `json:"includedQuantity,omitempty"`
	// Returns true if paid billing is enabled on the resource. Returns false for non-azure subscriptions, disabled azure subscriptions or explicitly disabled by user
	IsPaidBillingEnabled *bool `json:"isPaidBillingEnabled,omitempty"`
	// Returns true if resource is can be used otherwise returns false. DisabledReason can be used to identify why resource is disabled.
	IsUseable *bool `json:"isUseable,omitempty"`
	// Returns an integer representing the maximum quantity that can be billed for this resource. Any usage submitted over this number is automatically excluded from being sent to azure.
	MaximumQuantity *int `json:"maximumQuantity,omitempty"`
	// Gets or sets the name of this resource.
	Name *ResourceName `json:"name,omitempty"`
	// Returns a Date of UTC kind indicating when the next reset of quantities is going to happen. On this day at UTC 2:00 AM is when the reset will occur.
	ResetDate *azuredevops.Time `json:"resetDate,omitempty"`
}

type SubscriptionSource string

type subscriptionSourceValuesType struct {
	Normal              SubscriptionSource
	EnterpriseAgreement SubscriptionSource
	Internal            SubscriptionSource
	Unknown             SubscriptionSource
	FreeTier            SubscriptionSource
}

var SubscriptionSourceValues = subscriptionSourceValuesType{
	Normal:              "normal",
	EnterpriseAgreement: "enterpriseAgreement",
	Internal:            "internal",
	Unknown:             "unknown",
	FreeTier:            "freeTier",
}

// Azure subscription status
type SubscriptionStatus string

type subscriptionStatusValuesType struct {
	Unknown      SubscriptionStatus
	Active       SubscriptionStatus
	Disabled     SubscriptionStatus
	Deleted      SubscriptionStatus
	Unregistered SubscriptionStatus
}

var SubscriptionStatusValues = subscriptionStatusValuesType{
	Unknown:      "unknown",
	Active:       "active",
	Disabled:     "disabled",
	Deleted:      "deleted",
	Unregistered: "unregistered",
}

// Class that represents common set of properties for a raw usage event reported by TFS services.
type UsageEvent struct {
	// Gets or sets account id of the event. Note: This is for backward compat with BI.
	AccountId *uuid.UUID `json:"accountId,omitempty"`
	// Account name associated with the usage event
	AccountName *string `json:"accountName,omitempty"`
	// User GUID associated with the usage event
	AssociatedUser *uuid.UUID `json:"associatedUser,omitempty"`
	// Timestamp when this billing event is billable
	BillableDate *azuredevops.Time `json:"billableDate,omitempty"`
	// Unique event identifier
	EventId *string `json:"eventId,omitempty"`
	// Receiving Timestamp of the billing event by metering service
	EventTimestamp *azuredevops.Time `json:"eventTimestamp,omitempty"`
	// Gets or sets the event unique identifier.
	EventUniqueId *uuid.UUID `json:"eventUniqueId,omitempty"`
	// Meter Id.
	MeterName *string `json:"meterName,omitempty"`
	// Partition id of the account
	PartitionId *int `json:"partitionId,omitempty"`
	// Quantity of the usage event
	Quantity *int `json:"quantity,omitempty"`
	// Gets or sets the billing mode for the resource involved in the usage
	ResourceBillingMode *ResourceBillingMode `json:"resourceBillingMode,omitempty"`
	// Service context GUID associated with the usage event
	ServiceIdentity *uuid.UUID `json:"serviceIdentity,omitempty"`
	// Gets or sets subscription anniversary day of the subscription
	SubscriptionAnniversaryDay *int `json:"subscriptionAnniversaryDay,omitempty"`
	// Gets or sets subscription guid of the associated account of the event
	SubscriptionId *uuid.UUID `json:"subscriptionId,omitempty"`
}

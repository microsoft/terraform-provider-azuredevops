// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package accounts

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
)

type Account struct {
	// Identifier for an Account
	AccountId *uuid.UUID `json:"accountId,omitempty"`
	// Name for an account
	AccountName *string `json:"accountName,omitempty"`
	// Owner of account
	AccountOwner *uuid.UUID `json:"accountOwner,omitempty"`
	// Current account status
	AccountStatus *AccountStatus `json:"accountStatus,omitempty"`
	// Type of account: Personal, Organization
	AccountType *AccountType `json:"accountType,omitempty"`
	// Uri for an account
	AccountUri *string `json:"accountUri,omitempty"`
	// Who created the account
	CreatedBy *uuid.UUID `json:"createdBy,omitempty"`
	// Date account was created
	CreatedDate *azuredevops.Time `json:"createdDate,omitempty"`
	HasMoved    *bool             `json:"hasMoved,omitempty"`
	// Identity of last person to update the account
	LastUpdatedBy *uuid.UUID `json:"lastUpdatedBy,omitempty"`
	// Date account was last updated
	LastUpdatedDate *azuredevops.Time `json:"lastUpdatedDate,omitempty"`
	// Namespace for an account
	NamespaceId     *uuid.UUID `json:"namespaceId,omitempty"`
	NewCollectionId *uuid.UUID `json:"newCollectionId,omitempty"`
	// Organization that created the account
	OrganizationName *string `json:"organizationName,omitempty"`
	// Extended properties
	Properties interface{} `json:"properties,omitempty"`
	// Reason for current status
	StatusReason *string `json:"statusReason,omitempty"`
}

type AccountCreateInfoInternal struct {
	AccountName        *string                     `json:"accountName,omitempty"`
	Creator            *uuid.UUID                  `json:"creator,omitempty"`
	Organization       *string                     `json:"organization,omitempty"`
	Preferences        *AccountPreferencesInternal `json:"preferences,omitempty"`
	Properties         interface{}                 `json:"properties,omitempty"`
	ServiceDefinitions *[]azuredevops.KeyValuePair `json:"serviceDefinitions,omitempty"`
}

type AccountPreferencesInternal struct {
	Culture  interface{} `json:"culture,omitempty"`
	Language interface{} `json:"language,omitempty"`
	TimeZone interface{} `json:"timeZone,omitempty"`
}

type AccountStatus string

type accountStatusValuesType struct {
	None     AccountStatus
	Enabled  AccountStatus
	Disabled AccountStatus
	Deleted  AccountStatus
	Moved    AccountStatus
}

var AccountStatusValues = accountStatusValuesType{
	None: "none",
	// This hosting account is active and assigned to a customer.
	Enabled: "enabled",
	// This hosting account is disabled.
	Disabled: "disabled",
	// This account is part of deletion batch and scheduled for deletion.
	Deleted: "deleted",
	// This account is not mastered locally and has physically moved.
	Moved: "moved",
}

type AccountType string

type accountTypeValuesType struct {
	Personal     AccountType
	Organization AccountType
}

var AccountTypeValues = accountTypeValuesType{
	Personal:     "personal",
	Organization: "organization",
}

type AccountUserStatus string

type accountUserStatusValuesType struct {
	None            AccountUserStatus
	Active          AccountUserStatus
	Disabled        AccountUserStatus
	Deleted         AccountUserStatus
	Pending         AccountUserStatus
	Expired         AccountUserStatus
	PendingDisabled AccountUserStatus
}

var AccountUserStatusValues = accountUserStatusValuesType{
	None: "none",
	// User has signed in at least once to the VSTS account
	Active: "active",
	// User cannot sign in; primarily used by admin to temporarily remove a user due to absence or license reallocation
	Disabled: "disabled",
	// User is removed from the VSTS account by the VSTS account admin
	Deleted: "deleted",
	// User is invited to join the VSTS account by the VSTS account admin, but has not signed up/signed in yet
	Pending: "pending",
	// User can sign in; primarily used when license is in expired state and we give a grace period
	Expired: "expired",
	// User is disabled; if reenabled, they will still be in the Pending state
	PendingDisabled: "pendingDisabled",
}

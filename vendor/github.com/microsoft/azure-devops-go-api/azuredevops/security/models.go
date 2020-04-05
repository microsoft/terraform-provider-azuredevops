// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package security

import (
	"github.com/google/uuid"
)

// Class for encapsulating the allowed and denied permissions for a given IdentityDescriptor.
type AccessControlEntry struct {
	// The set of permission bits that represent the actions that the associated descriptor is allowed to perform.
	Allow *int `json:"allow,omitempty"`
	// The set of permission bits that represent the actions that the associated descriptor is not allowed to perform.
	Deny *int `json:"deny,omitempty"`
	// The descriptor for the user this AccessControlEntry applies to.
	Descriptor *string `json:"descriptor,omitempty"`
	// This value, when set, reports the inherited and effective information for the associated descriptor. This value is only set on AccessControlEntries returned by the QueryAccessControlList(s) call when its includeExtendedInfo parameter is set to true.
	ExtendedInfo *AceExtendedInformation `json:"extendedInfo,omitempty"`
}

// The AccessControlList class is meant to associate a set of AccessControlEntries with a security token and its inheritance settings.
type AccessControlList struct {
	// Storage of permissions keyed on the identity the permission is for.
	AcesDictionary *map[string]AccessControlEntry `json:"acesDictionary,omitempty"`
	// True if this ACL holds ACEs that have extended information.
	IncludeExtendedInfo *bool `json:"includeExtendedInfo,omitempty"`
	// True if the given token inherits permissions from parents.
	InheritPermissions *bool `json:"inheritPermissions,omitempty"`
	// The token that this AccessControlList is for.
	Token *string `json:"token,omitempty"`
}

// A list of AccessControlList. An AccessControlList is meant to associate a set of AccessControlEntries with a security token and its inheritance settings.
type AccessControlListsCollection struct {
}

// Holds the inherited and effective permission information for a given AccessControlEntry.
type AceExtendedInformation struct {
	// This is the combination of all of the explicit and inherited permissions for this identity on this token.  These are the permissions used when determining if a given user has permission to perform an action.
	EffectiveAllow *int `json:"effectiveAllow,omitempty"`
	// This is the combination of all of the explicit and inherited permissions for this identity on this token.  These are the permissions used when determining if a given user has permission to perform an action.
	EffectiveDeny *int `json:"effectiveDeny,omitempty"`
	// These are the permissions that are inherited for this identity on this token.  If the token does not inherit permissions this will be 0.  Note that any permissions that have been explicitly set on this token for this identity, or any groups that this identity is a part of, are not included here.
	InheritedAllow *int `json:"inheritedAllow,omitempty"`
	// These are the permissions that are inherited for this identity on this token.  If the token does not inherit permissions this will be 0.  Note that any permissions that have been explicitly set on this token for this identity, or any groups that this identity is a part of, are not included here.
	InheritedDeny *int `json:"inheritedDeny,omitempty"`
}

type ActionDefinition struct {
	// The bit mask integer for this action. Must be a power of 2.
	Bit *int `json:"bit,omitempty"`
	// The localized display name for this action.
	DisplayName *string `json:"displayName,omitempty"`
	// The non-localized name for this action.
	Name *string `json:"name,omitempty"`
	// The namespace that this action belongs to.  This will only be used for reading from the database.
	NamespaceId *uuid.UUID `json:"namespaceId,omitempty"`
}

// Represents an evaluated permission.
type PermissionEvaluation struct {
	// Permission bit for this evaluated permission.
	Permissions *int `json:"permissions,omitempty"`
	// Security namespace identifier for this evaluated permission.
	SecurityNamespaceId *uuid.UUID `json:"securityNamespaceId,omitempty"`
	// Security namespace-specific token for this evaluated permission.
	Token *string `json:"token,omitempty"`
	// Permission evaluation value.
	Value *bool `json:"value,omitempty"`
}

// Represents a set of evaluated permissions.
type PermissionEvaluationBatch struct {
	// True if members of the Administrators group should always pass the security check.
	AlwaysAllowAdministrators *bool `json:"alwaysAllowAdministrators,omitempty"`
	// Array of permission evaluations to evaluate.
	Evaluations *[]PermissionEvaluation `json:"evaluations,omitempty"`
}

// Class for describing the details of a TeamFoundationSecurityNamespace.
type SecurityNamespaceDescription struct {
	// The list of actions that this Security Namespace is responsible for securing.
	Actions *[]ActionDefinition `json:"actions,omitempty"`
	// This is the dataspace category that describes where the security information for this SecurityNamespace should be stored.
	DataspaceCategory *string `json:"dataspaceCategory,omitempty"`
	// This localized name for this namespace.
	DisplayName *string `json:"displayName,omitempty"`
	// If the security tokens this namespace will be operating on need to be split on certain character lengths to determine its elements, that length should be specified here. If not, this value will be -1.
	ElementLength *int `json:"elementLength,omitempty"`
	// This is the type of the extension that should be loaded from the plugins directory for extending this security namespace.
	ExtensionType *string `json:"extensionType,omitempty"`
	// If true, the security namespace is remotable, allowing another service to proxy the namespace.
	IsRemotable *bool `json:"isRemotable,omitempty"`
	// This non-localized for this namespace.
	Name *string `json:"name,omitempty"`
	// The unique identifier for this namespace.
	NamespaceId *uuid.UUID `json:"namespaceId,omitempty"`
	// The permission bits needed by a user in order to read security data on the Security Namespace.
	ReadPermission *int `json:"readPermission,omitempty"`
	// If the security tokens this namespace will be operating on need to be split on certain characters to determine its elements that character should be specified here. If not, this value will be the null character.
	SeparatorValue *rune `json:"separatorValue,omitempty"`
	// Used to send information about the structure of the security namespace over the web service.
	StructureValue *int `json:"structureValue,omitempty"`
	// The bits reserved by system store
	SystemBitMask *int `json:"systemBitMask,omitempty"`
	// If true, the security service will expect an ISecurityDataspaceTokenTranslator plugin to exist for this namespace
	UseTokenTranslator *bool `json:"useTokenTranslator,omitempty"`
	// The permission bits needed by a user in order to modify security data on the Security Namespace.
	WritePermission *int `json:"writePermission,omitempty"`
}

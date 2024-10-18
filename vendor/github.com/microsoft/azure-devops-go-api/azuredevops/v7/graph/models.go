// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package graph

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
)

type AadGraphMember struct {
	// This field contains zero or more interesting links about the graph subject. These links may be invoked to obtain additional relationships or more detailed information about this graph subject.
	Links interface{} `json:"_links,omitempty"`
	// The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the same graph subject across both Accounts and Organizations.
	Descriptor *string `json:"descriptor,omitempty"`
	// This is the non-unique display name of the graph subject. To change this field, you must alter its value in the source provider.
	DisplayName *string `json:"displayName,omitempty"`
	// This url is the full route to the source resource of this graph subject.
	Url *string `json:"url,omitempty"`
	// [Internal Use Only] The legacy descriptor is here in case you need to access old version IMS using identity descriptor.
	LegacyDescriptor *string `json:"legacyDescriptor,omitempty"`
	// The type of source provider for the origin identifier (ex:AD, AAD, MSA)
	Origin *string `json:"origin,omitempty"`
	// The unique identifier from the system of origin. Typically a sid, object id or Guid. Linking and unlinking operations can cause this value to change for a user because the user is not backed by a different provider and has a different unique id in the new provider.
	OriginId *string `json:"originId,omitempty"`
	// This field identifies the type of the graph subject (ex: Group, Scope, User).
	SubjectKind *string `json:"subjectKind,omitempty"`
	// This represents the name of the container of origin for a graph member. (For MSA this is "Windows Live ID", for AD the name of the domain, for AAD the tenantID of the directory, for VSTS groups the ScopeId, etc)
	Domain *string `json:"domain,omitempty"`
	// The email address of record for a given graph member. This may be different than the principal name.
	MailAddress *string `json:"mailAddress,omitempty"`
	// This is the PrincipalName of this graph member from the source provider. The source provider may change this field over time and it is not guaranteed to be immutable for the life of the graph member by VSTS.
	PrincipalName *string `json:"principalName,omitempty"`
	// The short, generally unique name for the user in the backing directory. For AAD users, this corresponds to the mail nickname, which is often but not necessarily similar to the part of the user's mail address before the @ sign. For GitHub users, this corresponds to the GitHub user handle.
	DirectoryAlias *string `json:"directoryAlias,omitempty"`
	// When true, the group has been deleted in the identity provider
	IsDeletedInOrigin *bool `json:"isDeletedInOrigin,omitempty"`
	// The meta type of the user in the origin, such as "member", "guest", etc. See UserMetaType for the set of possible values.
	MetaType *string `json:"metaType,omitempty"`
}

type GraphCachePolicies struct {
	// Size of the cache
	CacheSize *int `json:"cacheSize,omitempty"`
}

// Subject descriptor of a Graph entity
type GraphDescriptorResult struct {
	// This field contains zero or more interesting links about the graph descriptor. These links may be invoked to obtain additional relationships or more detailed information about this graph descriptor.
	Links interface{} `json:"_links,omitempty"`
	Value *string     `json:"value,omitempty"`
}

type GraphGlobalExtendedPropertyBatch struct {
	PropertyNameFilters *[]string `json:"propertyNameFilters,omitempty"`
	SubjectDescriptors  *[]string `json:"subjectDescriptors,omitempty"`
}

// Graph group entity
type GraphGroup struct {
	// This field contains zero or more interesting links about the graph subject. These links may be invoked to obtain additional relationships or more detailed information about this graph subject.
	Links interface{} `json:"_links,omitempty"`
	// The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the same graph subject across both Accounts and Organizations.
	Descriptor *string `json:"descriptor,omitempty"`
	// This is the non-unique display name of the graph subject. To change this field, you must alter its value in the source provider.
	DisplayName *string `json:"displayName,omitempty"`
	// This url is the full route to the source resource of this graph subject.
	Url *string `json:"url,omitempty"`
	// [Internal Use Only] The legacy descriptor is here in case you need to access old version IMS using identity descriptor.
	LegacyDescriptor *string `json:"legacyDescriptor,omitempty"`
	// The type of source provider for the origin identifier (ex:AD, AAD, MSA)
	Origin *string `json:"origin,omitempty"`
	// The unique identifier from the system of origin. Typically a sid, object id or Guid. Linking and unlinking operations can cause this value to change for a user because the user is not backed by a different provider and has a different unique id in the new provider.
	OriginId *string `json:"originId,omitempty"`
	// This field identifies the type of the graph subject (ex: Group, Scope, User).
	SubjectKind *string `json:"subjectKind,omitempty"`
	// This represents the name of the container of origin for a graph member. (For MSA this is "Windows Live ID", for AD the name of the domain, for AAD the tenantID of the directory, for VSTS groups the ScopeId, etc)
	Domain *string `json:"domain,omitempty"`
	// The email address of record for a given graph member. This may be different than the principal name.
	MailAddress *string `json:"mailAddress,omitempty"`
	// This is the PrincipalName of this graph member from the source provider. The source provider may change this field over time and it is not guaranteed to be immutable for the life of the graph member by VSTS.
	PrincipalName *string `json:"principalName,omitempty"`
	// A short phrase to help human readers disambiguate groups with similar names
	Description *string `json:"description,omitempty"`
	// Whether the group has been deleted
	IsDeleted *bool `json:"isDeleted,omitempty"`
}

// Do not attempt to use this type to create a new group. This type does not contain sufficient fields to create a new group.
type GraphGroupCreationContext struct {
	// Optional: If provided, we will use this identifier for the storage key of the created group
	StorageKey *uuid.UUID `json:"storageKey,omitempty"`
}

// Use this type to create a new group using the mail address as a reference to an existing group from an external AD or AAD backed provider. This is the subset of GraphGroup fields required for creation of a group for the AAD and AD use case.
type GraphGroupMailAddressCreationContext struct {
	// Optional: If provided, we will use this identifier for the storage key of the created group
	StorageKey *uuid.UUID `json:"storageKey,omitempty"`
	// This should be the mail address or the group in the source AD or AAD provider. Example: jamal@contoso.com Team Services will communicate with the source provider to fill all other fields on creation.
	MailAddress *string `json:"mailAddress,omitempty"`
}

// Use this type to create a new group using the OriginID as a reference to an existing group from an external AD or AAD backed provider. This is the subset of GraphGroup fields required for creation of a group for the AD and AAD use case.
type GraphGroupOriginIdCreationContext struct {
	// Optional: If provided, we will use this identifier for the storage key of the created group
	StorageKey *uuid.UUID `json:"storageKey,omitempty"`
	// This should be the object id or sid of the group from the source AD or AAD provider. Example: d47d025a-ce2f-4a79-8618-e8862ade30dd Team Services will communicate with the source provider to fill all other fields on creation.
	OriginId *string `json:"originId,omitempty"`
}

// Use this type to create a new Vsts group that is not backed by an external provider.
type GraphGroupVstsCreationContext struct {
	// Optional: If provided, we will use this identifier for the storage key of the created group
	StorageKey *uuid.UUID `json:"storageKey,omitempty"`
	// For internal use only in back compat scenarios.
	CrossProject *bool `json:"crossProject,omitempty"`
	// Used by VSTS groups; if set this will be the group description, otherwise ignored
	Description *string `json:"description,omitempty"`
	Descriptor  *string `json:"descriptor,omitempty"`
	// Used by VSTS groups; if set this will be the group DisplayName, otherwise ignored
	DisplayName *string `json:"displayName,omitempty"`
	// For internal use only in back compat scenarios.
	RestrictedVisibility *bool `json:"restrictedVisibility,omitempty"`
	// For internal use only in back compat scenarios.
	SpecialGroupType *string `json:"specialGroupType,omitempty"`
}

type GraphMember struct {
	// This field contains zero or more interesting links about the graph subject. These links may be invoked to obtain additional relationships or more detailed information about this graph subject.
	Links interface{} `json:"_links,omitempty"`
	// The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the same graph subject across both Accounts and Organizations.
	Descriptor *string `json:"descriptor,omitempty"`
	// This is the non-unique display name of the graph subject. To change this field, you must alter its value in the source provider.
	DisplayName *string `json:"displayName,omitempty"`
	// This url is the full route to the source resource of this graph subject.
	Url *string `json:"url,omitempty"`
	// [Internal Use Only] The legacy descriptor is here in case you need to access old version IMS using identity descriptor.
	LegacyDescriptor *string `json:"legacyDescriptor,omitempty"`
	// The type of source provider for the origin identifier (ex:AD, AAD, MSA)
	Origin *string `json:"origin,omitempty"`
	// The unique identifier from the system of origin. Typically a sid, object id or Guid. Linking and unlinking operations can cause this value to change for a user because the user is not backed by a different provider and has a different unique id in the new provider.
	OriginId *string `json:"originId,omitempty"`
	// This field identifies the type of the graph subject (ex: Group, Scope, User).
	SubjectKind *string `json:"subjectKind,omitempty"`
	// This represents the name of the container of origin for a graph member. (For MSA this is "Windows Live ID", for AD the name of the domain, for AAD the tenantID of the directory, for VSTS groups the ScopeId, etc)
	Domain *string `json:"domain,omitempty"`
	// The email address of record for a given graph member. This may be different than the principal name.
	MailAddress *string `json:"mailAddress,omitempty"`
	// This is the PrincipalName of this graph member from the source provider. The source provider may change this field over time and it is not guaranteed to be immutable for the life of the graph member by VSTS.
	PrincipalName *string `json:"principalName,omitempty"`
}

// Relationship between a container and a member
type GraphMembership struct {
	// This field contains zero or more interesting links about the graph membership. These links may be invoked to obtain additional relationships or more detailed information about this graph membership.
	Links               interface{} `json:"_links,omitempty"`
	ContainerDescriptor *string     `json:"containerDescriptor,omitempty"`
	MemberDescriptor    *string     `json:"memberDescriptor,omitempty"`
}

// Status of a Graph membership (active/inactive)
type GraphMembershipState struct {
	// This field contains zero or more interesting links about the graph membership state. These links may be invoked to obtain additional relationships or more detailed information about this graph membership state.
	Links interface{} `json:"_links,omitempty"`
	// When true, the membership is active
	Active *bool `json:"active,omitempty"`
}

type GraphMembershipTraversal struct {
	// Reason why the subject could not be traversed completely
	IncompletenessReason *string `json:"incompletenessReason,omitempty"`
	// When true, the subject is traversed completely
	IsComplete *bool `json:"isComplete,omitempty"`
	// The traversed subject descriptor
	SubjectDescriptor *string `json:"subjectDescriptor,omitempty"`
	// Subject descriptor ids of the traversed members
	TraversedSubjectIds *[]uuid.UUID `json:"traversedSubjectIds,omitempty"`
	// Subject descriptors of the traversed members
	TraversedSubjects *[]string `json:"traversedSubjects,omitempty"`
}

// Who is the provider for this user and what is the identifier and domain that is used to uniquely identify the user.
type GraphProviderInfo struct {
	// The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the same graph subject across both Accounts and Organizations.
	Descriptor *string `json:"descriptor,omitempty"`
	// This represents the name of the container of origin for a graph member. (For MSA this is "Windows Live ID", for AAD the tenantID of the directory.)
	Domain *string `json:"domain,omitempty"`
	// The type of source provider for the origin identifier (ex: "aad", "msa")
	Origin *string `json:"origin,omitempty"`
	// The unique identifier from the system of origin. (For MSA this is the PUID in hex notation, for AAD this is the object id.)
	OriginId *string `json:"originId,omitempty"`
}

// Container where a graph entity is defined (organization, project, team)
type GraphScope struct {
	// This field contains zero or more interesting links about the graph subject. These links may be invoked to obtain additional relationships or more detailed information about this graph subject.
	Links interface{} `json:"_links,omitempty"`
	// The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the same graph subject across both Accounts and Organizations.
	Descriptor *string `json:"descriptor,omitempty"`
	// This is the non-unique display name of the graph subject. To change this field, you must alter its value in the source provider.
	DisplayName *string `json:"displayName,omitempty"`
	// This url is the full route to the source resource of this graph subject.
	Url *string `json:"url,omitempty"`
	// [Internal Use Only] The legacy descriptor is here in case you need to access old version IMS using identity descriptor.
	LegacyDescriptor *string `json:"legacyDescriptor,omitempty"`
	// The type of source provider for the origin identifier (ex:AD, AAD, MSA)
	Origin *string `json:"origin,omitempty"`
	// The unique identifier from the system of origin. Typically a sid, object id or Guid. Linking and unlinking operations can cause this value to change for a user because the user is not backed by a different provider and has a different unique id in the new provider.
	OriginId *string `json:"originId,omitempty"`
	// This field identifies the type of the graph subject (ex: Group, Scope, User).
	SubjectKind *string `json:"subjectKind,omitempty"`
	// The subject descriptor that references the administrators group for this scope. Only members of this group can change the contents of this scope or assign other users permissions to access this scope.
	AdministratorDescriptor *string `json:"administratorDescriptor,omitempty"`
	// When true, this scope is also a securing host for one or more scopes.
	IsGlobal *bool `json:"isGlobal,omitempty"`
	// The subject descriptor for the closest account or organization in the ancestor tree of this scope.
	ParentDescriptor *string `json:"parentDescriptor,omitempty"`
	// The type of this scope. Typically ServiceHost or TeamProject.
	ScopeType *identity.GroupScopeType `json:"scopeType,omitempty"`
	// The subject descriptor for the containing organization in the ancestor tree of this scope.
	SecuringHostDescriptor *string `json:"securingHostDescriptor,omitempty"`
}

// This type is the subset of fields that can be provided by the user to create a Vsts scope. Scope creation is currently limited to internal back-compat scenarios. End users that attempt to create a scope with this API will fail.
type GraphScopeCreationContext struct {
	// Set this field to override the default description of this scope's admin group.
	AdminGroupDescription *string `json:"adminGroupDescription,omitempty"`
	// All scopes have an Administrator Group that controls access to the contents of the scope. Set this field to use a non-default group name for that administrators group.
	AdminGroupName *string `json:"adminGroupName,omitempty"`
	// Set this optional field if this scope is created on behalf of a user other than the user making the request. This should be the Id of the user that is not the requester.
	CreatorId *uuid.UUID `json:"creatorId,omitempty"`
	// The scope must be provided with a unique name within the parent scope. This means the created scope can have a parent or child with the same name, but no siblings with the same name.
	Name *string `json:"name,omitempty"`
	// The type of scope being created.
	ScopeType *identity.GroupScopeType `json:"scopeType,omitempty"`
	// An optional ID that uniquely represents the scope within it's parent scope. If this parameter is not provided, Vsts will generate on automatically.
	StorageKey *uuid.UUID `json:"storageKey,omitempty"`
}

type GraphServicePrincipal struct {
	// This field contains zero or more interesting links about the graph subject. These links may be invoked to obtain additional relationships or more detailed information about this graph subject.
	Links interface{} `json:"_links,omitempty"`
	// The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the same graph subject across both Accounts and Organizations.
	Descriptor *string `json:"descriptor,omitempty"`
	// This is the non-unique display name of the graph subject. To change this field, you must alter its value in the source provider.
	DisplayName *string `json:"displayName,omitempty"`
	// This url is the full route to the source resource of this graph subject.
	Url *string `json:"url,omitempty"`
	// [Internal Use Only] The legacy descriptor is here in case you need to access old version IMS using identity descriptor.
	LegacyDescriptor *string `json:"legacyDescriptor,omitempty"`
	// The type of source provider for the origin identifier (ex:AD, AAD, MSA)
	Origin *string `json:"origin,omitempty"`
	// The unique identifier from the system of origin. Typically a sid, object id or Guid. Linking and unlinking operations can cause this value to change for a user because the user is not backed by a different provider and has a different unique id in the new provider.
	OriginId *string `json:"originId,omitempty"`
	// This field identifies the type of the graph subject (ex: Group, Scope, User).
	SubjectKind *string `json:"subjectKind,omitempty"`
	// This represents the name of the container of origin for a graph member. (For MSA this is "Windows Live ID", for AD the name of the domain, for AAD the tenantID of the directory, for VSTS groups the ScopeId, etc)
	Domain *string `json:"domain,omitempty"`
	// The email address of record for a given graph member. This may be different than the principal name.
	MailAddress *string `json:"mailAddress,omitempty"`
	// This is the PrincipalName of this graph member from the source provider. The source provider may change this field over time and it is not guaranteed to be immutable for the life of the graph member by VSTS.
	PrincipalName *string `json:"principalName,omitempty"`
	// The short, generally unique name for the user in the backing directory. For AAD users, this corresponds to the mail nickname, which is often but not necessarily similar to the part of the user's mail address before the @ sign. For GitHub users, this corresponds to the GitHub user handle.
	DirectoryAlias *string `json:"directoryAlias,omitempty"`
	// When true, the group has been deleted in the identity provider
	IsDeletedInOrigin *bool `json:"isDeletedInOrigin,omitempty"`
	// The meta type of the user in the origin, such as "member", "guest", etc. See UserMetaType for the set of possible values.
	MetaType      *string `json:"metaType,omitempty"`
	ApplicationId *string `json:"applicationId,omitempty"`
}

// Do not attempt to use this type to create a new service principal. Use one of the subclasses instead. This type does not contain sufficient fields to create a new service principal.
type GraphServicePrincipalCreationContext struct {
	// Optional: If provided, we will use this identifier for the storage key of the created service principal
	StorageKey *uuid.UUID `json:"storageKey,omitempty"`
}

// Use this type to create a new service principal using the OriginID as a reference to an existing service principal from an external AAD backed provider. This is the subset of GraphServicePrincipal fields required for creation of a GraphServicePrincipal for the AAD use case when looking up the service principal by its unique ID in the backing provider.
type GraphServicePrincipalOriginIdCreationContext struct {
	// Optional: If provided, we will use this identifier for the storage key of the created service principal
	StorageKey *uuid.UUID `json:"storageKey,omitempty"`
	// This should be the object id of the service principal from the AAD provider. Example: d47d025a-ce2f-4a79-8618-e8862ade30dd Team Services will communicate with the source provider to fill all other fields on creation.
	OriginId *string `json:"originId,omitempty"`
}

// Use this type to update an existing service principal using the OriginID as a reference to an existing service principal from an external AAD backed provider. This is the subset of GraphServicePrincipal fields required for creation of a GraphServicePrincipal for AAD use case when looking up the service principal by its unique ID in the backing provider.
type GraphServicePrincipalOriginIdUpdateContext struct {
	// Storage key should not be specified in case of updating service principal
	StorageKey *uuid.UUID `json:"storageKey,omitempty"`
	// This should be the object id or sid of the service principal from the source AAD provider. Example: d47d025a-ce2f-4a79-8618-e8862ade30dd Azure Devops will communicate with the source provider to fill all other fields on creation.
	OriginId *string `json:"originId,omitempty"`
}

// Do not attempt to use this type to update service principal. Use one of the subclasses instead. This type does not contain sufficient fields to create a new service principal.
type GraphServicePrincipalUpdateContext struct {
	// Deprecated:
	StorageKey *uuid.UUID `json:"storageKey,omitempty"`
}

// Storage key of a Graph entity
type GraphStorageKeyResult struct {
	// This field contains zero or more interesting links about the graph storage key. These links may be invoked to obtain additional relationships or more detailed information about this graph storage key.
	Links interface{} `json:"_links,omitempty"`
	Value *uuid.UUID  `json:"value,omitempty"`
}

// Top-level graph entity
type GraphSubject struct {
	// This field contains zero or more interesting links about the graph subject. These links may be invoked to obtain additional relationships or more detailed information about this graph subject.
	Links interface{} `json:"_links,omitempty"`
	// The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the same graph subject across both Accounts and Organizations.
	Descriptor *string `json:"descriptor,omitempty"`
	// This is the non-unique display name of the graph subject. To change this field, you must alter its value in the source provider.
	DisplayName *string `json:"displayName,omitempty"`
	// This url is the full route to the source resource of this graph subject.
	Url *string `json:"url,omitempty"`
	// [Internal Use Only] The legacy descriptor is here in case you need to access old version IMS using identity descriptor.
	LegacyDescriptor *string `json:"legacyDescriptor,omitempty"`
	// The type of source provider for the origin identifier (ex:AD, AAD, MSA)
	Origin *string `json:"origin,omitempty"`
	// The unique identifier from the system of origin. Typically a sid, object id or Guid. Linking and unlinking operations can cause this value to change for a user because the user is not backed by a different provider and has a different unique id in the new provider.
	OriginId *string `json:"originId,omitempty"`
	// This field identifies the type of the graph subject (ex: Group, Scope, User).
	SubjectKind *string `json:"subjectKind,omitempty"`
}

type GraphSubjectBase struct {
	// This field contains zero or more interesting links about the graph subject. These links may be invoked to obtain additional relationships or more detailed information about this graph subject.
	Links interface{} `json:"_links,omitempty"`
	// The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the same graph subject across both Accounts and Organizations.
	Descriptor *string `json:"descriptor,omitempty"`
	// This is the non-unique display name of the graph subject. To change this field, you must alter its value in the source provider.
	DisplayName *string `json:"displayName,omitempty"`
	// This url is the full route to the source resource of this graph subject.
	Url *string `json:"url,omitempty"`
}

// Batching of subjects to lookup using the Graph API
type GraphSubjectLookup struct {
	LookupKeys *[]GraphSubjectLookupKey `json:"lookupKeys,omitempty"`
}

type GraphSubjectLookupKey struct {
	Descriptor *string `json:"descriptor,omitempty"`
}

// Subject to search using the Graph API
type GraphSubjectQuery struct {
	// Search term to search for Azure Devops users or/and groups
	Query *string `json:"query,omitempty"`
	// Optional parameter. Specify a non-default scope (collection, project) to search for users or groups within the scope.
	ScopeDescriptor *string `json:"scopeDescriptor,omitempty"`
	// "User" or "Group" can be specified, both or either
	SubjectKind *[]string `json:"subjectKind,omitempty"`
}

type GraphSystemSubject struct {
	// This field contains zero or more interesting links about the graph subject. These links may be invoked to obtain additional relationships or more detailed information about this graph subject.
	Links interface{} `json:"_links,omitempty"`
	// The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the same graph subject across both Accounts and Organizations.
	Descriptor *string `json:"descriptor,omitempty"`
	// This is the non-unique display name of the graph subject. To change this field, you must alter its value in the source provider.
	DisplayName *string `json:"displayName,omitempty"`
	// This url is the full route to the source resource of this graph subject.
	Url *string `json:"url,omitempty"`
	// [Internal Use Only] The legacy descriptor is here in case you need to access old version IMS using identity descriptor.
	LegacyDescriptor *string `json:"legacyDescriptor,omitempty"`
	// The type of source provider for the origin identifier (ex:AD, AAD, MSA)
	Origin *string `json:"origin,omitempty"`
	// The unique identifier from the system of origin. Typically a sid, object id or Guid. Linking and unlinking operations can cause this value to change for a user because the user is not backed by a different provider and has a different unique id in the new provider.
	OriginId *string `json:"originId,omitempty"`
	// This field identifies the type of the graph subject (ex: Group, Scope, User).
	SubjectKind *string `json:"subjectKind,omitempty"`
}

type GraphTraversalDirection string

type graphTraversalDirectionValuesType struct {
	Unknown GraphTraversalDirection
	Down    GraphTraversalDirection
	Up      GraphTraversalDirection
}

var GraphTraversalDirectionValues = graphTraversalDirectionValuesType{
	Unknown: "unknown",
	Down:    "down",
	Up:      "up",
}

type GraphUser struct {
	// This field contains zero or more interesting links about the graph subject. These links may be invoked to obtain additional relationships or more detailed information about this graph subject.
	Links interface{} `json:"_links,omitempty"`
	// The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the same graph subject across both Accounts and Organizations.
	Descriptor *string `json:"descriptor,omitempty"`
	// This is the non-unique display name of the graph subject. To change this field, you must alter its value in the source provider.
	DisplayName *string `json:"displayName,omitempty"`
	// This url is the full route to the source resource of this graph subject.
	Url *string `json:"url,omitempty"`
	// [Internal Use Only] The legacy descriptor is here in case you need to access old version IMS using identity descriptor.
	LegacyDescriptor *string `json:"legacyDescriptor,omitempty"`
	// The type of source provider for the origin identifier (ex:AD, AAD, MSA)
	Origin *string `json:"origin,omitempty"`
	// The unique identifier from the system of origin. Typically a sid, object id or Guid. Linking and unlinking operations can cause this value to change for a user because the user is not backed by a different provider and has a different unique id in the new provider.
	OriginId *string `json:"originId,omitempty"`
	// This field identifies the type of the graph subject (ex: Group, Scope, User).
	SubjectKind *string `json:"subjectKind,omitempty"`
	// This represents the name of the container of origin for a graph member. (For MSA this is "Windows Live ID", for AD the name of the domain, for AAD the tenantID of the directory, for VSTS groups the ScopeId, etc)
	Domain *string `json:"domain,omitempty"`
	// The email address of record for a given graph member. This may be different than the principal name.
	MailAddress *string `json:"mailAddress,omitempty"`
	// This is the PrincipalName of this graph member from the source provider. The source provider may change this field over time and it is not guaranteed to be immutable for the life of the graph member by VSTS.
	PrincipalName *string `json:"principalName,omitempty"`
	// The short, generally unique name for the user in the backing directory. For AAD users, this corresponds to the mail nickname, which is often but not necessarily similar to the part of the user's mail address before the @ sign. For GitHub users, this corresponds to the GitHub user handle.
	DirectoryAlias *string `json:"directoryAlias,omitempty"`
	// When true, the group has been deleted in the identity provider
	IsDeletedInOrigin *bool `json:"isDeletedInOrigin,omitempty"`
	// The meta type of the user in the origin, such as "member", "guest", etc. See UserMetaType for the set of possible values.
	MetaType *string `json:"metaType,omitempty"`
}

// Do not attempt to use this type to create a new user. Use one of the subclasses instead. This type does not contain sufficient fields to create a new user.
type GraphUserCreationContext struct {
	// Optional: If provided, we will use this identifier for the storage key of the created user
	StorageKey *uuid.UUID `json:"storageKey,omitempty"`
}

// Use this type to create a new user using the mail address as a reference to an existing user from an external AD or AAD backed provider. This is the subset of GraphUser fields required for creation of a GraphUser for the AD and AAD use case when looking up the user by its mail address in the backing provider.
type GraphUserMailAddressCreationContext struct {
	// Optional: If provided, we will use this identifier for the storage key of the created user
	StorageKey *uuid.UUID `json:"storageKey,omitempty"`
	// This should be the mail address of the user in the source AD or AAD provider. Example: Jamal.Hartnett@contoso.com Team Services will communicate with the source provider to fill all other fields on creation.
	MailAddress *string `json:"mailAddress,omitempty"`
}

// Use this type to create a new user using the OriginID as a reference to an existing user from an external AD or AAD backed provider. This is the subset of GraphUser fields required for creation of a GraphUser for the AD and AAD use case when looking up the user by its unique ID in the backing provider.
type GraphUserOriginIdCreationContext struct {
	// Optional: If provided, we will use this identifier for the storage key of the created user
	StorageKey *uuid.UUID `json:"storageKey,omitempty"`
	// This should be the name of the origin provider. Example: github.com
	Origin *string `json:"origin,omitempty"`
	// This should be the object id or sid of the user from the source AD or AAD provider. Example: d47d025a-ce2f-4a79-8618-e8862ade30dd Team Services will communicate with the source provider to fill all other fields on creation.
	OriginId *string `json:"originId,omitempty"`
}

// Use this type to update an existing user using the OriginID as a reference to an existing user from an external AD or AAD backed provider. This is the subset of GraphUser fields required for creation of a GraphUser for the AD and AAD use case when looking up the user by its unique ID in the backing provider.
type GraphUserOriginIdUpdateContext struct {
	// Storage key should not be specified in case of updating user
	StorageKey *uuid.UUID `json:"storageKey,omitempty"`
	// This should be the object id or sid of the user from the source AD or AAD provider. Example: d47d025a-ce2f-4a79-8618-e8862ade30dd Azure Devops will communicate with the source provider to fill all other fields on creation.
	OriginId *string `json:"originId,omitempty"`
}

// Use this type to create a new user using the principal name as a reference to an existing user from an external AD or AAD backed provider. This is the subset of GraphUser fields required for creation of a GraphUser for the AD and AAD use case when looking up the user by its principal name in the backing provider.
type GraphUserPrincipalNameCreationContext struct {
	// Optional: If provided, we will use this identifier for the storage key of the created user
	StorageKey *uuid.UUID `json:"storageKey,omitempty"`
	// This should be the principal name or upn of the user in the source AD or AAD provider. Example: jamal@contoso.com Team Services will communicate with the source provider to fill all other fields on creation.
	PrincipalName *string `json:"principalName,omitempty"`
}

// Use this type for transfering identity rights, for instance after performing a Tenant switch.
type GraphUserPrincipalNameUpdateContext struct {
	// Storage key should not be specified in case of updating user
	StorageKey *uuid.UUID `json:"storageKey,omitempty"`
	// This should be Principal Name (UPN) to which we want to transfer rights. Example: destination@email.com
	PrincipalName *string `json:"principalName,omitempty"`
}

// Do not attempt to use this type to update user. Use one of the subclasses instead. This type does not contain sufficient fields to create a new user.
type GraphUserUpdateContext struct {
	// Deprecated:
	StorageKey *uuid.UUID `json:"storageKey,omitempty"`
}

type IdentityMapping struct {
	Source *UserPrincipalName `json:"source,omitempty"`
	Target *UserPrincipalName `json:"target,omitempty"`
}

type IdentityMappings struct {
	Mappings *[]IdentityMapping `json:"mappings,omitempty"`
}

type MappingResult struct {
	Code         *string `json:"code,omitempty"`
	ErrorMessage *string `json:"errorMessage,omitempty"`
}

type PagedGraphGroups struct {
	// This will be non-null if there is another page of data. There will never be more than one continuation token returned by a request.
	ContinuationToken *[]string `json:"continuationToken,omitempty"`
	// The enumerable list of groups found within a page.
	GraphGroups *[]GraphGroup `json:"graphGroups,omitempty"`
}

type PagedGraphMembers struct {
	// This will be non-null if there is another page of data. There will never be more than one continuation token returned by a request.
	ContinuationToken *[]string `json:"continuationToken,omitempty"`
	// The enumerable list of members found within a page.
	GraphMembers *[]GraphMember `json:"graphMembers,omitempty"`
}

type PagedGraphServicePrincipals struct {
	// This will be non-null if there is another page of data. There will never be more than one continuation token returned by a request.
	ContinuationToken *[]string `json:"continuationToken,omitempty"`
	// The enumerable list of service principals found within a page.
	GraphServicePrincipals *[]GraphServicePrincipal `json:"graphServicePrincipals,omitempty"`
}

type PagedGraphUsers struct {
	// This will be non-null if there is another page of data. There will never be more than one continuation token returned by a request.
	ContinuationToken *[]string `json:"continuationToken,omitempty"`
	// The enumerable set of users found within a page.
	GraphUsers *[]GraphUser `json:"graphUsers,omitempty"`
}

type RequestAccessPayLoad struct {
	Message      *string `json:"message,omitempty"`
	ProjectUri   *string `json:"projectUri,omitempty"`
	UrlRequested *string `json:"urlRequested,omitempty"`
}

type ResolveDisconnectedUsersResponse struct {
	Code           *string          `json:"code,omitempty"`
	ErrorMessage   *string          `json:"errorMessage,omitempty"`
	MappingResults *[]MappingResult `json:"mappingResults,omitempty"`
}

type UserPrincipalName struct {
	PrincipalName *string `json:"principalName,omitempty"`
}

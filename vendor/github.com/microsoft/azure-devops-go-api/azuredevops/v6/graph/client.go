// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package graph

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/profile"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/webapi"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var ResourceAreaId, _ = uuid.Parse("bb1e7ec9-e901-4b68-999a-de7012b920f8")

type Client interface {
	// [Preview API] Create a new membership between a container and subject.
	AddMembership(context.Context, AddMembershipArgs) (*GraphMembership, error)
	// [Preview API] Check to see if a membership relationship between a container and subject exists.
	CheckMembershipExistence(context.Context, CheckMembershipExistenceArgs) error
	// [Preview API] Create a new Azure DevOps group or materialize an existing AAD group.
	CreateGroup(context.Context, CreateGroupArgs) (*GraphGroup, error)
	// [Preview API] Materialize an existing AAD or MSA user into the VSTS account.
	CreateUser(context.Context, CreateUserArgs) (*GraphUser, error)
	// [Preview API]
	DeleteAvatar(context.Context, DeleteAvatarArgs) error
	// [Preview API] Removes an Azure DevOps group from all of its parent groups.
	DeleteGroup(context.Context, DeleteGroupArgs) error
	// [Preview API] Disables a user.
	DeleteUser(context.Context, DeleteUserArgs) error
	// [Preview API]
	GetAvatar(context.Context, GetAvatarArgs) (*profile.Avatar, error)
	// [Preview API] Resolve a storage key to a descriptor
	GetDescriptor(context.Context, GetDescriptorArgs) (*GraphDescriptorResult, error)
	// [Preview API] Get a group by its descriptor.
	GetGroup(context.Context, GetGroupArgs) (*GraphGroup, error)
	// [Preview API] Get a membership relationship between a container and subject.
	GetMembership(context.Context, GetMembershipArgs) (*GraphMembership, error)
	// [Preview API] Check whether a subject is active or inactive.
	GetMembershipState(context.Context, GetMembershipStateArgs) (*GraphMembershipState, error)
	// [Preview API]
	GetProviderInfo(context.Context, GetProviderInfoArgs) (*GraphProviderInfo, error)
	// [Preview API] Resolve a descriptor to a storage key.
	GetStorageKey(context.Context, GetStorageKeyArgs) (*GraphStorageKeyResult, error)
	// [Preview API] Get a user by its descriptor.
	GetUser(context.Context, GetUserArgs) (*GraphUser, error)
	// [Preview API] Gets a list of all groups in the current scope (usually organization or account).
	ListGroups(context.Context, ListGroupsArgs) (*PagedGraphGroups, error)
	// [Preview API] Get all the memberships where this descriptor is a member in the relationship.
	ListMemberships(context.Context, ListMembershipsArgs) (*[]GraphMembership, error)
	// [Preview API] Get a list of all users in a given scope.
	ListUsers(context.Context, ListUsersArgs) (*PagedGraphUsers, error)
	// [Preview API] Resolve descriptors to users, groups or scopes (Subjects) in a batch.
	LookupSubjects(context.Context, LookupSubjectsArgs) (*map[string]GraphSubject, error)
	// [Preview API] Search for Azure Devops users, or/and groups. Results will be returned in a batch with no more than 100 graph subjects.
	QuerySubjects(context.Context, QuerySubjectsArgs) (*[]GraphSubject, error)
	// [Preview API] Deletes a membership between a container and subject.
	RemoveMembership(context.Context, RemoveMembershipArgs) error
	// [Preview API]
	RequestAccess(context.Context, RequestAccessArgs) error
	// [Preview API]
	SetAvatar(context.Context, SetAvatarArgs) error
	// [Preview API] Update the properties of an Azure DevOps group.
	UpdateGroup(context.Context, UpdateGroupArgs) (*GraphGroup, error)
	// [Preview API] Map an existing user to a different identity
	UpdateUser(context.Context, UpdateUserArgs) (*GraphUser, error)
}

type ClientImpl struct {
	Client azuredevops.Client
}

func NewClient(ctx context.Context, connection *azuredevops.Connection) (Client, error) {
	client, err := connection.GetClientByResourceAreaId(ctx, ResourceAreaId)
	if err != nil {
		return nil, err
	}
	return &ClientImpl{
		Client: *client,
	}, nil
}

// [Preview API] Create a new membership between a container and subject.
func (client *ClientImpl) AddMembership(ctx context.Context, args AddMembershipArgs) (*GraphMembership, error) {
	routeValues := make(map[string]string)
	if args.SubjectDescriptor == nil || *args.SubjectDescriptor == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubjectDescriptor"}
	}
	routeValues["subjectDescriptor"] = *args.SubjectDescriptor
	if args.ContainerDescriptor == nil || *args.ContainerDescriptor == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ContainerDescriptor"}
	}
	routeValues["containerDescriptor"] = *args.ContainerDescriptor

	locationId, _ := uuid.Parse("3fd2e6ca-fb30-443a-b579-95b19ed0934c")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GraphMembership
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the AddMembership function
type AddMembershipArgs struct {
	// (required) A descriptor to a group or user that can be the child subject in the relationship.
	SubjectDescriptor *string
	// (required) A descriptor to a group that can be the container in the relationship.
	ContainerDescriptor *string
}

// [Preview API] Check to see if a membership relationship between a container and subject exists.
func (client *ClientImpl) CheckMembershipExistence(ctx context.Context, args CheckMembershipExistenceArgs) error {
	routeValues := make(map[string]string)
	if args.SubjectDescriptor == nil || *args.SubjectDescriptor == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubjectDescriptor"}
	}
	routeValues["subjectDescriptor"] = *args.SubjectDescriptor
	if args.ContainerDescriptor == nil || *args.ContainerDescriptor == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ContainerDescriptor"}
	}
	routeValues["containerDescriptor"] = *args.ContainerDescriptor

	locationId, _ := uuid.Parse("3fd2e6ca-fb30-443a-b579-95b19ed0934c")
	_, err := client.Client.Send(ctx, http.MethodHead, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the CheckMembershipExistence function
type CheckMembershipExistenceArgs struct {
	// (required) The group or user that is a child subject of the relationship.
	SubjectDescriptor *string
	// (required) The group that is the container in the relationship.
	ContainerDescriptor *string
}

// [Preview API] Create a new Azure DevOps group or materialize an existing AAD group.
func (client *ClientImpl) CreateGroup(ctx context.Context, args CreateGroupArgs) (*GraphGroup, error) {
	if args.CreationContext == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.CreationContext"}
	}
	queryParams := url.Values{}
	if args.ScopeDescriptor != nil {
		queryParams.Add("scopeDescriptor", *args.ScopeDescriptor)
	}
	if args.GroupDescriptors != nil {
		listAsString := strings.Join((*args.GroupDescriptors)[:], ",")
		queryParams.Add("groupDescriptors", listAsString)
	}
	body, marshalErr := json.Marshal(*args.CreationContext)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("ebbe6af8-0b91-4c13-8cf1-777c14858188")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0-preview.1", nil, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GraphGroup
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateGroup function
type CreateGroupArgs struct {
	// (required) The subset of the full graph group used to uniquely find the graph subject in an external provider.
	CreationContext *GraphGroupCreationContext
	// (optional) A descriptor referencing the scope (collection, project) in which the group should be created. If omitted, will be created in the scope of the enclosing account or organization. Valid only for VSTS groups.
	ScopeDescriptor *string
	// (optional) A comma separated list of descriptors referencing groups you want the graph group to join
	GroupDescriptors *[]string
}

// [Preview API] Materialize an existing AAD or MSA user into the VSTS account.
func (client *ClientImpl) CreateUser(ctx context.Context, args CreateUserArgs) (*GraphUser, error) {
	if args.CreationContext == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.CreationContext"}
	}
	queryParams := url.Values{}
	if args.GroupDescriptors != nil {
		listAsString := strings.Join((*args.GroupDescriptors)[:], ",")
		queryParams.Add("groupDescriptors", listAsString)
	}
	body, marshalErr := json.Marshal(*args.CreationContext)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("005e26ec-6b77-4e4f-a986-b3827bf241f5")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0-preview.1", nil, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GraphUser
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateUser function
type CreateUserArgs struct {
	// (required) The subset of the full graph user used to uniquely find the graph subject in an external provider.
	CreationContext *GraphUserCreationContext
	// (optional) A comma separated list of descriptors of groups you want the graph user to join
	GroupDescriptors *[]string
}

// [Preview API]
func (client *ClientImpl) DeleteAvatar(ctx context.Context, args DeleteAvatarArgs) error {
	routeValues := make(map[string]string)
	if args.SubjectDescriptor == nil || *args.SubjectDescriptor == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubjectDescriptor"}
	}
	routeValues["subjectDescriptor"] = *args.SubjectDescriptor

	locationId, _ := uuid.Parse("801eaf9c-0585-4be8-9cdb-b0efa074de91")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteAvatar function
type DeleteAvatarArgs struct {
	// (required)
	SubjectDescriptor *string
}

// [Preview API] Removes an Azure DevOps group from all of its parent groups.
func (client *ClientImpl) DeleteGroup(ctx context.Context, args DeleteGroupArgs) error {
	routeValues := make(map[string]string)
	if args.GroupDescriptor == nil || *args.GroupDescriptor == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.GroupDescriptor"}
	}
	routeValues["groupDescriptor"] = *args.GroupDescriptor

	locationId, _ := uuid.Parse("ebbe6af8-0b91-4c13-8cf1-777c14858188")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteGroup function
type DeleteGroupArgs struct {
	// (required) The descriptor of the group to delete.
	GroupDescriptor *string
}

// [Preview API] Disables a user.
func (client *ClientImpl) DeleteUser(ctx context.Context, args DeleteUserArgs) error {
	routeValues := make(map[string]string)
	if args.UserDescriptor == nil || *args.UserDescriptor == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.UserDescriptor"}
	}
	routeValues["userDescriptor"] = *args.UserDescriptor

	locationId, _ := uuid.Parse("005e26ec-6b77-4e4f-a986-b3827bf241f5")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteUser function
type DeleteUserArgs struct {
	// (required) The descriptor of the user to delete.
	UserDescriptor *string
}

// [Preview API]
func (client *ClientImpl) GetAvatar(ctx context.Context, args GetAvatarArgs) (*profile.Avatar, error) {
	routeValues := make(map[string]string)
	if args.SubjectDescriptor == nil || *args.SubjectDescriptor == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubjectDescriptor"}
	}
	routeValues["subjectDescriptor"] = *args.SubjectDescriptor

	queryParams := url.Values{}
	if args.Size != nil {
		queryParams.Add("size", string(*args.Size))
	}
	if args.Format != nil {
		queryParams.Add("format", *args.Format)
	}
	locationId, _ := uuid.Parse("801eaf9c-0585-4be8-9cdb-b0efa074de91")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue profile.Avatar
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetAvatar function
type GetAvatarArgs struct {
	// (required)
	SubjectDescriptor *string
	// (optional)
	Size *profile.AvatarSize
	// (optional)
	Format *string
}

// [Preview API] Resolve a storage key to a descriptor
func (client *ClientImpl) GetDescriptor(ctx context.Context, args GetDescriptorArgs) (*GraphDescriptorResult, error) {
	routeValues := make(map[string]string)
	if args.StorageKey == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.StorageKey"}
	}
	routeValues["storageKey"] = (*args.StorageKey).String()

	locationId, _ := uuid.Parse("048aee0a-7072-4cde-ab73-7af77b1e0b4e")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GraphDescriptorResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetDescriptor function
type GetDescriptorArgs struct {
	// (required) Storage key of the subject (user, group, scope, etc.) to resolve
	StorageKey *uuid.UUID
}

// [Preview API] Get a group by its descriptor.
func (client *ClientImpl) GetGroup(ctx context.Context, args GetGroupArgs) (*GraphGroup, error) {
	routeValues := make(map[string]string)
	if args.GroupDescriptor == nil || *args.GroupDescriptor == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.GroupDescriptor"}
	}
	routeValues["groupDescriptor"] = *args.GroupDescriptor

	locationId, _ := uuid.Parse("ebbe6af8-0b91-4c13-8cf1-777c14858188")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GraphGroup
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetGroup function
type GetGroupArgs struct {
	// (required) The descriptor of the desired graph group.
	GroupDescriptor *string
}

// [Preview API] Get a membership relationship between a container and subject.
func (client *ClientImpl) GetMembership(ctx context.Context, args GetMembershipArgs) (*GraphMembership, error) {
	routeValues := make(map[string]string)
	if args.SubjectDescriptor == nil || *args.SubjectDescriptor == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubjectDescriptor"}
	}
	routeValues["subjectDescriptor"] = *args.SubjectDescriptor
	if args.ContainerDescriptor == nil || *args.ContainerDescriptor == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ContainerDescriptor"}
	}
	routeValues["containerDescriptor"] = *args.ContainerDescriptor

	locationId, _ := uuid.Parse("3fd2e6ca-fb30-443a-b579-95b19ed0934c")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GraphMembership
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetMembership function
type GetMembershipArgs struct {
	// (required) A descriptor to the child subject in the relationship.
	SubjectDescriptor *string
	// (required) A descriptor to the container in the relationship.
	ContainerDescriptor *string
}

// [Preview API] Check whether a subject is active or inactive.
func (client *ClientImpl) GetMembershipState(ctx context.Context, args GetMembershipStateArgs) (*GraphMembershipState, error) {
	routeValues := make(map[string]string)
	if args.SubjectDescriptor == nil || *args.SubjectDescriptor == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubjectDescriptor"}
	}
	routeValues["subjectDescriptor"] = *args.SubjectDescriptor

	locationId, _ := uuid.Parse("1ffe5c94-1144-4191-907b-d0211cad36a8")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GraphMembershipState
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetMembershipState function
type GetMembershipStateArgs struct {
	// (required) Descriptor of the subject (user, group, scope, etc.) to check state of
	SubjectDescriptor *string
}

// [Preview API]
func (client *ClientImpl) GetProviderInfo(ctx context.Context, args GetProviderInfoArgs) (*GraphProviderInfo, error) {
	routeValues := make(map[string]string)
	if args.UserDescriptor == nil || *args.UserDescriptor == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.UserDescriptor"}
	}
	routeValues["userDescriptor"] = *args.UserDescriptor

	locationId, _ := uuid.Parse("1e377995-6fa2-4588-bd64-930186abdcfa")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GraphProviderInfo
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetProviderInfo function
type GetProviderInfoArgs struct {
	// (required)
	UserDescriptor *string
}

// [Preview API] Resolve a descriptor to a storage key.
func (client *ClientImpl) GetStorageKey(ctx context.Context, args GetStorageKeyArgs) (*GraphStorageKeyResult, error) {
	routeValues := make(map[string]string)
	if args.SubjectDescriptor == nil || *args.SubjectDescriptor == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubjectDescriptor"}
	}
	routeValues["subjectDescriptor"] = *args.SubjectDescriptor

	locationId, _ := uuid.Parse("eb85f8cc-f0f6-4264-a5b1-ffe2e4d4801f")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GraphStorageKeyResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetStorageKey function
type GetStorageKeyArgs struct {
	// (required)
	SubjectDescriptor *string
}

// [Preview API] Get a user by its descriptor.
func (client *ClientImpl) GetUser(ctx context.Context, args GetUserArgs) (*GraphUser, error) {
	routeValues := make(map[string]string)
	if args.UserDescriptor == nil || *args.UserDescriptor == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.UserDescriptor"}
	}
	routeValues["userDescriptor"] = *args.UserDescriptor

	locationId, _ := uuid.Parse("005e26ec-6b77-4e4f-a986-b3827bf241f5")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GraphUser
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetUser function
type GetUserArgs struct {
	// (required) The descriptor of the desired user.
	UserDescriptor *string
}

// [Preview API] Gets a list of all groups in the current scope (usually organization or account).
func (client *ClientImpl) ListGroups(ctx context.Context, args ListGroupsArgs) (*PagedGraphGroups, error) {
	queryParams := url.Values{}
	if args.ScopeDescriptor != nil {
		queryParams.Add("scopeDescriptor", *args.ScopeDescriptor)
	}
	if args.SubjectTypes != nil {
		listAsString := strings.Join((*args.SubjectTypes)[:], ",")
		queryParams.Add("subjectTypes", listAsString)
	}
	if args.ContinuationToken != nil {
		queryParams.Add("continuationToken", *args.ContinuationToken)
	}
	locationId, _ := uuid.Parse("ebbe6af8-0b91-4c13-8cf1-777c14858188")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseBodyValue []GraphGroup
	err = client.Client.UnmarshalCollectionBody(resp, &responseBodyValue)

	var responseValue *PagedGraphGroups
	if err == nil {
		responseValue = &PagedGraphGroups{
			GraphGroups:       &responseBodyValue,
			ContinuationToken: &[]string{resp.Header.Get("X-MS-ContinuationToken")},
		}
	}

	return responseValue, err
}

// Arguments for the ListGroups function
type ListGroupsArgs struct {
	// (optional) Specify a non-default scope (collection, project) to search for groups.
	ScopeDescriptor *string
	// (optional) A comma separated list of user subject subtypes to reduce the retrieved results, e.g. Microsoft.IdentityModel.Claims.ClaimsIdentity
	SubjectTypes *[]string
	// (optional) An opaque data blob that allows the next page of data to resume immediately after where the previous page ended. The only reliable way to know if there is more data left is the presence of a continuation token.
	ContinuationToken *string
}

// [Preview API] Get all the memberships where this descriptor is a member in the relationship.
func (client *ClientImpl) ListMemberships(ctx context.Context, args ListMembershipsArgs) (*[]GraphMembership, error) {
	routeValues := make(map[string]string)
	if args.SubjectDescriptor == nil || *args.SubjectDescriptor == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubjectDescriptor"}
	}
	routeValues["subjectDescriptor"] = *args.SubjectDescriptor

	queryParams := url.Values{}
	if args.Direction != nil {
		queryParams.Add("direction", string(*args.Direction))
	}
	if args.Depth != nil {
		queryParams.Add("depth", strconv.Itoa(*args.Depth))
	}
	locationId, _ := uuid.Parse("e34b6394-6b30-4435-94a9-409a5eef3e31")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []GraphMembership
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ListMemberships function
type ListMembershipsArgs struct {
	// (required) Fetch all direct memberships of this descriptor.
	SubjectDescriptor *string
	// (optional) Defaults to Up.
	Direction *GraphTraversalDirection
	// (optional) The maximum number of edges to traverse up or down the membership tree. Currently the only supported value is '1'.
	Depth *int
}

// [Preview API] Get a list of all users in a given scope.
func (client *ClientImpl) ListUsers(ctx context.Context, args ListUsersArgs) (*PagedGraphUsers, error) {
	queryParams := url.Values{}
	if args.SubjectTypes != nil {
		listAsString := strings.Join((*args.SubjectTypes)[:], ",")
		queryParams.Add("subjectTypes", listAsString)
	}
	if args.ContinuationToken != nil {
		queryParams.Add("continuationToken", *args.ContinuationToken)
	}
	if args.ScopeDescriptor != nil {
		queryParams.Add("scopeDescriptor", *args.ScopeDescriptor)
	}
	locationId, _ := uuid.Parse("005e26ec-6b77-4e4f-a986-b3827bf241f5")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseBodyValue []GraphUser
	err = client.Client.UnmarshalCollectionBody(resp, &responseBodyValue)

	var responseValue *PagedGraphUsers
	if err == nil {
		responseValue = &PagedGraphUsers{
			GraphUsers:        &responseBodyValue,
			ContinuationToken: &[]string{resp.Header.Get("X-MS-ContinuationToken")},
		}
	}

	return responseValue, err
}

// Arguments for the ListUsers function
type ListUsersArgs struct {
	// (optional) A comma separated list of user subject subtypes to reduce the retrieved results, e.g. msa’, ‘aad’, ‘svc’ (service identity), ‘imp’ (imported identity), etc.
	SubjectTypes *[]string
	// (optional) An opaque data blob that allows the next page of data to resume immediately after where the previous page ended. The only reliable way to know if there is more data left is the presence of a continuation token.
	ContinuationToken *string
	// (optional) Specify a non-default scope (collection, project) to search for users.
	ScopeDescriptor *string
}

// [Preview API] Resolve descriptors to users, groups or scopes (Subjects) in a batch.
func (client *ClientImpl) LookupSubjects(ctx context.Context, args LookupSubjectsArgs) (*map[string]GraphSubject, error) {
	if args.SubjectLookup == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SubjectLookup"}
	}
	body, marshalErr := json.Marshal(*args.SubjectLookup)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("4dd4d168-11f2-48c4-83e8-756fa0de027c")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue map[string]GraphSubject
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the LookupSubjects function
type LookupSubjectsArgs struct {
	// (required) A list of descriptors that specifies a subset of subjects to retrieve. Each descriptor uniquely identifies the subject across all instance scopes, but only at a single point in time.
	SubjectLookup *GraphSubjectLookup
}

// [Preview API] Search for Azure Devops users, or/and groups. Results will be returned in a batch with no more than 100 graph subjects.
func (client *ClientImpl) QuerySubjects(ctx context.Context, args QuerySubjectsArgs) (*[]GraphSubject, error) {
	if args.SubjectQuery == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SubjectQuery"}
	}
	body, marshalErr := json.Marshal(*args.SubjectQuery)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("05942c89-006a-48ce-bb79-baeb8abf99c6")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []GraphSubject
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QuerySubjects function
type QuerySubjectsArgs struct {
	// (required) The query that we'll be using to search includes the following: Query: the search term. The search will be prefix matching only. SubjectKind: "User" or "Group" can be specified, both or either ScopeDescriptor: Non-default scope can be specified, i.e. project scope descriptor
	SubjectQuery *GraphSubjectQuery
}

// [Preview API] Deletes a membership between a container and subject.
func (client *ClientImpl) RemoveMembership(ctx context.Context, args RemoveMembershipArgs) error {
	routeValues := make(map[string]string)
	if args.SubjectDescriptor == nil || *args.SubjectDescriptor == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubjectDescriptor"}
	}
	routeValues["subjectDescriptor"] = *args.SubjectDescriptor
	if args.ContainerDescriptor == nil || *args.ContainerDescriptor == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ContainerDescriptor"}
	}
	routeValues["containerDescriptor"] = *args.ContainerDescriptor

	locationId, _ := uuid.Parse("3fd2e6ca-fb30-443a-b579-95b19ed0934c")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the RemoveMembership function
type RemoveMembershipArgs struct {
	// (required) A descriptor to a group or user that is the child subject in the relationship.
	SubjectDescriptor *string
	// (required) A descriptor to a group that is the container in the relationship.
	ContainerDescriptor *string
}

// [Preview API]
func (client *ClientImpl) RequestAccess(ctx context.Context, args RequestAccessArgs) error {
	if args.Jsondocument == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.Jsondocument"}
	}
	body, marshalErr := json.Marshal(args.Jsondocument)
	if marshalErr != nil {
		return marshalErr
	}
	locationId, _ := uuid.Parse("8d54bf92-8c99-47f2-9972-b21341f1722e")
	_, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the RequestAccess function
type RequestAccessArgs struct {
	// (required)
	Jsondocument interface{}
}

// [Preview API]
func (client *ClientImpl) SetAvatar(ctx context.Context, args SetAvatarArgs) error {
	if args.Avatar == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.Avatar"}
	}
	routeValues := make(map[string]string)
	if args.SubjectDescriptor == nil || *args.SubjectDescriptor == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubjectDescriptor"}
	}
	routeValues["subjectDescriptor"] = *args.SubjectDescriptor

	body, marshalErr := json.Marshal(*args.Avatar)
	if marshalErr != nil {
		return marshalErr
	}
	locationId, _ := uuid.Parse("801eaf9c-0585-4be8-9cdb-b0efa074de91")
	_, err := client.Client.Send(ctx, http.MethodPut, locationId, "6.0-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the SetAvatar function
type SetAvatarArgs struct {
	// (required)
	Avatar *profile.Avatar
	// (required)
	SubjectDescriptor *string
}

// [Preview API] Update the properties of an Azure DevOps group.
func (client *ClientImpl) UpdateGroup(ctx context.Context, args UpdateGroupArgs) (*GraphGroup, error) {
	if args.PatchDocument == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PatchDocument"}
	}
	routeValues := make(map[string]string)
	if args.GroupDescriptor == nil || *args.GroupDescriptor == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.GroupDescriptor"}
	}
	routeValues["groupDescriptor"] = *args.GroupDescriptor

	body, marshalErr := json.Marshal(*args.PatchDocument)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("ebbe6af8-0b91-4c13-8cf1-777c14858188")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "6.0-preview.1", routeValues, nil, bytes.NewReader(body), "application/json-patch+json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GraphGroup
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateGroup function
type UpdateGroupArgs struct {
	// (required) The descriptor of the group to modify.
	GroupDescriptor *string
	// (required) The JSON+Patch document containing the fields to alter.
	PatchDocument *[]webapi.JsonPatchOperation
}

// [Preview API] Map an existing user to a different identity
func (client *ClientImpl) UpdateUser(ctx context.Context, args UpdateUserArgs) (*GraphUser, error) {
	if args.UpdateContext == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UpdateContext"}
	}
	routeValues := make(map[string]string)
	if args.UserDescriptor == nil || *args.UserDescriptor == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.UserDescriptor"}
	}
	routeValues["userDescriptor"] = *args.UserDescriptor

	body, marshalErr := json.Marshal(*args.UpdateContext)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("005e26ec-6b77-4e4f-a986-b3827bf241f5")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "6.0-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GraphUser
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateUser function
type UpdateUserArgs struct {
	// (required) The subset of the full graph user used to uniquely find the graph subject in an external provider.
	UpdateContext *GraphUserUpdateContext
	// (required) the descriptor of the user to update
	UserDescriptor *string
}

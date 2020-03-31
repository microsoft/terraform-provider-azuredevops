// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package memberentitlementmanagement

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/licensingrule"
	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
	"net/http"
	"net/url"
	"strconv"
)

var ResourceAreaId, _ = uuid.Parse("68ddce18-2501-45f1-a17b-7931a9922690")

type Client interface {
	// [Preview API] Create a group entitlement with license rule, extension rule.
	AddGroupEntitlement(context.Context, AddGroupEntitlementArgs) (*GroupEntitlementOperationReference, error)
	// [Preview API] Add a member to a Group.
	AddMemberToGroup(context.Context, AddMemberToGroupArgs) error
	// [Preview API] Add a user, assign license and extensions and make them a member of a project group in an account.
	AddUserEntitlement(context.Context, AddUserEntitlementArgs) (*UserEntitlementsPostResponse, error)
	// [Preview API] Delete a group entitlement.
	DeleteGroupEntitlement(context.Context, DeleteGroupEntitlementArgs) (*GroupEntitlementOperationReference, error)
	// [Preview API] Delete a user from the account.
	DeleteUserEntitlement(context.Context, DeleteUserEntitlementArgs) error
	// [Preview API] Get a group entitlement.
	GetGroupEntitlement(context.Context, GetGroupEntitlementArgs) (*GroupEntitlement, error)
	// [Preview API] Get the group entitlements for an account.
	GetGroupEntitlements(context.Context, GetGroupEntitlementsArgs) (*[]GroupEntitlement, error)
	// [Preview API] Get direct members of a Group.
	GetGroupMembers(context.Context, GetGroupMembersArgs) (*PagedGraphMemberList, error)
	// [Preview API] Get User Entitlement for a user.
	GetUserEntitlement(context.Context, GetUserEntitlementArgs) (*UserEntitlement, error)
	// [Preview API] Get a paged set of user entitlements matching the filter criteria. If no filter is is passed, a page from all the account users is returned.
	GetUserEntitlements(context.Context, GetUserEntitlementsArgs) (*PagedGraphMemberList, error)
	// [Preview API] Get summary of Licenses, Extension, Projects, Groups and their assignments in the collection.
	GetUsersSummary(context.Context, GetUsersSummaryArgs) (*UsersSummary, error)
	// [Preview API] Remove a member from a Group.
	RemoveMemberFromGroup(context.Context, RemoveMemberFromGroupArgs) error
	// [Preview API] Update entitlements (License Rule, Extensions Rule, Project memberships etc.) for a group.
	UpdateGroupEntitlement(context.Context, UpdateGroupEntitlementArgs) (*GroupEntitlementOperationReference, error)
	// [Preview API] Edit the entitlements (License, Extensions, Projects, Teams etc) for a user.
	UpdateUserEntitlement(context.Context, UpdateUserEntitlementArgs) (*UserEntitlementsPatchResponse, error)
	// [Preview API] Edit the entitlements (License, Extensions, Projects, Teams etc) for one or more users.
	UpdateUserEntitlements(context.Context, UpdateUserEntitlementsArgs) (*UserEntitlementOperationReference, error)
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

// [Preview API] Create a group entitlement with license rule, extension rule.
func (client *ClientImpl) AddGroupEntitlement(ctx context.Context, args AddGroupEntitlementArgs) (*GroupEntitlementOperationReference, error) {
	if args.GroupEntitlement == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.GroupEntitlement"}
	}
	queryParams := url.Values{}
	if args.RuleOption != nil {
		queryParams.Add("ruleOption", string(*args.RuleOption))
	}
	body, marshalErr := json.Marshal(*args.GroupEntitlement)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("2280bffa-58a2-49da-822e-0764a1bb44f7")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "5.1-preview.1", nil, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GroupEntitlementOperationReference
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the AddGroupEntitlement function
type AddGroupEntitlementArgs struct {
	// (required) GroupEntitlement object specifying License Rule, Extensions Rule for the group. Based on the rules the members of the group will be given licenses and extensions. The Group Entitlement can be used to add the group to another project level groups
	GroupEntitlement *GroupEntitlement
	// (optional) RuleOption [ApplyGroupRule/TestApplyGroupRule] - specifies if the rules defined in group entitlement should be created and applied to it’s members (default option) or just be tested
	RuleOption *licensingrule.RuleOption
}

// [Preview API] Add a member to a Group.
func (client *ClientImpl) AddMemberToGroup(ctx context.Context, args AddMemberToGroupArgs) error {
	routeValues := make(map[string]string)
	if args.GroupId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.GroupId"}
	}
	routeValues["groupId"] = (*args.GroupId).String()
	if args.MemberId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.MemberId"}
	}
	routeValues["memberId"] = (*args.MemberId).String()

	locationId, _ := uuid.Parse("45a36e53-5286-4518-aa72-2d29f7acc5d8")
	_, err := client.Client.Send(ctx, http.MethodPut, locationId, "5.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the AddMemberToGroup function
type AddMemberToGroupArgs struct {
	// (required) Id of the Group.
	GroupId *uuid.UUID
	// (required) Id of the member to add.
	MemberId *uuid.UUID
}

// [Preview API] Add a user, assign license and extensions and make them a member of a project group in an account.
func (client *ClientImpl) AddUserEntitlement(ctx context.Context, args AddUserEntitlementArgs) (*UserEntitlementsPostResponse, error) {
	if args.UserEntitlement == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UserEntitlement"}
	}
	body, marshalErr := json.Marshal(*args.UserEntitlement)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("387f832c-dbf2-4643-88e9-c1aa94dbb737")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "5.1-preview.2", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue UserEntitlementsPostResponse
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the AddUserEntitlement function
type AddUserEntitlementArgs struct {
	// (required) UserEntitlement object specifying License, Extensions and Project/Team groups the user should be added to.
	UserEntitlement *UserEntitlement
}

// [Preview API] Delete a group entitlement.
func (client *ClientImpl) DeleteGroupEntitlement(ctx context.Context, args DeleteGroupEntitlementArgs) (*GroupEntitlementOperationReference, error) {
	routeValues := make(map[string]string)
	if args.GroupId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.GroupId"}
	}
	routeValues["groupId"] = (*args.GroupId).String()

	queryParams := url.Values{}
	if args.RuleOption != nil {
		queryParams.Add("ruleOption", string(*args.RuleOption))
	}
	if args.RemoveGroupMembership != nil {
		queryParams.Add("removeGroupMembership", strconv.FormatBool(*args.RemoveGroupMembership))
	}
	locationId, _ := uuid.Parse("2280bffa-58a2-49da-822e-0764a1bb44f7")
	resp, err := client.Client.Send(ctx, http.MethodDelete, locationId, "5.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GroupEntitlementOperationReference
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the DeleteGroupEntitlement function
type DeleteGroupEntitlementArgs struct {
	// (required) ID of the group to delete.
	GroupId *uuid.UUID
	// (optional) RuleOption [ApplyGroupRule/TestApplyGroupRule] - specifies if the rules defined in group entitlement should be deleted and the changes are applied to it’s members (default option) or just be tested
	RuleOption *licensingrule.RuleOption
	// (optional) Optional parameter that specifies whether the group with the given ID should be removed from all other groups
	RemoveGroupMembership *bool
}

// [Preview API] Delete a user from the account.
func (client *ClientImpl) DeleteUserEntitlement(ctx context.Context, args DeleteUserEntitlementArgs) error {
	routeValues := make(map[string]string)
	if args.UserId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.UserId"}
	}
	routeValues["userId"] = (*args.UserId).String()

	locationId, _ := uuid.Parse("8480c6eb-ce60-47e9-88df-eca3c801638b")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "5.1-preview.2", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteUserEntitlement function
type DeleteUserEntitlementArgs struct {
	// (required) ID of the user.
	UserId *uuid.UUID
}

// [Preview API] Get a group entitlement.
func (client *ClientImpl) GetGroupEntitlement(ctx context.Context, args GetGroupEntitlementArgs) (*GroupEntitlement, error) {
	routeValues := make(map[string]string)
	if args.GroupId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.GroupId"}
	}
	routeValues["groupId"] = (*args.GroupId).String()

	locationId, _ := uuid.Parse("2280bffa-58a2-49da-822e-0764a1bb44f7")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GroupEntitlement
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetGroupEntitlement function
type GetGroupEntitlementArgs struct {
	// (required) ID of the group.
	GroupId *uuid.UUID
}

// [Preview API] Get the group entitlements for an account.
func (client *ClientImpl) GetGroupEntitlements(ctx context.Context, args GetGroupEntitlementsArgs) (*[]GroupEntitlement, error) {
	locationId, _ := uuid.Parse("2280bffa-58a2-49da-822e-0764a1bb44f7")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1-preview.1", nil, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []GroupEntitlement
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetGroupEntitlements function
type GetGroupEntitlementsArgs struct {
}

// [Preview API] Get direct members of a Group.
func (client *ClientImpl) GetGroupMembers(ctx context.Context, args GetGroupMembersArgs) (*PagedGraphMemberList, error) {
	routeValues := make(map[string]string)
	if args.GroupId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.GroupId"}
	}
	routeValues["groupId"] = (*args.GroupId).String()

	queryParams := url.Values{}
	if args.MaxResults != nil {
		queryParams.Add("maxResults", strconv.Itoa(*args.MaxResults))
	}
	if args.PagingToken != nil {
		queryParams.Add("pagingToken", *args.PagingToken)
	}
	locationId, _ := uuid.Parse("45a36e53-5286-4518-aa72-2d29f7acc5d8")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PagedGraphMemberList
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetGroupMembers function
type GetGroupMembersArgs struct {
	// (required) Id of the Group.
	GroupId *uuid.UUID
	// (optional) Maximum number of results to retrieve.
	MaxResults *int
	// (optional) Paging Token from the previous page fetched. If the 'pagingToken' is null, the results would be fetched from the beginning of the Members List.
	PagingToken *string
}

// [Preview API] Get User Entitlement for a user.
func (client *ClientImpl) GetUserEntitlement(ctx context.Context, args GetUserEntitlementArgs) (*UserEntitlement, error) {
	routeValues := make(map[string]string)
	if args.UserId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UserId"}
	}
	routeValues["userId"] = (*args.UserId).String()

	locationId, _ := uuid.Parse("8480c6eb-ce60-47e9-88df-eca3c801638b")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1-preview.2", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue UserEntitlement
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetUserEntitlement function
type GetUserEntitlementArgs struct {
	// (required) ID of the user.
	UserId *uuid.UUID
}

// [Preview API] Get a paged set of user entitlements matching the filter criteria. If no filter is is passed, a page from all the account users is returned.
func (client *ClientImpl) GetUserEntitlements(ctx context.Context, args GetUserEntitlementsArgs) (*PagedGraphMemberList, error) {
	queryParams := url.Values{}
	if args.Top != nil {
		queryParams.Add("top", strconv.Itoa(*args.Top))
	}
	if args.Skip != nil {
		queryParams.Add("skip", strconv.Itoa(*args.Skip))
	}
	if args.Filter != nil {
		queryParams.Add("filter", *args.Filter)
	}
	if args.SortOption != nil {
		queryParams.Add("sortOption", *args.SortOption)
	}
	locationId, _ := uuid.Parse("387f832c-dbf2-4643-88e9-c1aa94dbb737")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1-preview.2", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PagedGraphMemberList
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetUserEntitlements function
type GetUserEntitlementsArgs struct {
	// (optional) Maximum number of the user entitlements to return. Max value is 10000. Default value is 100
	Top *int
	// (optional) Offset: Number of records to skip. Default value is 0
	Skip *int
	// (optional) Comma (",") separated list of properties and their values to filter on. Currently, the API only supports filtering by ExtensionId. An example parameter would be filter=extensionId eq search.
	Filter *string
	// (optional) PropertyName and Order (separated by a space ( )) to sort on (e.g. LastAccessDate Desc)
	SortOption *string
}

// [Preview API] Get summary of Licenses, Extension, Projects, Groups and their assignments in the collection.
func (client *ClientImpl) GetUsersSummary(ctx context.Context, args GetUsersSummaryArgs) (*UsersSummary, error) {
	queryParams := url.Values{}
	if args.Select != nil {
		queryParams.Add("select", *args.Select)
	}
	locationId, _ := uuid.Parse("5ae55b13-c9dd-49d1-957e-6e76c152e3d9")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1-preview.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue UsersSummary
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetUsersSummary function
type GetUsersSummaryArgs struct {
	// (optional) Comma (",") separated list of properties to select. Supported property names are {AccessLevels, Licenses, Extensions, Projects, Groups}.
	Select *string
}

// [Preview API] Remove a member from a Group.
func (client *ClientImpl) RemoveMemberFromGroup(ctx context.Context, args RemoveMemberFromGroupArgs) error {
	routeValues := make(map[string]string)
	if args.GroupId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.GroupId"}
	}
	routeValues["groupId"] = (*args.GroupId).String()
	if args.MemberId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.MemberId"}
	}
	routeValues["memberId"] = (*args.MemberId).String()

	locationId, _ := uuid.Parse("45a36e53-5286-4518-aa72-2d29f7acc5d8")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "5.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the RemoveMemberFromGroup function
type RemoveMemberFromGroupArgs struct {
	// (required) Id of the group.
	GroupId *uuid.UUID
	// (required) Id of the member to remove.
	MemberId *uuid.UUID
}

// [Preview API] Update entitlements (License Rule, Extensions Rule, Project memberships etc.) for a group.
func (client *ClientImpl) UpdateGroupEntitlement(ctx context.Context, args UpdateGroupEntitlementArgs) (*GroupEntitlementOperationReference, error) {
	if args.Document == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Document"}
	}
	routeValues := make(map[string]string)
	if args.GroupId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.GroupId"}
	}
	routeValues["groupId"] = (*args.GroupId).String()

	queryParams := url.Values{}
	if args.RuleOption != nil {
		queryParams.Add("ruleOption", string(*args.RuleOption))
	}
	body, marshalErr := json.Marshal(*args.Document)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("2280bffa-58a2-49da-822e-0764a1bb44f7")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "5.1-preview.1", routeValues, queryParams, bytes.NewReader(body), "application/json-patch+json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GroupEntitlementOperationReference
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateGroupEntitlement function
type UpdateGroupEntitlementArgs struct {
	// (required) JsonPatchDocument containing the operations to perform on the group.
	Document *[]webapi.JsonPatchOperation
	// (required) ID of the group.
	GroupId *uuid.UUID
	// (optional) RuleOption [ApplyGroupRule/TestApplyGroupRule] - specifies if the rules defined in group entitlement should be updated and the changes are applied to it’s members (default option) or just be tested
	RuleOption *licensingrule.RuleOption
}

// [Preview API] Edit the entitlements (License, Extensions, Projects, Teams etc) for a user.
func (client *ClientImpl) UpdateUserEntitlement(ctx context.Context, args UpdateUserEntitlementArgs) (*UserEntitlementsPatchResponse, error) {
	if args.Document == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Document"}
	}
	routeValues := make(map[string]string)
	if args.UserId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UserId"}
	}
	routeValues["userId"] = (*args.UserId).String()

	body, marshalErr := json.Marshal(*args.Document)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("8480c6eb-ce60-47e9-88df-eca3c801638b")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "5.1-preview.2", routeValues, nil, bytes.NewReader(body), "application/json-patch+json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue UserEntitlementsPatchResponse
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateUserEntitlement function
type UpdateUserEntitlementArgs struct {
	// (required) JsonPatchDocument containing the operations to perform on the user.
	Document *[]webapi.JsonPatchOperation
	// (required) ID of the user.
	UserId *uuid.UUID
}

// [Preview API] Edit the entitlements (License, Extensions, Projects, Teams etc) for one or more users.
func (client *ClientImpl) UpdateUserEntitlements(ctx context.Context, args UpdateUserEntitlementsArgs) (*UserEntitlementOperationReference, error) {
	if args.Document == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Document"}
	}
	queryParams := url.Values{}
	if args.DoNotSendInviteForNewUsers != nil {
		queryParams.Add("doNotSendInviteForNewUsers", strconv.FormatBool(*args.DoNotSendInviteForNewUsers))
	}
	body, marshalErr := json.Marshal(*args.Document)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("387f832c-dbf2-4643-88e9-c1aa94dbb737")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "5.1-preview.2", nil, queryParams, bytes.NewReader(body), "application/json-patch+json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue UserEntitlementOperationReference
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateUserEntitlements function
type UpdateUserEntitlementsArgs struct {
	// (required) JsonPatchDocument containing the operations to perform.
	Document *[]webapi.JsonPatchOperation
	// (optional) Whether to send email invites to new users or not
	DoNotSendInviteForNewUsers *bool
}

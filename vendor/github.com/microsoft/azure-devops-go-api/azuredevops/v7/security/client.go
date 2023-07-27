// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package security

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"net/http"
	"net/url"
	"strconv"
)

type Client interface {
	// [Preview API] Evaluates whether the caller has the specified permissions on the specified set of security tokens.
	HasPermissions(context.Context, HasPermissionsArgs) (*[]bool, error)
	// [Preview API] Evaluates multiple permissions for the calling user.  Note: This method does not aggregate the results, nor does it short-circuit if one of the permissions evaluates to false.
	HasPermissionsBatch(context.Context, HasPermissionsBatchArgs) (*PermissionEvaluationBatch, error)
	// [Preview API] Return a list of access control lists for the specified security namespace and token. All ACLs in the security namespace will be retrieved if no optional parameters are provided.
	QueryAccessControlLists(context.Context, QueryAccessControlListsArgs) (*[]AccessControlList, error)
	// [Preview API] List all security namespaces or just the specified namespace.
	QuerySecurityNamespaces(context.Context, QuerySecurityNamespacesArgs) (*[]SecurityNamespaceDescription, error)
	// [Preview API] Remove the specified ACEs from the ACL belonging to the specified token.
	RemoveAccessControlEntries(context.Context, RemoveAccessControlEntriesArgs) (*bool, error)
	// [Preview API] Remove access control lists under the specfied security namespace.
	RemoveAccessControlLists(context.Context, RemoveAccessControlListsArgs) (*bool, error)
	// [Preview API] Removes the specified permissions on a security token for a user or group.
	RemovePermission(context.Context, RemovePermissionArgs) (*AccessControlEntry, error)
	// [Preview API] Add or update ACEs in the ACL for the provided token. The request body contains the target token, a list of [ACEs](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/access%20control%20entries/set%20access%20control%20entries?#accesscontrolentry) and a optional merge parameter. In the case of a collision (by identity descriptor) with an existing ACE in the ACL, the "merge" parameter determines the behavior. If set, the existing ACE has its allow and deny merged with the incoming ACE's allow and deny. If unset, the existing ACE is displaced.
	SetAccessControlEntries(context.Context, SetAccessControlEntriesArgs) (*[]AccessControlEntry, error)
	// [Preview API] Create or update one or more access control lists. All data that currently exists for the ACLs supplied will be overwritten.
	SetAccessControlLists(context.Context, SetAccessControlListsArgs) error
}

type ClientImpl struct {
	Client azuredevops.Client
}

func NewClient(ctx context.Context, connection *azuredevops.Connection) Client {
	client := connection.GetClientByUrl(connection.BaseUrl)
	return &ClientImpl{
		Client: *client,
	}
}

// [Preview API] Evaluates whether the caller has the specified permissions on the specified set of security tokens.
func (client *ClientImpl) HasPermissions(ctx context.Context, args HasPermissionsArgs) (*[]bool, error) {
	routeValues := make(map[string]string)
	if args.SecurityNamespaceId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SecurityNamespaceId"}
	}
	routeValues["securityNamespaceId"] = (*args.SecurityNamespaceId).String()
	if args.Permissions != nil {
		routeValues["permissions"] = strconv.Itoa(*args.Permissions)
	}

	queryParams := url.Values{}
	if args.Tokens != nil {
		queryParams.Add("tokens", *args.Tokens)
	}
	if args.AlwaysAllowAdministrators != nil {
		queryParams.Add("alwaysAllowAdministrators", strconv.FormatBool(*args.AlwaysAllowAdministrators))
	}
	if args.Delimiter != nil {
		queryParams.Add("delimiter", *args.Delimiter)
	}
	locationId, _ := uuid.Parse("dd3b8bd6-c7fc-4cbd-929a-933d9c011c9d")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.2", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []bool
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the HasPermissions function
type HasPermissionsArgs struct {
	// (required) Security namespace identifier.
	SecurityNamespaceId *uuid.UUID
	// (optional) Permissions to evaluate.
	Permissions *int
	// (optional) One or more security tokens to evaluate.
	Tokens *string
	// (optional) If true and if the caller is an administrator, always return true.
	AlwaysAllowAdministrators *bool
	// (optional) Optional security token separator. Defaults to ",".
	Delimiter *string
}

// [Preview API] Evaluates multiple permissions for the calling user.  Note: This method does not aggregate the results, nor does it short-circuit if one of the permissions evaluates to false.
func (client *ClientImpl) HasPermissionsBatch(ctx context.Context, args HasPermissionsBatchArgs) (*PermissionEvaluationBatch, error) {
	if args.EvalBatch == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.EvalBatch"}
	}
	body, marshalErr := json.Marshal(*args.EvalBatch)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("cf1faa59-1b63-4448-bf04-13d981a46f5d")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PermissionEvaluationBatch
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the HasPermissionsBatch function
type HasPermissionsBatchArgs struct {
	// (required) The set of evaluation requests.
	EvalBatch *PermissionEvaluationBatch
}

// [Preview API] Return a list of access control lists for the specified security namespace and token. All ACLs in the security namespace will be retrieved if no optional parameters are provided.
func (client *ClientImpl) QueryAccessControlLists(ctx context.Context, args QueryAccessControlListsArgs) (*[]AccessControlList, error) {
	routeValues := make(map[string]string)
	if args.SecurityNamespaceId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SecurityNamespaceId"}
	}
	routeValues["securityNamespaceId"] = (*args.SecurityNamespaceId).String()

	queryParams := url.Values{}
	if args.Token != nil {
		queryParams.Add("token", *args.Token)
	}
	if args.Descriptors != nil {
		queryParams.Add("descriptors", *args.Descriptors)
	}
	if args.IncludeExtendedInfo != nil {
		queryParams.Add("includeExtendedInfo", strconv.FormatBool(*args.IncludeExtendedInfo))
	}
	if args.Recurse != nil {
		queryParams.Add("recurse", strconv.FormatBool(*args.Recurse))
	}
	locationId, _ := uuid.Parse("18a2ad18-7571-46ae-bec7-0c7da1495885")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []AccessControlList
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryAccessControlLists function
type QueryAccessControlListsArgs struct {
	// (required) Security namespace identifier.
	SecurityNamespaceId *uuid.UUID
	// (optional) Security token
	Token *string
	// (optional) An optional filter string containing a list of identity descriptors separated by ',' whose ACEs should be retrieved. If this is left null, entire ACLs will be returned.
	Descriptors *string
	// (optional) If true, populate the extended information properties for the access control entries contained in the returned lists.
	IncludeExtendedInfo *bool
	// (optional) If true and this is a hierarchical namespace, return child ACLs of the specified token.
	Recurse *bool
}

// [Preview API] List all security namespaces or just the specified namespace.
func (client *ClientImpl) QuerySecurityNamespaces(ctx context.Context, args QuerySecurityNamespacesArgs) (*[]SecurityNamespaceDescription, error) {
	routeValues := make(map[string]string)
	if args.SecurityNamespaceId != nil {
		routeValues["securityNamespaceId"] = (*args.SecurityNamespaceId).String()
	}

	queryParams := url.Values{}
	if args.LocalOnly != nil {
		queryParams.Add("localOnly", strconv.FormatBool(*args.LocalOnly))
	}
	locationId, _ := uuid.Parse("ce7b9f95-fde9-4be8-a86d-83b366f0b87a")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []SecurityNamespaceDescription
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QuerySecurityNamespaces function
type QuerySecurityNamespacesArgs struct {
	// (optional) Security namespace identifier.
	SecurityNamespaceId *uuid.UUID
	// (optional) If true, retrieve only local security namespaces.
	LocalOnly *bool
}

// [Preview API] Remove the specified ACEs from the ACL belonging to the specified token.
func (client *ClientImpl) RemoveAccessControlEntries(ctx context.Context, args RemoveAccessControlEntriesArgs) (*bool, error) {
	routeValues := make(map[string]string)
	if args.SecurityNamespaceId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SecurityNamespaceId"}
	}
	routeValues["securityNamespaceId"] = (*args.SecurityNamespaceId).String()

	queryParams := url.Values{}
	if args.Token != nil {
		queryParams.Add("token", *args.Token)
	}
	if args.Descriptors != nil {
		queryParams.Add("descriptors", *args.Descriptors)
	}
	locationId, _ := uuid.Parse("ac08c8ff-4323-4b08-af90-bcd018d380ce")
	resp, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue bool
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the RemoveAccessControlEntries function
type RemoveAccessControlEntriesArgs struct {
	// (required) Security namespace identifier.
	SecurityNamespaceId *uuid.UUID
	// (optional) The token whose ACL should be modified.
	Token *string
	// (optional) String containing a list of identity descriptors separated by ',' whose entries should be removed.
	Descriptors *string
}

// [Preview API] Remove access control lists under the specfied security namespace.
func (client *ClientImpl) RemoveAccessControlLists(ctx context.Context, args RemoveAccessControlListsArgs) (*bool, error) {
	routeValues := make(map[string]string)
	if args.SecurityNamespaceId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SecurityNamespaceId"}
	}
	routeValues["securityNamespaceId"] = (*args.SecurityNamespaceId).String()

	queryParams := url.Values{}
	if args.Tokens != nil {
		queryParams.Add("tokens", *args.Tokens)
	}
	if args.Recurse != nil {
		queryParams.Add("recurse", strconv.FormatBool(*args.Recurse))
	}
	locationId, _ := uuid.Parse("18a2ad18-7571-46ae-bec7-0c7da1495885")
	resp, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue bool
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the RemoveAccessControlLists function
type RemoveAccessControlListsArgs struct {
	// (required) Security namespace identifier.
	SecurityNamespaceId *uuid.UUID
	// (optional) One or more comma-separated security tokens
	Tokens *string
	// (optional) If true and this is a hierarchical namespace, also remove child ACLs of the specified tokens.
	Recurse *bool
}

// [Preview API] Removes the specified permissions on a security token for a user or group.
func (client *ClientImpl) RemovePermission(ctx context.Context, args RemovePermissionArgs) (*AccessControlEntry, error) {
	routeValues := make(map[string]string)
	if args.SecurityNamespaceId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SecurityNamespaceId"}
	}
	routeValues["securityNamespaceId"] = (*args.SecurityNamespaceId).String()
	if args.Permissions != nil {
		routeValues["permissions"] = strconv.Itoa(*args.Permissions)
	}

	queryParams := url.Values{}
	if args.Descriptor == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "descriptor"}
	}
	queryParams.Add("descriptor", *args.Descriptor)
	if args.Token != nil {
		queryParams.Add("token", *args.Token)
	}
	locationId, _ := uuid.Parse("dd3b8bd6-c7fc-4cbd-929a-933d9c011c9d")
	resp, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.2", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue AccessControlEntry
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the RemovePermission function
type RemovePermissionArgs struct {
	// (required) Security namespace identifier.
	SecurityNamespaceId *uuid.UUID
	// (required) Identity descriptor of the user to remove permissions for.
	Descriptor *string
	// (optional) Permissions to remove.
	Permissions *int
	// (optional) Security token to remove permissions for.
	Token *string
}

// [Preview API] Add or update ACEs in the ACL for the provided token. The request body contains the target token, a list of [ACEs](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/access%20control%20entries/set%20access%20control%20entries?#accesscontrolentry) and a optional merge parameter. In the case of a collision (by identity descriptor) with an existing ACE in the ACL, the "merge" parameter determines the behavior. If set, the existing ACE has its allow and deny merged with the incoming ACE's allow and deny. If unset, the existing ACE is displaced.
func (client *ClientImpl) SetAccessControlEntries(ctx context.Context, args SetAccessControlEntriesArgs) (*[]AccessControlEntry, error) {
	if args.Container == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Container"}
	}
	routeValues := make(map[string]string)
	if args.SecurityNamespaceId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SecurityNamespaceId"}
	}
	routeValues["securityNamespaceId"] = (*args.SecurityNamespaceId).String()

	body, marshalErr := json.Marshal(args.Container)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("ac08c8ff-4323-4b08-af90-bcd018d380ce")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []AccessControlEntry
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the SetAccessControlEntries function
type SetAccessControlEntriesArgs struct {
	// (required)
	Container interface{}
	// (required) Security namespace identifier.
	SecurityNamespaceId *uuid.UUID
}

// [Preview API] Create or update one or more access control lists. All data that currently exists for the ACLs supplied will be overwritten.
func (client *ClientImpl) SetAccessControlLists(ctx context.Context, args SetAccessControlListsArgs) error {
	if args.AccessControlLists == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.AccessControlLists"}
	}
	routeValues := make(map[string]string)
	if args.SecurityNamespaceId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.SecurityNamespaceId"}
	}
	routeValues["securityNamespaceId"] = (*args.SecurityNamespaceId).String()

	body, marshalErr := json.Marshal(*args.AccessControlLists)
	if marshalErr != nil {
		return marshalErr
	}
	locationId, _ := uuid.Parse("18a2ad18-7571-46ae-bec7-0c7da1495885")
	_, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the SetAccessControlLists function
type SetAccessControlListsArgs struct {
	// (required) A list of ACLs to create or update.
	AccessControlLists *azuredevops.VssJsonCollectionWrapper
	// (required) Security namespace identifier.
	SecurityNamespaceId *uuid.UUID
}

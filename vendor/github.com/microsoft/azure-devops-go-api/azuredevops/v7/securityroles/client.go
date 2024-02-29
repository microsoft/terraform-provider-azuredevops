// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package securityroles

import (
	"bytes"
	"encoding/json"
	"context"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
)

type Client interface {
	ListSecurityRoleDefinitions(ctx context.Context, args *ListSecurityRoleDefinitionsArgs) (*[]SecurityRoleDefinition, error)
	ListSecurityRoleAssignments(ctx context.Context, args *ListSecurityRoleAssignmentsArgs) (*[]SecurityRoleAssignment, error)
	GetSecurityRoleAssignment(ctx context.Context, args *GetSecurityRoleAssignmentArgs) (*SecurityRoleAssignment, error)
	SetSecurityRoleAssignment(ctx context.Context, args *SetSecurityRoleAssignmentArgs) error
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

// Arguments for the ListSecurityRoleDefinitions function
type ListSecurityRoleDefinitionsArgs struct {
	// (required) Scope identifier.
	Scope *string
}

func (client *ClientImpl) ListSecurityRoleDefinitions(ctx context.Context, args *ListSecurityRoleDefinitionsArgs) (*[]SecurityRoleDefinition, error) {
	routeValues := make(map[string]string)
	if args == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ScopeId"}
	}
	routeValues["scopeId"] = *args.Scope

	queryParams := url.Values{}

	locationId, _ := uuid.Parse("f4cc9a86-453c-48d2-b44d-d3bd5c105f4f")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []SecurityRoleDefinition
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ListSecurityRoleAssignmentsArgs function
type ListSecurityRoleAssignmentsArgs struct {
	// (required) Scope identifier.
	Scope *string
	ResourceId *string
}

func (client *ClientImpl) ListSecurityRoleAssignments(ctx context.Context, args *ListSecurityRoleAssignmentsArgs) (*[]SecurityRoleAssignment, error) {
	routeValues := make(map[string]string)
	if args == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ScopeId"}
	}

	resId := args.ResourceId

	routeValues["scopeId"] = *args.Scope
	routeValues["resource"] = "roleassignments"
	routeValues["resourceId"] = *resId

	queryParams := url.Values{}

	locationId, _ := uuid.Parse("9461c234-c84c-4ed2-b918-2f0f92ad0a35")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []SecurityRoleAssignment
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)

	return &responseValue, err
}

// Arguments for the GetSecurityRoleAssignmentArgs function
type GetSecurityRoleAssignmentArgs struct {
	// (required) Scope identifier.
	Scope *string
	ResourceId *string
	IdentityId *uuid.UUID
}

func (client *ClientImpl) GetSecurityRoleAssignment(ctx context.Context, args *GetSecurityRoleAssignmentArgs) (*SecurityRoleAssignment, error) {
	if args == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ScopeId"}
	}

	assignments, err := client.ListSecurityRoleAssignments(ctx, &ListSecurityRoleAssignmentsArgs{
		Scope:      args.Scope,
		ResourceId: args.ResourceId,
	})
	if err != nil {
		return nil, err
	}

	var responseValue SecurityRoleAssignment
	for _, assignment := range *assignments {
		if *assignment.Identity.ID == args.IdentityId.String() {
			responseValue = assignment
		}
	}

	return &responseValue, err
}

// Arguments for the GetSecurityRoleAssignmentArgs function
type SetSecurityRoleAssignmentArgs struct {
	// (required) Scope identifier.
	Scope *string
	ResourceId *string
	IdentityId *uuid.UUID
	RoleName *string
}

func (client *ClientImpl) SetSecurityRoleAssignment(ctx context.Context, args *SetSecurityRoleAssignmentArgs) error {
	if args == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.ScopeId"}
	}
	routeValues := make(map[string]string)
	resId := args.ResourceId

	routeValues["scopeId"] = *args.Scope
	routeValues["resource"] = "roleassignments"
	routeValues["resourceId"] = *resId

	bodyParams := []SetRoleAssignmentPayload{}
	bodyParams = append(bodyParams, SetRoleAssignmentPayload{
		UserID: args.IdentityId,
		RoleName: args.RoleName,
	})


	body, marshalErr := json.Marshal(bodyParams)
	if marshalErr != nil {
		return marshalErr
	}

	locationId, _ := uuid.Parse("9461c234-c84c-4ed2-b918-2f0f92ad0a35")
	_, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "", nil)
	if err != nil {
		return err
	}

	return nil
}

type SetRoleAssignmentPayload struct {
	UserID *uuid.UUID `json:"userId"`
	RoleName *string `json:"roleName"`
}

type SecurityRoleIdentity struct {
	DisplayName *string `json:"displayName"`
	ID          *string `json:"id"`
	UniqueName  *string `json:"uniqueName"`
}

type SecurityRoleDefinition struct {
	DisplayName      *string `json:"displayName"`
	Name             *string `json:"name"`
	AllowPermissions *int `json:"allowPermissions"`
	DenyPermissions  *int `json:"denyPermissions"`
	Identifier       *string `json:"identifier"`
	Description      *string `json:"description"`
	Scope            *string `json:"scope"`
}

type SecurityRoleAssignment struct {
	Identity          *SecurityRoleIdentity `json:"identity"`
	Role              *SecurityRoleDefinition `json:"role"`
	Access            *string `json:"access"`
	AccessDisplayName *string `json:"accessDisplayName"`
}

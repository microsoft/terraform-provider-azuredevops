// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package pipelinepermissions

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"net/http"
)

var ResourceAreaId, _ = uuid.Parse("a81a0441-de52-4000-aa15-ff0e07bfbbaa")

type Client interface {
	// [Preview API] Given a ResourceType and ResourceId, returns authorized definitions for that resource.
	GetPipelinePermissionsForResource(context.Context, GetPipelinePermissionsForResourceArgs) (*ResourcePipelinePermissions, error)
	// [Preview API] Authorizes/Unauthorizes a list of definitions for a given resource.
	UpdatePipelinePermisionsForResource(context.Context, UpdatePipelinePermisionsForResourceArgs) (*ResourcePipelinePermissions, error)
	// [Preview API] Batch API to authorize/unauthorize a list of definitions for a multiple resources.
	UpdatePipelinePermisionsForResources(context.Context, UpdatePipelinePermisionsForResourcesArgs) (*[]ResourcePipelinePermissions, error)
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

// [Preview API] Given a ResourceType and ResourceId, returns authorized definitions for that resource.
func (client *ClientImpl) GetPipelinePermissionsForResource(ctx context.Context, args GetPipelinePermissionsForResourceArgs) (*ResourcePipelinePermissions, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.ResourceType == nil || *args.ResourceType == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ResourceType"}
	}
	routeValues["resourceType"] = *args.ResourceType
	if args.ResourceId == nil || *args.ResourceId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ResourceId"}
	}
	routeValues["resourceId"] = *args.ResourceId

	locationId, _ := uuid.Parse("b5b9a4a4-e6cd-4096-853c-ab7d8b0c4eb2")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ResourcePipelinePermissions
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetPipelinePermissionsForResource function
type GetPipelinePermissionsForResourceArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required)
	ResourceType *string
	// (required)
	ResourceId *string
}

// [Preview API] Authorizes/Unauthorizes a list of definitions for a given resource.
func (client *ClientImpl) UpdatePipelinePermisionsForResource(ctx context.Context, args UpdatePipelinePermisionsForResourceArgs) (*ResourcePipelinePermissions, error) {
	if args.ResourceAuthorization == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ResourceAuthorization"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.ResourceType == nil || *args.ResourceType == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ResourceType"}
	}
	routeValues["resourceType"] = *args.ResourceType
	if args.ResourceId == nil || *args.ResourceId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ResourceId"}
	}
	routeValues["resourceId"] = *args.ResourceId

	body, marshalErr := json.Marshal(*args.ResourceAuthorization)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("b5b9a4a4-e6cd-4096-853c-ab7d8b0c4eb2")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ResourcePipelinePermissions
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdatePipelinePermisionsForResource function
type UpdatePipelinePermisionsForResourceArgs struct {
	// (required)
	ResourceAuthorization *ResourcePipelinePermissions
	// (required) Project ID or project name
	Project *string
	// (required)
	ResourceType *string
	// (required)
	ResourceId *string
}

// [Preview API] Batch API to authorize/unauthorize a list of definitions for a multiple resources.
func (client *ClientImpl) UpdatePipelinePermisionsForResources(ctx context.Context, args UpdatePipelinePermisionsForResourcesArgs) (*[]ResourcePipelinePermissions, error) {
	if args.ResourceAuthorizations == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ResourceAuthorizations"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	body, marshalErr := json.Marshal(*args.ResourceAuthorizations)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("b5b9a4a4-e6cd-4096-853c-ab7d8b0c4eb2")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []ResourcePipelinePermissions
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdatePipelinePermisionsForResources function
type UpdatePipelinePermisionsForResourcesArgs struct {
	// (required)
	ResourceAuthorizations *[]ResourcePipelinePermissions
	// (required) Project ID or project name
	Project *string
}

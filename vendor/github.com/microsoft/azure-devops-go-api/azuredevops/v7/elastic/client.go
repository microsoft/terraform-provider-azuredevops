// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package elastic

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
	// [Preview API] Create a new elastic pool. This will create a new TaskAgentPool at the organization level. If a project id is provided, this will create a new TaskAgentQueue in the specified project.
	CreateElasticPool(context.Context, CreateElasticPoolArgs) (*ElasticPoolCreationResult, error)
	// [Preview API] Get a list of ElasticNodes currently in the ElasticPool
	GetElasticNodes(context.Context, GetElasticNodesArgs) (*[]ElasticNode, error)
	// [Preview API] Returns the Elastic Pool with the specified Pool Id.
	GetElasticPool(context.Context, GetElasticPoolArgs) (*ElasticPool, error)
	// [Preview API] Get elastic pool diagnostics logs for a specified Elastic Pool.
	GetElasticPoolLogs(context.Context, GetElasticPoolLogsArgs) (*[]ElasticPoolLog, error)
	// [Preview API] Get a list of all Elastic Pools.
	GetElasticPools(context.Context, GetElasticPoolsArgs) (*[]ElasticPool, error)
	// [Preview API] Update properties on a specified ElasticNode
	UpdateElasticNode(context.Context, UpdateElasticNodeArgs) (*ElasticNode, error)
	// [Preview API] Update settings on a specified Elastic Pool.
	UpdateElasticPool(context.Context, UpdateElasticPoolArgs) (*ElasticPool, error)
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

// [Preview API] Create a new elastic pool. This will create a new TaskAgentPool at the organization level. If a project id is provided, this will create a new TaskAgentQueue in the specified project.
func (client *ClientImpl) CreateElasticPool(ctx context.Context, args CreateElasticPoolArgs) (*ElasticPoolCreationResult, error) {
	if args.ElasticPool == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ElasticPool"}
	}
	queryParams := url.Values{}
	if args.PoolName == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "poolName"}
	}
	queryParams.Add("poolName", *args.PoolName)
	if args.AuthorizeAllPipelines != nil {
		queryParams.Add("authorizeAllPipelines", strconv.FormatBool(*args.AuthorizeAllPipelines))
	}
	if args.AutoProvisionProjectPools != nil {
		queryParams.Add("autoProvisionProjectPools", strconv.FormatBool(*args.AutoProvisionProjectPools))
	}
	if args.ProjectId != nil {
		queryParams.Add("projectId", (*args.ProjectId).String())
	}
	body, marshalErr := json.Marshal(*args.ElasticPool)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("dd3c938f-835b-4971-b99a-db75a47aad43")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ElasticPoolCreationResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateElasticPool function
type CreateElasticPoolArgs struct {
	// (required) Elastic pool to create. Contains the properties necessary for configuring a new ElasticPool.
	ElasticPool *ElasticPool
	// (required) Name to use for the new TaskAgentPool
	PoolName *string
	// (optional) Setting to determine if all pipelines are authorized to use this TaskAgentPool by default.
	AuthorizeAllPipelines *bool
	// (optional) Setting to automatically provision TaskAgentQueues in every project for the new pool.
	AutoProvisionProjectPools *bool
	// (optional) Optional: If provided, a new TaskAgentQueue will be created in the specified project.
	ProjectId *uuid.UUID
}

// [Preview API] Get a list of ElasticNodes currently in the ElasticPool
func (client *ClientImpl) GetElasticNodes(ctx context.Context, args GetElasticNodesArgs) (*[]ElasticNode, error) {
	routeValues := make(map[string]string)
	if args.PoolId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PoolId"}
	}
	routeValues["poolId"] = strconv.Itoa(*args.PoolId)

	queryParams := url.Values{}
	if args.State != nil {
		queryParams.Add("$state", string(*args.State))
	}
	locationId, _ := uuid.Parse("1b232402-5ff0-42ad-9703-d76497835eb6")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []ElasticNode
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetElasticNodes function
type GetElasticNodesArgs struct {
	// (required) Pool id of the ElasticPool
	PoolId *int
	// (optional) Optional: Filter to only retrieve ElasticNodes in the given ElasticNodeState
	State *ElasticNodeState
}

// [Preview API] Returns the Elastic Pool with the specified Pool Id.
func (client *ClientImpl) GetElasticPool(ctx context.Context, args GetElasticPoolArgs) (*ElasticPool, error) {
	routeValues := make(map[string]string)
	if args.PoolId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PoolId"}
	}
	routeValues["poolId"] = strconv.Itoa(*args.PoolId)

	locationId, _ := uuid.Parse("dd3c938f-835b-4971-b99a-db75a47aad43")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ElasticPool
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetElasticPool function
type GetElasticPoolArgs struct {
	// (required) Pool Id of the associated TaskAgentPool
	PoolId *int
}

// [Preview API] Get elastic pool diagnostics logs for a specified Elastic Pool.
func (client *ClientImpl) GetElasticPoolLogs(ctx context.Context, args GetElasticPoolLogsArgs) (*[]ElasticPoolLog, error) {
	routeValues := make(map[string]string)
	if args.PoolId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PoolId"}
	}
	routeValues["poolId"] = strconv.Itoa(*args.PoolId)

	queryParams := url.Values{}
	if args.Top != nil {
		queryParams.Add("$top", strconv.Itoa(*args.Top))
	}
	locationId, _ := uuid.Parse("595b1769-61d5-4076-a72a-98a02105ca9a")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []ElasticPoolLog
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetElasticPoolLogs function
type GetElasticPoolLogsArgs struct {
	// (required) Pool Id of the Elastic Pool
	PoolId *int
	// (optional) Number of elastic pool logs to retrieve
	Top *int
}

// [Preview API] Get a list of all Elastic Pools.
func (client *ClientImpl) GetElasticPools(ctx context.Context, args GetElasticPoolsArgs) (*[]ElasticPool, error) {
	locationId, _ := uuid.Parse("dd3c938f-835b-4971-b99a-db75a47aad43")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", nil, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []ElasticPool
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetElasticPools function
type GetElasticPoolsArgs struct {
}

// [Preview API] Update properties on a specified ElasticNode
func (client *ClientImpl) UpdateElasticNode(ctx context.Context, args UpdateElasticNodeArgs) (*ElasticNode, error) {
	if args.ElasticNodeSettings == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ElasticNodeSettings"}
	}
	routeValues := make(map[string]string)
	if args.PoolId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PoolId"}
	}
	routeValues["poolId"] = strconv.Itoa(*args.PoolId)
	if args.ElasticNodeId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ElasticNodeId"}
	}
	routeValues["elasticNodeId"] = strconv.Itoa(*args.ElasticNodeId)

	body, marshalErr := json.Marshal(*args.ElasticNodeSettings)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("1b232402-5ff0-42ad-9703-d76497835eb6")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ElasticNode
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateElasticNode function
type UpdateElasticNodeArgs struct {
	// (required)
	ElasticNodeSettings *ElasticNodeSettings
	// (required)
	PoolId *int
	// (required)
	ElasticNodeId *int
}

// [Preview API] Update settings on a specified Elastic Pool.
func (client *ClientImpl) UpdateElasticPool(ctx context.Context, args UpdateElasticPoolArgs) (*ElasticPool, error) {
	if args.ElasticPoolSettings == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ElasticPoolSettings"}
	}
	routeValues := make(map[string]string)
	if args.PoolId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PoolId"}
	}
	routeValues["poolId"] = strconv.Itoa(*args.PoolId)

	body, marshalErr := json.Marshal(*args.ElasticPoolSettings)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("dd3c938f-835b-4971-b99a-db75a47aad43")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ElasticPool
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateElasticPool function
type UpdateElasticPoolArgs struct {
	// (required) New Elastic Pool settings data
	ElasticPoolSettings *ElasticPoolSettings
	// (required)
	PoolId *int
}

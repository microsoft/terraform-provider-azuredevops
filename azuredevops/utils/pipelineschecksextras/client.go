// This is a copy of github.com/microsoft/azure-devops-go-api/azuredevops/pipelineschecks/client.go
// The existing version does not contain the "Timeout" property on the CheckConfiguration struct

// This file cannot be under "internal", because azdosdkmocks/pipelines_checks_v5_extras_mock.go depends on it.

package pipelineschecksextras

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
)

var ResourceAreaId, _ = uuid.Parse("4a933897-0488-45af-bd82-6fd3ad33f46a")

type Client interface {
	// [Preview API] Add a check configuration
	AddCheckConfiguration(context.Context, AddCheckConfigurationArgs) (*CheckConfiguration, error)
	// [Preview API] Delete check configuration by id
	DeleteCheckConfiguration(context.Context, DeleteCheckConfigurationArgs) error
	// [Preview API] Initiate an evaluation for a check in a pipeline
	EvaluateCheckSuite(context.Context, EvaluateCheckSuiteArgs) (*CheckSuite, error)
	// [Preview API] Get Check configuration by Id
	GetCheckConfiguration(context.Context, GetCheckConfigurationArgs) (*CheckConfiguration, error)
	// [Preview API] Get Check configuration by resource type and id
	GetCheckConfigurationsOnResource(context.Context, GetCheckConfigurationsOnResourceArgs) (*[]CheckConfiguration, error)
	// [Preview API] Get details for a specific check evaluation
	GetCheckSuite(context.Context, GetCheckSuiteArgs) (*CheckSuite, error)
	// [Preview API] Get check configurations for multiple resources by resource type and id.
	QueryCheckConfigurationsOnResources(context.Context, QueryCheckConfigurationsOnResourcesArgs) (*[]CheckConfiguration, error)
	// [Preview API] Update check configuration
	UpdateCheckConfiguration(context.Context, UpdateCheckConfigurationArgs) (*CheckConfiguration, error)
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

// [Preview API] Add a check configuration
func (client *ClientImpl) AddCheckConfiguration(ctx context.Context, args AddCheckConfigurationArgs) (*CheckConfiguration, error) {
	if args.Configuration == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Configuration"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	body, marshalErr := json.Marshal(*args.Configuration)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("86c8381e-5aee-4cde-8ae4-25c0c7f5eaea")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue CheckConfiguration
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the AddCheckConfiguration function
type AddCheckConfigurationArgs struct {
	// (required)
	Configuration *CheckConfiguration
	// (required) Project ID or project name
	Project *string
}

// [Preview API]
func (client *ClientImpl) DeleteCheckConfiguration(ctx context.Context, args DeleteCheckConfigurationArgs) error {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Id == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	locationId, _ := uuid.Parse("86c8381e-5aee-4cde-8ae4-25c0c7f5eaea")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteCheckConfiguration function
type DeleteCheckConfigurationArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) check configuration id
	Id *int
}

// [Preview API]
func (client *ClientImpl) EvaluateCheckSuite(ctx context.Context, args EvaluateCheckSuiteArgs) (*CheckSuite, error) {
	if args.Request == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Request"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	queryParams := url.Values{}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	body, marshalErr := json.Marshal(*args.Request)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("91282c1d-c183-444f-9554-1485bfb3879d")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue CheckSuite
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the EvaluateCheckSuite function
type EvaluateCheckSuiteArgs struct {
	// (required)
	Request *CheckSuiteRequest
	// (required) Project ID or project name
	Project *string
	// (optional)
	Expand *CheckSuiteExpandParameter
}

// [Preview API] Get Check configuration by Id
func (client *ClientImpl) GetCheckConfiguration(ctx context.Context, args GetCheckConfigurationArgs) (*CheckConfiguration, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	queryParams := url.Values{}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	locationId, _ := uuid.Parse("86c8381e-5aee-4cde-8ae4-25c0c7f5eaea")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue CheckConfiguration
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetCheckConfiguration function
type GetCheckConfigurationArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required)
	Id *int
	// (optional)
	Expand *CheckConfigurationExpandParameter
}

// [Preview API] Get Check configuration by resource type and id
func (client *ClientImpl) GetCheckConfigurationsOnResource(ctx context.Context, args GetCheckConfigurationsOnResourceArgs) (*[]CheckConfiguration, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	queryParams := url.Values{}
	if args.ResourceType != nil {
		queryParams.Add("resourceType", *args.ResourceType)
	}
	if args.ResourceId != nil {
		queryParams.Add("resourceId", *args.ResourceId)
	}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	locationId, _ := uuid.Parse("86c8381e-5aee-4cde-8ae4-25c0c7f5eaea")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []CheckConfiguration
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetCheckConfigurationsOnResource function
type GetCheckConfigurationsOnResourceArgs struct {
	// (required) Project ID or project name
	Project *string
	// (optional) resource type
	ResourceType *string
	// (optional) resource id
	ResourceId *string
	// (optional)
	Expand *CheckConfigurationExpandParameter
}

// [Preview API]
func (client *ClientImpl) GetCheckSuite(ctx context.Context, args GetCheckSuiteArgs) (*CheckSuite, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.CheckSuiteId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.CheckSuiteId"}
	}
	routeValues["checkSuiteId"] = (*args.CheckSuiteId).String()

	queryParams := url.Values{}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	locationId, _ := uuid.Parse("91282c1d-c183-444f-9554-1485bfb3879d")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue CheckSuite
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetCheckSuite function
type GetCheckSuiteArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required)
	CheckSuiteId *uuid.UUID
	// (optional)
	Expand *CheckSuiteExpandParameter
}

// [Preview API] Update check configuration
func (client *ClientImpl) UpdateCheckConfiguration(ctx context.Context, args UpdateCheckConfigurationArgs) (*CheckConfiguration, error) {
	if args.Configuration == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Configuration"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	body, marshalErr := json.Marshal(*args.Configuration)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("86c8381e-5aee-4cde-8ae4-25c0c7f5eaea")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue CheckConfiguration
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateCheckConfiguration function
type UpdateCheckConfigurationArgs struct {
	// (required) check configuration
	Configuration *CheckConfiguration
	// (required) Project ID or project name
	Project *string
	// (required) check configuration id
	Id *int
}

// [Preview API] Get check configurations for multiple resources by resource type and id.
func (client *ClientImpl) QueryCheckConfigurationsOnResources(ctx context.Context, args QueryCheckConfigurationsOnResourcesArgs) (*[]CheckConfiguration, error) {
	if args.Resources == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Resources"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	queryParams := url.Values{}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	body, marshalErr := json.Marshal(*args.Resources)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("5f3d0e64-f943-4584-8811-77eb495e831e")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []CheckConfiguration
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryCheckConfigurationsOnResources function
type QueryCheckConfigurationsOnResourcesArgs struct {
	// (required) List of resources.
	Resources *[]Resource
	// (required) Project ID or project name
	Project *string
	// (optional) The properties that should be expanded in the list of check configurations.
	Expand *CheckConfigurationExpandParameter
}

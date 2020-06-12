// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package featuremanagement

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"net/http"
	"net/url"
)

type Client interface {
	// [Preview API] Get a specific feature by its id
	GetFeature(context.Context, GetFeatureArgs) (*ContributedFeature, error)
	// [Preview API] Get a list of all defined features
	GetFeatures(context.Context, GetFeaturesArgs) (*[]ContributedFeature, error)
	// [Preview API] Get the state of the specified feature for the given user/all-users scope
	GetFeatureState(context.Context, GetFeatureStateArgs) (*ContributedFeatureState, error)
	// [Preview API] Get the state of the specified feature for the given named scope
	GetFeatureStateForScope(context.Context, GetFeatureStateForScopeArgs) (*ContributedFeatureState, error)
	// [Preview API] Get the effective state for a list of feature ids
	QueryFeatureStates(context.Context, QueryFeatureStatesArgs) (*ContributedFeatureStateQuery, error)
	// [Preview API] Get the states of the specified features for the default scope
	QueryFeatureStatesForDefaultScope(context.Context, QueryFeatureStatesForDefaultScopeArgs) (*ContributedFeatureStateQuery, error)
	// [Preview API] Get the states of the specified features for the specific named scope
	QueryFeatureStatesForNamedScope(context.Context, QueryFeatureStatesForNamedScopeArgs) (*ContributedFeatureStateQuery, error)
	// [Preview API] Set the state of a feature
	SetFeatureState(context.Context, SetFeatureStateArgs) (*ContributedFeatureState, error)
	// [Preview API] Set the state of a feature at a specific scope
	SetFeatureStateForScope(context.Context, SetFeatureStateForScopeArgs) (*ContributedFeatureState, error)
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

// [Preview API] Get a specific feature by its id
func (client *ClientImpl) GetFeature(ctx context.Context, args GetFeatureArgs) (*ContributedFeature, error) {
	routeValues := make(map[string]string)
	if args.FeatureId == nil || *args.FeatureId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeatureId"}
	}
	routeValues["featureId"] = *args.FeatureId

	locationId, _ := uuid.Parse("c4209f25-7a27-41dd-9f04-06080c7b6afd")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ContributedFeature
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetFeature function
type GetFeatureArgs struct {
	// (required) The contribution id of the feature
	FeatureId *string
}

// [Preview API] Get a list of all defined features
func (client *ClientImpl) GetFeatures(ctx context.Context, args GetFeaturesArgs) (*[]ContributedFeature, error) {
	queryParams := url.Values{}
	if args.TargetContributionId != nil {
		queryParams.Add("targetContributionId", *args.TargetContributionId)
	}
	locationId, _ := uuid.Parse("c4209f25-7a27-41dd-9f04-06080c7b6afd")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1-preview.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []ContributedFeature
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetFeatures function
type GetFeaturesArgs struct {
	// (optional) Optional target contribution. If null/empty, return all features. If specified include the features that target the specified contribution.
	TargetContributionId *string
}

// [Preview API] Get the state of the specified feature for the given user/all-users scope
func (client *ClientImpl) GetFeatureState(ctx context.Context, args GetFeatureStateArgs) (*ContributedFeatureState, error) {
	routeValues := make(map[string]string)
	if args.FeatureId == nil || *args.FeatureId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeatureId"}
	}
	routeValues["featureId"] = *args.FeatureId
	if args.UserScope == nil || *args.UserScope == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.UserScope"}
	}
	routeValues["userScope"] = *args.UserScope

	locationId, _ := uuid.Parse("98911314-3f9b-4eaf-80e8-83900d8e85d9")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ContributedFeatureState
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetFeatureState function
type GetFeatureStateArgs struct {
	// (required) Contribution id of the feature
	FeatureId *string
	// (required) User-Scope at which to get the value. Should be "me" for the current user or "host" for all users.
	UserScope *string
}

// [Preview API] Get the state of the specified feature for the given named scope
func (client *ClientImpl) GetFeatureStateForScope(ctx context.Context, args GetFeatureStateForScopeArgs) (*ContributedFeatureState, error) {
	routeValues := make(map[string]string)
	if args.FeatureId == nil || *args.FeatureId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeatureId"}
	}
	routeValues["featureId"] = *args.FeatureId
	if args.UserScope == nil || *args.UserScope == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.UserScope"}
	}
	routeValues["userScope"] = *args.UserScope
	if args.ScopeName == nil || *args.ScopeName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ScopeName"}
	}
	routeValues["scopeName"] = *args.ScopeName
	if args.ScopeValue == nil || *args.ScopeValue == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ScopeValue"}
	}
	routeValues["scopeValue"] = *args.ScopeValue

	locationId, _ := uuid.Parse("dd291e43-aa9f-4cee-8465-a93c78e414a4")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ContributedFeatureState
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetFeatureStateForScope function
type GetFeatureStateForScopeArgs struct {
	// (required) Contribution id of the feature
	FeatureId *string
	// (required) User-Scope at which to get the value. Should be "me" for the current user or "host" for all users.
	UserScope *string
	// (required) Scope at which to get the feature setting for (e.g. "project" or "team")
	ScopeName *string
	// (required) Value of the scope (e.g. the project or team id)
	ScopeValue *string
}

// [Preview API] Get the effective state for a list of feature ids
func (client *ClientImpl) QueryFeatureStates(ctx context.Context, args QueryFeatureStatesArgs) (*ContributedFeatureStateQuery, error) {
	if args.Query == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Query"}
	}
	body, marshalErr := json.Marshal(*args.Query)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("2b4486ad-122b-400c-ae65-17b6672c1f9d")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "5.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ContributedFeatureStateQuery
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryFeatureStates function
type QueryFeatureStatesArgs struct {
	// (required) Features to query along with current scope values
	Query *ContributedFeatureStateQuery
}

// [Preview API] Get the states of the specified features for the default scope
func (client *ClientImpl) QueryFeatureStatesForDefaultScope(ctx context.Context, args QueryFeatureStatesForDefaultScopeArgs) (*ContributedFeatureStateQuery, error) {
	if args.Query == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Query"}
	}
	routeValues := make(map[string]string)
	if args.UserScope == nil || *args.UserScope == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.UserScope"}
	}
	routeValues["userScope"] = *args.UserScope

	body, marshalErr := json.Marshal(*args.Query)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("3f810f28-03e2-4239-b0bc-788add3005e5")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "5.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ContributedFeatureStateQuery
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryFeatureStatesForDefaultScope function
type QueryFeatureStatesForDefaultScopeArgs struct {
	// (required) Query describing the features to query.
	Query *ContributedFeatureStateQuery
	// (required)
	UserScope *string
}

// [Preview API] Get the states of the specified features for the specific named scope
func (client *ClientImpl) QueryFeatureStatesForNamedScope(ctx context.Context, args QueryFeatureStatesForNamedScopeArgs) (*ContributedFeatureStateQuery, error) {
	if args.Query == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Query"}
	}
	routeValues := make(map[string]string)
	if args.UserScope == nil || *args.UserScope == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.UserScope"}
	}
	routeValues["userScope"] = *args.UserScope
	if args.ScopeName == nil || *args.ScopeName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ScopeName"}
	}
	routeValues["scopeName"] = *args.ScopeName
	if args.ScopeValue == nil || *args.ScopeValue == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ScopeValue"}
	}
	routeValues["scopeValue"] = *args.ScopeValue

	body, marshalErr := json.Marshal(*args.Query)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("f29e997b-c2da-4d15-8380-765788a1a74c")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "5.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ContributedFeatureStateQuery
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryFeatureStatesForNamedScope function
type QueryFeatureStatesForNamedScopeArgs struct {
	// (required) Query describing the features to query.
	Query *ContributedFeatureStateQuery
	// (required)
	UserScope *string
	// (required)
	ScopeName *string
	// (required)
	ScopeValue *string
}

// [Preview API] Set the state of a feature
func (client *ClientImpl) SetFeatureState(ctx context.Context, args SetFeatureStateArgs) (*ContributedFeatureState, error) {
	if args.Feature == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Feature"}
	}
	routeValues := make(map[string]string)
	if args.FeatureId == nil || *args.FeatureId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeatureId"}
	}
	routeValues["featureId"] = *args.FeatureId
	if args.UserScope == nil || *args.UserScope == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.UserScope"}
	}
	routeValues["userScope"] = *args.UserScope

	queryParams := url.Values{}
	if args.Reason != nil {
		queryParams.Add("reason", *args.Reason)
	}
	if args.ReasonCode != nil {
		queryParams.Add("reasonCode", *args.ReasonCode)
	}
	body, marshalErr := json.Marshal(*args.Feature)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("98911314-3f9b-4eaf-80e8-83900d8e85d9")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "5.1-preview.1", routeValues, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ContributedFeatureState
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the SetFeatureState function
type SetFeatureStateArgs struct {
	// (required) Posted feature state object. Should specify the effective value.
	Feature *ContributedFeatureState
	// (required) Contribution id of the feature
	FeatureId *string
	// (required) User-Scope at which to set the value. Should be "me" for the current user or "host" for all users.
	UserScope *string
	// (optional) Reason for changing the state
	Reason *string
	// (optional) Short reason code
	ReasonCode *string
}

// [Preview API] Set the state of a feature at a specific scope
func (client *ClientImpl) SetFeatureStateForScope(ctx context.Context, args SetFeatureStateForScopeArgs) (*ContributedFeatureState, error) {
	if args.Feature == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Feature"}
	}
	routeValues := make(map[string]string)
	if args.FeatureId == nil || *args.FeatureId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeatureId"}
	}
	routeValues["featureId"] = *args.FeatureId
	if args.UserScope == nil || *args.UserScope == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.UserScope"}
	}
	routeValues["userScope"] = *args.UserScope
	if args.ScopeName == nil || *args.ScopeName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ScopeName"}
	}
	routeValues["scopeName"] = *args.ScopeName
	if args.ScopeValue == nil || *args.ScopeValue == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ScopeValue"}
	}
	routeValues["scopeValue"] = *args.ScopeValue

	queryParams := url.Values{}
	if args.Reason != nil {
		queryParams.Add("reason", *args.Reason)
	}
	if args.ReasonCode != nil {
		queryParams.Add("reasonCode", *args.ReasonCode)
	}
	body, marshalErr := json.Marshal(*args.Feature)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("dd291e43-aa9f-4cee-8465-a93c78e414a4")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "5.1-preview.1", routeValues, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ContributedFeatureState
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the SetFeatureStateForScope function
type SetFeatureStateForScopeArgs struct {
	// (required) Posted feature state object. Should specify the effective value.
	Feature *ContributedFeatureState
	// (required) Contribution id of the feature
	FeatureId *string
	// (required) User-Scope at which to set the value. Should be "me" for the current user or "host" for all users.
	UserScope *string
	// (required) Scope at which to get the feature setting for (e.g. "project" or "team")
	ScopeName *string
	// (required) Value of the scope (e.g. the project or team id)
	ScopeValue *string
	// (optional) Reason for changing the state
	Reason *string
	// (optional) Short reason code
	ReasonCode *string
}
